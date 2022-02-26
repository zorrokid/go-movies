package ui

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/disintegration/imaging"
	"github.com/otiai10/gosseract"
	"github.com/zorrokid/go-movies/scanner"
	"github.com/zorrokid/go-movies/scanner-ui/util"
)

type SelectMode int8

const (
	Word       SelectMode = 0
	FourPoints SelectMode = 1
)

const defaultSelectMode = Word

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
			w/r.widget.scale,
			h/r.widget.scale,
		)

		return minSize
	}
	return fyne.NewSize(300, 300)
}

func (r *imageWidgetRenderer) Refresh() {
	fmt.Println("Refresh")

	if r.widget.Image == nil {
		return
	}

	objects := []fyne.CanvasObject{}

	w := float32(r.widget.Image.Image.Bounds().Dx())
	h := float32(r.widget.Image.Image.Bounds().Dy())

	imageContainer := container.NewWithoutLayout(r.widget.Image)
	r.widget.Image.Move(fyne.NewPos(0, 0))
	s := fyne.NewSize(w/r.widget.scale, h/r.widget.scale)
	r.widget.Image.Resize(s)
	objects = append(objects, imageContainer)

	for _, bb := range r.widget.Boxes {
		// draw boxes by selected confidence
		if bb.Confidence < r.widget.confidence {
			continue
		}
		rect := bb.Box
		marker := canvas.NewRectangle(color.RGBA{
			R: 1,
			G: 50,
			B: 12,
			A: 50,
		})
		s := fyne.NewSize(
			float32(rect.Max.X-rect.Min.X)/r.widget.scale,
			float32(rect.Max.Y-rect.Min.Y)/r.widget.scale,
		)
		marker.Resize(s)
		pos := fyne.NewPos(
			float32(rect.Min.X)/r.widget.scale,
			float32(rect.Min.Y)/r.widget.scale,
		)
		marker.Move(pos)
		objects = append(objects, marker)
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
	Image              *canvas.Image
	Boxes              []gosseract.BoundingBox
	selected           func(words []string)
	tap                *fyne.Position
	scale              float32
	confidence         float64
	selectMode         SelectMode
	positionSelections []fyne.Position
}

func NewImageWidget(selected func(word []string)) *ImageWidget {
	i := &ImageWidget{
		selected:           selected,
		scale:              1.0,
		Boxes:              []gosseract.BoundingBox{},
		selectMode:         Word,
		positionSelections: []fyne.Position{},
	}
	i.ExtendBaseWidget(i)
	return i
}

func (i *ImageWidget) SetImage(image image.Image) {
	i.Image = canvas.NewImageFromImage(image)
	i.Rescan()
	i.Refresh()
}

func (i *ImageWidget) setBoxes(boxes []gosseract.BoundingBox) {
	fmt.Println("SetBoxes")
	i.Boxes = boxes
	i.Refresh()
}

func (i *ImageWidget) tappedWords(event *fyne.PointEvent) {

	positionX := event.Position.X * i.scale
	positionY := event.Position.Y * i.scale

	for _, b := range i.Boxes {
		if positionX > float32(b.Box.Min.X) &&
			positionX < float32(b.Box.Max.X) &&
			positionY > float32(b.Box.Min.Y) &&
			positionY < float32(b.Box.Max.Y) {
			words := make([]string, 1)
			words[0] = b.Word
			i.selected(words)
		}
	}
}

func (i *ImageWidget) tappedPoints(event *fyne.PointEvent) {

	scaledPosition := fyne.NewPos(event.Position.X*i.scale, event.Position.Y*i.scale)
	i.positionSelections = append(i.positionSelections, scaledPosition)

	if len(i.positionSelections) == 4 {
		for _, pos := range i.positionSelections {
			fmt.Printf("Position (%f, %f) selected.\n", pos.X, pos.Y)
		}
		imgTr, err := util.Transform(i.Image.Image, i.positionSelections)
		if err != nil {
			log.Fatal(err)
		}
		i.SetImage(imgTr)
		i.positionSelections = []fyne.Position{}
		i.selectMode = defaultSelectMode

	}
}

func (i *ImageWidget) Tapped(event *fyne.PointEvent) {

	switch i.selectMode {
	case FourPoints:
		i.tappedPoints(event)
	default:
		i.tappedWords(event)
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
	i.SetImage(imgRt)
	i.Refresh()
}

func (i *ImageWidget) RotateRight() {
	fmt.Println("rotate right")
	imgRt := imaging.Rotate270(i.Image.Image)
	i.Image.Image = imgRt

	i.SetImage(imgRt)
	i.Refresh()
}

func (i *ImageWidget) Rescan() {
	fmt.Println("rescan")
	imageBytes := util.ImageToByteBuffer(i.Image.Image)
	if bbs, err := scanner.ScanFromBytes(imageBytes, "fin"); err != nil {
		log.Fatal(err)
	} else {
		i.setBoxes(bbs)
	}
}

func (d *ImageWidget) setScale(scale string) {
	fmt.Printf("Scale %s\n", scale)
	if sc, err := strconv.Atoi(scale); err != nil {
		log.Fatal(err)
	} else {
		d.scale = float32(sc)
	}
	d.Refresh()
}

func (d *ImageWidget) setConfidence(confidence string) {
	fmt.Printf("Confidence %s\n", confidence)
	if c, err := strconv.Atoi(confidence); err != nil {
		log.Fatal(err)
	} else {
		d.confidence = float64(c)
	}
	d.Refresh()
}

func (d *ImageWidget) SetSelectPoints() {
	d.SetSelect(FourPoints)
}

func (d *ImageWidget) SetSelect(mode SelectMode) {
	d.selectMode = mode
}

func (i *ImageWidget) getWordsBetween(posA fyne.Position, posB fyne.Position) []string {

	minX := math.Min(float64(posA.X), float64(posB.X)) * float64(i.scale)
	minY := math.Min(float64(posA.Y), float64(posB.Y)) * float64(i.scale)

	maxX := math.Max(float64(posA.X), float64(posB.X)) * float64(i.scale)
	maxY := math.Max(float64(posA.Y), float64(posB.Y)) * float64(i.scale)

	words := make([]string, 5)
	for _, b := range i.Boxes {
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
