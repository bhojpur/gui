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
	"sync"
	"testing"
	"time"

	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/stretchr/testify/assert"
)

func TestBindPreference_DataRace(t *testing.T) {
	a := test.NewApp()
	p := a.Preferences()
	key := "test-key"
	const n = 100

	var wg sync.WaitGroup
	binds := make([]Int, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(index int) {
			bind := BindPreferenceInt(key, p)
			binds[index] = bind
			wg.Done()
		}(i)
	}

	wg.Wait()
	expectedBind := binds[0]
	for i := 0; i < n; i++ {
		assert.Equal(t, expectedBind, binds[i])
	}
}

func TestBindPreferenceBool(t *testing.T) {
	a := test.NewApp()
	p := a.Preferences()
	key := "test-bool"

	p.SetBool(key, true)
	bind := BindPreferenceBool(key, p)
	v, err := bind.Get()
	assert.Nil(t, err)
	assert.True(t, v)

	err = bind.Set(false)
	assert.Nil(t, err)
	v, err = bind.Get()
	assert.Nil(t, err)
	assert.False(t, v)
	assert.False(t, p.Bool(key))
}

func TestBindPreferenceFloat(t *testing.T) {
	a := test.NewApp()
	p := a.Preferences()
	key := "test-float"

	p.SetFloat(key, 1.3)
	bind := BindPreferenceFloat(key, p)
	v, err := bind.Get()
	assert.Nil(t, err)
	assert.Equal(t, 1.3, v)

	err = bind.Set(2.5)
	assert.Nil(t, err)
	v, err = bind.Get()
	assert.Nil(t, err)
	assert.Equal(t, 2.5, v)
	assert.Equal(t, 2.5, p.Float(key))
}

func TestBindPreferenceInt(t *testing.T) {
	a := test.NewApp()
	p := a.Preferences()
	key := "test-int"

	p.SetInt(key, 4)
	bind := BindPreferenceInt(key, p)
	v, err := bind.Get()
	assert.Nil(t, err)
	assert.Equal(t, 4, v)

	err = bind.Set(7)
	assert.Nil(t, err)
	v, err = bind.Get()
	assert.Nil(t, err)
	assert.Equal(t, 7, v)
	assert.Equal(t, 7, p.Int(key))
}

func TestBindPreferenceString(t *testing.T) {
	a := test.NewApp()
	p := a.Preferences()
	key := "test-string"

	p.SetString(key, "aString")
	bind := BindPreferenceString(key, p)
	v, err := bind.Get()
	assert.Nil(t, err)
	assert.Equal(t, "aString", v)

	err = bind.Set("overwritten")
	assert.Nil(t, err)
	v, err = bind.Get()
	assert.Nil(t, err)
	assert.Equal(t, "overwritten", v)
	assert.Equal(t, "overwritten", p.String(key))
}

func TestPreferenceBindingCopies(t *testing.T) {
	a := test.NewApp()
	p := a.Preferences()
	key := "test-string"

	p.SetString(key, "aString")
	bind1 := BindPreferenceString(key, p)
	bind2 := BindPreferenceString(key, p)
	v1, err := bind1.Get()
	assert.Nil(t, err)
	v2, err := bind2.Get()
	assert.Nil(t, err)
	assert.Equal(t, v2, v1)

	err = bind1.Set("overwritten")
	assert.Nil(t, err)
	v1, err = bind1.Get()
	assert.Nil(t, err)
	v2, err = bind2.Get()
	assert.Nil(t, err)
	assert.Equal(t, v2, v1)
}

func TestPreferenceBindingTriggers(t *testing.T) {
	a := test.NewApp()
	p := a.Preferences()
	key := "test-string"

	p.SetString(key, "aString")
	bind1 := BindPreferenceString(key, p)
	bind2 := BindPreferenceString(key, p)

	ch := make(chan interface{}, 2)
	bind1.AddListener(NewDataListener(func() {
		ch <- struct{}{}
	}))

	select {
	case <-ch: // bind1 gets initial value
	case <-time.After(time.Millisecond * 100):
		t.Errorf("Timed out waiting for data binding to send initial value")
	}

	err := bind2.Set("overwritten") // write on a different listener, preferences should trigger all
	assert.Nil(t, err)
	select {
	case <-ch: // bind1 triggered by bind2 changing the same key
	case <-time.After(time.Millisecond * 100):
		t.Errorf("Timed out waiting for data binding change to trigger")
	}

	p.SetString(key, "overwritten2") // changing preference should trigger as well
	select {
	case <-ch: // bind1 triggered by preferences changing the same key directly
	case <-time.After(time.Millisecond * 300):
		t.Errorf("Timed out waiting for data binding change to trigger")
	}
}
