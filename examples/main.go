// Example: drive the virtual mouse and keyboard.
//
// Prerequisites:
//   - The loki virtual HID driver is installed (see ../windows/driver).
//   - loki.dll is loadable (next to this exe, on PATH, or via LOKI_DLL).
//
// Build & run on Windows:
//
//	set LOKI_DLL=C:\path\to\loki.dll
//	go run ./examples
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hkloudou/loki"
)

func main() {
	v, err := loki.Version()
	if err != nil {
		log.Fatalf("loki not available: %v", err)
	}
	fmt.Printf("loki.dll C-API version: %d\n", v)

	mouse, err := loki.NewMouse()
	if err != nil {
		log.Fatalf("mouse: %v", err)
	}
	defer mouse.Close()

	kb, err := loki.NewKeyboard()
	if err != nil {
		log.Fatalf("keyboard: %v", err)
	}
	defer kb.Close()

	// Give yourself a moment to focus a target window.
	time.Sleep(2 * time.Second)

	// Move to an absolute position and left-click.
	if err := mouse.Move(400, 300); err != nil {
		log.Fatalf("move: %v", err)
	}
	if err := mouse.LeftClick(); err != nil {
		log.Fatalf("click: %v", err)
	}

	// Type "Hi" — capital H via left-shift held for one report.
	if err := kb.Send(loki.ModLshift, loki.KeyH); err != nil {
		log.Fatalf("send: %v", err)
	}
	if err := kb.Send(0); err != nil { // release
		log.Fatalf("release: %v", err)
	}
	if err := kb.Type(loki.KeyI); err != nil {
		log.Fatalf("type: %v", err)
	}

	fmt.Println("done")
}
