package graphics

import (
	"fmt"
	"github.com/racerxdl/gonx/nx/internal"
	"image"
	"unsafe"
)

const framebufferStructLen = 80 // Bytes

type Framebuffer struct {
	nativeBuff unsafe.Pointer
	height     uint32
	width      uint32
	format     PixelFormat
}

func makeEmptyFB() *Framebuffer {
	return &Framebuffer{
		nativeBuff: internal.Alloc(framebufferStructLen),
	}
}

func (fb *Framebuffer) ptr() unsafe.Pointer {
	return fb.nativeBuff
}

// MakeLinear Enables linear framebuffer mode in a Framebuffer, allocating a shadow buffer in the process.
func (fb *Framebuffer) MakeLinear() {
	internal.FramebufferMakeLinear(fb.ptr())
}

// Close Closes a \ref Framebuffer object, freeing all resources associated with it.
func (fb *Framebuffer) Close() error {
	internal.FramebufferClose(fb.ptr())
	internal.Free(fb.ptr())
	return nil
}

// StartFrameAsRGBA starts the frame as RGBA Image. Returns an error if image is not RGBA8888
// You must call .End() on generated frame to swap buffers
func (fb *Framebuffer) StartFrameAsRGBA() (RGBAFBImage, error) {
	if fb.format != PixelFormatRgba8888 {
		return RGBAFBImage{}, fmt.Errorf("RGBA image requires RGBA8888 framebuffer format")
	}

	v := RGBAFBImage{
		parentPtr: fb.ptr(),
	}
	stride := uint32(0)

	vptr := internal.FramebufferBegin(fb.ptr(), unsafe.Pointer(&stride))
	if vptr == nil {
		return RGBAFBImage{}, fmt.Errorf("error starting framebuffer")
	}

	bufLen := fb.height * stride
	v.RGBA = &image.RGBA{
		Pix:    internal.PointerToByteSlice(vptr, bufLen),
		Stride: int(stride),
		Rect:   image.Rect(0, 0, int(fb.width), int(fb.height)),
	}

	return v, nil
}
