package keyboard

import "github.com/nsf/termbox-go"

// This map binds the value returned by the keyboard to the corresponding chip-8 value.
var KeyMap map[int]byte = map[int]byte{
	48:  0x00, // "0"
	49:  0x01, // "1"
	50:  0x02, // "2"
	51:  0x03, // "3"
	52:  0x04, // "4"
	53:  0x05, // "5"
	54:  0x06, // "6"
	55:  0x07, // "7"
	56:  0x08, // "8"
	57:  0x09, // "9"
	97:  0x0A, // "a"
	98:  0x0B, // "b"
	99:  0x0C, // "c"
	100: 0x0D, // "d"
	101: 0x0E, // "e"
	102: 0x0F, // "f"
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

	go kb.ListenKeyboardInput()
}

func (kb *Keyboard) KeyStateToFalse() {
	for k := range kb.KeyState {
		kb.KeyState[k] = false
	}
}

func (kb *Keyboard) ActivateKey(key byte) {
	for k := range kb.KeyState {
		if k == key {
			kb.KeyState[k] = true
		} else {
			kb.KeyState[k] = false
		}
	}
}

func (kb *Keyboard) ListenKeyboardInput() {
	for {
		key := int(termbox.PollEvent().Ch)
		if v, ok := KeyMap[key]; ok {
			kb.ActivateKey(key)
		} else {
			kb.KeyStateToFalse()
		}
	}
}

func WaitForNextKey() byte {
	for {
		key := int(termbox.PollEvent().Ch)
		if v, ok := KeyMap[key]; ok {
			return v
		}
	}
}
