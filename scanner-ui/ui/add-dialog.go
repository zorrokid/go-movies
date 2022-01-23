package ui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/zorrokid/go-movies/scanner"
)

type AddDialog struct {
}

func NewAddDialog() *AddDialog {
	dialog := &AddDialog{}
	return dialog
}

func (d *AddDialog) ShowDialog(w *fyne.Window) {
	selectImageButton := widget.NewButton("Select image file (jpg or png)", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, *w)
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
				for _, bb := range bbs {
					fmt.Printf("%s, %d, %d\n", bb.Word, bb.Box.Dx(), bb.Box.Dy())
				}
			}

		}, *w)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
		fd.Show()
	})

	items := []*widget.FormItem{
		widget.NewFormItem("Select image", selectImageButton),
	}

	dialog.ShowForm("Scan new item", "Scan", "Cancel", items, func(b bool) {
		if !b {
			return
		}
	}, *w)
}
