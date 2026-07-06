#include "keyboard.h"

#include <stdexcept>
#include <cstring>

Keyboard::Keyboard()
    : Device{L"\\\\?\\HID#HIDRIVER&Col04#1"}
{
}

void Keyboard::initialize()
{
    if (isInitialized())
        throw std::runtime_error{"ERROR_DOUBLE_INITIALIZATION"};
    Device::initialize();
}

void Keyboard::type(BYTE keyCode)
{
    BYTE keyCodes[6] = {KEY_NONE};

    keyCodes[0] = keyCode;
    sendKeyboardReport(keyCodes);

    keyCodes[0] = KEY_NONE;
    sendKeyboardReport(keyCodes);
}

// In keyboard.cpp
void Keyboard::keyDown(BYTE keyCode)
{
    switch (keyCode)
    {
    case KEY_LEFTCTRL:
    case KEY_RIGHTCTRL:
        m_modifiers |= (keyCode == KEY_LEFTCTRL) ? KEY_MOD_LCTRL : KEY_MOD_RCTRL;
        break;
    case KEY_LEFTSHIFT:
    case KEY_RIGHTSHIFT:
        m_modifiers |= (keyCode == KEY_LEFTSHIFT) ? KEY_MOD_LSHIFT : KEY_MOD_RSHIFT;
        break;
    default:
        BYTE keyCodes[6] = {KEY_NONE};
        keyCodes[0] = keyCode;
        sendKeyboardReport(keyCodes);
        break;
    }
}

void Keyboard::keyUp(BYTE keyCode)
{
    switch (keyCode)
    {
    case KEY_LEFTCTRL:
    case KEY_RIGHTCTRL:
        m_modifiers &= ~((keyCode == KEY_LEFTCTRL) ? KEY_MOD_LCTRL : KEY_MOD_RCTRL);
        break;
    case KEY_LEFTSHIFT:
    case KEY_RIGHTSHIFT:
        m_modifiers &= ~((keyCode == KEY_LEFTSHIFT) ? KEY_MOD_LSHIFT : KEY_MOD_RSHIFT);
        break;
    // Add cases for other modifiers if necessary
    default:
        // For non-modifier keys, you might want to keep track
        // of which keys are currently pressed and release only those
        break;
    }
    sendKeyboardReport(nullptr);
}

void Keyboard::abort()
{
    Device::abort();
}

void Keyboard::sendKeyboardReport(BYTE *keyCodes)
{
    Report report;
    report.reportId = REPORT_ID;
    report.modifiers = m_modifiers;
    report._reserved = 0x00;
    if (keyCodes != nullptr)
    {
        memcpy(report.keyCodes, keyCodes, 6);
    }
    else
    {
        memset(report.keyCodes, KEY_NONE, 6);
    }

    setOutputReport(&report, static_cast<DWORD>(sizeof(Report)));
}
