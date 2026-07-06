@echo off
"C:\Program Files (x86)\Windows Kits\10\bin\10.0.22621.0\x64\stampinf.exe"  -f n\hidriver.inf -d "*" -v "1.2.0.1"
@REM -d 01/01/2022 -v 1.0.0

del "n\hidriver.cat"
"C:\Program Files (x86)\Windows Kits\10\bin\10.0.22621.0\x86\Inf2Cat.exe" /driver:.\n /os:10_X64

"C:\Program Files (x86)\Windows Kits\10\bin\10.0.20348.0\x64\Signtool.exe" sign /v /fd SHA256 /a /f .\cert.pfx /p password .\n\hidriver.sys

"C:\Program Files (x86)\Windows Kits\10\bin\10.0.20348.0\x64\Signtool.exe" sign /v /fd SHA256 /a /f .\cert.pfx /p password .\n\hidriver.cat

@REM "C:\Program Files (x86)\Windows Kits\10\Tools\10.0.22621.0\x64\devcon.exe" remove "root\hidriver"

@REM "C:\Program Files (x86)\Windows Kits\10\Tools\10.0.22621.0\x64\devcon.exe" find "root\hidriver" && "C:\Program Files (x86)\Windows Kits\10\Tools\10.0.22621.0\x64\devcon.exe" remove "root\hidriver" && "C:\Program Files (x86)\Windows Kits\10\Tools\10.0.22621.0\x64\devcon.exe" -r install .\hidriver.inf "root\hidriver"

pause