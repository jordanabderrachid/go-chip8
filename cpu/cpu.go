package cpu

import (
	"fmt"
	"github.com/jordanabderrachid/go-chip8/display"
	"github.com/jordanabderrachid/go-chip8/keyboard"
	"github.com/jordanabderrachid/go-chip8/mmu"
	"github.com/jordanabderrachid/go-chip8/timer"
	"math/rand"
	"time"
)

type Registers struct {
	V [16]byte // The last byte VF is the flag register
	I rune

	PC rune // program counter

	Stack [16]rune // stack
	SP    byte     // stack pointer

	DT byte // dhigelay timer
	ST byte // sound timer
}

func (r *Registers) Reset() {
	for i := range r.V {
		r.V[i] = 0x00
	}
	r.I = 0x0000

	r.PC = 0x200 // program start at address 0x200

	for i := range r.Stack {
		r.Stack[i] = 0x0000
	}
	r.SP = 0x00

	r.DT = 0x00
	r.ST = 0x00
}

type CPU struct {
	R                      *Registers
	Memory                 *mmu.Memory
	SoundTimer, DelayTimer timer.Timer
	Display                *display.Display
	Keyboard               *keyboard.Keyboard
}

func (cpu *CPU) Reset() {
	cpu.Memory = new(mmu.Memory)
	cpu.Keyboard = new(keyboard.Keyboard)
	cpu.Display = new(display.Display)
	cpu.R = new(Registers)

	cpu.R.Reset()
	cpu.Memory.Reset()
	cpu.Display.Reset()
	cpu.Keyboard.Reset()

	go cpu.SoundTimer.Run(&cpu.R.ST)
	go cpu.DelayTimer.Run(&cpu.R.DT)
}

func (cpu *CPU) LoadData(b []byte) {
	if err := cpu.Memory.AllocateWithBuffer(b, 0x200); err != nil {
		panic(err)
	}
}

func (cpu *CPU) GetOpcode(addr rune) (opcode rune) {
	var high byte
	var low byte
	var err error
	// instructions are stored as big-endian
	high, err = cpu.Memory.GetByte(addr)
	if err != nil {
		panic(err)
	}

	low, err = cpu.Memory.GetByte(addr + 1)
	if err != nil {
		panic(err)
	}

	fmt.Printf("opcode at %04x is %04x\n", addr, rune(high)<<8+rune(low))
	return rune(high)<<8 + rune(low)
}

func (cpu *CPU) Run() {
	ticker := time.NewTicker(time.Duration(int64(time.Second) / timer.Frenquency))

	for {
		select {
		case <-ticker.C:
			cpu.ExecuteOpcode(cpu.GetOpcode(cpu.R.PC))
		}
	}
}

