package desktop

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
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"
)

func TestCustomShortcut_Shortcut(t *testing.T) {
	type fields struct {
		KeyName  gui.KeyName
		Modifier Modifier
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Ctrl+C",
			fields: fields{
				KeyName:  gui.KeyC,
				Modifier: ControlModifier,
			},
			want: "CustomDesktop:Control+C",
		},
		{
			name: "Ctrl+Alt+Esc",
			fields: fields{
				KeyName:  gui.KeyEscape,
				Modifier: ControlModifier + AltModifier,
			},
			want: "CustomDesktop:Control+Alt+Escape",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &CustomShortcut{
				KeyName:  tt.fields.KeyName,
				Modifier: tt.fields.Modifier,
			}
			if got := cs.ShortcutName(); got != tt.want {
				t.Errorf("CustomShortcut.ShortcutName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_modifierToString(t *testing.T) {
	tests := []struct {
		name string
		mods Modifier
		want string
	}{
		{
			name: "None",
			mods: 0,
			want: "",
		},
		{
			name: "Ctrl",
			mods: ControlModifier,
			want: "Control",
		},
		{
			name: "Shift+Ctrl",
			mods: ShiftModifier + ControlModifier,
			want: "Shift+Control",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := modifierToString(tt.mods); got != tt.want {
				t.Errorf("modifierToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
