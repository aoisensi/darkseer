package internal

import (
	"bufio"
	"fmt"
	"io"
)

type Decoder struct {
	r        *bufio.Reader
	names    []string
	elements []*Element
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: bufio.NewReader(r)}
}

func (d *Decoder) Decode() (*Element, error) {
	header, err := d.readHeader()
	if err != nil {
		return nil, err
	}
	switch header.EncodingName {
	case "text":
		return nil, fmt.Errorf("dmx: text file format is not supported")
	case "binary":
		return d.decodeBinary()
	default:
		return nil, fmt.Errorf("dmx: unknown encoding: %v", header.EncodingName)
	}
}
