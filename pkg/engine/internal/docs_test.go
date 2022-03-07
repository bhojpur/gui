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

	intRepo "github.com/bhojpur/gui/pkg/engine/internal/repository"
	"github.com/bhojpur/gui/pkg/engine/storage"
	"github.com/bhojpur/gui/pkg/engine/storage/repository"

	"github.com/stretchr/testify/assert"
)

func TestDocs_Create(t *testing.T) {
	repository.Register("file", intRepo.NewInMemoryRepository("file"))
	docs := &Docs{storage.NewFileURI("/tmp/docs/create")}
	exist, err := storage.Exists(docs.RootDocURI)
	assert.Nil(t, err)
	assert.False(t, exist)

	w, err := docs.Create("test")
	assert.Nil(t, err)
	_ = w.Close()

	exist, err = storage.Exists(docs.RootDocURI)
	assert.Nil(t, err)
	assert.True(t, exist)
}

func TestDocs_List(t *testing.T) {
	repository.Register("file", intRepo.NewInMemoryRepository("file"))
	docs := &Docs{storage.NewFileURI("/tmp/docs/list")}
	list := docs.List()
	assert.Zero(t, len(list))

	w, _ := docs.Create("1")
	_, _ = w.Write([]byte{})
	_ = w.Close()
	w, _ = docs.Create("2")
	_, _ = w.Write([]byte{})
	_ = w.Close()

	list = docs.List()
	assert.Equal(t, 2, len(list))
}

func TestDocs_Save(t *testing.T) {
	r := intRepo.NewInMemoryRepository("file")
	repository.Register("file", r)
	docs := &Docs{storage.NewFileURI("/tmp/docs/save")}
	w, err := docs.Create("save.txt")
	assert.Nil(t, err)
	_, _ = w.Write([]byte{})
	_ = w.Close()
	u := w.URI()

	exist, err := r.Exists(u)
	assert.Nil(t, err)
	assert.True(t, exist)

	w, err = docs.Save("save.txt")
	assert.Nil(t, err)
	n, err := w.Write([]byte("save"))
	assert.Nil(t, err)
	assert.Equal(t, 4, n)
	err = w.Close()
	assert.Nil(t, err)
}
