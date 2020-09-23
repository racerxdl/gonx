package vi

import (
	"encoding/binary"
	"github.com/racerxdl/gonx/nx/internal"
	"github.com/racerxdl/gonx/nx/nxerrors"
	"unsafe"
)

// Parcel Represents a parcel
// 	Bounds checking is the caller's responsibility.
// 	Objects aren't currently supported very well.
type Parcel struct {
	Contents struct {
		DataSize      uint32
		DataOffset    uint32
		ObjectsSize   uint32
		ObjectsOffset uint32
		Payload       [0x200]byte
	}
	ReadHead         int
	WriteHead        int
	WritingFinalized bool
}

func (p *Parcel) ReadString() string {
	length := p.ReadU32()
	if length == 0xFFFFFFFF {
		return ""
	}

	size := 2 * (length + 1)

	u16 := p.ReadInPlace(int(size))
	u8 := make([]byte, len(u16)/2)
	for i := range u8 {
		u8[i] = uint8(binary.LittleEndian.Uint16(u16[i*2 : i*2+1]))
	}

	return string(u8)
}

func (p *Parcel) ReadU32() uint32 {
	d := p.ReadInPlace(4)
	return binary.LittleEndian.Uint32(d)
}

func (p *Parcel) ReadBinder() (*Binder, error) {
	fbo := FlatBinderObject{}
	d := p.ReadInPlace(int(unsafe.Sizeof(fbo)))
	internal.Memcpy(unsafe.Pointer(&fbo), unsafe.Pointer(&d[0]), uintptr(len(d)))

	binder := &Binder{}
	binder.handle = fbo.GetHandle()

	err := binder.AdjustRefCount(1, 0)
	if err != nil {
		return nil, err
	}
	err = binder.AdjustRefCount(1, 1)

	if err != nil {
		return nil, err
	}

	return binder, nil
}

func (p *Parcel) ReadInPlace(length int) []byte {
	d := p.Contents.Payload[p.ReadHead : p.ReadHead+length]
	p.ReadHead += length
	return d
}

func (p *Parcel) Remaining() int {
	return p.WriteHead - p.ReadHead
}

func (p *Parcel) WriteRemaining() int {
	return len(p.Contents.Payload) - p.WriteHead
}

func (p *Parcel) FinalizeWriting() ([]byte, int) {
	p.WritingFinalized = true
	p.Contents.DataSize = uint32(p.WriteHead)
	p.Contents.DataOffset = 0x10
	p.Contents.ObjectsSize = 0
	p.Contents.ObjectsOffset = uint32(0x10 + p.WriteHead)

	buff := make([]byte, unsafe.Sizeof(p.Contents))
	internal.Memcpy(unsafe.Pointer(&buff[0]), unsafe.Pointer(&p.Contents), uintptr(len(buff)))

	return buff, int(0x10 + p.WriteHead)
}

func ParcelLoad(flattened []byte) (*Parcel, error) {
	p := &Parcel{}
	p.Contents.DataSize = binary.LittleEndian.Uint32(flattened[0:4])
	p.Contents.DataOffset = binary.LittleEndian.Uint32(flattened[4:8])
	p.Contents.ObjectsSize = binary.LittleEndian.Uint32(flattened[8:12])
	p.Contents.ObjectsOffset = binary.LittleEndian.Uint32(flattened[12:16])

	if p.Contents.DataSize > 0x200 { // bigger than payload
		return nil, nxerrors.ParcelDataTooBig
	}

	copy(p.Contents.Payload[:p.Contents.ObjectsSize], flattened[p.Contents.DataOffset:])

	p.WriteHead = int(p.Contents.DataSize)
	p.ReadHead = 0
	p.WritingFinalized = true

	return p, nil
}
