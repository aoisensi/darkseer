package dmx

import "github.com/aoisensi/darkseer/dmx/internal"

type DmeModelRoot struct {
	Model    DmeModel
	Skeleton DmeModel
	// TODO CombinationOperator DmeCombinationOperator
}

func parseModelRoot(e *internal.Element) *DmeModelRoot {
	if e == nil {
		return nil
	}
	if e.Type != "DmeModelRoot" {
		panic("dmx: invalid element type")
	}
	return &DmeModelRoot{
		Model:    *parseModel(e.Attributes["model"].(*internal.Element)),
		Skeleton: *parseModel(e.Attributes["skeleton"].(*internal.Element)),
	}
}

type DmeModel struct {
	Name      string
	Transform *DmeTransform
	Visible   bool
	Children  []IDag
}

func parseModel(e *internal.Element) *DmeModel {
	if e == nil {
		return nil
	}
	if e.Type != "DmeModel" {
		panic("dmx: invalid element type")
	}
	return &DmeModel{
		Name:      e.Name,
		Transform: parseTransform(e.Attributes["transform"].(*internal.Element)),
		Visible:   e.Attributes["visible"].(bool),
		Children:  parseDagList(e.Attributes["children"]),
	}
}

func parseModelList(e []*internal.Element) []*DmeModel {
	if e == nil {
		return nil
	}
	list := make([]*DmeModel, len(e))
	for i, v := range e {
		list[i] = parseModel(v)
	}
	return list
}

type DmeAttachment struct {
	*DmeDag
	Visible bool
	// IsRigid        bool
	// IsWorldAligned bool
}

func parseAttachment(e *internal.Element) *DmeAttachment {
	if e == nil {
		return nil
	}
	if e.Type != "DmeAttachment" {
		panic("dmx: invalid element type")
	}
	return &DmeAttachment{
		DmeDag:  parseOnlyDag(e),
		Visible: e.Attributes["visible"].(bool),
		// IsRigid:        e.Attributes["isrigid"].(bool),
		// IsWorldAligned: e.Attributes["isworldaligned"].(bool),
	}
}

func parseAttachmentList(e []*internal.Element) []*DmeAttachment {
	if e == nil {
		return nil
	}
	list := make([]*DmeAttachment, len(e))
	for i, v := range e {
		list[i] = parseAttachment(v)
	}
	return list
}

type DmeMesh struct {
	*DmeDag
	Visible      bool
	CurrentState *DmeVertexData
	BaseStates   []*DmeVertexData
	DeltaStates  []*DmeVertexData
	FaceSets     []*DmeFaceSet
}

func parseMesh(e *internal.Element) *DmeMesh {
	if e == nil {
		return nil
	}
	if e.Type != "DmeMesh" {
		panic("dmx: invalid element type")
	}
	return &DmeMesh{
		DmeDag:       parseOnlyDag(e),
		Visible:      e.Attributes["visible"].(bool),
		CurrentState: parseVertexData(e.Attributes["currentState"].(*internal.Element)),
		BaseStates:   parseVertexDataList(e.Attributes["baseStates"].([]*internal.Element)),
		DeltaStates:  parseVertexDataList(e.Attributes["deltaStates"].([]*internal.Element)),
		FaceSets:     parseFaceSetList(e.Attributes["faceSets"].([]*internal.Element)),
	}
}

func parseMeshList(e []*internal.Element) []*DmeMesh {
	if e == nil {
		return nil
	}
	list := make([]*DmeMesh, len(e))
	for i, v := range e {
		list[i] = parseMesh(v)
	}
	return list
}

type DmeVertexData struct {
	VertexFormat              []string
	JointCount                int32
	Positions                 [][3]float32
	PositionIndices           []int32
	Normals                   [][3]float32
	NormalsIndices            []int32
	TextureCoordinates        [][2]float32
	TextureCoordinatesIndices []int32
	// TODO Balance
	// TODO BalanceIndices
	JointWeights []float32
	JointIndices []int32
}

