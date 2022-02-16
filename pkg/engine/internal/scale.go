package internal

import (
	"math"

	gui "github.com/bhojpur/gui/pkg/engine"
)

// ScaleInt converts a Bhojpur GUI coordinate in the given canvas to a screen coordinate
func ScaleInt(c gui.Canvas, v float32) int {
	return int(math.Round(float64(v * c.Scale())))
}

// UnscaleInt converts a screen coordinate for a given canvas to a Bhojpur GUI coordinate
func UnscaleInt(c gui.Canvas, v int) float32 {
	switch c.Scale() {
	case 1.0:
		return float32(v)
	default:
		return float32(v) / c.Scale()
	}
}
