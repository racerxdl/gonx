package env

import (
	"github.com/racerxdl/gonx/nx/nxtypes"
	"github.com/racerxdl/gonx/svc"
	"unsafe"
)

type configEntry struct {
	Key   uint32
	Flags uint32
	Value [2]uint64
}

var argvPtr uintptr
var syscallHints [2]uint64
var appletType = nxtypes.AppletTypeDefault
var processHandle uint64
var lastLoadResult uint64
var randomSeed [2]uint64
var userIdStorage uintptr
var hosVersion uint64

func LoadEnv() error {
	ctxPtr := svc.GetContextPtr()
	if ctxPtr != 0 {
		ptr := ctxPtr
		// See https://switchbrew.org/w/index.php?title=Homebrew_ABI
		entry := (*configEntry)(unsafe.Pointer(ptr))
		for entry.Key != envEntryTypeEndOfList {
			switch entry.Key {
			case envEntryTypeMainThreadHandle: // Handled by tinyGo
			case envEntryTypeNextLoadPath:
			case envEntryTypeOverrideHeap: // Handled by tinyGo
			case envEntryTypeOverrideService:
			case envEntryTypeArgv:
				argvPtr = uintptr(entry.Value[1])
			case envEntryTypeSyscallAvailableHint:
				syscallHints[0] = entry.Value[0]
				syscallHints[1] = entry.Value[1]
			case envEntryTypeAppletType:
				appletType = nxtypes.AppletType(entry.Value[0])
			case envEntryTypeAppletWorkaround:
			case envEntryTypeReserved9:
			case envEntryTypeProcessHandle:
				processHandle = entry.Value[0]
			case envEntryTypeLastLoadResult:
				lastLoadResult = entry.Value[0]
			case envEntryTypeRandomSeed:
				randomSeed[0] = entry.Value[0]
				randomSeed[1] = entry.Value[1]
			case envEntryTypeUserIdStorage:
				userIdStorage = uintptr(entry.Value[0])
			case envEntryTypeHosVersion:
				hosVersion = entry.Value[0]
			}
			ptr += unsafe.Sizeof(configEntry{})
			entry = (*configEntry)(unsafe.Pointer(ptr))
		}
	}

	return nil
}

// GetAppletType returns the Applet Type reported by homebrew launcher
func GetAppletType() nxtypes.AppletType {
	return nxtypes.AppletType(appletType)
}
