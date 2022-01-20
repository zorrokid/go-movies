package scanner

import (
	"github.com/otiai10/gosseract"
)

type PageDiv struct {
}

func Scan(path string, languages ...string) ([]gosseract.BoundingBox, error) {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetLanguage(languages...)
	client.SetImage(path)
	bb, _ := client.GetBoundingBoxes(gosseract.RIL_WORD)
	return bb, nil
}
