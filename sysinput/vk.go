package sysinput

// Common Windows Virtual-Key codes for KeyDown/KeyUp/KeyTap. This is the VK code
// space used by SendInput — distinct from the HID usage codes used by the driver
// backend's keys.go. Letters and digits equal their ASCII uppercase value
// (e.g. 'A' == VKA == 0x41), so you can also pass a rune like uint16('A').
const (
	VKBack     = 0x08
	VKTab      = 0x09
	VKReturn   = 0x0D
	VKShift    = 0x10
	VKControl  = 0x11
	VKMenu     = 0x12 // Alt
	VKPause    = 0x13
	VKCapital  = 0x14 // Caps Lock
	VKEscape   = 0x1B
	VKSpace    = 0x20
	VKPrior    = 0x21 // Page Up
	VKNext     = 0x22 // Page Down
	VKEnd      = 0x23
	VKHome     = 0x24
	VKLeft     = 0x25
	VKUp       = 0x26
	VKRight    = 0x27
	VKDown     = 0x28
	VKSnapshot = 0x2C // Print Screen
	VKInsert   = 0x2D
	VKDelete   = 0x2E

	VKLWin = 0x5B
	VKRWin = 0x5C

	VKF1  = 0x70
	VKF2  = 0x71
	VKF3  = 0x72
	VKF4  = 0x73
	VKF5  = 0x74
	VKF6  = 0x75
	VKF7  = 0x76
	VKF8  = 0x77
	VKF9  = 0x78
	VKF10 = 0x79
	VKF11 = 0x7A
	VKF12 = 0x7B

	VKLShift   = 0xA0
	VKRShift   = 0xA1
	VKLControl = 0xA2
	VKRControl = 0xA3
	VKLMenu    = 0xA4 // Left Alt
	VKRMenu    = 0xA5 // Right Alt
)
