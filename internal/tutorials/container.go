package tutorials

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
	"fmt"
	"image/color"
	"strconv"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/container"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

// containerScreen loads a tab panel for containers
func containerScreen(_ gui.Window) gui.CanvasObject {
	content := container.NewBorder(
		widget.NewLabelWithStyle("Top", gui.TextAlignCenter, gui.TextStyle{}),
		widget.NewLabelWithStyle("Bottom", gui.TextAlignCenter, gui.TextStyle{}),
		widget.NewLabel("Left"),
		widget.NewLabel("Right"),
		widget.NewLabel("Border Container"))
	return container.NewCenter(content)
}

func makeAppTabsTab(_ gui.Window) gui.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Tab 1", widget.NewLabel("Content of tab 1")),
		container.NewTabItem("Tab 2 bigger", widget.NewLabel("Content of tab 2")),
		container.NewTabItem("Tab 3", widget.NewLabel("Content of tab 3")),
	)
	for i := 4; i <= 12; i++ {
		tabs.Append(container.NewTabItem(fmt.Sprintf("Tab %d", i), widget.NewLabel(fmt.Sprintf("Content of tab %d", i))))
	}
	locations := makeTabLocationSelect(tabs.SetTabLocation)
	return container.NewBorder(locations, nil, nil, nil, tabs)
}

func makeBorderLayout(_ gui.Window) gui.CanvasObject {
	top := makeCell()
	bottom := makeCell()
	left := makeCell()
	right := makeCell()
	middle := widget.NewLabelWithStyle("BorderLayout", gui.TextAlignCenter, gui.TextStyle{})

	return container.NewBorder(top, bottom, left, right, middle)
}

func makeBoxLayout(_ gui.Window) gui.CanvasObject {
	top := makeCell()
	bottom := makeCell()
	middle := widget.NewLabel("BoxLayout")
	center := makeCell()
	right := makeCell()

	col := container.NewVBox(top, middle, bottom)

	return container.NewHBox(col, center, right)
}

func makeButtonList(count int) []gui.CanvasObject {
	var items []gui.CanvasObject
	for i := 1; i <= count; i++ {
		index := i // capture
		items = append(items, widget.NewButton("Button "+strconv.Itoa(index), func() {
			fmt.Println("Tapped", index)
		}))
	}

	return items
}

func makeCell() gui.CanvasObject {
	rect := canvas.NewRectangle(&color.NRGBA{128, 128, 128, 255})
	rect.SetMinSize(gui.NewSize(30, 30))
	return rect
}

func makeCenterLayout(_ gui.Window) gui.CanvasObject {
	middle := widget.NewButton("CenterLayout", func() {})

	return container.NewCenter(middle)
}

func makeDocTabsTab(_ gui.Window) gui.CanvasObject {
	tabs := container.NewDocTabs(
		container.NewTabItem("Doc 1", widget.NewLabel("Content of document 1")),
		container.NewTabItem("Doc 2 bigger", widget.NewLabel("Content of document 2")),
		container.NewTabItem("Doc 3", widget.NewLabel("Content of document 3")),
	)
	i := 3
	tabs.CreateTab = func() *container.TabItem {
		i++
		return container.NewTabItem(fmt.Sprintf("Doc %d", i), widget.NewLabel(fmt.Sprintf("Content of document %d", i)))
	}
	locations := makeTabLocationSelect(tabs.SetTabLocation)
	return container.NewBorder(locations, nil, nil, nil, tabs)
}

func makeGridLayout(_ gui.Window) gui.CanvasObject {
	box1 := makeCell()
	box2 := widget.NewLabel("Grid")
	box3 := makeCell()
	box4 := makeCell()

	return container.NewGridWithColumns(2,
		box1, box2, box3, box4)
}

func makeScrollTab(_ gui.Window) gui.CanvasObject {
	hlist := makeButtonList(20)
	vlist := makeButtonList(50)

	horiz := container.NewHScroll(container.NewHBox(hlist...))
	vert := container.NewVScroll(container.NewVBox(vlist...))

	return container.NewAdaptiveGrid(2,
		container.NewBorder(horiz, nil, nil, nil, vert),
		makeScrollBothTab())
}

func makeScrollBothTab() gui.CanvasObject {
	logo := canvas.NewImageFromResource(theme.BhojpurLogo())
	logo.SetMinSize(gui.NewSize(800, 800))

	scroll := container.NewScroll(logo)
	scroll.Resize(gui.NewSize(400, 400))

	return scroll
}

func makeSplitTab(_ gui.Window) gui.CanvasObject {
	left := widget.NewMultiLineEntry()
	left.Wrapping = gui.TextWrapWord
	left.SetText("Long text is looooooooooooooong")
	right := container.NewVSplit(
		widget.NewLabel("Label"),
		widget.NewButton("Button", func() { fmt.Println("button tapped!") }),
	)
	return container.NewHSplit(container.NewVScroll(left), right)
}

func makeTabLocationSelect(callback func(container.TabLocation)) *widget.Select {
	locations := widget.NewSelect([]string{"Top", "Bottom", "Leading", "Trailing"}, func(s string) {
		callback(map[string]container.TabLocation{
			"Top":      container.TabLocationTop,
			"Bottom":   container.TabLocationBottom,
			"Leading":  container.TabLocationLeading,
			"Trailing": container.TabLocationTrailing,
		}[s])
	})
	locations.SetSelected("Top")
	return locations
}
