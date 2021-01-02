package graphics

import (
	"github.com/racerxdl/gonx/internal"
	"image"
	"image/color"
	"unsafe"
)

// Displayer struct to be compatible with tinygo.org/drivers/displayer

type Displayer struct {
	fb *Framebuffer
	m  *image.RGBA
}

// Size returns the size of the framebuffer
// Compatible with drivers.Displayer
func (d *Displayer) Size() (x, y int16) {
	return int16(d.m.Bounds().Max.X), int16(d.m.Bounds().Max.Y)
}

// SetPixel modifies the framebuffer setting the pixel at X,Y to specified color
// Compatible with drivers.Displayer
func (d *Displayer) SetPixel(x, y int16, c color.RGBA) {
	i := uintptr(int(x)*4 + int(y)*d.m.Stride)
	shade32 := uint32(c.A)<<24 + uint32(c.B)<<16 + uint32(c.G)<<8 + uint32(c.R)

	*(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(&d.m.Pix[0])) + i)) = shade32
}

// Display sends the buffer (if any) to the screen.
// Compatible with drivers.Displayer
func (d *Displayer) Display() error {
	buff, err := d.fb.StartFrameAsRGBA()
	if err != nil {
		return err
	}

	internal.Memcpy(buff.imgPtr, unsafe.Pointer(&d.m.Pix[0]), uintptr(int(d.fb.height)*d.m.Stride))
	buff.End()

	return nil
}
