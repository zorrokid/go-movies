package ui

import (
	"errors"
	"fmt"
	_ "image/jpeg"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/zorrokid/go-movies/scanner-ui/util"
)

type AddDialog struct {
	imageContainer *fyne.Container
	scanWindow     fyne.Window
	text           *widget.Entry
	fileListData   []fyne.URI
	imageWidget    *ImageWidget
	progressBar    *widget.ProgressBarInfinite
}

func NewAddDialog(app *fyne.App) *AddDialog {
	text := widget.NewMultiLineEntry()
	text.MultiLine = true
	text.Wrapping = fyne.TextWrapBreak
	scanWindow := (*app).NewWindow("Scan new item")
	progressBar := widget.NewProgressBarInfinite()
	progressBar.Stop()
	dialog := &AddDialog{
		scanWindow:   scanWindow,
		text:         text,
		progressBar:  progressBar,
		fileListData: []fyne.URI{},
	}
	return dialog
}

func (d *AddDialog) readFiles(lu fyne.ListableURI, err error) {

	if err != nil {
		dialog.ShowError(err, d.scanWindow)
		return
	}

	if lu == nil {
		return
	}

	uriList, err := lu.List()
	if err != nil {
		log.Fatal(err)
		return
	}

	d.fileListData = util.FilterByExtension(uriList, ".jpg", ".jpeg")
}

func (d *AddDialog) setSelectedImage(i int) {
	if i >= len(d.fileListData) {
		dialog.ShowError(errors.New("selected index out of bounds"), d.scanWindow)
	}
	selectedImageURI := d.fileListData[i]
	d.setImage(selectedImageURI)
}

func (d *AddDialog) setImage(uri fyne.URI) {
	d.progressBar.Start()
	filePath := uri.Path()

	img, err := util.ReadImage(filePath)
	if err != nil {
		log.Fatal(err)
		return
	}

	imgFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer imgFile.Close()

	d.imageWidget.SetImage(img)
	d.progressBar.Stop()
}

func (d *AddDialog) createFileDialogButton() *widget.Button {
	return widget.NewButton("Select image file (jpg or png)", func() {
		fd := dialog.NewFolderOpen(d.readFiles, d.scanWindow)
		fd.Show()
	})
}

func (d *AddDialog) ShowDialog() {

	d.imageWidget = NewImageWidget(d.selected)

	btnRotateLeft := widget.NewButton("RoL", d.imageWidget.Rotate)
	btnRotateRight := widget.NewButton("RoR", d.imageWidget.RotateRight)

	btnRescan := widget.NewButton("Re", d.imageWidget.Rescan)
	selectScale := widget.NewSelect([]string{"1", "2", "4", "8"}, d.imageWidget.setScale)
	selectConfidence := widget.NewSelect([]string{"20", "40", "60", "80"}, d.imageWidget.setConfidence)

	btnSelectPoints := widget.NewButton("S", d.imageWidget.SetSelectPoints)

	imageActionsContainer := container.New(
		layout.NewHBoxLayout(),
		btnRotateLeft,
		btnRotateRight,
		btnRescan,
		selectScale,
		selectConfidence,
		btnSelectPoints,
	)

	content := container.New(layout.NewBorderLayout(imageActionsContainer, d.progressBar, nil, nil), imageActionsContainer, d.progressBar)
	content.Add(d.imageWidget)
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

	grid := container.NewAdaptiveGrid(3, d.createFileListContainer(), scroll, fieldsForm)

	d.scanWindow.SetContent(grid)
	d.scanWindow.Resize(fyne.NewSize(800, 800))
	d.scanWindow.Show()
}

func (d *AddDialog) createFileList() *widget.List {

	list := widget.NewList(
		// length
		func() int {
			return len(d.fileListData)
		},
		// create item
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
		},
		// update item
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(d.fileListData[id].Name())
		},
	)

	list.OnSelected = d.setSelectedImage
	list.OnUnselected = func(id widget.ListItemID) {
		fmt.Println("OnUnselected")
	}

	return list
}

func (d *AddDialog) createFileListContainer() fyne.CanvasObject {

	selectImageButton := d.createFileDialogButton()
	selectImageButton.Resize(fyne.NewSize(100, 100))
	content := container.New(layout.NewBorderLayout(selectImageButton, nil, nil, nil),
		selectImageButton)

	fileList := d.createFileList()
	content.Add(fileList)
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
