package cmd

import (
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRenderPBar is a test function for RenderPBar()
func TestRenderPBar(t *testing.T) {
	tests := []struct {
		testName string
		pBar     *PBar
	}{
		{
			testName: "Test Render Progress Bar",
			pBar: &PBar{ // Since no TTY is available, we can't test the output nor the window change
				Total:       uint16(rand.Intn(100)),
				header:      uint16(rand.Intn(100)),
				wscol:       uint16(rand.Intn(100)),
				wsrow:       uint16(rand.Intn(100)),
				doneStr:     "#",
				ongoingStr:  ".",
				signalWinch: make(chan os.Signal, 1),
				signalTerm:  make(chan os.Signal, 1),
			},
		},
	}

	for _, test := range tests {
		test.pBar.SignalHandler() // Handle the signals

		for count := 0; uint16(count) <= test.pBar.Total; count++ {
			assert.NotPanics(t, func() { test.pBar.RenderPBar(count) }, test.testName)
		}

		test.pBar.CleanUp() // Restore reserved bottom line
	}
}
