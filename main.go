package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	input := flag.String("in", "ttyrecord", "input ttyrec file")
	output := flag.String("out", "tty.gif", "output gif file")
	speed := flag.Float64("s", 1.0, "play speed")
	row := flag.Int("row", 24, "rows")
	col := flag.Int("col", 80, "columns")
	noloop := flag.Bool("noloop", false, "play only once")
	help := flag.Bool("help", false, "usage")
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	generator := NewGifGenerator()
	generator.Speed = *speed
	generator.Row = *row
	generator.Col = *col
	generator.NoLoop = *noloop
	err := generator.Generate(*input, *output)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	absPath, err := filepath.Abs(*output)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("%s created!\n", absPath)
}