func (cpu *CPU) ExecuteOpcode(opcode rune) {
	switch opcode & 0xF000 {
	case 0x0000: // 0x0xxx
		switch opcode {
		case 0x00E0: // 0x00E0
			cpu.instr_00E0()
		case 0x00EE: // 0x00EE
			cpu.instr_00EE()
		default:
			panic(fmt.Sprintf("Unknown opcode %04x", opcode))
		}
	case 0x1000: // 0x1xxx
		addr := opcode & 0x0FFF
		cpu.instr_1nnn(addr)
	case 0x2000: // 0x2xxx
		addr := opcode & 0x0FFF
		cpu.instr_2nnn(addr)
	case 0x3000: // 0x3xxx
		x := byte(opcode & 0x0F00)
		value := byte(opcode & 0x00FF)
		cpu.instr_3xkk(x, value)
	case 0x4000: // 0x4xxx
		x := byte(opcode & 0x0F00)
		value := byte(opcode & 0x00FF)
		cpu.instr_4xkk(x, value)
	case 0x5000: // 0x5xxx
		x := byte(opcode & 0x0F00)
		y := byte(opcode & 0x00F0)
		switch opcode & 0x000F {
		case 0x0000: // 0x5xx0
			cpu.instr_5xy0(x, y)
		default:
			panic(fmt.Sprintf("Unknown opcode %04x", opcode))
		}
	case 0x6000: // 0x6xxx
		x := byte(opcode & 0x0F00)
		value := byte(opcode & 0x00FF)
		cpu.instr_6xkk(x, value)
	case 0x7000: // 0x7xxx
		x := byte(opcode & 0x0F00)
		value := byte(opcode & 0x00FF)
		cpu.instr_7xkk(x, value)
	case 0x8000: // 0x8xxx
		x := byte(opcode & 0x0F00)
		y := byte(opcode & 0x00F0)
		switch opcode & 0x000F {
		case 0x0000: // 0x8xx0
			cpu.instr_8xy0(x, y)
		case 0x0001: // 0x8xx1
			cpu.instr_8xy1(x, y)
		case 0x0002: // 0x8xx2
			cpu.instr_8xy2(x, y)
		case 0x0003: // 0x8xx3
			cpu.instr_8xy3(x, y)
		case 0x0004: // 0x8xx4
			cpu.instr_8xy4(x, y)
		case 0x0005: // 0x8xx5
			cpu.instr_8xy5(x, y)
		case 0x0006: // 0x8xx6
			cpu.instr_8xy6(x)
		case 0x0007: // 0x8xx7
			cpu.instr_8xy7(x, y)
		case 0x000E: // 0x8xxE
			cpu.instr_8xyE(x)
		default:
			panic(fmt.Sprintf("Unknown opcode %04x", opcode))
		}
	case 0x9000: // 0x9xxx
		x := byte(opcode & 0x0F00)
		y := byte(opcode & 0x00F0)
		switch opcode & 0x000F {
		case 0x0000: // 0x9xx0
			cpu.instr_9xy0(x, y)
		default:
			panic(fmt.Sprintf("Unknown opcode %04x", opcode))
		}
	case 0xA000: // 0xAxxx
		addr := opcode & 0x0FFF
		cpu.instr_Annn(addr)
	case 0xB000: // 0xBxxx
		addr := opcode & 0x0FFF
		cpu.instr_Bnnn(addr)
	case 0xC000: // 0xCxxx
		x := byte(opcode & 0x0F00)
		value := byte(opcode & 0x00FF)
		cpu.instr_Cxkk(x, value)
	case 0xD000: // 0xDxxx
		x := byte(opcode & 0x0F00)
		y := byte(opcode & 0x00F0)
		n := byte(opcode & 0x000F)
		cpu.instr_Dxyn(x, y, n)
	case 0xE000: // 0xExxx
		x := byte(opcode & 0x0F00)
		switch opcode & 0x00FF {
		case 0x009E: // 0xEx9E
			cpu.instr_Ex9E(x)
		case 0x00A1: // 0xExA1
			cpu.instr_ExA1(x)
		default:
			panic(fmt.Sprintf("Unknown opcode %04x", opcode))
		}
	case 0xF000: // 0xF000
		x := byte(opcode & 0x0F00)
		switch opcode & 0x00FF {
		case 0x0007: // 0xFx07
			cpu.instr_Fx07(x)
		case 0x000A: // 0xFx0A
			cpu.instr_Fx0A(x)
		case 0x0015: // 0xFx15
			cpu.instr_Fx15(x)
		case 0x0018: // 0xFx18
			cpu.instr_Fx18(x)
		case 0x001E: // 0xFx1E
			cpu.instr_Fx1E(x)
		case 0x0029: // 0xFx29
			cpu.instr_Fx29(x)
		case 0x0033: // 0xFx33
			cpu.instr_Fx33(x)
		case 0x0055: // 0xFx55
			cpu.instr_Fx55(x)
		case 0x0065: // 0xFx65
			cpu.instr_Fx65(x)
		default:
			panic(fmt.Sprintf("Unknown opcode %04x", opcode))
		}
	}
}

// 0x00E0 - CLS
// Clear the display.
// Increment the PC.
func (cpu *CPU) instr_00E0() {
	cpu.Display.Clear()
	cpu.R.PC += 2
}

// 0x00EE - RET
// Return from a subroutine.
//
// The interpreter sets the program counter to the address at the top of the stack, then substracts 1 from the stack pointer.
func (cpu *CPU) instr_00EE() {
	cpu.R.PC = cpu.R.Stack[cpu.R.SP]
	cpu.R.SP--
}

// 0x1nnn - JP addr
// Jump to location nnn.
//
// The interpreter sets the program counter to nnn.
func (cpu *CPU) instr_1nnn(addr rune) {
	cpu.R.PC = addr
}

