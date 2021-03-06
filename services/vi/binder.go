package vi

// Binder Represents a remote interface
type Binder struct {
	Handle int32
}

func (b *Binder) AdjustRefCount(addVal, Type int32) error {
	return AdjustRefCount(b.Handle, addVal, Type)
}

// FlatBinderObject Binder object as included in a Parcel
type FlatBinderObject struct {
	Type    uint32
	Flags   uint32
	Content uint64 // union of void *binder and int32 Handle
	Cookie  uint64
}

func (fb *FlatBinderObject) GetBinder() uintptr {
	return uintptr(fb.Content)
}

func (fb *FlatBinderObject) GetHandle() int32 {
	return int32(fb.Content)
}

func BinderTransactParcel(binder Binder, transaction, flags uint32, in *Parcel) (*Parcel, error) {
	inFlattened, _ := in.FinalizeWriting()

	buff := make([]byte, 0x210)

	err := TransactParcel(binder.Handle, transaction, flags, inFlattened, buff)
	if err != nil {
		return nil, err
	}

	return ParcelLoad(buff)
}
