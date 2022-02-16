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
	"os"
	"path/filepath"
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

func TestSettingsBuildType(t *testing.T) {
	set := test.NewApp().Settings()
	assert.Equal(t, gui.BuildStandard, set.BuildType()) // during test we should have a normal build

	set = &settings{}
	assert.Equal(t, buildMode, set.BuildType()) // when testing this package only it could be debug or release
}

func TestSettingsLoad(t *testing.T) {
	settings := &settings{}

	err := settings.loadFromFile(filepath.Join("testdata", "light-theme.json"))
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "light", settings.schema.ThemeName)

	err = settings.loadFromFile(filepath.Join("testdata", "dark-theme.json"))
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "dark", settings.schema.ThemeName)
}

func TestOverrideTheme(t *testing.T) {
	set := &settings{}
	set.setupTheme()
	assert.Equal(t, defaultVariant(), set.ThemeVariant())

	set.schema.ThemeName = "light"
	set.setupTheme()
	assert.Equal(t, theme.LightTheme(), set.Theme())

	set.schema.ThemeName = "dark"
	set.setupTheme()
	assert.Equal(t, theme.DarkTheme(), set.Theme())

	set = &settings{}
	set.setupTheme()
	assert.Equal(t, defaultVariant(), set.ThemeVariant())

	err := os.Setenv("BHOJPUR_THEME", "light")
	if err != nil {
		t.Error(err)
	}
	set.setupTheme()
	assert.Equal(t, theme.LightTheme(), set.Theme())
	err = os.Setenv("BHOJPUR_THEME", "")
	if err != nil {
		t.Error(err)
	}
}

func TestOverrideTheme_IgnoresSettingsChange(t *testing.T) {
	// check that a file-load does not overwrite our value
	set := &settings{}
	err := os.Setenv("BHOJPUR_THEME", "light")
	if err != nil {
		t.Error(err)
	}
	set.setupTheme()
	assert.Equal(t, theme.LightTheme(), set.Theme())

	err = set.loadFromFile(filepath.Join("testdata", "dark-theme.json"))
	if err != nil {
		t.Error(err)
	}
	set.setupTheme()
	assert.Equal(t, theme.LightTheme(), set.Theme())
	err = os.Setenv("BHOJPUR_THEME", "")
	if err != nil {
		t.Error(err)
	}
}

func TestCustomTheme(t *testing.T) {
	type customTheme struct {
		gui.Theme
	}
	set := &settings{}
	ctheme := &customTheme{theme.LightTheme()}
	set.SetTheme(ctheme)

	set.setupTheme()
	assert.True(t, set.Theme() == ctheme)
	assert.Equal(t, defaultVariant(), set.ThemeVariant())

	err := set.loadFromFile(filepath.Join("testdata", "light-theme.json"))
	if err != nil {
		t.Error(err)
	}
	set.setupTheme()
	assert.True(t, set.Theme() == ctheme)
	assert.Equal(t, theme.VariantLight, set.ThemeVariant())

	err = set.loadFromFile(filepath.Join("testdata", "dark-theme.json"))
	if err != nil {
		t.Error(err)
	}
	set.setupTheme()
	assert.True(t, set.Theme() == ctheme)
	assert.Equal(t, theme.VariantDark, set.ThemeVariant())

	err = os.Setenv("BHOJPUR_THEME", "light")
	if err != nil {
		t.Error(err)
	}
	set.setupTheme()
	assert.True(t, set.Theme() == ctheme)
	assert.Equal(t, theme.VariantLight, set.ThemeVariant())

	err = os.Setenv("BHOJPUR_THEME", "dark")
	if err != nil {
		t.Error(err)
	}
	set.setupTheme()
	assert.True(t, set.Theme() == ctheme)
	assert.Equal(t, theme.VariantDark, set.ThemeVariant())

	err = os.Setenv("BHOJPUR_THEME", "")
	if err != nil {
		t.Error(err)
	}
	set.setupTheme()
	assert.True(t, set.Theme() == ctheme)
	assert.Equal(t, theme.VariantDark, set.ThemeVariant())
}
