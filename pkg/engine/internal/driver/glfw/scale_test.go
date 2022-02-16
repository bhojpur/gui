//go:build !mobile
// +build !mobile

package glfw

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
	"os"
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"
	_ "github.com/bhojpur/gui/pkg/engine/test"

	"github.com/stretchr/testify/assert"
)

func TestCalculateDetectedScale(t *testing.T) {
	lowDPI := calculateDetectedScale(300, 1600)
	assert.Equal(t, float32(1.1288888), lowDPI)

	hiDPI := calculateDetectedScale(420, 3800)
	assert.Equal(t, float32(1.9150794), hiDPI)
}

func TestCalculateDetectedScale_Min(t *testing.T) {
	lowDPI := calculateDetectedScale(300, 1280)
	assert.Equal(t, float32(1), lowDPI)
}

func TestCalculateScale(t *testing.T) {
	one := calculateScale(1.0, 1.0, 1.0)
	assert.Equal(t, float32(1.0), one)

	larger := calculateScale(1.5, 1.0, 1.0)
	assert.Equal(t, float32(1.5), larger)

	smaller := calculateScale(0.8, 1.0, 1.0)
	assert.Equal(t, float32(0.8), smaller)

	hiDPI := calculateScale(0.8, 2.0, 1.0)
	assert.Equal(t, float32(1.6), hiDPI)

	hiDPIAuto := calculateScale(0.8, scaleAuto, 2.0)
	assert.Equal(t, float32(1.6), hiDPIAuto)

	large := calculateScale(1.5, 2.0, 2.0)
	assert.Equal(t, float32(3.0), large)
}

func TestCalculateScale_Round(t *testing.T) {
	trunc := calculateScale(1.04321, 1.0, 1.0)
	assert.Equal(t, float32(1.0), trunc)

	round := calculateScale(1.1, 1.1, 1.0)
	assert.Equal(t, float32(1.2), round)
}

func TestUserScale(t *testing.T) {
	envVal := os.Getenv(scaleEnvKey)
	defer os.Setenv(scaleEnvKey, envVal)

	_ = os.Setenv(scaleEnvKey, "auto")
	set := gui.CurrentApp().Settings().Scale()
	if set == float32(0.0) { // no config set
		assert.Equal(t, float32(1.0), userScale())
	} else {
		assert.Equal(t, set, userScale())
	}

	_ = os.Setenv(scaleEnvKey, "1.2")
	assert.Equal(t, float32(1.2), userScale())
}
