package display

import (
	"github.com/racerxdl/gonx/nx/am"
	"github.com/racerxdl/gonx/nx/gpu"
	"github.com/racerxdl/gonx/nx/vi"
)

var display *vi.Display
var displayInitializations = 0
var displayInitializedAM = false

func Init() (err error) {
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

	display, err = vi.OpenDisplay("Default")
	if err != nil {
		return err
	}

	err = am.Init()
	if err == nil {
		displayInitializedAM = true
	}

	return nil
}

func forceFinalize() {
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
	displayInitializations--
	if displayInitializations <= 0 {
		forceFinalize()
	}
}
