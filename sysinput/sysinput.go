// Package sysinput is the driverless (B1) input backend: it injects mouse and
// keyboard input with the Win32 SendInput API from an agent process running
// inside the interactive user session.
//
// It is deliberately separate from the root loki package (the signed/virtual
// HID *driver* backend). Use this one when you already have a service running as
// LocalSystem and want to drive the real desktop — including the secure desktop
// (UAC prompts, the logon screen, Ctrl+Alt+Del) — without shipping a kernel
// driver. Trade-off: SendInput sets the LLMHF_INJECTED flag, so some
// games/anti-cheat ignore or detect it; for those, use the driver or an external
// USB-HID board instead.
//
// Topology:
//
//	LocalSystem service (Session 0)          user session (Session N)
//	┌─────────────────────────┐   stdin pipe  ┌──────────────────────────┐
//	│ LaunchAgentInActiveSession├─────────────►│ RunAgent: OpenInputDesktop│
//	│ AgentProcess.MoveTo/Key…  │   JSON lines  │ + SetThreadDesktop +      │
//	│                          │               │ SendInput (follows secure │
//	│                          │               │ desktop, runs as SYSTEM)  │
//	└─────────────────────────┘               └──────────────────────────┘
//
// The service links this package and calls LaunchAgentInActiveSession; the agent
// binary (cmd/loki-inputagent) calls RunAgent.
package sysinput

import (
	"encoding/json"
	"errors"
	"io"
)

// ErrUnsupported is returned by all entry points on non-Windows platforms.
var ErrUnsupported = errors.New("sysinput: only supported on Windows")

// Button is a mouse button bitmask value.
type Button int

const (
	Left   Button = 1
	Right  Button = 2
	Middle Button = 4
)

// wire command types.
const (
	cmdMoveTo  = "mt" // absolute move, pixels in virtual-screen space
	cmdMoveRel = "mr" // relative move
	cmdButton  = "b"  // mouse button up/down
	cmdWheel   = "w"  // mouse wheel
	cmdKey     = "k"  // key up/down by VK or scancode
	cmdText    = "t"  // type a unicode string
)

// command is the JSON message sent from the service to the agent, one per line.
type command struct {
	T    string `json:"t"`
	X    int    `json:"x,omitempty"`
	Y    int    `json:"y,omitempty"`
	B    int    `json:"b,omitempty"`
	Down bool   `json:"d,omitempty"`
	VK   uint16 `json:"vk,omitempty"`
	Scan uint16 `json:"sc,omitempty"`
	Amt  int    `json:"a,omitempty"`
	S    string `json:"s,omitempty"`
}

// writeCommand encodes one command as a JSON line.
func writeCommand(w io.Writer, c command) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	b = append(b, '\n')
	_, err = w.Write(b)
	return err
}
