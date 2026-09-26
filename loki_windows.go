//go:build windows

package loki

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"syscall"
	"unsafe"
)

// dll is loaded lazily on first use. The name can be overridden with the
// LOKI_DLL environment variable (absolute path recommended); otherwise the
// standard Windows DLL search order applies to "loki.dll".
var (
	dllOnce sync.Once
	dll     *syscall.LazyDLL

	procLastError *syscall.LazyProc
	procVersion   *syscall.LazyProc

	procMouseCreate     *syscall.LazyProc
	procMouseInit       *syscall.LazyProc
	procMouseMove       *syscall.LazyProc
	procMouseMoveRel    *syscall.LazyProc
	procMouseLeftDown   *syscall.LazyProc
	procMouseLeftUp     *syscall.LazyProc
	procMouseLeftClick  *syscall.LazyProc
	procMouseRightDown  *syscall.LazyProc
	procMouseRightUp    *syscall.LazyProc
	procMouseRightClick *syscall.LazyProc
	procMouseMidDown    *syscall.LazyProc
	procMouseMidUp      *syscall.LazyProc
	procMouseMidClick   *syscall.LazyProc
	procMouseDestroy    *syscall.LazyProc

	procKbCreate  *syscall.LazyProc
	procKbInit    *syscall.LazyProc
	procKbSetMods *syscall.LazyProc
	procKbType    *syscall.LazyProc
	procKbSend    *syscall.LazyProc
	procKbDestroy *syscall.LazyProc
)

func loadDLL() {
	dllOnce.Do(func() {
		name := os.Getenv("LOKI_DLL")
		if name == "" {
			name = "loki.dll"
		}
		dll = syscall.NewLazyDLL(name)

		procLastError = dll.NewProc("loki_last_error")
		procVersion = dll.NewProc("loki_version")

		procMouseCreate = dll.NewProc("loki_mouse_create")
		procMouseInit = dll.NewProc("loki_mouse_initialize")
		procMouseMove = dll.NewProc("loki_mouse_move")
		procMouseMoveRel = dll.NewProc("loki_mouse_move_relative")
		procMouseLeftDown = dll.NewProc("loki_mouse_left_down")
		procMouseLeftUp = dll.NewProc("loki_mouse_left_up")
		procMouseLeftClick = dll.NewProc("loki_mouse_left_click")
		procMouseRightDown = dll.NewProc("loki_mouse_right_down")
		procMouseRightUp = dll.NewProc("loki_mouse_right_up")
		procMouseRightClick = dll.NewProc("loki_mouse_right_click")
		procMouseMidDown = dll.NewProc("loki_mouse_middle_down")
		procMouseMidUp = dll.NewProc("loki_mouse_middle_up")
		procMouseMidClick = dll.NewProc("loki_mouse_middle_click")
		procMouseDestroy = dll.NewProc("loki_mouse_destroy")

		procKbCreate = dll.NewProc("loki_keyboard_create")
		procKbInit = dll.NewProc("loki_keyboard_initialize")
		procKbSetMods = dll.NewProc("loki_keyboard_set_modifiers")
		procKbType = dll.NewProc("loki_keyboard_type")
		procKbSend = dll.NewProc("loki_keyboard_send")
		procKbDestroy = dll.NewProc("loki_keyboard_destroy")
	})
}

// lastError copies the DLL's thread-local error message into a local buffer.
// The DLL stores it in thread_local storage, so this read is only guaranteed to
// match the immediately preceding failing call when both run on the same OS
// thread; for robust messages callers can wrap a sequence in
// runtime.LockOSThread. Copying by value here avoids any pointer-lifetime issue.
func lastError() string {
	if procLastError == nil {
		return ""
	}
	var buf [256]byte
	n, _, _ := procLastError.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return string(buf[:n])
}

// statusErr turns a non-zero C status into a Go error, attaching the DLL's
// last-error message when available.
func statusErr(op string, status uintptr) error {
	if status == 0 {
		return nil
	}
	if msg := lastError(); msg != "" {
		return fmt.Errorf("loki: %s failed: %s", op, msg)
	}
	return fmt.Errorf("loki: %s failed (status %d)", op, status)
}

// Version returns the loki.dll C-API version, or an error if the DLL cannot be
// loaded. Useful as a quick "is everything wired up" probe.
func Version() (int, error) {
	loadDLL()
	if err := dll.Load(); err != nil {
		return 0, fmt.Errorf("loki: cannot load DLL: %w", err)
	}
	r, _, _ := procVersion.Call()
	return int(r), nil
}

// ---- Mouse ---------------------------------------------------------------

// Mouse is a handle to the virtual mouse. It is not safe for concurrent use;
// guard it with your own mutex if multiple goroutines drive one Mouse.
type Mouse struct {
	h uintptr
}

