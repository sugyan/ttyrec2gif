package main

import (
	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
	"fmt"
	"github.com/errnoh/term.color"
	"image"
	"image/color/palette"
	"image/draw"
	"j4k.co/terminal"
	"os"
)

var font *truetype.Font

const fontSize = 18

func init() {
	fontData, err := Asset("font/Anonymous Pro Minus.ttf")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	font, err = freetype.ParseFont(fontData)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Capture draws virtual terminal and return paletted image
func (g *GifGenerator) Capture(state *terminal.State) (paletted *image.Paletted, err error) {
	fb := font.Bounds(fontSize)
	cursorX, cursorY := state.Cursor()
	paletted = image.NewPaletted(image.Rect(0, 0, g.Col*int(fb.XMax-fb.XMin)+10, g.Row*int(fb.YMax-fb.YMin)), palette.WebSafe)

	c := freetype.NewContext()
	c.SetFontSize(fontSize)
	c.SetFont(font)
	c.SetDst(paletted)
	c.SetClip(paletted.Bounds())
	for row := 0; row < g.Row; row++ {
		for col := 0; col < g.Col; col++ {
			ch, fg, bg := state.Cell(col, row)
			var uniform *image.Uniform
			// background color
			if bg != terminal.DefaultBG {
				if bg == terminal.DefaultFG {
					uniform = image.White
				} else {
					uniform = image.NewUniform(color.Term256{Val: uint8(bg)})
				}
			}
			// cursor
			if state.CursorVisible() && (row == cursorY && col == cursorX) {
				uniform = image.White
			}
			if uniform != nil {
				draw.Draw(paletted, image.Rect(5+col*int(fb.XMax-fb.XMin), row*int(fb.YMax-fb.YMin)-int(fb.YMin), 5+(col+1)*int(fb.XMax-fb.XMin), (row+1)*int(fb.YMax-fb.YMin)-int(fb.YMin)), uniform, image.ZP, draw.Src)
			}
			// foreground color
			switch fg {
			case terminal.DefaultFG:
				c.SetSrc(image.White)
			case terminal.DefaultBG:
				c.SetSrc(image.Black)
			default:
				c.SetSrc(image.NewUniform(color.Term256{Val: uint8(fg)}))
			}
			_, err = c.DrawString(string(ch), freetype.Pt(5+col*int(fb.XMax-fb.XMin), (row+1)*int(fb.YMax-fb.YMin)))
			if err != nil {
				return
			}
		}
	}
	return paletted, nil
}
