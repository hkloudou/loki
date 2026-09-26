//go:build !windows

package loki

// This file provides no-op stubs so the package compiles on non-Windows
// platforms (for editors, CI, and cross-platform builds). Every entry point
// returns ErrUnsupported.

// Mouse is unavailable on this platform.
type Mouse struct{}

// Keyboard is unavailable on this platform.
type Keyboard struct{}

// Version reports ErrUnsupported on non-Windows platforms.
func Version() (int, error) { return 0, ErrUnsupported }

// NewMouse returns ErrUnsupported on non-Windows platforms.
func NewMouse() (*Mouse, error) { return nil, ErrUnsupported }

func (m *Mouse) Move(x, y int) error           { return ErrUnsupported }
func (m *Mouse) MoveRelative(dx, dy int) error { return ErrUnsupported }
func (m *Mouse) LeftDown() error               { return ErrUnsupported }
func (m *Mouse) LeftUp() error                 { return ErrUnsupported }
func (m *Mouse) LeftClick() error              { return ErrUnsupported }
func (m *Mouse) RightDown() error              { return ErrUnsupported }
func (m *Mouse) RightUp() error                { return ErrUnsupported }
func (m *Mouse) RightClick() error             { return ErrUnsupported }
func (m *Mouse) MiddleDown() error             { return ErrUnsupported }
func (m *Mouse) MiddleUp() error               { return ErrUnsupported }
func (m *Mouse) MiddleClick() error            { return ErrUnsupported }
func (m *Mouse) Close() error                  { return nil }

// NewKeyboard returns ErrUnsupported on non-Windows platforms.
func NewKeyboard() (*Keyboard, error) { return nil, ErrUnsupported }

func (k *Keyboard) SetModifiers(mods byte) error       { return ErrUnsupported }
func (k *Keyboard) Type(key byte) error                { return ErrUnsupported }
func (k *Keyboard) Send(mods byte, keys ...byte) error { return ErrUnsupported }
func (k *Keyboard) Close() error                       { return nil }
