package engine

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
)

func testFunc1() {}
func testFunc2() {}
func testFunc3() {}
func testFunc4() {}

var menuItemsTest = []struct {
	Label  string
	Action func()
	Item   *MenuItem
}{
	{"item1", testFunc1, &MenuItem{Label: "item1", Action: testFunc1}},
	{"item2", testFunc2, &MenuItem{Label: "item2", Action: testFunc2}},
	{"item3", testFunc3, &MenuItem{Label: "item3", Action: testFunc3}},
	{"item4", testFunc4, &MenuItem{Label: "item4", Action: testFunc4}},
}

var menuTest = []struct {
	Label string
}{
	{"menu1"},
	{"menu2"},
}

func TestNewMainMenu(t *testing.T) {
	var items []*MenuItem
	var menus []*Menu

	for _, tt := range menuItemsTest {
		t.Run(tt.Label, func(t *testing.T) {
			item := NewMenuItem(tt.Label, tt.Action)
			// Compare sprinted address
			if reflect.ValueOf(item.Action) != reflect.ValueOf(tt.Action) {
				t.Errorf("Expected %v but got %v", reflect.ValueOf(tt.Action), reflect.ValueOf(item.Action))
			}
			if item.Label != tt.Label {
				t.Errorf("Expected %v but got %v", tt.Label, item.Label)
			}
			items = append(items, item)
		})
	}

	if len(items) < 4 {
		t.Errorf("Expected %d menu items but got %d", len(menuItemsTest), len(items))
	}

	for _, tt := range menuTest {
		t.Run(tt.Label, func(t *testing.T) {
			menu := NewMenu(tt.Label, items...)

			if menu.Label != tt.Label {
				t.Errorf("Expected menu label %s but got %s", tt.Label, menu.Label)
			}

			if !reflect.DeepEqual(menu.Items, items) {
				t.Errorf("Expected items to resemble what was inputted, but got %v", menu.Items)
			}

			menus = append(menus, menu)
		})
	}

	mm := NewMainMenu(menus...)

	if !reflect.DeepEqual(mm.Items, menus) {
		t.Errorf("Expected main menu to contain all submenus but got %v", mm.Items)
	}

}
