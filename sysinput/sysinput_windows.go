//go:build windows

package sysinput

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"syscall"
	"unsafe"
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	advapi32 = syscall.NewLazyDLL("advapi32.dll")

	procSendInput        = user32.NewProc("SendInput")
	procOpenInputDesktop = user32.NewProc("OpenInputDesktop")
	procSetThreadDesktop = user32.NewProc("SetThreadDesktop")
	procCloseDesktop     = user32.NewProc("CloseDesktop")
	procGetSystemMetrics = user32.NewProc("GetSystemMetrics")

	procWTSActiveSession     = kernel32.NewProc("WTSGetActiveConsoleSessionId")
	procCreatePipe           = kernel32.NewProc("CreatePipe")
	procSetHandleInfo        = kernel32.NewProc("SetHandleInformation")
	procCloseHandle          = kernel32.NewProc("CloseHandle")
	procGetCurrentProcess    = kernel32.NewProc("GetCurrentProcess")
	procProcessIdToSessionId = kernel32.NewProc("ProcessIdToSessionId")

	procOpenProcessToken     = advapi32.NewProc("OpenProcessToken")
	procDuplicateTokenEx     = advapi32.NewProc("DuplicateTokenEx")
	procSetTokenInformation  = advapi32.NewProc("SetTokenInformation")
	procCreateProcessAsUserW = advapi32.NewProc("CreateProcessAsUserW")
)

// Win32 constants.
const (
	inputMouseType    = 0
	inputKeyboardType = 1

	mouseeventfMove        = 0x0001
	mouseeventfLeftDown    = 0x0002
	mouseeventfLeftUp      = 0x0004
	mouseeventfRightDown   = 0x0008
	mouseeventfRightUp     = 0x0010
	mouseeventfMiddleDown  = 0x0020
	mouseeventfMiddleUp    = 0x0040
	mouseeventfWheel       = 0x0800
	mouseeventfAbsolute    = 0x8000
	mouseeventfVirtualDesk = 0x4000

	keyeventfExtended = 0x0001
	keyeventfKeyUp    = 0x0002
	keyeventfUnicode  = 0x0004
	keyeventfScancode = 0x0008

	smXVirtualScreen  = 76
	smYVirtualScreen  = 77
	smCXVirtualScreen = 78
	smCYVirtualScreen = 79

	tokenDuplicate        = 0x0002
	tokenQuery            = 0x0008
	tokenAssignPrimary    = 0x0001
	tokenAdjustDefault    = 0x0080
	tokenAdjustSession    = 0x0100
	securityImpersonation = 2
	tokenPrimary          = 1
	tokenSessionID        = 12
	maximumAllowed        = 0x02000000

	createUnicodeEnvironment = 0x00000400
	createNoWindow           = 0x08000000
	startfUseStdHandles      = 0x00000100
	handleFlagInherit        = 0x00000001

	wheelDelta = 120
)

type inputMouse struct {
	typ         uint32
	_           uint32
	dx          int32
	dy          int32
	mouseData   uint32
	dwFlags     uint32
	time        uint32
	_           uint32
	dwExtraInfo uintptr
}

type inputKbd struct {
	typ         uint32
	_           uint32
	wVk         uint16
	wScan       uint16
	dwFlags     uint32
	time        uint32
	_           uint32
	dwExtraInfo uintptr
	_tail       uint64
}

type startupInfoW struct {
	cb            uint32
	lpReserved    *uint16
	lpDesktop     *uint16
	lpTitle       *uint16
	dwX           uint32
	dwY           uint32
	dwXSize       uint32
	dwYSize       uint32
	dwXCountChars uint32
	dwYCountChars uint32
	dwFillAttr    uint32
	dwFlags       uint32
	wShowWindow   uint16
	cbReserved2   uint16
	lpReserved2   *byte
	hStdInput     syscall.Handle
	hStdOutput    syscall.Handle
	hStdError     syscall.Handle
}

type processInformation struct {
	hProcess    syscall.Handle
	hThread     syscall.Handle
	dwProcessID uint32
	dwThreadID  uint32
}

