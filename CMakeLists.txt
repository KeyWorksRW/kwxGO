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
    LANGUAGES CXX
    VERSION ${PROJECT_VERSION_MAJOR}.${PROJECT_VERSION_MINOR}.${PROJECT_VERSION_PATCH}.${PROJECT_VERSION_BUILD}
    DESCRIPTION "kwxGO - Go bindings for wxWidgets via kwxFFI"
    HOMEPAGE_URL "https://github.com/KeyWorksRW/kwxGO" )

# ###################### General Settings #######################
set( CMAKE_CXX_STANDARD 23 )
set( CMAKE_CXX_STANDARD_REQUIRED True )
set( CMAKE_CXX_EXTENSIONS OFF )

# ###################### Options #######################

# Set default build type to Release (kwxFFI.dll is always Release)
if(NOT CMAKE_BUILD_TYPE)
    set(CMAKE_BUILD_TYPE Release CACHE STRING "Choose the type of build." FORCE)
endif()

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

include( FetchContent )

# Suppress verbose FetchContent git output (CMake 3.30+)
set(FETCHCONTENT_QUIET ON)

# kwxFFI builds as a STATIC library — we'll absorb it into kwxGO.dll via /WHOLEARCHIVE
# so there's exactly one DLL with one copy of wxWidgets globals.

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

# kwxFFI is static, but compile with SHARED_LIBRARY so dllexport decorations are
# baked into the .obj files. When /WHOLEARCHIVE pulls them into kwxGO.dll, the
# linker will export them automatically.
target_compile_definitions( kwxFFI PRIVATE SHARED_LIBRARY )

message( NOTICE "kwxFFI has been fetched and configured." )
message( STATUS "kwxFFI sources: " ${kwxFFI_SOURCE_DIR} )

# Define wxWidgets paths (kwxFFI fetches wxWidgets as a subdependency)
# These match the FetchContent directory structure
set( WX_SRC_DIR ${CMAKE_BINARY_DIR}/_deps/wxwidgets-src )
set( WX_BUILD_DIR ${CMAKE_BINARY_DIR}/_deps/wxwidgets-build )
message( STATUS "wxWidgets sources: " ${WX_SRC_DIR} )
message( STATUS "wxWidgets build: " ${WX_BUILD_DIR} )

# ###################### kwxGO Library Target #######################

# kwxGO.dll = kwxApp.cpp + ALL of kwxFFI (via /WHOLEARCHIVE) + wxWidgets.
# This produces a single DLL with one copy of wxWidgets globals so the
# event loop and window creation share the same wxTheApp instance.
add_library( kwxGO SHARED
    ${kwxFFI_SOURCE_DIR}/examples/CApp/kwxApp.cpp
)

# Export macros for kwxApp.cpp; SHARED_LIBRARY enables dllexport in wxffi_def.h
target_compile_definitions( kwxGO PRIVATE KWXAPP_EXPORTS SHARED_LIBRARY )

# Link wxWidgets and kwxFFI. /WHOLEARCHIVE forces every kwxFFI symbol into the
# DLL even though kwxApp.cpp doesn't reference them directly.
target_link_libraries( kwxGO PRIVATE kwxFFI wxbase wxcore )
if( MSVC )
    target_link_options( kwxGO PRIVATE /WHOLEARCHIVE:kwxFFI )
endif()

# kwxApp.h is in the same directory as kwxApp.cpp (uses #include "kwxApp.h")
# wrapper.h is in the kwxFFI include directory (needed by go_events.cpp)
target_include_directories( kwxGO PRIVATE
    ${kwxFFI_SOURCE_DIR}/examples/CApp
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
endif()

# Output kwxGO.dll to bin/$<CONFIG>
set_target_properties( kwxGO PROPERTIES
    ARCHIVE_OUTPUT_DIRECTORY ${PROJECT_SOURCE_DIR}/bin/$<CONFIG>
    LIBRARY_OUTPUT_DIRECTORY ${PROJECT_SOURCE_DIR}/bin/$<CONFIG>
    RUNTIME_OUTPUT_DIRECTORY ${PROJECT_SOURCE_DIR}/bin/$<CONFIG>
)

# ###################### Build Summary #######################

message( STATUS "" )
message( STATUS "kwxGO CMake Configuration:" )
message( STATUS "  Library output:        ${PROJECT_SOURCE_DIR}/bin/$<CONFIG>" )
message( STATUS "    - kwxGO.dll (contains kwxFFI + kwxApp)" )
message( STATUS "  kwxFFI sources:        ${kwxFFI_SOURCE_DIR}" )
message( STATUS "  wxWidgets sources:     ${WX_SRC_DIR}" )
message( STATUS "  wxWidgets build:       ${WX_BUILD_DIR}" )
