package ui

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/zorrokid/go-movies/scanner"
)

type AddDialog struct {
	imageContainer *fyne.Container
	mainWindow     *fyne.Window
	app            fyne.App
}

func NewAddDialog(w *fyne.Window, app fyne.App) *AddDialog {
	text1 := canvas.NewText("Add image", color.White)
	container := container.New(layout.NewBorderLayout(text1, nil, nil, nil), text1)
	dialog := &AddDialog{
		imageContainer: container,
		mainWindow:     w,
		app:            app,
	}
	return dialog
}

func (d *AddDialog) update() {
	d.imageContainer.Refresh()
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

	image := d.loadImage(reader)
	image.Resize(fyne.NewSize(500, 500))
	d.imageContainer.Add(image)
	d.update()

	filePath := reader.URI().Path()

	if bbs, err := scanner.Scan(filePath, "fin"); err != nil {
		log.Fatal(err)
	} else {
		for _, bb := range bbs {
			fmt.Printf("%s, %d, %d\n", bb.Word, bb.Box.Dx(), bb.Box.Dy())
		}
	}

}

func (d *AddDialog) createFileDialogButton(w fyne.Window) *widget.Button {
	return widget.NewButton("Select image file (jpg or png)", func() {
		fd := dialog.NewFileOpen(d.setImage, w)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
		fd.Show()
	})
}

func (d *AddDialog) ShowDialog() {
	w := d.app.NewWindow("Scan new item")

	selectImageButton := d.createFileDialogButton(w)
	content := container.New(layout.NewBorderLayout(selectImageButton, nil, nil, nil), selectImageButton, d.imageContainer)

	content.Resize(fyne.NewSize(800, 800))

	w.SetContent(content)
	w.Resize(fyne.NewSize(800, 800))
	w.Show()
}

func (d *AddDialog) loadImage(f fyne.URIReadCloser) *canvas.Image {
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fyne.LogError("Failed to load image data", err)
		return nil
	}
	res := fyne.NewStaticResource(f.URI().Name(), data)

	return canvas.NewImageFromResource(res)
}
