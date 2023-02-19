package dmx

import (
	"image/color"

	"github.com/aoisensi/darkseer/dmx/internal"
)

type DmeAnimationList struct {
	Name       string
	Animations []*DmeChannelsClip
}

func parseAnimationList(e *internal.Element) *DmeAnimationList {
	if e == nil {
		return nil
	}
	return &DmeAnimationList{
		Name:       e.Name,
		Animations: parseChannelsClipList(e.Attributes["animations"].([]*internal.Element)),
	}
}

type DmeChannelsClip struct {
	Name      string
	TimeFrame *DmeTimeFrame
	Color     color.RGBA
	Text      string
	Mute      bool
	Channels  []*DmeChannel
	FrameRate int32
}

func parseChannelsClip(e *internal.Element) *DmeChannelsClip {
	if e == nil {
		return nil
	}
	return &DmeChannelsClip{
		Name:      e.Name,
		TimeFrame: parseTimeFrame(e.Attributes["timeFrame"].(*internal.Element)),
		Color:     e.Attributes["color"].(color.RGBA),
		Text:      e.Attributes["text"].(string),
		Mute:      e.Attributes["mute"].(bool),
		Channels:  parseChannelList(e.Attributes["channels"].([]*internal.Element)),
		FrameRate: e.Attributes["frameRate"].(int32),
	}
}

func parseChannelsClipList(e []*internal.Element) []*DmeChannelsClip {
	if e == nil {
		return nil
	}
	list := make([]*DmeChannelsClip, len(e))
	for i, v := range e {
		list[i] = parseChannelsClip(v)
	}
	return list
}

type DmeTimeFrame struct {
	Name         string
	StartTime    int32
	DurationTime int32
	OffsetTime   int32
	Scale        float32
}

func parseTimeFrame(e *internal.Element) *DmeTimeFrame {
	if e == nil {
		return nil
	}
	return &DmeTimeFrame{
		Name:         e.Name,
		StartTime:    e.Attributes["startTime"].(int32),
		DurationTime: e.Attributes["durationTime"].(int32),
		OffsetTime:   e.Attributes["offsetTime"].(int32),
		Scale:        e.Attributes["scale"].(float32),
	}
}

type DmeChannel struct {
	Name          string
	FromAttribute string
	FromIndex     int32
	ToElement     *DmeTransform
	ToAttribute   string
	ToIndex       int32
	LogQuaternion *DmeLog[[4]float32]
	LogVector3    *DmeLog[[3]float32]
}

func parseChannel(e *internal.Element) *DmeChannel {
	if e == nil {
		return nil
	}
	channel := &DmeChannel{
		Name:          e.Name,
		FromAttribute: e.Attributes["fromAttribute"].(string),
		FromIndex:     e.Attributes["fromIndex"].(int32),
		ToElement:     parseTransform(e.Attributes["toElement"].(*internal.Element)),
		ToAttribute:   e.Attributes["toAttribute"].(string),
		ToIndex:       e.Attributes["toIndex"].(int32),
	}
	if log, ok := e.Attributes["log"].(*internal.Element); ok {
		switch log.Type {
		case "DmeQuaternionLog":
			channel.LogQuaternion = parseLog[[4]float32](log)
		case "DmeVector3Log":
			channel.LogVector3 = parseLog[[3]float32](log)
		}
	}
	return channel
}

func parseChannelList(e []*internal.Element) []*DmeChannel {
	if e == nil {
		return nil
	}
	list := make([]*DmeChannel, len(e))
	for i, v := range e {
		list[i] = parseChannel(v)
	}
	return list
}

// Log //

type LogType interface {
	[4]float32 | [3]float32
}

type DmeLog[T LogType] struct {
	Name            string
	Layers          []*DmeLogLayer[T]
	UseDefaultValue bool
	DefaultValue    T
}

func parseLog[T LogType](e *internal.Element) *DmeLog[T] {
	if e == nil {
		return nil
	}
	return &DmeLog[T]{
		Name:            e.Name,
		Layers:          parseLayerList[T](e.Attributes["layers"].([]*internal.Element)),
		UseDefaultValue: e.Attributes["usedefaultvalue"].(bool),
		DefaultValue:    e.Attributes["defaultvalue"].(T),
	}
}

// Layer //
type DmeLogLayer[T LogType] struct {
	Name   string
	Times  []int32
	Values []T
}

func parseLayer[T LogType](e *internal.Element) *DmeLogLayer[T] {
	if e == nil {
		return nil
	}
	return &DmeLogLayer[T]{
		Name:   e.Name,
		Times:  e.Attributes["times"].([]int32),
		Values: e.Attributes["values"].([]T),
	}
}

func parseLayerList[T LogType](e []*internal.Element) []*DmeLogLayer[T] {
	if e == nil {
		return nil
	}
	list := make([]*DmeLogLayer[T], len(e))
	for i, v := range e {
		list[i] = parseLayer[T](v)
	}
	return list
}
