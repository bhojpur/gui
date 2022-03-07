package util

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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// AndroidBuildToolsPath tries to find the location of the "build-tools" directory.
// This depends on ANDROID_HOME variable, so you should call RequireAndroidSDK() first.
func AndroidBuildToolsPath() string {
	env := os.Getenv("ANDROID_HOME")
	dir := filepath.Join(env, "build-tools")

	// this may contain a version subdir
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return dir
	}

	childDir := ""
	for _, f := range files {
		if f.Name() == "zipalign" { // no version subdir
			return dir
		}
		if f.IsDir() && childDir == "" {
			childDir = f.Name()
		}
	}

	if childDir == "" {
		return dir
	}
	return filepath.Join(dir, childDir)
}

// IsAndroid returns true if the given os parameter represents one of the Android targets.
func IsAndroid(os string) bool {
	return strings.HasPrefix(os, "android")
}

// IsIOS returns true if the given os parameter represents one of the iOS targets (ios, iossimulator)
func IsIOS(os string) bool {
	return strings.HasPrefix(os, "ios")
}

// IsMobile returns true if the given os parameter represents a platform handled by gomobile.
func IsMobile(os string) bool {
	return IsIOS(os) || IsAndroid(os)
}

// RequireAndroidSDK will return an error if it cannot establish the location of a valid Android SDK installation.
// This is currently deduced using ANDROID_HOME environment variable.
func RequireAndroidSDK() error {
	if env, ok := os.LookupEnv("ANDROID_HOME"); !ok || env == "" {
		return fmt.Errorf("could not find android tools, missing ANDROID_HOME")
	}

	return nil
}
