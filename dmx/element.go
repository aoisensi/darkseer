package dmx

import (
	"github.com/aoisensi/darkseer/dmx/internal"
	"github.com/google/uuid"
)

var elementMap map[uuid.UUID]any

type DmElement struct {
	Name     string
	Model    *DmeModel
	Skeleton *DmeModel
}

func parseElement(e *internal.Element) *DmElement {
	if e == nil {
		return nil
	}
	if e.Type != "DmElement" {
		panic("dmx: invalid element type")
	}
	return &DmElement{
		Name:     e.Name,
		Model:    parseModel(e.Attributes["model"].(*internal.Element)),
		Skeleton: parseModel(e.Attributes["skeleton"].(*internal.Element)),
	}
	// elementMap[e.ID] = result
	// return result
}
