package common

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
	"runtime"
	"sync"
	"sync/atomic"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/internal"
	"github.com/bhojpur/gui/pkg/engine/internal/app"
	"github.com/bhojpur/gui/pkg/engine/internal/async"
	"github.com/bhojpur/gui/pkg/engine/internal/cache"
	"github.com/bhojpur/gui/pkg/engine/internal/driver"
	"github.com/bhojpur/gui/pkg/engine/internal/painter/gl"
)

// SizeableCanvas defines a canvas with size related functions.
type SizeableCanvas interface {
	gui.Canvas
	Resize(gui.Size)
	MinSize() gui.Size
}

// Canvas defines common canvas implementation.
type Canvas struct {
	sync.RWMutex

	OnFocus   func(obj gui.Focusable)
	OnUnfocus func()

	impl SizeableCanvas

	contentFocusMgr *app.FocusManager
	menuFocusMgr    *app.FocusManager
	overlays        *overlayStack

	shortcut gui.ShortcutHandler

	painter gl.Painter

	// Any object that requestes to enter to the refresh queue should
	// not be omitted as it is always a rendering task's decision
	// for skipping frames or drawing calls.
	//
	// If an object failed to ender the refresh queue, the object may
	// disappear or blink from the view at any frames. As of this reason,
	// the refreshQueue is an unbounded channel which is bale to cache
	// arbitrary number of gui.CanvasObject for the rendering.
	refreshQueue *async.UnboundedCanvasObjectChan
	refreshCount uint32 // atomic
	dirty        uint32 // atomic

	mWindowHeadTree, contentTree, menuTree *renderCacheTree
}

// AddShortcut adds a shortcut to the canvas.
func (c *Canvas) AddShortcut(shortcut gui.Shortcut, handler func(shortcut gui.Shortcut)) {
	c.shortcut.AddShortcut(shortcut, handler)
}

// EnsureMinSize ensure canvas min size.
//
// This function uses lock.
func (c *Canvas) EnsureMinSize() bool {
	if c.impl.Content() == nil {
		return false
	}
	var lastParent gui.CanvasObject

	windowNeedsMinSizeUpdate := false
	csize := c.impl.Size()
	min := c.impl.MinSize()

	ensureMinSize := func(node *RenderCacheNode) {
		obj := node.obj
		cache.SetCanvasForObject(obj, c.impl)

		if !obj.Visible() {
			return
		}
		minSize := obj.MinSize()
		minSizeChanged := node.minSize != minSize
		if minSizeChanged {
			objToLayout := obj
			node.minSize = minSize
			if node.parent != nil {
				objToLayout = node.parent.obj
			} else {
				windowNeedsMinSizeUpdate = true
				size := obj.Size()
				expectedSize := minSize.Max(size)
				if expectedSize != size && size != csize {
					objToLayout = nil
					obj.Resize(expectedSize)
				}
			}

			if objToLayout != lastParent {
				updateLayout(lastParent)
				lastParent = objToLayout
			}
		}
	}
	c.WalkTrees(nil, ensureMinSize)

	shouldResize := windowNeedsMinSizeUpdate && (csize.Width < min.Width || csize.Height < min.Height)
	if shouldResize {
		c.impl.Resize(csize.Max(min))
	}

	if lastParent != nil {
		c.RLock()
		updateLayout(lastParent)
		c.RUnlock()
	}
	return windowNeedsMinSizeUpdate
}

// Focus makes the provided item focused.
func (c *Canvas) Focus(obj gui.Focusable) {
	focusMgr := c.focusManager()
	if focusMgr != nil && focusMgr.Focus(obj) { // fast path – probably >99.9% of all cases
		if c.OnFocus != nil {
			c.OnFocus(obj)
		}
		return
	}

	c.RLock()
	focusMgrs := append([]*app.FocusManager{c.contentFocusMgr, c.menuFocusMgr}, c.overlays.ListFocusManagers()...)
	c.RUnlock()

	for _, mgr := range focusMgrs {
		if mgr == nil {
			continue
		}
		if focusMgr != mgr {
			if mgr.Focus(obj) {
				if c.OnFocus != nil {
					c.OnFocus(obj)
				}
				return
			}
		}
	}

	gui.LogError("Failed to focus object which is not part of the canvas’ content, menu or overlays.", nil)
}

