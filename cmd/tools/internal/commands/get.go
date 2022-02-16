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
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bhojpur/gui/cmd/tools/internal/util"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"golang.org/x/sys/execabs"
)

// Get returns the command which downloads and installs Bhojpur GUI applications.
func Get() *cli.Command {
	g := &Getter{}
	return &cli.Command{
		Name:        "get",
		Usage:       "Downloads and installs a Bhojpur GUI application",
		Description: "A single parameter is required to specify the Go package, as with \"go get\".",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "icon",
				Usage:       "The name of the application icon file.",
				Value:       "",
				Destination: &g.icon,
			},
			&cli.StringFlag{
				Name:        "appID",
				Aliases:     []string{"id"},
				Usage:       "For darwin and Windows targets an appID in the form of a reversed domain name is required, for ios this must match a valid provisioning profile",
				Destination: &g.appID,
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() != 1 {
				return errors.New("missing \"package\" argument to define the application to get")
			}

			pkg := ctx.Args().Slice()[0]
			return g.Get(pkg)
		},
	}
}

// Getter is the command that can handle downloading and installing Bhojpur GUI apps to the current platform.
type Getter struct {
	icon, appID string
}

// NewGetter returns a command that can handle the download and install of GUI apps built using Bhojpur GUI.
// It depends on a Go and C compiler installed at this stage and takes a single, package, parameter to identify the app.
func NewGetter() *Getter {
	return &Getter{}
}

// Get automates the download and install of a named Bhojpur GUI application package.
func (g *Getter) Get(pkg string) error {
	cmd := execabs.Command("go", "get", "-u", "-d", pkg)
	cmd.Env = append(os.Environ(), "GO111MODULE=off") // cache the downloaded code
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	path := filepath.Join(goPath(), "src", pkg)
	if !util.Exists(path) { // the error above may be ignorable, unless the path was not found
		return err
	}

	install := &Installer{srcDir: path, icon: g.icon, appID: g.appID, release: true}
	if err := install.validate(); err != nil {
		return errors.Wrap(err, "Failed to set up installer")
	}

	return install.install()
}

// SetAppID Allows the Get operation to specify an appID for the application that is being downloaded.
//
// Since: 2.1
func (g *Getter) SetAppID(id string) {
	g.appID = id
}

// SetIcon allows you to set the app icon path that will be used for the next Get operation.
func (g *Getter) SetIcon(path string) {
	g.icon = path
}

// AddFlags adds available flags to the current flags parser
//
// Deprecated: Get does not define any flags.
func (g *Getter) AddFlags() {
	flag.StringVar(&g.icon, "icon", "Icon.png", "The name of the application icon file")
}

// PrintHelp prints help for this command when used in a command-line context
//
// Deprecated: Use Get() to get the cli and help messages instead.
func (g *Getter) PrintHelp(indent string) {
	fmt.Println(indent, "The get command downloads and installs a Bhojpur GUI application.")
	fmt.Println(indent, "A single parameter is required to specify the Go package, as with \"go get\"")
}

// Run is the command-line version of the Get(pkg) command, the first of the passed arguments will be used
// as the package to get.
//
// Deprecated: Use Get() for the urfave/cli command or Getter.Get() to download and install an application.
func (g *Getter) Run(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Missing \"package\" argument to define the application to get")
		os.Exit(1)
	}

	if err := g.Get(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func goPath() string {
	cmd := execabs.Command("go", "env", "GOPATH")
	out, err := cmd.CombinedOutput()

	if err != nil {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "go")
	}

	return strings.TrimSpace(string(out[0 : len(out)-1]))
}
