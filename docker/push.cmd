@echo off
wsl --cd "%~dp0" bash push.sh
setlocal enabledelayedexpansion
set "CMDLINE=!CMDCMDLINE!"
if "!CMDLINE:~-2,1!"==" " pause
endlocal
