package graphics

import (
	"fmt"
	"github.com/racerxdl/gonx/internal"
	"image"
	"unsafe"
)

const framebufferStructLen = 80 // Bytes

type Framebuffer struct {
	nativeBuff unsafe.Pointer
	height     uint32
	width      uint32
	format     PixelFormat
	frame      *RGBAFBImage
}

func makeEmptyFB() *Framebuffer {
	return &Framebuffer{
		nativeBuff: internal.Alloc(framebufferStructLen),
	}
}

func (fb *Framebuffer) ptr() unsafe.Pointer {
	return fb.nativeBuff
}

func (fb *Framebuffer) buildBuffer() {
	fb.frame = &RGBAFBImage{
		parentPtr: fb.ptr(),
		m: &image.RGBA{
			Rect: image.Rect(0, 0, int(fb.width), int(fb.height)),
		},
	}
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
func (fb *Framebuffer) StartFrameAsRGBA() (*RGBAFBImage, error) {
	if fb.format != PixelFormatRgba8888 {
		return nil, fmt.Errorf("RGBA image requires RGBA8888 framebuffer format")
	}

	stride := uint32(0)

	vptr := internal.FramebufferBegin(fb.ptr(), unsafe.Pointer(&stride))
	if vptr == nil {
		return nil, fmt.Errorf("error starting framebuffer")
	}

	bufLen := fb.height * stride
	fb.frame.imgPtr = vptr
	fb.frame.m.Pix = internal.PointerToByteSlice(vptr, bufLen)
	fb.frame.m.Stride = int(stride)

	return fb.frame, nil
}

// GetDisplayer returns a tinygo-draw compatible displayer
func (fb *Framebuffer) GetDisplayer() *Displayer {
	return &Displayer{
		fb: fb,
		m:  image.NewRGBA(image.Rect(0, 0, int(fb.width), int(fb.height))),
	}
}
