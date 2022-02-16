// !build test

package theme

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
	"net/url"

	gui "github.com/bhojpur/gui/pkg/engine"
)

type themedApp struct {
	theme gui.Theme
}

func (t *themedApp) BuildType() gui.BuildType {
	return gui.BuildStandard
}

func (t *themedApp) NewWindow(title string) gui.Window {
	return nil
}

func (t *themedApp) OpenURL(url *url.URL) error {
	return nil
}

func (t *themedApp) Icon() gui.Resource {
	return nil
}

func (t *themedApp) SetIcon(gui.Resource) {
}

func (t *themedApp) Run() {
}

func (t *themedApp) Quit() {
}

func (t *themedApp) Driver() gui.Driver {
	return nil
}

func (t *themedApp) UniqueID() string {
	return ""
}

func (t *themedApp) SendNotification(notification *gui.Notification) {
}

func (t *themedApp) Settings() gui.Settings {
	return t
}

func (t *themedApp) Storage() gui.Storage {
	return nil
}

func (t *themedApp) Preferences() gui.Preferences {
	return nil
}

func (t *themedApp) Lifecycle() gui.Lifecycle {
	return nil
}

func (t *themedApp) PrimaryColor() string {
	return ColorBlue
}

func (t *themedApp) Theme() gui.Theme {
	return t.theme
}

func (t *themedApp) SetTheme(theme gui.Theme) {
	t.theme = theme
}

func (t *themedApp) ThemeVariant() gui.ThemeVariant {
	return VariantDark
}

func (t *themedApp) Scale() float32 {
	return 1.0
}

func (t *themedApp) AddChangeListener(chan gui.Settings) {
}
