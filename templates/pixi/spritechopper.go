package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"path/filepath"
)

type FrameData struct {
	Frame            Rect `json:"frame"`
	Rotated          bool `json:"rotated"`
	Trimmed          bool `json:"trimmed"`
	SpriteSourceSize Rect `json:"spriteSourceSize"`
	SourceSize       Size `json:"sourceSize"`
}

type Rect struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type Size struct {
	W int `json:"w"`
	H int `json:"h"`
}

type Meta struct {
	Image string `json:"image"`
	Size  Size   `json:"size"`
	Scale string `json:"scale"`
}

type SpriteSheet struct {
	Frames map[string]FrameData `json:"frames"`
	Meta   Meta                 `json:"meta"`
}

func main() {
	filePath := flag.String("file", "", "Path to the PNG tileset")
	tileW := flag.Int("tilewidth", 0, "Width of each tile")
	tileH := flag.Int("tileheight", 0, "Height of each tile")
	output := flag.String("output", "spritesheet.json", "Output JSON filename")
	flag.Parse()

	if *filePath == "" || *tileW <= 0 || *tileH <= 0 {
		fmt.Println("Usage: -file tiles.png -tilewidth 64 -tileheight 64 [-output out.json]")
		return
	}

	f, err := os.Open(*filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.DecodeConfig(f)
	if err != nil {
		panic(err)
	}

	imgW, imgH := img.Width, img.Height
	cols := imgW / *tileW
	rows := imgH / *tileH

	sheet := SpriteSheet{
		Frames: map[string]FrameData{},
		Meta: Meta{
			Image: filepath.Base(*filePath),
			Size:  Size{W: imgW, H: imgH},
			Scale: "1",
		},
	}

	count := 0
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			name := fmt.Sprintf("tile_%d.png", count)
			frame := FrameData{
				Frame: Rect{
					X: x * (*tileW),
					Y: y * (*tileH),
					W: *tileW,
					H: *tileH,
				},
				Rotated:          false,
				Trimmed:          false,
				SpriteSourceSize: Rect{X: 0, Y: 0, W: *tileW, H: *tileH},
				SourceSize:       Size{W: *tileW, H: *tileH},
			}
			sheet.Frames[name] = frame
			count++
		}
	}

	jsonData, err := json.MarshalIndent(sheet, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(*output, jsonData, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Spritesheet JSON written to %s (%d tiles)\n", *output, count)
}
