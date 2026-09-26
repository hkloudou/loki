// loki_c_api.cpp — implementation of the C ABI wrapper. See loki_c_api.h.
#include "pch.h"
#include "framework.h"

#include "loki_c_api.h"

#include "mouse.h"
#include "keyboard.h"

#include <cstring>
#include <exception>
#include <string>

namespace {

// Per-thread last error message so concurrent callers don't clobber each other.
thread_local std::string g_lastError;

void setError(const char *msg) { g_lastError = msg ? msg : ""; }
void clearError() { g_lastError.clear(); }

// Runs a callable, translating C++ exceptions into a status code + message.
template <typename Fn>
int guard(Fn &&fn) {
    try {
        fn();
        clearError();
        return LOKI_OK;
    } catch (const std::exception &e) {
        setError(e.what());
        return LOKI_ERR_EXCEPTION;
    } catch (...) {
        setError("unknown error");
        return LOKI_ERR_EXCEPTION;
    }
}

} // namespace

extern "C" {

int loki_last_error(char *buf, int len) {
    if (!buf || len <= 0) return 0;
    int n = (int)g_lastError.size();
    if (n > len - 1) n = len - 1;
    if (n > 0) memcpy(buf, g_lastError.data(), (size_t)n);
    buf[n] = '\0';
    return n;
}

int loki_version(void) {
    // 0.1.0
    return 100;
}

// ---- Mouse ---------------------------------------------------------------

loki_mouse_t loki_mouse_create(void) {
    try {
        return static_cast<loki_mouse_t>(new Mouse());
    } catch (...) {
        setError("failed to allocate Mouse");
        return nullptr;
    }
}

int loki_mouse_initialize(loki_mouse_t m) {
    if (!m) { setError("null mouse handle"); return LOKI_ERR_NULL_HANDLE; }
    return guard([&] { static_cast<Mouse *>(m)->initialize(); });
}

int loki_mouse_move(loki_mouse_t m, long x, long y) {
    if (!m) { setError("null mouse handle"); return LOKI_ERR_NULL_HANDLE; }
    return guard([&] { static_cast<Mouse *>(m)->moveCursor((LONG)x, (LONG)y); });
}

int loki_mouse_move_relative(loki_mouse_t m, int dx, int dy) {
    if (!m) { setError("null mouse handle"); return LOKI_ERR_NULL_HANDLE; }
    if (dx > 127) dx = 127; else if (dx < -127) dx = -127;
    if (dy > 127) dy = 127; else if (dy < -127) dy = -127;
    return guard([&] {
        static_cast<Mouse *>(m)->sendMouseReport((CHAR)dx, (CHAR)dy);
    });
}

int loki_mouse_left_down(loki_mouse_t m) {
    if (!m) { setError("null mouse handle"); return LOKI_ERR_NULL_HANDLE; }
    return guard([&] { static_cast<Mouse *>(m)->leftButtonDown(); });
}
int loki_mouse_left_up(loki_mouse_t m) {
    if (!m) { setError("null mouse handle"); return LOKI_ERR_NULL_HANDLE; }
    return guard([&] { static_cast<Mouse *>(m)->leftButtonUp(); });
}
int loki_mouse_left_click(loki_mouse_t m) {
    if (!m) { setError("null mouse handle"); return LOKI_ERR_NULL_HANDLE; }
    return guard([&] { static_cast<Mouse *>(m)->leftButtonClick(); });
}

int loki_mouse_right_down(loki_mouse_t m) {
    if (!m) { setError("null mouse handle"); return LOKI_ERR_NULL_HANDLE; }
    return guard([&] { static_cast<Mouse *>(m)->rightButtonDown(); });
}
int loki_mouse_right_up(loki_mouse_t m) {
    if (!m) { setError("null mouse handle"); return LOKI_ERR_NULL_HANDLE; }
    return guard([&] { static_cast<Mouse *>(m)->rightButtonUp(); });
}
int loki_mouse_right_click(loki_mouse_t m) {
    if (!m) { setError("null mouse handle"); return LOKI_ERR_NULL_HANDLE; }
    return guard([&] { static_cast<Mouse *>(m)->rightButtonClick(); });
}

int loki_mouse_middle_down(loki_mouse_t m) {
    if (!m) { setError("null mouse handle"); return LOKI_ERR_NULL_HANDLE; }
    return guard([&] { static_cast<Mouse *>(m)->middleButtonDown(); });
}
int loki_mouse_middle_up(loki_mouse_t m) {
    if (!m) { setError("null mouse handle"); return LOKI_ERR_NULL_HANDLE; }
    return guard([&] { static_cast<Mouse *>(m)->middleButtonUp(); });
}
int loki_mouse_middle_click(loki_mouse_t m) {
    if (!m) { setError("null mouse handle"); return LOKI_ERR_NULL_HANDLE; }
    return guard([&] { static_cast<Mouse *>(m)->middleButtonClick(); });
}

void loki_mouse_destroy(loki_mouse_t m) {
    if (!m) return;
    Mouse *mouse = static_cast<Mouse *>(m);
    try { mouse->abort(); } catch (...) {}
    delete mouse;
}

// ---- Keyboard ------------------------------------------------------------

loki_keyboard_t loki_keyboard_create(void) {
    try {
        return static_cast<loki_keyboard_t>(new Keyboard());
    } catch (...) {
        setError("failed to allocate Keyboard");
        return nullptr;
    }
}

int loki_keyboard_initialize(loki_keyboard_t k) {
    if (!k) { setError("null keyboard handle"); return LOKI_ERR_NULL_HANDLE; }
    return guard([&] { static_cast<Keyboard *>(k)->initialize(); });
}

int loki_keyboard_set_modifiers(loki_keyboard_t k, unsigned char modifiers) {
    if (!k) { setError("null keyboard handle"); return LOKI_ERR_NULL_HANDLE; }
    return guard([&] { static_cast<Keyboard *>(k)->setModifiers((BYTE)modifiers); });
}

int loki_keyboard_type(loki_keyboard_t k, unsigned char keyCode) {
    if (!k) { setError("null keyboard handle"); return LOKI_ERR_NULL_HANDLE; }
    return guard([&] { static_cast<Keyboard *>(k)->type((BYTE)keyCode); });
}

int loki_keyboard_send(loki_keyboard_t k, unsigned char modifiers,
                       const unsigned char *keyCodes, int count) {
    if (!k) { setError("null keyboard handle"); return LOKI_ERR_NULL_HANDLE; }
    return guard([&] {
        static_cast<Keyboard *>(k)->send((BYTE)modifiers, (const BYTE *)keyCodes, count);
    });
}

void loki_keyboard_destroy(loki_keyboard_t k) {
    if (!k) return;
    Keyboard *keyboard = static_cast<Keyboard *>(k);
    try { keyboard->abort(); } catch (...) {}
    delete keyboard;
}

} // extern "C"
