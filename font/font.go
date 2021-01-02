package font

import (
	"image/color"
	"strings"
	"unsafe"
)

// Generate all font files from bdf in this folder
//go:generate go run generate.go

type Glyph struct {
	Name string
	Data [][]byte
}

func (g *Glyph) DrawAt(x, y int, c color.RGBA, pixels []byte, imgWidth int) {
	u32color := (uint32(c.A) << 24) + (uint32(c.B) << 16) + (uint32(c.G) << 8) + uint32(c.R)
	p := y*imgWidth + x

	rows := len(g.Data)
	columns := len(g.Data[0])

	for cy := 0; cy < rows; cy++ {
		row := g.Data[cy]
		for cx := 0; cx < columns; cx++ {
			if row[cx] > 0 {
				cp := p + cx + cy*imgWidth
				*(*uint32)(unsafe.Pointer(&pixels[cp*4])) = u32color
			}
		}
	}
}

type Data struct {
	CharWidth  int
	CharHeight int
	XOffset    int
	YOffset    int
	Glyphs     map[uint32]*Glyph
}

var defaultGlyph = &Glyph{
	Name: "EMPTY SPACE",
	Data: [][]byte{
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
	},
}

func (d *Data) GetGlyph(unicode uint32) *Glyph {
	if glyph, ok := d.Glyphs[unicode]; ok {
		return glyph
	}

	return defaultGlyph
}

var fonts = map[string]*Data{}

func GetFontByName(name string) *Data {
	name = strings.ToLower(name)
	if font, ok := fonts[name]; ok {
		return font
	}

	return nil
}
