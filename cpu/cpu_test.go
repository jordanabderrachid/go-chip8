package cpu

import "testing"

func TestRegisterReset(t *testing.T) {
	r := Registers{}
	r.Reset()
}

func TestCPUReset(t *testing.T) {
	cpu := &CPU{}
	cpu.Reset()
}

// 0x00EE - RET
// Return from a subroutine.
//
// The interpreter sets the program counter to the address at the top of the stack, then substracts 1 from the stack pointer.
func TestInstr_00EE(t *testing.T) {
	tc := []struct {
		Stack      [16]rune
		SP         byte
		PC         rune
		ExpectedSP byte
		ExpectedPC rune
	}{
		{[16]rune{0x0000, 0x0001, 0x0002, 0x0003, 0x0004}, 0x04, 0x0034, 0x03, 0x0004},
		{[16]rune{0x1234}, 0x00, 0x034, 0x00, 0x1234},
	}

	for _, c := range tc {
		cpu := &CPU{}
		r := &Registers{}
		r.Reset()
		cpu.R = r

		cpu.R.Stack = c.Stack
		cpu.R.SP = c.SP
		cpu.R.PC = c.PC

		cpu.instr_00EE()
		if cpu.R.SP != c.ExpectedSP {
			t.Errorf("stack pointer should be 0x%02x, actual: 0x%02x\n", c.ExpectedSP, cpu.R.SP)
		}

		if cpu.R.PC != c.ExpectedPC {
			t.Errorf("program counter should be 0x%04x, actual: 0x%04x\n", c.ExpectedPC, cpu.R.PC)
		}
	}
}

// 0x1nnn - JP addr
// Jump to location nnn.
//
// The interpreter sets the program counter to nnn.
func TestInstr_1nnn(t *testing.T) {
	tc := []struct {
		addr rune
	}{
		{0x000},
		{0x00F},
		{0x0F0},
		{0xF00},
		{0xFFF},
	}

	for _, c := range tc {
		cpu := &CPU{}
		r := &Registers{}
		r.Reset()

		cpu.R = r
		cpu.instr_1nnn(c.addr)
		if cpu.R.PC != c.addr {
			t.Errorf("program counter should be 0x%04x, actual: 0x%04x\n", c.addr, cpu.R.PC)
		}
	}
}

// 0x2nnn - CALL addr
// Call subroutine at nnn.
//
// The interpreter increments the stack pointer, then puts the current PC on the top of the stack. The PC is then set to nnn.
func TestInstr_2nnn(t *testing.T) {
	tc := []struct {
		Stack         [16]rune
		SP            byte
		PC            rune
		addr          rune
		ExpectedStack [16]rune
	}{
		{[16]rune{0x0000}, 0x0000, 0x0034, 0x1234, [16]rune{0x0000, 0x0036}},
	}

	for _, c := range tc {
		cpu := &CPU{}
		r := &Registers{}
		r.Reset()

		cpu.R = r
		cpu.R.Stack = c.Stack
		cpu.R.SP = c.SP
		cpu.R.PC = c.PC
		cpu.instr_2nnn(c.addr)
		if cpu.R.PC != c.addr {
			t.Errorf("program counter should be 0x%04x, actual: 0x%04x\n", c.addr, cpu.R.PC)
		}

		if cpu.R.SP != c.SP+1 {
			t.Errorf("stack pointer should be 0x%02x, actual: 0x%02x\n", c.SP+1, cpu.R.SP)
		}

		if cpu.R.Stack != c.ExpectedStack {
			t.Errorf("stack should be %s, actual: %s", c.ExpectedStack, cpu.R.Stack)
		}
	}
}
