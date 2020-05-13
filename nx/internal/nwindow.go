// +build nintendoswitch

package internal

import "unsafe"

//go:export nwindowGetDefault
func NWindowGetDefault() unsafe.Pointer
