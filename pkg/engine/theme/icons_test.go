package theme

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
	"bytes"
	"fmt"
	"image"
	"image/color"
	"path/filepath"
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/internal/test"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gui.SetCurrentApp(&themedApp{})
}

func helperNewStaticResource() *gui.StaticResource {
	return &gui.StaticResource{
		StaticName: "content-remove.svg",
		StaticContent: []byte{
			60, 115, 118, 103, 32, 120, 109, 108, 110, 115, 61, 34, 104, 116, 116, 112, 58, 47, 47, 119, 119, 119, 46, 119, 51, 46, 111, 114, 103, 47, 50, 48, 48, 48, 47, 115, 118, 103, 34, 32, 119, 105, 100, 116, 104, 61, 34, 50, 52, 34, 32, 104, 101, 105, 103, 104, 116, 61, 34, 50, 52, 34, 32, 118, 105, 101, 119, 66, 111, 120, 61, 34, 48, 32, 48, 32, 50, 52, 32, 50, 52, 34, 62, 60, 112, 97, 116, 104, 32, 102, 105, 108, 108, 61, 34, 35, 102, 102, 102, 102, 102, 102, 34, 32, 100, 61, 34, 77, 49, 57, 32, 49, 51, 72, 53, 118, 45, 50, 104, 49, 52, 118, 50, 122, 34, 47, 62, 60, 112, 97, 116, 104, 32, 100, 61, 34, 77, 48, 32, 48, 104, 50, 52, 118, 50, 52, 72, 48, 122, 34, 32, 102, 105, 108, 108, 61, 34, 110, 111, 110, 101, 34, 47, 62, 60, 47, 115, 118, 103, 62},
	}
}

func helperLoadRes(t *testing.T, name string) gui.Resource {
	path := filepath.Join("testdata", name) // pathObj relative to this file
	res, err := gui.LoadResourceFromPath(path)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func helperDrawSVG(t *testing.T, data []byte) image.Image {
	icon, err := oksvg.ReadIconStream(bytes.NewReader(data))
	require.NoError(t, err, "failed to read SVG data")

	width := int(icon.ViewBox.W) * 2
	height := int(icon.ViewBox.H) * 2
	icon.SetTarget(0, 0, float64(width), float64(height))
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(width, height, img, img.Bounds())
	raster := rasterx.NewDasher(width, height, scanner)
	icon.Draw(raster, 1)
	return img
}

func TestIconThemeChangeDoesNotChangeName(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	cancel := CancelIcon()
	name := cancel.Name()

	gui.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, name, cancel.Name())
}

func TestIconThemeChangeContent(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	cancel := CancelIcon()
	content := cancel.Content()

	gui.CurrentApp().Settings().SetTheme(LightTheme())
	assert.NotEqual(t, content, cancel.Content())
}

func TestNewThemedResource_StaticResourceSupport(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	custom := NewThemedResource(helperNewStaticResource())
	content := custom.Content()
	name := custom.Name()

	gui.CurrentApp().Settings().SetTheme(LightTheme())
	assert.NotEqual(t, content, custom.Content())
	assert.Equal(t, name, custom.Name())
}

func TestNewDisabledResource(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	source := helperNewStaticResource()
	custom := NewDisabledResource(source)
	name := custom.Name()

	assert.Equal(t, name, fmt.Sprintf("disabled_%v", source.Name()))
}

func TestThemedResource_Invert(t *testing.T) {
	staticResource := helperLoadRes(t, "cancel_Paths.svg")
	inverted := NewInvertedThemedResource(staticResource)
	assert.Equal(t, "inverted-"+staticResource.Name(), inverted.Name())
}

func TestThemedResource_Name(t *testing.T) {
	staticResource := helperLoadRes(t, "cancel_Paths.svg")
	themedResource := &ThemedResource{
		source: staticResource,
	}
	assert.Equal(t, staticResource.Name(), themedResource.Name())
}

func TestThemedResource_Content_NoGroupsFile(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	staticResource := helperLoadRes(t, "cancel_Paths.svg")
	themedResource := &ThemedResource{
		source: staticResource,
	}
	assert.NotEqual(t, staticResource.Content(), themedResource.Content())
}

