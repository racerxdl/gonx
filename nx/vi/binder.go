package vi

// Binder Represents a remote interface
type Binder struct {
	handle int32
}

func (b *Binder) AdjustRefCount(addVal, Type int32) error {
	return AdjustRefCount(b.handle, addVal, Type)
}

// FlatBinderObject Binder object as included in a Parcel
type FlatBinderObject struct {
	Type    uint32
	Flags   uint32
	Content uintptr // union of void *binder and int32 handle
	Cookie  uintptr
}

func (fb *FlatBinderObject) GetBinder() uintptr {
	return fb.Content
}

func (fb *FlatBinderObject) GetHandle() int32 {
	return int32(fb.Content)
}

func BinderTransactParcel(binder Binder, transaction, flags uint32, in *Parcel) (*Parcel, error) {
	inFlattened, _ := in.FinalizeWriting()

	buff := make([]byte, 0x210)

	err := TransactParcel(binder.handle, transaction, flags, inFlattened, buff)
	if err != nil {
		return nil, err
	}

	return ParcelLoad(buff)
}
