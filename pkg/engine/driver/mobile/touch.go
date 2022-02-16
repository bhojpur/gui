package mobile

import gui "github.com/bhojpur/gui/pkg/engine"

// TouchEvent contains data relating to mobile touch events
type TouchEvent struct {
	gui.PointEvent
}

// Touchable represents mobile touch events that can be sent to CanvasObjects
type Touchable interface {
	TouchDown(*TouchEvent)
	TouchUp(*TouchEvent)
	TouchCancel(*TouchEvent)
}
