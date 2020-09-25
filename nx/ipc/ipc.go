package ipc

import (
	"encoding/binary"
	"fmt"
	"github.com/racerxdl/gonx/nx/internal"
	"github.com/racerxdl/gonx/nx/nxerrors"
	"github.com/racerxdl/gonx/nx/nxtypes"
	"github.com/racerxdl/gonx/nx/svc"
	"unsafe"
)

// PackMessage is equivalent to libtransistor ipc_pack_message call
func PackMessage(msg *PackedMessage, buffer *[0x40]uint32) error {
	if debug {
		println("ipc: pack ipc message")
	}

	aDescriptors := make([]*Buffer, 0, maxIPCDescriptors)
	bDescriptors := make([]*Buffer, 0, maxIPCDescriptors)
	cDescriptors := make([]*Buffer, 0, maxIPCDescriptors)
	xDescriptors := make([]*Buffer, 0, maxIPCDescriptors)

	if debug {
		println("ipc: pack: prepare")
	}
	for _, ipcBuffer := range msg.Buffers {
		if ipcBuffer.Type&0x20 == 0 {
			if ipcBuffer.Direction() == DirectionInput { // AX
				if ipcBuffer.Family() == FamilyA {
					aDescriptors = append(aDescriptors, ipcBuffer)
				} else if ipcBuffer.Family() == FamilyX { // X
					xDescriptors = append(xDescriptors, ipcBuffer)
				} else {
					return nxerrors.UnsupportedBufferType
				}
			} else if ipcBuffer.Direction() == DirectionOutput { // BC
				if ipcBuffer.Family() == FamilyB {
					bDescriptors = append(bDescriptors, ipcBuffer)
				} else if ipcBuffer.Family() == FamilyC { // C
					cDescriptors = append(cDescriptors, ipcBuffer)
				} else {
					return nxerrors.UnsupportedBufferType
				}
			} else {
				return nxerrors.UnsupportedBufferType
			}
		} else { // flag 0x20 set
			if ipcBuffer.Type == 0x21 { // IN (ax)
				aDescriptors = append(aDescriptors, ipcBuffer)
				xDescriptors = append(xDescriptors, &nullBuffer)
			} else if ipcBuffer.Type == 0x22 { // OUT (bc)
				bDescriptors = append(bDescriptors, ipcBuffer)
				cDescriptors = append(cDescriptors, &nullBuffer)
			} else {
				return nxerrors.UnsupportedBufferType
			}
		}

		// Check for overflow
		if len(aDescriptors) >= maxIPCDescriptors ||
			len(bDescriptors) >= maxIPCDescriptors ||
			len(cDescriptors) >= maxIPCDescriptors ||
			len(xDescriptors) >= maxIPCDescriptors {
			return nxerrors.TooManyBuffers
		}
	}

	if debug {
		println("ipc: pack: packing")
	}
	h := 0 // H for header count

	buffer[h] = uint32(msg.Type) |
		(uint32(len(xDescriptors)) << 16) |
		(uint32(len(aDescriptors)) << 20) |
		(uint32(len(bDescriptors)) << 24) |
		(0 << 28) // "w" descriptors
	h++

	cDescriptorFlags := 0
	if len(cDescriptors) == 1 {
		cDescriptorFlags = 2
	} else if len(cDescriptors) > 1 {
		cDescriptorFlags = len(cDescriptors) + 2
	}
	handleDescriptorEnabled := len(msg.CopyHandles) > 0 || len(msg.MoveHandles) > 0 || msg.SendPID
	sizeFieldOffset := h

	handleDescriptorEnabledNum := 0
	if handleDescriptorEnabled {
		handleDescriptorEnabledNum = 1
	}

	// header field 2
	buffer[h] = 0 | // Size to be filled
		(uint32(cDescriptorFlags) << 10) |
		(uint32(handleDescriptorEnabledNum) << 31)
	h++

	if handleDescriptorEnabled {
		if len(msg.CopyHandles) >= maxIPCDescriptors || len(msg.MoveHandles) >= maxIPCDescriptors {
			return nxerrors.TooManyHandles
		}

		sendPid := 0
		if msg.SendPID {
			sendPid = 1
		}

		buffer[h] = uint32(sendPid) |
			(uint32(len(msg.CopyHandles)) << 1) |
			(uint32(len(msg.MoveHandles)) << 5)
		h++

		if msg.SendPID {
			h += 2
		}

		for _, v := range msg.CopyHandles {
			buffer[h] = uint32(v)
			h++
		}
		for _, v := range msg.MoveHandles {
			buffer[h] = uint32(v)
			h++
		}
	}

	// X Descriptors
	for i, x := range xDescriptors {
		if x.Addr>>39 > 0 {
			return nxerrors.InvalidBufferAddress
		}

		if x.Size>>16 > 0 {
			return nxerrors.InvalidBufferSize
		}

		// This mess -> https://switchbrew.org/wiki/IPC_Marshalling#Buffer_descriptor_X_.22Pointer.22
		addr64 := uint64(x.Addr)

		buffer[h] = uint32(i)
		buffer[h] |= uint32(((addr64 >> 36) & 7) << 6)
		buffer[h] |= uint32(((i >> 9) & 7) << 9)
		buffer[h] |= uint32(((addr64 >> 32) & 0xF) << 12)
		buffer[h] |= uint32(x.Size) << 16
		h++

		buffer[h] = uint32(addr64 & 0xFFFFFFFF)
		h++
	}

	// A descriptors
	for _, a := range aDescriptors {
		if a.Addr>>39 > 0 {
			return nxerrors.InvalidBufferAddress
		}

		if a.Size>>35 > 0 {
			return nxerrors.InvalidBufferSize
		}

		if a.Type>>8 > 4 {
			return nxerrors.InvalidBufferFlags
		}

		buffer[h] = uint32(a.Size & 0xFFFFFFFF)
		h++
		buffer[h] = uint32(a.Addr & 0xFFFFFFFF)
		h++

		addr64 := uint64(a.Addr)

		buffer[h] = uint32(a.Type) >> 6
		buffer[h] |= uint32((addr64>>36)&0x7) << 2
		buffer[h] |= uint32((a.Size>>32)&0xF) << 24
		buffer[h] |= uint32((addr64>>32)&0xF) << 28

		h++
	}

	// B Descriptors
	for _, b := range bDescriptors {
		if b.Addr>>39 > 0 {
			return nxerrors.InvalidBufferAddress
		}

		if b.Size>>35 > 0 {
			return nxerrors.InvalidBufferSize
		}

		if b.Type>>8 > 4 {
			return nxerrors.InvalidBufferFlags
		}

		buffer[h] = uint32(b.Size & 0xFFFFFFFF)
		h++
		buffer[h] = uint32(b.Addr & 0xFFFFFFFF)
		h++

		addr64 := uint64(b.Addr)

		buffer[h] = uint32(b.Type) >> 6
		buffer[h] |= uint32((addr64>>36)&0x7) << 2
		buffer[h] |= uint32((b.Size>>32)&0xF) << 24
		buffer[h] |= uint32((addr64>>32)&0xF) << 28

		h++
	}

	// "w" descriptors would go here

	rawDataStart := h
	h = int(uint32(h+3) & ^uint32(3))

	prePadding := h - rawDataStart

	if len(msg.DataSection) > 0 {
		internal.Memcpy(unsafe.Pointer(&buffer[h]), unsafe.Pointer(&msg.DataSection[0]), uintptr(len(msg.DataSection)))
		paddedSize := ipcPadSize(uint64(len(msg.DataSection)))
		paddedSize /= 4
		h += int(paddedSize)
	}

	h += 4 - prePadding

	u16LengthList := make([]uint16, 0)

	// c descriptor u16 length list
	for _, buf := range cDescriptors {
		if buf.Type&0x10 == 0 { // u16 length list flag
			if buf.Size>>16 > 0 {
				return nxerrors.InvalidBufferSize
			}
			u16LengthList = append(u16LengthList, uint16(buf.Size))
		}
	}

	if len(u16LengthList) > 0 {
		// Copy to IPC Buffer
		internal.Memcpy(unsafe.Pointer(&buffer[h]), unsafe.Pointer(&u16LengthList[0]), uintptr(len(u16LengthList))*unsafe.Sizeof(uint16(0)))
	}

	// Move header to point to right position
	h += (len(u16LengthList) + 1) >> 1

	buffer[sizeFieldOffset] |= uint32(h - rawDataStart) // raw data section size

	// C Descriptors
	for _, c := range cDescriptors {
		if c.Addr>>48 > 0 {
			return nxerrors.InvalidBufferAddress
		}

		if c.Size>>16 > 0 {
			return nxerrors.InvalidBufferSize
		}

		addr64 := uint64(c.Addr)
		buffer[h] = uint32(addr64 & 0xFFFFFFFF)
		h++
		buffer[h] = uint32(addr64 >> 32)
		buffer[h] |= uint32(c.Size) << 16
		h++
	}

	return nil
}

