# loki virtual HID driver

A KMDF kernel-mode driver that registers a virtual HID device with a **mouse**
collection and a **keyboard** collection. User mode writes HID *output* reports
(via `loki.dll`); the driver forwards them as HID *input* reports, so Windows
delivers real input events.

The C sources are unchanged from
[hedgar2017/loki-hidriver](https://github.com/hedgar2017/loki-hidriver) /
[dengqizhou30/HIDDriver](https://github.com/dengqizhou30/HIDDriver); see
[../../docs/comparison.md](../../docs/comparison.md).

## Build

Requires **Visual Studio** + the **Windows Driver Kit (WDK)** matching your VS
version. Build `loki-driver.vcxproj` as `Release|x64`. Output: `hidriver.sys`,
`hidriver.inf`, and (when signed) `hidriver.cat`.

## Install (developer / test mode)

The default build is **test-signed only**, so it needs signature enforcement
relaxed. In an **elevated** command prompt:

```bat
bcdedit /set testsigning on
:: reboot
devcon install hidriver.inf root\hidriver
devcon status root\hidriver
```

Remove / disable:

```bat
devcon remove root\hidriver
bcdedit /set testsigning off
```

Install logs: `C:\Windows\INF\setupapi.dev.log`.

## Shipping to real users

Test mode is **not** shippable. To load on normal machines the driver must be
properly signed — for a software driver like this that means **Microsoft
attestation signing** via Partner Center, using an **EV code-signing
certificate**. See [../../docs/signing.md](../../docs/signing.md), which also
covers the trade-offs vs. going driverless (SYSTEM service + `SendInput`, or an
external USB-HID microcontroller) the way ToDesk / 影刀 / UiPath do.

## Notes

- `VENDOR_ID` / `PRODUCT_ID` / `VERSION_NUMBER` are all `0x00` in `device.h`.
- The device is discovered from user mode by the string `HID#HIDRIVER&Col02`
  (mouse) / `&Col04` (keyboard), which is derived from `[Strings] ClassName` in
  `hidriver.inf`. Keep the INF `ClassName` and the lib's device paths in sync.