// 0x2nnn - CALL addr
// Call subroutine at nnn.
//
// The interpreter increments the stack pointer, then puts the current PC on the top of the stack. The PC is then set to nnn.
func (cpu *CPU) instr_2nnn(addr rune) {
	cpu.R.SP++
	cpu.R.Stack[cpu.R.SP] = cpu.R.PC
	cpu.R.PC = addr
}

// 0x3xkk - SE Vx, byte
// Skip next instruction if Vx == kk.
//
// The interpreter compares register Vx to kk, and if they are equal, increments the program counter by 4, else increments by 2.
func (cpu *CPU) instr_3xkk(x, value byte) {
	if cpu.R.V[x] == value {
		cpu.R.PC += 4
	} else {
		cpu.R.PC += 2
	}
}

// 0x4xkk - SNE Vx, byte
// Skip next instruction if Vx != kk.
//
// The interpreter compares register Vx to kk, and if they are not equal, increments the program counter by 4, else increments by 2.
func (cpu *CPU) instr_4xkk(x, value byte) {
	if cpu.R.V[x] != value {
		cpu.R.PC += 4
	} else {
		cpu.R.PC += 2
	}
}

// 0x5xy0 - SE Vx, Vy
// Skip next instruction if Vx == Vy.
//
// The interpreter compares register Vx to register Vy, and if they are equal, increments the program counter by 4, else increments by 2.
func (cpu *CPU) instr_5xy0(x, y byte) {
	if cpu.R.V[x] == cpu.R.V[y] {
		cpu.R.PC += 4
	} else {
		cpu.R.PC += 2
	}
}

// 0x6xkk - LD Vx, byte
// Set Vx = kk.
//
// The interpreter puts the value kk into register Vx.
func (cpu *CPU) instr_6xkk(x, value byte) {
	cpu.R.V[x] = value
	cpu.R.PC += 2
}

// 0x7xkk - ADD Vx, byte
// Set Vx = Vx + kk.
//
// Adds the value kk to the value of register Vx, then stores the result in Vx.
func (cpu *CPU) instr_7xkk(x, value byte) {
	cpu.R.V[x] += value
	cpu.R.PC += 2
}

// 0x8xy0 - LD Vx, Vy
// Set Vx = Vy.
//
// Stores the value of the register Vy in register Vx.
func (cpu *CPU) instr_8xy0(x, y byte) {
	cpu.R.V[x] = cpu.R.V[y]
	cpu.R.PC += 2
}

// 0x8xy1 - OR Vx, Vy
// Set Vx = Vx OR Vy.
//
// Performs a bitwise OR on the values of Vx and Vy, then stores the result in Vx.
func (cpu *CPU) instr_8xy1(x, y byte) {
	cpu.R.V[x] = cpu.R.V[x] | cpu.R.V[y]
	cpu.R.PC += 2
}

// 0x8xy2 - AND Vx, Vy
// Set Vx = Vx AND Vy.
//
// Performs a bitwise AND on the values of Vx and Vy, then stores the result in Vx.
func (cpu *CPU) instr_8xy2(x, y byte) {
	cpu.R.V[x] = cpu.R.V[x] & cpu.R.V[y]
	cpu.R.PC += 2
}

// 0x8xy3 - XOR Vx, Vy
// Set Vx = Vx XOR Vy.
//
// Performs a bitwise exclusive OR on the values of Vx and Vy, then stores the result in Vx.
func (cpu *CPU) instr_8xy3(x, y byte) {
	cpu.R.V[x] = cpu.R.V[x] ^ cpu.R.V[y]
	cpu.R.PC += 2
}

// 0x8xy4 - ADD Vx, Vy
// Set Vx = Vx + Vy, set VF = carry.
//
// The values of Vx and Vy are added together. If the result is greater than 8bits (>255), VF is set to 1, otherwise 0.
// Only the lowest 8 bits of the result are kept, and stored in Vx.
func (cpu *CPU) instr_8xy4(x, y byte) {
	result := rune(cpu.R.V[x] + cpu.R.V[y])

	if result > 0xFF {
		cpu.R.V[0xF] = 1
	} else {
		cpu.R.V[0xF] = 0
	}

	cpu.R.V[x] = byte(result & 0xFF)
	cpu.R.PC += 2
}