// PackIPCRequest is equivalent to libtransistor ipc_pack_request call
func PackIPCRequest(rq *Request, object Object, marshalBuffer *[0x40]uint32) error {
	if debug {
		println("ipc: pack ipc request")
	}

	msg := PackedMessage{}

	toDomain := rq.Type == 4 && object.ObjectID >= 0
	msg.Buffers = rq.Buffers

	if uint32(rq.Type) & ^uint32(0xFFFF) > 0 {
		return nxerrors.InvalidRequestType
	}

	msg.Type = rq.Type
	moveHandles := make([]nxtypes.Handle, 0, maxIPCDescriptors)

	if !toDomain {
		for _, obj := range rq.Objects {
			if obj.ObjectID >= 0 {
				return nxerrors.CantSendDomainObjectToSession
			}
			moveHandles = append(moveHandles, nxtypes.Handle(obj.GetSession()))
		}
	}

	for _, handle := range rq.MoveHandles {
		moveHandles = append(moveHandles, handle)
	}

	if len(moveHandles) > maxIPCDescriptors {
		return nxerrors.TooManyHandles
	}

	msg.CopyHandles = rq.CopyHandles
	msg.SendPID = rq.SendPID

	var buff [0x200 >> 2]uint32

	h := 0
	dataSectionLen := 0 // In bytes

	if toDomain {
		if len(rq.Objects) > 8 {
			return nxerrors.TooManyObjects
		}
		v := uint32(1)
		if rq.CloseObject {
			v = 2
		}
		v |= uint32(len(rq.Objects) << 8)
		buff[h] = v
		h++
		buff[h] = uint32(object.ObjectID)
		h++

		h += 2 // alignment

		dataSectionLen += 0x10
	}

	payloadSize := 0

	if !rq.CloseObject {
		buff[h] = sfci
		h++

		buff[h] = 0
		h++

		buff[h] = rq.RequestID
		h++

		buff[h] = 0
		h++

		payloadSize += 0x10
		dataSectionLen += 0x10

		if len(rq.RawData) > 0x200 {
			return nxerrors.InvalidRawDataSize
		}
		if len(rq.RawData) > 0 {
			internal.Memcpy(unsafe.Pointer(&buff[dataSectionLen/4]), unsafe.Pointer(&rq.RawData[0]), uintptr(len(rq.RawData)))
			payloadSize += len(rq.RawData)
			dataSectionLen += len(rq.RawData)
		}
	} else {
		if !toDomain {
			return nxerrors.CantCloseSessionLikeDomainObjects
		}

		if rq.Type != 4 ||
			len(rq.Buffers) != 0 ||
			len(rq.RawData) != 0 ||
			rq.SendPID != false ||
			len(rq.CopyHandles) != 0 ||
			len(rq.MoveHandles) != 0 ||
			len(rq.Objects) != 0 {
			return nxerrors.MalformedCloseRequest
		}
	}

	if toDomain {
		buff[0] |= uint32(payloadSize) << 16
		for _, obj := range rq.Objects {
			if obj.GetDomain() != object.GetDomain() {
				return nxerrors.CantSendObjectAcrossDomains
			}
			internal.Memcpy(unsafe.Pointer(&buff[dataSectionLen/4]), unsafe.Pointer(&obj.ObjectID), 4)
			dataSectionLen += 4
		}
	}

	if dataSectionLen > 0 {
		msg.DataSection = make([]byte, dataSectionLen)
		internal.Memcpy(unsafe.Pointer(&msg.DataSection[0]), unsafe.Pointer(&buff[0]), uintptr(dataSectionLen))
	}

	return PackMessage(&msg, marshalBuffer)
}