// Focused returns the current focused object.
func (c *Canvas) Focused() gui.Focusable {
	mgr := c.focusManager()
	if mgr == nil {
		return nil
	}
	return mgr.Focused()
}

// FocusGained signals to the manager that its content got focus.
// Valid only on Desktop.
func (c *Canvas) FocusGained() {
	mgr := c.focusManager()
	if mgr == nil {
		return
	}
	mgr.FocusGained()
}

// FocusLost signals to the manager that its content lost focus.
// Valid only on Desktop.
func (c *Canvas) FocusLost() {
	mgr := c.focusManager()
	if mgr == nil {
		return
	}
	mgr.FocusLost()
}

// FocusNext focuses the next focusable item.
func (c *Canvas) FocusNext() {
	mgr := c.focusManager()
	if mgr == nil {
		return
	}
	mgr.FocusNext()
}

// FocusPrevious focuses the previous focusable item.
func (c *Canvas) FocusPrevious() {
	mgr := c.focusManager()
	if mgr == nil {
		return
	}
	mgr.FocusPrevious()
}

// FreeDirtyTextures frees dirty textures and returns the number of freed textures.
func (c *Canvas) FreeDirtyTextures() uint64 {
	freed := uint64(0)

	// Within a frame, refresh tasks are requested from the Refresh method,
	// and we desire to process all requested operations as much as possible
	// in a frame. Use a counter to guarantee that all desired tasks are
	// processed.
	for atomic.LoadUint32(&c.refreshCount) > 0 {
		var object gui.CanvasObject
		select {
		case object = <-c.refreshQueue.Out():
		default:
			// If refreshCount is positive but we cannot receive any object
			// from the refreshQueue, this means that the refresh task is
			// not yet ready to receive, continue until we can receive it.
			// Furthermore, we use Gosched to avoid CPU spin.
			runtime.Gosched()
			continue
		}
		atomic.AddUint32(&c.refreshCount, ^uint32(0))
		freed++
		freeWalked := func(obj gui.CanvasObject, _ gui.Position, _ gui.Position, _ gui.Size) bool {
			if c.painter != nil {
				c.painter.Free(obj)
			}
			return false
		}
		driver.WalkCompleteObjectTree(object, freeWalked, nil)
	}

	cache.RangeExpiredTexturesFor(c.impl, func(obj gui.CanvasObject) {
		if c.painter != nil {
			c.painter.Free(obj)
		}
	})
	return freed
}

// Initialize initializes the canvas.
func (c *Canvas) Initialize(impl SizeableCanvas, onOverlayChanged func()) {
	c.impl = impl
	c.refreshQueue = async.NewUnboundedCanvasObjectChan()
	c.overlays = &overlayStack{
		OverlayStack: internal.OverlayStack{
			OnChange: onOverlayChanged,
			Canvas:   impl,
		},
	}
}

// ObjectTrees return canvas object trees.
//
// This function uses lock.
func (c *Canvas) ObjectTrees() []gui.CanvasObject {
	c.RLock()
	var content, menu gui.CanvasObject
	if c.contentTree != nil && c.contentTree.root != nil {
		content = c.contentTree.root.obj
	}
	if c.menuTree != nil && c.menuTree.root != nil {
		menu = c.menuTree.root.obj
	}
	c.RUnlock()
	trees := make([]gui.CanvasObject, 0, len(c.Overlays().List())+2)
	trees = append(trees, content)
	if menu != nil {
		trees = append(trees, menu)
	}
	trees = append(trees, c.Overlays().List()...)
	return trees
}

// Overlays returns the overlay stack.
func (c *Canvas) Overlays() gui.OverlayStack {
	// we don't need to lock here, because overlays never changes
	return c.overlays
}

