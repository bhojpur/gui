package container

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
	"sync"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/driver/desktop"
	"github.com/bhojpur/gui/pkg/engine/internal"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

// TabItem represents a single view in a tab view.
// The Text and Icon are used for the tab button and the Content is shown when the corresponding tab is active.
//
// Since: 1.4
type TabItem struct {
	Text    string
	Icon    gui.Resource
	Content gui.CanvasObject
}

// TabLocation is the location where the tabs of a tab container should be rendered
//
// Since: 1.4
type TabLocation int

// TabLocation values
const (
	TabLocationTop TabLocation = iota
	TabLocationLeading
	TabLocationBottom
	TabLocationTrailing
)

// NewTabItem creates a new item for a tabbed widget - each item specifies the content and a label for its tab.
//
// Since: 1.4
func NewTabItem(text string, content gui.CanvasObject) *TabItem {
	return &TabItem{Text: text, Content: content}
}

// NewTabItemWithIcon creates a new item for a tabbed widget - each item specifies the content and a label with an icon for its tab.
//
// Since: 1.4
func NewTabItemWithIcon(text string, icon gui.Resource, content gui.CanvasObject) *TabItem {
	return &TabItem{Text: text, Icon: icon, Content: content}
}

type baseTabs interface {
	onUnselected() func(*TabItem)
	onSelected() func(*TabItem)

	items() []*TabItem
	setItems([]*TabItem)

	selected() int
	setSelected(int)

	tabLocation() TabLocation

	transitioning() bool
	setTransitioning(bool)
}

func tabsAdjustedLocation(l TabLocation) TabLocation {
	// Mobile has limited screen space, so don't put app tab bar on long edges
	if d := gui.CurrentDevice(); d.IsMobile() {
		if o := d.Orientation(); gui.IsVertical(o) {
			if l == TabLocationLeading {
				return TabLocationTop
			} else if l == TabLocationTrailing {
				return TabLocationBottom
			}
		} else {
			if l == TabLocationTop {
				return TabLocationLeading
			} else if l == TabLocationBottom {
				return TabLocationTrailing
			}
		}
	}

	return l
}

func buildPopUpMenu(t baseTabs, button *widget.Button, items []*gui.MenuItem) *widget.PopUpMenu {
	d := gui.CurrentApp().Driver()
	c := d.CanvasForObject(button)
	popUpMenu := widget.NewPopUpMenu(gui.NewMenu("", items...), c)
	buttonPos := d.AbsolutePositionForObject(button)
	buttonSize := button.Size()
	popUpMin := popUpMenu.MinSize()
	var popUpPos gui.Position
	switch t.tabLocation() {
	case TabLocationLeading:
		popUpPos.X = buttonPos.X + buttonSize.Width
		popUpPos.Y = buttonPos.Y + buttonSize.Height - popUpMin.Height
	case TabLocationTrailing:
		popUpPos.X = buttonPos.X - popUpMin.Width
		popUpPos.Y = buttonPos.Y + buttonSize.Height - popUpMin.Height
	case TabLocationTop:
		popUpPos.X = buttonPos.X + buttonSize.Width - popUpMin.Width
		popUpPos.Y = buttonPos.Y + buttonSize.Height
	case TabLocationBottom:
		popUpPos.X = buttonPos.X + buttonSize.Width - popUpMin.Width
		popUpPos.Y = buttonPos.Y - popUpMin.Height
	}
	popUpMenu.ShowAtPosition(popUpPos)
	return popUpMenu
}

func removeIndex(t baseTabs, index int) {
	items := t.items()
	if index < 0 || index >= len(items) {
		return
	}
	setItems(t, append(items[:index], items[index+1:]...))
	if s := t.selected(); index < s {
		t.setSelected(s - 1)
	}
}

func removeItem(t baseTabs, item *TabItem) {
	for index, existingItem := range t.items() {
		if existingItem == item {
			removeIndex(t, index)
			break
		}
	}
}

func selected(t baseTabs) *TabItem {
	selected := t.selected()
	items := t.items()
	if selected < 0 || selected >= len(items) {
		return nil
	}
	return items[selected]
}

func selectIndex(t baseTabs, index int) {
	selected := t.selected()

	if selected == index {
		// No change, so do nothing
		return
	}

	items := t.items()

	if f := t.onUnselected(); f != nil && selected >= 0 && selected < len(items) {
		// Notification of unselected
		f(items[selected])
	}

	if index < 0 || index >= len(items) {
		// Out of bounds, so do nothing
		return
	}

	t.setTransitioning(true)
	t.setSelected(index)

	if f := t.onSelected(); f != nil {
		// Notification of selected
		f(items[index])
	}
}

