package metadata

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
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Save attempts to write a BhojpurApp metadata to the provided writer.
// If the encoding fails an error will be returned.
func Save(f *BhojpurApp, w io.Writer) error {
	var buf bytes.Buffer
	e := toml.NewEncoder(&buf)
	err := e.Encode(f)
	if err != nil {
		return err
	}

	_, err = w.Write(buf.Bytes())
	return err
}

// SaveStandard attempts to save a BhojpurApp metadata to the `BhojpurApp.toml` file in the specified dir.
// If the file cannot be written or encoding fails an error will be returned.
func SaveStandard(f *BhojpurApp, dir string) error {
	path := filepath.Join(dir, "BhojpurApp.toml")
	w, err := os.Create(path)
	if err != nil {
		return err
	}

	defer w.Close()
	return Save(f, w)
}
