package internal

import (
	"encoding/binary"
	"fmt"
	"image"
	"io"
	"time"

	"github.com/google/uuid"
)

func (d *Decoder) decodeBinary() (*Element, error) {
	d.r.ReadByte()
	var numNames int16
	binary.Read(d.r, binary.LittleEndian, &numNames)
	d.names = make([]string, numNames)
	for i := range d.names {
		var err error
		d.names[i], err = d.decodeString()
		if err != nil {
			return nil, fmt.Errorf("dmx: failed to read")
		}
	}
	var numHeaders int32
	binary.Read(d.r, binary.LittleEndian, &numHeaders) // num elements
	d.elements = make([]*Element, numHeaders)
	for i := range d.elements {
		var err error
		d.elements[i], err = d.decodeBinaryElmHeader()
		if err != nil {
			return nil, err
		}
	}
	for i := range d.elements {
		var err error
		d.elements[i].Attributes, err = d.decodeBinaryAttributes()
		if err != nil {
			return nil, err
		}
	}
	return d.elements[0], nil
}

func (d *Decoder) decodeBinaryElmHeader() (*Element, error) {
	e := new(Element)
	e.Type = d.decodeName()
	e.Name, _ = d.decodeString()
	e.ID, _ = d.decodeBinaryID()
	return e, nil
}

func (d *Decoder) decodeName() string {
	var nameID int16
	binary.Read(d.r, binary.LittleEndian, &nameID)
	return d.names[nameID]
}

func (d *Decoder) decodeBinaryID() (uuid.UUID, error) {
	b := make([]byte, 16)
	io.ReadFull(d.r, b)
	return uuid.FromBytes(b)
}

func (d *Decoder) decodeBinaryAttributes() (map[string]any, error) {
	var num int32
	binary.Read(d.r, binary.LittleEndian, &num)
	result := make(map[string]any)
	for i := 0; i < int(num); i++ {
		name := d.decodeName()
		typeID, err := d.r.ReadByte()
		if err != nil {
			return nil, err
		}
		result[name] = d.decodeBinaryAttribute(typeID)
	}
	return result, nil
}

func (d *Decoder) decodeBinaryAttribute(typeID byte) any {
	switch typeID {
	case 0: // nil
		return nil
	case 1: // *Element
		var elemID int32
		binary.Read(d.r, binary.LittleEndian, &elemID)
		if elemID == -1 {
			return nil
		}
		return d.elements[elemID]
	case 2: // int32
		var value int32
		binary.Read(d.r, binary.LittleEndian, &value)
		return value
	case 3: // float32
		var value float32
		binary.Read(d.r, binary.LittleEndian, &value)
		return value
	case 4: // bool
		value, _ := d.r.ReadByte()
		return value != 0
	case 5: // string
		value, _ := d.decodeString()
		return value
	case 6: // []byte
		var length int32
		binary.Read(d.r, binary.LittleEndian, &length)
		data := make([]byte, length)
		io.ReadFull(d.r, data)
		return data
	case 7: // time.Duration
		var value int32
		binary.Read(d.r, binary.LittleEndian, &value)
		data := time.Microsecond * 100 * time.Duration(value)
		return data
	case 8: // image.RGBA
		var value image.RGBA
		binary.Read(d.r, binary.LittleEndian, &value)
		return value
	case 9: // [2]float32 Vec2
		var value [2]float32
		binary.Read(d.r, binary.LittleEndian, &value)
		return value
	case 10: // [3]float32 Vec3
		var value [3]float32
		binary.Read(d.r, binary.LittleEndian, &value)
		return value
	case 11: // [4]float32 Vec4
		var value [4]float32
		binary.Read(d.r, binary.LittleEndian, &value)
		return value
	case 12: // [3]float32 Angle
		var value [3]float32
		binary.Read(d.r, binary.LittleEndian, &value)
		return value
	case 13: // [4]float32 Quat
		var value [4]float32
		binary.Read(d.r, binary.LittleEndian, &value)
		return value
	case 14: // [4][4]float32 Matrix
		var value [4][4]float32
		binary.Read(d.r, binary.LittleEndian, &value)
		return value
	}
	typeID -= 14
	// read length
	var length int32
	binary.Read(d.r, binary.LittleEndian, &length)

	switch typeID {
	case 1:
		result := make([]*Element, length)
		for i := range result {
			var elemID int32
			binary.Read(d.r, binary.LittleEndian, &elemID)
			if elemID == -1 {
				result[i] = nil
			} else {
				result[i] = d.elements[elemID]
			}
		}
		return result
	case 2:
		result := make([]int32, length)
		binary.Read(d.r, binary.LittleEndian, result)
		return result
	case 3:
		result := make([]float32, length)
		binary.Read(d.r, binary.LittleEndian, result)
		return result
	case 4:
		result := make([]bool, length)
		for i := range result {
			value, _ := d.r.ReadByte()
			result[i] = value != 0
		}
		return result
	case 5:
		result := make([]string, length)
		for i := range result {
			result[i], _ = d.decodeString()
		}
		return result
	case 6:
		result := make([][]byte, length)
		for i := range result {
			var length int32
			binary.Read(d.r, binary.LittleEndian, &length)
			data := make([]byte, length)
			io.ReadFull(d.r, data)
			result[i] = data
		}
		return result
	case 7:
		result := make([]time.Duration, length)
		for i := range result {
			var value int32
			binary.Read(d.r, binary.LittleEndian, &value)
			result[i] = time.Microsecond * 100 * time.Duration(value)
		}
		return result
	case 8:
		result := make([]image.RGBA, length)
		binary.Read(d.r, binary.LittleEndian, &result)
		return result
	case 9:
		result := make([][2]float32, length)
		binary.Read(d.r, binary.LittleEndian, &result)
		return result
	case 10:
		result := make([][3]float32, length)
		binary.Read(d.r, binary.LittleEndian, &result)
		return result
	case 11:
		result := make([][4]float32, length)
		binary.Read(d.r, binary.LittleEndian, &result)
		return result
	case 12:
		result := make([][3]float32, length)
		binary.Read(d.r, binary.LittleEndian, &result)
		return result
	case 13:
		result := make([][4]float32, length)
		binary.Read(d.r, binary.LittleEndian, &result)
		return result
	case 14:
		result := make([][4][4]float32, length)
		binary.Read(d.r, binary.LittleEndian, &result)
		return result
	}
	panic("unreachable")
}

func (d *Decoder) decodeString() (string, error) {
	str, err := d.r.ReadString(0)
	if err != nil {
		return "", err
	}
	return str[:len(str)-1], nil
}
