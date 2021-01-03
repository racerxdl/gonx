package nxtypes

// AppletType
// https://github.com/switchbrew/libnx/blob/master/nx/include/switch/services/applet.h
type AppletType int64

const (
	AppletTypeNone              AppletType = -2
	AppletTypeDefault           AppletType = -1
	AppletTypeApplication       AppletType = 0
	AppletTypeSystemApplet      AppletType = 1
	AppletTypeLibraryApplet     AppletType = 2
	AppletTypeOverlayApplet     AppletType = 3
	AppletTypeSystemApplication AppletType = 4
)
