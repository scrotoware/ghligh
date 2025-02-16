package document

import (
	"github.com/scrotadamus/ghligh/go-poppler"

	"encoding/json"
)

func unmarshallHighlights(jsonData string) (AnnotsMap, error) {
	var annotsMap AnnotsMap

	err := json.Unmarshal([]byte(jsonData), &struct {
		Highlights *AnnotsMap `json:"highlights"`
	}{
		Highlights: &annotsMap,
	})

	return annotsMap, err
}

func isInPage(a *poppler.Annot, p *poppler.Page) bool {
	annots := p.GetAnnots()
	for _, annot := range annots {
		if popplerAnnotsMatch(a, annot) {
			return true
		}
	}

	return false
}
