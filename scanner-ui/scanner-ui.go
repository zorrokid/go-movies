package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type ScannerUi struct {
	app    fyne.App
	window fyne.Window
}

func NewScannerUi() *ScannerUi {
	app := app.New()
	window := app.NewWindow("Go-Retro!")
	ui := &ScannerUi{
		app:    app,
		window: window,
	}
	return ui
}

func (ui *ScannerUi) InitAndRun() {
	ui.window.SetMainMenu(ui.makeMenu(ui.app, ui.window))
	ui.window.ShowAndRun()
}

func (ui *ScannerUi) makeMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {
	newItem := fyne.NewMenuItem("New", nil)

	file := fyne.NewMenu("File", newItem)
	return fyne.NewMainMenu(
		file,
	)
}
