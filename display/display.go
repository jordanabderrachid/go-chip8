package display

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"unicode/utf8"
)

const (
	X int = 64
	Y int = 32
)

var activeColor termbox.Attribute = termbox.ColorWhite
var inactiveColor termbox.Attribute = termbox.ColorBlack

var Sprites map[string][5]byte = map[string][5]byte{
	"0": {0xF0, 0x90, 0x90, 0x90, 0xF0},
	"1": {0x20, 0x60, 0x20, 0x20, 0x70},
	"2": {0xF0, 0x10, 0xF0, 0x80, 0xF0},
	"3": {0xF0, 0x10, 0xF0, 0x10, 0xF0},
	"4": {0x90, 0x90, 0xF0, 0x10, 0x10},
	"5": {0xF0, 0x80, 0xF0, 0x10, 0xF0},
	"6": {0xF0, 0x80, 0xF0, 0x90, 0xF0},
	"7": {0xF0, 0x10, 0x20, 0x40, 0x40},
	"8": {0xF0, 0x90, 0xF0, 0x90, 0xF0},
	"9": {0xF0, 0x90, 0xF0, 0x10, 0xF0},
	"A": {0xF0, 0x90, 0xF0, 0x90, 0x90},
	"B": {0xE0, 0x90, 0xE0, 0x90, 0xE0},
	"C": {0xF0, 0x80, 0x80, 0x80, 0xF0},
	"D": {0xE0, 0x90, 0x90, 0x90, 0xE0},
	"E": {0xF0, 0x80, 0xF0, 0x80, 0xF0},
	"F": {0xF0, 0x80, 0xF0, 0x80, 0x80},
}

var SpritesAddresses map[byte]rune = map[byte]rune{
	0x00: 0x0000,
	0x01: 0x0005,
	0x02: 0x000A,
	0x03: 0x000F,
	0x04: 0x0014,
	0x05: 0x0019,
	0x06: 0x001E,
	0x07: 0x0023,
	0x08: 0x0028,
	0x09: 0x002D,
	0x0A: 0x0032,
	0x0B: 0x0037,
	0x0C: 0x003C,
	0x0E: 0x0041,
	0x0F: 0x0046,
}

type Sprite struct {
	Cells []byte
}

type Display struct {
	Cells [][]byte
}

func (d *Display) Reset() {
	termbox.HideCursor()
	d.Cells = make([][]byte, Y)
	for i := range d.Cells {
		d.Cells[i] = make([]byte, X)
	}

	d.Clear()
}

func (d *Display) Clear() {
	for y := range d.Cells {
		for x := range d.Cells[y] {
			d.Cells[y][x] = 0
		}
	}

	d.draw()
}

func (d *Display) DrawSprite(x, y int, s Sprite) (bool, error) {
	_ = "breakpoint"
	coll := false
	for iy := 0; iy < len(s.Cells); iy++ {
		barr := [8]byte{
			(s.Cells[iy] & 0x80) >> 7,
			(s.Cells[iy] & 0x40) >> 6,
			(s.Cells[iy] & 0x20) >> 5,
			(s.Cells[iy] & 0x10) >> 4,
			(s.Cells[iy] & 0x08) >> 3,
			(s.Cells[iy] & 0x04) >> 2,
			(s.Cells[iy] & 0x02) >> 1,
			(s.Cells[iy] & 0x01),
		}

		for ix := 0; ix < len(barr); ix++ {
			c, err := d.setPixel((x+ix)%X, (y+iy)%Y, barr[ix])
			if err != nil {
				return false, err
			}

			if c == true {
				coll = true
			}
		}
	}

	d.draw()
	return coll, nil
}

func (d *Display) setPixel(x, y int, b byte) (bool, error) {
	if y < 0 || y > len(d.Cells)-1 || x < 0 || x > len(d.Cells[0])-1 {
		return false, fmt.Errorf("(%d, %d) out of range of display\n", x, y)
	}

	coll := false
	if d.Cells[y][x] == 1 && b == 1 {
		coll = true
	}

	d.Cells[y][x] = b
	return coll, nil
}

func (d *Display) draw() {
	for y := range d.Cells {
		for x := range d.Cells[y] {
			var color termbox.Attribute
			if d.Cells[y][x] == 1 {
				color = activeColor
			} else {
				color = inactiveColor
			}
			r, _ := utf8.DecodeRuneInString(" ")
			termbox.SetCell(x, y, r, color, color)
		}
	}
	if err := termbox.Flush(); err != nil {
		panic(err)
	}
}
