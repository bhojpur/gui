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
	"net/url"
	"strings"
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"
	_ "github.com/bhojpur/gui/pkg/engine/test"
	"github.com/stretchr/testify/assert"

	"golang.org/x/sys/execabs"
)

func TestDummyApp(t *testing.T) {
	app := NewWithID("net.bhojpur.test")

	app.Quit()
}

func TestCurrentApp(t *testing.T) {
	app := NewWithID("net.bhojpur.test")

	assert.Equal(t, app, gui.CurrentApp())
}

func TestBhojpurApp_UniqueID(t *testing.T) {
	appID := "net.bhojpur.test"
	app := NewWithID(appID)

	assert.Equal(t, appID, app.UniqueID())
}

func TestBhojpurApp_OpenURL(t *testing.T) {
	opened := ""
	app := NewWithID("net.bhojpur.test")
	app.(*bhojpurApp).exec = func(cmd string, arg ...string) *execabs.Cmd {
		opened = arg[len(arg)-1]
		return execabs.Command("")
	}

	urlStr := "https://bhojpur.net"
	u, _ := url.Parse(urlStr)
	err := app.OpenURL(u)

	if err != nil && strings.Contains(err.Error(), "unknown operating system") {
		return // when running in CI mode we don't actually open URLs...
	}

	assert.Equal(t, urlStr, opened)
}
