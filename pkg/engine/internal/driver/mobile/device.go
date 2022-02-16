package mobile

import (
	"github.com/bhojpur/gui/pkg/engine/driver/mobile"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/mobile/event/size"

	gui "github.com/bhojpur/gui/pkg/engine"
)

type device struct {
	safeTop, safeLeft, safeWidth, safeHeight int
}

//lint:file-ignore U1000 Var currentDPI is used in other files, but not here
var (
	currentOrientation size.Orientation
	currentDPI         float32
)

// Declare conformity with Device
var _ gui.Device = (*device)(nil)

func (*device) Orientation() gui.DeviceOrientation {
	switch currentOrientation {
	case size.OrientationLandscape:
		return gui.OrientationHorizontalLeft
	default:
		return gui.OrientationVertical
	}
}

func (*device) IsMobile() bool {
	return true
}

func (*device) HasKeyboard() bool {
	return false
}

func (*device) ShowVirtualKeyboard() {
	showVirtualKeyboard(mobile.DefaultKeyboard)
}

func (*device) ShowVirtualKeyboardType(keyboard mobile.KeyboardType) {
	showVirtualKeyboard(keyboard)
}

func (*device) HideVirtualKeyboard() {
	hideVirtualKeyboard()
}