type securityAttributes struct {
	length             uint32
	securityDescriptor uintptr
	inheritHandle      int32
}

// ==========================================================================
// Service side (Session 0): launch and drive the agent.
// ==========================================================================

// AgentProcess is a handle, held by the LocalSystem service, to an input agent
// running in the interactive session. Methods are safe for concurrent use.
type AgentProcess struct {
	mu    sync.Mutex
	stdin *os.File
	proc  syscall.Handle
	pid   uint32
}

// PID returns the agent process id.
func (p *AgentProcess) PID() uint32 { return p.pid }

// ActiveConsoleSessionID returns the session id of the physical console, or
// 0xFFFFFFFF when no user is logged on.
func ActiveConsoleSessionID() uint32 {
	r, _, _ := procWTSActiveSession.Call()
	return uint32(r)
}

// LaunchAgentInActiveSession starts agentPath inside the active console session,
// as SYSTEM (so it can follow the secure desktop), with its stdin wired to a
// pipe this process keeps. Call it from a process already running as LocalSystem.
// extraArgs are passed after the exe path on the agent command line.
func LaunchAgentInActiveSession(agentPath string, extraArgs ...string) (*AgentProcess, error) {
	session := ActiveConsoleSessionID()
	if session == 0xFFFFFFFF {
		return nil, fmt.Errorf("sysinput: no active console session (no user logged on)")
	}

	hProc, _, _ := procGetCurrentProcess.Call()

	var hTok syscall.Handle
	if r, _, err := procOpenProcessToken.Call(hProc,
		uintptr(tokenDuplicate|tokenQuery|tokenAssignPrimary|tokenAdjustDefault|tokenAdjustSession),
		uintptr(unsafe.Pointer(&hTok))); r == 0 {
		return nil, fmt.Errorf("sysinput: OpenProcessToken: %w (are you running as LocalSystem?)", err)
	}
	defer procCloseHandle.Call(uintptr(hTok))

	var hDup syscall.Handle
	if r, _, err := procDuplicateTokenEx.Call(uintptr(hTok), maximumAllowed, 0,
		securityImpersonation, tokenPrimary, uintptr(unsafe.Pointer(&hDup))); r == 0 {
		return nil, fmt.Errorf("sysinput: DuplicateTokenEx: %w", err)
	}
	defer procCloseHandle.Call(uintptr(hDup))

	sid := session
	if r, _, err := procSetTokenInformation.Call(uintptr(hDup), tokenSessionID,
		uintptr(unsafe.Pointer(&sid)), unsafe.Sizeof(sid)); r == 0 {
		return nil, fmt.Errorf("sysinput: SetTokenInformation(session=%d): %w", session, err)
	}

	// Command pipe: child reads stdin, we keep the write end (non-inheritable).
	sa := securityAttributes{length: uint32(unsafe.Sizeof(securityAttributes{})), inheritHandle: 1}
	var rd, wr syscall.Handle
	if r, _, err := procCreatePipe.Call(uintptr(unsafe.Pointer(&rd)),
		uintptr(unsafe.Pointer(&wr)), uintptr(unsafe.Pointer(&sa)), 0); r == 0 {
		return nil, fmt.Errorf("sysinput: CreatePipe: %w", err)
	}
	procSetHandleInfo.Call(uintptr(wr), handleFlagInherit, 0)

	desktop, err := syscall.UTF16PtrFromString(`winsta0\default`)
	if err != nil {
		return nil, err
	}
	cmdline := quoteArg(agentPath)
	for _, a := range extraArgs {
		cmdline += " " + quoteArg(a)
	}
	cmdlinePtr, err := syscall.UTF16PtrFromString(cmdline)
	if err != nil {
		return nil, err
	}

	si := startupInfoW{}
	si.cb = uint32(unsafe.Sizeof(si))
	si.lpDesktop = desktop
	si.dwFlags = startfUseStdHandles
	si.hStdInput = rd
	var pi processInformation

	r, _, cerr := procCreateProcessAsUserW.Call(
		uintptr(hDup),
		0, // lpApplicationName (use cmdline)
		uintptr(unsafe.Pointer(cmdlinePtr)),
		0, 0, // process/thread attrs
		1, // bInheritHandles
		uintptr(createUnicodeEnvironment|createNoWindow),
		0, // environment (inherit) — see note in package doc
		0, // current dir
		uintptr(unsafe.Pointer(&si)),
		uintptr(unsafe.Pointer(&pi)),
	)
	// Child owns the read end now.
	procCloseHandle.Call(uintptr(rd))
	if r == 0 {
		procCloseHandle.Call(uintptr(wr))
		return nil, fmt.Errorf("sysinput: CreateProcessAsUser: %w", cerr)
	}
	procCloseHandle.Call(uintptr(pi.hThread))

	return &AgentProcess{
		stdin: os.NewFile(uintptr(wr), "loki-agent-stdin"),
		proc:  pi.hProcess,
		pid:   pi.dwProcessID,
	}, nil
}

