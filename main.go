package main

import (
	"flag"
	"github.com/jordanabderrachid/go-chip8/cpu"
	"log"
	"os"
	"runtime"
)

func main() {
	runtime.LockOSThread()
	logfile, err := os.Create("log")
	if err != nil {
		panic(err)
	}
	log.SetOutput(logfile)

	romFile := flag.String("r", "", "rom file")
	flag.Parse()

	CPU := new(cpu.CPU)
	CPU.Reset()

	f, _ := os.Open(*romFile)
	b := make([]byte, 3584)
	f.Read(b)
	CPU.LoadData(b)
	CPU.Run()
}
