package ui

import (
	"fmt"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/otiai10/gosseract"
)

type imageWidgetRenderer struct {
	fyne.WidgetRenderer
	//widget  *ImageWidget
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
}

func NewImageWidget(image *canvas.Image, boxes []gosseract.BoundingBox, imageConfig image.Config) *ImageWidget {
	i := &ImageWidget{
		image:       image,
		boxes:       boxes,
		imageConfig: imageConfig,
	}
	i.ExtendBaseWidget(i)
	return i
}

func (i *ImageWidget) Tapped(event *fyne.PointEvent) {

	positionX := event.Position.X * 4
	positionY := event.Position.Y * 4

	fmt.Printf("Tapped position (%f, %f)\n", positionX, positionY)

	for _, b := range i.boxes {
		if positionX > float32(b.Box.Min.X) && positionX < float32(b.Box.Max.X) && positionY > float32(b.Box.Min.Y) && positionY < float32(b.Box.Max.Y) {
			fmt.Println(b.Word)
		}
	}
}

func (i *ImageWidget) TappedSecondary(event *fyne.PointEvent) {
	fmt.Printf("Tapped %f %f\n", event.Position.X, event.Position.Y)
}

func (i *ImageWidget) CreateRenderer() fyne.WidgetRenderer {
	i.ExtendBaseWidget(i)

	// rect := canvas.NewRectangle(color.White)
	// rect.Resize(fyne.NewSize(10, 10))
	// rect.Move(fyne.NewPos(100, 100))

	imageContainer := container.NewWithoutLayout(i.image)
	i.image.Move(fyne.NewPos(0, 0))
	i.image.Resize(fyne.NewSize(float32(i.imageConfig.Width)/4, float32(i.imageConfig.Height)/4))
	objects := []fyne.CanvasObject{imageContainer}
	r := &imageWidgetRenderer{
		objects: objects,
		boxes:   i.boxes,
	}
	return r
}
