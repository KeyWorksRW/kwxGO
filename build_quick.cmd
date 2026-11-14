@echo off
setlocal

rem Build the kwxGO quick test app into bin\<config>\quick.exe. The config defaults
rem to Debug; Release must be requested explicitly as a second argument.
rem
rem   build_quick.cmd                  build Debug only -- the verification gate
rem   build_quick.cmd run              build Debug, then launch the executable
rem   build_quick.cmd build Release    build the Release config instead
rem   build_quick.cmd probe            build the probe lane (tests\quick_probe\gen)
rem   build_quick.cmd help             this text
rem
rem The sources are tests\quick: quick_main.go (hand-written) plus the forms
rem quick.wxui generates there (quick_mainframe.go, quick_dlg.go).
rem
rem The result is a normal standalone executable: every wxWidgets / kwxFFI /
rem kwxGO archive is linked statically, and libc++/libunwind are static too, so
rem the .exe imports only Windows system DLLs. It can be copied anywhere and run
rem without Go, CMake or msys64 on the machine.
rem
rem This mirrors cmake/functions.cmake (go_program): CGo needs the same
rem environment, otherwise the link fails with undefined kwxApp_*/wxApp::*
rem symbols. A quick.rc file in tests\quick is compiled to
rem quick_windows.syso, which `go build` links automatically; without it
rem wxApp::Initialize() shows a "no Common Controls v6 manifest" warning dialog.
rem ---------------------------------------------------------------------------

cd /D "%~dp0"
set "ROOT=%~dp0"

set "TARGET=quick"
set "ACTION=%~1"
set "CONFIG=%~2"
if "%CONFIG%"=="" set "CONFIG=Debug"
if /I not "%CONFIG%"=="Debug" if /I not "%CONFIG%"=="Release" goto :bad_config

rem `probe` builds the probe lane (tests\quick_probe\gen) with the same cgo
rem environment and link line; only the package and the output name differ.
if /I "%ACTION%"=="probe" (
    set "TARGET=probe"
    set "ACTION=build"
)
set "SAMPLE=quick"
if /I "%TARGET%"=="probe" set "SAMPLE=probe"
set "SAMPLE_DIR=%ROOT%tests\quick"
if /I "%TARGET%"=="probe" set "SAMPLE_DIR=%ROOT%tests\quick_probe\gen"
set "PKG=./tests/quick"
if /I "%TARGET%"=="probe" set "PKG=./tests/quick_probe/gen"
rem Debug links the wxWidgets *_d.a archives and lib\Debug; Release the unsigned ones.
if /I "%CONFIG%"=="Debug" (set "WXSUF=d") else (set "WXSUF=")

if /I "%ACTION%"=="help" goto :usage
if /I "%ACTION%"=="-h" goto :usage
if /I "%ACTION%"=="/?" goto :usage

if "%ACTION%"=="" set "ACTION=build"
if /I not "%ACTION%"=="build" if /I not "%ACTION%"=="run" goto :bad_action

if not exist "%SAMPLE_DIR%\" goto :no_sample
dir /b "%SAMPLE_DIR%\*.go" >nul 2>&1
if errorlevel 1 goto :no_go

if /I "%TARGET%"=="probe" (
    rem Go requires one directory per package, so the hand-written probe entry
    rem joins the generated forms in gen\ before the package is built.
    copy /Y "%ROOT%tests\quick_probe\probe_main.go" "%SAMPLE_DIR%\probe_main.go" >nul
    if errorlevel 1 goto :copy_failed
)

where go >nul 2>&1
if errorlevel 1 (
    echo ERROR: go not found on PATH.
    exit /b 1
)

rem Use forward slashes from here on: cgo shell-splits CGO_LDFLAGS and treats a
rem backslash as an escape, which corrupts the flags on the way to the linker.
set "KWXGO_ROOT=%ROOT:\=/%"

set "CGO_ENABLED=1"
set "CC=C:/msys64/clang64/bin/clang.exe"
set "CXX=C:/msys64/clang64/bin/clang++.exe"

rem kwxFFI lives beside this repository (../kwxFFI); fall back to the checkout
rem path used by the sibling workspace when the relative one is not present.
for %%I in ("%KWXGO_ROOT%../kwxFFI/include") do set "KWXFFI_INC=%%~fI"
if not exist "%KWXFFI_INC%\kwx_classes.h" set "KWXFFI_INC=C:\rwCode\wxLanguages\kwxFFI\include"
if not exist "%KWXFFI_INC%\kwx_classes.h" (
    echo ERROR: kwxFFI include directory not found ^(kwx_classes.h missing^).
    echo        Expected "%KWXFFI_INC%".
    exit /b 1
)
set "KWXFFI_INC=%KWXFFI_INC:\=/%"
set "CGO_CFLAGS=-I%KWXFFI_INC%"

