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

func ListenKeyboardInput(c chan byte) {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)
	for {
		key := int(termbox.PollEvent().Ch)
		if v, ok := KeyMap[key]; ok {
			c <- v
		}
	}
}
