// Copyright 2018 Silke Hofstra
//
// Licensed under the EUPL
//

package main

import (
	"encoding/binary"
	"fmt"
	"image/color"
)

const (
	bwThreshold = 3 * 0xffff // pure white
)

type Mode struct {
	Density int
	Bits    int
	Convert func(pixels []color.Color) string
}

var modes = map[string]Mode{
	"bw":    {8, 8, pixelToBW},
	"565":   {1, 16, pixelTo565},
	"565le": {1, 16, pixelTo565LE},
	"rgb":   {1, 24, pixelToRGB},
	"rgba":  {1, 32, pixelToRGBA},
}

// pixelToBW converts a black/white pixel to
func pixelToBW(pixels []color.Color) string {
	var px byte
	for _, p := range pixels {
		px = px << 1
		if r, g, b, _ := p.RGBA(); r+g+b >= bwThreshold {
			px |= 1
		}
	}

	return fmt.Sprintf("0x%02x", px)
}

func pixelTo565(pixels []color.Color) string {
	r, g, b, _ := pixels[0].RGBA()
	var px uint16
	px |= uint16(r) & 0xF800
	px |= (uint16(g) >> 5) & 0x07E0
	px |= (uint16(b) >> 11) & 0x001F
	return fmt.Sprintf("0x%04x", px)
}

func pixelTo565LE(pixels []color.Color) string {
	r, g, b, _ := pixels[0].RGBA()
	var px uint16
	px |= uint16(r) & 0xF800
	px |= (uint16(g) >> 5) & 0x07E0
	px |= (uint16(b) >> 11) & 0x001F
	le := make([]byte, 2)
	binary.LittleEndian.PutUint16(le, px)
	return fmt.Sprintf("0x%02x", le)
}

func pixelToRGB(pixels []color.Color) string {
	r, g, b, _ := pixels[0].RGBA()
	return fmt.Sprintf("0x%02x%02x%02x", r>>8, g>>8, b>>8)
}

func pixelToRGBA(pixels []color.Color) string {
	r, g, b, a := pixels[0].RGBA()

	return fmt.Sprintf("0x%02x%02x%02x%02x", r>>8, g>>8, b>>8, a>>8)
}
