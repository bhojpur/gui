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
	"errors"
	"flag"
	"fmt"
	_ "image/jpeg" // import image encodings
	_ "image/png"  // import image encodings
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"

	"github.com/bhojpur/gui/cmd/tools/internal/metadata"
)

const (
	defaultAppBuild   = 1
	defaultAppVersion = "1.0.0"
)

// Package returns the CLI command for packaging Bhojpur GUI applications
func Package() *cli.Command {
	p := &Packager{}

	return &cli.Command{
		Name:        "package",
		Usage:       "Packages a Bhojpur GUI application for distribution.",
		Description: "You may specify the -executable to package, otherwise -sourceDir will be built.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "target",
				Aliases:     []string{"os"},
				Usage:       "The mobile platform to target (android, android/arm, android/arm64, android/amd64, android/386, ios, iossimulator).",
				Destination: &p.os,
			},
			&cli.StringFlag{
				Name:        "executable",
				Aliases:     []string{"exe"},
				Usage:       "The path to the executable, default is the current dir main binary",
				Destination: &p.exe,
			},
			&cli.StringFlag{
				Name:        "name",
				Usage:       "The name of the application, default is the executable file name",
				Destination: &p.name,
			},
			&cli.StringFlag{
				Name:        "tags",
				Usage:       "A comma-separated list of build tags.",
				Destination: &p.tags,
			},
			&cli.StringFlag{
				Name:        "appVersion",
				Usage:       "Version number in the form x, x.y or x.y.z semantic version",
				Destination: &p.appVersion,
			},
			&cli.IntFlag{
				Name:        "appBuild",
				Usage:       "Build number, should be greater than 0 and incremented for each build",
				Destination: &p.appBuild,
			},
			&cli.StringFlag{
				Name:        "sourceDir",
				Aliases:     []string{"src"},
				Usage:       "The directory to package, if executable is not set.",
				Destination: &p.srcDir,
			},
			&cli.StringFlag{
				Name:        "icon",
				Usage:       "The name of the application icon file.",
				Value:       "",
				Destination: &p.icon,
			},
			&cli.StringFlag{
				Name:        "appID",
				Aliases:     []string{"id"},
				Usage:       "For Android, darwin, iOS and Windows targets an appID in the form of a reversed domain name is required, for ios this must match a valid provisioning profile",
				Destination: &p.appID,
			},
			&cli.StringFlag{
				Name:        "certificate",
				Aliases:     []string{"cert"},
				Usage:       "iOS/macOS/Windows: name of the certificate to sign the build",
				Destination: &p.certificate,
			},
			&cli.StringFlag{
				Name:        "profile",
				Usage:       "iOS/macOS: name of the provisioning profile for this build",
				Destination: &p.profile,
			},
			&cli.BoolFlag{
				Name:        "release",
				Usage:       "Enable installation in release mode (disable debug etc).",
				Destination: &p.release,
			},
		},
		Action: func(_ *cli.Context) error {
			return p.Package()
		},
	}
}

// Packager wraps executables into full GUI app packages.
type Packager struct {
	name, srcDir, dir, exe, icon string
	os, appID, appVersion        string
	appBuild                     int
	install, release             bool
	certificate, profile         string // optional flags for releasing
	tags, category               string
}

// AddFlags adds the flags for interacting with the package command.
//
// Deprecated: Access to the individual cli commands are being removed.
func (p *Packager) AddFlags() {
	flag.StringVar(&p.os, "os", "", "The operating system to target (android, android/arm, android/arm64, android/amd64, android/386, darwin, freebsd, ios, linux, netbsd, openbsd, windows, wasm)")
	flag.StringVar(&p.exe, "executable", "", "Specify an existing binary instead of building before package")
	flag.StringVar(&p.srcDir, "sourceDir", "", "The directory to package, if executable is not set")
	flag.StringVar(&p.name, "name", "", "The name of the application, default is the executable file name")
	flag.StringVar(&p.icon, "icon", "", "The name of the application icon file")
	flag.StringVar(&p.appID, "appID", "", "For ios or darwin targets an appID is required, for ios this must \nmatch a valid provisioning profile")
	flag.StringVar(&p.appVersion, "appVersion", "", "Version number in the form x, x.y or x.y.z semantic version")
	flag.IntVar(&p.appBuild, "appBuild", 0, "Build number, should be greater than 0 and incremented for each build")
	flag.BoolVar(&p.release, "release", false, "Should this package be prepared for release? (disable debug etc)")
	flag.StringVar(&p.tags, "tags", "", "A comma-separated list of build tags")
}

