package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ScannerUi struct {
	app          fyne.App
	window       fyne.Window
	addDialog    *AddDialog
	fileListData []string
}

func NewScannerUi() *ScannerUi {
	app := app.NewWithID("github.com/zorrokid/go-movies/scanner-ui")
	window := app.NewWindow("Go-Movie!")
	addDialog := NewAddDialog(&app)
	ui := &ScannerUi{
		app:          app,
		window:       window,
		addDialog:    addDialog,
		fileListData: []string{"aaa", "bbb", "ccc"},
	}
	return ui
}

func (ui *ScannerUi) InitAndRun() {
	ui.window.SetMainMenu(ui.makeMenu(ui.app, ui.window))
	ui.window.Resize(fyne.NewSize(1200, 800))
	ui.window.SetContent(ui.createFileList())
	ui.window.ShowAndRun()
}

func (ui *ScannerUi) makeMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {
	newItem := fyne.NewMenuItem("New", func() {
		ui.addDialog.ShowDialog()
	})

	file := fyne.NewMenu("File", newItem)
	return fyne.NewMainMenu(
		file,
	)
}

func (ui *ScannerUi) createFileList() *widget.List {

	list := widget.NewList(
		// length
		func() int {
			return len(ui.fileListData)
		},
		// create item
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
		},
		// update item
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(ui.fileListData[id])
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		fmt.Println(ui.fileListData[id])

	}
	list.OnUnselected = func(id widget.ListItemID) {
		fmt.Println("OnUnselected")
	}

	return list
}
