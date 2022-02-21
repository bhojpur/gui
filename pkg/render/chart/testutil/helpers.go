package testutil

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
	"math"
	"reflect"
	"strings"
	"testing"
)

// AssertNil asserts an argument is nil.
func AssertNil(t *testing.T, actual interface{}) {
	if !isNil(actual) {
		t.Errorf("assertion failed; expected actual to be nil")
		t.FailNow()
	}
}

// AssertNotNil asserts an argument is not nil.
func AssertNotNil(t *testing.T, actual interface{}, message ...interface{}) {
	if isNil(actual) {
		t.Error("assertion failed; expected actual to not be nil")
		if len(message) > 0 {
			t.Error(message...)
		}
		t.FailNow()
	}
}

// AssertEqual asserts two arguments are equal.
func AssertEqual(t *testing.T, expected, actual interface{}, message ...interface{}) {
	if !equal(expected, actual) {
		t.Errorf("assertion failed; expected %v to equal %v", actual, expected)
		if len(message) > 0 {
			t.Error(message...)
		}
		t.FailNow()
	}
}

// AssertNotEqual asserts two arguments are not equal.
func AssertNotEqual(t *testing.T, expected, actual interface{}, message ...interface{}) {
	if equal(expected, actual) {
		t.Errorf("assertion failed; expected %v to not equal %v", actual, expected)
		if len(message) > 0 {
			t.Error(message...)
		}
		t.FailNow()
	}
}

// AssertZero asserts an argument is zero.
func AssertZero(t *testing.T, actual interface{}, message ...interface{}) {
	AssertEqual(t, 0, actual)
}

// AssertNotZero asserts an argument is not zero.
func AssertNotZero(t *testing.T, actual interface{}, message ...interface{}) {
	AssertNotEqual(t, 0, actual)
}

// AssertTrue asserts an argument is true.
func AssertTrue(t *testing.T, arg bool, message ...interface{}) {
	if !arg {
		t.Errorf("assertion failed; expected actual to be true")
		if len(message) > 0 {
			t.Error(message...)
		}
		t.FailNow()
	}
}

// AssertFalse asserts an argument is false.
func AssertFalse(t *testing.T, arg bool, message ...interface{}) {
	if arg {
		t.Errorf("assertion failed; expected actual to be false")
		if len(message) > 0 {
			t.Error(message...)
		}
		t.FailNow()
	}
}

// AssertInDelta asserts a two arguments are within a delta of eachother.
//
// This delta will be determined absolute, and the delta should always be positive.
func AssertInDelta(t *testing.T, from, to, delta float64, message ...interface{}) {
	if diff := math.Abs(from - to); diff > delta {
		t.Errorf("assertion failed; expected absolute difference of %f and %f to be %f", from, to, delta)
		if len(message) > 0 {
			t.Error(message...)
		}
		t.FailNow()
	}
}

// AssertEmpty asserts an argument is empty.
func AssertEmpty(t *testing.T, arg interface{}, message ...interface{}) {
	if getLength(arg) != 0 {
		t.Errorf("assertion failed; expected actual to be empty")
		if len(message) > 0 {
			t.Error(message...)
		}
		t.FailNow()
	}
}

// AssertNotEmpty asserts an argument is not empty.
func AssertNotEmpty(t *testing.T, arg interface{}, message ...interface{}) {
	if getLength(arg) == 0 {
		t.Errorf("assertion failed; expected actual to not be empty")
		if len(message) > 0 {
			t.Error(message...)
		}
		t.FailNow()
	}
}

// AssertLen asserts an argument has a given length.
func AssertLen(t *testing.T, arg interface{}, length int, message ...interface{}) {
	if getLength(arg) != length {
		t.Errorf("assertion failed; expected actual to have length %d", length)
		if len(message) > 0 {
			t.Error(message...)
		}
		t.FailNow()
	}
}

// AssertContains asserts an argument contains a given substring.
func AssertContains(t *testing.T, s, substr string, message ...interface{}) {
	if !strings.Contains(s, substr) {
		t.Errorf("assertion failed; expected actual to contain %q", substr)
		if len(message) > 0 {
			t.Error(message...)
		}
		t.FailNow()
	}
}

// AssertNotContains asserts an argument does not contain a given substring.
func AssertNotContains(t *testing.T, s, substr string, message ...interface{}) {
	if strings.Contains(s, substr) {
		t.Errorf("assertion failed; expected actual to not contain %q", substr)
		if len(message) > 0 {
			t.Error(message...)
		}
		t.FailNow()
	}
}

func equal(expected, actual interface{}) bool {
	if expected == nil && actual == nil {
		return true
	}
	if (expected == nil && actual != nil) || (expected != nil && actual == nil) {
		return false
	}

	actualType := reflect.TypeOf(actual)
	if actualType == nil {
		return false
	}
	expectedValue := reflect.ValueOf(expected)
	if expectedValue.IsValid() && expectedValue.Type().ConvertibleTo(actualType) {
		return reflect.DeepEqual(expectedValue.Convert(actualType).Interface(), actual)
	}

	return reflect.DeepEqual(expected, actual)
}

func isNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}
	return false
}

func getLength(object interface{}) int {
	if object == nil {
		return 0
	} else if object == "" {
		return 0
	}

	objValue := reflect.ValueOf(object)

	switch objValue.Kind() {
	case reflect.Map:
		fallthrough
	case reflect.Slice, reflect.Chan, reflect.String:
		{
			return objValue.Len()
		}
	}
	return 0
}
