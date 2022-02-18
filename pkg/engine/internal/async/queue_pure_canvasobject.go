//go:build js
// +build js

// Code generated by Bhojpur GUI Foundation Framework (go run gen.go); DO NOT EDIT.
package async

import (
	gui "github.com/bhojpur/gui/pkg/engine"
)

// CanvasObjectQueue implements lock-free FIFO freelist based queue.
//
// Reference: https://dl.acm.org/citation.cfm?doid=248052.248106
type CanvasObjectQueue struct {
	head *itemCanvasObject
	tail *itemCanvasObject
	len  uint64
}

// NewCanvasObjectQueue returns a queue for caching values.
func NewCanvasObjectQueue() *CanvasObjectQueue {
	head := &itemCanvasObject{next: nil, v: nil}
	return &CanvasObjectQueue{
		tail: head,
		head: head,
	}
}

type itemCanvasObject struct {
	next *itemCanvasObject
	v    gui.CanvasObject
}

func loadCanvasObjectItem(p **itemCanvasObject) *itemCanvasObject {
	return *p
}
func casCanvasObjectItem(p **itemCanvasObject, _, new *itemCanvasObject) bool {
	*p = new
	return true
}
