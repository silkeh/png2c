// Copyright 2018 Silke Hofstra
//
// Licensed under the EUPL
//

// png2c is a simple tool for converting PNG images to a C definition.
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func imageToHex(img image.Image, mode Mode) (hex [][]string) {
	bounds := img.Bounds()
	hex = make([][]string, bounds.Max.Y)

	for y := 0; y < bounds.Max.Y; y++ {
		hex[y] = make([]string, bounds.Max.X/mode.Density)
		for x := 0; x < bounds.Max.X; x += mode.Density {
			pixels := make([]color.Color, mode.Density)
			for i := range pixels {
				pixels[i] = img.At(x+i, y)
			}
			hex[y][x/mode.Density] = mode.Convert(pixels)
		}
	}

	return
}

func hexToString(hex [][]string) string {
	rows := make([]string, len(hex))
	for i, r := range hex {
		rows[i] = fmt.Sprintf("{ %s }", strings.Join(r, ", "))
	}
	return strings.Join(rows, ",\n    ")
}

func main() {
	var fileName, varName, modeName, brief string
	var listModes bool
	flag.StringVar(&fileName, "file", "", "Image to convert")
	flag.StringVar(&varName, "var", "picture", "Name of the variable")
	flag.StringVar(&modeName, "mode", "bw", "Conversion mode")
	flag.StringVar(&brief, "brief", "", "Doxygen brief")
	flag.BoolVar(&listModes, "list-modes", false, "List available modes")
	flag.Parse()

	if listModes {
		for m := range modes {
			fmt.Printf("%s\n", m)
		}
		return
	}

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}

	img, err := png.Decode(file)
	if err != nil {
		log.Fatalf("Error decoding PNG: %s", err)
	}

	mode, ok := modes[modeName]
	if !ok {
		log.Fatalf("Unknown mode: %s", modeName)
	}

	data := imageToHex(img, mode)

	if brief != "" {
		fmt.Printf("/**\n * @brief   %s (%vx%v)\n */\n",
			brief, len(data[0])*mode.Density, len(data))
	}

	fmt.Printf("const uint%v_t %s[%v][%v] = {\n    ",
		mode.Bits,
		varName,
		len(data),
		len(data[0]),
	)

	fmt.Printf("%s\n};\n\n", hexToString(data))
}
