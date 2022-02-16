package binding

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

	"github.com/bhojpur/gui/pkg/engine/storage"
	"github.com/stretchr/testify/assert"
)

func TestBindFloat(t *testing.T) {
	val := 0.5
	f := BindFloat(&val)
	v, err := f.Get()
	assert.Nil(t, err)
	assert.Equal(t, 0.5, v)

	called := false
	fn := NewDataListener(func() {
		called = true
	})
	f.AddListener(fn)
	waitForItems()
	assert.True(t, called)

	called = false
	err = f.Set(0.3)
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, 0.3, val)
	assert.True(t, called)

	called = false
	val = 1.2
	f.Reload()
	waitForItems()
	assert.True(t, called)
	v, err = f.Get()
	assert.Nil(t, err)
	assert.Equal(t, 1.2, v)
}

func TestNewFloat(t *testing.T) {
	f := NewFloat()
	v, err := f.Get()
	assert.Nil(t, err)
	assert.Equal(t, 0.0, v)

	err = f.Set(0.3)
	assert.Nil(t, err)
	v, err = f.Get()
	assert.Nil(t, err)
	assert.Equal(t, 0.3, v)
}

func TestBindFloat_TriggerOnlyWhenChange(t *testing.T) {
	v := 9.8
	f := BindFloat(&v)
	triggered := 0
	f.AddListener(NewDataListener(func() {
		triggered++
	}))
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	for i := 0; i < 5; i++ {
		err := f.Set(9.0)
		assert.Nil(t, err)
		waitForItems()
		assert.Equal(t, 1, triggered)
	}

	triggered = 0
	err := f.Set(9.1)
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	err = f.Set(8.0)
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	v = 8.0
	err = f.Reload()
	assert.NoError(t, err)
	waitForItems()
	assert.Equal(t, 0, triggered)

	triggered = 0
	v = 9.2
	err = f.Reload()
	assert.NoError(t, err)
	waitForItems()
	assert.Equal(t, 1, triggered)
}

func TestNewFloat_TriggerOnlyWhenChange(t *testing.T) {
	f := NewFloat()
	triggered := 0
	f.AddListener(NewDataListener(func() {
		triggered++
	}))
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	for i := 0; i < 5; i++ {
		err := f.Set(9.0)
		assert.Nil(t, err)
		waitForItems()
		assert.Equal(t, 1, triggered)
	}

	triggered = 0
	err := f.Set(9.1)
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	err = f.Set(8.0)
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, 1, triggered)
}

func TestBindURI(t *testing.T) {
	val := storage.NewFileURI("/tmp")
	f := BindURI(&val)
	v, err := f.Get()
	assert.Nil(t, err)
	assert.Equal(t, "file:///tmp", v.String())

	called := false
	fn := NewDataListener(func() {
		called = true
	})
	f.AddListener(fn)
	waitForItems()
	assert.True(t, called)

	called = false
	err = f.Set(storage.NewFileURI("/tmp/test.txt"))
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, "file:///tmp/test.txt", val.String())
	assert.True(t, called)

	called = false
	val = storage.NewFileURI("/hello")
	_ = f.Reload()
	waitForItems()
	assert.True(t, called)
	v, err = f.Get()
	assert.Nil(t, err)
	assert.Equal(t, "file:///hello", v.String())
}

func TestBindURI_TriggerOnlyWhenChange(t *testing.T) {
	v := storage.NewFileURI("first")
	b := BindURI(&v)
	triggered := 0
	b.AddListener(NewDataListener(func() {
		triggered++
	}))
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	for i := 0; i < 5; i++ {
		err := b.Set(storage.NewFileURI("second"))
		assert.Nil(t, err)
		waitForItems()
		assert.Equal(t, 1, triggered)
	}

	triggered = 0
	err := b.Set(storage.NewFileURI("third"))
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	err = b.Set(storage.NewFileURI("fourth"))
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	v = storage.NewFileURI("fourth")
	err = b.Reload()
	assert.NoError(t, err)
	waitForItems()
	assert.Equal(t, 0, triggered)

	triggered = 0
	v = storage.NewFileURI("fifth")
	err = b.Reload()
	assert.NoError(t, err)
	waitForItems()
	assert.Equal(t, 1, triggered)
}

func TestNewURI(t *testing.T) {
	f := NewURI()
	v, err := f.Get()
	assert.Nil(t, err)
	assert.Equal(t, nil, v)

	err = f.Set(storage.NewFileURI("/var"))
	assert.Nil(t, err)
	v, err = f.Get()
	assert.Nil(t, err)
	assert.Equal(t, "file:///var", v.String())
}

func TestNewURI_TriggerOnlyWhenChange(t *testing.T) {
	b := NewURI()
	triggered := 0
	b.AddListener(NewDataListener(func() {
		triggered++
	}))
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	for i := 0; i < 5; i++ {
		err := b.Set(storage.NewFileURI("first"))
		assert.Nil(t, err)
		waitForItems()
		assert.Equal(t, 1, triggered)
	}

	triggered = 0
	err := b.Set(storage.NewFileURI("second"))
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	err = b.Set(storage.NewFileURI("third"))
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, 1, triggered)
}
