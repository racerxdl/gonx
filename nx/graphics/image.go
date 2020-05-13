package graphics

import (
	"github.com/racerxdl/gonx/nx/internal"
	"image"
	"unsafe"
)

type RGBAFBImage struct {
	parentPtr unsafe.Pointer
	*image.RGBA
}

func (img RGBAFBImage) AsImage() image.Image {
	return img.RGBA
}

func (img RGBAFBImage) AsRGBA() *image.RGBA {
	return img.RGBA
}

func (img RGBAFBImage) End() {
	internal.FramebufferEnd(img.parentPtr)
	img.RGBA = nil
	img.parentPtr = nil
}
