package nxerrors

type IPCError struct {
	Message string
	Result  uint64
}

func (i IPCError) Error() string {
	return i.Message
}

func (i IPCError) String() string {
	return i.Message
}
