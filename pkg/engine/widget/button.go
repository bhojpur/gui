package widget

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
	"image/color"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/driver/desktop"
	col "github.com/bhojpur/gui/pkg/engine/internal/color"
	"github.com/bhojpur/gui/pkg/engine/internal/widget"
	"github.com/bhojpur/gui/pkg/engine/layout"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

// ButtonAlign represents the horizontal alignment of a button.
type ButtonAlign int

// ButtonIconPlacement represents the ordering of icon & text within a button.
type ButtonIconPlacement int

// ButtonImportance represents how prominent the button should appear
//
// Since: 1.4
type ButtonImportance int

// ButtonStyle determines the behaviour and rendering of a button.
type ButtonStyle int

const (
	// ButtonAlignCenter aligns the icon and the text centrally.
	ButtonAlignCenter ButtonAlign = iota
	// ButtonAlignLeading aligns the icon and the text with the leading edge.
	ButtonAlignLeading
	// ButtonAlignTrailing aligns the icon and the text with the trailing edge.
	ButtonAlignTrailing
)

const (
	// ButtonIconLeadingText aligns the icon on the leading edge of the text.
	ButtonIconLeadingText ButtonIconPlacement = iota
	// ButtonIconTrailingText aligns the icon on the trailing edge of the text.
	ButtonIconTrailingText
)

const (
	// MediumImportance applies a standard appearance.
	MediumImportance ButtonImportance = iota
	// HighImportance applies a prominent appearance.
	HighImportance
	// LowImportance applies a subtle appearance.
	LowImportance
)

var _ gui.Focusable = (*Button)(nil)

// Button widget has a text label and triggers an event func when clicked
type Button struct {
	DisableableWidget
	Text string
	Icon gui.Resource
	// Specify how prominent the button should be, High will highlight the button and Low will remove some decoration.
	//
	// Since: 1.4
	Importance    ButtonImportance
	Alignment     ButtonAlign
	IconPlacement ButtonIconPlacement

	OnTapped func() `json:"-"`

	hovered, focused bool
	tapAnim          *gui.Animation
}

// NewButton creates a new button widget with the set label and tap handler
func NewButton(label string, tapped func()) *Button {
	button := &Button{
		Text:     label,
		OnTapped: tapped,
	}

	button.ExtendBaseWidget(button)
	return button
}

// NewButtonWithIcon creates a new button widget with the specified label, themed icon and tap handler
func NewButtonWithIcon(label string, icon gui.Resource, tapped func()) *Button {
	button := &Button{
		Text:     label,
		Icon:     icon,
		OnTapped: tapped,
	}

	button.ExtendBaseWidget(button)
	return button
}

// CreateRenderer is a private method to Bhojpur GUI which links this widget to its renderer
func (b *Button) CreateRenderer() gui.WidgetRenderer {
	b.ExtendBaseWidget(b)
	seg := &TextSegment{Text: b.Text, Style: RichTextStyleStrong}
	seg.Style.Alignment = gui.TextAlignCenter
	text := NewRichText(seg)
	text.inset = gui.NewSize(theme.Padding()*2, theme.Padding()*2)

	background := canvas.NewRectangle(theme.ButtonColor())
	tapBG := canvas.NewRectangle(color.Transparent)
	b.tapAnim = newButtonTapAnimation(tapBG, b)
	b.tapAnim.Curve = gui.AnimationEaseOut
	objects := []gui.CanvasObject{
		background,
		tapBG,
		text,
	}
	shadowLevel := widget.ButtonLevel
	if b.Importance == LowImportance {
		shadowLevel = widget.BaseLevel
	}
	r := &buttonRenderer{
		ShadowingRenderer: widget.NewShadowingRenderer(objects, shadowLevel),
		background:        background,
		tapBG:             tapBG,
		button:            b,
		label:             text,
		layout:            layout.NewHBoxLayout(),
	}
	r.updateIconAndText()
	r.applyTheme()
	return r
}

// Cursor returns the cursor type of this widget
func (b *Button) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

// FocusGained is a hook called by the focus handling logic after this object gained the focus.
func (b *Button) FocusGained() {
	b.focused = true
	b.Refresh()
}

// FocusLost is a hook called by the focus handling logic after this object lost the focus.
func (b *Button) FocusLost() {
	b.focused = false
	b.Refresh()
}

// MinSize returns the size that this widget should not shrink below
func (b *Button) MinSize() gui.Size {
	b.ExtendBaseWidget(b)
	return b.BaseWidget.MinSize()
}