// UnpackIPCMessage is equivalent to libtransistor ipc_unpack call
func UnpackIPCMessage(msg *Message, buffer *[0x40]uint32) error {
	if debug {
		println("ipc: unpack ipc message")
	}
	h := 0 // HEAD position

	header0 := buffer[h]
	h++
	header1 := buffer[h]
	h++

	msg.MessageType = uint16(header0 & 0xFFFF)

	msg.NumXDescriptors = (header0 >> 16) & 0xF
	msg.NumADescriptors = (header0 >> 20) & 0xF
	msg.NumBDescriptors = (header0 >> 24) & 0xF
	msg.NumWDescriptors = (header0 >> 28) & 0xF

	msg.RawDataSectionSize = header1 & 0xFFFFF //  0b11 1111 1111

	msg.CDescriptorFlags = (header1 >> 10) & 0xF
	hasHandleDescriptor := (header1 >> 31) > 0

	msg.NumCopyHandles = 0
	msg.NumMoveHandles = 0
	msg.CopyHandles = nil
	msg.MoveHandles = nil
	msg.HasPID = false
	msg.PID = 0

	if hasHandleDescriptor {
		handleDescriptor := buffer[h]
		h++

		if handleDescriptor&1 > 0 {
			msg.HasPID = true
			msg.PID = *(*uint64)(unsafe.Pointer(&buffer[h]))
			h += 2
		}

		msg.NumCopyHandles = (handleDescriptor >> 1) & 0xF
		msg.NumMoveHandles = (handleDescriptor >> 5) & 0xF

		if msg.NumCopyHandles > 0 {
			msg.CopyHandles = buffer[h:]
			h += int(msg.NumCopyHandles)
		}

		if msg.NumMoveHandles > 0 {
			msg.MoveHandles = buffer[h:]
			h += int(msg.NumMoveHandles)
		}
	}

	// Descriptors

	if msg.NumXDescriptors > 0 {
		msg.XDescriptors = buffer[h:]
		h += int(msg.NumXDescriptors * 2)
	}

	if msg.NumADescriptors > 0 {
		msg.ADescriptors = buffer[h:]
		h += int(msg.NumADescriptors * 3)
	}
	if msg.NumBDescriptors > 0 {
		msg.BDescriptors = buffer[h:]
		h += int(msg.NumBDescriptors * 3)
	}

	if msg.NumWDescriptors > 0 {
		msg.WDescriptors = buffer[h:]
		h += int(msg.NumWDescriptors * 3)
	}

	before := h

	// Align head to 4 words
	h = int(uint32(h+3) & ^uint32(3))

	msg.PrePadding = h - before
	msg.PostPadding = 4 - msg.PrePadding
	if msg.RawDataSectionSize > 0 {
		msg.DataSection = buffer[h:]
	}

	h = before + int(msg.RawDataSectionSize)

	msg.CDescriptors = buffer[h:]

	return nil
}

