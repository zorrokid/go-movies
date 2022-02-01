package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type imageWidgetRenderer struct {
	fyne.WidgetRenderer
	//widget  *ImageWidget
	objects []fyne.CanvasObject
}

func (r *imageWidgetRenderer) Layout(s fyne.Size) {
	r.objects[0].Resize(s)
}

func (r *imageWidgetRenderer) MinSize() fyne.Size {
	s := fyne.NewSize(500, 500)
	return s
}

func (i *imageWidgetRenderer) Refresh() {
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
	image *canvas.Image
}

func NewImageWidget(image *canvas.Image) *ImageWidget {
	i := &ImageWidget{
		image: image,
	}
	i.ExtendBaseWidget(i)
	return i
}

func (i *ImageWidget) Tapped(event *fyne.PointEvent) {
	fmt.Printf("Tapped %f %f", event.Position.X, event.Position.Y)
}

func (i *ImageWidget) TappedSecondary(event *fyne.PointEvent) {
	fmt.Printf("Tapped %f %f", event.Position.X, event.Position.Y)
}

func (i *ImageWidget) CreateRenderer() fyne.WidgetRenderer {
	i.ExtendBaseWidget(i)
	objects := []fyne.CanvasObject{i.image}
	r := &imageWidgetRenderer{
		objects: objects,
	}
	return r
}
