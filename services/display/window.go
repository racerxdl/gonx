package display

import (
	"github.com/racerxdl/gonx/nx/nxtypes"
	"github.com/racerxdl/gonx/services/vi"
)

const (
	// From Android window.h.
	/* attributes queriable with query() */
	NativeWindowWidth  = 0
	NativeWindowHeight = 1
	NativeWindowFormat = 2
)

// From Android window.h.
/* parameter for NATIVE_WINDOW_[API_][DIS]CONNECT */
//...
/* Buffers will be queued after being filled using the CPU
 */
const NativeWindowAPICPU = 2

type BqRect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

type NativeWindow struct {
	Magic                   uint32
	bq                      vi.Binder
	event                   nxtypes.Event
	SlotsConfigured         uint64
	SlotsRequested          uint64
	CurSlot                 int32
	Width                   uint32
	Height                  uint32
	Format                  uint32
	Usage                   uint32
	Crop                    BqRect
	ScalingMode             uint32
	Transform               uint32
	StickyTransform         uint32
	DefaultWidth            uint32
	DefaultHeight           uint32
	SwapInterval            uint32
	IsConnected             bool
	ProducerControlledByApp bool
	ConsumerRunningBehind   bool
}
