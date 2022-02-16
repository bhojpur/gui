package dialog

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
	"errors"
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/widget"
	"github.com/stretchr/testify/assert"
)

// formDialogResult is the result of the test form dialog callback.
type formDialogResult int

const (
	formDialogNoAction formDialogResult = iota
	formDialogConfirm
	formDialogCancel
)

func TestFormDialog_Control(t *testing.T) {
	var result formDialogResult
	fd := controlFormDialog(&result, test.NewWindow(nil))
	fd.Show()
	test.Tap(fd.confirm)

	assert.Equal(t, formDialogConfirm, result, "Control form should be confirmed with no validation")
}

func TestFormDialog_InvalidCannotSubmit(t *testing.T) {
	var result formDialogResult
	fd := validatingFormDialog(&result, test.NewWindow(nil))
	fd.Show()

	assert.False(t, fd.win.Hidden)
	assert.True(t, fd.confirm.Disabled(), "Confirm button should be disabled due to validation state")
	test.Tap(fd.confirm)

	assert.Equal(t, formDialogNoAction, result, "Callback should not have ran with invalid form")
}

func TestFormDialog_ValidCanSubmit(t *testing.T) {
	var result formDialogResult
	fd := validatingFormDialog(&result, test.NewWindow(nil))
	fd.Show()

	assert.False(t, fd.win.Hidden)
	assert.True(t, fd.confirm.Disabled(), "Confirm button should be disabled due to validation state")

	if validatingEntry, ok := fd.items[0].Widget.(*widget.Entry); ok {
		validatingEntry.SetText("abc")
		assert.False(t, fd.confirm.Disabled())
		test.Tap(fd.confirm)

		assert.Equal(t, formDialogConfirm, result, "Valid form should be able to be confirmed")
	} else {
		assert.Fail(t, "First item's widget should be an Entry (check validatingFormDialog)")
	}
}

func TestFormDialog_CanCancelInvalid(t *testing.T) {
	var result formDialogResult
	fd := validatingFormDialog(&result, test.NewWindow(nil))
	fd.Show()
	assert.False(t, fd.win.Hidden)

	test.Tap(fd.dismiss)

	assert.Equal(t, formDialogCancel, result, "Expected cancel result")
}

func TestFormDialog_CanCancelNoValidation(t *testing.T) {
	var result formDialogResult
	fd := controlFormDialog(&result, test.NewWindow(nil))
	fd.Show()
	assert.False(t, fd.win.Hidden)

	test.Tap(fd.dismiss)

	assert.Equal(t, formDialogCancel, result, "Expected cancel result")
}

func validatingFormDialog(result *formDialogResult, parent gui.Window) *formDialog {
	validatingEntry := widget.NewEntry()
	validatingEntry.Validator = func(input string) error {
		if input != "abc" {
			return errors.New("only accepts 'abc'")
		}
		return nil
	}
	validatingItem := &widget.FormItem{
		Text:   "Only accepts 'abc'",
		Widget: validatingEntry,
	}
	controlEntry := widget.NewPasswordEntry()
	controlItem := &widget.FormItem{
		Text:   "I accept anything",
		Widget: controlEntry,
	}

	items := []*widget.FormItem{validatingItem, controlItem}
	return NewForm("Validating Form Dialog", "Submit", "Cancel", items, func(confirm bool) {
		if confirm {
			*result = formDialogConfirm
		} else {
			*result = formDialogCancel
		}
	}, parent).(*formDialog)
}

func controlFormDialog(result *formDialogResult, parent gui.Window) *formDialog {
	controlEntry := widget.NewEntry()
	controlItem := &widget.FormItem{
		Text:   "I accept anything",
		Widget: controlEntry,
	}
	controlEntry2 := widget.NewPasswordEntry()
	controlItem2 := &widget.FormItem{
		Text:   "I accept anything",
		Widget: controlEntry2,
	}
	items := []*widget.FormItem{controlItem, controlItem2}
	return NewForm("Validating Form Dialog", "Submit", "Cancel", items, func(confirm bool) {
		if confirm {
			*result = formDialogConfirm
		} else {
			*result = formDialogCancel
		}
	}, parent).(*formDialog)
}
