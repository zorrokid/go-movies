package ui

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/theme"
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
	text := widget.NewMultiLineEntry()
	text.MultiLine = true
	text.Wrapping = fyne.TextWrapBreak
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

		img, err := util.ReadImage(filePath)
		if err != nil {
			log.Fatal(err)
			return
		}

		image := canvas.NewImageFromImage(img)
		image.FillMode = canvas.ImageFillContain
		imageWidget := NewImageWidget(image, bbs, imgConfig, d.selected)
		imageWidget.Resize(fyne.NewSize(800, 500))
		d.imageContainer.Add(imageWidget)
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
	content := container.New(layout.NewBorderLayout(nil, nil, nil, nil))
	d.imageContainer = content
	d.imageContainer.Resize(fyne.NewSize(1800, 1500))
	scroll := container.NewScroll(d.imageContainer)
	scroll.Resize(fyne.NewSize(1800, 1500))

	fieldsForm := container.New(layout.NewGridLayoutWithColumns(3))
	label := widget.NewLabel("Title")
	clearBtn := widget.NewButton("Clear", d.clearText)
	fieldsForm.Add(label)
	fieldsForm.Add(d.text)
	fieldsForm.Add(clearBtn)

	grid := container.NewAdaptiveGrid(3, d.createFileList(), scroll, fieldsForm)

	d.scanWindow.SetContent(grid)
	d.scanWindow.Show()
}

func (d *AddDialog) createFileList() fyne.CanvasObject {

	// icon := widget.NewIcon(nil)
	// label := widget.NewLabel("Select An Item From The List")
	// hbox := container.NewHBox(icon, label)

	data := make([]string, 1000)
	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[id])
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		// label.SetText(data[id])
		// icon.SetResource(theme.DocumentIcon())
	}
	list.OnUnselected = func(id widget.ListItemID) {
		// label.SetText("Select An Item From The List")
		// icon.SetResource(nil)
	}

	selectImageButton := d.createFileDialogButton()
	selectImageButton.Resize(fyne.NewSize(100, 100))
	content := container.New(layout.NewBorderLayout(selectImageButton, nil, nil, nil),
		selectImageButton)
	list.Resize(fyne.NewSize(100, 100))
	content.Add(list)

	return content
}

func (d *AddDialog) clearText() {
	d.text.Text = ""
	d.text.Refresh()
}

func (d *AddDialog) selected(word []string) {
	if len(d.text.Text) > 0 {
		d.text.Text += " "
	}
	d.text.Text += strings.Join(word, " ")
	d.text.Refresh()
	fmt.Printf("Word %s selected\n", word)
}
