# CMake 3.30 required for FetchContent performance improvements (CMP0168)
cmake_minimum_required(VERSION 3.30)

string( TIMESTAMP YEAR "%y" ) # e.g., "25"
string( TIMESTAMP DOY "%j" ) # e.g., "192"
string( TIMESTAMP PROJECT_YEAR "%Y" )
math( EXPR BUILD_NUM "${YEAR}%100 * 1000 + ${DOY}" )

set( PROJECT_VERSION_MAJOR 0 )
set( PROJECT_VERSION_MINOR 1 )
set( PROJECT_VERSION_PATCH 0 )
set( PROJECT_VERSION_BUILD ${BUILD_NUM} )
set( PROJECT_VERSION_YEAR ${PROJECT_YEAR} )

project( kwxGO
    LANGUAGES C CXX
    VERSION ${PROJECT_VERSION_MAJOR}.${PROJECT_VERSION_MINOR}.${PROJECT_VERSION_PATCH}.${PROJECT_VERSION_BUILD}
    DESCRIPTION "kwxGO - Go bindings for wxWidgets via kwxFFI"
    HOMEPAGE_URL "https://github.com/KeyWorksRW/kwxGO" )

# ###################### sccache #######################
find_program(SCCACHE_PROGRAM sccache)
if(SCCACHE_PROGRAM)
    message(STATUS "sccache found: ${SCCACHE_PROGRAM}")
    set(CMAKE_C_COMPILER_LAUNCHER "${SCCACHE_PROGRAM}")
    set(CMAKE_CXX_COMPILER_LAUNCHER "${SCCACHE_PROGRAM}")
else()
    message(STATUS "sccache not found; building without a compiler cache")
endif()

# ###################### Compiler Guard #######################
# CGo requires GCC or Clang with MinGW-w64 ABI on Windows.
# MSVC (cl.exe) produces incompatible object format and C++ ABI.
if( MSVC )
    message( FATAL_ERROR
        "MSVC (cl.exe) is not supported. kwxGO requires GCC or Clang targeting the "
        "MinGW-w64 ABI for CGo compatibility. "
        "Use the 'ninja-gcc' or 'ninja-clang' CMake preset." )
endif()

# ###################### General Settings #######################
# kwxFFI's public headers (kwx_wrapper.h) require C++23 (std::bit_cast and
# friends), so kwxGO's own C++ sources must use the same standard.
set( CMAKE_CXX_STANDARD 23 )
set( CMAKE_CXX_STANDARD_REQUIRED True )
set( CMAKE_CXX_EXTENSIONS OFF )

# ###################### Options #######################

# kwxFFI and wxWidgets are always Release — no debug CRT needed.
# Go debugging (Delve + -gcflags) is independent of C/C++ debug info.
set(CMAKE_BUILD_TYPE Release CACHE STRING "Always Release for C/C++ code." FORCE)

# Only generate compile_commands.json for local development
if ( NOT DEFINED ENV{CI} )
    # creates a compile_commands.json file in the build directory, suitable for clangd
    set( CMAKE_EXPORT_COMPILE_COMMANDS ON )
endif()

# ###################### ccache support #######################
find_program( CCACHE_PROGRAM ccache )

if( CCACHE_PROGRAM )
    message( STATUS "Using ccache: ${CCACHE_PROGRAM}" )
    set( CMAKE_C_COMPILER_LAUNCHER "${CCACHE_PROGRAM}" )
    set( CMAKE_CXX_COMPILER_LAUNCHER "${CCACHE_PROGRAM}" )

    message( NOTICE "ccache detected. Configure ccache globally for your development environment:" )
    message( NOTICE "  Recommended for all platforms:" )
    message( NOTICE "    ccache --set-config=compression=true" )
    message( NOTICE "    ccache --set-config=max_size=<appropriate for your projects>" )
    message( NOTICE "    ccache --set-config=sloppiness=pch_defines,time_macros,include_file_mtime,include_file_ctime" )

    if( MSVC )
        message( NOTICE "  Required for MSVC:" )

        message( NOTICE "    ccache --set-config=pch_external_checksum=true" )
        message( NOTICE "    ccache --set-config=depend_mode=true" )
        message( NOTICE "    ccache --set-config=\"ignore_options=/Zi /Z7 /ZI /FS /Fd*\"" )
    endif()

    message( NOTICE "  Run 'ccache --show-config' to verify current settings" )
    message( NOTICE "  Run 'ccache -s' to monitor cache usage" )
