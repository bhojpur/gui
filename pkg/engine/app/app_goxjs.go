//go:build !ci && (!android || !ios || !mobile) && (js || wasm || test_web_driver)
// +build !ci
// +build !android !ios !mobile
// +build js wasm test_web_driver

package app

import (
	"errors"
	"net/url"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

func defaultVariant() gui.ThemeVariant {
	return theme.VariantDark
}

func (app *fbhojpurApp) OpenURL(url *url.URL) error {
	// TODO #2736
	return errors.New("OpenURL is not supported yet with GopherJS backend.")
}

func (app *bhojpurApp) SendNotification(_ *gui.Notification) {
	// TODO #2735
	gui.LogError("Sending notification is not supported yet.", nil)
}

func rootConfigDir() string {
	return "/data/"
}

func defaultTheme() gui.Theme {
	return theme.DarkTheme()
}
