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
	"encoding/xml"
	"image/color"

	gui "github.com/bhojpur/gui/pkg/engine"
)

const (
	// IconNameCancel is the name of theme lookup for cancel icon.
	//
	// Since: 2.0
	IconNameCancel gui.ThemeIconName = "cancel"

	// IconNameConfirm is the name of theme lookup for confirm icon.
	//
	// Since: 2.0
	IconNameConfirm gui.ThemeIconName = "confirm"

	// IconNameDelete is the name of theme lookup for delete icon.
	//
	// Since: 2.0
	IconNameDelete gui.ThemeIconName = "delete"

	// IconNameSearch is the name of theme lookup for search icon.
	//
	// Since: 2.0
	IconNameSearch gui.ThemeIconName = "search"

	// IconNameSearchReplace is the name of theme lookup for search and replace icon.
	//
	// Since: 2.0
	IconNameSearchReplace gui.ThemeIconName = "searchReplace"

	// IconNameMenu is the name of theme lookup for menu icon.
	//
	// Since: 2.0
	IconNameMenu gui.ThemeIconName = "menu"

	// IconNameMenuExpand is the name of theme lookup for menu expansion icon.
	//
	// Since: 2.0
	IconNameMenuExpand gui.ThemeIconName = "menuExpand"

	// IconNameCheckButtonChecked is the name of theme lookup for checked check button icon.
	//
	// Since: 2.0
	IconNameCheckButtonChecked gui.ThemeIconName = "checked"

	// IconNameCheckButton is the name of theme lookup for  unchecked check button icon.
	//
	// Since: 2.0
	IconNameCheckButton gui.ThemeIconName = "unchecked"

	// IconNameRadioButton is the name of theme lookup for radio button unchecked icon.
	//
	// Since: 2.0
	IconNameRadioButton gui.ThemeIconName = "radioButton"

	// IconNameRadioButtonChecked is the name of theme lookup for radio button checked icon.
	//
	// Since: 2.0
	IconNameRadioButtonChecked gui.ThemeIconName = "radioButtonChecked"

	// IconNameColorAchromatic is the name of theme lookup for greyscale color icon.
	//
	// Since: 2.0
	IconNameColorAchromatic gui.ThemeIconName = "colorAchromatic"

	// IconNameColorChromatic is the name of theme lookup for full color icon.
	//
	// Since: 2.0
	IconNameColorChromatic gui.ThemeIconName = "colorChromatic"

	// IconNameColorPalette is the name of theme lookup for color palette icon.
	//
	// Since: 2.0
	IconNameColorPalette gui.ThemeIconName = "colorPalette"

	// IconNameContentAdd is the name of theme lookup for content add icon.
	//
	// Since: 2.0
	IconNameContentAdd gui.ThemeIconName = "contentAdd"

	// IconNameContentRemove is the name of theme lookup for content remove icon.
	//
	// Since: 2.0
	IconNameContentRemove gui.ThemeIconName = "contentRemove"

	// IconNameContentCut is the name of theme lookup for content cut icon.
	//
	// Since: 2.0
	IconNameContentCut gui.ThemeIconName = "contentCut"

	// IconNameContentCopy is the name of theme lookup for content copy icon.
	//
	// Since: 2.0
	IconNameContentCopy gui.ThemeIconName = "contentCopy"

	// IconNameContentPaste is the name of theme lookup for content paste icon.
	//
	// Since: 2.0
	IconNameContentPaste gui.ThemeIconName = "contentPaste"

	// IconNameContentClear is the name of theme lookup for content clear icon.
	//
	// Since: 2.0
	IconNameContentClear gui.ThemeIconName = "contentClear"

	// IconNameContentRedo is the name of theme lookup for content redo icon.
	//
	// Since: 2.0
	IconNameContentRedo gui.ThemeIconName = "contentRedo"

	// IconNameContentUndo is the name of theme lookup for content undo icon.
	//
	// Since: 2.0
	IconNameContentUndo gui.ThemeIconName = "contentUndo"

	// IconNameInfo is the name of theme lookup for info icon.
	//
	// Since: 2.0
	IconNameInfo gui.ThemeIconName = "info"

	// IconNameQuestion is the name of theme lookup for question icon.
	//
	// Since: 2.0
	IconNameQuestion gui.ThemeIconName = "question"

	// IconNameWarning is the name of theme lookup for warning icon.
	//
	// Since: 2.0
	IconNameWarning gui.ThemeIconName = "warning"

	// IconNameError is the name of theme lookup for error icon.
	//
	// Since: 2.0
	IconNameError gui.ThemeIconName = "error"

	// IconNameDocument is the name of theme lookup for document icon.
	//
	// Since: 2.0
	IconNameDocument gui.ThemeIconName = "document"

	// IconNameDocumentCreate is the name of theme lookup for document create icon.
	//
	// Since: 2.0
	IconNameDocumentCreate gui.ThemeIconName = "documentCreate"

	// IconNameDocumentPrint is the name of theme lookup for document print icon.
	//
	// Since: 2.0
	IconNameDocumentPrint gui.ThemeIconName = "documentPrint"

	// IconNameDocumentSave is the name of theme lookup for document save icon.
	//
	// Since: 2.0
	IconNameDocumentSave gui.ThemeIconName = "documentSave"

	// IconNameMoreHorizontal is the name of theme lookup for horizontal more.
	//
	// Since 2.0
	IconNameMoreHorizontal gui.ThemeIconName = "moreHorizontal"

	// IconNameMoreVertical is the name of theme lookup for vertical more.
	//
	// Since 2.0
	IconNameMoreVertical gui.ThemeIconName = "moreVertical"

	// IconNameMailAttachment is the name of theme lookup for mail attachment icon.
	//
	// Since: 2.0
	IconNameMailAttachment gui.ThemeIconName = "mailAttachment"

	// IconNameMailCompose is the name of theme lookup for mail compose icon.
	//
	// Since: 2.0
	IconNameMailCompose gui.ThemeIconName = "mailCompose"

	// IconNameMailForward is the name of theme lookup for mail forward icon.
	//
	// Since: 2.0
	IconNameMailForward gui.ThemeIconName = "mailForward"

	// IconNameMailReply is the name of theme lookup for mail reply icon.
	//
	// Since: 2.0
	IconNameMailReply gui.ThemeIconName = "mailReply"

	// IconNameMailReplyAll is the name of theme lookup for mail reply-all icon.
	//
	// Since: 2.0
	IconNameMailReplyAll gui.ThemeIconName = "mailReplyAll"

	// IconNameMailSend is the name of theme lookup for mail send icon.
	//
	// Since: 2.0
	IconNameMailSend gui.ThemeIconName = "mailSend"

	// IconNameMediaMusic is the name of theme lookup for media music icon.
	//
	// Since: 2.1
	IconNameMediaMusic gui.ThemeIconName = "mediaMusic"

	// IconNameMediaPhoto is the name of theme lookup for media photo icon.
	//
	// Since: 2.1
	IconNameMediaPhoto gui.ThemeIconName = "mediaPhoto"

	// IconNameMediaVideo is the name of theme lookup for media video icon.
	//
	// Since: 2.1
	IconNameMediaVideo gui.ThemeIconName = "mediaVideo"

	// IconNameMediaFastForward is the name of theme lookup for media fast-forward icon.
	//
	// Since: 2.0
	IconNameMediaFastForward gui.ThemeIconName = "mediaFastForward"

	// IconNameMediaFastRewind is the name of theme lookup for media fast-rewind icon.
	//
	// Since: 2.0
	IconNameMediaFastRewind gui.ThemeIconName = "mediaFastRewind"

	// IconNameMediaPause is the name of theme lookup for media pause icon.
	//
	// Since: 2.0
	IconNameMediaPause gui.ThemeIconName = "mediaPause"

	// IconNameMediaPlay is the name of theme lookup for media play icon.
	//
	// Since: 2.0
	IconNameMediaPlay gui.ThemeIconName = "mediaPlay"

	// IconNameMediaRecord is the name of theme lookup for media record icon.
	//
	// Since: 2.0
	IconNameMediaRecord gui.ThemeIconName = "mediaRecord"

	// IconNameMediaReplay is the name of theme lookup for media replay icon.
	//
	// Since: 2.0
	IconNameMediaReplay gui.ThemeIconName = "mediaReplay"

	// IconNameMediaSkipNext is the name of theme lookup for media skip next icon.
	//
	// Since: 2.0
	IconNameMediaSkipNext gui.ThemeIconName = "mediaSkipNext"

	// IconNameMediaSkipPrevious is the name of theme lookup for media skip previous icon.
	//
	// Since: 2.0
	IconNameMediaSkipPrevious gui.ThemeIconName = "mediaSkipPrevious"

	// IconNameMediaStop is the name of theme lookup for media stop icon.
	//
	// Since: 2.0
	IconNameMediaStop gui.ThemeIconName = "mediaStop"

	// IconNameMoveDown is the name of theme lookup for move down icon.
	//
	// Since: 2.0
	IconNameMoveDown gui.ThemeIconName = "arrowDown"

	// IconNameMoveUp is the name of theme lookup for move up icon.
	//
	// Since: 2.0
	IconNameMoveUp gui.ThemeIconName = "arrowUp"

	// IconNameNavigateBack is the name of theme lookup for navigate back icon.
	//
	// Since: 2.0
	IconNameNavigateBack gui.ThemeIconName = "arrowBack"

	// IconNameNavigateNext is the name of theme lookup for navigate next icon.
	//
	// Since: 2.0
	IconNameNavigateNext gui.ThemeIconName = "arrowForward"

	// IconNameArrowDropDown is the name of theme lookup for drop-down arrow icon.
	//
	// Since: 2.0
	IconNameArrowDropDown gui.ThemeIconName = "arrowDropDown"

	// IconNameArrowDropUp is the name of theme lookup for drop-up arrow icon.
	//
	// Since: 2.0
	IconNameArrowDropUp gui.ThemeIconName = "arrowDropUp"

	// IconNameFile is the name of theme lookup for file icon.
	//
	// Since: 2.0
	IconNameFile gui.ThemeIconName = "file"

	// IconNameFileApplication is the name of theme lookup for file application icon.
	//
	// Since: 2.0
	IconNameFileApplication gui.ThemeIconName = "fileApplication"

	// IconNameFileAudio is the name of theme lookup for file audio icon.
	//
	// Since: 2.0
	IconNameFileAudio gui.ThemeIconName = "fileAudio"

	// IconNameFileImage is the name of theme lookup for file image icon.
	//
	// Since: 2.0
	IconNameFileImage gui.ThemeIconName = "fileImage"

	// IconNameFileText is the name of theme lookup for file text icon.
	//
	// Since: 2.0
	IconNameFileText gui.ThemeIconName = "fileText"

	// IconNameFileVideo is the name of theme lookup for file video icon.
	//
	// Since: 2.0
	IconNameFileVideo gui.ThemeIconName = "fileVideo"

	// IconNameFolder is the name of theme lookup for folder icon.
	//
	// Since: 2.0
	IconNameFolder gui.ThemeIconName = "folder"

	// IconNameFolderNew is the name of theme lookup for folder new icon.
	//
	// Since: 2.0
	IconNameFolderNew gui.ThemeIconName = "folderNew"

	// IconNameFolderOpen is the name of theme lookup for folder open icon.
	//
	// Since: 2.0
	IconNameFolderOpen gui.ThemeIconName = "folderOpen"

	// IconNameHelp is the name of theme lookup for help icon.
	//
	// Since: 2.0
	IconNameHelp gui.ThemeIconName = "help"

	// IconNameHistory is the name of theme lookup for history icon.
	//
	// Since: 2.0
	IconNameHistory gui.ThemeIconName = "history"

	// IconNameHome is the name of theme lookup for home icon.
	//
	// Since: 2.0
	IconNameHome gui.ThemeIconName = "home"

	// IconNameSettings is the name of theme lookup for settings icon.
	//
	// Since: 2.0
	IconNameSettings gui.ThemeIconName = "settings"

	// IconNameStorage is the name of theme lookup for storage icon.
	//
	// Since: 2.0
	IconNameStorage gui.ThemeIconName = "storage"

	// IconNameUpload is the name of theme lookup for upload icon.
	//
	// Since: 2.0
	IconNameUpload gui.ThemeIconName = "upload"

	// IconNameViewFullScreen is the name of theme lookup for view fullscreen icon.
	//
	// Since: 2.0
	IconNameViewFullScreen gui.ThemeIconName = "viewFullScreen"

	// IconNameViewRefresh is the name of theme lookup for view refresh icon.
	//
	// Since: 2.0
	IconNameViewRefresh gui.ThemeIconName = "viewRefresh"

	// IconNameViewZoomFit is the name of theme lookup for view zoom fit icon.
	//
	// Since: 2.0
	IconNameViewZoomFit gui.ThemeIconName = "viewZoomFit"

	// IconNameViewZoomIn is the name of theme lookup for view zoom in icon.
	//
	// Since: 2.0
	IconNameViewZoomIn gui.ThemeIconName = "viewZoomIn"

	// IconNameViewZoomOut is the name of theme lookup for view zoom out icon.
	//
	// Since: 2.0
	IconNameViewZoomOut gui.ThemeIconName = "viewZoomOut"

	// IconNameViewRestore is the name of theme lookup for view restore icon.
	//
	// Since: 2.0
	IconNameViewRestore gui.ThemeIconName = "viewRestore"

	// IconNameVisibility is the name of theme lookup for visibility icon.
	//
	// Since: 2.0
	IconNameVisibility gui.ThemeIconName = "visibility"

	// IconNameVisibilityOff is the name of theme lookup for invisibility icon.
	//
	// Since: 2.0
	IconNameVisibilityOff gui.ThemeIconName = "visibilityOff"

	// IconNameVolumeDown is the name of theme lookup for volume down icon.
	//
	// Since: 2.0
	IconNameVolumeDown gui.ThemeIconName = "volumeDown"

	// IconNameVolumeMute is the name of theme lookup for volume mute icon.
	//
	// Since: 2.0
	IconNameVolumeMute gui.ThemeIconName = "volumeMute"

	// IconNameVolumeUp is the name of theme lookup for volume up icon.
	//
	// Since: 2.0
	IconNameVolumeUp gui.ThemeIconName = "volumeUp"

	// IconNameDownload is the name of theme lookup for download icon.
	//
	// Since: 2.0
	IconNameDownload gui.ThemeIconName = "download"

	// IconNameComputer is the name of theme lookup for computer icon.
	//
	// Since: 2.0
	IconNameComputer gui.ThemeIconName = "computer"

	// IconNameAccount is the name of theme lookup for account icon.
	//
	// Since: 2.1
	IconNameAccount gui.ThemeIconName = "account"

	// IconNameLogin is the name of theme lookup for login icon.
	//
	// Since: 2.1
	IconNameLogin gui.ThemeIconName = "login"

	// IconNameLogout is the name of theme lookup for logout icon.
	//
	// Since: 2.1
	IconNameLogout gui.ThemeIconName = "logout"

	// IconNameList is the name of theme lookup for list icon.
	//
	// Since: 2.1
	IconNameList gui.ThemeIconName = "list"

	// IconNameGrid is the name of theme lookup for grid icon.
	//
	// Since: 2.1
	IconNameGrid gui.ThemeIconName = "grid"
)

