package display

import (
	"github.com/racerxdl/gonx/nx/am"
	"github.com/racerxdl/gonx/nx/gpu"
	"github.com/racerxdl/gonx/nx/nxerrors"
	"github.com/racerxdl/gonx/nx/nxtypes"
	"github.com/racerxdl/gonx/nx/vi"
)

const debugDisplay = true

var display *vi.Display
var displayInitializations = 0
var displayInitializedAM = false

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
	err = am.Init()
	if err == nil {
		displayInitializedAM = true
	}

	return nil
}

func forceFinalize() {
	if debugDisplay {
		println("Display::ForceFinalize()")
	}
	if displayInitializedAM {
		am.Finalize()
		displayInitializedAM = false
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

	if !displayInitializedAM {
		// Homebrew override
		// TODO: vi_create_managed_layer
		if debugDisplay {
			println("Display::OpenLayer() - AM NOT INITIALIZED")
		}
		return surface, nxerrors.NotImplemented
	}

	if debugDisplay {
		println("Display::OpenLayer() - IwcAcquireForegroundRights")
	}
	err = am.IwcAcquireForegroundRights()
	if err != nil {
		return surface, err
	}

	if debugDisplay {
		println("Display::OpenLayer() - IwcGetAppletResourceUserId")
	}
	aruid, err = am.IwcGetAppletResourceUserId()
	if err != nil {
		return surface, err
	}

	if debugDisplay {
		println("Display::OpenLayer() - IscCreateManagedDisplayLayer")
	}
	layerId, err = am.IscCreateManagedDisplayLayer()
	if err != nil {
		return surface, err
	}
	managedLayerInit = true

	if debugDisplay {
		println("Display::OpenLayer() - OpenLayer")
	}
	igbp, err = vi.OpenLayer("Default", layerId, aruid)
	if err != nil {
		return surface, err
	}
	igbpInit = true

	if debugDisplay {
		println("Display::OpenLayer() - Surface Create")
	}
	surface, _, err = SurfaceCreate(layerId, *igbp)
	if err != nil {
		return surface, err
	}
	surfaceInit = true

	if debugDisplay {
		println("Display::OpenLayer() - IadsSetLayerScalingMode")
	}
	err = vi.IadsSetLayerScalingMode(2, layerId)
	if err != nil {
		return surface, err
	}

	//if !displayInitializedAM {
	//	// Homebrew, TODO
	//}

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
