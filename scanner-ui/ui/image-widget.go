package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type imageWidgetRenderer struct {
	image *canvas.Image
	// widget  *ImageWidget
	objects []fyne.CanvasObject
}

func (r *imageWidgetRenderer) Layout(s fyne.Size) {

}

func (r *imageWidgetRenderer) MinSize() fyne.Size {
	return fyne.Size{}
}

func (r *imageWidgetRenderer) Refresh() {

}

func (r *imageWidgetRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *imageWidgetRenderer) Destroy() {

}

type ImageWidget struct {
	widget.BaseWidget
	ImagePath string
	OnTapped  func(event *fyne.PointEvent)
}

func NewImageWidget(imagePath string) *ImageWidget {
	i := &ImageWidget{
		ImagePath: imagePath,
		OnTapped:  onTapped,
	}
	i.ExtendBaseWidget(i)
	return i
}

func onTapped(event *fyne.PointEvent) {
	fmt.Printf("Tapped %f %f", event.Position.X, event.Position.Y)
}

func (i *ImageWidget) CreateRenderer() fyne.WidgetRenderer {
	i.ExtendBaseWidget(i)

	image := canvas.NewImageFromFile(i.ImagePath)
	objects := []fyne.CanvasObject{image}
	r := &imageWidgetRenderer{
		image:   image,
		objects: objects,
	}
	return r
}