var (
	icons = map[gui.ThemeIconName]gui.Resource{
		IconNameCancel:        NewThemedResource(cancelIconRes),
		IconNameConfirm:       NewThemedResource(checkIconRes),
		IconNameDelete:        NewThemedResource(deleteIconRes),
		IconNameSearch:        NewThemedResource(searchIconRes),
		IconNameSearchReplace: NewThemedResource(searchreplaceIconRes),
		IconNameMenu:          NewThemedResource(menuIconRes),
		IconNameMenuExpand:    NewThemedResource(menuexpandIconRes),

		IconNameCheckButton:        NewThemedResource(checkboxblankIconRes),
		IconNameCheckButtonChecked: NewThemedResource(checkboxIconRes),
		IconNameRadioButton:        NewThemedResource(radiobuttonIconRes),
		IconNameRadioButtonChecked: NewThemedResource(radiobuttoncheckedIconRes),

		IconNameContentAdd:    NewThemedResource(contentaddIconRes),
		IconNameContentClear:  NewThemedResource(cancelIconRes),
		IconNameContentRemove: NewThemedResource(contentremoveIconRes),
		IconNameContentCut:    NewThemedResource(contentcutIconRes),
		IconNameContentCopy:   NewThemedResource(contentcopyIconRes),
		IconNameContentPaste:  NewThemedResource(contentpasteIconRes),
		IconNameContentRedo:   NewThemedResource(contentredoIconRes),
		IconNameContentUndo:   NewThemedResource(contentundoIconRes),

		IconNameColorAchromatic: NewThemedResource(colorachromaticIconRes),
		IconNameColorChromatic:  NewThemedResource(colorchromaticIconRes),
		IconNameColorPalette:    NewThemedResource(colorpaletteIconRes),

		IconNameDocument:       NewThemedResource(documentIconRes),
		IconNameDocumentCreate: NewThemedResource(documentcreateIconRes),
		IconNameDocumentPrint:  NewThemedResource(documentprintIconRes),
		IconNameDocumentSave:   NewThemedResource(documentsaveIconRes),

		IconNameMoreHorizontal: NewThemedResource(morehorizontalIconRes),
		IconNameMoreVertical:   NewThemedResource(moreverticalIconRes),

		IconNameInfo:     NewThemedResource(infoIconRes),
		IconNameQuestion: NewThemedResource(questionIconRes),
		IconNameWarning:  NewThemedResource(warningIconRes),
		IconNameError:    NewThemedResource(errorIconRes),

		IconNameMailAttachment: NewThemedResource(mailattachmentIconRes),
		IconNameMailCompose:    NewThemedResource(mailcomposeIconRes),
		IconNameMailForward:    NewThemedResource(mailforwardIconRes),
		IconNameMailReply:      NewThemedResource(mailreplyIconRes),
		IconNameMailReplyAll:   NewThemedResource(mailreplyallIconRes),
		IconNameMailSend:       NewThemedResource(mailsendIconRes),

		IconNameMediaMusic:        NewThemedResource(mediamusicIconRes),
		IconNameMediaPhoto:        NewThemedResource(mediaphotoIconRes),
		IconNameMediaVideo:        NewThemedResource(mediavideoIconRes),
		IconNameMediaFastForward:  NewThemedResource(mediafastforwardIconRes),
		IconNameMediaFastRewind:   NewThemedResource(mediafastrewindIconRes),
		IconNameMediaPause:        NewThemedResource(mediapauseIconRes),
		IconNameMediaPlay:         NewThemedResource(mediaplayIconRes),
		IconNameMediaRecord:       NewThemedResource(mediarecordIconRes),
		IconNameMediaReplay:       NewThemedResource(mediareplayIconRes),
		IconNameMediaSkipNext:     NewThemedResource(mediaskipnextIconRes),
		IconNameMediaSkipPrevious: NewThemedResource(mediaskippreviousIconRes),
		IconNameMediaStop:         NewThemedResource(mediastopIconRes),

		IconNameNavigateBack:  NewThemedResource(arrowbackIconRes),
		IconNameMoveDown:      NewThemedResource(arrowdownIconRes),
		IconNameNavigateNext:  NewThemedResource(arrowforwardIconRes),
		IconNameMoveUp:        NewThemedResource(arrowupIconRes),
		IconNameArrowDropDown: NewThemedResource(arrowdropdownIconRes),
		IconNameArrowDropUp:   NewThemedResource(arrowdropupIconRes),

		IconNameFile:            NewThemedResource(fileIconRes),
		IconNameFileApplication: NewThemedResource(fileapplicationIconRes),
		IconNameFileAudio:       NewThemedResource(fileaudioIconRes),
		IconNameFileImage:       NewThemedResource(fileimageIconRes),
		IconNameFileText:        NewThemedResource(filetextIconRes),
		IconNameFileVideo:       NewThemedResource(filevideoIconRes),
		IconNameFolder:          NewThemedResource(folderIconRes),
		IconNameFolderNew:       NewThemedResource(foldernewIconRes),
		IconNameFolderOpen:      NewThemedResource(folderopenIconRes),
		IconNameHelp:            NewThemedResource(helpIconRes),
		IconNameHistory:         NewThemedResource(historyIconRes),
		IconNameHome:            NewThemedResource(homeIconRes),
		IconNameSettings:        NewThemedResource(settingsIconRes),

		IconNameViewFullScreen: NewThemedResource(viewfullscreenIconRes),
		IconNameViewRefresh:    NewThemedResource(viewrefreshIconRes),
		IconNameViewRestore:    NewThemedResource(viewzoomfitIconRes),
		IconNameViewZoomFit:    NewThemedResource(viewzoomfitIconRes),
		IconNameViewZoomIn:     NewThemedResource(viewzoominIconRes),
		IconNameViewZoomOut:    NewThemedResource(viewzoomoutIconRes),

		IconNameVisibility:    NewThemedResource(visibilityIconRes),
		IconNameVisibilityOff: NewThemedResource(visibilityoffIconRes),

		IconNameVolumeDown: NewThemedResource(volumedownIconRes),
		IconNameVolumeMute: NewThemedResource(volumemuteIconRes),
		IconNameVolumeUp:   NewThemedResource(volumeupIconRes),

		IconNameDownload: NewThemedResource(downloadIconRes),
		IconNameComputer: NewThemedResource(computerIconRes),
		IconNameStorage:  NewThemedResource(storageIconRes),
		IconNameUpload:   NewThemedResource(uploadIconRes),

		IconNameAccount: NewThemedResource(accountIconRes),
		IconNameLogin:   NewThemedResource(loginIconRes),
		IconNameLogout:  NewThemedResource(logoutIconRes),

		IconNameList: NewThemedResource(listIconRes),
		IconNameGrid: NewThemedResource(gridIconRes),
	}
)

