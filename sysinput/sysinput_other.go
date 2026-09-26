//go:build !windows

package sysinput

import (
	"io"
	"time"
)

// AgentProcess is unavailable on non-Windows platforms.
type AgentProcess struct{}

// ActiveConsoleSessionID returns 0xFFFFFFFF on non-Windows platforms.
func ActiveConsoleSessionID() uint32 { return 0xFFFFFFFF }

// LaunchAgentInActiveSession returns ErrUnsupported on non-Windows platforms.
func LaunchAgentInActiveSession(agentPath string, extraArgs ...string) (*AgentProcess, error) {
	return nil, ErrUnsupported
}

// RunAgent returns ErrUnsupported on non-Windows platforms.
func RunAgent(r io.Reader) error { return ErrUnsupported }

func (p *AgentProcess) PID() uint32                   { return 0 }
func (p *AgentProcess) Alive() bool                   { return false }
func (p *AgentProcess) MoveTo(x, y int) error         { return ErrUnsupported }
func (p *AgentProcess) Move(dx, dy int) error         { return ErrUnsupported }
func (p *AgentProcess) ButtonDown(b Button) error     { return ErrUnsupported }
func (p *AgentProcess) ButtonUp(b Button) error       { return ErrUnsupported }
func (p *AgentProcess) Click(b Button) error          { return ErrUnsupported }
func (p *AgentProcess) Wheel(notches int) error       { return ErrUnsupported }
func (p *AgentProcess) KeyDown(vk, scan uint16) error { return ErrUnsupported }
func (p *AgentProcess) KeyUp(vk, scan uint16) error   { return ErrUnsupported }
func (p *AgentProcess) KeyTap(vk uint16) error        { return ErrUnsupported }
func (p *AgentProcess) Type(s string) error           { return ErrUnsupported }
func (p *AgentProcess) Close() error                  { return nil }

// Manager is unavailable on non-Windows platforms.
type Manager struct{}

// NewManager returns a no-op Manager on non-Windows platforms.
func NewManager(agentPath string, agentArgs ...string) *Manager { return &Manager{} }

func (m *Manager) SetPollInterval(d time.Duration) {}
func (m *Manager) SetLogger(fn func(string))       {}
func (m *Manager) Start()                          {}
func (m *Manager) Stop()                           {}
func (m *Manager) Available() bool                 { return false }
func (m *Manager) MoveTo(x, y int) error           { return ErrUnsupported }
func (m *Manager) Move(dx, dy int) error           { return ErrUnsupported }
func (m *Manager) ButtonDown(b Button) error       { return ErrUnsupported }
func (m *Manager) ButtonUp(b Button) error         { return ErrUnsupported }
func (m *Manager) Click(b Button) error            { return ErrUnsupported }
func (m *Manager) Wheel(n int) error               { return ErrUnsupported }
func (m *Manager) KeyDown(vk, scan uint16) error   { return ErrUnsupported }
func (m *Manager) KeyUp(vk, scan uint16) error     { return ErrUnsupported }
func (m *Manager) KeyTap(vk uint16) error          { return ErrUnsupported }
func (m *Manager) Type(s string) error             { return ErrUnsupported }
