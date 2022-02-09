package util

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

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