// UnflattenResponse is equivalent to libtransistor ipc_unflatten_response
func UnflattenResponse(msg *Message, rs *ResponseFmt, object Object) error {
	if debug {
		println("ipc: unflatten ipc response")
	}
	fromDomain := object.ObjectID >= 0

	if msg.MessageType != 0 && msg.MessageType != 4 {
		return nxerrors.InvalidIPCResponseType
	}

	h := 0

	if fromDomain {
		h += 4 // skip domain header
	}

	if msg.DataSection[h] != sfco {
		return nxerrors.InvalidIPCResponseMagic
	}
	h += 2

	responseCode := msg.DataSection[h]
	h++

	if responseCode != nxtypes.ResultOK {
		return nxerrors.IPCError{
			Result:  uint64(responseCode),
			Message: "response error",
		}
	}
	h++

	rawData := msg.DataSection[h:]

	nObjs := int(0)

	if fromDomain {
		nObjs = 0x10 + len(rs.Objects)*4
	}

	// RawDataSectionLength - SFCI, Command ID - Padding - nObjs
	if (int(msg.RawDataSectionSize*4) - 0x10 - 0x10 - nObjs) != int(ipcPadSize(uint64(len(rs.RawData)))) {
		if debug {
			v := ipcPadSize(uint64(len(rs.RawData)))
			println("expected", int(msg.RawDataSectionSize*4)-0x10-0x10-nObjs, "got", v)
			println("raw data section size", msg.RawDataSectionSize)
		}
		return nxerrors.UnexpectedRawDataSize
	}

	if msg.HasPID != rs.HasPID {
		return nxerrors.UnexpectedPID
	}

	if int(msg.NumCopyHandles) != len(rs.CopyHandles) {
		return nxerrors.UnexpectedCopyHandles
	}

	numObjs := 0
	if !fromDomain {
		numObjs = len(rs.Objects)
	}

	if int(msg.NumMoveHandles) != len(rs.MoveHandles)+numObjs {
		return nxerrors.UnexpectedMoveHandles
	}

	if fromDomain {
		type responseDomainHeader struct {
			numObjects uint32
			unknown1   [2]uint32
			unknown2   uint32
		}
		ptr := unsafe.Pointer(&msg.DataSection[0])
		domainHeader := (*responseDomainHeader)(ptr)

		if int(domainHeader.numObjects) != len(rs.Objects) {
			return nxerrors.UnexpectedObjects
		}

		// this is a pointer to a uint32_t array, but it is allowed to be unaligned
		ptru := uintptr(ptr)
		ptru += unsafe.Sizeof(responseDomainHeader{})
		ptru += 0x10 // SFCO, result code
		ptru += uintptr(len(rs.RawData))

		domainIds := ptru

		for i := range rs.Objects {
			rs.Objects[i].Content = object.Content
			internal.Memcpy(unsafe.Pointer(&rs.Objects[i].ObjectID), unsafe.Pointer(domainIds+uintptr(i*4)), 4)
			rs.Objects[i].IsBorrowed = false
		}
	}

	for i := range rs.CopyHandles {
		rs.CopyHandles[i] = nxtypes.Handle(msg.CopyHandles[i])
	}

	mhi := 0 // move handle index

	if !fromDomain {
		for i := range rs.Objects {
			rs.Objects[i].Content = uint64(msg.MoveHandles[mhi])
			rs.Objects[i].ObjectID = -1
			rs.Objects[i].IsBorrowed = false
			mhi++
		}
	}

	for i := range rs.MoveHandles {
		rs.MoveHandles[i] = nxtypes.Handle(msg.MoveHandles[mhi])
		mhi++
	}

	if rs.HasPID {
		*rs.PID = msg.PID
	}

	if len(rs.RawData) > 0 {
		internal.Memcpy(unsafe.Pointer(&rs.RawData[0]), unsafe.Pointer(&rawData[0]), uintptr(len(rs.RawData)))
	}

	return nil
}