func selectItem(t baseTabs, item *TabItem) {
	for i, child := range t.items() {
		if child == item {
			selectIndex(t, i)
			return
		}
	}
}

func setItems(t baseTabs, items []*TabItem) {
	if mismatchedTabItems(items) {
		internal.LogHint("Tab items should all have the same type of content (text, icons or both)")
	}
	t.setItems(items)
	selected := t.selected()
	count := len(items)
	switch {
	case count == 0:
		// No items available to be selected
		selectIndex(t, -1) // Unsure OnUnselected gets called if applicable
		t.setSelected(-1)
	case selected < 0:
		// Current is first tab item
		selectIndex(t, 0)
	case selected >= count:
		// Current doesn't exist, select last tab
		selectIndex(t, count-1)
	}
}

type baseTabsRenderer struct {
	positionAnimation, sizeAnimation *gui.Animation

	lastIndicatorMutex sync.RWMutex
	lastIndicatorPos   gui.Position
	lastIndicatorSize  gui.Size

	action             *widget.Button
	bar                *gui.Container
	divider, indicator *canvas.Rectangle

	buttonCache map[*TabItem]*tabButton
}

func (r *baseTabsRenderer) Destroy() {
}

func (r *baseTabsRenderer) applyTheme(t baseTabs) {
	if r.action != nil {
		r.action.SetIcon(moreIcon(t))
	}
	r.divider.FillColor = theme.ShadowColor()
	r.indicator.FillColor = theme.PrimaryColor()
}

func (r *baseTabsRenderer) layout(t baseTabs, size gui.Size) {
	var (
		barPos, dividerPos, contentPos    gui.Position
		barSize, dividerSize, contentSize gui.Size
	)

	barMin := r.bar.MinSize()

	switch t.tabLocation() {
	case TabLocationTop:
		barHeight := barMin.Height
		barPos = gui.NewPos(0, 0)
		barSize = gui.NewSize(size.Width, barHeight)
		dividerPos = gui.NewPos(0, barHeight)
		dividerSize = gui.NewSize(size.Width, theme.Padding())
		contentPos = gui.NewPos(0, barHeight+theme.Padding())
		contentSize = gui.NewSize(size.Width, size.Height-barHeight-theme.Padding())
	case TabLocationLeading:
		barWidth := barMin.Width
		barPos = gui.NewPos(0, 0)
		barSize = gui.NewSize(barWidth, size.Height)
		dividerPos = gui.NewPos(barWidth, 0)
		dividerSize = gui.NewSize(theme.Padding(), size.Height)
		contentPos = gui.NewPos(barWidth+theme.Padding(), 0)
		contentSize = gui.NewSize(size.Width-barWidth-theme.Padding(), size.Height)
	case TabLocationBottom:
		barHeight := barMin.Height
		barPos = gui.NewPos(0, size.Height-barHeight)
		barSize = gui.NewSize(size.Width, barHeight)
		dividerPos = gui.NewPos(0, size.Height-barHeight-theme.Padding())
		dividerSize = gui.NewSize(size.Width, theme.Padding())
		contentPos = gui.NewPos(0, 0)
		contentSize = gui.NewSize(size.Width, size.Height-barHeight-theme.Padding())
	case TabLocationTrailing:
		barWidth := barMin.Width
		barPos = gui.NewPos(size.Width-barWidth, 0)
		barSize = gui.NewSize(barWidth, size.Height)
		dividerPos = gui.NewPos(size.Width-barWidth-theme.Padding(), 0)
		dividerSize = gui.NewSize(theme.Padding(), size.Height)
		contentPos = gui.NewPos(0, 0)
		contentSize = gui.NewSize(size.Width-barWidth-theme.Padding(), size.Height)
	}

	r.bar.Move(barPos)
	r.bar.Resize(barSize)
	r.divider.Move(dividerPos)
	r.divider.Resize(dividerSize)
	selected := t.selected()
	for i, ti := range t.items() {
		if i == selected {
			ti.Content.Move(contentPos)
			ti.Content.Resize(contentSize)
			ti.Content.Show()
		} else {
			ti.Content.Hide()
		}
	}
}