func TestThemedResource_Content_GroupPathFile(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	staticResource := helperLoadRes(t, "check_GroupPaths.svg")
	themedResource := &ThemedResource{
		source: staticResource,
	}
	assert.NotEqual(t, staticResource.Content(), themedResource.Content())
}

func TestThemedResource_Content_GroupRectFile(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	staticResource := helperLoadRes(t, "info_GroupRects.svg")
	themedResource := &ThemedResource{
		source: staticResource,
	}
	assert.NotEqual(t, staticResource.Content(), themedResource.Content())
}

func TestThemedResource_Content_GroupPolygonsFile(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	staticResource := helperLoadRes(t, "warning_GroupPolygons.svg")
	themedResource := &ThemedResource{
		source: staticResource,
	}
	assert.NotEqual(t, staticResource.Content(), themedResource.Content())
}

// a black svg object omits the fill tag, this checks it it still properly updated
func TestThemedResource_Content_BlackFillIsUpdated(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	staticResource := helperLoadRes(t, "cancel_PathsBlackFill.svg")
	themedResource := &ThemedResource{
		source: staticResource,
	}
	assert.NotEqual(t, staticResource.Content(), themedResource.Content())
}

func TestDisabledResource_Name(t *testing.T) {
	staticResource := helperLoadRes(t, "cancel_Paths.svg")
	disabledResource := &DisabledResource{
		source: staticResource,
	}
	assert.Equal(t, fmt.Sprintf("disabled_%v", staticResource.Name()), disabledResource.Name())
}

func TestDisabledResource_Content_NoGroupsFile(t *testing.T) {
	gui.CurrentApp().Settings().SetTheme(DarkTheme())
	staticResource := helperLoadRes(t, "cancel_Paths.svg")
	disabledResource := &DisabledResource{
		source: staticResource,
	}
	assert.NotEqual(t, staticResource.Content(), disabledResource.Content())
}