func Send(object Object, rq *Request, rs *ResponseFmt) error {
	svc.ClearIPCBuffer()
	ipcBuff := svc.GetIPCBuffer()

	err := PackIPCRequest(rq, object, ipcBuff)
	if err != nil {
		if debug {
			println("ipc: packing error:", err.Error())
		}
		return err
	}

	if debug {
		println("ipc: send sync request")
	}

	if debugDumpBeforeSend {
		svc.DumpIPCBuffer()
	}

	r := svc.SendSyncRequest(object.Content)
	if r > 0 {
		if debug {
			fmt.Printf("ipc: bad request with return code %x\n", r)
			svc.DumpIPCBuffer()
		}

		return nxerrors.IPCError{
			Result:  r,
			Message: "bad request",
		}
	}

	if debug {
		println("ipc: processing response")
	}

	msg := Message{}

	err = UnpackIPCMessage(&msg, ipcBuff)

	if err != nil {
		if debug {
			println("ipc: unpacking error:", err.Error())
			svc.DumpIPCBuffer()
		}
		return err
	}

	err = UnflattenResponse(&msg, rs, object)

	if err != nil {
		if debug {
			println("ipc: unflatten error:", err.Error())
		}
		return err
	}

	return nil
}

func CloseSession(session nxtypes.SessionHandle) error {
	var err error

	rq := MakeDefaultRequest(0)
	rq.Type = TypeIPCClose

	obj := Object{
		ObjectID: -1,
	}

	obj.SetSession(nxtypes.Handle(session))

	svc.ClearIPCBuffer()
	ipcBuff := svc.GetIPCBuffer()

	err = PackIPCRequest(&rq, obj, ipcBuff)
	if err != nil {
		return err
	}

	if debug {
		println("ipc: send sync request")
	}

	if debugDumpBeforeSend {
		svc.DumpIPCBuffer()
	}

	r := svc.SendSyncRequest(uint64(session))
	if r != 0xf601 {
		if debug {
			fmt.Printf("ipc: expected session closure, got %x\n", r)
		}
		err = nxerrors.IPCError{
			Message: nxerrors.ExpectedSessionClosure.String(),
			Result:  r,
		}
	}

	if debug {
		println("ipc: svc close handle")
	}

	svc.CloseHandle(nxtypes.Handle(session))

	return err
}

