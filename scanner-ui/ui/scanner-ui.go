package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type ScannerUi struct {
	app       fyne.App
	window    fyne.Window
	addDialog *AddDialog
}

func NewScannerUi() *ScannerUi {
	app := app.New()
	window := app.NewWindow("Go-Movie!")
	addDialog := NewAddDialog(&window, app)
	ui := &ScannerUi{
		app:       app,
		window:    window,
		addDialog: addDialog,
	}
	return ui
}

func (ui *ScannerUi) InitAndRun() {
	ui.window.SetMainMenu(ui.makeMenu(ui.app, ui.window))
	ui.window.Resize(fyne.NewSize(1200, 800))
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
