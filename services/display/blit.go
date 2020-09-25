package display

func pdep(mask, value uint32) uint32 {
	out := uint32(0)
	for shift := uint32(0); shift < 32; shift++ {
		bit := uint32(1 << shift)
		if mask&bit > 0 {
			if value&1 > 0 {
				out |= bit
			}
			value >>= 1
		}
	}
	return out
}

func swizzleX(v uint32) uint32 {
	return pdep(^uint32(0x7B4), v)
}

func swizzleY(v uint32) uint32 {
	return pdep(0x7B4, v)
}

const tileHeight = 128
const paddedWidth = tileHeight * 10

func setUint32(buffer []byte, idx uint32, v uint32) {
	idx *= 4          // Uint32 size
	_ = buffer[idx+3] // early bounds check to guarantee safety of writes below
	buffer[idx] = byte(v)
	buffer[idx+1] = byte(v >> 8)
	buffer[idx+2] = byte(v >> 16)
	buffer[idx+3] = byte(v >> 24)
}

func getUint32(buffer []byte, idx uint32) uint32 {
	idx *= 4          // Uint32 size
	_ = buffer[idx+3] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(buffer[idx+0]) | uint32(buffer[idx+1])<<8 | uint32(buffer[idx+2])<<16 | uint32(buffer[idx+3])<<24
}

func GFXSlowSwizzlingBlit(buffer []byte, image []byte, w, h, tx, ty int) {
	x0 := uint32(tx)
	y0 := uint32(ty)
	x1 := x0 + uint32(w)
	y1 := y0 + uint32(h)

	// we're doing this in pixels - should just shift the swizzles instead
	offsX0 := swizzleX(x0)
	offsY := swizzleY(y0)
	XMask := swizzleX(^uint32(0))
	YMask := swizzleY(^uint32(0))
	IncrY := swizzleX(paddedWidth)

	// step offs_x0 to the right row of tiles
	offsX0 += IncrY * (y0 / tileHeight)

	srcPos := uint32(0)

	for y := uint32(0); y < y1; y++ {
		offsX := offsX0

		for x := uint32(0); x < x1; x++ {
			pixel := getUint32(image, srcPos)
			srcPos++
			setUint32(buffer, offsY+offsX, pixel)
			offsX = (offsX - XMask) & XMask
		}

		offsY = (offsY - YMask) & YMask
		if offsY == 0 {
			offsX0 += IncrY // wrap into next tile row
		}
	}
}
