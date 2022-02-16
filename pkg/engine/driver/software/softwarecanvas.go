package software

import (
	"github.com/bhojpur/gui/pkg/engine/internal/painter/software"
	"github.com/bhojpur/gui/pkg/engine/test"
)

// NewCanvas creates a new canvas in memory that can render without hardware support
func NewCanvas() test.WindowlessCanvas {
	return test.NewCanvasWithPainter(software.NewPainter())
}
