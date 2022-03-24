package cmd

import (
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
			pBar:     NewPBar(),
		},
	}

	for idx, test := range tests {
		test.pBar.SignalHandler() // Handle the signals
		test.pBar.Total = uint16(idx + 1)

		for count := 0; uint16(count) <= test.pBar.Total; count++ {
			assert.NotPanics(t, func() { test.pBar.RenderPBar(count) }, test.testName)
		}

		test.pBar.CleanUp() // Restore reserved bottom line
	}
}
