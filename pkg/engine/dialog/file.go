package dialog

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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/container"
	"github.com/bhojpur/gui/pkg/engine/storage"
	"github.com/bhojpur/gui/pkg/engine/storage/repository"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

type viewLayout int

const (
	gridView viewLayout = iota
	listView
)

type textWidget interface {
	gui.Widget
	SetText(string)
}

type favoriteItem struct {
	locName string
	locIcon gui.Resource
	loc     gui.URI
}

type fileDialog struct {
	file             *FileDialog
	fileName         textWidget
	dismiss          *widget.Button
	open             *widget.Button
	breadcrumb       *gui.Container
	breadcrumbScroll *container.Scroll
	files            *gui.Container
	filesScroll      *container.Scroll
	favorites        []favoriteItem
	favoritesList    *widget.List
	showHidden       bool

	view viewLayout

	win      *widget.PopUp
	selected *fileDialogItem
	dir      gui.ListableURI
	// this will be the initial filename in a FileDialog in save mode
	initialFileName string
}

// FileDialog is a dialog containing a file picker for use in opening or saving files.
type FileDialog struct {
	save             bool
	callback         interface{}
	onClosedCallback func(bool)
	filter           storage.FileFilter
	parent           gui.Window
	dialog           *fileDialog
	dismissText      string
	desiredSize      gui.Size
	// this will be applied to dialog.dir when it's loaded
	startingLocation gui.ListableURI
	// this will be the initial filename in a FileDialog in save mode
	initialFileName string
}

// Declare conformity to Dialog interface
var _ Dialog = (*FileDialog)(nil)

