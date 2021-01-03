package nxerrors

type constError string

func (err constError) Error() string {
	return string(err)
}

func (err constError) String() string {
	return string(err)
}

// Generic errors
const (
	NotImplemented            = constError("not implemented")
	OutOfMemory               = constError("out of memory")
	CannotSetMemoryAttributes = constError("cannot set memory attributes")
	Timeout                   = constError("timeout")
)

// IPC Errors
const (
	TooManyHandles                    = constError("too many handles")
	TooManyBuffers                    = constError("too many buffers")
	UnsupportedBufferType             = constError("unsupported buffer type")
	InvalidBufferAddress              = constError("invalid buffer address")
	InvalidBufferSize                 = constError("invalid buffer size")
	InvalidBufferFlags                = constError("invalid buffer flags")
	InvalidRequestType                = constError("invalid request type")
	InvalidDomain                     = constError("invalid domain")
	InvalidHandle                     = constError("invalid handle")
	CantSendDomainObjectToSession     = constError("cant send domain object to session")
	TooManyObjects                    = constError("too many objects")
	InvalidRawDataSize                = constError("invalid raw data size")
	CantCloseSessionLikeDomainObjects = constError("can't close sessions like domain objects")
	MalformedCloseRequest             = constError("malformed close request")
	CantSendObjectAcrossDomains       = constError("can't send object across domains")
	InvalidIPCResponseType            = constError("invalid ipc response type")
	InvalidIPCResponseMagic           = constError("invalid ipc response magic")
	UnexpectedRawDataSize             = constError("unexpected raw data size")
	UnexpectedPID                     = constError("unexpected pid")
	UnexpectedCopyHandles             = constError("unexpected copy handles")
	UnexpectedMoveHandles             = constError("unexpected move handles")
	UnexpectedObjects                 = constError("unexpected objects")
	ExpectedSessionClosure            = constError("expected session closure")
	RefusalToConvertBorrowedObject    = constError("refusal to convert borrowed object")
	AlreadyADomain                    = constError("already a domain")
)

// SM Errors
const (
	SMNotInitialized     = constError("sm not initialized")
	SMServiceNameTooLong = constError("sm service name too long")
)

// NV Errors
const (
	NVNotInitialized = constError("nv not initialized")
)

// GPU Errors
const (
	GPUNotInitialized  = constError("gpu not initialized")
	GPUBufferUnaligned = constError("gpu buffer unaligned")
)

// VI Errors
const (
	VINotInitialized = constError("vi not initialized")
	ParcelDataTooBig = constError("parcel data too big")
)

// AM Errors
const (
	AMNotInitialized  = constError("am not initialized")
	UnknownAppletType = constError("unknown applet type")

	AMBusy = 0x19280
)

// Display Errors
const (
	DisplayNotInitialized              = constError("display not initialized")
	ParcelDataUnderrun                 = constError("parcel data underrun")
	DisplayInvalidFence                = constError("invalid display fence")
	DisplayFenceTooManyFds             = constError("too many display fence file descriptors")
	DisplayGraphicBufferLengthMismatch = constError("display graphic buffer length mismatch")
	SurfaceInvalidState                = constError("surface invalid state")
	SurfaceBufferDequeueFailed         = constError("surface buffer dequeue failed")
	SurfaceBufferQueueFailed           = constError("surface buffer queue failed")
)

// SM Errors
const (
	SmNotInitialized     = constError("sm not initialized")
	SmServiceNameTooLong = constError("sm service name too long")
)
