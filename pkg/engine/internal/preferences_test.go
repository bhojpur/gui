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
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPrefs_SetBool(t *testing.T) {
	p := NewInMemoryPreferences()
	p.SetBool("testBool", true)

	assert.Equal(t, true, p.Bool("testBool"))
}

func TestPrefs_Bool(t *testing.T) {
	p := NewInMemoryPreferences()
	p.WriteValues(func(val map[string]interface{}) {
		val["testBool"] = true
	})

	assert.Equal(t, true, p.Bool("testBool"))
}

func TestPrefs_BoolWithFallback(t *testing.T) {
	p := NewInMemoryPreferences()

	assert.Equal(t, true, p.BoolWithFallback("testBool", true))
	p.SetBool("testBool", false)
	assert.Equal(t, false, p.BoolWithFallback("testBool", true))
}

func TestPrefs_Bool_Zero(t *testing.T) {
	p := NewInMemoryPreferences()

	assert.Equal(t, false, p.Bool("testBool"))
}

func TestPrefs_SetFloat(t *testing.T) {
	p := NewInMemoryPreferences()
	p.SetFloat("testFloat", 1.7)

	assert.Equal(t, 1.7, p.Float("testFloat"))
}

func TestPrefs_Float(t *testing.T) {
	p := NewInMemoryPreferences()
	p.WriteValues(func(val map[string]interface{}) {
		val["testFloat"] = 1.2
	})

	assert.Equal(t, 1.2, p.Float("testFloat"))
}

func TestPrefs_FloatWithFallback(t *testing.T) {
	p := NewInMemoryPreferences()

	assert.Equal(t, 1.0, p.FloatWithFallback("testFloat", 1.0))
	p.WriteValues(func(val map[string]interface{}) {
		val["testFloat"] = 1.2
	})
	assert.Equal(t, 1.2, p.FloatWithFallback("testFloat", 1.0))
}

func TestPrefs_Float_Zero(t *testing.T) {
	p := NewInMemoryPreferences()

	assert.Equal(t, 0.0, p.Float("testFloat"))
}

func TestPrefs_SetInt(t *testing.T) {
	p := NewInMemoryPreferences()
	p.SetInt("testInt", 5)

	assert.Equal(t, 5, p.Int("testInt"))
}

func TestPrefs_Int(t *testing.T) {
	p := NewInMemoryPreferences()
	p.WriteValues(func(val map[string]interface{}) {
		val["testInt"] = 5
	})
	assert.Equal(t, 5, p.Int("testInt"))
}

func TestPrefs_IntWithFallback(t *testing.T) {
	p := NewInMemoryPreferences()

	assert.Equal(t, 2, p.IntWithFallback("testInt", 2))
	p.WriteValues(func(val map[string]interface{}) {
		val["testInt"] = 5
	})
	assert.Equal(t, 5, p.IntWithFallback("testInt", 2))
}

func TestPrefs_Int_Zero(t *testing.T) {
	p := NewInMemoryPreferences()

	assert.Equal(t, 0, p.Int("testInt"))
}

func TestPrefs_SetString(t *testing.T) {
	p := NewInMemoryPreferences()
	p.SetString("test", "value")

	assert.Equal(t, "value", p.String("test"))
}

func TestPrefs_String(t *testing.T) {
	p := NewInMemoryPreferences()
	p.WriteValues(func(val map[string]interface{}) {
		val["test"] = "value"
	})

	assert.Equal(t, "value", p.String("test"))
}

func TestPrefs_StringWithFallback(t *testing.T) {
	p := NewInMemoryPreferences()

	assert.Equal(t, "default", p.StringWithFallback("test", "default"))
	p.WriteValues(func(val map[string]interface{}) {
		val["test"] = "value"
	})
	assert.Equal(t, "value", p.StringWithFallback("test", "default"))
}

func TestPrefs_String_Zero(t *testing.T) {
	p := NewInMemoryPreferences()

	assert.Equal(t, "", p.String("test"))
}

func TestInMemoryPreferences_OnChange(t *testing.T) {
	p := NewInMemoryPreferences()
	called := false
	p.AddChangeListener(func() {
		called = true
	})

	p.SetString("dummy", "another")
	time.Sleep(time.Millisecond * 100)

	assert.True(t, called)
}

func TestRemoveValue(t *testing.T) {
	p := NewInMemoryPreferences()

	p.SetBool("dummy", true)
	p.SetFloat("pi", 3.14)
	p.SetInt("number", 2)
	p.SetString("month", "January")

	p.RemoveValue("dummy")
	p.RemoveValue("pi")
	p.RemoveValue("number")
	p.RemoveValue("month")

	assert.Equal(t, false, p.Bool("dummy"))
	assert.Equal(t, float64(0), p.Float("pi"))
	assert.Equal(t, 0, p.Int("number"))
	assert.Equal(t, "", p.String("month"))
}
