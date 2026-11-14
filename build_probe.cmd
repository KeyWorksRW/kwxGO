@echo off
setlocal
rem Probe lane for kwxGO: the same cgo environment and link line as build_quick.cmd,
rem with the package (tests\quick_probe\gen) and the output name (probe.exe) changed.
rem Keeping one environment definition is deliberate -- a second copy of the link
rem line is a second source of truth for the build.
rem
rem Probe output never goes to tests\quick.
rem
rem   build_probe.cmd        build tests\quick_probe\gen (the probe gate)
rem   build_probe.cmd help   this text
call "%~dp0build_quick.cmd" probe
exit /b %errorlevel%
exit /b 1