// Painter returns the canvas painter.
func (c *Canvas) Painter() gl.Painter {
	return c.painter
}

// Refresh refreshes a canvas object.
func (c *Canvas) Refresh(obj gui.CanvasObject) {
	atomic.AddUint32(&c.refreshCount, 1)
	c.refreshQueue.In() <- obj // never block
	c.SetDirty(true)
}

// RemoveShortcut removes a shortcut from the canvas.
func (c *Canvas) RemoveShortcut(shortcut gui.Shortcut) {
	c.shortcut.RemoveShortcut(shortcut)
}

// SetContentTreeAndFocusMgr sets content tree and focus manager.
//
// This function does not use the canvas lock.
func (c *Canvas) SetContentTreeAndFocusMgr(content gui.CanvasObject) {
	c.contentTree = &renderCacheTree{root: &RenderCacheNode{obj: content}}
	var focused gui.Focusable
	if c.contentFocusMgr != nil {
		focused = c.contentFocusMgr.Focused() // keep old focus if possible
	}
	c.contentFocusMgr = app.NewFocusManager(content)
	if focused != nil {
		c.contentFocusMgr.Focus(focused)
	}
}

const (
	dirtyTrue  = 1
	dirtyFalse = 0
)

// IsDirty checks if the canvas is dirty.
func (c *Canvas) IsDirty() bool {
	return atomic.LoadUint32(&c.dirty) == dirtyTrue
}

// SetDirty sets canvas dirty flag.
func (c *Canvas) SetDirty(dirty bool) {
	if dirty {
		atomic.StoreUint32(&c.dirty, dirtyTrue)
	} else {
		atomic.StoreUint32(&c.dirty, dirtyFalse)
	}
}

// SetMenuTreeAndFocusMgr sets menu tree and focus manager.
//
// This function does not use the canvas lock.
func (c *Canvas) SetMenuTreeAndFocusMgr(menu gui.CanvasObject) {
	c.menuTree = &renderCacheTree{root: &RenderCacheNode{obj: menu}}
	if menu != nil {
		c.menuFocusMgr = app.NewFocusManager(menu)
	} else {
		c.menuFocusMgr = nil
	}
}

// SetMobileWindowHeadTree sets window head tree.
//
// This function does not use the canvas lock.
func (c *Canvas) SetMobileWindowHeadTree(head gui.CanvasObject) {
	c.mWindowHeadTree = &renderCacheTree{root: &RenderCacheNode{obj: head}}
}

// SetPainter sets the canvas painter.
func (c *Canvas) SetPainter(p gl.Painter) {
	c.painter = p
}

// TypedShortcut handle the registered shortcut.
func (c *Canvas) TypedShortcut(shortcut gui.Shortcut) {
	c.shortcut.TypedShortcut(shortcut)
}

// Unfocus unfocuses all the objects in the canvas.
func (c *Canvas) Unfocus() {
	mgr := c.focusManager()
	if mgr == nil {
		return
	}
	if mgr.Focus(nil) && c.OnUnfocus != nil {
		c.OnUnfocus()
	}
}

// WalkTrees walks over the trees.
func (c *Canvas) WalkTrees(
	beforeChildren func(*RenderCacheNode, gui.Position),
	afterChildren func(*RenderCacheNode),
) {
	c.walkTree(c.contentTree, beforeChildren, afterChildren)
	if c.mWindowHeadTree != nil && c.mWindowHeadTree.root.obj != nil {
		c.walkTree(c.mWindowHeadTree, beforeChildren, afterChildren)
	}
	if c.menuTree != nil && c.menuTree.root.obj != nil {
		c.walkTree(c.menuTree, beforeChildren, afterChildren)
	}
	for _, tree := range c.overlays.renderCaches {
		if tree != nil {
			c.walkTree(tree, beforeChildren, afterChildren)
		}
	}
}

func (c *Canvas) focusManager() *app.FocusManager {
	if focusMgr := c.overlays.TopFocusManager(); focusMgr != nil {
		return focusMgr
	}
	c.RLock()
	defer c.RUnlock()
	if c.isMenuActive() {
		return c.menuFocusMgr
	}
	return c.contentFocusMgr
}