// NewMouse creates and initializes the virtual mouse. It fails if the driver is
// not installed or the DLL cannot be loaded.
func NewMouse() (*Mouse, error) {
	loadDLL()
	if err := dll.Load(); err != nil {
		return nil, fmt.Errorf("loki: cannot load DLL: %w", err)
	}
	h, _, _ := procMouseCreate.Call()
	if h == 0 {
		return nil, fmt.Errorf("loki: mouse create failed: %s", lastError())
	}
	m := &Mouse{h: h}
	if status, _, _ := procMouseInit.Call(h); status != 0 {
		err := statusErr("mouse initialize", status)
		procMouseDestroy.Call(h)
		m.h = 0
		return nil, err
	}
	return m, nil
}

// Move walks the cursor towards absolute screen coordinates (x, y).
func (m *Mouse) Move(x, y int) error {
	s, _, _ := procMouseMove.Call(m.h, uintptr(int32(x)), uintptr(int32(y)))
	return statusErr("mouse move", s)
}

// MoveRelative sends one report with signed deltas, each clamped to [-127, 127].
func (m *Mouse) MoveRelative(dx, dy int) error {
	s, _, _ := procMouseMoveRel.Call(m.h, uintptr(int32(dx)), uintptr(int32(dy)))
	return statusErr("mouse move relative", s)
}

func (m *Mouse) LeftDown() error {
	s, _, _ := procMouseLeftDown.Call(m.h)
	return statusErr("left down", s)
}
func (m *Mouse) LeftUp() error { s, _, _ := procMouseLeftUp.Call(m.h); return statusErr("left up", s) }
func (m *Mouse) LeftClick() error {
	s, _, _ := procMouseLeftClick.Call(m.h)
	return statusErr("left click", s)
}
func (m *Mouse) RightDown() error {
	s, _, _ := procMouseRightDown.Call(m.h)
	return statusErr("right down", s)
}
func (m *Mouse) RightUp() error {
	s, _, _ := procMouseRightUp.Call(m.h)
	return statusErr("right up", s)
}
func (m *Mouse) RightClick() error {
	s, _, _ := procMouseRightClick.Call(m.h)
	return statusErr("right click", s)
}
func (m *Mouse) MiddleDown() error {
	s, _, _ := procMouseMidDown.Call(m.h)
	return statusErr("middle down", s)
}
func (m *Mouse) MiddleUp() error {
	s, _, _ := procMouseMidUp.Call(m.h)
	return statusErr("middle up", s)
}
func (m *Mouse) MiddleClick() error {
	s, _, _ := procMouseMidClick.Call(m.h)
	return statusErr("middle click", s)
}

// Close releases the mouse handle. The Mouse must not be used afterwards.
func (m *Mouse) Close() error {
	if m.h != 0 {
		procMouseDestroy.Call(m.h)
		m.h = 0
	}
	return nil
}

// ---- Keyboard ------------------------------------------------------------

// Keyboard is a handle to the virtual keyboard. Not safe for concurrent use.
type Keyboard struct {
	h uintptr
}

// NewKeyboard creates and initializes the virtual keyboard.
func NewKeyboard() (*Keyboard, error) {
	loadDLL()
	if err := dll.Load(); err != nil {
		return nil, fmt.Errorf("loki: cannot load DLL: %w", err)
	}
	h, _, _ := procKbCreate.Call()
	if h == 0 {
		return nil, fmt.Errorf("loki: keyboard create failed: %s", lastError())
	}
	k := &Keyboard{h: h}
	if status, _, _ := procKbInit.Call(h); status != 0 {
		err := statusErr("keyboard initialize", status)
		procKbDestroy.Call(h)
		k.h = 0
		return nil, err
	}
	return k, nil
}

// SetModifiers sets the modifier bitmask (Mod* constants) applied to subsequent
// reports until changed again.
func (k *Keyboard) SetModifiers(mods byte) error {
	s, _, _ := procKbSetMods.Call(k.h, uintptr(mods))
	return statusErr("keyboard set modifiers", s)
}

// Type presses and releases a single HID usage code (a Key* constant).
func (k *Keyboard) Type(key byte) error {
	s, _, _ := procKbType.Call(k.h, uintptr(key))
	return statusErr("keyboard type", s)
}

// Send emits one report with the given modifier bitmask and up to 6 key codes
// held simultaneously. Call Send(0) to release all keys.
func (k *Keyboard) Send(mods byte, keys ...byte) error {
	if len(keys) > 6 {
		keys = keys[:6]
	}
	var ptr uintptr
	if len(keys) > 0 {
		ptr = uintptr(unsafe.Pointer(&keys[0]))
	}
	s, _, _ := procKbSend.Call(k.h, uintptr(mods), ptr, uintptr(len(keys)))
	// keys must stay alive across the syscall above so the GC does not move or
	// collect the backing array while the DLL reads from ptr.
	runtime.KeepAlive(keys)
	return statusErr("keyboard send", s)
}

// Close releases the keyboard handle.
func (k *Keyboard) Close() error {
	if k.h != 0 {
		procKbDestroy.Call(k.h)
		k.h = 0
	}
	return nil
}