func (r *baseTabsRenderer) minSize(t baseTabs) gui.Size {
	barMin := r.bar.MinSize()

	contentMin := gui.NewSize(0, 0)
	for _, content := range t.items() {
		contentMin = contentMin.Max(content.Content.MinSize())
	}

	switch t.tabLocation() {
	case TabLocationLeading, TabLocationTrailing:
		return gui.NewSize(barMin.Width+contentMin.Width+theme.Padding(), contentMin.Height)
	default:
		return gui.NewSize(contentMin.Width, barMin.Height+contentMin.Height+theme.Padding())
	}
}

func (r *baseTabsRenderer) moveIndicator(pos gui.Position, siz gui.Size, animate bool) {
	r.lastIndicatorMutex.RLock()
	isSameState := r.lastIndicatorPos.Subtract(pos).IsZero() && r.lastIndicatorSize.Subtract(siz).IsZero()
	r.lastIndicatorMutex.RUnlock()
	if isSameState {
		return
	}

	if r.positionAnimation != nil {
		r.positionAnimation.Stop()
		r.positionAnimation = nil
	}
	if r.sizeAnimation != nil {
		r.sizeAnimation.Stop()
		r.sizeAnimation = nil
	}

	r.indicator.Show()
	if r.indicator.Position().IsZero() {
		r.indicator.Move(pos)
		r.indicator.Resize(siz)
		r.indicator.Refresh()
		return
	}

	r.lastIndicatorMutex.Lock()
	r.lastIndicatorPos = pos
	r.lastIndicatorSize = siz
	r.lastIndicatorMutex.Unlock()

	if animate {
		r.positionAnimation = canvas.NewPositionAnimation(r.indicator.Position(), pos, canvas.DurationShort, func(p gui.Position) {
			r.indicator.Move(p)
			r.indicator.Refresh()
			if pos == p {
				r.positionAnimation.Stop()
				r.positionAnimation = nil
			}
		})
		r.sizeAnimation = canvas.NewSizeAnimation(r.indicator.Size(), siz, canvas.DurationShort, func(s gui.Size) {
			r.indicator.Resize(s)
			r.indicator.Refresh()
			if siz == s {
				r.sizeAnimation.Stop()
				r.sizeAnimation = nil
			}
		})

		r.positionAnimation.Start()
		r.sizeAnimation.Start()
	} else {
		r.indicator.Move(pos)
		r.indicator.Resize(siz)
		r.indicator.Refresh()
	}
}

func (r *baseTabsRenderer) objects(t baseTabs) []gui.CanvasObject {
	objects := []gui.CanvasObject{r.bar, r.divider, r.indicator}
	if i, is := t.selected(), t.items(); i >= 0 && i < len(is) {
		objects = append(objects, is[i].Content)
	}
	return objects
}

func (r *baseTabsRenderer) refresh(t baseTabs) {
	r.applyTheme(t)

	r.bar.Refresh()
	r.divider.Refresh()
	r.indicator.Refresh()
}

type buttonIconPosition int

const (
	buttonIconInline buttonIconPosition = iota
	buttonIconTop
)

var _ gui.Widget = (*tabButton)(nil)
var _ gui.Tappable = (*tabButton)(nil)
var _ desktop.Hoverable = (*tabButton)(nil)

type tabButton struct {
	widget.BaseWidget
	hovered       bool
	icon          gui.Resource
	iconPosition  buttonIconPosition
	importance    widget.ButtonImportance
	onTapped      func()
	onClosed      func()
	text          string
	textAlignment gui.TextAlign
}

func (b *tabButton) CreateRenderer() gui.WidgetRenderer {
	b.ExtendBaseWidget(b)
	background := canvas.NewRectangle(theme.HoverColor())
	background.Hide()
	var icon *canvas.Image
	if b.icon != nil {
		icon = canvas.NewImageFromResource(b.icon)
	}

	label := canvas.NewText(b.text, theme.TextColor())
	label.TextStyle.Bold = true

	close := &tabCloseButton{
		parent: b,
		onTapped: func() {
			if f := b.onClosed; f != nil {
				f()
			}
		},
	}
	close.ExtendBaseWidget(close)
	close.Hide()

	objects := []gui.CanvasObject{background, label, close}
	if icon != nil {
		objects = append(objects, icon)
	}

	r := &tabButtonRenderer{
		button:     b,
		background: background,
		icon:       icon,
		label:      label,
		close:      close,
		objects:    objects,
	}
	r.Refresh()
	return r
}

