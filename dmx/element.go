package dmx

import (
	"github.com/aoisensi/darkseer/dmx/internal"
	"github.com/google/uuid"
)

var elementMap map[uuid.UUID]any

type DmElement struct {
	Name          string
	Model         *DmeModel
	Skeleton      *DmeModel
	AnimationList *DmeAnimationList
}

func parseElement(e *internal.Element) *DmElement {
	if e == nil {
		return nil
	}
	if e.Type != "DmElement" {
		panic("dmx: invalid element type")
	}
	element := &DmElement{Name: e.Name}
	if model, ok := e.Attributes["model"]; ok && model != nil {
		element.Model = parseModel(model.(*internal.Element))
	}
	if skeleton, ok := e.Attributes["skeleton"]; ok && skeleton != nil {
		element.Skeleton = parseModel(skeleton.(*internal.Element))
	}
	if animationList, ok := e.Attributes["animationList"]; ok && animationList != nil {
		element.AnimationList = parseAnimationList(animationList.(*internal.Element))
	}
	return element
	// elementMap[e.ID] = result
	// return result
}
