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
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBindTime(t *testing.T) {
	val := time.Now()
	f := bindTime(&val)
	v, err := f.Get()
	assert.Nil(t, err)
	assert.Equal(t, val.Unix(), v.Unix())

	called := false
	fn := NewDataListener(func() {
		called = true
	})
	f.AddListener(fn)
	waitForItems()
	assert.True(t, called)

	newTime := val.Add(time.Hour)
	called = false
	err = f.Set(newTime)
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, newTime.Unix(), val.Unix())
	assert.True(t, called)

	newTime = newTime.Add(time.Minute)
	called = false
	val = newTime
	_ = f.Reload()
	waitForItems()
	assert.True(t, called)
	v, err = f.Get()
	assert.Nil(t, err)
	assert.Equal(t, newTime.Unix(), v.Unix())
}

func TestNewTime(t *testing.T) {
	f := newTime()
	v, err := f.Get()
	assert.Nil(t, err)
	assert.Equal(t, time.Unix(0, 0), v)

	now := time.Now()
	err = f.Set(now)
	assert.Nil(t, err)
	v, err = f.Get()
	assert.Nil(t, err)
	assert.Equal(t, now.Unix(), v.Unix())
}

type timeBinding struct {
	Int
	src *time.Time
}

func bindTime(t *time.Time) *timeBinding {
	return &timeBinding{Int: NewInt(), src: t}
}

func newTime() *timeBinding {
	return &timeBinding{Int: NewInt()}
}

func (t *timeBinding) Get() (time.Time, error) {
	if t.src != nil {
		return *t.src, nil
	}

	i, err := t.Int.Get()
	return time.Unix(int64(i), 0), err
}

func (t *timeBinding) Reload() error {
	if t.src == nil {
		return nil
	}

	return t.Set(*t.src)
}

func (t *timeBinding) Set(time time.Time) error {
	if t.src != nil {
		*t.src = time
	}

	i := time.Unix()
	return t.Int.Set(int(i))
}