func (f *fileDialog) makeUI() gui.CanvasObject {
	if f.file.save {
		saveName := widget.NewEntry()
		saveName.OnChanged = func(s string) {
			if s == "" {
				f.open.Disable()
			} else {
				f.open.Enable()
			}
		}
		saveName.SetPlaceHolder("Enter filename")
		f.fileName = saveName
	} else {
		f.fileName = widget.NewLabel("")
	}

	label := "Open"
	if f.file.save {
		label = "Save"
	}
	f.open = widget.NewButton(label, func() {
		if f.file.callback == nil {
			f.win.Hide()
			if f.file.onClosedCallback != nil {
				f.file.onClosedCallback(false)
			}
			return
		}

		if f.file.save {
			callback := f.file.callback.(func(gui.URIWriteCloser, error))
			name := f.fileName.(*widget.Entry).Text
			location, _ := storage.Child(f.dir, name)

			exists, _ := storage.Exists(location)

			// check if a directory is selected
			listable, err := storage.CanList(location)

			if !exists {
				f.win.Hide()
				if f.file.onClosedCallback != nil {
					f.file.onClosedCallback(true)
				}
				callback(storage.Writer(location))
				return
			} else if err == nil && listable {
				// a directory has been selected
				ShowInformation("Cannot overwrite",
					"Files cannot replace a directory,\ncheck the file name and try again", f.file.parent)
				return
			}

			ShowConfirm("Overwrite?", "Are you sure you want to overwrite the file\n"+name+"?",
				func(ok bool) {
					if !ok {
						return
					}
					f.win.Hide()

					callback(storage.Writer(location))
					if f.file.onClosedCallback != nil {
						f.file.onClosedCallback(true)
					}
				}, f.file.parent)
		} else if f.selected != nil {
			callback := f.file.callback.(func(gui.URIReadCloser, error))
			f.win.Hide()
			if f.file.onClosedCallback != nil {
				f.file.onClosedCallback(true)
			}
			callback(storage.Reader(f.selected.location))
		} else if f.file.isDirectory() {
			callback := f.file.callback.(func(gui.ListableURI, error))
			f.win.Hide()
			if f.file.onClosedCallback != nil {
				f.file.onClosedCallback(true)
			}
			callback(f.dir, nil)
		}
	})
	f.open.Importance = widget.HighImportance
	f.open.Disable()
	if f.file.save {
		f.fileName.SetText(f.initialFileName)
	}
	dismissLabel := "Cancel"
	if f.file.dismissText != "" {
		dismissLabel = f.file.dismissText
	}
	f.dismiss = widget.NewButton(dismissLabel, func() {
		f.win.Hide()
		if f.file.onClosedCallback != nil {
			f.file.onClosedCallback(false)
		}
		if f.file.callback != nil {
			if f.file.save {
				f.file.callback.(func(gui.URIWriteCloser, error))(nil, nil)
			} else if f.file.isDirectory() {
				f.file.callback.(func(gui.ListableURI, error))(nil, nil)
			} else {
				f.file.callback.(func(gui.URIReadCloser, error))(nil, nil)
			}
		}
	})
	buttons := container.NewGridWithRows(1, f.dismiss, f.open)

	f.filesScroll = container.NewScroll(nil) // filesScroll's content will be set by setView function.
	verticalExtra := float32(float64(fileIconSize) * 0.25)
	f.filesScroll.SetMinSize(gui.NewSize(fileIconCellWidth*2+theme.Padding(),
		(fileIconSize+fileTextSize)+theme.Padding()*2+verticalExtra))

	f.breadcrumb = container.NewHBox()
	f.breadcrumbScroll = container.NewHScroll(container.NewPadded(f.breadcrumb))
	title := label + " File"
	if f.file.isDirectory() {
		title = label + " Folder"
	}

	f.setView(gridView)
	f.loadFavorites()

	f.favoritesList = widget.NewList(
		func() int {
			return len(f.favorites)
		},
		func() gui.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item gui.CanvasObject) {
			item.(*gui.Container).Objects[0].(*widget.Icon).SetResource(f.favorites[id].locIcon)
			item.(*gui.Container).Objects[1].(*widget.Label).SetText(f.favorites[id].locName)
		},
	)
	f.favoritesList.OnSelected = func(id widget.ListItemID) {
		f.setLocation(f.favorites[id].loc)
	}

	var optionsButton *widget.Button
	optionsButton = widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		f.optionsMenu(gui.CurrentApp().Driver().AbsolutePositionForObject(optionsButton), optionsButton.Size())
	})

	var toggleViewButton *widget.Button
	toggleViewButton = widget.NewButtonWithIcon("", theme.ListIcon(), func() {
		if f.view == gridView {
			f.setView(listView)
			toggleViewButton.SetIcon(theme.GridIcon())
		} else {
			f.setView(gridView)
			toggleViewButton.SetIcon(theme.ListIcon())
		}
	})

	optionsbuttons := container.NewHBox(
		toggleViewButton,
		optionsButton,
	)

	header := container.NewBorder(nil, nil, nil, optionsbuttons,
		optionsbuttons, widget.NewLabelWithStyle(title, gui.TextAlignLeading, gui.TextStyle{Bold: true}),
	)

	footer := container.NewBorder(nil, nil, nil, buttons,
		buttons, container.NewHScroll(f.fileName),
	)

	body := container.NewHSplit(
		f.favoritesList,
		container.NewBorder(f.breadcrumbScroll, nil, nil, nil,
			f.breadcrumbScroll, f.filesScroll,
		),
	)
	body.SetOffset(0) // Set the minimum offset so that the favoritesList takes only it's minimal width

	return container.NewBorder(header, footer, nil, nil, body)
}

func (f *fileDialog) optionsMenu(position gui.Position, buttonSize gui.Size) {
	hiddenFiles := widget.NewCheck("Show Hidden Files", func(changed bool) {
		f.showHidden = changed
		f.refreshDir(f.dir)
	})
	hiddenFiles.Checked = f.showHidden
	hiddenFiles.Refresh()
	content := container.NewVBox(hiddenFiles)

	p := position.Add(buttonSize)
	pos := gui.NewPos(p.X-content.MinSize().Width-theme.Padding()*2, p.Y+theme.Padding()*2)
	widget.ShowPopUpAtPosition(content, f.win.Canvas, pos)
}

func (f *fileDialog) loadFavorites() {
	favoriteLocations, err := getFavoriteLocations()
	if err != nil {
		gui.LogError("Getting favorite locations", err)
	}
	favoriteIcons := getFavoriteIcons()
	favoriteOrder := getFavoriteOrder()

	f.favorites = []favoriteItem{
		{locName: "Home", locIcon: theme.HomeIcon(), loc: favoriteLocations["Home"]}}
	app := gui.CurrentApp()
	if hasAppFiles(app) {
		f.favorites = append(f.favorites,
			favoriteItem{locName: "App Files", locIcon: theme.FileIcon(), loc: storageURI(app)})
	}
	f.favorites = append(f.favorites, f.getPlaces()...)

	for _, locName := range favoriteOrder {
		loc, ok := favoriteLocations[locName]
		if !ok {
			continue
		}
		locIcon := favoriteIcons[locName]
		f.favorites = append(f.favorites,
			favoriteItem{locName: locName, locIcon: locIcon, loc: loc})
	}
}

