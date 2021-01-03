package nv

import "github.com/racerxdl/gonx/nx/nxtypes"

type Surface struct {
	Width             uint32
	Height            uint32
	ColorFormat       ColorFormat
	Layout            Layout
	Pitch             uint32
	Unused            uint32
	Offset            uint32
	Kind              Kind
	BlockHeightLog2   uint32
	Scan              DisplayScanFormat
	SecondFieldOffset uint32
	Flags             uint64
	Size              uint64
	Unk               [6]uint32
}

type GraphicBuffer struct {
	Header    nxtypes.NativeHandle
	Unk0      int32 // -1
	NVMapID   int32
	Unk2      uint32 // 0
	Magic     uint32 // 0xDAFFCAFF
	PID       uint32 // 42
	Type      uint32
	Usage     uint32 // GRALLOC_USAGE_* bitmask
	Format    uint32 // PIXEL_FORMAT_*
	ExtFormat uint32 // Copy the value in Format field
	Stride    uint32 // in pixels
	TotalSize uint32 // in bytes
	NumPlanes uint32 // 1
	Unk12     uint32 // 0
	Planes    [3]Surface
	Unused    uint64
}
