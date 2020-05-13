// +build nintendoswitch

package nx

type HidControllerId uint64
type HidControllerKeys uint64

// HidControllerID
const (
    ControllerPlayer1  HidControllerId = 0
    ControllerPlayer2  HidControllerId = 1
    ControllerPlayer3  HidControllerId = 2
    ControllerPlayer4  HidControllerId = 3
    ControllerPlayer5  HidControllerId = 4
    ControllerPlayer6  HidControllerId = 5
    ControllerPlayer7  HidControllerId = 6
    ControllerPlayer8  HidControllerId = 7
    ControllerHandheld HidControllerId = 8
    ControllerUnknown  HidControllerId = 9
    // Not an actual HID-sysmodule ID. Only for hidKeys*()/hidJoystickRead()/hidSixAxisSensorValuesRead()/hidGetControllerType()/hidGetControllerColors()/hidIsControllerConnected().
    // Automatically uses CONTROLLER_PLAYER_1 when connected, otherwise uses CONTROLLER_HANDHELD.
    ControllerP1Auto HidControllerId = 10
)

const (
    KeyA           HidControllerKeys = 1 << 0  // A
    KeyB           HidControllerKeys = 1 << 1  // B
    KeyX           HidControllerKeys = 1 << 2  // X
    KeyY           HidControllerKeys = 1 << 3  // Y
    KeyLStick      HidControllerKeys = 1 << 4  // Left Stick Button
    KeyRStick      HidControllerKeys = 1 << 5  // Right Stick Button
    KeyL           HidControllerKeys = 1 << 6  // L
    KeyR           HidControllerKeys = 1 << 7  // R
    KeyZL          HidControllerKeys = 1 << 8  // ZL
    KeyZR          HidControllerKeys = 1 << 9  // ZR
    KeyPlus        HidControllerKeys = 1 << 10 // Plus
    KeyMinus       HidControllerKeys = 1 << 11 // Minus
    KeyDLeft       HidControllerKeys = 1 << 12 // D-Pad Left
    KeyDup         HidControllerKeys = 1 << 13 // D-Pad Up
    KeyDRight      HidControllerKeys = 1 << 14 // D-Pad Right
    KeyDDown       HidControllerKeys = 1 << 15 // D-Pad Down
    KeyLStickLeft  HidControllerKeys = 1 << 16 // Left Stick Left
    KeyLStickUp    HidControllerKeys = 1 << 17 // Left Stick Up
    KeyLStickRight HidControllerKeys = 1 << 18 // Left Stick Right
    KeyLStickDown  HidControllerKeys = 1 << 19 // Left Stick Down
    KeyRStickLeft  HidControllerKeys = 1 << 20 // Right Stick Left
    KeyRStickUp    HidControllerKeys = 1 << 21 // Right Stick Up
    KeyRStickRight HidControllerKeys = 1 << 22 // Right Stick Right
    KeyRStickDown  HidControllerKeys = 1 << 23 // Right Stick Down
    KeySLLeft      HidControllerKeys = 1 << 24 // SL on Left Joy-Con
    KeySRLeft      HidControllerKeys = 1 << 25 // SR on Left Joy-Con
    KeySLRight     HidControllerKeys = 1 << 26 // SL on Right Joy-Con
    KeySRRight     HidControllerKeys = 1 << 27 // SR on Right Joy-Con

    KeyHome    HidControllerKeys = 1 << 18 // HOME button, only available for use with HiddbgHdlsState::buttons.
    KeyCapture HidControllerKeys = 1 << 19 // Capture button, only available for use with HiddbgHdlsState::buttons.

    // Pseudo-key for at least one finger on the touch screen
    KeyTouch HidControllerKeys = 1 << 28

    // Buttons by orientation (for single Joy-Con), also works with Joy-Con pairs, Pro Controller
    KeyJoyconRight HidControllerKeys = 1 << 0
    KeyJoyconDown  HidControllerKeys = 1 << 1
    KeyJoyconUp    HidControllerKeys = 1 << 2
    KeyJoyconLeft  HidControllerKeys = 1 << 3

    // Generic catch-all directions, also works for single Joy-Con
    KeyUp    = KeyDup | KeyLStickUp | KeyRStickUp          // D-Pad Up or Sticks Up
    KeyDown  = KeyDDown | KeyLStickDown | KeyRStickDown    // D-Pad Down or Sticks Down
    KeyLeft  = KeyDLeft | KeyLStickLeft | KeyRStickLeft    // D-Pad Left or Sticks Left
    KeyRight = KeyDRight | KeyLStickRight | KeyRStickRight // D-Pad Right or Sticks Right
    KeySl    = KeySLLeft | KeySLRight                      // SL on Left or Right Joy-Con
    KeySr    = KeySRLeft | KeySRRight                      // SR on Left or Right Joy-Con
)

func (k HidControllerKeys) IsKeyDown(key HidControllerKeys) bool {
    return k & key > 0
}

func (k HidControllerKeys) IsKeyUp(key HidControllerKeys) bool {
    return !k.IsKeyDown(key)
}

//go:export hidScanInput
func HidScanInput()

//u64 hidKeysDown(HidControllerID id);
//go:export hidKeysDown
func HidKeysDown(id HidControllerId) HidControllerKeys
