package data

//go:generate guiutl bundle -package data -o bundled.go assets

import (
	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

// ThemedResource is a resource wrapper that will return an appropriate resource
// for the currently selected theme.
type ThemedResource struct {
	dark, light gui.Resource
}

func isLight() bool {
	r, g, b, _ := theme.ForegroundColor().RGBA()
	return r < 0xaaaa && g < 0xaaaa && b < 0xaaaa
}

// Name returns the underlying resource name (used for caching)
func (res *ThemedResource) Name() string {
	if isLight() {
		return res.light.Name()
	}
	return res.dark.Name()
}

// Content returns the underlying content of the correct resource for the current theme
func (res *ThemedResource) Content() []byte {
	if isLight() {
		return res.light.Content()
	}
	return res.dark.Content()
}

// NewThemedResource creates a resource that adapts to the current theme setting.
func NewThemedResource(dark, light gui.Resource) *ThemedResource {
	return &ThemedResource{dark, light}
}

// BhojpurScene contains the full Bhojpur GUI logo with background design
var BhojpurScene = NewThemedResource(resourceBhojpurscenedarkPng, resourceBhojpurscenelightPng)