func (f *fileDialog) refreshDir(dir gui.ListableURI) {
	f.files.Objects = nil

	files, err := dir.List()
	if err != nil {
		gui.LogError("Unable to read ListableURI "+dir.String(), err)
		return
	}

	var icons []gui.CanvasObject
	parent, err := storage.Parent(dir)
	if err != nil && err != repository.ErrURIRoot {
		gui.LogError("Unable to get parent of "+dir.String(), err)
		return
	}
	if parent != nil && parent.String() != dir.String() {
		fi := &fileDialogItem{picker: f, name: "(Parent)", location: parent, dir: true}
		fi.ExtendBaseWidget(fi)
		icons = append(icons, fi)
	}

	for _, file := range files {
		if !f.showHidden && isHidden(file) {
			continue
		}

		listable, err := storage.CanList(file)
		if f.file.isDirectory() && err != nil {
			continue
		} else if err == nil && listable { // URI points to a directory
			icons = append(icons, f.newFileItem(file, true)) // Pass the listable URI to avoid doing the same check in FileIcon
		} else if f.file.filter == nil || f.file.filter.Matches(file) {
			icons = append(icons, f.newFileItem(file, false))
		}
	}

	f.files.Objects = icons
	f.files.Refresh()
	f.filesScroll.Offset = gui.NewPos(0, 0)
	f.filesScroll.Refresh()
}

func (f *fileDialog) setLocation(dir gui.URI) error {
	if dir == nil {
		return fmt.Errorf("failed to open nil directory")
	}
	list, err := storage.ListerForURI(dir)
	if err != nil {
		return err
	}

	isFav := false
	for i, fav := range f.favorites {
		if fav.loc == nil {
			continue
		}
		if fav.loc.Path() == dir.Path() {
			f.favoritesList.Select(i)
			isFav = true
			break
		}
	}
	if !isFav {
		f.favoritesList.UnselectAll()
	}

	f.setSelected(nil)
	f.dir = list

	f.breadcrumb.Objects = nil

	localdir := dir.String()[len(dir.Scheme())+3:]

	buildDir := filepath.VolumeName(localdir)
	for i, d := range strings.Split(localdir, "/") {
		if d == "" {
			if i > 0 { // what we get if we split "/"
				break
			}
			buildDir = "/"
			d = "/"
		} else if i > 0 {
			buildDir = filepath.Join(buildDir, d)
		} else {
			d = buildDir
			buildDir = d + string(os.PathSeparator)
		}

		newDir := storage.NewFileURI(buildDir)
		isDir, err := storage.CanList(newDir)
		if err != nil {
			return err
		}

		if !isDir {
			return errors.New("location was not a listable URI")
		}
		f.breadcrumb.Add(
			widget.NewButton(d, func() {
				err := f.setLocation(newDir)
				if err != nil {
					gui.LogError("Failed to set directory", err)
				}
			}),
		)
	}

	f.breadcrumbScroll.Refresh()
	f.breadcrumbScroll.Offset.X = f.breadcrumbScroll.Content.Size().Width - f.breadcrumbScroll.Size().Width
	f.breadcrumbScroll.Refresh()

	if f.file.isDirectory() {
		f.fileName.SetText(dir.Name())
		f.open.Enable()
	}
	f.refreshDir(list)

	return nil
}

func (f *fileDialog) setSelected(file *fileDialogItem) {
	if f.selected != nil {
		f.selected.isCurrent = false
		f.selected.Refresh()
	}
	if file != nil && file.isDirectory() {
		listable, err := storage.CanList(file.location)
		if err != nil || !listable {
			gui.LogError("Failed to create lister for URI"+file.location.String(), err)
		}
		f.setLocation(file.location)
		return
	}
	f.selected = file

	if file == nil || file.location.String()[len(file.location.Scheme())+3:] == "" {
		// keep user input while navigating
		// in a FileSave dialog
		if !f.file.save {
			f.fileName.SetText("")
			f.open.Disable()
		}
	} else {
		file.isCurrent = true
		f.fileName.SetText(file.location.Name())
		f.open.Enable()
	}
}

