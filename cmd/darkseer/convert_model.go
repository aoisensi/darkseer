package main

import (
	"strings"

	"github.com/aoisensi/darkseer/dmx"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
	"github.com/samber/lo"
)

func convertModel(element *dmx.DmElement) (*gltf.Document, error) {
	doc := gltf.NewDocument()
	doc.Scene = gltf.Index(0)
	scene := &gltf.Scene{}
	doc.Scenes = append(doc.Scenes, scene)

	materialMap := make(map[string]*uint32)

	getMaterialID := func(mtlName string) *uint32 {
		if id, ok := materialMap[mtlName]; ok {
			return id
		}
		id := uint32(len(doc.Materials))
		doc.Materials = append(doc.Materials, &gltf.Material{
			Name: lo.Must(lo.Last(strings.Split(mtlName, "/"))),
		})
		materialMap[mtlName] = &id
		return &id
	}

	for _, dag := range element.Model.Children {
		switch dmxDag := dag.(type) {
		case *dmx.DmeDag:
			if dmxDag.Mesh == nil {
				continue
			}
			dmxMesh := dmxDag.Mesh
			dmxVertexData := dmxMesh.CurrentState
			mesh := &gltf.Mesh{Name: strings.TrimSuffix(dmxDag.Name, "_mesh")}
			attribute := gltf.Attribute{
				"POSITION":   modeler.WritePosition(doc, dmxIndicesSort(dmxVertexData.PositionIndices, mulGlobalScale(dmxVertexData.Positions))),
				"NORMAL":     modeler.WriteNormal(doc, dmxIndicesSort(dmxVertexData.NormalsIndices, dmxVertexData.Normals)),
				"TEXCOORD_0": modeler.WriteTextureCoord(doc, dmxIndicesSort(dmxVertexData.TextureCoordinatesIndices, dmxUVToGLTFUV(dmxVertexData.TextureCoordinates))),
			}
			for _, dmxFaceSet := range dmxMesh.FaceSets {
				primitive := &gltf.Primitive{
					Attributes: attribute,
					Indices: gltf.Index(
						modeler.WriteIndices(
							doc,
							int32SliceTouint16Slice(
								dmxFacesetToGLTFIndices(dmxFaceSet.Faces),
							),
						),
					),
					Material: getMaterialID(dmxFaceSet.Material.MtlName),
				}
				mesh.Primitives = append(mesh.Primitives, primitive)
			}
			scene.Nodes = append(scene.Nodes, uint32(len(doc.Nodes)))
			node := &gltf.Node{
				Mesh: gltf.Index(uint32(len(doc.Meshes))),
			}
			doc.Nodes = append(doc.Nodes, node)
			doc.Meshes = append(doc.Meshes, mesh)
		}
	}
	return doc, nil
}

func dmxFacesetToGLTFIndices(faceset []int32) []int32 {
	result := make([]int32, 0, len(faceset))
	first := int32(-1)
	second := int32(-1)
	for _, i := range faceset {
		if first == -1 {
			first = i
			continue
		}
		if second == -1 {
			second = i
			continue
		}
		if i == -1 {
			first = -1
			second = -1
			continue
		}
		result = append(result, first, second, int32(i))
		second = i
	}
	return result
}

func dmxUVToGLTFUV(uv [][2]float32) [][2]float32 {
	for i := range uv {
		uv[i][1] = 1 - uv[i][1]
	}
	return uv
}

func dmxIndicesSort[T any](indices []int32, value []T) []T {
	result := make([]T, len(indices))
	for i, index := range indices {
		result[i] = value[index]
	}
	return result
}

func int32SliceTouint16Slice(values []int32) []uint16 {
	result := make([]uint16, len(values))
	for i, v := range values {
		result[i] = uint16(v)
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
