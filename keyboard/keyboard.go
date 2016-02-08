package keyboard

import "github.com/veandco/go-sdl2/sdl"

// This map binds the value returned by the keyboard to the corresponding chip-8 value.
var KeyMap map[byte]sdl.Scancode = map[byte]sdl.Scancode{
	0x00: sdl.SCANCODE_0, // "0"
	0x01: sdl.SCANCODE_1, // "1"
	0x02: sdl.SCANCODE_2, // "2"
	0x03: sdl.SCANCODE_3, // "3"
	0x04: sdl.SCANCODE_4, // "4"
	0x05: sdl.SCANCODE_5, // "5"
	0x06: sdl.SCANCODE_6, // "6"
	0x07: sdl.SCANCODE_7, // "7"
	0x08: sdl.SCANCODE_8, // "8"
	0x09: sdl.SCANCODE_9, // "9"
	0x0A: sdl.SCANCODE_A, // "a"
	0x0B: sdl.SCANCODE_B, // "b"
	0x0C: sdl.SCANCODE_C, // "c"
	0x0D: sdl.SCANCODE_D, // "d"
	0x0E: sdl.SCANCODE_E, // "e"
	0x0F: sdl.SCANCODE_F, // "f"
}

type Keyboard struct {
	KeyState map[byte]bool
}

func (kb *Keyboard) Reset() {
	kb.KeyState = map[byte]bool{
		0x00: false,
		0x01: false,
		0x02: false,
		0x03: false,
		0x04: false,
		0x05: false,
		0x06: false,
		0x07: false,
		0x08: false,
		0x09: false,
		0x0A: false,
		0x0B: false,
		0x0C: false,
		0x0D: false,
		0x0E: false,
		0x0F: false,
	}
}

func (kb *Keyboard) KeyStateToFalse() {
	for k := range kb.KeyState {
		kb.KeyState[k] = false
	}
}

func IsKeyPressed(b byte) bool {
	keyboardState := sdl.GetKeyboardState()
	code := KeyMap[b]
	pressed := keyboardState[code]
	if pressed == 1 {
		return true
	} else {
		return false
	}
}
