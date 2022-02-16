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
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Load attempts to read a BhojpurApp metadata from the provided reader.
// If this cannot be done an error will be returned.
func Load(r io.Reader) (*BhojpurApp, error) {
	str, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var data BhojpurApp
	if _, err := toml.Decode(string(str), &data); err != nil {
		return nil, err
	}

	return &data, nil
}

// LoadStandard attempts to read a BhojpurApp metadata from the `BhojpurApp.toml` file in the specified dir.
// If the file cannot be found or parsed an error will be returned.
func LoadStandard(dir string) (*BhojpurApp, error) {
	path := filepath.Join(dir, "BhojpurApp.toml")
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer r.Close()
	return Load(r)
}
