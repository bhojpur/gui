package app

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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func loadPreferences(id string) *preferences {
	p := newPreferences(&bhojpurApp{uniqueID: id})
	p.load()

	return p
}

func TestPreferences_Save(t *testing.T) {
	p := loadPreferences("dummy")
	p.WriteValues(func(val map[string]interface{}) {
		val["keyString"] = "value"
		val["keyInt"] = 4
		val["keyFloat"] = 3.5
		val["keyBool"] = true
	})

	path := filepath.Join(os.TempDir(), "bhojpurPrefs.json")
	defer os.Remove(path)
	p.saveToFile(path)

	expected, err := ioutil.ReadFile(filepath.Join("testdata", "preferences.json"))
	if err != nil {
		assert.Fail(t, "Failed to load, %v", err)
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		assert.Fail(t, "Failed to load, %v", err)
	}
	assert.JSONEq(t, string(expected), string(content))
}

func TestPreferences_Save_OverwriteFast(t *testing.T) {
	p := loadPreferences("dummy2")
	p.WriteValues(func(val map[string]interface{}) {
		val["key"] = "value"
	})

	path := filepath.Join(os.TempDir(), "bhojpurPrefs2.json")
	defer os.Remove(path)
	p.saveToFile(path)

	p.WriteValues(func(val map[string]interface{}) {
		val["key2"] = "value2"
	})
	p.saveToFile(path)

	p2 := loadPreferences("dummy")
	p2.loadFromFile(path)
	assert.Equal(t, "value2", p2.String("key2"))
}

func TestPreferences_Load(t *testing.T) {
	p := loadPreferences("dummy")
	p.loadFromFile(filepath.Join("testdata", "preferences.json"))

	assert.Equal(t, "value", p.String("keyString"))
	assert.Equal(t, 4, p.Int("keyInt"))
	assert.Equal(t, 3.5, p.Float("keyFloat"))
	assert.Equal(t, true, p.Bool("keyBool"))
}
