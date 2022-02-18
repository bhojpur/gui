package render

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
	"sort"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/container"
	"github.com/bhojpur/gui/pkg/engine/storage"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

// FileTree extends widget.Tree to display a file system hierarchy.
type FileTree struct {
	widget.Tree
	Filter       storage.FileFilter
	ShowRootPath bool
	Sorter       func(gui.URI, gui.URI) bool

	listableCache map[widget.TreeNodeID]gui.ListableURI
	uriCache      map[widget.TreeNodeID]gui.URI
}

// NewFileTree creates a new FileTree from the given root URI.
func NewFileTree(root gui.URI) *FileTree {
	tree := &FileTree{
		Tree: widget.Tree{
			Root: root.String(),
			CreateNode: func(branch bool) gui.CanvasObject {
				var icon gui.CanvasObject
				if branch {
					icon = widget.NewIcon(nil)
				} else {
					icon = widget.NewFileIcon(nil)
				}
				return container.NewBorder(nil, nil, icon, nil, widget.NewLabel("Template Object"))
			},
		},
		listableCache: make(map[widget.TreeNodeID]gui.ListableURI),
		uriCache:      make(map[widget.TreeNodeID]gui.URI),
	}
	tree.IsBranch = func(id widget.TreeNodeID) bool {
		_, err := tree.toListable(id)
		return err == nil
	}
	tree.ChildUIDs = func(id widget.TreeNodeID) (c []string) {
		listable, err := tree.toListable(id)
		if err != nil {
			gui.LogError("Unable to get lister for "+id, err)
			return
		}

		uris, err := listable.List()
		if err != nil {
			gui.LogError("Unable to list "+listable.String(), err)
			return
		}

		for _, u := range tree.sort(tree.filter(uris)) {
			// Convert to String
			c = append(c, u.String())
		}
		return
	}
	tree.UpdateNode = func(id widget.TreeNodeID, branch bool, node gui.CanvasObject) {
		uri, err := tree.toURI(id)
		if err != nil {
			gui.LogError("Unable to parse URI", err)
			return
		}

		c := node.(*gui.Container)
		if branch {
			var r gui.Resource
			if tree.IsBranchOpen(id) {
				// Set open folder icon
				r = theme.FolderOpenIcon()
			} else {
				// Set folder icon
				r = theme.FolderIcon()
			}
			c.Objects[1].(*widget.Icon).SetResource(r)
		} else {
			// Set file uri to update icon
			c.Objects[1].(*widget.FileIcon).SetURI(uri)
		}

		var l string
		if tree.Root == id && tree.ShowRootPath {
			l = id
		} else {
			l = uri.Name()
		}
		c.Objects[0].(*widget.Label).SetText(l)
	}
	tree.ExtendBaseWidget(tree)
	return tree
}

func (t *FileTree) filter(uris []gui.URI) []gui.URI {
	filter := t.Filter
	if filter == nil {
		return uris
	}
	var filtered []gui.URI
	for _, u := range uris {
		if filter.Matches(u) {
			filtered = append(filtered, u)
		}
	}
	return filtered
}

func (t *FileTree) sort(uris []gui.URI) []gui.URI {
	if sorter := t.Sorter; sorter != nil {
		sort.Slice(uris, func(i, j int) bool {
			return sorter(uris[i], uris[j])
		})
	}
	return uris
}

func (t *FileTree) toListable(id widget.TreeNodeID) (gui.ListableURI, error) {
	listable, ok := t.listableCache[id]
	if ok {
		return listable, nil
	}
	uri, err := t.toURI(id)
	if err != nil {
		return nil, err
	}

	listable, err = storage.ListerForURI(uri)
	if err != nil {
		return nil, err
	}
	t.listableCache[id] = listable
	return listable, nil
}

func (t *FileTree) toURI(id widget.TreeNodeID) (gui.URI, error) {
	uri, ok := t.uriCache[id]
	if ok {
		return uri, nil
	}

	uri, err := storage.ParseURI(id)
	if err != nil {
		return nil, err
	}
	t.uriCache[id] = uri
	return uri, nil
}