// MouseIn is called when a desktop pointer enters the widget
func (b *Button) MouseIn(*desktop.MouseEvent) {
	b.hovered = true
	b.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (b *Button) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut is called when a desktop pointer exits the widget
func (b *Button) MouseOut() {
	b.hovered = false
	b.Refresh()
}

// SetIcon updates the icon on a label - pass nil to hide an icon
func (b *Button) SetIcon(icon gui.Resource) {
	b.Icon = icon

	b.Refresh()
}

// SetText allows the button label to be changed
func (b *Button) SetText(text string) {
	b.Text = text

	b.Refresh()
}

// Tapped is called when a pointer tapped event is captured and triggers any tap handler
func (b *Button) Tapped(*gui.PointEvent) {
	if b.Disabled() {
		return
	}

	b.tapAnimation()
	b.Refresh()

	if b.OnTapped != nil {
		b.OnTapped()
	}
}

// TypedRune is a hook called by the input handling logic on text input events if this object is focused.
func (b *Button) TypedRune(rune) {
}

// TypedKey is a hook called by the input handling logic on key events if this object is focused.
func (b *Button) TypedKey(ev *gui.KeyEvent) {
	if ev.Name == gui.KeySpace {
		b.Tapped(nil)
	}
}

func (b *Button) tapAnimation() {
	if b.tapAnim == nil {
		return
	}
	b.tapAnim.Stop()
	b.tapAnim.Start()
}

type buttonRenderer struct {
	*widget.ShadowingRenderer

	icon       *canvas.Image
	label      *RichText
	background *canvas.Rectangle
	tapBG      *canvas.Rectangle
	button     *Button
	layout     gui.Layout
}

// Layout the components of the button widget
func (r *buttonRenderer) Layout(size gui.Size) {
	var inset gui.Position
	bgSize := size
	if r.button.Importance != LowImportance {
		inset = gui.NewPos(theme.Padding()/2, theme.Padding()/2)
		bgSize = size.Subtract(gui.NewSize(theme.Padding(), theme.Padding()))
	}
	r.LayoutShadow(bgSize, inset)

	r.background.Move(inset)
	r.background.Resize(bgSize)

	hasIcon := r.icon != nil
	hasLabel := r.label.Segments[0].(*TextSegment).Text != ""
	if !hasIcon && !hasLabel {
		// Nothing to layout
		return
	}
	iconSize := gui.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
	labelSize := r.label.MinSize()
	padding := r.padding()
	if hasLabel {
		if hasIcon {
			// Both
			var objects []gui.CanvasObject
			if r.button.IconPlacement == ButtonIconLeadingText {
				objects = append(objects, r.icon, r.label)
			} else {
				objects = append(objects, r.label, r.icon)
			}
			r.icon.SetMinSize(iconSize)
			min := r.layout.MinSize(objects)
			r.layout.Layout(objects, min)
			pos := alignedPosition(r.button.Alignment, padding, min, size)
			r.label.Move(r.label.Position().Add(pos))
			r.icon.Move(r.icon.Position().Add(pos))
		} else {
			// Label Only
			r.label.Move(alignedPosition(r.button.Alignment, padding, labelSize, size))
			r.label.Resize(labelSize)
		}
	} else {
		// Icon Only
		r.icon.Move(alignedPosition(r.button.Alignment, padding, iconSize, size))
		r.icon.Resize(iconSize)
	}
}

// MinSize calculates the minimum size of a button.
// This is based on the contained text, any icon that is set and a standard
// amount of padding added.
func (r *buttonRenderer) MinSize() (size gui.Size) {
	hasIcon := r.icon != nil
	hasLabel := r.label.Segments[0].(*TextSegment).Text != ""
	iconSize := gui.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
	labelSize := r.label.MinSize()
	if hasLabel {
		size.Width = labelSize.Width
	}
	if hasIcon {
		if hasLabel {
			size.Width += theme.Padding()
		}
		size.Width += iconSize.Width
	}
	size.Height = gui.Max(labelSize.Height, iconSize.Height)
	size = size.Add(r.padding())
	return
}

func (r *buttonRenderer) Refresh() {
	r.label.inset = gui.NewSize(theme.Padding()*2, theme.Padding()*2)
	r.label.Segments[0].(*TextSegment).Text = r.button.Text
	r.updateIconAndText()
	r.applyTheme()
	r.background.Refresh()
	r.Layout(r.button.Size())
	canvas.Refresh(r.button.super())
}

// applyTheme updates this button to match the current theme
func (r *buttonRenderer) applyTheme() {
	r.background.FillColor = r.buttonColor()
	r.label.Segments[0].(*TextSegment).Style.ColorName = theme.ColorNameForeground
	switch {
	case r.button.disabled:
		r.label.Segments[0].(*TextSegment).Style.ColorName = theme.ColorNameDisabled
	case r.button.Importance == HighImportance:
		r.label.Segments[0].(*TextSegment).Style.ColorName = theme.ColorNameBackground
	}
	r.label.Refresh()
	if r.icon != nil && r.icon.Resource != nil {
		switch res := r.icon.Resource.(type) {
		case *theme.ThemedResource:
			if r.button.Importance == HighImportance {
				r.icon.Resource = theme.NewInvertedThemedResource(res)
				r.icon.Refresh()
			}
		case *theme.InvertedThemedResource:
			if r.button.Importance != HighImportance {
				r.icon.Resource = res.Original()
				r.icon.Refresh()
			}
		}
	}
	r.ShadowingRenderer.RefreshShadow()
}

func (r *buttonRenderer) buttonColor() color.Color {
	switch {
	case r.button.Disabled():
		return theme.DisabledButtonColor()
	case r.button.focused:
		return blendColor(theme.ButtonColor(), theme.FocusColor())
	case r.button.hovered:
		bg := theme.ButtonColor()
		if r.button.Importance == HighImportance {
			bg = theme.PrimaryColor()
		}

		return blendColor(bg, theme.HoverColor())
	case r.button.Importance == HighImportance:
		return theme.PrimaryColor()
	default:
		return theme.ButtonColor()
	}
}

func (r *buttonRenderer) padding() gui.Size {
	if r.button.Text == "" {
		return gui.NewSize(theme.Padding()*4, theme.Padding()*4)
	}
	return gui.NewSize(theme.Padding()*6, theme.Padding()*4)
}

func (r *buttonRenderer) updateIconAndText() {
	if r.button.Icon != nil && r.button.Visible() {
		if r.icon == nil {
			r.icon = canvas.NewImageFromResource(r.button.Icon)
			r.icon.FillMode = canvas.ImageFillContain
			r.SetObjects([]gui.CanvasObject{r.background, r.tapBG, r.label, r.icon})
		}
		if r.button.Disabled() {
			r.icon.Resource = theme.NewDisabledResource(r.button.Icon)
		} else {
			r.icon.Resource = r.button.Icon
		}
		r.icon.Refresh()
		r.icon.Show()
	} else if r.icon != nil {
		r.icon.Hide()
	}
	if r.button.Text == "" {
		r.label.Hide()
	} else {
		r.label.Show()
	}
	r.label.Refresh()
}

func alignedPosition(align ButtonAlign, padding, objectSize, layoutSize gui.Size) (pos gui.Position) {
	pos.Y = (layoutSize.Height - objectSize.Height) / 2
	switch align {
	case ButtonAlignCenter:
		pos.X = (layoutSize.Width - objectSize.Width) / 2
	case ButtonAlignLeading:
		pos.X = padding.Width / 2
	case ButtonAlignTrailing:
		pos.X = layoutSize.Width - objectSize.Width - padding.Width/2
	}
	return
}

func blendColor(under, over color.Color) color.Color {
	// This alpha blends with the over operator, and accounts for RGBA() returning alpha-premultiplied values
	dstR, dstG, dstB, dstA := under.RGBA()
	srcR, srcG, srcB, srcA := over.RGBA()

	srcAlpha := float32(srcA) / 0xFFFF
	dstAlpha := float32(dstA) / 0xFFFF

	outAlpha := srcAlpha + dstAlpha*(1-srcAlpha)
	outR := srcR + uint32(float32(dstR)*(1-srcAlpha))
	outG := srcG + uint32(float32(dstG)*(1-srcAlpha))
	outB := srcB + uint32(float32(dstB)*(1-srcAlpha))
	// We create an RGBA64 here because the color components are already alpha-premultiplied 16-bit values (they're just stored in uint32s).
	return color.RGBA64{R: uint16(outR), G: uint16(outG), B: uint16(outB), A: uint16(outAlpha * 0xFFFF)}

}

func newButtonTapAnimation(bg *canvas.Rectangle, w gui.Widget) *gui.Animation {
	return gui.NewAnimation(canvas.DurationStandard, func(done float32) {
		mid := (w.Size().Width - theme.Padding()) / 2
		size := mid * done
		bg.Resize(gui.NewSize(size*2, w.Size().Height-theme.Padding()))
		bg.Move(gui.NewPos(mid-size, theme.Padding()/2))

		r, g, bb, a := col.ToNRGBA(theme.PressedColor())
		aa := uint8(a)
		fade := aa - uint8(float32(aa)*done)
		bg.FillColor = &color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(bb), A: fade}
		canvas.Refresh(bg)
	})
}
