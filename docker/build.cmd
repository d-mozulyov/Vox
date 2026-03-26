@echo off
wsl --cd "%~dp0" bash build.sh
setlocal enabledelayedexpansion
set "CMDLINE=!CMDCMDLINE!"
if "!CMDLINE:~-2,1!"==" " pause
endlocal
