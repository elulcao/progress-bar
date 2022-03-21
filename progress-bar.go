package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

var (
	total      = 100 // Total number of iterations to sum 100%
	header     = 0   // Header length, to be used to calculate the bar width "Progress: [100%] []"
	count      = 0   // Current iteration
	wscol      = 0   // Window width
	wsrow      = 0   // Window height
	doneStr    = "#" // Progress bar done string
	ongoingStr = "." // Progress bar ongoing string
)

// init do the task we want to do before doing all other things
func init() {
	if err := updateWSize(); err != nil {
		panic(err)
	}
}

// Restore reserved bottom line
func cleanUp() {
	fmt.Print("\x1B7")              // Save the cursor position
	fmt.Printf("\x1B[0;%dr", wsrow) // Drop margin reservation
	fmt.Printf("\x1B[%d;0f", wsrow) // Move the cursor to the bottom line
	fmt.Print("\x1B[0K")            // Erase the entire line
	fmt.Print("\x1B8")              // Restore the cursor position
}

// updateWSize update the window size
func updateWSize() error {
	fmt.Printf("\x1B[0;%dr", wsrow) // Drop margin reservation

	ws, err := unix.IoctlGetWinsize(syscall.Stdout, unix.TIOCGWINSZ)
	if err != nil {
		return err
	}

	wscol = int(ws.Col)
	wsrow = int(ws.Row)

	switch {
	case wscol >= 0 && wscol <= 9:
		header = 6 // len("[100%]")
	case wscol >= 10 && wscol <= 20:
		header = 9 // len("[100%] []")
	default:
		header = 19 // len("Progress: [100%] []")
	}

	fmt.Print("\x1BD")                // Return carriage
	fmt.Print("\x1B7")                // Save the cursor position
	fmt.Printf("\x1B[0;%dr", wsrow-1) // Reserve the bottom line
	fmt.Print("\x1B8")                // Restore the cursor position
	fmt.Print("\x1B[1A")              // Moves cursor up # lines

	return nil
}

// renderPBar render the progress bar
func renderPBar() {
	fmt.Print("\x1B7")       // Save the cursor position
	fmt.Print("\x1B[2K")     // Erase the entire line
	fmt.Print("\x1B[0J")     // Erase from cursor to end of screen
	fmt.Print("\x1B[?47h")   // Save screen
	fmt.Print("\x1B[1J")     // Erase from cursor to beginning of screen
	fmt.Print("\x1B[?47l")   // Restore screen
	defer fmt.Print("\x1B8") // Restore the cursor position

	barWidth := int(math.Abs(float64(wscol - header)))                  // Calculate the bar width
	barDone := int(float64(barWidth) * float64(count) / float64(total)) // Calculate the bar done length
	done := strings.Repeat(doneStr, barDone)                            // Fill the bar with done string
	todo := strings.Repeat(ongoingStr, barWidth-barDone)                // Fill the bar with todo string
	bar := fmt.Sprintf("[%s%s]", done, todo)                            // Combine the done and todo string

	fmt.Printf("\x1B[%d;%dH", wsrow, 0) // move cursor to row #, col #

	switch {
	case wscol >= 0 && wscol <= 9:
		fmt.Printf("[\x1B[33m%3d%%\x1B[0m]", count*100/total)
	case wscol >= 10 && wscol <= 20:
		fmt.Printf("[\x1B[33m%3d%%\x1B[0m] %s", count*100/total, bar)
	default:
		fmt.Printf("Progress: [\x1B[33m%3d%%\x1B[0m] %s", count*100/total, bar)
	}
}

// main do the task we want to do
func main() {
	sigwinch := make(chan os.Signal, 1) // Set signal handler
	defer close(sigwinch)
	signal.Notify(sigwinch, syscall.SIGWINCH)
	go func() {
		for {
			if _, ok := <-sigwinch; !ok {
				return
			}

			err := updateWSize()
			if err != nil {
				panic(err) // The window size could not be updated
			}
		}
	}()

	sigterm := make(chan os.Signal, 1)
	defer close(sigterm)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			if _, ok := <-sigterm; !ok {
				return
			}

			cleanUp() // Restore reserved bottom line
			os.Exit(0)
		}
	}()

	for count = 1; count <= total; count++ {
		renderPBar()
		time.Sleep(time.Second)
		fmt.Println(count) // Action to be performed after 1 second
	}

	cleanUp() // Restore reserved bottom line
}
