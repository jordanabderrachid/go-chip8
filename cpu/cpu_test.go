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
