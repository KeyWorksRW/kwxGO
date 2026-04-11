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
set( CMAKE_CXX_STANDARD 17 )
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
    GIT_REPOSITORY https://github.com/KeyWorksRW/kwxFFI.git
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

add_subdirectory( ${kwxFFI_SOURCE_DIR}/tools/kwxgen kwxgen-build )

set( KWXFFI_HEADERS
    ${kwxFFI_SOURCE_DIR}/include/kwx_classes.h
    ${kwxFFI_SOURCE_DIR}/include/kwx_events.h
    ${kwxFFI_SOURCE_DIR}/include/kwx_keys.h
    ${kwxFFI_SOURCE_DIR}/include/kwx_constants.h
)

add_custom_command(
    OUTPUT  ${CMAKE_SOURCE_DIR}/wx/events_gen.go
            ${CMAKE_SOURCE_DIR}/wx/keys_gen.go
            ${CMAKE_SOURCE_DIR}/wx/constants_gen.go
    COMMAND kwxgen generate
    DEPENDS kwxgen ${KWXFFI_HEADERS}
    COMMENT "Generating Go bindings via kwxgen"
)

add_custom_target( generate_go_bindings ALL
    DEPENDS
        ${CMAKE_SOURCE_DIR}/wx/events_gen.go
        ${CMAKE_SOURCE_DIR}/wx/keys_gen.go
        ${CMAKE_SOURCE_DIR}/wx/constants_gen.go
)

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

go_example(minimal)
go_example(demo)

# ###################### Build Summary #######################

message( STATUS "" )
message( STATUS "kwxGO CMake Configuration:" )
message( STATUS "  Library output:        ${PROJECT_SOURCE_DIR}/lib/$<CONFIG>" )
message( STATUS "    - kwxGO.a (static: kwxApp + transitive kwxFFI/wxWidgets)" )
message( STATUS "  kwxFFI sources:        ${kwxFFI_SOURCE_DIR}" )
message( STATUS "  wxWidgets sources:     ${WX_SRC_DIR}" )
message( STATUS "  wxWidgets build:       ${WX_BUILD_DIR}" )
