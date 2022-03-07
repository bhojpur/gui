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
	"fmt"
	"image"
	"os"
	"path/filepath"
	"runtime"

	ico "github.com/Kodeworks/golang-image-ico"
	"github.com/bhojpur/gui/cmd/tools/internal/templates"
	"github.com/josephspurrier/goversioninfo"
	"golang.org/x/sys/execabs"
)

type windowsData struct {
	Name            string
	CombinedVersion string
}

func (p *Packager) packageWindows() error {
	exePath := filepath.Dir(p.exe)

	// convert icon
	img, err := os.Open(p.icon)
	if err != nil {
		return fmt.Errorf("failed to open source image: %w", err)
	}
	defer img.Close()
	srcImg, _, err := image.Decode(img)
	if err != nil {
		return fmt.Errorf("failed to decode source image: %w", err)
	}

	icoPath := filepath.Join(exePath, p.name+".ico")
	file, err := os.Create(icoPath)
	if err != nil {
		return fmt.Errorf("failed to open image file: %w", err)
	}

	err = ico.Encode(file, srcImg)
	if err != nil {
		return fmt.Errorf("failed to encode icon: %w", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("failed to close image file: %w", err)
	}

	// write manifest
	manifest := p.exe + ".manifest"
	manifestGenerated := false
	if _, err := os.Stat(manifest); os.IsNotExist(err) {
		manifestGenerated = true
		manifestFile, _ := os.Create(manifest)

		tplData := windowsData{
			Name:            p.name,
			CombinedVersion: p.combinedVersion(),
		}
		err := templates.ManifestWindows.Execute(manifestFile, tplData)
		if err != nil {
			return fmt.Errorf("failed to write manifest template: %w", err)
		}
	}

	// launch rsrc to generate the object file
	outPath := filepath.Join(exePath, "bhojpur.syso")

	vi := &goversioninfo.VersionInfo{}
	vi.ProductName = p.name
	vi.IconPath = icoPath
	vi.ManifestPath = manifest

	vi.Build()
	vi.Walk()

	arch, ok := os.LookupEnv("GOARCH")
	if !ok {
		arch = runtime.GOARCH
	}

	err = vi.WriteSyso(outPath, arch)
	if err != nil {
		return fmt.Errorf("failed to write .syso file: %w", err)
	}
	defer os.Remove(outPath)

	err = os.Remove(icoPath)
	if err != nil {
		return fmt.Errorf("failed to remove icon: %w", err)
	} else if manifestGenerated {
		err := os.Remove(manifest)
		if err != nil {
			return fmt.Errorf("failed to remove manifest: %w", err)
		}
	}

	_, err = p.buildPackage(nil)
	if err != nil {
		return fmt.Errorf("failed to rebuild after adding metadata: %w", err)
	}

	appPath := p.exe
	appName := filepath.Base(p.exe)
	if filepath.Base(p.exe) != p.name {
		appName = p.name
		if filepath.Ext(p.name) != ".exe" {
			appName = appName + ".exe"
		}
		appPath = filepath.Join(filepath.Dir(p.exe), appName)
		os.Rename(filepath.Base(p.exe), appName)
	}

	if p.install {
		err := runAsAdminWindows("copy", "\"\""+appPath+"\"\"", "\"\""+filepath.Join(p.dir, appName)+"\"\"")
		if err != nil {
			return fmt.Errorf("failed to run as administrator: %w", err)
		}
	}
	return nil
}

func runAsAdminWindows(args ...string) error {
	cmd := "\"/c\""

	for _, arg := range args {
		cmd += ",\"" + arg + "\""
	}

	return execabs.Command("powershell.exe", "Start-Process", "cmd.exe", "-Verb", "runAs", "-ArgumentList", cmd).Run()
}
