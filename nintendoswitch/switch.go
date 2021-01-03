package nintendoswitch

import (
	"fmt"
	"github.com/racerxdl/gonx/nx/env"
	"github.com/racerxdl/gonx/services/am"
	"github.com/racerxdl/gonx/services/display"
	"github.com/racerxdl/gonx/services/gpu"
	"github.com/racerxdl/gonx/services/sm"
	"github.com/racerxdl/gonx/services/vi"
	"github.com/racerxdl/gonx/svc"
)

// Configure initializes all necessary services for gonx works
// Make sure you call Cleanup before exiting the program or it may crash the HorizonOS
func Configure() (err error) {
	x0 := svc.GetContextPtr()
	x1 := svc.GetMainThreadHandle()

	if x0 != 0 && x1 == (1<<64-1) {
		fmt.Println("TODO: exception to handle")
		svc.Break(0x5EF0D30, 0x1234, 0x4321)
		select {}
	}

	smInit := false
	amInit := false
	gpuInit := false
	viInit := false
	displayInit := false

	defer func() {
		if err != nil {
			if displayInit {
				display.Finalize()
			}
			if viInit {
				vi.Finalize()
			}
			if gpuInit {
				gpu.Finalize()
			}
			if amInit {
				am.Finalize()
			}
			if smInit {
				sm.Finalize()
			}
		}
	}()

	err = env.LoadEnv()
	if err != nil {
		return err
	}

	err = sm.Init()
	if err != nil {
		return err
	}
	smInit = true

	err = am.Init()
	if err != nil {
		return err
	}
	amInit = true

	err = gpu.Init()
	if err != nil {
		return err
	}
	gpuInit = true

	err = vi.Init()
	if err != nil {
		return err
	}
	viInit = true

	err = display.Init()
	if err != nil {
		return err
	}
	displayInit = true

	return nil
}

// Cleanup cleans up all resources created by Configure
// It is important to cleanup resources otherwise Horizon might crash!
func Cleanup() {
	display.Finalize()
	vi.Finalize()
	gpu.Finalize()
	sm.Finalize()
}