func (t *builtinTheme) Icon(n gui.ThemeIconName) gui.Resource {
	return icons[n]
}

// ThemedResource is a resource wrapper that will return a version of the resource with the main color changed
// for the currently selected theme.
type ThemedResource struct {
	source gui.Resource
}

// NewThemedResource creates a resource that adapts to the current theme setting.
func NewThemedResource(src gui.Resource) *ThemedResource {
	return &ThemedResource{
		source: src,
	}
}

// Name returns the underlying resource name (used for caching).
func (res *ThemedResource) Name() string {
	return res.source.Name()
}

// Content returns the underlying content of the resource adapted to the current text color.
func (res *ThemedResource) Content() []byte {
	return colorizeResource(res.source, ForegroundColor())
}

// Error returns a different resource for indicating an error.
func (res *ThemedResource) Error() *ErrorThemedResource {
	return NewErrorThemedResource(res)
}

// InvertedThemedResource is a resource wrapper that will return a version of the resource with the main color changed
// for use over highlighted elements.
type InvertedThemedResource struct {
	source gui.Resource
}

// NewInvertedThemedResource creates a resource that adapts to the current theme for use over highlighted elements.
func NewInvertedThemedResource(orig gui.Resource) *InvertedThemedResource {
	res := &InvertedThemedResource{source: orig}
	return res
}

