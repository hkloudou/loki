// loki_c_api.h
//
// C ABI wrapper around the C++ Mouse / Keyboard classes so that the driver can
// be driven from any language with a plain C FFI (Go, Rust, Python, C#, ...).
//
// The whole surface is `extern "C"` and uses only primitive types + opaque
// handles, so no C++ runtime detail leaks across the boundary. Every call that
// can fail returns an int status: 0 == LOKI_OK, non-zero == error. Use
// loki_last_error() to retrieve a human readable message for the current
// thread.
//
// Build target: loki.dll
#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#if defined(_WIN32)
#define LOKI_API __declspec(dllexport)
#else
#define LOKI_API
#endif

typedef void *loki_mouse_t;
typedef void *loki_keyboard_t;

// Status codes.
enum {
    LOKI_OK = 0,
    LOKI_ERR_NULL_HANDLE = 1,
    LOKI_ERR_NOT_INITIALIZED = 2,
    LOKI_ERR_EXCEPTION = 3,
};

// Copies the last error message set on the current thread into buf (NUL
// terminated, truncated to len). Returns the number of bytes written excluding
// the terminator. Pass a buffer of a few hundred bytes; messages are short.
LOKI_API int loki_last_error(char *buf, int len);

// Returns the C-API version as an integer (major*10000 + minor*100 + patch).
LOKI_API int loki_version(void);

// ---- Mouse ---------------------------------------------------------------

// Allocates a Mouse object. Never returns NULL under normal conditions.
LOKI_API loki_mouse_t loki_mouse_create(void);

// Opens the virtual mouse HID collection. Must be called before any action.
LOKI_API int loki_mouse_initialize(loki_mouse_t m);

// Absolute move: walks the cursor towards (x, y) in screen coordinates, the
// same behaviour as the original loki example.
LOKI_API int loki_mouse_move(loki_mouse_t m, long x, long y);

// Relative nudge: sends a single HID report with the given signed deltas
// (clamped to [-127, 127] by the HID report format).
LOKI_API int loki_mouse_move_relative(loki_mouse_t m, int dx, int dy);

LOKI_API int loki_mouse_left_down(loki_mouse_t m);
LOKI_API int loki_mouse_left_up(loki_mouse_t m);
LOKI_API int loki_mouse_left_click(loki_mouse_t m);

LOKI_API int loki_mouse_right_down(loki_mouse_t m);
LOKI_API int loki_mouse_right_up(loki_mouse_t m);
LOKI_API int loki_mouse_right_click(loki_mouse_t m);

LOKI_API int loki_mouse_middle_down(loki_mouse_t m);
LOKI_API int loki_mouse_middle_up(loki_mouse_t m);
LOKI_API int loki_mouse_middle_click(loki_mouse_t m);

// Releases the object and closes the device handle.
LOKI_API void loki_mouse_destroy(loki_mouse_t m);

// ---- Keyboard ------------------------------------------------------------

LOKI_API loki_keyboard_t loki_keyboard_create(void);
LOKI_API int loki_keyboard_initialize(loki_keyboard_t k);

// Sets the current modifier byte (bitmask of KEY_MOD_* from keys.h). It is
// applied to every following report until changed again.
LOKI_API int loki_keyboard_set_modifiers(loki_keyboard_t k, unsigned char modifiers);

// Presses then releases a single HID usage code (a KEY_* value from keys.h).
LOKI_API int loki_keyboard_type(loki_keyboard_t k, unsigned char keyCode);

// Sends up to 6 simultaneous key codes plus a modifier byte. Pass count 0 with
// modifiers 0 to release everything. keyCodes may be NULL when count is 0.
LOKI_API int loki_keyboard_send(loki_keyboard_t k, unsigned char modifiers,
                                const unsigned char *keyCodes, int count);

LOKI_API void loki_keyboard_destroy(loki_keyboard_t k);

#ifdef __cplusplus
}
#endif
