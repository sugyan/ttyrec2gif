package main

import (
	"github.com/sugyan/ttyread"
	"io"
	"os"
)

// Terminal interface
type Terminal interface {
	Write([]byte) (int, error)
}

// TtyPlay reads ttyrecord file and play
func (g *GifGenerator) TtyPlay(input string, term Terminal, capture func(int32) error) (err error) {
	inFile, err := os.Open(input)
	if err != nil {
		return
	}
	defer inFile.Close()

	reader := ttyread.NewTtyReader(inFile)
	var prevTv *ttyread.TimeVal
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

		err = capture(diff.Sec*1000000 + diff.Usec)
		if err != nil {
			return
		}

		_, err = term.Write(*data.Buffer)
		if err != nil {
			return
		}
	}
	return nil
}
