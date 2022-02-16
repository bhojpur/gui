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
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/bhojpur/gui/cmd/tools/internal/templates"
	"github.com/bhojpur/gui/cmd/tools/internal/util"
	"github.com/jackmordaunt/icns"
	"github.com/pkg/errors"
)

type darwinData struct {
	Name, ExeName string
	AppID         string
	Version       string
	Build         int
	Category      string
}

func (p *Packager) packageDarwin() (err error) {
	appDir := util.EnsureSubDir(p.dir, p.name+".app")
	exeName := filepath.Base(p.exe)

	contentsDir := util.EnsureSubDir(appDir, "Contents")
	info := filepath.Join(contentsDir, "Info.plist")
	infoFile, err := os.Create(info)
	if err != nil {
		return errors.Wrap(err, "Failed to write plist template")
	}
	defer func() {
		if r := infoFile.Close(); r != nil && err == nil {
			err = r
		}
	}()

	tplData := darwinData{Name: p.name, ExeName: exeName, AppID: p.appID, Version: p.appVersion, Build: p.appBuild,
		Category: strings.ToLower(p.category)}
	if err := templates.InfoPlistDarwin.Execute(infoFile, tplData); err != nil {
		return errors.Wrap(err, "Failed to write plist template")
	}

	macOSDir := util.EnsureSubDir(contentsDir, "MacOS")
	binName := filepath.Join(macOSDir, exeName)
	if err := util.CopyExeFile(p.exe, binName); err != nil {
		return errors.Wrap(err, "Failed to copy exe file")
	}

	resDir := util.EnsureSubDir(contentsDir, "Resources")
	icnsPath := filepath.Join(resDir, "icon.icns")

	img, err := os.Open(p.icon)
	if err != nil {
		return errors.Wrapf(err, "Failed to open source image \"%s\"", p.icon)
	}
	defer img.Close()
	srcImg, _, err := image.Decode(img)
	if err != nil {
		return errors.Wrapf(err, "Failed to decode source image")
	}
	dest, err := os.Create(icnsPath)
	if err != nil {
		return errors.Wrap(err, "Failed to open destination file")
	}
	defer func() {
		if r := dest.Close(); r != nil && err == nil {
			err = r
		}
	}()
	if err := icns.Encode(dest, srcImg); err != nil {
		return errors.Wrap(err, "Failed to encode icns")
	}

	return nil
}
