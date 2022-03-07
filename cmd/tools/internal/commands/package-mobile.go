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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/bhojpur/gui/cmd/tools/internal/mobile"
	"github.com/bhojpur/gui/cmd/tools/internal/templates"
	gui "github.com/bhojpur/gui/pkg/engine"
	"golang.org/x/sys/execabs"
)

func (p *Packager) packageAndroid(arch string) error {
	return mobile.RunNewBuild(arch, p.appID, p.icon, p.name, p.appVersion, p.appBuild, p.release, "", "")
}

func (p *Packager) packageIOS(target string) error {
	err := mobile.RunNewBuild(target, p.appID, p.icon, p.name, p.appVersion, p.appBuild, p.release, p.certificate, p.profile)
	if err != nil {
		return err
	}

	assetDir := util.EnsureSubDir(p.dir, "Images.xcassets")
	defer os.RemoveAll(assetDir)
	err = ioutil.WriteFile(filepath.Join(assetDir, "Contents.json"), []byte(`{
  "info" : {
    "author" : "xcode",
    "version" : 1
  }
}`), 0644)
	if err != nil {
		gui.LogError("Content err", err)
	}

	iconDir := util.EnsureSubDir(assetDir, "AppIcon.appiconset")
	contentFile, _ := os.Create(filepath.Join(iconDir, "Contents.json"))

	err = templates.XCAssetsDarwin.Execute(contentFile, nil)
	if err != nil {
		return fmt.Errorf("failed to write xcassets content template: %w", err)
	}

	if err = copyResizeIcon(1024, iconDir, p.icon); err != nil {
		return err
	}
	if err = copyResizeIcon(180, iconDir, p.icon); err != nil {
		return err
	}
	if err = copyResizeIcon(120, iconDir, p.icon); err != nil {
		return err
	}
	if err = copyResizeIcon(76, iconDir, p.icon); err != nil {
		return err
	}
	if err = copyResizeIcon(152, iconDir, p.icon); err != nil {
		return err
	}

	appDir := filepath.Join(p.dir, mobile.AppOutputName(p.os, p.name))
	return runCmdCaptureOutput("xcrun", "actool", "Images.xcassets", "--compile", appDir, "--platform",
		"iphoneos", "--target-device", "iphone", "--minimum-deployment-target", "9.0", "--app-icon", "AppIcon",
		"--output-format", "human-readable-text", "--output-partial-info-plist", "/dev/null")
}

func copyResizeIcon(size int, dir, source string) error {
	strSize := strconv.Itoa(size)
	path := dir + "/Icon_" + strSize + ".png"
	return runCmdCaptureOutput("sips", "-o", path, "-Z", strSize, source)
}

// runCmdCaptureOutput is a exec.Command wrapper that offers better error messages from stdout and stderr.
func runCmdCaptureOutput(name string, args ...string) error {
	var (
		outbuf = &bytes.Buffer{}
		errbuf = &bytes.Buffer{}
	)
	cmd := execabs.Command(name, args...)
	cmd.Stdout = outbuf
	cmd.Stderr = errbuf
	err := cmd.Run()
	if err != nil {
		outstr := outbuf.String()
		errstr := errbuf.String()
		if outstr != "" {
			err = fmt.Errorf(outbuf.String()+": %w", err)
		}
		if errstr != "" {
			err = fmt.Errorf(outbuf.String()+": %w", err)
		}
		return err
	}
	return nil
}
