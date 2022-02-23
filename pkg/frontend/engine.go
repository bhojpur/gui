package frontend

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
	"io"
	"log"
	"net"

	"github.com/bhojpur/gui/pkg/engine/widget"
)

func ConfirmSend(form *widget.Form, btnSave *widget.Button, lblSuccess *widget.Label, lblError *widget.Label, data string) {
	localhost := defaultConfig.GetDefaultServer()
	if err := sendData(localhost, data); err != nil {
		log.Println(err)
		lblError.SetText(err.Error())
	} else {
		lblSuccess.SetText("Saved successfully")
	}
}

//sendData sends data to backend server
func sendData(addr, data string) error {
	// connect with server
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("Error connecting to backend server: %v\n", err)
	}

	// send data to backend server
	if _, err := conn.Write([]byte(data)); err != nil {
		return fmt.Errorf("Error sending data to backend server: %v\n", err)
	}

	log.Printf("Send data: %v\n", data)

	resp := make([]byte, 512)
	respLen, err := conn.Read(resp)
	if err != nil && err != io.EOF {
		return fmt.Errorf("Response error: %v\n", err)
	}

	log.Println("Packet processing: ", string(resp[:respLen]))

	// close connection
	return conn.Close()
}
