package ui

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/zorrokid/go-movies/scanner"
	"github.com/zorrokid/go-movies/scanner-ui/util"
)

type AddDialog struct {
	imageContainer *fyne.Container
	mainWindow     *fyne.Window
	scanWindow     fyne.Window
	app            fyne.App
	text           *widget.Entry
}

func NewAddDialog(w *fyne.Window, app fyne.App) *AddDialog {
	text := widget.NewEntry()
	scanWindow := app.NewWindow("Scan new item")
	dialog := &AddDialog{
		mainWindow: w,
		app:        app,
		scanWindow: scanWindow,
		text:       text,
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

	filePath := reader.URI().Path()

	if bbs, err := scanner.Scan(filePath, "fin"); err != nil {
		log.Fatal(err)
	} else {

		imgFile, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer imgFile.Close()

		imgConfig, imgType, err := image.DecodeConfig(imgFile)

		fmt.Printf("image type %s\n", imgType)
		fmt.Printf("image %d x %d\n", imgConfig.Width, imgConfig.Height)

		if err != nil {
			log.Fatal(err)
			return
		}

		// rects := make([]image.Rectangle, len(bbs))
		// for _, r := range bbs {
		// 	rects = append(rects, r.Box)
		// }

		img, err := util.ReadImage(filePath)
		if err != nil {
			log.Fatal(err)
			return
		}
		//img2 := util.DrawBoxes(img, rects)

		image := canvas.NewImageFromImage(img)
		//image := canvas.NewImageFromReader(reader, "test")
		image.FillMode = canvas.ImageFillContain
		imageWidget := NewImageWidget(image, bbs, imgConfig, d.selected)
		d.imageContainer.Add(imageWidget)
		d.imageContainer.Add(widget.NewLabel("Test"))
	}

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

	content := container.New(layout.NewBorderLayout(selectImageButton, nil, nil, nil), selectImageButton)
	d.imageContainer = content

	// fieldsForm := container.New(layout.NewFormLayout())
	// label := widget.NewLabel("Title")
	// fieldsForm.Add(label)
	// fieldsForm.Add(d.text)

	// addForm := container.New(layout.NewHBoxLayout())

	// content := container.New(layout.NewVBoxLayout())
	// content.Add(selectImageButton)
	// d.imageContainer = content

	// //addForm.Add(d.imageContainer)
	// addForm.Add(fieldsForm)

	d.scanWindow.SetContent(d.imageContainer)
	d.scanWindow.Show()
}

func (d *AddDialog) selected(word string) {
	fmt.Printf("Word %s selected\n", word)
}
