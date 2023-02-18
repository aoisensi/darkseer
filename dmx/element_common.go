package dmx

import "github.com/aoisensi/darkseer/dmx/internal"

type DmeDag struct {
	Name       string
	Transform  *DmeTransform
	Visible    bool
	Children   []IDag
	Mesh       *DmeMesh
	Attachment *DmeAttachment
}

func (d *DmeDag) Dag() *DmeDag {
	return d
}

type IDag interface {
	Dag() *DmeDag
}

func parseDag(e *internal.Element) IDag {
	if e == nil {
		return nil
	}
	switch e.Type {
	default:
		panic("dmx: invalid element type")
	case "DmeJoint":
		return parseJoint(e)
	case "DmeAttachment":
		return parseAttachment(e)
	case "DmeMesh":
		return parseMesh(e)
	case "DmeDag":
		return parseOnlyDag(e)
	}
}

func parseOnlyDag(e *internal.Element) *DmeDag {
	if e == nil {
		return nil
	}
	result := &DmeDag{
		Name:    e.Name,
		Visible: e.Attributes["visible"].(bool),
	}
	if transform, ok := e.Attributes["transform"]; ok && transform != nil {
		result.Transform = parseTransform(transform.(*internal.Element))
	}
	if children, ok := e.Attributes["children"]; ok && children != nil {
		eChildlen := children.([]*internal.Element)
		children := make([]IDag, len(eChildlen))
		for i, c := range eChildlen {
			children[i] = parseDag(c)
		}
		result.Children = children
	}
	if shape, ok := e.Attributes["shape"]; ok && shape != nil {
		shape := shape.(*internal.Element)
		switch shape.Type {
		case "DmeMesh":
			result.Mesh = parseMesh(shape)
		case "DmeAttachment":
			result.Attachment = parseAttachment(shape)
		}
	}
	return result
}

func parseDagList(e any) []IDag {
	if e == nil {
		return nil
	}
	eChildlen := e.([]*internal.Element)
	children := make([]IDag, len(eChildlen))
	for i, c := range eChildlen {
		children[i] = parseDag(c)
	}
	return children
}

type DmeJoint struct {
	*DmeDag
	Transform            *DmeTransform
	Visible              bool
	Children             []IDag
	LockInfluenceWeights bool
}

func parseJoint(e *internal.Element) *DmeJoint {
	if e == nil {
		return nil
	}
	// if result, ok := elementMap[e.ID]; ok {
	// 	return result.(*DmeJoint)
	// }
	if e.Type != "DmeJoint" {
		panic("dmx: invalid element type")
	}
	eChildlen := e.Attributes["children"].([]*internal.Element)
	children := make([]IDag, len(eChildlen))
	for i, c := range eChildlen {
		children[i] = parseDag(c)
	}
	joint := &DmeJoint{
		DmeDag:    parseOnlyDag(e),
		Transform: parseTransform(e.Attributes["transform"].(*internal.Element)),
		Visible:   e.Attributes["visible"].(bool),
		Children:  children,
	}
	if lockInfluenceWeights, ok := e.Attributes["lockInfluenceWeights"]; ok {
		joint.LockInfluenceWeights = lockInfluenceWeights.(bool)
	}
	return joint
}

type DmeTransformList struct {
	Transforms []*DmeTransform
}

func parseTransformList(e *internal.Element) *DmeTransformList {
	if e == nil {
		return nil
	}
	if e.Type != "DmeTransformList" {
		panic("dmx: invalid element type")
	}
	eTransforms := e.Attributes["transforms"].([]*internal.Element)
	transforms := make([]*DmeTransform, len(eTransforms))
	for i, t := range eTransforms {
		transforms[i] = parseTransform(t)
	}
	return &DmeTransformList{
		Transforms: transforms,
	}
}

type DmeTransform struct {
	Name        string
	Position    [3]float32
	Orientation [4]float32
}

func parseTransform(e *internal.Element) *DmeTransform {
	if e == nil {
		return nil
	}
	if e.Type != "DmeTransform" {
		panic("dmx: invalid element type")
	}
	return &DmeTransform{
		Name:        e.Name,
		Position:    e.Attributes["position"].([3]float32),
		Orientation: e.Attributes["orientation"].([4]float32),
	}
}