func parseVertexData(e *internal.Element) *DmeVertexData {
	if e == nil {
		return nil
	}
	if e.Type != "DmeVertexData" {
		panic("dmx: invalid element type")
	}
	result := &DmeVertexData{
		VertexFormat:              e.Attributes["vertexFormat"].([]string),
		JointCount:                e.Attributes["jointCount"].(int32),
		Positions:                 e.Attributes["positions"].([][3]float32),
		PositionIndices:           e.Attributes["positionsIndices"].([]int32),
		Normals:                   e.Attributes["normals"].([][3]float32),
		NormalsIndices:            e.Attributes["normalsIndices"].([]int32),
		TextureCoordinates:        e.Attributes["textureCoordinates"].([][2]float32),
		TextureCoordinatesIndices: e.Attributes["textureCoordinatesIndices"].([]int32),
	}
	if e.Attributes["jointWeights"] != nil {
		result.JointWeights = e.Attributes["jointWeights"].([]float32)
	}
	if e.Attributes["jointIndeices"] != nil {
		result.JointIndices = e.Attributes["jointIndices"].([]int32)
	}
	return result
}

type DmeFaceSet struct {
	Material *DmeMaterial
	Faces    []int32
}

func parseFaceSet(e *internal.Element) *DmeFaceSet {
	if e == nil {
		return nil
	}
	if e.Type != "DmeFaceSet" {
		panic("dmx: invalid element type")
	}
	return &DmeFaceSet{
		Material: parseMaterial(e.Attributes["material"].(*internal.Element)),
		Faces:    e.Attributes["faces"].([]int32),
	}
}

func parseFaceSetList(e []*internal.Element) []*DmeFaceSet {
	if e == nil {
		return nil
	}
	list := make([]*DmeFaceSet, len(e))
	for i, v := range e {
		list[i] = parseFaceSet(v)
	}
	return list
}

type DmeMaterial struct {
	MtlName string
}

func parseMaterial(e *internal.Element) *DmeMaterial {
	if e == nil {
		return nil
	}
	if e.Type != "DmeMaterial" {
		panic("dmx: invalid element type")
	}
	return &DmeMaterial{
		MtlName: e.Attributes["mtlName"].(string),
	}
}

type DmeVertexDeltaData struct {
	VertexFormat     []string
	FlipVCoordinates bool
	Corrected        bool
	Positions        [][3]float32
	PositionIndices  []int
	Normals          [][3]float32
	NormalsIndices   []int
	Wrinkle          []float32
	WrinkleIndices   []int
}

func parseVertexDeltaData(e *internal.Element) *DmeVertexDeltaData {
	if e == nil {
		return nil
	}
	if e.Type != "DmeVertexDeltaData" {
		panic("dmx: invalid element type")
	}
	return &DmeVertexDeltaData{
		VertexFormat:     e.Attributes["vertexformat"].([]string),
		FlipVCoordinates: e.Attributes["flipvcoordinates"].(bool),
		Corrected:        e.Attributes["corrected"].(bool),
		Positions:        e.Attributes["positions"].([][3]float32),
		PositionIndices:  e.Attributes["positionindices"].([]int),
		Normals:          e.Attributes["normals"].([][3]float32),
		NormalsIndices:   e.Attributes["normalsindices"].([]int),
		Wrinkle:          e.Attributes["wrinkle"].([]float32),
		WrinkleIndices:   e.Attributes["wrinkleindices"].([]int),
	}
}

func parseVertexDeltaDataList(e []*internal.Element) []*DmeVertexDeltaData {
	if e == nil {
		return nil
	}
	list := make([]*DmeVertexDeltaData, len(e))
	for i, v := range e {
		list[i] = parseVertexDeltaData(v)
	}
	return list
}

func parseVertexDataList(e []*internal.Element) []*DmeVertexData {
	if e == nil {
		return nil
	}
	list := make([]*DmeVertexData, len(e))
	for i, v := range e {
		list[i] = parseVertexData(v)
	}
	return list
}
