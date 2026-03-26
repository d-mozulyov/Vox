@echo off
wsl --cd "%~dp0" bash make.sh %*
setlocal enabledelayedexpansion
set "CMDLINE=!CMDCMDLINE!"
if "!CMDLINE:~-2,1!"==" " pause
endlocal
