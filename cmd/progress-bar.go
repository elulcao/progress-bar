package cmd

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

// pBar is the progress bar model
type pBar struct {
	Total      uint16         // Total number of iterations to sum 100%
	Header     uint16         // Header length, to be used to calculate the bar width "Progress: [100%] []"
	Wscol      uint16         // Window width
	Wsrow      uint16         // Window height
	DoneStr    string         // Progress bar done string
	OngoingStr string         // Progress bar ongoing string
	Sigwinch   chan os.Signal // Signal handler: SIGWINCH
	Sigterm    chan os.Signal // Signal handler: SIGTERM
	once       sync.Once      // Close the signal channel only once
}

type winSize struct {
	Row    uint16 // row
	Col    uint16 // column
	Xpixel uint16 // X pixel
	Ypixel uint16 // Y pixel
}

// init do the task we want to do before doing all other things
func init() {}

// NewPBar create a new progress bar
// After NewPBar() is called:
// 	- initialize SignalHandler()
// 	- update pBar.Total for new number of iterations to sum 100%
// After progressBar() is finished:
//	- do a CleanUp()
func NewPBar() *pBar {
	pb := &pBar{
		Total:      100,
		Header:     0,
		Wscol:      0,
		Wsrow:      0,
		DoneStr:    "#",
		OngoingStr: ".",
		Sigwinch:   make(chan os.Signal, 1),
		Sigterm:    make(chan os.Signal, 1),
	}

	signal.Notify(pb.Sigwinch, syscall.SIGWINCH)               // Register SIGWINCH signal
	signal.Notify(pb.Sigterm, syscall.SIGINT, syscall.SIGTERM) // Register SIGINT and SIGTERM signal

	pb.UpdateWSize()

	return pb
}

// CleanUp restore reserved bottom line and restore cursor position
func (pb *pBar) CleanUp() {
	fmt.Print("\x1B7")                 // Save the cursor position
	fmt.Printf("\x1B[0;%dr", pb.Wsrow) // Drop margin reservation
	fmt.Printf("\x1B[%d;0f", pb.Wsrow) // Move the cursor to the bottom line
	fmt.Print("\x1B[0K")               // Erase the entire line
	fmt.Print("\x1B8")                 // Restore the cursor position

	pb.once.Do(func() { close(pb.Sigwinch) }) // Close the signal channel politely
	pb.once.Do(func() { close(pb.Sigterm) })  // Close the signal channel politely
}

// UpdateWSize update the window size
func (pb *pBar) UpdateWSize() error {
	fmt.Printf("\x1B[0;%dr", pb.Wsrow) // Drop margin reservation

	ws := &winSize{}
	ret, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(ws)))
	if int(ret) == -1 {
		panic(err)
	}

	pb.Wscol = ws.Col
	pb.Wsrow = ws.Row

	switch {
	case pb.Wscol >= 0 && pb.Wscol <= 9:
		pb.Header = 6 // len("[100%]") is the minimum header length
	case pb.Wscol >= 10 && pb.Wscol <= 20:
		pb.Header = 9 // len("[100%] []") is the midium header length
	default:
		pb.Header = 19 // len("Progress: [100%] []") is the maximum header length
	}

	fmt.Print("\x1BD")                   // Return carriage
	fmt.Print("\x1B7")                   // Save the cursor position
	fmt.Printf("\x1B[0;%dr", pb.Wsrow-1) // Reserve the bottom line
	fmt.Print("\x1B8")                   // Restore the cursor position
	fmt.Print("\x1B[1A")                 // Moves cursor up # lines

	return nil
}

// SignalHandler handle the signals, like SIGWINCH and SIGTERM
func (pb *pBar) SignalHandler() {
	go func() {
		for {
			select {
			case <-pb.Sigwinch:
				if err := pb.UpdateWSize(); err != nil {
					panic(err) // The window size could not be updated
				}
			case <-pb.Sigterm:
				fmt.Printf("\nCaught SIGTERM, exiting...\n") // Print the message for SIGTERM
				pb.CleanUp()                                 // Restore reserved bottom line
				os.Exit(0)                                   // Exit gracefully
			}
		}
	}()
}

// RenderPBar render the progress bar. Receives the current iteration count
func (pb *pBar) RenderPBar(count int) {
	fmt.Print("\x1B7")       // Save the cursor position
	fmt.Print("\x1B[2K")     // Erase the entire line
	fmt.Print("\x1B[0J")     // Erase from cursor to end of screen
	fmt.Print("\x1B[?47h")   // Save screen
	fmt.Print("\x1B[1J")     // Erase from cursor to beginning of screen
	fmt.Print("\x1B[?47l")   // Restore screen
	defer fmt.Print("\x1B8") // Restore the cursor position util new size is calculated

	barWidth := int(math.Abs(float64(pb.Wscol - pb.Header)))               // Calculate the bar width
	barDone := int(float64(barWidth) * float64(count) / float64(pb.Total)) // Calculate the bar done length
	done := strings.Repeat(pb.DoneStr, barDone)                            // Fill the bar with done string
	todo := strings.Repeat(pb.OngoingStr, barWidth-barDone)                // Fill the bar with todo string
	bar := fmt.Sprintf("[%s%s]", done, todo)                               // Combine the done and todo string

	fmt.Printf("\x1B[%d;%dH", pb.Wsrow, 0) // move cursor to row #, col #

	switch {
	case pb.Wscol >= 0 && pb.Wscol <= 9:
		fmt.Printf("[\x1B[33m%3d%%\x1B[0m]", uint16(count)*100/pb.Total)
	case pb.Wscol >= 10 && pb.Wscol <= 20:
		fmt.Printf("[\x1B[33m%3d%%\x1B[0m] %s", uint16(count)*100/pb.Total, bar)
	default:
		fmt.Printf("Progress: [\x1B[33m%3d%%\x1B[0m] %s", uint16(count)*100/pb.Total, bar)
	}
}