rem Locate the wxWidgets static archives matching the requested config.
set "WXLIB="
for /D %%D in ("%ROOT%build\_deps\wxwidgets-build\lib\*") do (
    if exist "%%~D\wxbase33u%WXSUF%.a" set "WXLIB=%%~D"
)
if not defined WXLIB (
    echo ERROR: wxWidgets %CONFIG% static libs not found under build\_deps\wxwidgets-build\lib\.
    echo Configure and build the project first.
    exit /b 1
)
set "WXLIB=%WXLIB:\=/%"

rem A partially built config links badly, so name the exact missing archive.
set "WXA_MISSING="
for %%A in (wxbase33u%WXSUF%.a wxbase33u%WXSUF%_xml.a wxmsw33u%WXSUF%_core.a wxmsw33u%WXSUF%_xrc.a wxlunasvg%WXSUF%.a wxpng%WXSUF%.a wxjpeg%WXSUF%.a wxtiff%WXSUF%.a wxregexu%WXSUF%.a wxexpat%WXSUF%.a wxscintilla%WXSUF%.a wxwebp%WXSUF%.a wxwebpdemux%WXSUF%.a wxsharpyuv%WXSUF%.a) do (
    if not exist "%WXLIB%/%%A" set "WXA_MISSING=%%A"
)
if defined WXA_MISSING (
    echo ERROR: wxWidgets %CONFIG% archive %WXA_MISSING% is missing in "%WXLIB%".
    echo        Build the complete %CONFIG% config first, from the build folder:
    echo          ninja -f build-%CONFIG%.ninja
    exit /b 1
)

if not exist "%KWXGO_ROOT%lib/%CONFIG%/libkwxGO.a" (
    echo ERROR: lib\%CONFIG%\libkwxGO.a not found - build the %CONFIG% C/C++ side first.
    echo        ninja -f build-%CONFIG%.ninja kwxGO   ^(from the build folder^)
    exit /b 1
)
if not exist "%KWXGO_ROOT%lib/%CONFIG%/libkwxFFI.a" (
    echo ERROR: lib\%CONFIG%\libkwxFFI.a not found - build the %CONFIG% C/C++ side first.
    echo        ninja -f build-%CONFIG%.ninja kwxGO   ^(from the build folder^)
    exit /b 1
)

rem -static makes the link fully self-contained (libc++.a + libunwind.a instead
rem of libc++.dll). Drop -static and restore -lstdc++ for a smaller .exe that
rem needs C:\msys64\clang64\bin on PATH when it runs.
rem GCC preset users: replace "-lc++ -lunwind" with "-lstdc++ -lgcc_s -lgcc".
set "CGO_LDFLAGS=-L%KWXGO_ROOT%lib/%CONFIG% -Wl,--start-group -lkwxGO -lkwxFFI %WXLIB%/wxbase33u%WXSUF%.a %WXLIB%/wxmsw33u%WXSUF%_core.a %WXLIB%/wxlunasvg%WXSUF%.a %WXLIB%/wxmsw33u%WXSUF%_aui.a %WXLIB%/wxmsw33u%WXSUF%_html.a %WXLIB%/wxmsw33u%WXSUF%_ribbon.a %WXLIB%/wxmsw33u%WXSUF%_richtext.a %WXLIB%/wxmsw33u%WXSUF%_propgrid.a %WXLIB%/wxmsw33u%WXSUF%_webview.a %WXLIB%/wxmsw33u%WXSUF%_stc.a %WXLIB%/wxbase33u%WXSUF%_xml.a %WXLIB%/wxmsw33u%WXSUF%_xrc.a %WXLIB%/wxmsw33u%WXSUF%_gl.a %WXLIB%/wxmsw33u%WXSUF%_media.a %WXLIB%/wxbase33u%WXSUF%_net.a %WXLIB%/wxzlib%WXSUF%.a %WXLIB%/wxpng%WXSUF%.a %WXLIB%/wxjpeg%WXSUF%.a %WXLIB%/wxexpat%WXSUF%.a %WXLIB%/wxscintilla%WXSUF%.a %WXLIB%/wxlexilla%WXSUF%.a %WXLIB%/wxwebp%WXSUF%.a %WXLIB%/wxwebpdemux%WXSUF%.a %WXLIB%/wxsharpyuv%WXSUF%.a %WXLIB%/wxregexu%WXSUF%.a %WXLIB%/wxtiff%WXSUF%.a -Wl,--end-group -lgdiplus -lmsimg32 -lcomctl32 -luuid -loleacc -lgdi32 -luser32 -lole32 -loleaut32 -luxtheme -lwinspool -lwinmm -lshell32 -lcomdlg32 -ladvapi32 -lwsock32 -lws2_32 -lshlwapi -lversion -lopengl32 -static -lc++ -lunwind -Wl,--error-limit=0"

