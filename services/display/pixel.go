package display

type PixelFormat uint32

const (
	RGBA_8888 PixelFormat = 0x1 // Full RGBA channels
	RGBX_8888 PixelFormat = 0x2 // RGB channels normal, X always 255 (Alpha is ignored)
	RGB_888   PixelFormat = 0x3 // Only RGB channels, no alpha
)
