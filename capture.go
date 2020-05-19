package main

import (
	"fmt"
	"github.com/errnoh/term.color"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/james4k/terminal"
	"image"
	"image/color/palette"
	"image/draw"
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
	left, top, right, bottom := g.ScreenInfo.GetRedrawRange(g.Col, g.Row, state)
	paletted = image.NewPaletted(image.Rect(left*int(fb.Max.X-fb.Min.X), top*int(fb.Max.Y-fb.Min.Y), right*int(fb.Max.X-fb.Min.X)+10, bottom*int(fb.Max.Y-fb.Min.Y)+10), palette.WebSafe)

	c := freetype.NewContext()
	c.SetFontSize(fontSize)
	c.SetFont(font)
	c.SetDst(paletted)
	c.SetClip(image.Rect(0, 0, g.Col*int(fb.Max.X-fb.Min.X)+10, g.Row*int(fb.Max.Y-fb.Min.Y)+10))
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
				draw.Draw(paletted, image.Rect(5+col*int(fb.Max.X-fb.Min.X), row*int(fb.Max.Y-fb.Min.Y)-int(fb.Min.Y), 5+(col+1)*int(fb.Max.X-fb.Min.X), (row+1)*int(fb.Max.Y-fb.Min.Y)-int(fb.Min.Y)), uniform, image.ZP, draw.Src)
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
			_, err = c.DrawString(string(ch), freetype.Pt(5+col*int(fb.Max.X-fb.Min.X), (row+1)*int(fb.Max.Y-fb.Min.Y)))
			if err != nil {
				return
			}
		}
	}
	return paletted, nil
}
