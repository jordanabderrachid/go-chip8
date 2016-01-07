package mmu

import (
	"bytes"
	"testing"
)

func TestReset(t *testing.T) {
	mem := new(Memory)
	mem.Reset()

	for i, b := range mem.m {
		if b != 0 {
			t.Errorf("Memory not reseted at %04x got %02x", i, b)
			break
		}
	}

	if l := len(mem.m); l != int(memorySize) {
		t.Errorf("Memory has not the proper size, got %d", l)
	}
}

func TestLoadSprites(t *testing.T) {
	mem := new(Memory)
	mem.Reset()

	if err := mem.LoadSprites(); err != nil {
		t.Errorf("Failed loading sprites %s", err)
	}
}

func TestAllocate(t *testing.T) {
	mem := new(Memory)
	mem.Reset()

	buffer := make([]byte, 10)
	for i := range buffer {
		buffer[i] = byte(i)
	}

	mem.Allocate(buffer)

	if ok := bytes.Equal(buffer, mem.m[:len(buffer)]); !ok {
		t.Error("Allocate failed")
	}
}

func TestGetByte(t *testing.T) {
	mem := new(Memory)
	mem.Reset()

	buffer := make([]byte, 10)
	for i := range buffer {
		buffer[i] = byte(i)
	}

	mem.Allocate(buffer)

	b, err := mem.GetByte(rune(3))
	if err != nil {
		t.Error(err)
	}

	if b != buffer[3] {
		t.Errorf("Expected %02x Got %02x", 3, b)
	}

	if _, err := mem.GetByte(memorySize + 1); err == nil {
		t.Errorf("Expected error with call address %x", memorySize+1)
	}

	if _, err := mem.GetByte(-1); err == nil {
		t.Error("Expected error with negative call address")
	}
}

func TestSetByte(t *testing.T) {
	mem := new(Memory)
	mem.Reset()

	if err := mem.SetByte(-1, 0); err == nil {
		t.Error("Expected error with negative call address")
	}

	if err := mem.SetByte(memorySize+1, 0); err == nil {
		t.Errorf("Expected error with call address %x", memorySize+1)
	}

	var addr rune = 42
	var b byte = 16
	if err := mem.SetByte(addr, b); err != nil {
		t.Errorf("Got an error while setting byte %s", err)
	}

	if res, _ := mem.GetByte(addr); b != res {
		t.Errorf("Error setting byte, expected %x, got %x, at %04x", b, res, addr)
	}
}
