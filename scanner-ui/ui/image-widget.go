package ui

import (
	"fmt"
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/otiai10/gosseract"
)

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
	selected    func(word string)
}

func NewImageWidget(image *canvas.Image, boxes []gosseract.BoundingBox, imageConfig image.Config, selected func(word string)) *ImageWidget {
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

	positionX := event.Position.X * 4
	positionY := event.Position.Y * 4

	for _, b := range i.boxes {
		if positionX > float32(b.Box.Min.X) && positionX < float32(b.Box.Max.X) && positionY > float32(b.Box.Min.Y) && positionY < float32(b.Box.Max.Y) {
			i.selected(b.Word)
		}
	}
}

func (i *ImageWidget) TappedSecondary(event *fyne.PointEvent) {
	fmt.Printf("Tapped %f %f\n", event.Position.X, event.Position.Y)
}

func (i *ImageWidget) CreateRenderer() fyne.WidgetRenderer {
	i.ExtendBaseWidget(i)

	imageContainer := container.NewWithoutLayout(i.image)
	i.image.Move(fyne.NewPos(0, 0))
	i.image.Resize(fyne.NewSize(float32(i.imageConfig.Width)/4, float32(i.imageConfig.Height)/4))
	objects := []fyne.CanvasObject{imageContainer}

	for _, bb := range i.boxes {
		rect := bb.Box
		marker := canvas.NewRectangle(color.RGBA{
			R: 1,
			G: 50,
			B: 12,
			A: 50,
		})
		marker.Resize(fyne.NewSize(float32(rect.Max.X-rect.Min.X)/4, float32(rect.Max.Y-rect.Min.Y)/4))
		marker.Move(fyne.NewPos(float32(rect.Min.X)/4, float32(rect.Min.Y)/4))
		objects = append(objects, marker)
	}

	r := &imageWidgetRenderer{
		objects: objects,
		boxes:   i.boxes,
	}
	return r
}