func TestColorizeResource(t *testing.T) {
	tests := map[string]struct {
		svgFile   string
		color     color.Color
		wantImage string
	}{
		"paths": {
			svgFile:   "cancel_Paths.svg",
			color:     color.NRGBA{R: 100, G: 100, A: 200},
			wantImage: "colorized/paths.png",
		},
		"circles": {
			svgFile:   "circles.svg",
			color:     color.NRGBA{R: 100, B: 100, A: 200},
			wantImage: "colorized/circles.png",
		},
		"polygons": {
			svgFile:   "polygons.svg",
			color:     color.NRGBA{G: 100, B: 100, A: 200},
			wantImage: "colorized/polygons.png",
		},
		"rects": {
			svgFile:   "rects.svg",
			color:     color.NRGBA{R: 100, G: 100, B: 100, A: 200},
			wantImage: "colorized/rects.png",
		},
		"group of paths": {
			svgFile:   "check_GroupPaths.svg",
			color:     color.NRGBA{R: 100, G: 100, A: 100},
			wantImage: "colorized/group_paths.png",
		},
		"group of circles": {
			svgFile:   "group_circles.svg",
			color:     color.NRGBA{R: 100, B: 100, A: 100},
			wantImage: "colorized/group_circles.png",
		},
		"group of polygons": {
			svgFile:   "warning_GroupPolygons.svg",
			color:     color.NRGBA{G: 100, B: 100, A: 100},
			wantImage: "colorized/group_polygons.png",
		},
		"group of rects": {
			svgFile:   "info_GroupRects.svg",
			color:     color.NRGBA{R: 100, G: 100, B: 100, A: 100},
			wantImage: "colorized/group_rects.png",
		},
		"NRGBA64": {
			svgFile: "circles.svg",
			// If the low 8 bits of each component were used, this would look cyan instead of yellow.
			// When the MSB is used instead, it correctly looks yellow.
			color:     color.NRGBA64{R: 0xff00, G: 0xffff, B: 0x00ff, A: 0xffff},
			wantImage: "colorized/circles_yellow.png",
		},
		"translucent NRGBA64": {
			svgFile:   "circles.svg",
			color:     color.NRGBA64{R: 0xff00, G: 0xffff, B: 0x00ff, A: 0x7fff},
			wantImage: "colorized/circles_yellow_translucent.png",
		},
		"RGBA": {
			svgFile:   "circles.svg",
			color:     color.RGBAModel.Convert(color.NRGBA{R: 0xff, G: 0xff, B: 0x00, A: 0xff}),
			wantImage: "colorized/circles_yellow.png",
		},
		"transluscent RGBA": {
			svgFile:   "circles.svg",
			color:     color.RGBAModel.Convert(color.NRGBA{R: 0xff, G: 0xff, B: 0x00, A: 0x7f}),
			wantImage: "colorized/circles_yellow_translucent.png",
		},
		"RGBA64": {
			svgFile: "circles.svg",
			// If the least significant byte of each component was being used, this would look cyan instead of yellow.
			// Since alpha=0xffff, unmultiplyAlpha knows it does not need to unmultiply anything, and so it just
			// returns the MSB of each component.
			color:     color.RGBA64Model.Convert(color.NRGBA64{R: 0xff00, G: 0xffff, B: 0x00ff, A: 0xffff}),
			wantImage: "colorized/circles_yellow.png",
		},
		"transluscent RGBA64": {
			svgFile: "circles.svg",
			// Since alpha!=0xffff, if we were to use R:0xff00, G:0xffff, B:0x00ff like before,
			// this would end up being drawn with 0xfeff00 instead of 0xffff00, and we would need a separate image to test for that.
			// Instead, we use R:0xfff0, G:0xfff0, B:0x000f, A:0x7fff instead, which unmultiplyAlpha returns as 0xff, 0xff, 0x00, 0x7f,
			// so that we correctly get 0xffff00 with alpha 0x7f when ToRGBA is used.
			// The RGBA64's contents are 0x7ff7, 0x7ff7, 0x0007, 0x7fff, so:
			// If ToRGBA wasn't being called and instead the LSB of each component was being read, this would show up as 0xf7f707 with alpha 0xff.
			// If the MSB was being read without umultiplication, this would show up as 0x7f7f00 with alpha 0x7f.
			color:     color.RGBA64Model.Convert(color.NRGBA64{R: 0xfff0, G: 0xfff0, B: 0x000f, A: 0x7fff}),
			wantImage: "colorized/circles_yellow_translucent.png",
		},
		"Alpha": {
			svgFile:   "circles.svg",
			color:     color.Alpha{A: 0x7f},
			wantImage: "colorized/circles_white_translucent.png",
		},
		"Alpha16": {
			svgFile: "circles.svg",
			// If the LSB from components returned by RGBA() was being used, this would be black.
			// If the MSB from components returned by RGBA() was being used, this would be grey.
			// It is white when either we bypass RGBA() and directly make a 0xffffff color with the alpha's MSB (which is what ToRGBA does),
			// or if we call RBGA(), un-multiply the alpha from the non-alpha components, and use their MSB to get white (Or something very near it like 0xfefefe).
			color:     color.Alpha16{A: 0x7f00},
			wantImage: "colorized/circles_white_translucent.png",
		},
		"Gray": {
			svgFile:   "circles.svg",
			color:     color.Gray{Y: 0xff},
			wantImage: "colorized/circles_white.png",
		},
		"Gray16": {
			svgFile: "circles.svg",
			// If the LSB from components returned by RGBA() was being used, this would be black.
			// It is white when either we bypass RGBA() and directly make a 0xffffff color with the alpha's MSB (which is what ToRGBA does),
			// or if we call RBGA(), un-multiply the alpha from the non-alpha components, and use their MSB to get white (Or something very near it like 0xfefefe),
			// or if the MSB from components returned by RGBA() was being used (because Gray and Gray16 do not have alpha values).
			color:     color.Gray16{Y: 0xff00},
			wantImage: "colorized/circles_white.png",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			src := helperLoadRes(t, tt.svgFile)
			got := helperDrawSVG(t, colorizeResource(src, tt.color))
			test.AssertImageMatches(t, tt.wantImage, got)
		})
	}
}

// Test Asset Sources
func Test_BhojpurLogo_FileSource(t *testing.T) {
	result := BhojpurLogo().Name()
	assert.Equal(t, "bhojpur.png", result)
}

func Test_CancelIcon_FileSource(t *testing.T) {
	result := CancelIcon().Name()
	assert.Equal(t, "cancel.svg", result)
}

func Test_ConfirmIcon_FileSource(t *testing.T) {
	result := ConfirmIcon().Name()
	assert.Equal(t, "check.svg", result)
}