// Name returns the underlying resource name (used for caching).
func (res *InvertedThemedResource) Name() string {
	return "inverted-" + res.source.Name()
}

// Content returns the underlying content of the resource adapted to the current background color.
func (res *InvertedThemedResource) Content() []byte {
	clr := BackgroundColor()
	return colorizeResource(res.source, clr)
}

// Original returns the underlying resource that this inverted themed resource was adapted from
func (res *InvertedThemedResource) Original() gui.Resource {
	return res.source
}

// ErrorThemedResource is a resource wrapper that will return a version of the resource with the main color changed
// to indicate an error.
type ErrorThemedResource struct {
	source gui.Resource
}

// NewErrorThemedResource creates a resource that adapts to the error color for the current theme.
func NewErrorThemedResource(orig gui.Resource) *ErrorThemedResource {
	res := &ErrorThemedResource{source: orig}
	return res
}

// Name returns the underlying resource name (used for caching).
func (res *ErrorThemedResource) Name() string {
	return "error-" + res.source.Name()
}

// Content returns the underlying content of the resource adapted to the current background color.
func (res *ErrorThemedResource) Content() []byte {
	return colorizeResource(res.source, ErrorColor())
}

// Original returns the underlying resource that this error themed resource was adapted from
func (res *ErrorThemedResource) Original() gui.Resource {
	return res.source
}