// PrintHelp prints the help for the package command.
//
// Deprecated: Access to the individual cli commands are being removed.
func (*Packager) PrintHelp(indent string) {
	fmt.Println(indent, "The package command prepares an application for installation and testing.")
	fmt.Println(indent, "You may specify the -executable to package, otherwise -sourceDir will be built.")
	fmt.Println(indent, "Command usage: guiutl package [parameters]")
}

// Run runs the package command.
//
// Deprecated: A better version will be exposed in the future.
func (p *Packager) Run(_ []string) {
	err := p.validate()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	err = p.doPackage(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

// Package starts the packaging process
func (p *Packager) Package() error {
	err := p.validate()
	if err != nil {
		return err
	}

	return p.packageWithoutValidate()
}

func (p *Packager) packageWithoutValidate() error {
	err := p.doPackage(nil)
	if err != nil {
		return err
	}

	data, err := metadata.LoadStandard(p.srcDir)
	if err != nil {
		return nil // no metadata to update
	}
	data.Details.Build++
	return metadata.SaveStandard(data, p.srcDir)
}

func (p *Packager) buildPackage(runner runner) ([]string, error) {
	var tags []string
	if p.tags != "" {
		tags = strings.Split(p.tags, ",")
	}
	if p.os != "web" {
		b := &builder{
			os:      p.os,
			srcdir:  p.srcDir,
			target:  p.exe,
			release: p.release,
			tags:    tags,
			runner:  runner,
		}

		return []string{p.exe}, b.build()
	}

	bWasm := &builder{
		os:      "wasm",
		srcdir:  p.srcDir,
		target:  p.exe + ".wasm",
		release: p.release,
		tags:    tags,
		runner:  runner,
	}

	err := bWasm.build()
	if err != nil {
		return nil, err
	}

	bGopherJS := &builder{
		os:      "gopherjs",
		srcdir:  p.srcDir,
		target:  p.exe + ".js",
		release: p.release,
		tags:    tags,
		runner:  runner,
	}

	err = bGopherJS.build()
	if err != nil {
		return nil, err
	}

	return []string{bWasm.target, bGopherJS.target}, nil
}

func (p *Packager) combinedVersion() string {
	return fmt.Sprintf("%s.%d", p.appVersion, p.appBuild)
}

func (p *Packager) doPackage(runner runner) error {
	// sensible defaults - validation deemed them optional
	if p.appVersion == "" {
		p.appVersion = defaultAppVersion
	}
	if p.appBuild <= 0 {
		p.appBuild = defaultAppBuild
	}

	if !util.Exists(p.exe) && !util.IsMobile(p.os) {
		files, err := p.buildPackage(runner)
		if err != nil {
			return fmt.Errorf("error building application: %w", err)
		}
		for _, file := range files {
			if p.os != "web" && !util.Exists(file) {
				return fmt.Errorf("unable to build directory to expected executable, %s", file)
			}
		}
		if p.os != "windows" {
			defer p.removeBuild(files)
		}
	}

	switch p.os {
	case "darwin":
		return p.packageDarwin()
	case "linux", "openbsd", "freebsd", "netbsd":
		return p.packageUNIX()
	case "windows":
		return p.packageWindows()
	case "android/arm", "android/arm64", "android/amd64", "android/386", "android":
		return p.packageAndroid(p.os)
	case "ios", "iossimulator":
		return p.packageIOS(p.os)
	case "wasm":
		return p.packageWasm()
	case "gopherjs":
		return p.packageGopherJS()
	case "web":
		return p.packageWeb()
	default:
		return fmt.Errorf("unsupported target operating system \"%s\"", p.os)
	}
}

func (p *Packager) removeBuild(files []string) {
	for _, file := range files {
		err := os.RemoveAll(file)
		if err != nil {
			log.Println("Unable to remove temporary build file", p.exe)
		}
	}
}

func (p *Packager) validate() error {
	if p.os == "" {
		p.os = targetOS()
	}
	baseDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to get the current directory, needed to find main executable: %w", err)
	}
	if p.dir == "" {
		p.dir = baseDir
	}
	if p.srcDir == "" {
		p.srcDir = baseDir
	} else if p.os == "ios" || p.os == "android" {
		return errors.New("parameter -sourceDir is currently not supported for mobile builds. " +
			"Change directory to the main package and try again")
	}

	data, err := metadata.LoadStandard(p.srcDir)
	if err == nil {
		mergeMetadata(p, data)
	}

	exeName := calculateExeName(p.srcDir, p.os)

	if p.exe == "" {
		p.exe = filepath.Join(p.srcDir, exeName)

		if util.Exists(p.exe) { // the exe was not specified, assume stale
			p.removeBuild([]string{p.exe})
		}
	} else if p.os == "ios" || p.os == "android" {
		_, _ = fmt.Fprint(os.Stderr, "Parameter -executable is ignored for mobile builds.\n")
	}

	if p.name == "" || p.os == "wasm" || p.os == "gopherjs" || p.os == "web" {
		p.name = exeName
	}
	if p.icon == "" || p.icon == "Icon.png" {
		p.icon = filepath.Join(p.srcDir, "Icon.png")
	}
	if !util.Exists(p.icon) {
		return errors.New("Missing application icon at \"" + p.icon + "\"")
	}

	p.appID, err = validateAppID(p.appID, p.os, p.name, p.release)
	if err != nil {
		return err
	}
	if p.appVersion != "" && !isValidVersion(p.appVersion) {
		return errors.New("invalid -appVersion parameter, integer and '.' characters only up to x.y.z")
	}

	return nil
}

func calculateExeName(sourceDir, os string) string {
	exeName := filepath.Base(sourceDir)
	/* #nosec */
	if data, err := ioutil.ReadFile(filepath.Join(sourceDir, "go.mod")); err == nil {
		modulePath := modfile.ModulePath(data)
		moduleName, _, ok := module.SplitPathVersion(modulePath)
		if ok {
			paths := strings.Split(moduleName, "/")
			name := paths[len(paths)-1]
			if name != "" {
				exeName = name
			}
		}
	}

	if os == "windows" {
		exeName = exeName + ".exe"
	} else if os == "wasm" {
		exeName = exeName + ".wasm"
	} else if os == "gopherjs" {
		exeName = exeName + ".js"
	}

	return exeName
}

func isValidVersion(ver string) bool {
	nums := strings.Split(ver, ".")
	if len(nums) == 0 || len(nums) > 3 {
		return false
	}
	for _, num := range nums {
		if _, err := strconv.Atoi(num); err != nil {
			return false
		}
	}
	return true
}

func mergeMetadata(p *Packager, data *metadata.BhojpurApp) {
	if p.icon == "" {
		p.icon = data.Details.Icon
	}
	if p.name == "" {
		p.name = data.Details.Name
	}
	if p.appID == "" {
		p.appID = data.Details.ID
	}
	if p.appVersion == "" {
		p.appVersion = data.Details.Version
	}
	if p.appBuild == 0 {
		p.appBuild = data.Details.Build
	}
}

func validateAppID(appID, os, name string, release bool) (string, error) {
	// old darwin compatibility
	if os == "darwin" {
		if appID == "" {
			return "com.example." + name, nil
		}
	} else if os == "ios" || util.IsAndroid(os) || (os == "windows" && release) {
		// all mobile, and for windows when releasing, needs a unique id - usually reverse DNS style
		if appID == "" {
			return "", errors.New("missing appID parameter for package")
		} else if !strings.Contains(appID, ".") {
			return "", errors.New("appID must be globally unique and contain at least 1 '.'")
		}
	}

	return appID, nil
}
