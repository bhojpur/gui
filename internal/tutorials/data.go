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
	gui "github.com/bhojpur/gui/pkg/engine"
)

// Tutorial defines the data structure for a tutorial
type Tutorial struct {
	Title, Intro string
	View         func(w gui.Window) gui.CanvasObject
}

var (
	// Tutorials defines the metadata for each tutorial
	Tutorials = map[string]Tutorial{
		"welcome": {"Welcome", "", welcomeScreen},
		"canvas": {"Canvas",
			"See the canvas capabilities.",
			canvasScreen,
		},
		"animations": {"Animations",
			"See how to animate components.",
			makeAnimationScreen,
		},
		"icons": {"Theme Icons",
			"Browse the embedded icons.",
			iconScreen,
		},
		"containers": {"Containers",
			"Containers group other widgets and canvas objects, organising according to their layout.\n" +
				"Standard containers are illustrated in this section, but developers can also provide custom " +
				"layouts using the gui.NewContainerWithLayout() constructor.",
			containerScreen,
		},
		"apptabs": {"AppTabs",
			"A container to help divide up an application into functional areas.",
			makeAppTabsTab,
		},
		"border": {"Border",
			"A container that positions items around a central content.",
			makeBorderLayout,
		},
		"box": {"Box",
			"A container arranges items in horizontal or vertical list.",
			makeBoxLayout,
		},
		"center": {"Center",
			"A container to that centers child elements.",
			makeCenterLayout,
		},
		"doctabs": {"DocTabs",
			"A container to display a single document from a set of many.",
			makeDocTabsTab,
		},
		"grid": {"Grid",
			"A container that arranges all items in a grid.",
			makeGridLayout,
		},
		"split": {"Split",
			"A split container divides the container in two pieces that the user can resize.",
			makeSplitTab,
		},
		"scroll": {"Scroll",
			"A container that provides scrolling for it's content.",
			makeScrollTab,
		},
		"widgets": {"Widgets",
			"In this section you can see the features available in the toolkit widget set.\n" +
				"Expand the tree on the left to browse the individual tutorial elements.",
			widgetScreen,
		},
		"accordion": {"Accordion",
			"Expand or collapse content panels.",
			makeAccordionTab,
		},
		"button": {"Button",
			"Simple widget for user tap handling.",
			makeButtonTab,
		},
		"card": {"Card",
			"Group content and widgets.",
			makeCardTab,
		},
		"entry": {"Entry",
			"Different ways to use the entry widget.",
			makeEntryTab,
		},
		"form": {"Form",
			"Gathering input widgets for data submission.",
			makeFormTab,
		},
		"input": {"Input",
			"A collection of widgets for user input.",
			makeInputTab,
		},
		"text": {"Text",
			"Text handling widgets.",
			makeTextTab,
		},
		"toolbar": {"Toolbar",
			"A row of shortcut icons for common tasks.",
			makeToolbarTab,
		},
		"progress": {"Progress",
			"Show duration or the need to wait for a task.",
			makeProgressTab,
		},
		"collections": {"Collections",
			"Collection widgets provide an efficient way to present lots of content.\n" +
				"The List, Table, and Tree provide a cache and re-use mechanism that make it possible to scroll through thousands of elements.\n" +
				"Use this for large data sets or for collections that can expand as users scroll.",
			collectionScreen,
		},
		"list": {"List",
			"A vertical arrangement of cached elements with the same styling.",
			makeListTab,
		},
		"table": {"Table",
			"A two dimensional cached collection of cells.",
			makeTableTab,
		},
		"tree": {"Tree",
			"A tree based arrangement of cached elements with the same styling.",
			makeTreeTab,
		},
		"dialogs": {"Dialogs",
			"Work with dialogs.",
			dialogScreen,
		},
		"windows": {"Windows",
			"Window function demo.",
			windowScreen,
		},
		"binding": {"Data Binding",
			"Connecting widgets to a data source.",
			bindingScreen},
		"advanced": {"Advanced",
			"Debug and advanced information.",
			advancedScreen,
		},
	}

	// TutorialIndex  defines how our tutorials should be laid out in the index tree
	TutorialIndex = map[string][]string{
		"":            {"welcome", "canvas", "animations", "icons", "widgets", "collections", "containers", "dialogs", "windows", "binding", "advanced"},
		"collections": {"list", "table", "tree"},
		"containers":  {"apptabs", "border", "box", "center", "doctabs", "grid", "scroll", "split"},
		"widgets":     {"accordion", "button", "card", "entry", "form", "input", "progress", "text", "toolbar"},
	}
)
