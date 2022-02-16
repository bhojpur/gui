package mobile

//go:generate gomobile help documentation doc.go

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
	"flag"
	"fmt"
	"os"
)

var (
	gomobileName = "gomobile"
)

const documentationHeader = `// Copyright 2018 Bhojpur Consulting Private Limited, India. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by Bhojpur GUI - Foundation Framework 'gomobile help documentation doc.go'. DO NOT EDIT.
`

type command struct {
	run   func(*command) error
	Flag  flag.FlagSet
	Name  string
	Usage string
	Short string
	Long  string

	IconPath, AppName      string
	Version, Cert, Profile string
	Build                  int
}

func (cmd *command) usage() {
	fmt.Fprintf(os.Stdout, "usage: %s %s %s\n%s", gomobileName, cmd.Name, cmd.Usage, cmd.Long)
}
