@echo off
setlocal

rem ---------------------------------------------------------------------------
rem Build the calendar sample into bin\Debug\calendar.exe, then run it.
rem
rem   run_calender.cmd          build Debug + run
rem   run_calender.cmd build    build Debug only, do not run
rem
rem Thin wrapper over build_sample.cmd so there is exactly one build path. That
rem script owns the CGo environment (CGO_ENABLED, CC/CXX, CGO_CFLAGS, the full
rem CGO_LDFLAGS link line), the .rc -> _windows.syso step, and the Debug default:
rem it links lib\Debug and the wxWidgets *_d.a archives. Never build a sample
rem Release here - the sample gate and the GUI verification both use Debug.
rem ---------------------------------------------------------------------------

cd /D "%~dp0"

if /I "%~1"=="build" (
    call "%~dp0build_sample.cmd" calendar build Debug
    exit /b %ERRORLEVEL%
)

call "%~dp0build_sample.cmd" calendar run Debug
exit /b %ERRORLEVEL%
