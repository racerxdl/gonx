// +build nintendoswitch

package nx

import "unsafe"

const (
    consoleEsc     = "\x1b["
    ConsoleReset   = consoleEsc + "0m"
    ConsoleBlack   = consoleEsc + "30m"
    ConsoleRed     = consoleEsc + "31;1m"
    ConsoleGreen   = consoleEsc + "32;1m"
    ConsoleYellow  = consoleEsc + "33;1m"
    ConsoleBlue    = consoleEsc + "34;1m"
    ConsoleMagenta = consoleEsc + "35;1m"
    ConsoleCyan    = consoleEsc + "36;1m"
    ConsoleWhite   = consoleEsc + "37;1m"
)

func ConsoleESC(text string) string {
    return consoleEsc + text
}


//go:export consoleInit
func ConsoleInit(c unsafe.Pointer)

//go:export consoleUpdate
func ConsoleUpdate(c unsafe.Pointer)

//go:export consoleExit
func ConsoleExit(c unsafe.Pointer)
