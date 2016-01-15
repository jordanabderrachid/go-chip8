package mmu

import (
	"fmt"
	"github.com/jordanabderrachid/go-chip8/display"
	"log"
)

const memorySize rune = 0x1000 // 4096

type Memory struct {
	m [memorySize]byte // 4096 bytes
}

func (mem *Memory) Reset() {
	log.Println("reseting memory")
	for i := range mem.m {
		mem.SetByte(rune(i), 0x00)
	}
}

func (mem *Memory) LoadSprites() error {
	log.Println("loading sprites")
	var addr rune = 0x0000
	for i := 0; i <= 0x0F; i++ {
		s := display.Sprites[byte(i)]
		for _, b := range s {
			if err := mem.SetByte(addr, b); err != nil {
				return err
			}
			addr++
		}
	}

	return nil
}

func (mem *Memory) Allocate(buffer []byte) error {
	for i := 0; i < len(buffer); i++ {
		if err := mem.SetByte(rune(i), buffer[i]); err != nil {
			return err
		}
	}

	return nil
}

func (mem *Memory) AllocateWithBuffer(buffer []byte, offset rune) error {
	for i := 0; i < len(buffer); i++ {
		if err := mem.SetByte(offset+rune(i), buffer[i]); err != nil {
			return err
		}
	}

	return nil
}

func (mem *Memory) GetByte(addr rune) (byte, error) {
	if addr > memorySize || addr < 0 {
		return 0, fmt.Errorf("Illegal address %04x\n", addr)
	}

	return mem.m[addr], nil
}

func (mem *Memory) SetByte(addr rune, b byte) error {
	if addr > memorySize || addr < 0 {
		return fmt.Errorf("Illegal address %04x\n", addr)
	}

	mem.m[addr] = b
	return nil
}
