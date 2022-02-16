package internal

import (
	"sync"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/internal/app"
	"github.com/bhojpur/gui/pkg/engine/internal/widget"
)

// OverlayStack implements gui.OverlayStack
type OverlayStack struct {
	OnChange      func()
	Canvas        gui.Canvas
	focusManagers []*app.FocusManager
	overlays      []gui.CanvasObject
	propertyLock  sync.RWMutex
}

var _ gui.OverlayStack = (*OverlayStack)(nil)

// Add puts an overlay on the stack.
//
// Implements: gui.OverlayStack
func (s *OverlayStack) Add(overlay gui.CanvasObject) {
	if overlay == nil {
		return
	}

	if s.OnChange != nil {
		defer s.OnChange()
	}

	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()
	s.overlays = append(s.overlays, overlay)

	// TODO this should probably apply to all once #707 is addressed
	if _, ok := overlay.(*widget.OverlayContainer); ok {
		safePos, safeSize := s.Canvas.InteractiveArea()

		overlay.Resize(safeSize)
		overlay.Move(safePos)
	}

	s.focusManagers = append(s.focusManagers, app.NewFocusManager(overlay))
}

// List returns all overlays on the stack from bottom to top.
//
// Implements: gui.OverlayStack
func (s *OverlayStack) List() []gui.CanvasObject {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()

	return s.overlays
}

// ListFocusManagers returns all focus managers on the stack from bottom to top.
func (s *OverlayStack) ListFocusManagers() []*app.FocusManager {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()

	return s.focusManagers
}

// Remove deletes an overlay and all overlays above it from the stack.
//
// Implements: gui.OverlayStack
func (s *OverlayStack) Remove(overlay gui.CanvasObject) {
	if s.OnChange != nil {
		defer s.OnChange()
	}

	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()

	for i, o := range s.overlays {
		if o == overlay {
			s.overlays = s.overlays[:i]
			s.focusManagers = s.focusManagers[:i]
			break
		}
	}
}

// Top returns the top-most overlay of the stack.
//
// Implements: gui.OverlayStack
func (s *OverlayStack) Top() gui.CanvasObject {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()

	if len(s.overlays) == 0 {
		return nil
	}
	return s.overlays[len(s.overlays)-1]
}

// TopFocusManager returns the app.FocusManager assigned to the top-most overlay of the stack.
func (s *OverlayStack) TopFocusManager() *app.FocusManager {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()
	return s.topFocusManager()
}

func (s *OverlayStack) topFocusManager() *app.FocusManager {
	var fm *app.FocusManager
	if len(s.focusManagers) > 0 {
		fm = s.focusManagers[len(s.focusManagers)-1]
	}
	return fm
}
