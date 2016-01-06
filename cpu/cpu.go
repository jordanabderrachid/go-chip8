package cpu

import (
	"github.com/jordanabderrachid/go-chip8/mmu"
	"github.com/jordanabderrachid/go-chip8/timer"
)

type Registers struct {
	V [16]byte // The last byte VF is the flag register
	I rune

	PC rune // program counter
	SP byte // stack pointer

	DT byte // delay timer
	ST byte // sound timer
}

func (r Registers) Reset() {
	for i := range r.V {
		r.V[i] = 0
	}

	r.I = 0
	r.PC = 0
	r.SP = 0
	r.DT = 0
	r.ST = 0
}

type CPU struct {
	R                      Registers
	Memory                 *mmu.Memory
	SoundTimer, DelayTimer timer.Timer
}

func (cpu *CPU) Reset() {
	cpu.Memory = new(mmu.Memory)

	cpu.R.Reset()
	cpu.Memory.Reset()

	go cpu.SoundTimer.Run(&cpu.R.ST)
	go cpu.DelayTimer.Run(&cpu.R.DT)
}
