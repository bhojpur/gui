package tutorials

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
	"image/color"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/container"
	"github.com/bhojpur/gui/pkg/engine/layout"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

type iconInfo struct {
	name string
	icon gui.Resource
}

type browser struct {
	current int
	icons   []iconInfo

	name *widget.Select
	icon *widget.Icon
}

func (b *browser) setIcon(index int) {
	if index < 0 || index > len(b.icons)-1 {
		return
	}
	b.current = index

	b.name.SetSelected(b.icons[index].name)
	b.icon.SetResource(b.icons[index].icon)
}

// iconScreen loads a panel that shows the various icons available in Bhojpur GUI
func iconScreen(_ gui.Window) gui.CanvasObject {
	b := &browser{}
	b.icons = loadIcons()

	prev := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		b.setIcon(b.current - 1)
	})
	next := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		b.setIcon(b.current + 1)
	})
	b.name = widget.NewSelect(iconList(b.icons), func(name string) {
		for i, icon := range b.icons {
			if icon.name == name {
				if b.current != i {
					b.setIcon(i)
				}
				break
			}
		}
	})
	b.name.SetSelected(b.icons[b.current].name)
	buttons := container.NewHBox(prev, next)
	bar := container.NewBorder(nil, nil, buttons, nil, b.name)

	background := canvas.NewRasterWithPixels(checkerPattern)
	background.SetMinSize(gui.NewSize(280, 280))
	b.icon = widget.NewIcon(b.icons[b.current].icon)

	return gui.NewContainerWithLayout(layout.NewBorderLayout(
		bar, nil, nil, nil), bar, background, b.icon)
}

func checkerPattern(x, y, _, _ int) color.Color {
	x /= 20
	y /= 20

	if x%2 == y%2 {
		return theme.BackgroundColor()
	}

	return theme.ShadowColor()
}

func iconList(icons []iconInfo) []string {
	ret := make([]string, len(icons))
	for i, icon := range icons {
		ret[i] = icon.name
	}

	return ret
}

