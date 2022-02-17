package util

import (
	"testing"

	"fyne.io/fyne/v2"
)

func TestFilterByExtensionEmptyUriList(t *testing.T) {
	uriList := []fyne.URI{}
	extension := ".jpg"

	result := FilterByExtension(uriList, extension)

	if len(result) > 0 {
		t.Fatal("No results in filtered list expected")
	}
}