func (p *AgentProcess) send(c command) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.stdin == nil {
		return fmt.Errorf("sysinput: agent closed")
	}
	return writeCommand(p.stdin, c)
}

// MoveTo moves the cursor to absolute virtual-screen pixel coordinates.
func (p *AgentProcess) MoveTo(x, y int) error { return p.send(command{T: cmdMoveTo, X: x, Y: y}) }

// Move nudges the cursor by a relative pixel delta.
func (p *AgentProcess) Move(dx, dy int) error { return p.send(command{T: cmdMoveRel, X: dx, Y: dy}) }

// ButtonDown / ButtonUp press or release a mouse button.
func (p *AgentProcess) ButtonDown(b Button) error {
	return p.send(command{T: cmdButton, B: int(b), Down: true})
}
func (p *AgentProcess) ButtonUp(b Button) error {
	return p.send(command{T: cmdButton, B: int(b)})
}

// Click presses and releases a mouse button.
func (p *AgentProcess) Click(b Button) error {
	if err := p.ButtonDown(b); err != nil {
		return err
	}
	return p.ButtonUp(b)
}

// Wheel scrolls by notches (positive = up/away from user).
func (p *AgentProcess) Wheel(notches int) error { return p.send(command{T: cmdWheel, Amt: notches}) }

// KeyDown / KeyUp press or release a key. Provide a virtual-key code in vk, or
// pass vk=0 and a hardware scan code in scan to inject by scancode.
func (p *AgentProcess) KeyDown(vk, scan uint16) error {
	return p.send(command{T: cmdKey, VK: vk, Scan: scan, Down: true})
}
func (p *AgentProcess) KeyUp(vk, scan uint16) error {
	return p.send(command{T: cmdKey, VK: vk, Scan: scan})
}

// KeyTap presses and releases a virtual-key code.
func (p *AgentProcess) KeyTap(vk uint16) error {
	if err := p.KeyDown(vk, 0); err != nil {
		return err
	}
	return p.KeyUp(vk, 0)
}

// Type sends a Unicode string as keyboard input (layout-independent).
func (p *AgentProcess) Type(s string) error { return p.send(command{T: cmdText, S: s}) }

// Close shuts the command channel (the agent exits on EOF) and releases the
// process handle. It does not force-kill the agent.
func (p *AgentProcess) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.stdin != nil {
		p.stdin.Close()
		p.stdin = nil
	}
	if p.proc != 0 {
		procCloseHandle.Call(uintptr(p.proc))
		p.proc = 0
	}
	return nil
}

func quoteArg(s string) string { return `"` + s + `"` }

// ==========================================================================
// Agent side (interactive session): read commands and SendInput.
// ==========================================================================

// RunAgent reads newline-delimited JSON commands from r and injects them with
// SendInput on whatever desktop currently has input focus (following the user
// onto the secure desktop). It locks the OS thread for the lifetime of the loop
// because SetThreadDesktop and SendInput are thread-affine. Call this from the
// agent binary (see cmd/loki-inputagent), passing os.Stdin.
func RunAgent(r io.Reader) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var curDesk uintptr
	ensureDesktop := func() {
		h, _, _ := procOpenInputDesktop.Call(0, 0, maximumAllowed)
		if h == 0 {
			return
		}
		if ok, _, _ := procSetThreadDesktop.Call(h); ok == 0 {
			procCloseDesktop.Call(h)
			return
		}
		if curDesk != 0 {
			procCloseDesktop.Call(curDesk)
		}
		curDesk = h
	}

	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var c command
		if err := json.Unmarshal(line, &c); err != nil {
			continue // skip malformed lines rather than dying
		}
		ensureDesktop()
		applyCommand(c)
	}
	return sc.Err()
}

