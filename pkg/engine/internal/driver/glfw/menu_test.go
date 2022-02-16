//go:build !ci && !mobile
// +build !ci,!mobile

package glfw

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
	"reflect"
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"

	"github.com/stretchr/testify/assert"
)

func Test_Menu_Empty(t *testing.T) {
	w := createWindow("Menu Test").(*window)
	bar := buildMenuOverlay(gui.NewMainMenu(), w)
	assert.Nil(t, bar) // no bar but does not crash
}

func Test_Menu_AddsQuit(t *testing.T) {
	w := createWindow("Menu Test").(*window)
	mainMenu := gui.NewMainMenu(gui.NewMenu("File"))
	bar := buildMenuOverlay(mainMenu, w)
	assert.NotNil(t, bar)
	assert.Equal(t, 1, len(mainMenu.Items))
	assert.Equal(t, 2, len(mainMenu.Items[0].Items)) // separator+quit inserted
	assert.True(t, mainMenu.Items[0].Items[1].IsQuit)
}

func Test_Menu_LeaveQuit(t *testing.T) {
	w := createWindow("Menu Test").(*window)
	quitFunc := func() {}
	mainMenu := gui.NewMainMenu(gui.NewMenu("File", gui.NewMenuItem("Quit", quitFunc)))
	bar := buildMenuOverlay(mainMenu, w)
	assert.NotNil(t, bar)
	assert.Equal(t, 1, len(mainMenu.Items[0].Items)) // no separator added
	assert.Equal(t, reflect.ValueOf(quitFunc).Pointer(), reflect.ValueOf(mainMenu.Items[0].Items[0].Action).Pointer())
}
func Test_Menu_LeaveQuit_AddAction(t *testing.T) {
	w := createWindow("Menu Test").(*window)
	mainMenu := gui.NewMainMenu(gui.NewMenu("File", gui.NewMenuItem("Quit", nil)))
	bar := buildMenuOverlay(mainMenu, w)
	assert.NotNil(t, bar)
	assert.Equal(t, 1, len(mainMenu.Items[0].Items))    // no separator added
	assert.NotNil(t, mainMenu.Items[0].Items[0].Action) // quit action was added
}

func Test_Menu_CustomQuit(t *testing.T) {
	w := createWindow("Menu Test").(*window)

	quitFunc := func() {}
	quitItem := gui.NewMenuItem("Beenden", quitFunc)
	quitItem.IsQuit = true

	mainMenu := gui.NewMainMenu(gui.NewMenu("File", quitItem))
	bar := buildMenuOverlay(mainMenu, w)

	assert.NotNil(t, bar)
	assert.Equal(t, 1, len(mainMenu.Items[0].Items)) // no separator added
	assert.Equal(t, reflect.ValueOf(quitFunc).Pointer(), reflect.ValueOf(mainMenu.Items[0].Items[0].Action).Pointer())
}

func Test_Menu_CustomQuit_NoAction(t *testing.T) {
	w := createWindow("Menu Test").(*window)
	quitItem := gui.NewMenuItem("Beenden", nil)
	quitItem.IsQuit = true
	mainMenu := gui.NewMainMenu(gui.NewMenu("File", quitItem))
	bar := buildMenuOverlay(mainMenu, w)

	assert.NotNil(t, bar)
	assert.Equal(t, 1, len(mainMenu.Items[0].Items))    // no separator added
	assert.NotNil(t, mainMenu.Items[0].Items[0].Action) // quit action was added
}