func (c *Canvas) isMenuActive() bool {
	if c.menuTree == nil || c.menuTree.root == nil || c.menuTree.root.obj == nil {
		return false
	}
	menu := c.menuTree.root.obj
	if am, ok := menu.(activatableMenu); ok {
		return am.IsActive()
	}
	return true
}

func (c *Canvas) walkTree(
	tree *renderCacheTree,
	beforeChildren func(*RenderCacheNode, gui.Position),
	afterChildren func(*RenderCacheNode),
) {
	tree.Lock()
	defer tree.Unlock()
	var node, parent, prev *RenderCacheNode
	node = tree.root

	bc := func(obj gui.CanvasObject, pos gui.Position, _ gui.Position, _ gui.Size) bool {
		if node != nil && node.obj != obj {
			if parent.firstChild == node {
				parent.firstChild = nil
			}
			node = nil
		}
		if node == nil {
			node = &RenderCacheNode{parent: parent, obj: obj}
			if parent.firstChild == nil {
				parent.firstChild = node
			} else {
				prev.nextSibling = node
			}
		}
		if prev != nil && prev.parent != parent {
			prev = nil
		}

		if beforeChildren != nil {
			beforeChildren(node, pos)
		}

		parent = node
		node = parent.firstChild
		return false
	}
	ac := func(obj gui.CanvasObject, _ gui.CanvasObject) {
		node = parent
		parent = node.parent
		if prev != nil && prev.parent != parent {
			prev.nextSibling = nil
		}

		if afterChildren != nil {
			afterChildren(node)
		}

		prev = node
		node = node.nextSibling
	}
	driver.WalkVisibleObjectTree(tree.root.obj, bc, ac)
}

// RenderCacheNode represents a node in a render cache tree.
type RenderCacheNode struct {
	// structural data
	firstChild  *RenderCacheNode
	nextSibling *RenderCacheNode
	obj         gui.CanvasObject
	parent      *RenderCacheNode
	// cache data
	minSize gui.Size
	// painterData is some data from the painter associated with the drawed node
	// it may for instance point to a GL texture
	// it should free all associated resources when released
	// i.e. it should not simply be a texture reference integer
	painterData interface{}
}

// Obj returns the node object.
func (r *RenderCacheNode) Obj() gui.CanvasObject {
	return r.obj
}

type activatableMenu interface {
	IsActive() bool
}

type overlayStack struct {
	internal.OverlayStack

	propertyLock sync.RWMutex
	renderCaches []*renderCacheTree
}

func (o *overlayStack) Add(overlay gui.CanvasObject) {
	if overlay == nil {
		return
	}
	o.propertyLock.Lock()
	defer o.propertyLock.Unlock()
	o.add(overlay)
}

func (o *overlayStack) Remove(overlay gui.CanvasObject) {
	if overlay == nil || len(o.List()) == 0 {
		return
	}
	o.propertyLock.Lock()
	defer o.propertyLock.Unlock()
	o.remove(overlay)
}

func (o *overlayStack) add(overlay gui.CanvasObject) {
	o.renderCaches = append(o.renderCaches, &renderCacheTree{root: &RenderCacheNode{obj: overlay}})
	o.OverlayStack.Add(overlay)
}

func (o *overlayStack) remove(overlay gui.CanvasObject) {
	o.OverlayStack.Remove(overlay)
	overlayCount := len(o.List())
	o.renderCaches = o.renderCaches[:overlayCount]
}

type renderCacheTree struct {
	sync.RWMutex
	root *RenderCacheNode
}

func updateLayout(objToLayout gui.CanvasObject) {
	switch cont := objToLayout.(type) {
	case *gui.Container:
		if cont.Layout != nil {
			cont.Layout.Layout(cont.Objects, cont.Size())
		}
	case gui.Widget:
		cache.Renderer(cont).Layout(cont.Size())
	}
}
