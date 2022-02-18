//go:build !ci && !js && !wasm && !test_web_driver && (linux || openbsd || freebsd || netbsd) && !android
// +build !ci
// +build !js
// +build !wasm
// +build !test_web_driver
// +build linux openbsd freebsd netbsd
// +build !android

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
	"os"
	"path/filepath"

	"github.com/godbus/dbus/v5"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

func defaultVariant() gui.ThemeVariant {
	return theme.VariantDark
}

func (a *bhojpurApp) OpenURL(url *url.URL) error {
	cmd := a.exec("xdg-open", url.String())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Start()
}

func (a *bhojpurApp) SendNotification(n *gui.Notification) {
	conn, err := dbus.SessionBus() // shared connection, don't close
	if err != nil {
		gui.LogError("Unable to connect to session D-Bus", err)
		return
	}

	appName := gui.CurrentApp().UniqueID()
	appIcon := ""       // TODO in the future extract icon from app to cache and pass path
	timeout := int32(0) // we don't support this yet

	obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	call := obj.Call("org.freedesktop.Notifications.Notify", 0, appName, uint32(0),
		appIcon, n.Title, n.Content, []string{}, map[string]dbus.Variant{}, timeout)
	if call.Err != nil {
		gui.LogError("Failed to send message to bus", call.Err)
	}
}

func rootConfigDir() string {
	desktopConfig, _ := os.UserConfigDir()
	return filepath.Join(desktopConfig, "bhojpur")
}

func watchTheme() {
	// no-op, not able to read linux theme in a standard way
}
