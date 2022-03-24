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
		pBar     *pBar
	}{
		{
			testName: "Test Render Progress Bar",
			pBar: &pBar{ // Since no TTY is available, we can't test the output nor the window change
				Total:      uint16(rand.Intn(100)),
				Header:     uint16(rand.Intn(100)),
				Wscol:      uint16(rand.Intn(100)),
				Wsrow:      uint16(rand.Intn(100)),
				DoneStr:    "#",
				OngoingStr: ".",
				Sigwinch:   make(chan os.Signal, 1),
				Sigterm:    make(chan os.Signal, 1),
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
