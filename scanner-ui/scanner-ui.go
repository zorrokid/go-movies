package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
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
	newItem := fyne.NewMenuItem("New", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				log.Println("Cancelled")
				return
			}
			defer reader.Close()
		}, w)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
		fd.Show()
	})

	file := fyne.NewMenu("File", newItem)
	return fyne.NewMainMenu(
		file,
	)
}
