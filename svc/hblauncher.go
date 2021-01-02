package svc

import _ "unsafe" // for go:linkname

//go:linkname GetContextPtr runtime.getContextPointer
func GetContextPtr() uintptr

//go:linkname GetMainThreadHandle runtime.getMainThreadHandle
func GetMainThreadHandle() uintptr

//go:linkname GetHeapBase runtime.getHeapBase
func GetHeapBase() uintptr