// 0x8xy5 - SUB Vx, Vy
// Set Vx = Vx - Vy, set VF = NOT borrow.
//
// If Vx > Vy, then VF is set to 1, otherwise 0. The Vy is subtracted from Vx, and the result is stored in Vx.
func (cpu *CPU) instr_8xy5(x, y byte) {
	if cpu.R.V[x] > cpu.R.V[y] {
		cpu.R.V[0xF] = 1
	} else {
		cpu.R.V[0xF] = 0
	}

	cpu.R.V[x] -= cpu.R.V[y]
	cpu.R.PC += 2
}

// 0x8xy6 - SHR Vx {, Vy}
// Set Vx = Vx SHR 1.
//
// If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2.
func (cpu *CPU) instr_8xy6(x byte) {
	cpu.R.V[0xF] = cpu.R.V[x] & 0x1
	cpu.R.V[x] /= 2
	cpu.R.PC += 2
}

// 0x8xy7 - SUBN Vx, Vy
// Set Vx = Vy - Vy, set VF = NOT borrow.
//
// If Vy > Vx, then VF is set to 1, otherwise 0. Then Vx is subtracted from Vy, and the result is stored in Vx.
func (cpu *CPU) instr_8xy7(x, y byte) {
	if cpu.R.V[y] > cpu.R.V[x] {
		cpu.R.V[0xF] = 1
	} else {
		cpu.R.V[0xF] = 0
	}

	cpu.R.V[x] = cpu.R.V[y] - cpu.R.V[x]
	cpu.R.PC += 2
}

// 0x8xyE - SHL Vx, {, Vy}
// Set Vx = Vx SHL 1.
//
// If the most-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is multiplied by 2.
func (cpu *CPU) instr_8xyE(x byte) {
	cpu.R.V[0xF] = cpu.R.V[x] & 0x80
	cpu.R.V[x] *= 2
	cpu.R.PC += 2
}

// 0x9xy0 - SNE Vx, Vy
// Skip next instruction if Vx != Vy.
//
// The values of Vx and Vy are compared, and if they are not equal, the program counter is increased by 4, otherwise by 2.
func (cpu *CPU) instr_9xy0(x, y byte) {
	if cpu.R.V[x] != cpu.R.V[y] {
		cpu.R.PC += 4
	} else {
		cpu.R.PC += 2
	}
}

// 0xAnnn - LD I, addr
// Set I = addr.
//
// The value of register I is set to nnn.
func (cpu *CPU) instr_Annn(addr rune) {
	cpu.R.I = addr
	cpu.R.PC += 2
}

// 0xBnnn - JP V0, addr
// Jump to location nnn + V0.
//
// The program counter is set to nnn plus the value of V0.
func (cpu *CPU) instr_Bnnn(addr rune) {
	cpu.R.PC = rune(cpu.R.V[0x0]) + addr
}

// 0xCxkk - RND Vx, byte
// Set Vx = random byte AND kk.
//
// The interpreter generates a random number from 0 to 255, which is then ANDed with the value kk. The results are stored in Vx.
func (cpu *CPU) instr_Cxkk(x, value byte) {
	random := byte(rand.Intn(0xFF + 1)) // Intn exclude the last value.
	cpu.R.V[x] = random & value
	cpu.R.PC += 2
}

// 0xDxyn - DRW Vx, Vy, nibble
// Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
//
// The interpreter reads n bytes from memory, starting at the address stored in I. These bytes are the displayed as sprites on screen
// coordinates (Vx, Vy). Sprites are XORed onto the existing screen. If this causes any pixel to be erased, VF is set to 1, otherwise
// it is set to 0. If the sprite is positioned so part of it is outside the coordinates of the display, it wraps around to the
// opposite side of the screen.
func (cpu *CPU) instr_Dxyn(x, y, n byte) {
	panic("To implement opcode Dxyn")
}

// 0xEx9E - SKP Vx
// Skip next instruction if key with the value of Vx is pressed.
//
// Checks the keyboard, and if the key corresponding to the value of Vx is currently in the down position, PC is increased by 4,
// otherwise by 2.
func (cpu *CPU) instr_Ex9E(x byte) {
	b := cpu.R.V[x]
	if cpu.Keyboard.KeyState[b] {
		cpu.R.PC += 4
	} else {
		cpu.R.PC += 2
	}
}

