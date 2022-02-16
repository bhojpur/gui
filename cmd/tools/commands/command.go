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

// It provides functionality for managing Bhojpur GUI application and the build process

import "github.com/bhojpur/gui/cmd/tools/internal/commands"

// Command defines the required functionality to provide a subcommand to the "guiutl" tool.
//
// Deprecated: Access to the internal cli commands are being removed.
// Better versions of bundler, installer, packager and releaser will be avaliable in the future.
type Command interface {
	AddFlags()
	PrintHelp(string)
	Run(args []string)
}

//Getter is the command that can handle downloading and installing Bhojpur GUI applications to the current platform.
type Getter = commands.Getter

// NewGetter returns a command that can handle the download and install of GUI apps built using Bhojpur GUI.
// It depends on a Go and C compiler installed at this stage and takes a single, package, parameter to identify the app.
func NewGetter() *Getter {
	return &Getter{}
}

// NewBundler returns a command that can bundle resources into Go code.
//
// Deprecated: A better version will be exposed in the future.
func NewBundler() Command {
	return &commands.Bundler{}
}

// NewInstaller returns an install command that can install locally built Bhojpur GUI applications.
//
// Deprecated: A better version will be exposed in the future.
func NewInstaller() Command {
	return &commands.Installer{}
}

// NewPackager returns a packager command that can wrap executables into full Bhojpur GUI app packages.
//
// Deprecated: A better version will be exposed in the future.
func NewPackager() Command {
	return &commands.Packager{}
}

// NewReleaser returns a command that can adapt app packages for distribution.
//
// Deprecated: A better version will be exposed in the future.
func NewReleaser() Command {
	return &commands.Releaser{}
}
