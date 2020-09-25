package ipc

import (
	"encoding/binary"
	"github.com/racerxdl/gonx/nx/nxtypes"
)

var (
	sfci = binary.LittleEndian.Uint32([]byte("SFCI"))
	sfco = binary.LittleEndian.Uint32([]byte("SFCO"))
)

type Direction uint8
type Family uint8

type Type uint32
type BufferFamily int

const (
	BufferFamilyA BufferFamily = 0
	BufferFamilyB BufferFamily = 1
	BufferFamilyC BufferFamily = 2
	BufferFamilyX BufferFamily = 3

	DirectionInput  Direction = 1 // 0b01
	DirectionOutput Direction = 2 // 0b10

	FamilyA Family = 1 // 0b01
	FamilyB Family = 1 // 0b01
	FamilyX Family = 2 // 0b10
	FamilyC Family = 2 // 0b10
)

type Domain struct {
	Session nxtypes.SessionHandle
}

type Buffer struct {
	Addr uintptr
	Size uint64
	Type uint32
}

func (i Buffer) Direction() Direction {
	return Direction(i.Type & 3)
}

func (i Buffer) Family() Family {
	return Family((i.Type & 12) >> 2)
}

// nullBuffer is a empty buffer
// used in some "workarounds"
var nullBuffer = Buffer{
	Addr: 0,
	Type: 0,
	Size: 0,
}

func ipcPadSize(size uint64) uint64 {
	return (size + 3) & ^uint64(3)
}