// 0xExA1 - SKNP Vx
// Skip next instruction if key with the value of Vx is pressed.
//
// Checks the keyboard, and if the key corresponding to the value of Vx is currently in the up position, PC is increased by 4,
// otherwise by 2.
func (cpu *CPU) instr_ExA1(x byte) {
	b := cpu.R.V[x]
	if cpu.Keyboard.KeyState[b] {
		cpu.R.PC += 2
	} else {
		cpu.R.PC += 4
	}
}

// 0xFx07 - LD Vx, DT
// Set Vx = delay timer value.
//
// The value of DT is placed into Vx.
func (cpu *CPU) instr_Fx07(x byte) {
	cpu.R.V[x] = cpu.R.DT
	cpu.R.PC += 2
}

// 0xFx0A - LD Vx, K
// Wait for a key press, store the value of the key in Vx.
//
// All execution stops until a key is pressed, then the value of that key is stored in Vx.
func (cpu *CPU) instr_Fx0A(x byte) {
	b := keyboard.WaitForNextKey()
	cpu.R.V[x] = b
	cpu.R.PC += 2
}

// OxFx15 - LD DT, Vx
// Set delay timer = Vx.
//
// DT is set equal to the value of Vx.
func (cpu *CPU) instr_Fx15(x byte) {
	cpu.R.DT = cpu.R.V[x]
	cpu.R.PC += 2
}

// 0xFx18 - LD ST, Vx
// Set sound timer = Vx.
//
// ST is set equal to the value of Vx.
func (cpu *CPU) instr_Fx18(x byte) {
	cpu.R.ST = cpu.R.V[x]
	cpu.R.PC += 2
}

// 0xFx1E - ADD I, Vx
// Set I = I + Vx.
//
// The values of I and Vx are added, and the results are stored in I.
func (cpu *CPU) instr_Fx1E(x byte) {
	cpu.R.I = cpu.R.I + rune(cpu.R.V[x])
	cpu.R.PC += 2
}

// 0xFx29 - LD F, Vx
// Set I = location of the sprite for digit Vx.
//
// The value of I is set to the location for the hexadecimal sprite corresponding to the value of Vx.
func (cpu *CPU) instr_Fx29(x byte) {
	cpu.R.I = display.SpritesAddresses[cpu.R.V[x]]
	cpu.R.PC += 2
}

// 0xFx33 - LD B, Vx
// Stores BCD representation of Vx in memory locations I, I + 1, I + 2.
//
// The interpreter takes the decimal value of Vx, and places the hundreds digit in memory at location I, the tens digit at location I+1,
// and the ones digit at location I+2.
func (cpu *CPU) instr_Fx33(x byte) {
	value := cpu.R.V[x]
	ones := value % 10
	value /= 10
	tens := value % 10
	value /= 10
	hundreds := value % 10

	if err := cpu.Memory.SetByte(cpu.R.I, hundreds); err != nil {
		panic(err)
	}

	if err := cpu.Memory.SetByte(cpu.R.I+1, tens); err != nil {
		panic(err)
	}

	if err := cpu.Memory.SetByte(cpu.R.I+2, ones); err != nil {
		panic(err)
	}
}

// 0xFx55 - LD [I], Vx
// Store registers V0 through Vx in memory starting at location I.
//
// The interpreter copies the values of registers V0 through Vx into memory, starting at the address in I.
func (cpu *CPU) instr_Fx55(x byte) {
	for i := 0; i <= int(x); i++ {
		if err := cpu.Memory.SetByte(cpu.R.I+rune(i), cpu.R.V[i]); err != nil {
			panic(err)
		}
	}

	cpu.R.PC += 2
}

// 0xFx65 - LD Vx, [I]
// Read registers V0 through Vx from memory starting at location I.
//
// The interpreter reads values from memory starting at location I into registers V0 through Vx.
func (cpu *CPU) instr_Fx65(x byte) {
	for i := 0; i <= int(x); i++ {
		b, err := cpu.Memory.GetByte(cpu.R.I + rune(i))
		if err != nil {
			panic(err)
		}
		cpu.R.V[i] = b
	}

	cpu.R.PC += 2
}