func Test_DeleteIcon_FileSource(t *testing.T) {
	result := DeleteIcon().Name()
	assert.Equal(t, "delete.svg", result)
}

func Test_SearchIcon_FileSource(t *testing.T) {
	result := SearchIcon().Name()
	assert.Equal(t, "search.svg", result)
}

func Test_SearchReplaceIcon_FileSource(t *testing.T) {
	result := SearchReplaceIcon().Name()
	assert.Equal(t, "search-replace.svg", result)
}

func Test_CheckButtonIcon_FileSource(t *testing.T) {
	result := CheckButtonIcon().Name()
	assert.Equal(t, "check-box-blank.svg", result)
}

func Test_CheckButtonCheckedIcon_FileSource(t *testing.T) {
	result := CheckButtonCheckedIcon().Name()
	assert.Equal(t, "check-box.svg", result)
}

func Test_RadioButtonIcon_FileSource(t *testing.T) {
	result := RadioButtonIcon().Name()
	assert.Equal(t, "radio-button.svg", result)
}

func Test_RadioButtonCheckedIcon_FileSource(t *testing.T) {
	result := RadioButtonCheckedIcon().Name()
	assert.Equal(t, "radio-button-checked.svg", result)
}

func Test_ContentAddIcon_FileSource(t *testing.T) {
	result := ContentAddIcon().Name()
	assert.Equal(t, "content-add.svg", result)
}

func Test_ContentRemoveIcon_FileSource(t *testing.T) {
	result := ContentRemoveIcon().Name()
	assert.Equal(t, "content-remove.svg", result)
}

func Test_ContentClearIcon_FileSource(t *testing.T) {
	result := ContentClearIcon().Name()
	assert.Equal(t, "cancel.svg", result)
}

func Test_ContentCutIcon_FileSource(t *testing.T) {
	result := ContentCutIcon().Name()
	assert.Equal(t, "content-cut.svg", result)
}

func Test_ContentCopyIcon_FileSource(t *testing.T) {
	result := ContentCopyIcon().Name()
	assert.Equal(t, "content-copy.svg", result)
}

func Test_ContentPasteIcon_FileSource(t *testing.T) {
	result := ContentPasteIcon().Name()
	assert.Equal(t, "content-paste.svg", result)
}

func Test_ContentRedoIcon_FileSource(t *testing.T) {
	result := ContentRedoIcon().Name()
	assert.Equal(t, "content-redo.svg", result)
}

func Test_ContentUndoIcon_FileSource(t *testing.T) {
	result := ContentUndoIcon().Name()
	assert.Equal(t, "content-undo.svg", result)
}

func Test_DocumentCreateIcon_FileSource(t *testing.T) {
	result := DocumentCreateIcon().Name()
	assert.Equal(t, "document-create.svg", result)
}

func Test_DocumentPrintIcon_FileSource(t *testing.T) {
	result := DocumentPrintIcon().Name()
	assert.Equal(t, "document-print.svg", result)
}

func Test_DocumentSaveIcon_FileSource(t *testing.T) {
	result := DocumentSaveIcon().Name()
	assert.Equal(t, "document-save.svg", result)
}

func Test_InfoIcon_FileSource(t *testing.T) {
	result := InfoIcon().Name()
	assert.Equal(t, "info.svg", result)
}

func Test_QuestionIcon_FileSource(t *testing.T) {
	result := QuestionIcon().Name()
	assert.Equal(t, "question.svg", result)
}

func Test_WarningIcon_FileSource(t *testing.T) {
	result := WarningIcon().Name()
	assert.Equal(t, "warning.svg", result)
}

func Test_FolderIcon_FileSource(t *testing.T) {
	result := FolderIcon().Name()
	assert.Equal(t, "folder.svg", result)
}

func Test_FolderNewIcon_FileSource(t *testing.T) {
	result := FolderNewIcon().Name()
	assert.Equal(t, "folder-new.svg", result)
}

func Test_FolderOpenIcon_FileSource(t *testing.T) {
	result := FolderOpenIcon().Name()
	assert.Equal(t, "folder-open.svg", result)
}

func Test_HelpIcon_FileSource(t *testing.T) {
	result := HelpIcon().Name()
	assert.Equal(t, "help.svg", result)
}

