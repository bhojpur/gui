// It provides utility drivers for running Bhojpur GUI tests without rendering
package test // import gui "github.com/bhojpur/gui/pkg/engine/test"

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
	"sync"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/internal"
	"github.com/bhojpur/gui/pkg/engine/internal/app"
	"github.com/bhojpur/gui/pkg/engine/internal/cache"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

// ensure we have a dummy app loaded and ready to test
func init() {
	NewApp()
}

type testApp struct {
	driver       *testDriver
	settings     gui.Settings
	prefs        gui.Preferences
	propertyLock sync.RWMutex
	storage      gui.Storage
	lifecycle    gui.Lifecycle

	// user action variables
	appliedTheme     gui.Theme
	lastNotification *gui.Notification
}

func (a *testApp) Icon() gui.Resource {
	return nil
}

func (a *testApp) SetIcon(gui.Resource) {
	// no-op
}

func (a *testApp) NewWindow(title string) gui.Window {
	return a.driver.CreateWindow(title)
}

func (a *testApp) OpenURL(url *url.URL) error {
	// no-op
	return nil
}

func (a *testApp) Run() {
	// no-op
}

func (a *testApp) Quit() {
	// no-op
}

func (a *testApp) UniqueID() string {
	return "testApp" // TODO should this be randomised?
}

func (a *testApp) Driver() gui.Driver {
	return a.driver
}

func (a *testApp) SendNotification(notify *gui.Notification) {
	a.propertyLock.Lock()
	defer a.propertyLock.Unlock()

	a.lastNotification = notify
}

func (a *testApp) Settings() gui.Settings {
	return a.settings
}

func (a *testApp) Preferences() gui.Preferences {
	return a.prefs
}

func (a *testApp) Storage() gui.Storage {
	return a.storage
}

func (a *testApp) Lifecycle() gui.Lifecycle {
	return a.lifecycle
}

func (a *testApp) lastAppliedTheme() gui.Theme {
	a.propertyLock.Lock()
	defer a.propertyLock.Unlock()

	return a.appliedTheme
}

// NewApp returns a new dummy app used for testing.
// It loads a test driver which creates a virtual window in memory for testing.
func NewApp() gui.App {
	settings := &testSettings{scale: 1.0, theme: Theme()}
	prefs := internal.NewInMemoryPreferences()
	store := &testStorage{}
	test := &testApp{settings: settings, prefs: prefs, storage: store, driver: NewDriver().(*testDriver),
		lifecycle: &app.Lifecycle{}}
	root, _ := store.docRootURI()
	store.Docs = &internal.Docs{RootDocURI: root}
	cache.ResetThemeCaches()
	gui.SetCurrentApp(test)

	listener := make(chan gui.Settings)
	test.Settings().AddChangeListener(listener)
	go func() {
		for {
			<-listener
			cache.ResetThemeCaches()
			app.ApplySettings(test.Settings(), test)

			test.propertyLock.Lock()
			test.appliedTheme = test.Settings().Theme()
			test.propertyLock.Unlock()
		}
	}()

	return test
}

type testSettings struct {
	theme gui.Theme
	scale float32

	changeListeners []chan gui.Settings
	propertyLock    sync.RWMutex
}

func (s *testSettings) AddChangeListener(listener chan gui.Settings) {
	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()
	s.changeListeners = append(s.changeListeners, listener)
}

func (s *testSettings) BuildType() gui.BuildType {
	return gui.BuildStandard
}

func (s *testSettings) PrimaryColor() string {
	return theme.ColorBlue
}

func (s *testSettings) SetTheme(theme gui.Theme) {
	s.propertyLock.Lock()
	s.theme = theme
	s.propertyLock.Unlock()

	s.apply()
}

func (s *testSettings) Theme() gui.Theme {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()

	if s.theme == nil {
		return theme.DarkTheme()
	}

	return s.theme
}

func (s *testSettings) ThemeVariant() gui.ThemeVariant {
	return 2 // not a preference
}

func (s *testSettings) Scale() float32 {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()
	return s.scale
}

func (s *testSettings) apply() {
	s.propertyLock.RLock()
	listeners := s.changeListeners
	s.propertyLock.RUnlock()

	for _, listener := range listeners {
		listener <- s
	}
}
