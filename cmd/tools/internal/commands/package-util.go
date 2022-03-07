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
	"io/ioutil"

	realUtil "github.com/bhojpur/gui/cmd/tools/internal/util"
)

type packagerUtil interface {
	Exists(path string) bool
	CopyFile(source string, target string) error
	CopyExeFile(src, tgt string) error
	WriteFile(target string, data []byte) error
	EnsureSubDir(parent, name string) string

	RequireAndroidSDK() error
	AndroidBuildToolsPath() string

	IsAndroid(os string) bool
	IsIOS(os string) bool
	IsMobile(os string) bool
}

type defaultUtil struct{}

func (d defaultUtil) Exists(path string) bool {
	return realUtil.Exists(path)
}

func (d defaultUtil) CopyFile(source string, target string) error {
	return realUtil.CopyFile(source, target)
}

func (d defaultUtil) CopyExeFile(src, tgt string) error {
	return realUtil.CopyExeFile(src, tgt)
}

func (d defaultUtil) WriteFile(target string, data []byte) error {
	return ioutil.WriteFile(target, data, 0644)
}

func (d defaultUtil) EnsureSubDir(parent, name string) string {
	return realUtil.EnsureSubDir(parent, name)
}

func (d defaultUtil) RequireAndroidSDK() error {
	return realUtil.RequireAndroidSDK()
}

func (d defaultUtil) AndroidBuildToolsPath() string {
	return realUtil.AndroidBuildToolsPath()
}

func (d defaultUtil) IsAndroid(os string) bool {
	return realUtil.IsAndroid(os)
}

func (d defaultUtil) IsIOS(os string) bool {
	return realUtil.IsIOS(os)
}

func (d defaultUtil) IsMobile(os string) bool {
	return realUtil.IsMobile(os)
}

var util packagerUtil

func init() {
	util = defaultUtil{}
}
