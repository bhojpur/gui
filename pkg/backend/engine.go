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
	"io"
	"log"
	"net"
	"strings"
)

type Server struct {
	l net.Listener
}

//Init initialize Bhojpur GUI server instance
func (s *Server) Init(addr string) error {
	var err error

	s.l, err = net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("Error start server: %v\n", err)
	}
	log.Println("Bhojpur GUI - Server Listener start ", addr)

	return err
}

//Run start listen connection on server
func (s *Server) Run() {
	for {
		conn, err := s.l.Accept()
		if err != nil {
			log.Printf("Accept error: %v\n", err)
		} else {
			go s.handlerConn(conn)
		}
	}
}

//handlerConn handler incoming connection
func (s *Server) handlerConn(c net.Conn) {
	defer c.Close()
	buf := make([]byte, 2048)
	rcvPacketSize, err := c.Read(buf)
	if err != nil && err != io.EOF {
		log.Println("Read error: ", err)
		return
	}
	data := buf[:rcvPacketSize]

	rec := strings.Split(string(data), " ")
	log.Println("Received data from client: ", rec)

	// rec must have 3 field (as at form)
	if len(rec) <= 3 {
		log.Printf("Forward record in target system: %v\n", rec)
	}
}
