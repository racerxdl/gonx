package graphics

import (
	"fmt"
	"github.com/racerxdl/gonx/internal"
	"github.com/racerxdl/gonx/nx/nxtypes"
	"unsafe"
)

type Window struct {
	ptr unsafe.Pointer
}

func GetDefaultWindow() *Window {
	return &Window{
		ptr: internal.NWindowGetDefault(),
	}
}

func (w *Window) MakeFramebuffer(width, height, numFbs uint32, format PixelFormat) (*Framebuffer, error) {
	fb := makeEmptyFB()
	fb.width = width
	fb.height = height
	fb.format = format
	fb.buildBuffer()

	r := internal.FramebufferCreate(fb.ptr(), w.ptr, width, height, uint32(format), numFbs)
	if r != nxtypes.ResultOK {
		return nil, fmt.Errorf("error creating framebuffer")
	}

	return fb, nil
}
