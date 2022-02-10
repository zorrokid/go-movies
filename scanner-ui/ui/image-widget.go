package ui

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/otiai10/gosseract"
)

const Scale = 4

type imageWidgetRenderer struct {
	fyne.WidgetRenderer
	objects []fyne.CanvasObject
	boxes   []gosseract.BoundingBox
}

func (r *imageWidgetRenderer) Layout(s fyne.Size) {
	fmt.Println("Layout")
	r.objects[0].Move(fyne.NewPos(0, 0))
	r.objects[0].Resize(s)
}

func (r *imageWidgetRenderer) MinSize() fyne.Size {
	return r.objects[0].MinSize()
}

func (i *imageWidgetRenderer) Refresh() {
	fmt.Println("Refresh")
	i.objects[0].Move(fyne.NewPos(0, 0))
	i.objects[0].Refresh()
}

func (r *imageWidgetRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *imageWidgetRenderer) Destroy() {
	fmt.Println("Destroy")
}

type ImageWidget struct {
	widget.BaseWidget
	image       *canvas.Image
	boxes       []gosseract.BoundingBox
	imageConfig image.Config
	selected    func(words []string)
	tap         *fyne.Position
}

func NewImageWidget(image *canvas.Image, boxes []gosseract.BoundingBox, imageConfig image.Config, selected func(word []string)) *ImageWidget {
	i := &ImageWidget{
		image:       image,
		boxes:       boxes,
		imageConfig: imageConfig,
		selected:    selected,
	}
	i.ExtendBaseWidget(i)
	return i
}

func (i *ImageWidget) Tapped(event *fyne.PointEvent) {

	positionX := event.Position.X * Scale
	positionY := event.Position.Y * Scale

	for _, b := range i.boxes {
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

func (i *ImageWidget) getWordsBetween(posA fyne.Position, posB fyne.Position) []string {

	minX := math.Min(float64(posA.X), float64(posB.X)) * Scale
	minY := math.Min(float64(posA.Y), float64(posB.Y)) * Scale

	maxX := math.Max(float64(posA.X), float64(posB.X)) * Scale
	maxY := math.Max(float64(posA.Y), float64(posB.Y)) * Scale

	words := make([]string, 5)
	for _, b := range i.boxes {
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

	imageContainer := container.NewWithoutLayout(i.image)
	i.image.Move(fyne.NewPos(0, 0))
	s := fyne.NewSize(float32(i.imageConfig.Width)/Scale, float32(i.imageConfig.Height)/Scale)
	i.image.Resize(s)
	objects := []fyne.CanvasObject{imageContainer}

	for _, bb := range i.boxes {
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

	r := &imageWidgetRenderer{
		objects: objects,
		boxes:   i.boxes,
	}
	return r
}
