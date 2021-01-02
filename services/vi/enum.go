package vi

type ScalingMode uint32

const (
	ScalingMode_None                ScalingMode = 0x0
	ScalingMode_FitToLayer          ScalingMode = 0x2
	ScalingMode_PreserveAspectRatio ScalingMode = 0x4
	ScalingModeDefault                          = ScalingMode_FitToLayer
)
