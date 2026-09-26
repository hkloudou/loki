// Code generated from windows/lib/keys.h. DO NOT EDIT by hand.
// HID keyboard usage codes and modifier bitmasks (USB HID Usage Tables).
// Regenerate with: go generate ./...  (see gen_keys.sh) or edit keys.h + rerun.

package loki

// Modifier bitmask values for Keyboard.SetModifiers / Keyboard.Send.
const (
	ModLctrl  = 0x01
	ModLshift = 0x02
	ModLalt   = 0x04
	ModLmeta  = 0x08
	ModRctrl  = 0x10
	ModRshift = 0x20
	ModRalt   = 0x40
	ModRmeta  = 0x80
)

// HID keyboard usage codes for Keyboard.Type / Keyboard.Send.
const (
	KeyNone       = 0x00
	KeyA          = 0x04
	KeyB          = 0x05
	KeyC          = 0x06
	KeyD          = 0x07
	KeyE          = 0x08
	KeyF          = 0x09
	KeyG          = 0x0a
	KeyH          = 0x0b
	KeyI          = 0x0c
	KeyJ          = 0x0d
	KeyK          = 0x0e
	KeyL          = 0x0f
	KeyM          = 0x10
	KeyN          = 0x11
	KeyO          = 0x12
	KeyP          = 0x13
	KeyQ          = 0x14
	KeyR          = 0x15
	KeyS          = 0x16
	KeyT          = 0x17
	KeyU          = 0x18
	KeyV          = 0x19
	KeyW          = 0x1a
	KeyX          = 0x1b
	KeyY          = 0x1c
	KeyZ          = 0x1d
	Key1          = 0x1e
	Key2          = 0x1f
	Key3          = 0x20
	Key4          = 0x21
	Key5          = 0x22
	Key6          = 0x23
	Key7          = 0x24
	Key8          = 0x25
	Key9          = 0x26
	Key0          = 0x27
	KeyEnter      = 0x28
	KeyEsc        = 0x29
	KeyBackspace  = 0x2a
	KeyTab        = 0x2b
	KeySpace      = 0x2c
	KeyMinus      = 0x2d
	KeyEqual      = 0x2e
	KeyLeftbrace  = 0x2f
	KeyRightbrace = 0x30
	KeyBackslash  = 0x31
	KeyHashtilde  = 0x32
	KeySemicolon  = 0x33
	KeyApostrophe = 0x34
	KeyGrave      = 0x35
	KeyComma      = 0x36
	KeyDot        = 0x37
	KeySlash      = 0x38
	KeyCapslock   = 0x39
	KeyF1         = 0x3a
	KeyF2         = 0x3b
	KeyF3         = 0x3c
	KeyF4         = 0x3d
	KeyF5         = 0x3e
	KeyF6         = 0x3f
	KeyF7         = 0x40
	KeyF8         = 0x41
	KeyF9         = 0x42
	KeyF10        = 0x43
	KeyF11        = 0x44
	KeyF12        = 0x45
	KeySysrq      = 0x46
	KeyScrolllock = 0x47
	KeyPause      = 0x48
	KeyInsert     = 0x49
	KeyHome       = 0x4a
	KeyPageup     = 0x4b
	KeyDelete     = 0x4c
	KeyEnd        = 0x4d
	KeyPagedown   = 0x4e
	KeyRight      = 0x4f
	KeyLeft       = 0x50
	KeyDown       = 0x51
	KeyUp         = 0x52
	KeyNumlock    = 0x53
	KeyKpslash    = 0x54
	KeyKpasterisk = 0x55
	KeyKpminus    = 0x56
	KeyKpplus     = 0x57
	KeyKpenter    = 0x58
	KeyKp1        = 0x59
	KeyKp2        = 0x5a
	KeyKp3        = 0x5b
	KeyKp4        = 0x5c
	KeyKp5        = 0x5d
	KeyKp6        = 0x5e
	KeyKp7        = 0x5f
	KeyKp8        = 0x60
	KeyKp9        = 0x61
	KeyKp0        = 0x62
	KeyKpdot      = 0x63
	Key102nd      = 0x64
	KeyCompose    = 0x65
	KeyLeftctrl   = 0xe0
	KeyLeftshift  = 0xe1
	KeyLeftalt    = 0xe2
	KeyLeftmeta   = 0xe3
	KeyRightctrl  = 0xe4
	KeyRightshift = 0xe5
	KeyRightalt   = 0xe6
	KeyRightmeta  = 0xe7
)
