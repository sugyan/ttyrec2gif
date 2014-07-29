package main

import (
	"bytes"
	"github.com/sugyan/ttyread"
	"image"
	"image/gif"
	"io"
	"io/ioutil"
	"j4k.co/terminal"
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
	// input
	inFile, err := os.Open(input)
	if err != nil {
		return
	}
	defer inFile.Close()

	// virtual terminal
	var state = terminal.State{}
	vt, err := terminal.Create(&state, ioutil.NopCloser(bytes.NewBuffer([]byte{})))
	if err != nil {
		return
	}
	defer vt.Close()
	vt.Resize(g.Col, g.Row)

	// read ttyrecord
	reader := ttyread.NewTtyReader(inFile)
	var (
		prevTv *ttyread.TimeVal
		images []*image.Paletted
		delays []int
	)
	for {
		var data *ttyread.TtyData
		data, err = reader.ReadData()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return
			}
		}
		var diff ttyread.TimeVal
		if prevTv != nil {
			diff = data.TimeVal.Subtract(*prevTv)
		}
		prevTv = &data.TimeVal

		// calc delay and capture
		delay := int(float64(diff.Sec*1000000+diff.Usec)/g.Speed) / 10000
		if delay > 0 {
			var img *image.Paletted
			img, err = g.Capture(&state)
			if err != nil {
				return
			}
			images = append(images, img)
			delays = append(delays, delay)
		}
		// write to vt
		_, err = vt.Write(*data.Buffer)
		if err != nil {
			return
		}
	}

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
