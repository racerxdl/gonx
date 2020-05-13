package graphics

type PixelFormat uint32

const (
	PixelFormatRgba8888              PixelFormat = 1
	PixelFormatRgbx8888              PixelFormat = 2
	PixelFormatRgb888                PixelFormat = 3
	PixelFormatRgb565                PixelFormat = 4
	PixelFormatBgra8888              PixelFormat = 5
	PixelFormatRgba5551              PixelFormat = 6
	PixelFormatRgba4444              PixelFormat = 7
	PixelFormatYcrcb420Sp            PixelFormat = 17
	PixelFormatRaw16                 PixelFormat = 32
	PixelFormatBlob                  PixelFormat = 33
	PixelFormatImplementationDefined PixelFormat = 34
	PixelFormatYcbcr420888           PixelFormat = 35
	PixelFormatY8                    PixelFormat = 0x20203859
	PixelFormatY16                   PixelFormat = 0x20363159
	PixelFormatYv12                  PixelFormat = 0x32315659
)
