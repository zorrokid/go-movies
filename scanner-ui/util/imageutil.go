package util

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"gocv.io/x/gocv"
)

func ImageToByteBuffer(image image.Image) []byte {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, image, nil)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

func ReadImage(filePath string) (image.Image, error) {

	reader, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func DrawBoxes(img image.Image, rects []image.Rectangle) image.Image {

	m := image.NewRGBA(img.Bounds())
	draw.Draw(m, m.Bounds(), img, image.Point{0, 0}, draw.Src)

	blue := color.RGBA{0, 0, 255, 80}

	for _, r := range rects {
		m1 := image.NewNRGBA(r)
		draw.Draw(m1, m1.Bounds(), &image.Uniform{blue}, image.Point{0, 0}, draw.Src)
		draw.Draw(m, m1.Bounds(), m1, m1.Bounds().Min, draw.Over)
	}

	return m
}

func GetMinMaxXY(positions []fyne.Position) (minX float32, maxX float32, minY float32, maxY float32) {
	minX = positions[0].X
	maxX = positions[0].X
	minY = positions[0].Y
	maxY = positions[0].Y
	for _, pos := range positions {
		if pos.X < minX {
			minX = pos.X
		}
		if pos.X > maxX {
			maxX = pos.X
		}
		if pos.Y < minY {
			minY = pos.Y
		}
		if pos.Y > maxY {
			maxY = pos.Y
		}
	}
	return minX, maxX, minY, maxY
}

func Transform(img image.Image, positions []fyne.Position) (image.Image, error) {
	imageBuffer := ImageToByteBuffer(img)
	img2, err := gocv.IMDecode(imageBuffer, gocv.IMReadUnchanged)

	if len(positions) < 4 {
		return nil, fmt.Errorf("not enough positions for transform")
	}

	if err != nil {
		return nil, err
	}

	srcVect := gocv.NewPointVector()
	srcVect.Append(image.Point{int(positions[0].X), int(positions[0].Y)})
	srcVect.Append(image.Point{int(positions[1].X), int(positions[1].Y)})
	srcVect.Append(image.Point{int(positions[2].X), int(positions[2].Y)})
	srcVect.Append(image.Point{int(positions[3].X), int(positions[3].Y)})

	minX, maxX, minY, maxY := GetMinMaxXY(positions)
	fmt.Printf("minX: %f, maxX: %f, minY: %f, maxY: %f", minX, maxX, minY, maxY)
	destVect := gocv.NewPointVector()
	destVect.Append(image.Point{int(minX), int(minY)})
	destVect.Append(image.Point{int(maxX), int(minY)})
	destVect.Append(image.Point{int(maxX), int(maxY)})
	destVect.Append(image.Point{int(minX), int(maxY)})

	trmat := gocv.GetPerspectiveTransform(srcVect, destVect)
	gocv.WarpPerspective(img2, &img2, trmat, image.Point{X: img.Bounds().Dx(), Y: img.Bounds().Dy()})

	bufOut, err := gocv.IMEncode(".jpg", img2)
	if err != nil {
		return nil, err
	}
	defer bufOut.Close()
	bytesOut := bufOut.GetBytes()
	reader := bytes.NewReader(bytesOut)
	imgOut, _, err := image.Decode(reader)

	if err != nil {
		return nil, err
	}

	return imgOut, nil
}
