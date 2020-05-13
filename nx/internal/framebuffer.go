// +build nintendoswitch

package internal

import (
	"github.com/racerxdl/gonx/nx/system"
	"unsafe"
)

// Result framebufferCreate(Framebuffer* fb, NWindow *win, u32 width, u32 height, u32 format, u32 num_fbs);
//go:export framebufferCreate
func FramebufferCreate(fb unsafe.Pointer, window unsafe.Pointer, width, height, format, numFbs uint32) system.Result

// Enables linear framebuffer mode in a \ref Framebuffer, allocating a shadow buffer in the process.
// Result framebufferMakeLinear(Framebuffer* fb);
//go:export framebufferMakeLinear
func FramebufferMakeLinear(fb unsafe.Pointer)

/// Closes a \ref Framebuffer object, freeing all resources associated with it.
//void framebufferClose(Framebuffer* fb);
//go:export framebufferClose
func FramebufferClose(fb unsafe.Pointer)

// void* framebufferBegin(Framebuffer* fb, u32* out_stride);
//go:export framebufferBegin
func FramebufferBegin(fb unsafe.Pointer, outString unsafe.Pointer) unsafe.Pointer

// void framebufferEnd(Framebuffer* fb);
//go:export framebufferEnd
func FramebufferEnd(fb unsafe.Pointer)
