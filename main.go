package main

import (
	"flag"
	"github.com/jordanabderrachid/go-chip8/cpu"
	"github.com/nsf/termbox-go"
	"os"
)

func main() {
	romFile := flag.String("r", "", "rom file")
	flag.Parse()

	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	CPU := new(cpu.CPU)
	CPU.Reset()

	f, _ := os.Open(*romFile)
	b := make([]byte, 3584)
	f.Read(b)
	CPU.LoadData(b)
	CPU.Run()
}