// PrimaryThemedResource is a resource wrapper that will return a version of the resource with the main color changed
// to the theme primary color.
type PrimaryThemedResource struct {
	source gui.Resource
}

// NewPrimaryThemedResource creates a resource that adapts to the primary color for the current theme.
func NewPrimaryThemedResource(orig gui.Resource) *PrimaryThemedResource {
	res := &PrimaryThemedResource{source: orig}
	return res
}

// Name returns the underlying resource name (used for caching).
func (res *PrimaryThemedResource) Name() string {
	return "primary-" + res.source.Name()
}

// Content returns the underlying content of the resource adapted to the current background color.
func (res *PrimaryThemedResource) Content() []byte {
	return colorizeResource(res.source, PrimaryColor())
}

// Original returns the underlying resource that this primary themed resource was adapted from
func (res *PrimaryThemedResource) Original() gui.Resource {
	return res.source
}

// DisabledResource is a resource wrapper that will return an appropriate resource colorized by
// the current theme's `DisabledColor` color.
type DisabledResource struct {
	source gui.Resource
}

// Name returns the resource source name prefixed with `disabled_` (used for caching)
func (res *DisabledResource) Name() string {
	return "disabled_" + res.source.Name()
}

// Content returns the disabled style content of the correct resource for the current theme
func (res *DisabledResource) Content() []byte {
	return colorizeResource(res.source, DisabledColor())
}

// NewDisabledResource creates a resource that adapts to the current theme's DisabledColor setting.
func NewDisabledResource(res gui.Resource) *DisabledResource {
	return &DisabledResource{
		source: res,
	}
}

func colorizeResource(res gui.Resource, clr color.Color) []byte {
	rdr := bytes.NewReader(res.Content())
	s, err := svgFromXML(rdr)
	if err != nil {
		gui.LogError("could not load SVG, falling back to static content:", err)
		return res.Content()
	}
	if err := s.replaceFillColor(clr); err != nil {
		gui.LogError("could not replace fill color, falling back to static content:", err)
		return res.Content()
	}
	b, err := xml.Marshal(s)
	if err != nil {
		gui.LogError("could not marshal svg, falling back to static content:", err)
		return res.Content()
	}
	return b
}

// BhojpurLogo returns a resource containing the Bhojpur GUI logo
func BhojpurLogo() gui.Resource {
	return bhojpurlogo
}