rem ---------------------------------------------------------------------------
rem Windows resources: embed the application manifest (Common Controls v6).
rem windres runs in the resource folder so a relative path inside the .rc (for
rem example "1 RT_MANIFEST ""quick.manifest""") resolves. Only one .syso per
rem package may carry a .rsrc section. The quick target keeps its .rc beside the
rem generated forms; the probe target keeps them in tests\quick_probe because
rem gen\ is generated output. The .syso is always written into the package folder.
rem ---------------------------------------------------------------------------
set "RC_DIR=%SAMPLE_DIR%"
if /I "%TARGET%"=="probe" set "RC_DIR=%ROOT%tests\quick_probe"

set "RC_FILE="
for %%R in ("%RC_DIR%\*.rc") do if not defined RC_FILE set "RC_FILE=%%~nxR"
rem An unmatched wildcard is passed through literally - drop it.
if defined RC_FILE if not exist "%RC_DIR%\%RC_FILE%" set "RC_FILE="

if not defined RC_FILE goto :nores
set "WINDRES="
for %%T in (windres llvm-windres) do if not defined WINDRES if not "%%~$PATH:T"=="" set "WINDRES=%%~$PATH:T"
if not defined WINDRES if exist "C:\msys64\clang64\bin\windres.exe" set "WINDRES=C:\msys64\clang64\bin\windres.exe"
if not defined WINDRES goto :nowindres

pushd "%RC_DIR%"
"%WINDRES%" -i "%RC_FILE%" -o "%SAMPLE_DIR%\%SAMPLE%_windows.syso" -O coff
if errorlevel 1 (
    popd
    echo ERROR: windres failed to compile %RC_FILE%.
    exit /b 1
)
popd
goto :nores

:nowindres
echo WARNING: windres not found - the exe will lack the Common Controls v6
echo          manifest and wxWidgets will show a warning dialog at startup.

:nores

set "OUT=bin/%CONFIG%/%SAMPLE%.exe"
if not exist "bin\%CONFIG%" mkdir "bin\%CONFIG%"

rem Debug builds disable optimisation and inlining so Delve can step and show
rem accurate values. The quotes are part of the value so the two gcflags reach
rem `go build` as one argument; an empty value expands to nothing.
set GOFLAGS_CFG=
if /I "%CONFIG%"=="Debug" set GOFLAGS_CFG="-gcflags=all=-N -l"

echo Building %OUT% ...
go build %GOFLAGS_CFG% -o "%OUT%" "%PKG%"
if errorlevel 1 (
    echo ERROR: go build failed.
    exit /b 1
)
echo Built %OUT%

if /I not "%ACTION%"=="run" exit /b 0

echo Running %OUT% ...
pushd "%SAMPLE_DIR%"
"%ROOT%%OUT%"
set "RC=%ERRORLEVEL%"
popd
exit /b %RC%

:usage
call :usage_text
exit /b 0

:usage_fail
echo ERROR: nothing to build.
echo.
call :usage_text
exit /b 1

:usage_text
echo Usage: build_quick.cmd [build^|run^|probe] [Debug^|Release]
echo.
echo   build_quick.cmd             build tests\quick into bin\Debug\quick.exe
echo   build_quick.cmd run         build Debug, then launch the executable
echo   build_quick.cmd build Release  build the Release config instead
echo   build_quick.cmd probe       build tests\quick_probe\gen into bin\Debug\probe.exe
echo.
echo   build  build only - use this for verification (default)
echo   run    build, then launch the GUI
echo   Debug  default config; links lib\Debug and the wxWidgets *_d.a archives
echo   Release  opt-in config; links lib\Release and the wxWidgets *.a archives
exit /b 0

:bad_action
echo ERROR: unknown action "%ACTION%" ^(expected "build", "run" or "probe"^).
echo.
exit /b 1

:bad_config
echo ERROR: unknown config "%CONFIG%" ^(expected "Debug" or "Release"^).
echo.
exit /b 1

:no_sample
echo ERROR: %SAMPLE_DIR% does not exist.
exit /b 1

:copy_failed
echo ERROR: could not copy probe_main.go into %SAMPLE_DIR%.
exit /b 1

:no_go
echo ERROR: tests\quick\ contains no .go files.
exit /b 1
