package main

import (
	"flag"
	"github.com/jordanabderrachid/go-chip8/cpu"
	"github.com/nsf/termbox-go"
	"log"
	"os"
)

func main() {
	logfile, err := os.Create("log")
	if err != nil {
		panic(err)
	}
	log.SetOutput(logfile)

	romFile := flag.String("r", "", "rom file")
	flag.Parse()

	log.Println("initializing termbox")
	if err := termbox.Init(); err != nil {
		log.Panicln(err)
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
