package app

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"bytes"
	"os"
	"path/filepath"
	"sync"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

// SettingsSchema is used for loading and storing global settings
type SettingsSchema struct {
	// these items are used for global settings load
	ThemeName    string  `json:"theme"`
	Scale        float32 `json:"scale"`
	PrimaryColor string  `json:"primary_color"`
}

// StoragePath returns the location of the settings storage
func (sc *SettingsSchema) StoragePath() string {
	return filepath.Join(rootConfigDir(), "settings.json")
}

// Declare conformity with Settings interface
var _ gui.Settings = (*settings)(nil)

type settings struct {
	propertyLock   sync.RWMutex
	theme          gui.Theme
	themeSpecified bool
	variant        gui.ThemeVariant

	changeListeners sync.Map    // map[chan gui.Settings]bool
	watcher         interface{} // normally *fsnotify.Watcher or nil - avoid import in this file

	schema SettingsSchema
}

func (s *settings) BuildType() gui.BuildType {
	return buildMode
}

func (s *settings) PrimaryColor() string {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()
	return s.schema.PrimaryColor
}

// OverrideTheme allows the settings app to temporarily preview different theme details.
// Please make sure that you remember the original settings and call this again to revert the change.
func (s *settings) OverrideTheme(theme gui.Theme, name string) {
	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()
	s.schema.PrimaryColor = name
	s.theme = theme
}

func (s *settings) Theme() gui.Theme {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()
	return s.theme
}

func (s *settings) SetTheme(theme gui.Theme) {
	s.themeSpecified = true
	s.applyTheme(theme, s.variant)
}

func (s *settings) ThemeVariant() gui.ThemeVariant {
	return s.variant
}

func (s *settings) applyTheme(theme gui.Theme, variant gui.ThemeVariant) {
	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()
	s.variant = variant
	s.theme = theme
	s.apply()
}

func (s *settings) Scale() float32 {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()
	if s.schema.Scale < 0.0 {
		return 1.0 // catching any really old data still using the `-1`  value for "auto" scale
	}
	return s.schema.Scale
}

func (s *settings) AddChangeListener(listener chan gui.Settings) {
	s.changeListeners.Store(listener, true) // the boolean is just a dummy value here.
}

func (s *settings) apply() {
	s.changeListeners.Range(func(key, _ interface{}) bool {
		listener := key.(chan gui.Settings)
		select {
		case listener <- s:
		default:
			l := listener
			go func() { l <- s }()
		}
		return true
	})
}

func (s *settings) fileChanged() {
	s.load()
	s.apply()
}

func (s *settings) loadSystemTheme() gui.Theme {
	path := filepath.Join(rootConfigDir(), "theme.json")
	data, err := gui.LoadResourceFromPath(path)
	if err != nil {
		if !os.IsNotExist(err) {
			gui.LogError("Failed to load user theme file: "+path, err)
		}
		return theme.DefaultTheme()
	}
	if data != nil && data.Content() != nil {
		th, err := theme.FromJSONReader(bytes.NewReader(data.Content()))
		if err == nil {
			return th
		}
		gui.LogError("Failed to parse user theme file: "+path, err)
	}
	return theme.DefaultTheme()
}

func (s *settings) setupTheme() {
	name := s.schema.ThemeName
	if env := os.Getenv("BHOJPUR_THEME"); env != "" {
		name = env
	}

	variant := defaultVariant()
	effectiveTheme := s.theme
	if !s.themeSpecified {
		effectiveTheme = s.loadSystemTheme()
	}
	switch name {
	case "light":
		variant = theme.VariantLight
	case "dark":
		variant = theme.VariantDark
	}

	s.applyTheme(effectiveTheme, variant)
}

func loadSettings() *settings {
	s := &settings{}
	s.load()

	return s
}
