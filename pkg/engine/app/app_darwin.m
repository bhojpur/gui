// +build !ci

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

#import <Foundation/Foundation.h>
#if __MAC_OS_X_VERSION_MAX_ALLOWED >= 101400
#import <UserNotifications/UserNotifications.h>
#endif

static int notifyNum = 0;

extern void fallbackSend(char *cTitle, char *cBody);

bool isBundled() {
    return [[NSBundle mainBundle] bundleIdentifier] != nil;
}

#if __MAC_OS_X_VERSION_MAX_ALLOWED >= 101400
void doSendNotification(UNUserNotificationCenter *center, NSString *title, NSString *body) {
    UNMutableNotificationContent *content = [UNMutableNotificationContent new];
    [content autorelease];
    content.title = title;
    content.body = body;

    notifyNum++;
    NSString *identifier = [NSString stringWithFormat:@"bhojpur-notify-%d", notifyNum];
    UNNotificationRequest *request = [UNNotificationRequest requestWithIdentifier:identifier
        content:content trigger:nil];

    [center addNotificationRequest:request withCompletionHandler:^(NSError * _Nullable error) {
        if (error != nil) {
            NSLog(@"Could not send notification: %@", error);
        }
    }];
}

void sendNotification(char *cTitle, char *cBody) {
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    NSString *title = [NSString stringWithUTF8String:cTitle];
    NSString *body = [NSString stringWithUTF8String:cBody];

    UNAuthorizationOptions options = UNAuthorizationOptionAlert;
    [center requestAuthorizationWithOptions:options
        completionHandler:^(BOOL granted, NSError *_Nullable error) {
            if (!granted) {
                if (error != NULL) {
                    NSLog(@"Error asking for permission to send notifications %@", error);
                    // this happens if our app was not signed, so do it the old way
                    fallbackSend((char *)[title UTF8String], (char *)[body UTF8String]);
                } else {
                    NSLog(@"Unable to get permission to send notifications");
                }
            } else {
                doSendNotification(center, title, body);
            }
        }];
}
#else
void sendNotification(char *cTitle, char *cBody) {
	fallbackSend(cTitle, cBody);
}
#endif
