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

#ifdef os_android
// TODO: We could include <android/api-level.h> and
// condition on __ANDROID_API__ to get GLES3 headers. However
// we also need to add -lGLESv3 to LDFLAGS, which we cannot do
// from inside an ifdef.
#include <GLES2/gl2.h>
#elif os_linux
#include <GLES3/gl3.h> // install on Ubuntu with: sudo apt-get install libegl1-mesa-dev libgles2-mesa-dev libx11-dev
#elif os_openbsd
#include <GLES3/gl3.h>
#elif os_freebsd
#include <GLES3/gl3.h>
#endif

#ifdef os_ios
#include <OpenGLES/ES2/glext.h>
#endif

#ifdef os_macos
#include <OpenGL/gl3.h>
#define GL_ES_VERSION_3_0 1
#endif

#if defined(GL_ES_VERSION_3_0) && GL_ES_VERSION_3_0
#define GLES_VERSION "GL_ES_3_0"
#else
#define GLES_VERSION "GL_ES_2_0"
#endif

#include <stdint.h>
#include <stdlib.h>

// TODO: generate this enum from fn.go.
typedef enum {
	glfnUNDEFINED,
	glfnActiveTexture,
	glfnAttachShader,
	glfnBindBuffer,
	glfnBindTexture,
	glfnBindVertexArray,
	glfnBlendColor,
	glfnBlendFunc,
	glfnBufferData,
	glfnClear,
	glfnClearColor,
	glfnCompileShader,
	glfnCreateProgram,
	glfnCreateShader,
	glfnDeleteBuffer,
	glfnDeleteTexture,
	glfnDisable,
	glfnDrawArrays,
	glfnEnable,
	glfnEnableVertexAttribArray,
	glfnFlush,
	glfnGenBuffer,
	glfnGenTexture,
	glfnGenVertexArray,
	glfnGetAttribLocation,
	glfnGetError,
	glfnGetShaderInfoLog,
	glfnGetShaderSource,
	glfnGetShaderiv,
	glfnGetTexParameteriv,
	glfnGetUniformLocation,
	glfnLinkProgram,
	glfnReadPixels,
	glfnScissor,
	glfnShaderSource,
	glfnTexImage2D,
	glfnTexParameteri,
	glfnUniform1f,
	glfnUniform4f,
	glfnUniform4fv,
	glfnUseProgram,
	glfnVertexAttribPointer,
	glfnViewport,
} glfn;

// TODO: generate this type from fn.go.
struct fnargs {
	glfn fn;

	uintptr_t a0;
	uintptr_t a1;
	uintptr_t a2;
	uintptr_t a3;
	uintptr_t a4;
	uintptr_t a5;
	uintptr_t a6;
	uintptr_t a7;
	uintptr_t a8;
	uintptr_t a9;
};

extern uintptr_t processFn(struct fnargs* args, char* parg);
