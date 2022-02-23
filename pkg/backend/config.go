package backend

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
	"log"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/bhojpur/configure/pkg/toml"
)

type (
	hosting struct {
		Title      string
		Desc       string
		Integers   []int
		Floats     []float64
		Times      []fmtTime
		Duration   []duration
		Distros    []distro
		Servers    map[string]server
		Characters map[string][]struct {
			Name string
			Rank string
		}
		Log string
	}

	server struct {
		IP       string
		Hostname string
		Enabled  bool
	}

	distro struct {
		Name     string
		Packages string
	}

	duration struct{ time.Duration }
	fmtTime  struct{ time.Time }
)

var defaultConfig hosting

func (d *duration) UnmarshalText(text []byte) (err error) {
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func (t fmtTime) String() string {
	f := "2006-01-02 15:04:05.999999999"
	if t.Time.Hour() == 0 {
		f = "2006-01-02"
	}
	if t.Time.Year() == 0 {
		f = "15:04:05.999999999"
	}
	if t.Time.Location() == time.UTC {
		f += " UTC"
	} else {
		f += " +0530"
	}
	return t.Time.Format(`"` + f + `"`)
}

func (c *hosting) GetDefaultServer() string {
	server := c.Servers["servers.dev"]

	return server.Hostname
}

func (c *hosting) GetLogFileName() string {
	logname := c.Log

	return logname
}

//LoadConfig load configuration from file
func LoadConfig(confPath string) error {
	f := "server.toml"
	if _, err := os.Stat(confPath); err != nil {
		f = "./conf/server.toml"
	}
	meta, err := toml.DecodeFile(f, &defaultConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	indent := strings.Repeat(" ", 14)

	fmt.Print("Decoded")
	typ, val := reflect.TypeOf(defaultConfig), reflect.ValueOf(defaultConfig)
	for i := 0; i < typ.NumField(); i++ {
		indent := indent
		if i == 0 {
			indent = strings.Repeat(" ", 7)
		}
		fmt.Printf("%s%-11s → %v\n", indent, typ.Field(i).Name, val.Field(i).Interface())
	}

	fmt.Print("\nKeys")
	keys := meta.Keys()
	sort.Slice(keys, func(i, j int) bool { return keys[i].String() < keys[j].String() })
	for i, k := range keys {
		indent := indent
		if i == 0 {
			indent = strings.Repeat(" ", 10)
		}
		fmt.Printf("%s%-10s %s\n", indent, meta.Type(k...), k)
	}

	fmt.Print("\nUndecoded")
	keys = meta.Undecoded()
	sort.Slice(keys, func(i, j int) bool { return keys[i].String() < keys[j].String() })
	for i, k := range keys {
		indent := indent
		if i == 0 {
			indent = strings.Repeat(" ", 5)
		}
		fmt.Printf("%s%-10s %s\n", indent, meta.Type(k...), k)
	}

	//create log file
	logfname := defaultConfig.GetLogFileName()
	logFile, err := os.OpenFile(logfname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	return err
}
