@echo off
REM Build loki.dll from a "x64 Native Tools Command Prompt for VS".
REM Usage: build.bat [Debug|Release]   (default: Release)
setlocal
set CFG=%1
if "%CFG%"=="" set CFG=Release

set OUT=%~dp0..\..\build\x64\%CFG%
if not exist "%OUT%" mkdir "%OUT%"

set FLAGS=/nologo /std:c++17 /EHsc /W3 /DUNICODE /D_UNICODE /D_WINDOWS /D_USRDLL /LD
if /I "%CFG%"=="Debug" (
  set FLAGS=%FLAGS% /Zi /MDd /Od
) else (
  set FLAGS=%FLAGS% /MD /O2
)

pushd "%~dp0"
cl %FLAGS% ^
  device.cpp mouse.cpp keyboard.cpp registry.cpp loki_c_api.cpp ^
  /Fe:"%OUT%\loki.dll" ^
  /link user32.lib cfgmgr32.lib advapi32.lib hid.lib
set ERR=%ERRORLEVEL%
popd

if %ERR% NEQ 0 (
  echo Build FAILED with code %ERR%
  exit /b %ERR%
)
echo Built %OUT%\loki.dll
endlocal
