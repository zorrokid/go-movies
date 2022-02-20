package ui

import (
	"bytes"
	"fmt"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/disintegration/imaging"
	"github.com/otiai10/gosseract"
	"github.com/zorrokid/go-movies/scanner"
)

const Scale = 4

type imageWidgetRenderer struct {
	fyne.WidgetRenderer
	objects *[]fyne.CanvasObject
	widget  *ImageWidget
}

func (r *imageWidgetRenderer) Layout(s fyne.Size) {
	fmt.Println("Layout")
	// (*r.objects)[0].Move(fyne.NewPos(0, 0))
	// (*r.objects)[0].Resize(s)
}

func (r *imageWidgetRenderer) MinSize() fyne.Size {
	if r.widget.Image != nil {
		w := float32(r.widget.Image.Image.Bounds().Dx())
		h := float32(r.widget.Image.Image.Bounds().Dy())
		minSize := fyne.NewSize(
			w/Scale,
			h/Scale,
		)

		return minSize
	}
	return fyne.NewSize(300, 300)
}

func (r *imageWidgetRenderer) Refresh() {
	fmt.Println("Refresh")

	objects := []fyne.CanvasObject{}

	w := float32(r.widget.Image.Image.Bounds().Dx())
	h := float32(r.widget.Image.Image.Bounds().Dy())

	if r.widget.Image != nil {
		imageContainer := container.NewWithoutLayout(r.widget.Image)
		r.widget.Image.Move(fyne.NewPos(0, 0))
		s := fyne.NewSize(w/Scale, h/Scale)
		r.widget.Image.Resize(s)
		objects = append(objects, imageContainer)
	}

	if r.widget.Boxes != nil {
		for _, bb := range *r.widget.Boxes {
			rect := bb.Box
			marker := canvas.NewRectangle(color.RGBA{
				R: 1,
				G: 50,
				B: 12,
				A: 50,
			})
			s := fyne.NewSize(float32(rect.Max.X-rect.Min.X)/Scale, float32(rect.Max.Y-rect.Min.Y)/Scale)
			marker.Resize(s)
			pos := fyne.NewPos(float32(rect.Min.X)/Scale, float32(rect.Min.Y)/Scale)
			marker.Move(pos)
			objects = append(objects, marker)
		}
	}

	r.objects = &objects
}

func (r *imageWidgetRenderer) Objects() []fyne.CanvasObject {
	return *r.objects
}

func (r *imageWidgetRenderer) Destroy() {
	fmt.Println("Destroy")
}

type ImageWidget struct {
	widget.BaseWidget
	Image    *canvas.Image
	Boxes    *[]gosseract.BoundingBox
	selected func(words []string)
	tap      *fyne.Position
}

func NewImageWidget(selected func(word []string)) *ImageWidget {
	i := &ImageWidget{
		selected: selected,
	}
	i.ExtendBaseWidget(i)
	return i
}

func (i *ImageWidget) SetImage(image *canvas.Image) {
	i.Image = image
	i.Refresh()
}

func (i *ImageWidget) SetBoxes(boxes *[]gosseract.BoundingBox) {
	i.Boxes = boxes
	i.Refresh()
}

func (i *ImageWidget) Tapped(event *fyne.PointEvent) {

	positionX := event.Position.X * Scale
	positionY := event.Position.Y * Scale

	for _, b := range *i.Boxes {
		if positionX > float32(b.Box.Min.X) && positionX < float32(b.Box.Max.X) && positionY > float32(b.Box.Min.Y) && positionY < float32(b.Box.Max.Y) {
			words := make([]string, 1)
			words[0] = b.Word
			i.selected(words)
		}
	}
}

func (i *ImageWidget) TappedSecondary(event *fyne.PointEvent) {
	if i.tap == nil {
		i.tap = &event.Position
		return
	}

	words := i.getWordsBetween(*i.tap, event.Position)
	if len(words) > 0 {
		i.selected(words)
	}
	i.tap = nil
}

func (i *ImageWidget) Rotate() {
	fmt.Println("rotate")
	imgRt := imaging.Rotate90(i.Image.Image)
	i.Image.Image = imgRt

	i.SetImage(canvas.NewImageFromImage(imgRt))
	i.Refresh()
}

func (i *ImageWidget) Rescan() {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, i.Image.Image, nil)
	if err != nil {
		log.Fatal(err)
	}

	imageBytes := buf.Bytes()
	if bbs, err := scanner.ScanFromBytes(imageBytes, "fin"); err != nil {
		log.Fatal(err)
	} else {
		i.SetBoxes(&bbs)
	}
}

func (i *ImageWidget) getWordsBetween(posA fyne.Position, posB fyne.Position) []string {

	minX := math.Min(float64(posA.X), float64(posB.X)) * Scale
	minY := math.Min(float64(posA.Y), float64(posB.Y)) * Scale

	maxX := math.Max(float64(posA.X), float64(posB.X)) * Scale
	maxY := math.Max(float64(posA.Y), float64(posB.Y)) * Scale

	words := make([]string, 5)
	for _, b := range *i.Boxes {
		bMinX := float64(b.Box.Min.X)
		bMaxX := float64(b.Box.Max.X)
		bMinY := float64(b.Box.Min.Y)
		bMaxY := float64(b.Box.Max.Y)
		if bMinX > minX && bMaxX < maxX && bMinY > minY && bMaxY < maxY {
			w := strings.TrimSpace(b.Word)
			if len(w) > 0 {
				words = append(words, w)
			}
		}
	}
	return words
}

func (i *ImageWidget) CreateRenderer() fyne.WidgetRenderer {
	i.ExtendBaseWidget(i)

	objects := []fyne.CanvasObject{}

	r := &imageWidgetRenderer{
		objects: &objects,
		widget:  i,
	}
	return r
}
