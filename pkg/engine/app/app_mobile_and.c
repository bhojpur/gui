// +build !ci

// +build android

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

#include <android/log.h>
#include <jni.h>
#include <stdbool.h>
#include <stdlib.h>

#define LOG_FATAL(...) __android_log_print(ANDROID_LOG_FATAL, "Bhojpur", __VA_ARGS__)

static jclass find_class(JNIEnv *env, const char *class_name) {
	jclass clazz = (*env)->FindClass(env, class_name);
	if (clazz == NULL) {
		(*env)->ExceptionClear(env);
		LOG_FATAL("cannot find %s", class_name);
		return NULL;
	}
	return clazz;
}

static jmethodID find_method(JNIEnv *env, jclass clazz, const char *name, const char *sig) {
	jmethodID m = (*env)->GetMethodID(env, clazz, name, sig);
	if (m == 0) {
		(*env)->ExceptionClear(env);
		LOG_FATAL("cannot find method %s %s", name, sig);
		return 0;
	}
	return m;
}

static jmethodID find_static_method(JNIEnv *env, jclass clazz, const char *name, const char *sig) {
	jmethodID m = (*env)->GetStaticMethodID(env, clazz, name, sig);
	if (m == 0) {
		(*env)->ExceptionClear(env);
		LOG_FATAL("cannot find method %s %s", name, sig);
		return 0;
	}
	return m;
}

jobject getSystemService(uintptr_t jni_env, uintptr_t ctx, char *service) {
	JNIEnv *env = (JNIEnv*)jni_env;
	jstring serviceStr = (*env)->NewStringUTF(env, service);

	jclass ctxClass = (*env)->GetObjectClass(env, ctx);
	jmethodID getSystemService = find_method(env, ctxClass, "getSystemService", "(Ljava/lang/String;)Ljava/lang/Object;");

	return (jobject)(*env)->CallObjectMethod(env, ctx, getSystemService, serviceStr);
}

int nextId = 1;

bool isOreoOrLater(JNIEnv *env) {
    jclass versionClass = find_class(env, "android/os/Build$VERSION" );
    jfieldID sdkIntFieldID = (*env)->GetStaticFieldID(env, versionClass, "SDK_INT", "I" );
    int sdkVersion = (*env)->GetStaticIntField(env, versionClass, sdkIntFieldID );

    return sdkVersion >= 26; // O = Oreo, will not be defined for older builds
}

jobject parseURL(uintptr_t jni_env, uintptr_t ctx, char* uriCstr) {
	JNIEnv *env = (JNIEnv*)jni_env;

	jstring uriStr = (*env)->NewStringUTF(env, uriCstr);
	jclass uriClass = find_class(env, "android/net/Uri");
	jmethodID parse = find_static_method(env, uriClass, "parse", "(Ljava/lang/String;)Landroid/net/Uri;");

	return (jobject)(*env)->CallStaticObjectMethod(env, uriClass, parse, uriStr);
}

void openURL(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx, char *url) {
	JNIEnv *env = (JNIEnv*)jni_env;
	jobject uri = parseURL(jni_env, ctx, url);

	jclass intentClass = find_class(env, "android/content/Intent");
	jfieldID viewFieldID = (*env)->GetStaticFieldID(env, intentClass, "ACTION_VIEW", "Ljava/lang/String;" );
    jstring view = (*env)->GetStaticObjectField(env, intentClass, viewFieldID);

	jmethodID constructor = find_method(env, intentClass, "<init>", "(Ljava/lang/String;Landroid/net/Uri;)V");
	jobject intent = (*env)->NewObject(env, intentClass, constructor, view, uri);

	jclass contextClass = find_class(env, "android/content/Context");
	jmethodID start = find_method(env, contextClass, "startActivity", "(Landroid/content/Intent;)V");
	(*env)->CallVoidMethod(env, ctx, start, intent);
}

void sendNotification(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx, char *title, char *body) {
	JNIEnv *env = (JNIEnv*)jni_env;
	jstring titleStr = (*env)->NewStringUTF(env, title);
	jstring bodyStr = (*env)->NewStringUTF(env, body);

	jclass cls = find_class(env, "android/app/Notification$Builder");
	jmethodID constructor = find_method(env, cls, "<init>", "(Landroid/content/Context;)V");
	jobject builder = (*env)->NewObject(env, cls, constructor, ctx);

	jclass mgrCls = find_class(env, "android/app/NotificationManager");
	jobject mgr = getSystemService(env, ctx, "notification");

	if (isOreoOrLater(env)) {
		jstring channelId = (*env)->NewStringUTF(env, "bhojpur-notif");
		jstring name = (*env)->NewStringUTF(env, "Bhojpur Notification");
        int importance = 4; // IMPORTANCE_HIGH

		jclass chanCls = find_class(env, "android/app/NotificationChannel");
		jmethodID constructor = find_method(env, chanCls, "<init>", "(Ljava/lang/String;Ljava/lang/CharSequence;I)V");
		jobject channel = (*env)->NewObject(env, chanCls, constructor, channelId, name, importance);

		jmethodID createChannel = find_method(env, mgrCls, "createNotificationChannel", "(Landroid/app/NotificationChannel;)V");
		(*env)->CallVoidMethod(env, mgr, createChannel, channel);

		jmethodID setChannelId = find_method(env, cls, "setChannelId", "(Ljava/lang/String;)Landroid/app/Notification$Builder;");
		(*env)->CallObjectMethod(env, builder, setChannelId, channelId);
	}

	jmethodID setContentTitle = find_method(env, cls, "setContentTitle", "(Ljava/lang/CharSequence;)Landroid/app/Notification$Builder;");
	(*env)->CallObjectMethod(env, builder, setContentTitle, titleStr);

	jmethodID setContentText = find_method(env, cls, "setContentText", "(Ljava/lang/CharSequence;)Landroid/app/Notification$Builder;");
	(*env)->CallObjectMethod(env, builder, setContentText, bodyStr);

	int iconID = 17629184; // constant of "unknown app icon"
	jmethodID setSmallIcon = find_method(env, cls, "setSmallIcon", "(I)Landroid/app/Notification$Builder;");
	(*env)->CallObjectMethod(env, builder, setSmallIcon, iconID);

	jmethodID build = find_method(env, cls, "build", "()Landroid/app/Notification;");
	jobject notif = (*env)->CallObjectMethod(env, builder, build);

	jmethodID notify = find_method(env, mgrCls, "notify", "(ILandroid/app/Notification;)V");
	(*env)->CallVoidMethod(env, mgr, notify, nextId, notif);
	nextId++;
}