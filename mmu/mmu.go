package mmu

import (
	"fmt"
	"github.com/jordanabderrachid/go-chip8/display"
)

const memorySize rune = 4096

type Memory struct {
	m [memorySize]byte
}

func (mem *Memory) Reset() {
	for i := range mem.m {
		mem.SetByte(rune(i), 0)
	}
}

func (mem *Memory) LoadSprites() error {
	var addr rune = 0
	for _, v := range display.Sprites {
		for _, b := range v {
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
