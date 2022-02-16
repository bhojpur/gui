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
	"os"
	"runtime/debug"

	"github.com/lucor/goinfo"
	"github.com/lucor/goinfo/format"
	"github.com/lucor/goinfo/report"
	"github.com/urfave/cli/v2"
)

const guiModule = "github.com/bhojpur/gui/pkg/engine"

// Env returns the env command
func Env() *cli.Command {
	return &cli.Command{
		Name:  "env",
		Usage: "Prints Bhojpur GUI application module and environment information",
		Action: func(_ *cli.Context) error {
			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("could not get the path for the current working dir: %v", err)
			}

			reporters := []goinfo.Reporter{
				&guiReport{GoMod: &report.GoMod{WorkDir: workDir, Module: guiModule}},
				&report.GoVersion{},
				&report.GoEnv{Filter: []string{"GOOS", "GOARCH", "CGO_ENABLED", "GO111MODULE"}},
				&report.OS{},
			}

			err = goinfo.Write(os.Stdout, reporters, &format.Text{})
			if err != nil {
				return err
			}

			return nil
		},
	}
}

// guiReport defines a custom report for Bhojpur GUI application
type guiReport struct {
	*report.GoMod
}

// Info returns the collected info
func (r *guiReport) Info() (goinfo.Info, error) {
	info, err := r.GoMod.Info()
	if err != nil {
		return info, err
	}
	// remove info for the report
	delete(info, "module")

	binfo, ok := debug.ReadBuildInfo()
	if !ok {
		info["cli_version"] = "could not retrieve version information (ensure module support is activated and build again)"
	} else {
		info["cli_version"] = binfo.Main.Version
	}

	return info, nil
}