func applyCommand(c command) {
	switch c.T {
	case cmdMoveTo:
		nx, ny := toAbsolute(c.X, c.Y)
		sendMouse(mouseeventfMove|mouseeventfAbsolute|mouseeventfVirtualDesk, nx, ny, 0)
	case cmdMoveRel:
		sendMouse(mouseeventfMove, int32(c.X), int32(c.Y), 0)
	case cmdButton:
		down, up := buttonFlags(Button(c.B))
		if c.Down {
			sendMouse(down, 0, 0, 0)
		} else {
			sendMouse(up, 0, 0, 0)
		}
	case cmdWheel:
		sendMouse(mouseeventfWheel, 0, 0, uint32(int32(c.Amt*wheelDelta)))
	case cmdKey:
		flags := uint32(0)
		vk := c.VK
		scan := c.Scan
		if vk == 0 && scan != 0 {
			flags |= keyeventfScancode
		}
		if !c.Down {
			flags |= keyeventfKeyUp
		}
		sendKey(vk, scan, flags)
	case cmdText:
		for _, r := range c.S {
			typeRune(r)
		}
	}
}

func toAbsolute(x, y int) (int32, int32) {
	xv, _, _ := procGetSystemMetrics.Call(smXVirtualScreen)
	yv, _, _ := procGetSystemMetrics.Call(smYVirtualScreen)
	cx, _, _ := procGetSystemMetrics.Call(smCXVirtualScreen)
	cy, _, _ := procGetSystemMetrics.Call(smCYVirtualScreen)
	cxv := int(int32(cx))
	cyv := int(int32(cy))
	if cxv < 2 {
		cxv = 2
	}
	if cyv < 2 {
		cyv = 2
	}
	nx := int64(x-int(int32(xv))) * 65535 / int64(cxv-1)
	ny := int64(y-int(int32(yv))) * 65535 / int64(cyv-1)
	return int32(nx), int32(ny)
}

func buttonFlags(b Button) (down, up uint32) {
	switch b {
	case Right:
		return mouseeventfRightDown, mouseeventfRightUp
	case Middle:
		return mouseeventfMiddleDown, mouseeventfMiddleUp
	default:
		return mouseeventfLeftDown, mouseeventfLeftUp
	}
}

func sendMouse(flags uint32, dx, dy int32, data uint32) {
	in := inputMouse{typ: inputMouseType, dx: dx, dy: dy, mouseData: data, dwFlags: flags}
	procSendInput.Call(1, uintptr(unsafe.Pointer(&in)), unsafe.Sizeof(in))
}

func sendKey(vk, scan uint16, flags uint32) {
	in := inputKbd{typ: inputKeyboardType, wVk: vk, wScan: scan, dwFlags: flags}
	procSendInput.Call(1, uintptr(unsafe.Pointer(&in)), unsafe.Sizeof(in))
}

// typeRune injects one Unicode code point (surrogate pair for astral planes).
func typeRune(r rune) {
	if r > 0xFFFF {
		r -= 0x10000
		hi := uint16(0xD800 + (r >> 10))
		lo := uint16(0xDC00 + (r & 0x3FF))
		sendKey(0, hi, keyeventfUnicode)
		sendKey(0, hi, keyeventfUnicode|keyeventfKeyUp)
		sendKey(0, lo, keyeventfUnicode)
		sendKey(0, lo, keyeventfUnicode|keyeventfKeyUp)
		return
	}
	sendKey(0, uint16(r), keyeventfUnicode)
	sendKey(0, uint16(r), keyeventfUnicode|keyeventfKeyUp)
}
