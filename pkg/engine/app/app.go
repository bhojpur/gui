package app // import "github.com/bhojpur/gui/pkg/engine/app"

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

// It provides implementations for working with Bhojpur GUI application framework.
// The fastest way to get started is to call app.New() which will normally load a
// new desktop application. If the "ci" tag is passed to Go (go run -tags ci myapp.go)
// it will run an in-memory application.

import (
	"strconv"
	"sync/atomic"
	"time"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/internal"
	"github.com/bhojpur/gui/pkg/engine/internal/app"
	intRepo "github.com/bhojpur/gui/pkg/engine/internal/repository"
	"github.com/bhojpur/gui/pkg/engine/storage/repository"

	"golang.org/x/sys/execabs"
)

// Declare conformity with Bhojpur GUI application interface
var _ gui.App = (*bhojpurApp)(nil)

type bhojpurApp struct {
	driver   gui.Driver
	icon     gui.Resource
	uniqueID string

	lifecycle gui.Lifecycle
	settings  *settings
	storage   *store
	prefs     gui.Preferences

	running uint32 // atomic, 1 == running, 0 == stopped
	exec    func(name string, arg ...string) *execabs.Cmd
}

func (a *bhojpurApp) Icon() gui.Resource {
	return a.icon
}

func (a *bhojpurApp) SetIcon(icon gui.Resource) {
	a.icon = icon
}

func (a *bhojpurApp) UniqueID() string {
	if a.uniqueID != "" {
		return a.uniqueID
	}

	gui.LogError("Preferences API requires a unique ID, use app.NewWithID()", nil)
	a.uniqueID = "missing-id-" + strconv.FormatInt(time.Now().Unix(), 10) // This is a fake unique - it just has to not be reused...
	return a.uniqueID
}

func (a *bhojpurApp) NewWindow(title string) gui.Window {
	return a.driver.CreateWindow(title)
}

func (a *bhojpurApp) Run() {
	if atomic.CompareAndSwapUint32(&a.running, 0, 1) {
		a.driver.Run()
		return
	}
}

func (a *bhojpurApp) Quit() {
	for _, window := range a.driver.AllWindows() {
		window.Close()
	}

	a.driver.Quit()
	a.settings.stopWatching()
	atomic.StoreUint32(&a.running, 0)
}

func (a *bhojpurApp) Driver() gui.Driver {
	return a.driver
}

// Settings returns the application settings currently configured.
func (a *bhojpurApp) Settings() gui.Settings {
	return a.settings
}

func (a *bhojpurApp) Storage() gui.Storage {
	return a.storage
}

func (a *bhojpurApp) Preferences() gui.Preferences {
	if a.uniqueID == "" {
		gui.LogError("Preferences API requires a unique ID, use app.NewWithID()", nil)
	}
	return a.prefs
}

func (a *bhojpurApp) Lifecycle() gui.Lifecycle {
	return a.lifecycle
}

// New returns a new Bhojpur GUI application instance with the default driver and no unique ID
func New() gui.App {
	internal.LogHint("Applications should be created with a unique ID using app.NewWithID()")
	return NewWithID("")
}

func newAppWithDriver(d gui.Driver, id string) gui.App {
	newApp := &bhojpurApp{uniqueID: id, driver: d, exec: execabs.Command, lifecycle: &app.Lifecycle{}}
	gui.SetCurrentApp(newApp)
	newApp.settings = loadSettings()

	newApp.prefs = newPreferences(newApp)
	newApp.storage = &store{a: newApp}
	if id != "" {
		if pref, ok := newApp.prefs.(interface{ load() }); ok {
			pref.load()
		}

		root, _ := newApp.storage.docRootURI()
		newApp.storage.Docs = &internal.Docs{RootDocURI: root}
	} else {
		newApp.storage.Docs = &internal.Docs{} // an empty impl to avoid crashes
	}

	if !d.Device().IsMobile() {
		newApp.settings.watchSettings()
	}

	repository.Register("http", intRepo.NewHTTPRepository())
	repository.Register("https", intRepo.NewHTTPRepository())

	return newApp
}
