//go:build !ci && !js && !android && !ios && !wasm && !test_web_driver
// +build !ci,!js,!android,!ios,!wasm,!test_web_driver

package app

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
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/sys/windows/registry"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/theme"

	"golang.org/x/sys/execabs"
)

const notificationTemplate = `$title = "%s"
$content = "%s"
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] > $null
$template = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastText02)
$toastXml = [xml] $template.GetXml()
$toastXml.GetElementsByTagName("text")[0].AppendChild($toastXml.CreateTextNode($title)) > $null
$toastXml.GetElementsByTagName("text")[1].AppendChild($toastXml.CreateTextNode($content)) > $null
$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$xml.LoadXml($toastXml.OuterXml)
$toast = [Windows.UI.Notifications.ToastNotification]::new($xml)
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("%s").Show($toast);`

func isDark() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize`, registry.QUERY_VALUE)
	if err != nil { // older version of Windows will not have this key
		return false
	}
	defer k.Close()

	useLight, _, err := k.GetIntegerValue("AppsUseLightTheme")
	if err != nil { // older version of Windows will not have this value
		return false
	}

	return useLight == 0
}

func defaultVariant() gui.ThemeVariant {
	if isDark() {
		return theme.VariantDark
	}
	return theme.VariantLight
}

func rootConfigDir() string {
	homeDir, _ := os.UserHomeDir()

	desktopConfig := filepath.Join(filepath.Join(homeDir, "AppData"), "Roaming")
	return filepath.Join(desktopConfig, "Bhojpur")
}

func (a *bhojpurApp) OpenURL(url *url.URL) error {
	cmd := a.exec("rundll32", "url.dll,FileProtocolHandler", url.String())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

var scriptNum = 0

func (a *bhojpurApp) SendNotification(n *gui.Notification) {
	title := escapeNotificationString(n.Title)
	content := escapeNotificationString(n.Content)
	appID := a.UniqueID() // TODO once we have an app name compiled in this could be improved
	if appID == "" || strings.Index(appID, "missing-id") == 0 {
		appID = "Bhojpur GUI"
	}

	script := fmt.Sprintf(notificationTemplate, title, content, appID)
	go runScript("notify", script)
}

func escapeNotificationString(in string) string {
	noSlash := strings.ReplaceAll(in, "`", "``")
	return strings.ReplaceAll(noSlash, "\"", "`\"")
}

func runScript(name, script string) {
	scriptNum++
	appID := gui.CurrentApp().UniqueID()
	fileName := fmt.Sprintf("bhojpur-%s-%s-%d.ps1", appID, name, scriptNum)

	tmpFilePath := filepath.Join(os.TempDir(), fileName)
	err := ioutil.WriteFile(tmpFilePath, []byte(script), 0600)
	if err != nil {
		gui.LogError("Could not write script to show notification", err)
		return
	}
	defer os.Remove(tmpFilePath)

	launch := "(Get-Content -Encoding UTF8 -Path " + tmpFilePath + " -Raw) | Invoke-Expression"
	cmd := execabs.Command("PowerShell", "-ExecutionPolicy", "Bypass", launch)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err = cmd.Run()
	if err != nil {
		gui.LogError("Failed to launch windows notify script", err)
	}
}
func watchTheme() {
	// TODO monitor the Windows theme
}
