package binding

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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBindUserType(t *testing.T) {
	val := user{name: "Unnamed"}
	u := bindUserType(&val)
	v, err := u.Get()
	assert.Nil(t, err)
	assert.Equal(t, "User: Unnamed", v.String())

	called := false
	fn := NewDataListener(func() {
		called = true
	})
	u.AddListener(fn)
	waitForItems()
	assert.True(t, called)

	called = false
	err = u.Set(user{name: "Replace"})
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, "User: Replace", val.String())
	assert.True(t, called)

	called = false
	val = user{name: "Direct"}
	_ = u.Reload()
	waitForItems()
	assert.True(t, called)
	v, err = u.Get()
	assert.Nil(t, err)
	assert.Equal(t, "User: Direct", v.String())

	called = false
	val.name = "FieldSet"
	_ = u.Reload()
	waitForItems()
	assert.True(t, called)
	v, err = u.Get()
	assert.Nil(t, err)
	assert.Equal(t, "User: FieldSet", v.String())
}

func TestNewUserType(t *testing.T) {
	u := newUserType()
	v, err := u.Get()
	assert.Nil(t, err)
	assert.Equal(t, "User: ", v.String())

	err = u.Set(user{name: "Dave"})
	assert.Nil(t, err)
	v, err = u.Get()
	assert.Nil(t, err)
	assert.Equal(t, "User: Dave", v.String())
}

type user struct {
	name string
}

func (u *user) String() string {
	return "User: " + u.name
}

type userType struct {
	Untyped
}

func newUserType() *userType {
	ret := &userType{Untyped: NewUntyped()}
	_ = ret.Set(user{})
	return ret
}

func (t *userType) Get() (user, error) {
	val, err := t.Untyped.Get()
	return val.(user), err
}

func (t *userType) Set(u user) error {
	return t.Untyped.Set(u)
}

type externalUserType struct {
	ExternalUntyped
}

func bindUserType(u *user) *externalUserType {
	return &externalUserType{ExternalUntyped: BindUntyped(u)}
}

func (t *externalUserType) Get() (user, error) {
	val, err := t.ExternalUntyped.Get()
	return val.(user), err
}

func (t *externalUserType) Set(u user) error {
	return t.ExternalUntyped.Set(u)
}
