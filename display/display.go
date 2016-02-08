package display

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"log"
)

const (
	X             int    = 64
	Y             int    = 32
	activeColor   uint32 = 0xFFFFFFF
	inactiveColor uint32 = 0x000000
)

var Sprites map[byte][5]byte = map[byte][5]byte{
	0x00: {0xF0, 0x90, 0x90, 0x90, 0xF0}, // "0"
	0x01: {0x20, 0x60, 0x20, 0x20, 0x70}, // "1"
	0x02: {0xF0, 0x10, 0xF0, 0x80, 0xF0}, // "2"
	0x03: {0xF0, 0x10, 0xF0, 0x10, 0xF0}, // "3"
	0x04: {0x90, 0x90, 0xF0, 0x10, 0x10}, // "4"
	0x05: {0xF0, 0x80, 0xF0, 0x10, 0xF0}, // "5"
	0x06: {0xF0, 0x80, 0xF0, 0x90, 0xF0}, // "6"
	0x07: {0xF0, 0x10, 0x20, 0x40, 0x40}, // "7"
	0x08: {0xF0, 0x90, 0xF0, 0x90, 0xF0}, // "8"
	0x09: {0xF0, 0x90, 0xF0, 0x10, 0xF0}, // "9"
	0x0A: {0xF0, 0x90, 0xF0, 0x90, 0x90}, // "A"
	0x0B: {0xE0, 0x90, 0xE0, 0x90, 0xE0}, // "B"
	0x0C: {0xF0, 0x80, 0x80, 0x80, 0xF0}, // "C"
	0x0D: {0xE0, 0x90, 0x90, 0x90, 0xE0}, // "D"
	0x0E: {0xF0, 0x80, 0xF0, 0x80, 0xF0}, // "E"
	0x0F: {0xF0, 0x80, 0xF0, 0x80, 0x80}, // "F"
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
	0x0D: 0x0041,
	0x0E: 0x0046,
	0x0F: 0x004B,
}

type Sprite struct {
	Cells []byte
}

type Display struct {
	Cells   [][]byte
	Window  *sdl.Window
	Surface *sdl.Surface
}

func (d *Display) Reset() {
	var err error

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Panicln(err)
	}

	d.Window, err = sdl.CreateWindow("go chip8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, X*10, Y*10, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Panicln(err)
	}

	d.Surface, err = d.Window.GetSurface()
	if err != nil {
		log.Panicln(err)
	}

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

	d.Cells[y][x] = d.Cells[y][x] ^ b
	return coll, nil
}

func (d *Display) draw() {
	for y := range d.Cells {
		for x := range d.Cells[y] {
			var color uint32
			if d.Cells[y][x] == 1 {
				color = activeColor
			} else {
				color = inactiveColor
			}
			rect := &sdl.Rect{int32(x) * 10, int32(y) * 10, 10, 10}
			d.Surface.FillRect(rect, color)
		}
	}

	d.Window.UpdateSurface()
}