// CancelIcon returns a resource containing the standard cancel icon for the current theme
func CancelIcon() gui.Resource {
	return safeIconLookup(IconNameCancel)
}

// ConfirmIcon returns a resource containing the standard confirm icon for the current theme
func ConfirmIcon() gui.Resource {
	return safeIconLookup(IconNameConfirm)
}

// DeleteIcon returns a resource containing the standard delete icon for the current theme
func DeleteIcon() gui.Resource {
	return safeIconLookup(IconNameDelete)
}

// SearchIcon returns a resource containing the standard search icon for the current theme
func SearchIcon() gui.Resource {
	return safeIconLookup(IconNameSearch)
}

// SearchReplaceIcon returns a resource containing the standard search and replace icon for the current theme
func SearchReplaceIcon() gui.Resource {
	return safeIconLookup(IconNameSearchReplace)
}

// MenuIcon returns a resource containing the standard (mobile) menu icon for the current theme
func MenuIcon() gui.Resource {
	return safeIconLookup(IconNameMenu)
}

// MenuExpandIcon returns a resource containing the standard (mobile) expand "submenu icon for the current theme
func MenuExpandIcon() gui.Resource {
	return safeIconLookup(IconNameMenuExpand)
}

// CheckButtonIcon returns a resource containing the standard checkbox icon for the current theme
func CheckButtonIcon() gui.Resource {
	return safeIconLookup(IconNameCheckButton)
}

// CheckButtonCheckedIcon returns a resource containing the standard checkbox checked icon for the current theme
func CheckButtonCheckedIcon() gui.Resource {
	return safeIconLookup(IconNameCheckButtonChecked)
}

// RadioButtonIcon returns a resource containing the standard radio button icon for the current theme
func RadioButtonIcon() gui.Resource {
	return safeIconLookup(IconNameRadioButton)
}

// RadioButtonCheckedIcon returns a resource containing the standard radio button checked icon for the current theme
func RadioButtonCheckedIcon() gui.Resource {
	return safeIconLookup(IconNameRadioButtonChecked)
}

// ContentAddIcon returns a resource containing the standard content add icon for the current theme
func ContentAddIcon() gui.Resource {
	return safeIconLookup(IconNameContentAdd)
}

// ContentRemoveIcon returns a resource containing the standard content remove icon for the current theme
func ContentRemoveIcon() gui.Resource {
	return safeIconLookup(IconNameContentRemove)
}

// ContentClearIcon returns a resource containing the standard content clear icon for the current theme
func ContentClearIcon() gui.Resource {
	return safeIconLookup(IconNameContentClear)
}

// ContentCutIcon returns a resource containing the standard content cut icon for the current theme
func ContentCutIcon() gui.Resource {
	return safeIconLookup(IconNameContentCut)
}

// ContentCopyIcon returns a resource containing the standard content copy icon for the current theme
func ContentCopyIcon() gui.Resource {
	return safeIconLookup(IconNameContentCopy)
}

// ContentPasteIcon returns a resource containing the standard content paste icon for the current theme
func ContentPasteIcon() gui.Resource {
	return safeIconLookup(IconNameContentPaste)
}

// ContentRedoIcon returns a resource containing the standard content redo icon for the current theme
func ContentRedoIcon() gui.Resource {
	return safeIconLookup(IconNameContentRedo)
}

// ContentUndoIcon returns a resource containing the standard content undo icon for the current theme
func ContentUndoIcon() gui.Resource {
	return safeIconLookup(IconNameContentUndo)
}

// ColorAchromaticIcon returns a resource containing the standard achromatic color icon for the current theme
func ColorAchromaticIcon() gui.Resource {
	return safeIconLookup(IconNameColorAchromatic)
}

// ColorChromaticIcon returns a resource containing the standard chromatic color icon for the current theme
func ColorChromaticIcon() gui.Resource {
	return safeIconLookup(IconNameColorChromatic)
}

// ColorPaletteIcon returns a resource containing the standard color palette icon for the current theme
func ColorPaletteIcon() gui.Resource {
	return safeIconLookup(IconNameColorPalette)
}

// DocumentIcon returns a resource containing the standard document icon for the current theme
func DocumentIcon() gui.Resource {
	return safeIconLookup(IconNameDocument)
}

// DocumentCreateIcon returns a resource containing the standard document create icon for the current theme
func DocumentCreateIcon() gui.Resource {
	return safeIconLookup(IconNameDocumentCreate)
}

// DocumentPrintIcon returns a resource containing the standard document print icon for the current theme
func DocumentPrintIcon() gui.Resource {
	return safeIconLookup(IconNameDocumentPrint)
}

// DocumentSaveIcon returns a resource containing the standard document save icon for the current theme
func DocumentSaveIcon() gui.Resource {
	return safeIconLookup(IconNameDocumentSave)
}

// MoreHorizontalIcon returns a resource containing the standard horizontal more icon for the current theme
func MoreHorizontalIcon() gui.Resource {
	return current().Icon(IconNameMoreHorizontal)
}

// MoreVerticalIcon returns a resource containing the standard vertical more icon for the current theme
func MoreVerticalIcon() gui.Resource {
	return current().Icon(IconNameMoreVertical)
}

// InfoIcon returns a resource containing the standard dialog info icon for the current theme
func InfoIcon() gui.Resource {
	return safeIconLookup(IconNameInfo)
}

// QuestionIcon returns a resource containing the standard dialog question icon for the current theme
func QuestionIcon() gui.Resource {
	return safeIconLookup(IconNameQuestion)
}

