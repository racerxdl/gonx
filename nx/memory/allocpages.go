package memory

import "unsafe"

const pageSize = 0x1000

// AllocPages will alocate a page-aligned region of memory that is at least min bytes long and no more than max bytes long
// Returns null if can't allocate a page aligned byte slice
// By default Nintendo Switch page size is 4KB
func AllocPages(min, max int) (page []byte) {
	pageSizeMinus1 := uintptr(pageSize - 1)
	slice := make([]byte, max+int(pageSizeMinus1))
	ptr := uintptr(unsafe.Pointer(&slice[0]))
	slice = slice[pageSizeMinus1-(ptr+pageSizeMinus1)%pageSize:]
	if len(slice) > max {
		slice = slice[:max]
	}

	if len(slice) < min {
		return nil
	}

	return slice
}
