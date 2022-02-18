package settings

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
	"encoding/json"
	"image"
	"image/color"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/app"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/container"

	//"github.com/bhojpur/gui/pkg/engine/internal/cache"
	"github.com/bhojpur/gui/pkg/engine/layout"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/tools/playground"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

const (
	systemThemeName = "system default"
)

// Settings gives access to user interfaces to control Bhojpur GUI settings
type Settings struct {
	guiSettings app.SettingsSchema

	preview *canvas.Image
	colors  []gui.CanvasObject
}

// NewSettings returns a new settings instance with the current configuration loaded
func NewSettings() *Settings {
	s := &Settings{}
	s.load()
	if s.guiSettings.Scale == 0 {
		s.guiSettings.Scale = 1
	}
	return s
}

// AppearanceIcon returns the icon for appearance settings
func (s *Settings) AppearanceIcon() gui.Resource {
	return theme.NewThemedResource(resourceAppearanceSvg)
}

// LoadAppearanceScreen creates a new settings screen to handle appearance configuration
func (s *Settings) LoadAppearanceScreen(w gui.Window) gui.CanvasObject {
	s.preview = canvas.NewImageFromImage(s.createPreview())
	s.preview.FillMode = canvas.ImageFillContain

	def := s.guiSettings.ThemeName
	themeNames := []string{"dark", "light"}
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		themeNames = append(themeNames, systemThemeName)
		if s.guiSettings.ThemeName == "" {
			def = systemThemeName
		}
	}
	themes := widget.NewSelect(themeNames, s.chooseTheme)
	themes.SetSelected(def)

	scale := s.makeScaleGroup(w.Canvas().Scale())
	box := container.NewVBox(scale)

	for _, c := range theme.PrimaryColorNames() {
		b := newColorButton(c, theme.PrimaryColorNamed(c), s)
		s.colors = append(s.colors, b)
	}
	swatch := container.NewGridWithColumns(len(s.colors), s.colors...)
	appearance := widget.NewForm(widget.NewFormItem("Main Color", swatch),
		widget.NewFormItem("Theme", themes))

	box.Add(widget.NewCard("Appearance", "", appearance))
	bottom := container.NewHBox(layout.NewSpacer(),
		&widget.Button{Text: "Apply", Importance: widget.HighImportance, OnTapped: func() {
			if s.guiSettings.Scale == 0.0 {
				s.chooseScale(1.0)
			}
			err := s.save()
			if err != nil {
				gui.LogError("Failed on saving", err)
			}

			s.appliedScale(s.guiSettings.Scale)
		}})

	return container.NewBorder(box, bottom, nil, nil, s.preview)
}

func (s *Settings) chooseTheme(name string) {
	if name == systemThemeName {
		name = ""
	}
	s.guiSettings.ThemeName = name

	s.preview.Image = s.createPreview()
	canvas.Refresh(s.preview)
}

type overrideTheme interface {
	OverrideTheme(gui.Theme, string)
}

func (s *Settings) createPreview() image.Image {
	c := playground.NewSoftwareCanvas()
	oldTheme := gui.CurrentApp().Settings().Theme()
	oldColor := gui.CurrentApp().Settings().PrimaryColor()

	th := oldTheme
	if s.guiSettings.ThemeName == "light" {
		th = theme.LightTheme()
	} else if s.guiSettings.ThemeName == "dark" {
		th = theme.DarkTheme()
	}

	//cache.ResetThemeCaches() // reset icon cache
	gui.CurrentApp().Settings().(overrideTheme).OverrideTheme(th, s.guiSettings.PrimaryColor)

	empty := widget.NewLabel("")
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("Home")),
		container.NewTabItemWithIcon("Browse", theme.ComputerIcon(), empty),
		container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), empty),
		container.NewTabItemWithIcon("Help", theme.HelpIcon(), empty))
	tabs.SetTabLocation(container.TabLocationLeading)
	showOverlay(c)

	c.SetContent(tabs)
	c.Resize(gui.NewSize(380, 380))
	// wait for indicator animation
	time.Sleep(canvas.DurationShort)
	img := c.Capture()

	//cache.ResetThemeCaches() // ensure we re-create the correct cached assets
	gui.CurrentApp().Settings().(overrideTheme).OverrideTheme(oldTheme, oldColor)
	return img
}

func (s *Settings) load() {
	err := s.loadFromFile(s.guiSettings.StoragePath())
	if err != nil {
		gui.LogError("Settings load error:", err)
	}
}

func (s *Settings) loadFromFile(path string) error {
	file, err := os.Open(path) // #nosec
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(filepath.Dir(path), 0700)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	decode := json.NewDecoder(file)

	return decode.Decode(&s.guiSettings)
}

func (s *Settings) save() error {
	return s.saveToFile(s.guiSettings.StoragePath())
}

func (s *Settings) saveToFile(path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil { // this is not an exists error according to docs
		return err
	}

	data, err := json.Marshal(&s.guiSettings)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}

type colorButton struct {
	widget.BaseWidget
	name  string
	color color.Color

	s *Settings
}

func newColorButton(n string, c color.Color, s *Settings) *colorButton {
	b := &colorButton{name: n, color: c, s: s}
	b.ExtendBaseWidget(b)
	return b
}

func (c *colorButton) CreateRenderer() gui.WidgetRenderer {
	r := canvas.NewRectangle(c.color)
	r.StrokeWidth = 5

	if c.name == c.s.guiSettings.PrimaryColor {
		r.StrokeColor = theme.PrimaryColor()
	}

	return &colorRenderer{c: c, rect: r, objs: []gui.CanvasObject{r}}
}

func (c *colorButton) Tapped(_ *gui.PointEvent) {
	c.s.guiSettings.PrimaryColor = c.name
	for _, child := range c.s.colors {
		child.Refresh()
	}

	c.s.preview.Image = c.s.createPreview()
	canvas.Refresh(c.s.preview)
}

type colorRenderer struct {
	c    *colorButton
	rect *canvas.Rectangle
	objs []gui.CanvasObject
}

func (c *colorRenderer) Layout(s gui.Size) {
	c.rect.Resize(s)
}

func (c *colorRenderer) MinSize() gui.Size {
	return gui.NewSize(20, 20)
}

func (c *colorRenderer) Refresh() {
	if c.c.name == c.c.s.guiSettings.PrimaryColor {
		c.rect.StrokeColor = theme.PrimaryColor()
	} else {
		c.rect.StrokeColor = color.Transparent
	}
	c.rect.FillColor = c.c.color

	c.rect.Refresh()
}

func (c *colorRenderer) Objects() []gui.CanvasObject {
	return c.objs
}

func (c *colorRenderer) Destroy() {
}

func showOverlay(c gui.Canvas) {
	username := widget.NewEntry()
	password := widget.NewPasswordEntry()
	form := widget.NewForm(widget.NewFormItem("Username", username),
		widget.NewFormItem("Password", password))
	form.OnCancel = func() {}
	form.OnSubmit = func() {}
	content := container.NewVBox(
		widget.NewLabelWithStyle("Login demo", gui.TextAlignCenter, gui.TextStyle{Bold: true}), form)
	wrap := container.NewWithoutLayout(content)
	wrap.Resize(content.MinSize().Add(gui.NewSize(theme.Padding()*2, theme.Padding()*2)))
	content.Resize(content.MinSize())
	content.Move(gui.NewPos(theme.Padding(), theme.Padding()))

	over := container.NewMax(
		canvas.NewRectangle(theme.ShadowColor()), container.NewCenter(wrap),
	)

	c.Overlays().Add(over)
	c.Focus(username)
}
