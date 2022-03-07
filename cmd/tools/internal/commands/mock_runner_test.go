package commands

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
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type expectedValue struct {
	dir   interface{}
	env   []string
	osEnv bool
	args  []string
}

type mockReturn struct {
	ret []byte
	err error
}

type expectedCall struct {
	dirSet bool
	envSet bool
}

type mockRunner struct {
	expectedValue
	expectedCall
	mockReturn
}

type testCommandRuns struct {
	runs       []mockRunner
	currentRun int
	t          *testing.T
}

func (t *testCommandRuns) runOutput(args ...string) ([]byte, error) {
	require.Less(t.t, t.currentRun, len(t.runs))
	require.Equal(t.t, len(t.runs[t.currentRun].args), len(args))

	expectedArgs := t.runs[t.currentRun].args
	require.Equal(t.t, expectedArgs, args)

	ret, err := t.runs[t.currentRun].ret, t.runs[t.currentRun].err
	t.currentRun++

	return ret, err
}

func (t *testCommandRuns) setDir(dir string) {
	require.Less(t.t, t.currentRun, len(t.runs))

	require.Equal(t.t, t.runs[t.currentRun].dir.(string), dir)
	t.runs[t.currentRun].dirSet = true
}

func (t *testCommandRuns) setEnv(env []string) {
	require.Less(t.t, t.currentRun, len(t.runs))

	// Prepare array for comparison
	expectedEnv := t.runs[t.currentRun].env
	if t.runs[t.currentRun].osEnv {
		expectedEnv = append(expectedEnv, os.Environ()...)
	}
	sort.Strings(expectedEnv)
	sort.Strings(env)

	require.Equal(t.t, len(expectedEnv), len(env))
	require.Equal(t.t, expectedEnv, env)

	t.runs[t.currentRun].envSet = true
}

func (t *testCommandRuns) verifyExpectation() {
	require.Equal(t.t, len(t.runs), t.currentRun)

	for _, value := range t.runs {
		if value.dir != nil {
			assert.Equal(t.t, true, value.dirSet)
		}
		if len(value.env) > 0 {
			assert.Equal(t.t, true, value.envSet)
		}
	}
}
