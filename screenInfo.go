package main

import (
	"github.com/james4k/terminal"
)

// ScreenInfo type
type ScreenInfo struct {
	Width         int
	Height        int
	Chars         []rune
	Fcolors       []terminal.Color
	Bcolors       []terminal.Color
	CursorX       int
	CursorY       int
	CursorVisible bool
	top           int
	bottom        int
	left          int
	right         int
}

// NewScreenInfo returns ScreenInfo instance
func NewScreenInfo() *ScreenInfo {
	return &ScreenInfo{
		Width:         -1,
		Height:        -1,
		Chars:         []rune{},
		Fcolors:       []terminal.Color{},
		Bcolors:       []terminal.Color{},
		CursorX:       -1,
		CursorY:       -1,
		CursorVisible: false,
		top:           -1,
		bottom:        -1,
		left:          -1,
		right:         -1,
	}
}

func (s *ScreenInfo) save(width int, height int, state *terminal.State) {
	if s.Width != width || s.Height != height {
		s.Width = width
		s.Height = height
		s.Chars = make([]rune, width*height)
		s.Fcolors = make([]terminal.Color, width*height)
		s.Bcolors = make([]terminal.Color, width*height)
	}
	for row := 0; row < s.Height; row++ {
		for col := 0; col < s.Width; col++ {
			ch, fg, bg := state.Cell(col, row)
			s.Chars[row*s.Width+col] = ch
			s.Fcolors[row*s.Width+col] = fg
			s.Bcolors[row*s.Width+col] = bg
		}
	}
	s.CursorX, s.CursorY = state.Cursor()
	s.CursorVisible = state.CursorVisible()
}

func (s *ScreenInfo) updateRedrawRange(x int, y int) {
	if y < s.top {
		s.top = y
	}
	if s.bottom < y+1 {
		s.bottom = y + 1
	}
	if x < s.left {
		s.left = x
	}
	if s.right < x+1 {
		s.right = x + 1
	}
}

// GetRedrawRange returns redraw range.
func (s *ScreenInfo) GetRedrawRange(width int, height int, state *terminal.State) (left int, top int, right int, bottom int) {

	defer s.save(width, height, state)

	if s.Width != width || s.Height != height {
		return 0, 0, width, height
	}

	s.top = height
	s.bottom = 0
	s.left = width
	s.right = 0

	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			ch, fg, bg := state.Cell(col, row)
			ch0 := s.Chars[row*width+col]
			fg0 := s.Fcolors[row*width+col]
			bg0 := s.Bcolors[row*width+col]
			if ch != ch0 || fg != fg0 || bg != bg0 {
				s.updateRedrawRange(col, row)
			}
		}
	}
	cursorVisible := state.CursorVisible()
	cursorX, cursorY := state.Cursor()

	if s.CursorVisible && !cursorVisible {
		s.updateRedrawRange(s.CursorX, s.CursorY)
	}

	if !s.CursorVisible && cursorVisible {
		s.updateRedrawRange(cursorX, cursorY)
	}

	if s.CursorVisible && cursorVisible && (s.CursorX != cursorX || s.CursorY != cursorY) {
		s.updateRedrawRange(s.CursorX, s.CursorY)
		s.updateRedrawRange(cursorX, cursorY)
	}

	return s.left, s.top, s.right, s.bottom
}
