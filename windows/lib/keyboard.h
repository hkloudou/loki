#pragma once

#include "device.h"

#include <Windows.h>

#include "keys.h"

class Keyboard : public Device
{

public:
    explicit Keyboard();
    Keyboard(const Keyboard &) = delete;
    void operator =(const Keyboard &) = delete;
    virtual ~Keyboard() override = default;

    virtual void                        initialize() override;

    virtual void                        type(BYTE keyCode);

    // Sets the modifier byte applied to every subsequent report.
    virtual void                        setModifiers(BYTE modifiers);

    // Sends a raw report: a modifier bitmask plus up to 6 key codes. Pass
    // count 0 (with modifiers 0) to release all keys. Codes beyond index 5
    // are ignored.
    virtual void                        send(BYTE modifiers, const BYTE *keyCodes, int count);

    virtual void                        abort() override;

protected:
    BYTE                                m_modifiers             {KEY_NONE};

    static const BYTE                   REPORT_ID               {0x04};

private:
#pragma pack(1)
    struct Report {
        BYTE                            reportId;
        BYTE                            modifiers;
        BYTE                            _reserved;
        BYTE                            keyCodes[6];
    };
#pragma pack()

    virtual void                        sendKeyboardReport(BYTE *keyCodes);
};