func (f *fileDialog) setView(view viewLayout) {
	f.view = view
	if f.view == gridView {
		padding := gui.NewSize(fileIconCellWidth-fileIconSize, theme.Padding())
		f.files = container.NewGridWrap(
			gui.NewSize(fileIconSize, fileIconSize+fileTextSize).Add(padding),
		)
	} else {
		f.files = container.NewVBox()
	}
	if f.dir != nil {
		f.refreshDir(f.dir)
	}
	f.filesScroll.Content = container.NewPadded(f.files)
	f.filesScroll.Refresh()
}

// effectiveStartingDir calculates the directory at which the file dialog should
// open, based on the values of startingDirectory, CWD, home, and any error
// conditions which occur.
//
// Order of precedence is:
//
// * file.startingDirectory if non-empty, os.Stat()-able, and uses the file://
//   URI scheme
// * os.UserHomeDir()
// * os.Getwd()
// * "/" (should be filesystem root on all supported platforms)
//
func (f *FileDialog) effectiveStartingDir() gui.ListableURI {
	if f.startingLocation != nil {
		if f.startingLocation.Scheme() == "file" {
			path := f.startingLocation.Path()

			// the starting directory is set explicitly
			if _, err := os.Stat(path); err != nil {
				gui.LogError("Error with StartingLocation", err)
			} else {
				return f.startingLocation
			}
		}

	}

	// Try app storage
	app := gui.CurrentApp()
	if hasAppFiles(app) {
		list, _ := storage.ListerForURI(storageURI(app))
		return list
	}

	// Try home dir
	dir, err := os.UserHomeDir()
	if err == nil {
		lister, err := storage.ListerForURI(storage.NewFileURI(dir))
		if err == nil {
			return lister
		}
		gui.LogError("Could not create lister for user home dir", err)
	}
	gui.LogError("Could not load user home dir", err)

	// Try to get ./
	wd, err := os.Getwd()
	if err == nil {
		lister, err := storage.ListerForURI(storage.NewFileURI(wd))
		if err == nil {
			return lister
		}
		gui.LogError("Could not create lister for working dir", err)
	}

	lister, err := storage.ListerForURI(storage.NewFileURI("/"))
	if err != nil {
		gui.LogError("could not create lister for /", err)
		return nil
	}
	return lister
}

func showFile(file *FileDialog) *fileDialog {
	d := &fileDialog{file: file, initialFileName: file.initialFileName}
	ui := d.makeUI()
	size := ui.MinSize().Add(gui.NewSize(fileIconCellWidth*2+theme.Padding()*6+theme.Padding(),
		(fileIconSize+fileTextSize)+theme.Padding()*6))

	d.win = widget.NewModalPopUp(ui, file.parent.Canvas())
	d.win.Resize(size)

	d.setLocation(file.effectiveStartingDir())
	d.win.Show()
	return d
}

// MinSize returns the size that this dialog should not shrink below
//
// Since: 2.1
func (f *FileDialog) MinSize() gui.Size {
	return f.dialog.win.MinSize()
}

// Show shows the file dialog.
func (f *FileDialog) Show() {
	if f.save {
		if fileSaveOSOverride(f) {
			return
		}
	} else {
		if fileOpenOSOverride(f) {
			return
		}
	}
	if f.dialog != nil {
		f.dialog.win.Show()
		return
	}
	f.dialog = showFile(f)
	if !f.desiredSize.IsZero() {
		f.Resize(f.desiredSize)
	}
}

// Refresh causes this dialog to be updated
func (f *FileDialog) Refresh() {
	f.dialog.win.Refresh()
}

// Resize dialog to the requested size, if there is sufficient space.
// If the parent window is not large enough then the size will be reduced to fit.
func (f *FileDialog) Resize(size gui.Size) {
	f.desiredSize = size
	if f.dialog == nil {
		return
	}
	f.dialog.win.Resize(size)
}

// Hide hides the file dialog.
func (f *FileDialog) Hide() {
	if f.dialog == nil {
		return
	}
	f.dialog.win.Hide()
	if f.onClosedCallback != nil {
		f.onClosedCallback(false)
	}
}

// SetDismissText allows custom text to be set in the confirmation button
func (f *FileDialog) SetDismissText(label string) {
	if f.dialog == nil {
		return
	}
	f.dialog.dismiss.SetText(label)
	f.dialog.win.Refresh()
}

// SetLocation tells this FileDirectory which location to display.
// This is normally called before the dialog is shown.
//
// Since: 1.4
func (f *FileDialog) SetLocation(u gui.ListableURI) {
	f.startingLocation = u
	if f.dialog != nil {
		f.dialog.setLocation(u)
	}
}

