package display

import (
	"fmt"
	"github.com/racerxdl/gonx/nx/nxerrors"
	"github.com/racerxdl/gonx/nx/nxtypes"
	"github.com/racerxdl/gonx/services/am"
	"github.com/racerxdl/gonx/services/gpu"
	"github.com/racerxdl/gonx/services/vi"
)

const debugDisplay = false

var display *vi.Display
var displayInitializations = 0

func Init() (err error) {
	if debugDisplay {
		println("Display::Init()")
	}
	displayInitializations++
	if displayInitializations > 1 {
		return nil
	}

	gpuInit := false
	viInit := false

	defer func() {
		if err != nil {
			displayInitializations--
			if viInit {
				vi.Finalize()
			}
			if gpuInit {
				gpu.Finalize()
			}
		}
	}()

	if debugDisplay {
		println("Display::Init() - GPU Init")
	}
	err = gpu.Init()
	if err != nil {
		return err
	}
	gpuInit = true

	if debugDisplay {
		println("Display::Init() - VI Init")
	}
	err = vi.Init()
	if err != nil {
		return err
	}
	viInit = true

	if debugDisplay {
		println("Display::Init() - Open Display")
	}
	display, err = vi.OpenDisplay("Default")
	if err != nil {
		return err
	}

	if debugDisplay {
		println("Display::Init() - AM Init")
	}

	return nil
}

func forceFinalize() {
	if debugDisplay {
		println("Display::ForceFinalize()")
	}

	_ = vi.CloseDisplay(display)
	vi.Finalize()
	gpu.Finalize()
	displayInitializations = 0
}

func Finalize() {
	if debugDisplay {
		println("Display::Finalize()")
	}
	displayInitializations--
	if displayInitializations <= 0 {
		forceFinalize()
	}
}

func OpenLayer() (surface *Surface, err error) {
	if debugDisplay {
		println("Display::OpenLayer()")
	}
	if displayInitializations < 1 {
		return nil, nxerrors.DisplayNotInitialized
	}

	var layerId uint64
	var aruid nxtypes.ARUID
	var igbp *vi.IGBP

	managedLayerInit := false
	igbpInit := false
	surfaceInit := false
	layerId = 0

	defer func() {
		if err != nil {
			if surfaceInit {
				// Surface takes ownership of IGBP and layer
				surface.Destroy()
				surface = nil
				return
			}

			if igbpInit {
				_ = vi.AdjustRefCount(igbp.IgbpBinder.Handle, -1, 1)
				_ = vi.CloseLayer(layerId)
			}

			if managedLayerInit {
				_ = vi.DestroyManagedLayer(layerId)
			}
		}
	}()

	aruid, err = am.IwcGetAppletResourceUserId()
	if err != nil {
		return surface, err
	}

	if debugDisplay {
		println("Display::OpenLayer() - CreateManagedLayer")
	}
	if aruid > 0 {
		// Applet
		layerId, err = am.IscCreateManagedDisplayLayer()
		if err != nil {
			return surface, err
		}
	} else {
		layerId, err = vi.CreateManagedLayer(display, 0, aruid)
		if err != nil {
			return surface, err
		}
	}
	managedLayerInit = true

	if debugDisplay {
		fmt.Printf("Display::OpenLayer() - OpenLayer(\"Default\", %d, %d)\n", layerId, aruid)
	}
	igbp, err = vi.OpenLayer("Default", layerId, aruid)
	if err != nil {
		return surface, err
	}
	igbpInit = true

	if debugDisplay {
		println("Display::OpenLayer() - SurfaceCreate")
	}
	surface, _, err = SurfaceCreate(layerId, *igbp)
	if err != nil {
		return surface, err
	}
	surfaceInit = true

	if debugDisplay {
		println("Display::OpenLayer() - IadsSetLayerScalingMode")
	}
	err = vi.IadsSetLayerScalingMode(vi.ScalingMode_FitToLayer, layerId)
	if err != nil {
		return surface, err
	}

	return surface, nil
}

func GetVSyncEvent() (nxtypes.ReventHandle, error) {
	if displayInitializations <= 0 {
		return 0, nxerrors.DisplayNotInitialized
	}

	if display.VSync == 0 {
		err := vi.GetDisplayVsyncEvent(display)
		if err != nil {
			return 0, err
		}
	}

	return display.VSync, nil
}
