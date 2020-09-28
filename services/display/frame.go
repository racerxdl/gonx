package display

import (
	"github.com/racerxdl/gonx/font"
	"github.com/racerxdl/gonx/nx/nxtypes"
	"github.com/racerxdl/gonx/svc"
	"image"
	"image/color"
	"time"
	"unsafe"
)

const minScreenFps = 30 // 30 fps
const vsyncTimeout = time.Second / minScreenFps

type Frame struct {
	surface     *Surface
	buff        []byte // Local Image Buffer
	surfaceBuff []byte // GPU Remote Buffer
	bounds      image.Rectangle
}

func (f *Frame) Convert(c color.Color) color.Color {
	return c
}

// ColorModel returns the Image's color model.
func (f *Frame) ColorModel() color.Model {
	return f
}

// Bounds returns the domain for which At can return non-zero color.
// The bounds do not necessarily contain the point (0, 0).
func (f *Frame) Bounds() image.Rectangle {
	return f.bounds
}

// At returns the color of the pixel at (x, y).
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (f *Frame) At(x, y int) color.Color {
	w := f.bounds.Size().X
	off := (y*w + x) * 4
	_ = f.buff[off+3]
	return color.RGBA{
		R: f.buff[off+0],
		G: f.buff[off+1],
		B: f.buff[off+2],
		A: f.buff[off+3],
	}
}

func (f *Frame) SetPixel(x, y int, c color.RGBA) {
	w := f.bounds.Size().X
	off := (y*w + x) * 4
	_ = f.buff[off+3]
	f.buff[off+0] = c.R
	f.buff[off+1] = c.G
	f.buff[off+2] = c.B
	f.buff[off+3] = c.A
}

//go:nobounds
func (f *Frame) Clear(c color.RGBA) {
	u32color := (uint32(c.A) << 24) + (uint32(c.B) << 16) + (uint32(c.G) << 8) + uint32(c.R)
	u32len := len(f.buff) / 4

	for i := 0; i < u32len; i++ {
		*(*uint32)(unsafe.Pointer(&f.buff[i*4])) = u32color
	}
}

func (f *Frame) DrawStringAt(x, y int, data string, c color.RGBA, font *font.Data) {
	cx := x
	cy := y
	for _, v := range data {
		if v == '\n' {
			cy += font.CharHeight
			cx = x
			continue
		}
		if v == '\r' {
			continue
		}
		g := font.GetGlyph(uint32(v))
		g.DrawAt(cx, cy, c, f.buff, f.bounds.Size().X)
		cx += font.CharWidth
	}
}

func (f *Frame) Display() error {
	s := f.bounds.Size()
	GFXSlowSwizzlingBlit(f.surfaceBuff, f.buff, s.X, s.Y, 0, 0)
	err := f.surface.QueueBuffer()
	if err != nil {
		return err
	}

	return f.surface.refreshFrame(f)
}

func (f *Frame) WaitVSync() error {
	vsync, err := GetVSyncEvent()
	if err != nil {
		return err
	}

	return svc.WaitSynchronizationSingle(nxtypes.Handle(vsync), vsyncTimeout)
}

func (f *Frame) Destroy() error {
	if f.surface.State == SURFACE_STATE_DEQUEUED {
		return f.surface.QueueBuffer()
	}
	return nil
}