// SetOnClosed sets a callback function that is called when
// the dialog is closed.
func (f *FileDialog) SetOnClosed(closed func()) {
	if f.dialog == nil {
		return
	}
	// If there is already a callback set, remember it and call both.
	originalCallback := f.onClosedCallback

	f.onClosedCallback = func(response bool) {
		closed()
		if originalCallback != nil {
			originalCallback(response)
		}
	}
}

// SetFilter sets a filter for limiting files that can be chosen in the file dialog.
func (f *FileDialog) SetFilter(filter storage.FileFilter) {
	if f.isDirectory() {
		gui.LogError("Cannot set a filter for a folder dialog", nil)
		return
	}
	f.filter = filter
	if f.dialog != nil {
		f.dialog.refreshDir(f.dialog.dir)
	}
}

// SetFileName sets the filename in a FileDialog in save mode.
// This is normally called before the dialog is shown.
func (f *FileDialog) SetFileName(fileName string) {
	if f.save {
		f.initialFileName = fileName
		//Update entry if fileDialog has already been created
		if f.dialog != nil {
			f.dialog.fileName.SetText(fileName)
		}
	}
}

// NewFileOpen creates a file dialog allowing the user to choose a file to open.
// The callback function will run when the dialog closes. The URI will be nil
// when the user cancels or when nothing is selected.
//
// The dialog will appear over the window specified when Show() is called.
func NewFileOpen(callback func(gui.URIReadCloser, error), parent gui.Window) *FileDialog {
	dialog := &FileDialog{callback: callback, parent: parent}
	return dialog
}

// NewFileSave creates a file dialog allowing the user to choose a file to save
// to (new or overwrite). If the user chooses an existing file they will be
// asked if they are sure. The callback function will run when the dialog
// closes. The URI will be nil when the user cancels or when nothing is
// selected.
//
// The dialog will appear over the window specified when Show() is called.
func NewFileSave(callback func(gui.URIWriteCloser, error), parent gui.Window) *FileDialog {
	dialog := &FileDialog{callback: callback, parent: parent, save: true}
	return dialog
}

// ShowFileOpen creates and shows a file dialog allowing the user to choose a
// file to open. The callback function will run when the dialog closes. The URI
// will be nil when the user cancels or when nothing is selected.
//
// The dialog will appear over the window specified.
func ShowFileOpen(callback func(gui.URIReadCloser, error), parent gui.Window) {
	dialog := NewFileOpen(callback, parent)
	if fileOpenOSOverride(dialog) {
		return
	}
	dialog.Show()
}

// ShowFileSave creates and shows a file dialog allowing the user to choose a
// file to save to (new or overwrite). If the user chooses an existing file they
// will be asked if they are sure. The callback function will run when the
// dialog closes. The URI will be nil when the user cancels or when nothing is
// selected.
//
// The dialog will appear over the window specified.
func ShowFileSave(callback func(gui.URIWriteCloser, error), parent gui.Window) {
	dialog := NewFileSave(callback, parent)
	if fileSaveOSOverride(dialog) {
		return
	}
	dialog.Show()
}

func getFavoriteIcons() map[string]gui.Resource {
	if runtime.GOOS == "darwin" {
		return map[string]gui.Resource{
			"Documents": theme.DocumentIcon(),
			"Downloads": theme.DownloadIcon(),
			"Music":     theme.MediaMusicIcon(),
			"Pictures":  theme.MediaPhotoIcon(),
			"Movies":    theme.MediaVideoIcon(),
		}
	}

	return map[string]gui.Resource{
		"Documents": theme.DocumentIcon(),
		"Downloads": theme.DownloadIcon(),
		"Music":     theme.MediaMusicIcon(),
		"Pictures":  theme.MediaPhotoIcon(),
		"Videos":    theme.MediaVideoIcon(),
	}
}

func getFavoriteOrder() []string {
	order := []string{
		"Documents",
		"Downloads",
		"Music",
		"Pictures",
		"Videos",
	}

	if runtime.GOOS == "darwin" {
		order[4] = "Movies"
	}

	return order
}

func hasAppFiles(a gui.App) bool {
	return len(a.Storage().List()) > 0
}

func storageURI(a gui.App) gui.URI {
	dir, _ := storage.Child(a.Storage().RootURI(), "Documents")
	return dir
}
