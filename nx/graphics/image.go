package graphics

import (
	"github.com/racerxdl/gonx/internal"
	"image"
	"image/color"
	"unsafe"
)

type packedRGBA uint32

func (p packedRGBA) RGBA() (r, g, b, a uint32) {
	pv := uint32(p)
	r = pv & 0xFF
	g = (pv & (0xFF << 8)) >> 8
	b = (pv & (0xFF << 16)) >> 16
	a = (pv & (0xFF << 24)) >> 24
	return
}

type RGBAFBImage struct {
	parentPtr unsafe.Pointer
	imgPtr    unsafe.Pointer
	m         *image.RGBA
}

func (img RGBAFBImage) GetRawPointer() unsafe.Pointer {
	return img.imgPtr
}

func (img RGBAFBImage) GetStride() int {
	return img.m.Stride
}

func (img RGBAFBImage) AsImage() image.Image {
	return img.m
}

func (img RGBAFBImage) AsRGBA() *image.RGBA {
	return img.m
}

func (img RGBAFBImage) At(x, y int) color.Color {
	p := uintptr(x*4 + y*img.m.Stride)
	return *(*packedRGBA)(unsafe.Pointer(uintptr(img.imgPtr) + p))
}

// ColorModel returns the Image's color model.
func (img RGBAFBImage) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds returns the domain for which At can return non-zero color.
// The bounds do not necessarily contain the point (0, 0).
func (img RGBAFBImage) Bounds() image.Rectangle {
	return img.m.Bounds()
}

func (img *RGBAFBImage) Set(x, y int, c color.Color) {
	c1 := color.RGBAModel.Convert(c).(color.RGBA)
	img.SetRGBA(x, y, c1)
}

func (img *RGBAFBImage) SetRGBA(x, y int, c color.RGBA) {
	i := uintptr(x*4 + y*img.m.Stride)
	shade32 := uint32(c.A)<<24 + uint32(c.B)<<16 + uint32(c.G)<<8 + uint32(c.R)

	*(*uint32)(unsafe.Pointer(uintptr(img.imgPtr) + i)) = shade32
}

func (img *RGBAFBImage) SetRGBA32(x, y int, c uint32) {
	i := uintptr(x*4 + y*img.m.Stride)
	*(*uint32)(unsafe.Pointer(uintptr(img.imgPtr) + i)) = c
}

func (img RGBAFBImage) End() {
	internal.FramebufferEnd(img.parentPtr)
	img.m = nil
	img.parentPtr = nil
}
