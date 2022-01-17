package scanner

import (
	"github.com/otiai10/gosseract"
)

func Scan(path string) string {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetLanguage("fin")
	client.SetImage(path)
	hocrtext, _ := client.Text() // client.HOCRText()
	return hocrtext
}
