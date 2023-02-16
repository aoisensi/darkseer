package dmx

import (
	"io"

	"github.com/aoisensi/darkseer/dmx/internal"
)

type Decoder struct {
	d *internal.Decoder
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{d: internal.NewDecoder(r)}
}

func (d *Decoder) Decode() (*DmElement, error) {
	dmx, err := d.d.Decode()
	if err != nil {
		return nil, err
	}
	return parseElement(dmx), nil
}