func loadIcons() []iconInfo {
	return []iconInfo{
		{"CancelIcon", theme.CancelIcon()},
		{"ConfirmIcon", theme.ConfirmIcon()},
		{"DeleteIcon", theme.DeleteIcon()},
		{"SearchIcon", theme.SearchIcon()},
		{"SearchReplaceIcon", theme.SearchReplaceIcon()},

		{"CheckButtonIcon", theme.CheckButtonIcon()},
		{"CheckButtonCheckedIcon", theme.CheckButtonCheckedIcon()},
		{"RadioButtonIcon", theme.RadioButtonIcon()},
		{"RadioButtonCheckedIcon", theme.RadioButtonCheckedIcon()},

		{"ColorAchromaticIcon", theme.ColorAchromaticIcon()},
		{"ColorChromaticIcon", theme.ColorChromaticIcon()},
		{"ColorPaletteIcon", theme.ColorPaletteIcon()},

		{"ContentAddIcon", theme.ContentAddIcon()},
		{"ContentRemoveIcon", theme.ContentRemoveIcon()},
		{"ContentClearIcon", theme.ContentClearIcon()},
		{"ContentCutIcon", theme.ContentCutIcon()},
		{"ContentCopyIcon", theme.ContentCopyIcon()},
		{"ContentPasteIcon", theme.ContentPasteIcon()},
		{"ContentRedoIcon", theme.ContentRedoIcon()},
		{"ContentUndoIcon", theme.ContentUndoIcon()},

		{"InfoIcon", theme.InfoIcon()},
		{"ErrorIcon", theme.ErrorIcon()},
		{"QuestionIcon", theme.QuestionIcon()},
		{"WarningIcon", theme.WarningIcon()},

		{"DocumentIcon", theme.DocumentIcon()},
		{"DocumentCreateIcon", theme.DocumentCreateIcon()},
		{"DocumentPrintIcon", theme.DocumentPrintIcon()},
		{"DocumentSaveIcon", theme.DocumentSaveIcon()},

		{"FileIcon", theme.FileIcon()},
		{"FileApplicationIcon", theme.FileApplicationIcon()},
		{"FileAudioIcon", theme.FileAudioIcon()},
		{"FileImageIcon", theme.FileImageIcon()},
		{"FileTextIcon", theme.FileTextIcon()},
		{"FileVideoIcon", theme.FileVideoIcon()},
		{"FolderIcon", theme.FolderIcon()},
		{"FolderNewIcon", theme.FolderNewIcon()},
		{"FolderOpenIcon", theme.FolderOpenIcon()},
		{"ComputerIcon", theme.ComputerIcon()},
		{"HomeIcon", theme.HomeIcon()},
		{"HelpIcon", theme.HelpIcon()},
		{"HistoryIcon", theme.HistoryIcon()},
		{"SettingsIcon", theme.SettingsIcon()},
		{"StorageIcon", theme.StorageIcon()},
		{"DownloadIcon", theme.DownloadIcon()},
		{"UploadIcon", theme.UploadIcon()},

		{"ViewFullScreenIcon", theme.ViewFullScreenIcon()},
		{"ViewRestoreIcon", theme.ViewRestoreIcon()},
		{"ViewRefreshIcon", theme.ViewRefreshIcon()},
		{"VisibilityIcon", theme.VisibilityIcon()},
		{"VisibilityOffIcon", theme.VisibilityOffIcon()},
		{"ZoomFitIcon", theme.ZoomFitIcon()},
		{"ZoomInIcon", theme.ZoomInIcon()},
		{"ZoomOutIcon", theme.ZoomOutIcon()},

		{"MoreHorizontalIcon", theme.MoreHorizontalIcon()},
		{"MoreVerticalIcon", theme.MoreVerticalIcon()},

		{"MoveDownIcon", theme.MoveDownIcon()},
		{"MoveUpIcon", theme.MoveUpIcon()},

		{"NavigateBackIcon", theme.NavigateBackIcon()},
		{"NavigateNextIcon", theme.NavigateNextIcon()},

		{"Menu", theme.MenuIcon()},
		{"MenuExpand", theme.MenuExpandIcon()},
		{"MenuDropDown", theme.MenuDropDownIcon()},
		{"MenuDropUp", theme.MenuDropUpIcon()},

		{"MailAttachmentIcon", theme.MailAttachmentIcon()},
		{"MailComposeIcon", theme.MailComposeIcon()},
		{"MailForwardIcon", theme.MailForwardIcon()},
		{"MailReplyIcon", theme.MailReplyIcon()},
		{"MailReplyAllIcon", theme.MailReplyAllIcon()},
		{"MailSendIcon", theme.MailSendIcon()},

		{"MediaFastForward", theme.MediaFastForwardIcon()},
		{"MediaFastRewind", theme.MediaFastRewindIcon()},
		{"MediaPause", theme.MediaPauseIcon()},
		{"MediaPlay", theme.MediaPlayIcon()},
		{"MediaStop", theme.MediaStopIcon()},
		{"MediaRecord", theme.MediaRecordIcon()},
		{"MediaReplay", theme.MediaReplayIcon()},
		{"MediaSkipNext", theme.MediaSkipNextIcon()},
		{"MediaSkipPrevious", theme.MediaSkipPreviousIcon()},

		{"VolumeDown", theme.VolumeDownIcon()},
		{"VolumeMute", theme.VolumeMuteIcon()},
		{"VolumeUp", theme.VolumeUpIcon()},

		{"AccountIcon", theme.AccountIcon()},
		{"LoginIcon", theme.LoginIcon()},
		{"LogoutIcon", theme.LogoutIcon()},

		{"ListIcon", theme.ListIcon()},
		{"GridIcon", theme.GridIcon()},
	}
}
