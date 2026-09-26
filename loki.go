// Package loki drives a virtual HID mouse and keyboard on Windows through the
// loki KMDF driver, via the loki.dll C wrapper.
//
// The driver presents a fake HID mouse+keyboard; writing HID output reports to
// it makes Windows deliver real input events, so automation is indistinguishable
// from a physical device at the OS input layer. This is the same technique used
// by remote-desktop and RPA tools that need input that games and anti-cheat /
// anti-bot layers treat as genuine hardware.
//
// Requirements at runtime:
//   - The loki virtual HID driver must be installed (see windows/driver).
//   - loki.dll (built from windows/lib) must be loadable: either next to the
//     executable, on the DLL search path, or pointed to by the LOKI_DLL env var.
//
// All functionality is Windows-only. On other platforms every constructor
// returns ErrUnsupported so code still compiles for cross-platform builds.
package loki

import "errors"

// ErrUnsupported is returned by constructors on non-Windows platforms.
var ErrUnsupported = errors.New("loki: only supported on Windows")

//go:generate ./gen_keys.sh