endif()

# ###################### kwxFFI/wxWidgets via FetchContent #######################

include( FetchContent )

# Suppress verbose FetchContent git output (CMake 3.30+)
set(FETCHCONTENT_QUIET ON)

# kwxFFI is built as a STATIC library and linked into kwxGO.a.

message( STATUS "Fetching kwxFFI and wxWidgets (this may take awhile) ..." )

FetchContent_Declare(
    kwxFFI
    # TEMPORARY: use the local sibling checkout so source changes show up immediately.
    # Restore the GIT_REPOSITORY line below to fetch from GitHub again.
    # GIT_REPOSITORY https://github.com/KeyWorksRW/kwxFFI.git
    SOURCE_DIR ${CMAKE_CURRENT_SOURCE_DIR}/../kwxFFI
    GIT_TAG origin/main
    GIT_SHALLOW TRUE
    UPDATE_DISCONNECTED FALSE  # Explicitly enable updates
)

message( NOTICE "Configuring kwxFFI..." )

FetchContent_MakeAvailable( kwxFFI )

# Build kwxFFI without DLL-export decorations (static library, no dllexport needed).

message( NOTICE "kwxFFI has been fetched and configured." )
message( STATUS "kwxFFI sources: " ${kwxFFI_SOURCE_DIR} )

# ###################### kwxgen — Go binding generator #######################

# kwxFFI's own CMakeLists.txt already does add_subdirectory(kwxgen), so the
# `kwxgen` target exists here; do NOT add it again. The generator source moved
# from kwxFFI/tools/kwxgen to kwxFFI/kwxgen.

