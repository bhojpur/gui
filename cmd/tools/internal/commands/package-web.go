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
	"path/filepath"
	"runtime"

	"github.com/bhojpur/gui/cmd/tools/internal/templates"
)

func (p *Packager) packageWeb() error {
	appDir := util.EnsureSubDir(p.dir, "web")

	tpl := webData{
		GopherJSFile: p.name + ".js",
		WasmFile:     p.name + ".wasm",
		IsReleased:   p.release,
		HasGopherJS:  true,
		HasWasm:      true,
	}

	return tpl.packageWebInternal(appDir, p.exe+".wasm", p.exe+".js", p.icon, p.release)
}

func (p *Packager) packageWasm() error {
	appDir := util.EnsureSubDir(p.dir, "wasm")

	tpl := webData{
		WasmFile:   p.name,
		IsReleased: p.release,
		HasWasm:    true,
	}

	return tpl.packageWebInternal(appDir, p.exe, "", p.icon, p.release)
}

func (p *Packager) packageGopherJS() error {
	appDir := util.EnsureSubDir(p.dir, "gopherjs")

	tpl := webData{
		GopherJSFile: p.name,
		IsReleased:   p.release,
		HasGopherJS:  true,
	}

	return tpl.packageWebInternal(appDir, "", p.exe, p.icon, p.release)
}

type webData struct {
	WasmFile     string
	GopherJSFile string
	IsReleased   bool
	HasWasm      bool
	HasGopherJS  bool
}

func (w webData) packageWebInternal(appDir string, exeWasmSrc string, exeJSSrc string, icon string, release bool) error {
	var tpl bytes.Buffer
	err := templates.IndexHTML.Execute(&tpl, w)
	if err != nil {
		return err
	}

	index := filepath.Join(appDir, "index.html")
	err = util.WriteFile(index, tpl.Bytes())
	if err != nil {
		return err
	}

	iconDst := filepath.Join(appDir, "icon.png")
	err = util.CopyFile(icon, iconDst)
	if err != nil {
		return err
	}

	if w.HasGopherJS {
		exeJSDst := filepath.Join(appDir, w.GopherJSFile)
		err = util.CopyFile(exeJSSrc, exeJSDst)
		if err != nil {
			return err
		}
	}

	if w.HasWasm {
		wasmExecSrc := filepath.Join(runtime.GOROOT(), "misc", "wasm", "wasm_exec.js")
		wasmExecDst := filepath.Join(appDir, "wasm_exec.js")
		err = util.CopyFile(wasmExecSrc, wasmExecDst)
		if err != nil {
			return err
		}

		exeWasmDst := filepath.Join(appDir, w.WasmFile)
		err = util.CopyFile(exeWasmSrc, exeWasmDst)
		if err != nil {
			return err
		}
	}

	// Download webgl-debug.js directly from the KhronosGroup repository when needed
	if !release {
		webglDebugFile := filepath.Join(appDir, "webgl-debug.js")
		err := util.WriteFile(webglDebugFile, templates.WebGLDebugJs)
		if err != nil {
			return err
		}
	}

	return nil
}
