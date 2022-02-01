package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type AddDialog struct {
	imageContainer *fyne.Container
	mainWindow     *fyne.Window
	scanWindow     fyne.Window
	app            fyne.App
}

func NewAddDialog(w *fyne.Window, app fyne.App) *AddDialog {
	container := container.New(layout.NewBorderLayout(nil, nil, nil, nil))

	scanWindow := app.NewWindow("Scan new item")
	dialog := &AddDialog{
		imageContainer: container,
		mainWindow:     w,
		app:            app,
		scanWindow:     scanWindow,
	}
	return dialog
}

func (d *AddDialog) setImage(reader fyne.URIReadCloser, err error) {

	if err != nil {
		dialog.ShowError(err, *d.mainWindow)
		return
	}
	if reader == nil {
		log.Println("Cancelled")
		return
	}
	defer reader.Close()

	image := canvas.NewImageFromReader(reader, "test")
	imageWidget := NewImageWidget(image)
	d.imageContainer.Add(imageWidget)
}

func (d *AddDialog) createFileDialogButton() *widget.Button {
	return widget.NewButton("Select image file (jpg or png)", func() {
		fd := dialog.NewFileOpen(d.setImage, d.scanWindow)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
		fd.Show()
	})
}

func (d *AddDialog) ShowDialog() {

	selectImageButton := d.createFileDialogButton()

	content := container.New(layout.NewBorderLayout(selectImageButton, nil, nil, nil), selectImageButton, d.imageContainer)

	content.Resize(fyne.NewSize(800, 800))

	d.scanWindow.SetContent(content)
	d.scanWindow.Resize(fyne.NewSize(800, 800))
	d.scanWindow.Show()
}
