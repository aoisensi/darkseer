package internal

import (
	"fmt"
	"strconv"
	"strings"
)

type Header struct {
	EncodingName    string
	EncodingVersion int
	FormatName      string
	FormatVersion   int
}

func (d *Decoder) readHeader() (*Header, error) {
	line, _, err := d.r.ReadLine()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(line), " ")
	if len(lines) != 9 || lines[0] != "<!--" ||
		lines[1] != "dmx" || lines[2] != "encoding" ||
		lines[5] != "format" || lines[8] != "-->" {
		return nil, fmt.Errorf("dmx: this is not dmx file")
	}
	ev, err := strconv.Atoi(lines[4])
	if err != nil {
		return nil, fmt.Errorf("dmx: this is not dmx file")
	}
	fv, err := strconv.Atoi(lines[7])
	if err != nil {
		return nil, fmt.Errorf("dmx: this is not dmx file")
	}
	return &Header{
		EncodingName:    lines[3],
		EncodingVersion: ev,
		FormatName:      lines[6],
		FormatVersion:   fv,
	}, nil
}
