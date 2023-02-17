package main

import (
	"math"

	"github.com/aoisensi/darkseer/dmx"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

func convertModel(element *dmx.DmElement) (*gltf.Document, error) {
	doc := gltf.NewDocument()
	doc.Scene = gltf.Index(0)
	scene := &gltf.Scene{}
	doc.Scenes = append(doc.Scenes, scene)

	for _, dag := range element.Model.Children {
		switch dmxDag := dag.(type) {
		case *dmx.DmeDag:
			if dmxDag.Mesh == nil {
				continue
			}
			dmxMesh := dmxDag.Mesh
			dmxVertexData := dmxMesh.CurrentState
			dmxFaceSet := dmxMesh.FaceSets[0]

			primitive := &gltf.Primitive{
				Attributes: gltf.Attribute{
					"POSITION": modeler.WritePosition(doc, mulGlobalScale(dmxVertexData.Positions)),
				},
				Indices: gltf.Index(
					modeler.WriteIndices(doc, dmxIndicesToGLTFIndices(dmxVertexData.PositionIndices, dmxFaceSet.Faces)),
				),
			}
			scene.Nodes = append(scene.Nodes, uint32(len(doc.Nodes)))
			node := &gltf.Node{
				Mesh: gltf.Index(uint32(len(doc.Meshes))),
			}
			doc.Nodes = append(doc.Nodes, node)

			doc.Meshes = append(doc.Meshes, &gltf.Mesh{
				Name:       dmxDag.Name,
				Primitives: []*gltf.Primitive{primitive},
			})
		}
	}
	return doc, nil
}

func dmxIndicesToGLTFIndices(indices, faceset []int32) []uint16 {
	result := make([]uint16, 0, 256)
	first := math.MaxUint16
	second := math.MaxUint16
	for _, i := range faceset {
		if first == math.MaxUint16 {
			first = int(i)
			continue
		}
		if second == math.MaxUint16 {
			second = int(i)
			continue
		}
		if i == -1 {
			first = math.MaxUint16
			second = math.MaxUint16
			continue
		}
		result = append(result, uint16(indices[first]), uint16(indices[second]), uint16(indices[i]))
		second = int(i)
	}
	return result
}

func mulGlobalScale[T any](values T) T {
	switch values := any(values).(type) {
	case [][3]float32:
		for i := range values {
			values[i][0] *= float32(*argScale)
			values[i][1] *= float32(*argScale)
			values[i][2] *= float32(*argScale)
		}
		return any(values).(T)
	default:
		panic("fjdkslajfdklsa")
	}
}
