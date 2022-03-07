package internal

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

	"github.com/stretchr/testify/assert"

	gui "github.com/bhojpur/gui/pkg/engine"
)

func TestClipStack_Intersect(t *testing.T) {
	p1 := gui.NewPos(5, 25)
	s1 := gui.NewSize(100, 100)
	c := &ClipStack{
		clips: []*ClipItem{
			{p1, s1},
		},
	}

	p2 := gui.NewPos(25, 0)
	s2 := gui.NewSize(50, 50)
	i := c.Push(p2, s2)

	assert.Equal(t, gui.NewPos(25, 25), i.pos)
	assert.Equal(t, gui.NewSize(50, 25), i.size)
	assert.Equal(t, 2, len(c.clips))

	_ = c.Pop()
	p2 = gui.NewPos(50, 50)
	s2 = gui.NewSize(150, 50)
	i = c.Push(p2, s2)

	assert.Equal(t, gui.NewPos(50, 50), i.pos)
	assert.Equal(t, gui.NewSize(55, 50), i.size)
	assert.Equal(t, 2, len(c.clips))
}

func TestClipStack_Pop(t *testing.T) {
	p := gui.NewPos(5, 5)
	s := gui.NewSize(100, 100)
	c := &ClipStack{
		clips: []*ClipItem{
			{p, s},
		},
	}

	i := c.Pop()
	assert.Equal(t, p, i.pos)
	assert.Equal(t, s, i.size)
	assert.Equal(t, 0, len(c.clips))
}

func TestClipStack_Push(t *testing.T) {
	c := &ClipStack{}
	p := gui.NewPos(5, 5)
	s := gui.NewSize(100, 100)

	i := c.Push(p, s)
	assert.Equal(t, p, i.pos)
	assert.Equal(t, s, i.size)
	assert.Equal(t, 1, len(c.clips))
}