// WarningIcon returns a resource containing the standard dialog warning icon for the current theme
func WarningIcon() gui.Resource {
	return safeIconLookup(IconNameWarning)
}

// ErrorIcon returns a resource containing the standard dialog error icon for the current theme
func ErrorIcon() gui.Resource {
	return safeIconLookup(IconNameError)
}

// FileIcon returns a resource containing the appropriate file icon for the current theme
func FileIcon() gui.Resource {
	return safeIconLookup(IconNameFile)
}

// FileApplicationIcon returns a resource containing the file icon representing application files for the current theme
func FileApplicationIcon() gui.Resource {
	return safeIconLookup(IconNameFileApplication)
}

// FileAudioIcon returns a resource containing the file icon representing audio files for the current theme
func FileAudioIcon() gui.Resource {
	return safeIconLookup(IconNameFileAudio)
}

// FileImageIcon returns a resource containing the file icon representing image files for the current theme
func FileImageIcon() gui.Resource {
	return safeIconLookup(IconNameFileImage)
}

// FileTextIcon returns a resource containing the file icon representing text files for the current theme
func FileTextIcon() gui.Resource {
	return safeIconLookup(IconNameFileText)
}

// FileVideoIcon returns a resource containing the file icon representing video files for the current theme
func FileVideoIcon() gui.Resource {
	return safeIconLookup(IconNameFileVideo)
}

// FolderIcon returns a resource containing the standard folder icon for the current theme
func FolderIcon() gui.Resource {
	return safeIconLookup(IconNameFolder)
}

// FolderNewIcon returns a resource containing the standard folder creation icon for the current theme
func FolderNewIcon() gui.Resource {
	return safeIconLookup(IconNameFolderNew)
}

// FolderOpenIcon returns a resource containing the standard folder open icon for the current theme
func FolderOpenIcon() gui.Resource {
	return safeIconLookup(IconNameFolderOpen)
}

// HelpIcon returns a resource containing the standard help icon for the current theme
func HelpIcon() gui.Resource {
	return safeIconLookup(IconNameHelp)
}

// HistoryIcon returns a resource containing the standard history icon for the current theme
func HistoryIcon() gui.Resource {
	return safeIconLookup(IconNameHistory)
}

// HomeIcon returns a resource containing the standard home folder icon for the current theme
func HomeIcon() gui.Resource {
	return safeIconLookup(IconNameHome)
}

// SettingsIcon returns a resource containing the standard settings icon for the current theme
func SettingsIcon() gui.Resource {
	return safeIconLookup(IconNameSettings)
}

// MailAttachmentIcon returns a resource containing the standard mail attachment icon for the current theme
func MailAttachmentIcon() gui.Resource {
	return safeIconLookup(IconNameMailAttachment)
}

// MailComposeIcon returns a resource containing the standard mail compose icon for the current theme
func MailComposeIcon() gui.Resource {
	return safeIconLookup(IconNameMailCompose)
}

// MailForwardIcon returns a resource containing the standard mail forward icon for the current theme
func MailForwardIcon() gui.Resource {
	return safeIconLookup(IconNameMailForward)
}

// MailReplyIcon returns a resource containing the standard mail reply icon for the current theme
func MailReplyIcon() gui.Resource {
	return safeIconLookup(IconNameMailReply)
}

// MailReplyAllIcon returns a resource containing the standard mail reply all icon for the current theme
func MailReplyAllIcon() gui.Resource {
	return safeIconLookup(IconNameMailReplyAll)
}

// MailSendIcon returns a resource containing the standard mail send icon for the current theme
func MailSendIcon() gui.Resource {
	return safeIconLookup(IconNameMailSend)
}

// MediaMusicIcon returns a resource containing the standard media music icon for the current theme
//
// Since: 2.1
func MediaMusicIcon() gui.Resource {
	return safeIconLookup(IconNameMediaMusic)
}

// MediaPhotoIcon returns a resource containing the standard media photo icon for the current theme
//
// Since: 2.1
func MediaPhotoIcon() gui.Resource {
	return safeIconLookup(IconNameMediaPhoto)
}

// MediaVideoIcon returns a resource containing the standard media video icon for the current theme
//
// Since: 2.1
func MediaVideoIcon() gui.Resource {
	return safeIconLookup(IconNameMediaVideo)
}

// MediaFastForwardIcon returns a resource containing the standard media fast-forward icon for the current theme
func MediaFastForwardIcon() gui.Resource {
	return safeIconLookup(IconNameMediaFastForward)
}

// MediaFastRewindIcon returns a resource containing the standard media fast-rewind icon for the current theme
func MediaFastRewindIcon() gui.Resource {
	return safeIconLookup(IconNameMediaFastRewind)
}

// MediaPauseIcon returns a resource containing the standard media pause icon for the current theme
func MediaPauseIcon() gui.Resource {
	return safeIconLookup(IconNameMediaPause)
}

// MediaPlayIcon returns a resource containing the standard media play icon for the current theme
func MediaPlayIcon() gui.Resource {
	return safeIconLookup(IconNameMediaPlay)
}

// MediaRecordIcon returns a resource containing the standard media record icon for the current theme
func MediaRecordIcon() gui.Resource {
	return safeIconLookup(IconNameMediaRecord)
}

// MediaReplayIcon returns a resource containing the standard media replay icon for the current theme
func MediaReplayIcon() gui.Resource {
	return safeIconLookup(IconNameMediaReplay)
}

