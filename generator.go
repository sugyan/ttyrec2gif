package main

import (
	"bytes"
	"image"
	"image/gif"
	"io/ioutil"
	"github.com/james4k/terminal"
	"os"
)

// GifGenerator type
type GifGenerator struct {
	Speed  float64
	Col    int
	Row    int
	NoLoop bool
}

// NewGifGenerator returns GifGenerator instance
func NewGifGenerator() *GifGenerator {
	return &GifGenerator{
		Speed:  1.0,
		Col:    80,
		Row:    24,
		NoLoop: false,
	}
}

// Generate writes to outFile an animated GIF
func (g *GifGenerator) Generate(input string, output string) (err error) {
	// virtual terminal
	var state = terminal.State{}
	vt, err := terminal.Create(&state, ioutil.NopCloser(bytes.NewBuffer([]byte{})))
	if err != nil {
		return
	}
	defer vt.Close()
	vt.Resize(g.Col, g.Row)

	// play and capture
	var (
		images []*image.Paletted
		delays []int
	)
	err = g.TtyPlay(input, vt, func(diff int32) (err error) {
		delay := int(float64(diff)/g.Speed) / 10000
		if delay > 0 {
			var img *image.Paletted
			img, err = g.Capture(&state)
			if err != nil {
				return
			}
			images = append(images, img)
			delays = append(delays, delay)
		}
		return nil
	})
	if err != nil {
		return
	}

	// generate gif file
	outFile, err := os.Create(output)
	if err != nil {
		return
	}
	defer outFile.Close()
	opts := gif.GIF{
		Image: images,
		Delay: delays,
	}
	if g.NoLoop {
		opts.LoopCount = 1
	}
	err = gif.EncodeAll(outFile, &opts)
	if err != nil {
		return
	}
	return nil
}
