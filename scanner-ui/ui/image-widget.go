package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type imageWidgetRenderer struct {
	fyne.WidgetRenderer
	//	image *canvas.Image
	//	label *widget.Label
	// widget  *ImageWidget
	objects   []fyne.CanvasObject
	container *fyne.Container
}

func (r *imageWidgetRenderer) Layout(s fyne.Size) {
	fmt.Println("Layout")
}

func (r *imageWidgetRenderer) MinSize() fyne.Size {
	return r.container.MinSize()
}

func (r *imageWidgetRenderer) Refresh() {
	fmt.Println("Refresh")
	r.container.Refresh()
	//canvas.Refresh(r.container)
}

func (r *imageWidgetRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *imageWidgetRenderer) Destroy() {
	fmt.Println("Destroy")
}

type ImageWidget struct {
	widget.BaseWidget
	reader   fyne.URIReadCloser
	OnTapped func(event *fyne.PointEvent)
}

func NewImageWidget(reader fyne.URIReadCloser) *ImageWidget {
	i := &ImageWidget{
		reader:   reader,
		OnTapped: onTapped,
	}
	i.ExtendBaseWidget(i)
	return i
}

func onTapped(event *fyne.PointEvent) {
	fmt.Printf("Tapped %f %f", event.Position.X, event.Position.Y)
}

func (i *ImageWidget) CreateRenderer() fyne.WidgetRenderer {
	i.ExtendBaseWidget(i)

	//image := loadImage(i.reader)
	image := canvas.NewImageFromReader(i.reader, "test")
	image.Resize(fyne.NewSize(600, 600))
	image.FillMode = canvas.ImageFillOriginal
	// label := widget.NewLabel("Hello")
	// label2 := widget.NewLabel("Hello")
	container := container.NewHBox(image)
	//container.Add(label2)
	container.Resize(fyne.NewSize(500, 500))
	objects := []fyne.CanvasObject{container}
	r := &imageWidgetRenderer{
		//label:     label,
		//image:     image,
		objects:   objects,
		container: container,
	}
	return r
}