func Close(object Object) error {
	if object.IsBorrowed {
		return nil // we're not allowed to close borrowed objects,
		// and we would also like to handle this transparently
	}

	if object.ObjectID < 0 {
		return CloseSession(object.GetSession())
	}

	rq := MakeDefaultRequest(0)
	rq.CloseObject = true

	svc.ClearIPCBuffer()
	ipcBuff := svc.GetIPCBuffer()

	err := PackIPCRequest(&rq, object, ipcBuff)
	if err != nil {
		return err
	}

	if debug {
		println("ipc: send sync request")
	}

	if debugDumpBeforeSend {
		svc.DumpIPCBuffer()
	}

	d := object.GetDomain()

	if d == nil {
		return nxerrors.InvalidDomain
	}

	r := svc.SendSyncRequest(uint64(d.Session))
	if r > 0 {
		if debug {
			println("ipc: error sending request")
			svc.DumpIPCBuffer()
		}
		return nxerrors.IPCError{
			Result:  r,
			Message: "error sending request",
		}
	}

	return nil
}

func ConvertToDomain(object *Object) (*Domain, error) {
	if object.IsBorrowed {
		return nil, nxerrors.RefusalToConvertBorrowedObject
	}
	if object.ObjectID != -1 {
		return nil, nxerrors.AlreadyADomain
	}

	session := *object
	domain := &Domain{
		Session: session.GetSession(),
	}

	object.SetDomain(domain)

	rq := MakeDefaultRequest(0)
	rq.Type = TypeIPCControl

	rs := ResponseFmt{}
	rs.RawData = make([]byte, unsafe.Sizeof(object.ObjectID))

	err := Send(session, &rq, &rs)
	if err != nil {
		return nil, err
	}

	object.ObjectID = int32(binary.LittleEndian.Uint32(rs.RawData))

	return domain, nil
}
