package main

import (
	"fmt"
	"strings"

	"github.com/aoisensi/darkseer/dmx"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
	"github.com/samber/lo"
)

func convertModel(title string, dmxElement *dmx.DmElement) (*gltf.Document, error) {
	fmt.Println(title)
	doc := gltf.NewDocument()
	doc.Scene = gltf.Index(0)
	scene := &gltf.Scene{}
	doc.Scenes = append(doc.Scenes, scene)

	materialMap := make(map[string]*uint32)
	jointMap := make(map[string]uint32)

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

	var skinID *uint32

	// Find skins
	for _, dmxChild := range dmxElement.Skeleton.Children {
		var joints []uint32
		if dmxJoint, ok := dmxChild.(*dmx.DmeJoint); ok {
			var addJoint func(*dmx.DmeJoint) *uint32
			addJoint = func(dmxJoint *dmx.DmeJoint) *uint32 {
				if strings.Contains(dmxJoint.Name, "End") || strings.Contains(dmxJoint.Name, "parentConstraint") {
					return nil
				}
				nodeID := uint32(len(doc.Nodes))
				node := &gltf.Node{
					Name:        dmxJoint.Name,
					Translation: mulGlobalScale(dmxJoint.Transform.Position),
					Rotation:    dmxJoint.Transform.Orientation,
				}
				jointMap[dmxJoint.Name] = nodeID
				joints = append(joints, nodeID)
				doc.Nodes = append(doc.Nodes, node)
				for _, child := range dmxJoint.Children {
					if child, ok := child.(*dmx.DmeJoint); ok {
						childID := addJoint(child)
						if childID != nil {
							node.Children = append(node.Children, *childID)
						}

					}
				}
				return &nodeID
			}
			rootBoneID := addJoint(dmxJoint)
			scene.Nodes = append(scene.Nodes, *rootBoneID)
			_skinID := uint32(len(doc.Skins))
			skinID = &_skinID
			doc.Skins = append(doc.Skins, &gltf.Skin{
				Name:     "Armature",
				Skeleton: rootBoneID,
				Joints:   joints,
			})
		}
	}

	// Find meshes
	var findMesh func([]dmx.IDag)
	findMesh = func(dmxDags []dmx.IDag) {
		for _, dmxDag := range dmxDags {
			dmxDag, ok := dmxDag.(*dmx.DmeDag)
			if !ok {
				continue
			}
			meshName := strings.TrimSuffix(dmxDag.Name, "_mesh")
			if dmxMesh := dmxDag.Mesh; dmxMesh != nil {
				dmxVertexData := dmxMesh.CurrentState
				mesh := &gltf.Mesh{Name: meshName}
				attribute := gltf.Attribute{
					"POSITION":   modeler.WritePosition(doc, dmxIndicesSort(dmxVertexData.PositionIndices, mulGlobalScale(dmxVertexData.Positions))),
					"NORMAL":     modeler.WriteNormal(doc, dmxIndicesSort(dmxVertexData.NormalsIndices, dmxVertexData.Normals)),
					"TEXCOORD_0": modeler.WriteTextureCoord(doc, dmxIndicesSort(dmxVertexData.TextureCoordinatesIndices, dmxUVToGLTFUV(dmxVertexData.TextureCoordinates))),
				}
				if len(dmxElement.Model.JointTransforms) > 0 {
					jc := int(dmxVertexData.JointCount)
					jointIndeices := make([][4]uint8, 0)
					jointWeights := make([][4]float32, 0)
					for i := range dmxVertexData.Positions {
						var ji [4]uint8
						var jw [4]float32
						for j := 0; j < 4; j++ {
							if j < jc {
								_ji := dmxVertexData.JointIndices[i*jc+j]
								ji[j] = uint8(jointMap[dmxElement.Model.JointTransforms[_ji].Name])
								jw[j] = dmxVertexData.JointWeights[i*jc+j]
							}
						}
						jointIndeices = append(jointIndeices, ji)
						jointWeights = append(jointWeights, jw)
					}
					attribute["JOINTS_0"] = modeler.WriteJoints(doc, dmxIndicesSort(dmxVertexData.PositionIndices, jointIndeices))
					attribute["WEIGHTS_0"] = modeler.WriteWeights(doc, dmxIndicesSort(dmxVertexData.PositionIndices, jointWeights))
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
					Name: meshName,
					Mesh: gltf.Index(uint32(len(doc.Meshes))),
					Skin: skinID,
				}
				doc.Nodes = append(doc.Nodes, node)
				doc.Meshes = append(doc.Meshes, mesh)
			}
			findMesh(dmxDag.Dag().Children)
		}
	} // findMesh
	if dmxElement.Model != nil {
		findMesh(dmxElement.Model.Children)
	}

	// Find animations
	if dmxElement.AnimationList != nil {
		for _, dmxAnimation := range dmxElement.AnimationList.Animations {
			animation := &gltf.Animation{
				Name: title,
			}
			for _, dmxChannel := range dmxAnimation.Channels {
				samplerID := uint32(len(animation.Samplers))
				sampler := &gltf.AnimationSampler{
					Interpolation: gltf.InterpolationLinear,
				}
				joint, found := jointMap[dmxChannel.ToElement.Name]
				if !found {
					continue
				}
				channel := &gltf.Channel{
					Sampler: gltf.Index(samplerID),
					Target: gltf.ChannelTarget{
						Node: gltf.Index(joint),
					},
				}

				writeInput := func(input []int32) {
					times := mulTimes(input)
					sampler.Input = modeler.WriteAccessor(
						doc,
						gltf.TargetNone,
						times,
					)
					accessor := doc.Accessors[int(sampler.Input)]
					accessor.Min = []float32{lo.Min(times)}
					accessor.Max = []float32{lo.Max(times)}
				}

				if dmxChannel.LogVector3 != nil {
					channel.Target.Path = gltf.TRSTranslation
					writeInput(dmxChannel.LogVector3.Layers[0].Times)
					sampler.Output = modeler.WritePosition(
						doc,
						mulGlobalScale(dmxChannel.LogVector3.Layers[0].Values),
					)
				} else if dmxChannel.LogQuaternion != nil {
					channel.Target.Path = gltf.TRSRotation
					writeInput(dmxChannel.LogQuaternion.Layers[0].Times)
					sampler.Output = modeler.WriteAccessor(
						doc,
						gltf.TargetNone,
						dmxChannel.LogQuaternion.Layers[0].Values,
					)
				} else {
					continue
				}
				bvID := doc.Accessors[sampler.Output].BufferView
				doc.BufferViews[*bvID].Target = gltf.TargetNone
				animation.Samplers = append(animation.Samplers, sampler)
				animation.Channels = append(animation.Channels, channel)
			}
			doc.Animations = append(doc.Animations, animation)
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

type GlobalScaler interface {
	[][3]float32 | [3]float32
}

func mulGlobalScale[T GlobalScaler](values T) T {
	switch values := any(values).(type) {
	case [][3]float32:
		for i := range values {
			values[i][0] *= float32(*argScale)
			values[i][1] *= float32(*argScale)
			values[i][2] *= float32(*argScale)
		}
		return any(values).(T)
	case [3]float32:
		values[0] *= float32(*argScale)
		values[1] *= float32(*argScale)
		values[2] *= float32(*argScale)
		return any(values).(T)
	default:
		panic("unreachable")
	}
}

func mulTimes(s []int32) []float32 {
	result := make([]float32, len(s))
	for i, v := range s {
		result[i] = float32(v) * 0.0001
	}
	return result
}

/*
func makeMatrix(t *dmx.DmeTransform) [16]float32 {
	pos := t.Position
	quat := t.Orientation
	qx, qy, qz, qw := quat[0], quat[1], quat[2], quat[3]
	return [16]float32{
		1 - 2*qy*qy - 2*qz*qz, 2*qx*qy - 2*qz*qw, 2*qx*qz + 2*qy*qw, 0,
		2*qx*qy + 2*qz*qw, 1 - 2*qx*qx - 2*qz*qz, 2*qy*qz - 2*qx*qw, 0,
		2*qx*qz - 2*qy*qw, 2*qy*qz + 2*qx*qw, 1 - 2*qx*qx - 2*qy*qy, 0,
		pos[0], pos[1], pos[2], 1,
	}
}

func scaleMatrix(m [16]float32, scale float32) [16]float32 {
	for i := 0; i < 3; i++ {
		m[i*4+0] *= scale
		m[i*4+1] *= scale
		m[i*4+2] *= scale
	}
	return m
}
*/
