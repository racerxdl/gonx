package env

// Types of config Entry.
// https://github.com/switchbrew/libnx/blob/master/nx/include/switch/runtime/env.h#L24-L40
const (
	envEntryTypeEndOfList            = 0  // Entry list terminator.
	envEntryTypeMainThreadHandle     = 1  // Provides the handle to the main thread.
	envEntryTypeNextLoadPath         = 2  // Provides a buffer containing information about the next homebrew application to load.
	envEntryTypeOverrideHeap         = 3  // Provides heap override information.
	envEntryTypeOverrideService      = 4  // Provides service override information.
	envEntryTypeArgv                 = 5  // Provides argv.
	envEntryTypeSyscallAvailableHint = 6  // Provides syscall availability hints.
	envEntryTypeAppletType           = 7  // Provides APT applet type.
	envEntryTypeAppletWorkaround     = 8  // Indicates that APT is broken and should not be used.
	envEntryTypeReserved9            = 9  // Unused/reserved entry type, formerly used by StdioSockets.
	envEntryTypeProcessHandle        = 10 // Provides the process handle.
	envEntryTypeLastLoadResult       = 11 // Provides the last load result.
	envEntryTypeRandomSeed           = 14 // Provides random data used to seed the pseudo-random number generator.
	envEntryTypeUserIdStorage        = 15 // Provides persistent storage for the preselected user id.
	envEntryTypeHosVersion           = 16 // Provides the currently running Horizon OS version.
)