// MediaSkipNextIcon returns a resource containing the standard media skip next icon for the current theme
func MediaSkipNextIcon() gui.Resource {
	return safeIconLookup(IconNameMediaSkipNext)
}

// MediaSkipPreviousIcon returns a resource containing the standard media skip previous icon for the current theme
func MediaSkipPreviousIcon() gui.Resource {
	return safeIconLookup(IconNameMediaSkipPrevious)
}

// MediaStopIcon returns a resource containing the standard media stop icon for the current theme
func MediaStopIcon() gui.Resource {
	return safeIconLookup(IconNameMediaStop)
}

// MoveDownIcon returns a resource containing the standard down arrow icon for the current theme
func MoveDownIcon() gui.Resource {
	return safeIconLookup(IconNameMoveDown)
}

// MoveUpIcon returns a resource containing the standard up arrow icon for the current theme
func MoveUpIcon() gui.Resource {
	return safeIconLookup(IconNameMoveUp)
}

// NavigateBackIcon returns a resource containing the standard backward navigation icon for the current theme
func NavigateBackIcon() gui.Resource {
	return safeIconLookup(IconNameNavigateBack)
}

// NavigateNextIcon returns a resource containing the standard forward navigation icon for the current theme
func NavigateNextIcon() gui.Resource {
	return safeIconLookup(IconNameNavigateNext)
}

// MenuDropDownIcon returns a resource containing the standard menu drop down icon for the current theme
func MenuDropDownIcon() gui.Resource {
	return safeIconLookup(IconNameArrowDropDown)
}

// MenuDropUpIcon returns a resource containing the standard menu drop up icon for the current theme
func MenuDropUpIcon() gui.Resource {
	return safeIconLookup(IconNameArrowDropUp)
}

// ViewFullScreenIcon returns a resource containing the standard fullscreen icon for the current theme
func ViewFullScreenIcon() gui.Resource {
	return safeIconLookup(IconNameViewFullScreen)
}

// ViewRestoreIcon returns a resource containing the standard exit fullscreen icon for the current theme
func ViewRestoreIcon() gui.Resource {
	return safeIconLookup(IconNameViewRestore)
}

// ViewRefreshIcon returns a resource containing the standard refresh icon for the current theme
func ViewRefreshIcon() gui.Resource {
	return safeIconLookup(IconNameViewRefresh)
}

// ZoomFitIcon returns a resource containing the standard zoom fit icon for the current theme
func ZoomFitIcon() gui.Resource {
	return safeIconLookup(IconNameViewZoomFit)
}

// ZoomInIcon returns a resource containing the standard zoom in icon for the current theme
func ZoomInIcon() gui.Resource {
	return safeIconLookup(IconNameViewZoomIn)
}

// ZoomOutIcon returns a resource containing the standard zoom out icon for the current theme
func ZoomOutIcon() gui.Resource {
	return safeIconLookup(IconNameViewZoomOut)
}

// VisibilityIcon returns a resource containing the standard visibility icon for the current theme
func VisibilityIcon() gui.Resource {
	return safeIconLookup(IconNameVisibility)
}

// VisibilityOffIcon returns a resource containing the standard visibility off icon for the current theme
func VisibilityOffIcon() gui.Resource {
	return safeIconLookup(IconNameVisibilityOff)
}

// VolumeDownIcon returns a resource containing the standard volume down icon for the current theme
func VolumeDownIcon() gui.Resource {
	return safeIconLookup(IconNameVolumeDown)
}

// VolumeMuteIcon returns a resource containing the standard volume mute icon for the current theme
func VolumeMuteIcon() gui.Resource {
	return safeIconLookup(IconNameVolumeMute)
}

// VolumeUpIcon returns a resource containing the standard volume up icon for the current theme
func VolumeUpIcon() gui.Resource {
	return safeIconLookup(IconNameVolumeUp)
}

// ComputerIcon returns a resource containing the standard computer icon for the current theme
func ComputerIcon() gui.Resource {
	return safeIconLookup(IconNameComputer)
}

// DownloadIcon returns a resource containing the standard download icon for the current theme
func DownloadIcon() gui.Resource {
	return safeIconLookup(IconNameDownload)
}

// StorageIcon returns a resource containing the standard storage icon for the current theme
func StorageIcon() gui.Resource {
	return safeIconLookup(IconNameStorage)
}

// UploadIcon returns a resource containing the standard upload icon for the current theme
func UploadIcon() gui.Resource {
	return safeIconLookup(IconNameUpload)
}

// AccountIcon returns a resource containing the standard account icon for the current theme
func AccountIcon() gui.Resource {
	return safeIconLookup(IconNameAccount)
}

// LoginIcon returns a resource containing the standard login icon for the current theme
func LoginIcon() gui.Resource {
	return safeIconLookup(IconNameLogin)
}

// LogoutIcon returns a resource containing the standard logout icon for the current theme
func LogoutIcon() gui.Resource {
	return safeIconLookup(IconNameLogout)
}

// ListIcon returns a resource containing the standard list icon for the current theme
func ListIcon() gui.Resource {
	return safeIconLookup(IconNameList)
}

// GridIcon returns a resource containing the standard grid icon for the current theme
func GridIcon() gui.Resource {
	return safeIconLookup(IconNameGrid)
}

func safeIconLookup(n gui.ThemeIconName) gui.Resource {
	icon := current().Icon(n)
	if icon == nil {
		gui.LogError("Loaded theme returned nil icon", nil)
		return fallbackIcon
	}
	return icon
}
