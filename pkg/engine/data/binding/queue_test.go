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
	"os"
	"runtime"
	"sync"
	"testing"

	"github.com/bhojpur/gui/pkg/engine/internal/async"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// TestQueryLazyInit resets the current unbounded func queue, and tests
// if the queue is lazy initialized.
//
// Note that this test may fail, if any of other tests in this package
// calls t.Parallel().
func TestQueueLazyInit(t *testing.T) {
	if queue != nil { // Reset queues
		queue.Close()
		queue = nil
		once = sync.Once{}
	}

	initialGoRoutines := runtime.NumGoroutine()

	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		queueItem(func() { wg.Done() })
	}
	wg.Wait()

	n := runtime.NumGoroutine()
	if n > initialGoRoutines+2 {
		t.Fatalf("unexpected number of goroutines after initialization, probably leaking: got %v want %v", n, initialGoRoutines+2)
	}
}

func TestQueueItem(t *testing.T) {
	called := 0
	queueItem(func() { called++ })
	queueItem(func() { called++ })
	waitForItems()
	assert.Equal(t, 2, called)
}

func TestMakeInfiniteQueue(t *testing.T) {
	var wg sync.WaitGroup
	queue := async.NewUnboundedFuncChan()

	wg.Add(1)
	c := 0
	go func() {
		for range queue.Out() {
			c++
		}
		wg.Done()
	}()

	for i := 0; i < 2048; i++ {
		queue.In() <- func() {}
	}
	queue.Close()

	wg.Wait()
	assert.Equal(t, 2048, c)
}