func Test_HomeIcon_FileSource(t *testing.T) {
	result := HomeIcon().Name()
	assert.Equal(t, "home.svg", result)
}

func Test_SettingsIcon_FileSource(t *testing.T) {
	result := SettingsIcon().Name()
	assert.Equal(t, "settings.svg", result)
}

func Test_MailAttachmentIcon_FileSource(t *testing.T) {
	result := MailAttachmentIcon().Name()
	assert.Equal(t, "mail-attachment.svg", result)
}

func Test_MailComposeIcon_FileSource(t *testing.T) {
	result := MailComposeIcon().Name()
	assert.Equal(t, "mail-compose.svg", result)
}

func Test_MailForwardIcon_FileSource(t *testing.T) {
	result := MailForwardIcon().Name()
	assert.Equal(t, "mail-forward.svg", result)
}

func Test_MailReplyIcon_FileSource(t *testing.T) {
	result := MailReplyIcon().Name()
	assert.Equal(t, "mail-reply.svg", result)
}

func Test_MailReplyAllIcon_FileSource(t *testing.T) {
	result := MailReplyAllIcon().Name()
	assert.Equal(t, "mail-reply_all.svg", result)
}

func Test_MailSendIcon_FileSource(t *testing.T) {
	result := MailSendIcon().Name()
	assert.Equal(t, "mail-send.svg", result)
}

func Test_MoveDownIcon_FileSource(t *testing.T) {
	result := MoveDownIcon().Name()
	assert.Equal(t, "arrow-down.svg", result)
}

func Test_MoveUpIcon_FileSource(t *testing.T) {
	result := MoveUpIcon().Name()
	assert.Equal(t, "arrow-up.svg", result)
}

func Test_NavigateBackIcon_FileSource(t *testing.T) {
	result := NavigateBackIcon().Name()
	assert.Equal(t, "arrow-back.svg", result)
}

func Test_NavigateNextIcon_FileSource(t *testing.T) {
	result := NavigateNextIcon().Name()
	assert.Equal(t, "arrow-forward.svg", result)
}

func Test_ViewFullScreenIcon_FileSource(t *testing.T) {
	result := ViewFullScreenIcon().Name()
	assert.Equal(t, "view-fullscreen.svg", result)
}

func Test_ViewRestoreIcon_FileSource(t *testing.T) {
	result := ViewRestoreIcon().Name()
	assert.Equal(t, "view-zoom-fit.svg", result)
}

func Test_ViewRefreshIcon_FileSource(t *testing.T) {
	result := ViewRefreshIcon().Name()
	assert.Equal(t, "view-refresh.svg", result)
}

func Test_ZoomFitIcon_FileSource(t *testing.T) {
	result := ZoomFitIcon().Name()
	assert.Equal(t, "view-zoom-fit.svg", result)
}

func Test_ZoomInIcon_FileSource(t *testing.T) {
	result := ZoomInIcon().Name()
	assert.Equal(t, "view-zoom-in.svg", result)
}

func Test_ZoomOutIcon_FileSource(t *testing.T) {
	result := ZoomOutIcon().Name()
	assert.Equal(t, "view-zoom-out.svg", result)
}

func Test_VisibilityIcon_FileSource(t *testing.T) {
	result := VisibilityIcon().Name()
	assert.Equal(t, "visibility.svg", result)
}

func Test_VisibilityOffIcon_FileSource(t *testing.T) {
	result := VisibilityOffIcon().Name()
	assert.Equal(t, "visibility-off.svg", result)
}

func Test_AccountIcon_FileSource(t *testing.T) {
	result := AccountIcon().Name()
	assert.Equal(t, "account.svg", result)
}

func Test_LoginIcon_FileSource(t *testing.T) {
	result := LoginIcon().Name()
	assert.Equal(t, "login.svg", result)
}

func Test_LogoutIcon_FileSource(t *testing.T) {
	result := LogoutIcon().Name()
	assert.Equal(t, "logout.svg", result)
}

func Test_ListIcon_FileSource(t *testing.T) {
	result := ListIcon().Name()
	assert.Equal(t, "list.svg", result)
}

func Test_GridIcon_FileSource(t *testing.T) {
	result := GridIcon().Name()
	assert.Equal(t, "grid.svg", result)
}
