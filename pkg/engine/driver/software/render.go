package software

import (
	"image"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/internal/app"
)

// RenderCanvas takes a canvas and renders it to a regular Go image using the provided Theme.
// This is the same as setting the application theme and then calling Canvas.Capture().
func RenderCanvas(c gui.Canvas, t gui.Theme) image.Image {
	gui.CurrentApp().Settings().SetTheme(t)
	app.ApplyThemeTo(c.Content(), c)

	return c.Capture()
}

// Render takes a canvas object and renders it to a regular Go image using the provided Theme.
// The returned image will be set to the object's minimum size.
// Use the theme.LightTheme() or theme.DarkTheme() to access the builtin themes.
func Render(obj gui.CanvasObject, t gui.Theme) image.Image {
	c := NewCanvas()
	c.SetPadded(false)
	c.SetContent(obj)

	gui.CurrentApp().Settings().SetTheme(t)
	app.ApplyThemeTo(obj, c)
	return c.Capture()
}
