package svc

import _ "unsafe" // for go:linkname

//go:linkname GetContextPtr runtime.getContextPtr
func GetContextPtr() uintptr

//go:linkname GetMainThreadHandle runtime.getMainThreadHandle
func GetMainThreadHandle() uintptr

//go:linkname GetHeapBase runtime.getHeapBase
func GetHeapBase() uintptr