# Every kwxFFI include participates in the parsed surface (kwx_classes.h pulls
# in the per-area headers), so glob them to keep the dependency list honest.
file( GLOB KWXFFI_HEADERS CONFIGURE_DEPENDS ${kwxFFI_SOURCE_DIR}/include/*.h )
set( KWXFFI_DEFS ${kwxFFI_SOURCE_DIR}/src/kwx_defs.cpp )

# kwxgen loads kwxgen.config.json from the current working directory, so run it
# from the source root. Explicit arguments below override the config file.
set( KWXGEN_STAMP ${CMAKE_BINARY_DIR}/kwxgen.stamp )

add_custom_command(
    OUTPUT  ${KWXGEN_STAMP}
    COMMAND kwxgen generate
                --headers ${kwxFFI_SOURCE_DIR}/include
                --defs    ${KWXFFI_DEFS}
                --lang    go
                --out     ${CMAKE_SOURCE_DIR}/wx
    COMMAND ${CMAKE_COMMAND} -E touch ${KWXGEN_STAMP}
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    DEPENDS kwxgen ${KWXFFI_HEADERS} ${KWXFFI_DEFS}
            ${CMAKE_SOURCE_DIR}/kwxgen.config.json
    COMMENT "Generating Go bindings via kwxgen"
    VERBATIM
)

add_custom_target( generate_go_bindings ALL DEPENDS ${KWXGEN_STAMP} )

# Define wxWidgets paths (kwxFFI fetches wxWidgets as a subdependency)
# These match the FetchContent directory structure
set( WX_SRC_DIR ${CMAKE_BINARY_DIR}/_deps/wxwidgets-src )
set( WX_BUILD_DIR ${CMAKE_BINARY_DIR}/_deps/wxwidgets-build )
message( STATUS "wxWidgets sources: " ${WX_SRC_DIR} )
message( STATUS "wxWidgets build: " ${WX_BUILD_DIR} )

# ###################### kwxGO Library Target #######################

# kwxGO.a = kwxApp.cpp + go_events.cpp compiled against kwxFFI and wxWidgets headers.
# CGo links against kwxGO.a, kwxFFI.a, and wxWidgets .a libs at Go build time.
# kwxApp.cpp is kwxGO's own copy (originally from kwxFFI/examples/CApp), so we
# can modify it freely for Go-specific needs.
add_library( kwxGO STATIC
    src/kwxApp.cpp
    src/go_events.cpp
)

# No dllexport decorations needed for a static library.
target_compile_definitions( kwxGO PRIVATE KWXAPP_EXPORTS )

# Declare link dependencies so CMake propagates them to dependents.
target_link_libraries( kwxGO PRIVATE kwxFFI )

# kwxApp.h is in src/ alongside kwxApp.cpp (uses #include "kwxApp.h")
# kwx_wrapper.h is in the kwxFFI include directory (needed by go_events.cpp)
target_include_directories( kwxGO PRIVATE
    ${PROJECT_SOURCE_DIR}/src
    ${kwxFFI_SOURCE_DIR}/include
)

# Include wxWidgets headers for kwxApp.cpp
target_include_directories( kwxGO PRIVATE
    ${WX_SRC_DIR}/include
)
if( WIN32 )
    target_include_directories( kwxGO PRIVATE
        ${WX_BUILD_DIR}/lib/vc_x64_lib/mswu
    )
elseif( APPLE )
    target_include_directories( kwxGO PRIVATE
        ${WX_BUILD_DIR}/lib/osx_cocoa-unicode-static
    )
elseif( UNIX )
    target_include_directories( kwxGO PRIVATE
        ${WX_BUILD_DIR}/lib/gtk3_unixu-3.3/include
    )
else()
    message( FATAL_ERROR "Unsupported platform: Only WIN32, APPLE, or UNIX are supported." )
endif()
# Output kwxGO.a to lib/$<CONFIG>
set_target_properties( kwxGO PROPERTIES
    ARCHIVE_OUTPUT_DIRECTORY ${PROJECT_SOURCE_DIR}/lib/$<CONFIG>
)

# ###################### Examples ######################

include( cmake/functions.cmake )

# ###################### Packaging: Merged Archive #######################
#
# `libkwxall.a` is for release packaging only — it merges kwxGO, kwxFFI and every
# wxWidgets static archive into one file, so a released binding needs a single
# library to link against and the linker can resolve wx symbols without
# -Wl,--start-group rescanning.
# Member-level granularity is preserved, so the size of the final executable is
# identical to linking the separate archives.
# The merge is toolchain/ABI-locked, so one bundle must be produced per
# compiler + platform.
#
# Deliberately NOT part of ALL and not a dependency of any other target: ordinary
# development builds never pay for it. Build it explicitly with `ninja kwxall`.
if(TARGET kwxGO AND TARGET kwxFFI)
    kwx_wx_archive_genexes(_kwxall_wx_archives)

    set(_kwxall_inputs
        $<TARGET_LINKER_FILE:kwxGO>
        $<TARGET_LINKER_FILE:kwxFFI>
        ${_kwxall_wx_archives}
    )

    add_custom_target(kwxall
        COMMAND ${CMAKE_COMMAND}
            "-DOUT_ARCHIVE=${PROJECT_SOURCE_DIR}/lib/$<CONFIG>/libkwxall.a"
            "-DIN_ARCHIVES=${_kwxall_inputs}"
            "-DAR_TOOL=${CMAKE_AR}"
            -P ${PROJECT_SOURCE_DIR}/cmake/MergeArchives.cmake
        DEPENDS kwxGO kwxFFI
        WORKING_DIRECTORY ${CMAKE_BINARY_DIR}
        COMMENT "Merging kwxGO + kwxFFI + wxWidgets archives into libkwxall.a"
        VERBATIM
    )
endif()

go_example(minimal)
go_example(demo)

# Samples live under samples/ rather than examples/, but use the same helper.
# calendar.rc → calendar_windows.syso is part of the target, so `ninja calendar`
# produces a runnable bin/Release/calendar.exe whose embedded manifest enables
# Common Controls v6 (otherwise wxWidgets pops up a warning dialog at startup and
# every control looks like Windows 2000). On macOS/Linux no resource step runs
# and go build ignores the _windows-suffixed .syso entirely.
go_program(samples/calendar calendar bin/Release/calendar)

# ###################### Build Summary #######################

message( STATUS "" )
message( STATUS "kwxGO CMake Configuration:" )
message( STATUS "  Library output:        ${PROJECT_SOURCE_DIR}/lib/$<CONFIG>" )
message( STATUS "    - kwxGO.a (static: kwxApp + transitive kwxFFI/wxWidgets)" )
message( STATUS "    - kwxall.a (packaging only, target 'kwxall', not part of ALL)" )
message( STATUS "  kwxFFI sources:        ${kwxFFI_SOURCE_DIR}" )
message( STATUS "  wxWidgets sources:     ${WX_SRC_DIR}" )
message( STATUS "  wxWidgets build:       ${WX_BUILD_DIR}" )