func (b *tabButton) MinSize() gui.Size {
	b.ExtendBaseWidget(b)
	return b.BaseWidget.MinSize()
}

func (b *tabButton) MouseIn(*desktop.MouseEvent) {
	b.hovered = true
	b.Refresh()
}

func (b *tabButton) MouseMoved(*desktop.MouseEvent) {
}

func (b *tabButton) MouseOut() {
	b.hovered = false
	b.Refresh()
}

func (b *tabButton) Tapped(*gui.PointEvent) {
	b.onTapped()
}

type tabButtonRenderer struct {
	button     *tabButton
	background *canvas.Rectangle
	icon       *canvas.Image
	label      *canvas.Text
	close      *tabCloseButton
	objects    []gui.CanvasObject
}

func (r *tabButtonRenderer) Destroy() {
}

func (r *tabButtonRenderer) Layout(size gui.Size) {
	r.background.Resize(size)
	padding := r.padding()
	innerSize := size.Subtract(padding)
	innerOffset := gui.NewPos(padding.Width/2, padding.Height/2)
	labelShift := float32(0)
	if r.icon != nil {
		var iconOffset gui.Position
		if r.button.iconPosition == buttonIconTop {
			iconOffset = gui.NewPos((innerSize.Width-r.iconSize())/2, 0)
		} else {
			iconOffset = gui.NewPos(0, (innerSize.Height-r.iconSize())/2)
		}
		r.icon.Resize(gui.NewSize(r.iconSize(), r.iconSize()))
		r.icon.Move(innerOffset.Add(iconOffset))
		labelShift = r.iconSize() + theme.Padding()
	}
	if r.label.Text != "" {
		var labelOffset gui.Position
		var labelSize gui.Size
		if r.button.iconPosition == buttonIconTop {
			labelOffset = gui.NewPos(0, labelShift)
			labelSize = gui.NewSize(innerSize.Width, r.label.MinSize().Height)
		} else {
			labelOffset = gui.NewPos(labelShift, 0)
			labelSize = gui.NewSize(innerSize.Width-labelShift, innerSize.Height)
		}
		r.label.Resize(labelSize)
		r.label.Move(innerOffset.Add(labelOffset))
	}
	r.close.Move(gui.NewPos(size.Width-theme.IconInlineSize()-theme.Padding(), (size.Height-theme.IconInlineSize())/2))
	r.close.Resize(gui.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
}

func (r *tabButtonRenderer) MinSize() gui.Size {
	var contentWidth, contentHeight float32
	textSize := r.label.MinSize()
	if r.button.iconPosition == buttonIconTop {
		contentWidth = gui.Max(textSize.Width, r.iconSize())
		if r.icon != nil {
			contentHeight += r.iconSize()
		}
		if r.label.Text != "" {
			if r.icon != nil {
				contentHeight += theme.Padding()
			}
			contentHeight += textSize.Height
		}
	} else {
		contentHeight = gui.Max(textSize.Height, r.iconSize())
		if r.icon != nil {
			contentWidth += r.iconSize()
		}
		if r.label.Text != "" {
			if r.icon != nil {
				contentWidth += theme.Padding()
			}
			contentWidth += textSize.Width
		}
	}
	if r.button.onClosed != nil {
		contentWidth += theme.IconInlineSize() + theme.Padding()
		contentHeight = gui.Max(contentHeight, theme.IconInlineSize())
	}
	return gui.NewSize(contentWidth, contentHeight).Add(r.padding())
}

func (r *tabButtonRenderer) Objects() []gui.CanvasObject {
	return r.objects
}

func (r *tabButtonRenderer) Refresh() {
	if r.button.hovered {
		r.background.FillColor = theme.HoverColor()
		r.background.Show()
	} else {
		r.background.Hide()
	}
	r.background.Refresh()

	r.label.Text = r.button.text
	r.label.Alignment = r.button.textAlignment
	if r.button.importance == widget.HighImportance {
		r.label.Color = theme.PrimaryColor()
	} else {
		r.label.Color = theme.TextColor()
	}
	r.label.TextSize = theme.TextSize()
	if r.button.text == "" {
		r.label.Hide()
	} else {
		r.label.Show()
	}

	if r.icon != nil && r.icon.Resource != nil {
		switch res := r.icon.Resource.(type) {
		case *theme.ThemedResource:
			if r.button.importance == widget.HighImportance {
				r.icon.Resource = theme.NewPrimaryThemedResource(res)
				r.icon.Refresh()
			}
		case *theme.PrimaryThemedResource:
			if r.button.importance != widget.HighImportance {
				r.icon.Resource = res.Original()
				r.icon.Refresh()
			}
		}
	}

	if d := gui.CurrentDevice(); r.button.onClosed != nil && (d.IsMobile() || r.button.hovered || r.close.hovered) {
		r.close.Show()
	} else {
		r.close.Hide()
	}
	r.close.Refresh()

	canvas.Refresh(r.button)
}

func (r *tabButtonRenderer) iconSize() float32 {
	switch r.button.iconPosition {
	case buttonIconTop:
		return 2 * theme.IconInlineSize()
	default:
		return theme.IconInlineSize()
	}
}

func (r *tabButtonRenderer) padding() gui.Size {
	if r.label.Text != "" && r.button.iconPosition == buttonIconInline {
		return gui.NewSize(theme.Padding()*4, theme.Padding()*4)
	}
	return gui.NewSize(theme.Padding()*2, theme.Padding()*4)
}

var _ gui.Widget = (*tabCloseButton)(nil)
var _ gui.Tappable = (*tabCloseButton)(nil)
var _ desktop.Hoverable = (*tabCloseButton)(nil)

type tabCloseButton struct {
	widget.BaseWidget
	parent   *tabButton
	hovered  bool
	onTapped func()
}

func (b *tabCloseButton) CreateRenderer() gui.WidgetRenderer {
	b.ExtendBaseWidget(b)
	background := canvas.NewRectangle(theme.HoverColor())
	background.Hide()
	icon := canvas.NewImageFromResource(theme.CancelIcon())

	r := &tabCloseButtonRenderer{
		button:     b,
		background: background,
		icon:       icon,
		objects:    []gui.CanvasObject{background, icon},
	}
	r.Refresh()
	return r
}

func (b *tabCloseButton) MinSize() gui.Size {
	b.ExtendBaseWidget(b)
	return b.BaseWidget.MinSize()
}

func (b *tabCloseButton) MouseIn(*desktop.MouseEvent) {
	b.hovered = true
	b.parent.Refresh()
}

func (b *tabCloseButton) MouseMoved(*desktop.MouseEvent) {
}

func (b *tabCloseButton) MouseOut() {
	b.hovered = false
	b.parent.Refresh()
}

func (b *tabCloseButton) Tapped(*gui.PointEvent) {
	b.onTapped()
}

type tabCloseButtonRenderer struct {
	button     *tabCloseButton
	background *canvas.Rectangle
	icon       *canvas.Image
	objects    []gui.CanvasObject
}

func (r *tabCloseButtonRenderer) Destroy() {
}

func (r *tabCloseButtonRenderer) Layout(size gui.Size) {
	r.background.Resize(size)
	r.icon.Resize(size)
}

func (r *tabCloseButtonRenderer) MinSize() gui.Size {
	return gui.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
}

func (r *tabCloseButtonRenderer) Objects() []gui.CanvasObject {
	return r.objects
}

func (r *tabCloseButtonRenderer) Refresh() {
	if r.button.hovered {
		r.background.FillColor = theme.HoverColor()
		r.background.Show()
	} else {
		r.background.Hide()
	}
	r.background.Refresh()
	switch res := r.icon.Resource.(type) {
	case *theme.ThemedResource:
		if r.button.parent.importance == widget.HighImportance {
			r.icon.Resource = theme.NewPrimaryThemedResource(res)
		}
	case *theme.PrimaryThemedResource:
		if r.button.parent.importance != widget.HighImportance {
			r.icon.Resource = res.Original()
		}
	}
	r.icon.Refresh()
}

func mismatchedTabItems(items []*TabItem) bool {
	var hasText, hasIcon bool
	for _, tab := range items {
		hasText = hasText || tab.Text != ""
		hasIcon = hasIcon || tab.Icon != nil
	}

	mismatch := false
	for _, tab := range items {
		if (hasText && tab.Text == "") || (hasIcon && tab.Icon == nil) {
			mismatch = true
			break
		}
	}

	return mismatch
}

func moreIcon(t baseTabs) gui.Resource {
	if l := t.tabLocation(); l == TabLocationLeading || l == TabLocationTrailing {
		return theme.MoreVerticalIcon()
	}
	return theme.MoreHorizontalIcon()
}
