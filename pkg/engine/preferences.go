package engine

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

// Preferences describes the ways that an app can save and load user preferences
type Preferences interface {
	// Bool looks up a boolean value for the key
	Bool(key string) bool
	// BoolWithFallback looks up a boolean value and returns the given fallback if not found
	BoolWithFallback(key string, fallback bool) bool
	// SetBool saves a boolean value for the given key
	SetBool(key string, value bool)

	// Float looks up a float64 value for the key
	Float(key string) float64
	// FloatWithFallback looks up a float64 value and returns the given fallback if not found
	FloatWithFallback(key string, fallback float64) float64
	// SetFloat saves a float64 value for the given key
	SetFloat(key string, value float64)

	// Int looks up an integer value for the key
	Int(key string) int
	// IntWithFallback looks up an integer value and returns the given fallback if not found
	IntWithFallback(key string, fallback int) int
	// SetInt saves an integer value for the given key
	SetInt(key string, value int)

	// String looks up a string value for the key
	String(key string) string
	// StringWithFallback looks up a string value and returns the given fallback if not found
	StringWithFallback(key, fallback string) string
	// SetString saves a string value for the given key
	SetString(key string, value string)

	// RemoveValue removes a value for the given key (not currently supported on iOS)
	RemoveValue(key string)

	// AddChangeListener allows code to be notified when some preferences change. This will fire on any update.
	AddChangeListener(func())
}
