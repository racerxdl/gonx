package vi

// IGBP IGraphicBufferProducer object
// https://android.googlesource.com/platform/frameworks/native/+/f7a6758/include/gui/IGraphicBufferProducer.h#51
type IGBP struct {
	igbpBinder Binder // IGraphicBufferProducer
}

type PixelFormat int

const (
	RGBA_8888 PixelFormat = 0x1 // Full RGBA channels
	RGBX_8888 PixelFormat = 0x2 // RGB channels normal, X always 255 (Alpha is ignored)
	RGB_888   PixelFormat = 0x3 // Only RGB channels, no alpha
)
