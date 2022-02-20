package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

var (
	header     = len("Progress: [100%] []")
	total      = 100
	count      = 0
	wscol      = 0
	wsrow      = 0
	doneStr    = "#"
	ongoingStr = "."
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

	barWidth := wscol - header
	barDone := int(float64(barWidth) * float64(count) / float64(total))

	fmt.Printf("\x1B[%d;%dH", wsrow, 0) // move cursor to row #, col #
	fmt.Printf("Progress: [\x1B[33m%3d%%\x1B[0m] ", count*100/total)
	fmt.Printf("[%s%s]", strings.Repeat(doneStr, barDone), strings.Repeat(ongoingStr, barWidth-barDone))
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

			_ = updateWSize()
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
