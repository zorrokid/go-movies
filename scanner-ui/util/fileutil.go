package util

import "fyne.io/fyne/v2"

func FilterByExtension(uriList []fyne.URI, extension ...string) []fyne.URI {
	filteredUriList := make([]fyne.URI, 0)
	for _, uri := range uriList {
		for _, ext := range extension {
			if uri.Extension() == ext {
				filteredUriList = append(filteredUriList, uri)
			}
		}
	}
	return filteredUriList
}
