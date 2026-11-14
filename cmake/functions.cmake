# ###################### Go Program Helper ######################
#
# Registers a Go program as a CMake build target.
# Usage:  go_program(<source-dir> <name> [<output-path>])
#
#   <source-dir>   directory holding the package, relative to the project root
#   <output-path>  optional output path relative to the project root, without
#                  extension. Default: <source-dir>/<name>
#
#   On Windows, compiles <source-dir>/<name>.rc → <name>_windows.syso when the
#   RC exists, which embeds the Win32 manifest (Common Controls v6, DPI
#   awareness). Without that manifest wxApp::Initialize() shows a warning dialog
#   and all controls use the pre-XP comctl32 v5 look. The _windows suffix is a
#   GOOS constraint, so go build ignores the file on macOS and Linux, and the
#   whole resource step is skipped there.
#   Depends on the kwxGO static library.
#
#   CGO_CFLAGS points to the kwxFFI include directory so CGo can find kwx_classes.h
#   and other generated headers.
#   CGO_LDFLAGS links kwxGO and kwxFFI plus every wxWidgets static archive pulled in
#   through kwxFFI's PUBLIC link interface (libkwxFFI.a itself contains only the FFI
#   wrapper objects), and the required platform system libraries.
#
# ###################### wxWidgets Archive Genexes ######################
#
# Sets <outvar> in the caller's scope to one $<TARGET_LINKER_FILE:...> genex per
# wxWidgets static archive reachable through kwxFFI's PUBLIC link interface. Using
# genexes keeps the paths correct for every platform and configuration instead of
# hardcoding the FetchContent output layout.
#
# Usage:  kwx_wx_archive_genexes(<outvar>)
function(kwx_wx_archive_genexes outvar)
    get_target_property(_kwxffi_deps kwxFFI INTERFACE_LINK_LIBRARIES)
    set(_wx_lib_genex "")
    foreach(_dep ${_kwxffi_deps})
        if(TARGET ${_dep})
            list(APPEND _wx_lib_genex "$<TARGET_LINKER_FILE:${_dep}>")
        endif()
    endforeach()
    set(${outvar} ${_wx_lib_genex} PARENT_SCOPE)
endfunction()


function(go_program src_rel name)
    set(_dir ${PROJECT_SOURCE_DIR}/${src_rel})

    if(ARGC GREATER 2)
        set(_exe ${PROJECT_SOURCE_DIR}/${ARGV2})
    else()
        set(_exe ${_dir}/${name})
    endif()
    if(WIN32)
        set(_exe ${_exe}.exe)
    endif()
    get_filename_component(_exe_dir ${_exe} DIRECTORY)

    # Windows resources: .rc → <name>_windows.syso
    # Go auto-links .syso files in the package directory. Only one of them may
    # contribute a .rsrc section, so each program gets exactly one resource file.
    set(_res_dep "")
    if(WIN32 AND EXISTS ${_dir}/${name}.rc)
        find_program(_windres NAMES windres llvm-windres)
        if(_windres)
            set(_syso ${_dir}/${name}_windows.syso)
            set(_res_inputs ${_dir}/${name}.rc)
            if(EXISTS ${_dir}/${name}.manifest)
                list(APPEND _res_inputs ${_dir}/${name}.manifest)
            endif()
            add_custom_command(
                OUTPUT  ${_syso}
                COMMAND ${_windres} -i ${name}.rc -o ${name}_windows.syso -O coff
                WORKING_DIRECTORY ${_dir}
                DEPENDS ${_res_inputs}
                COMMENT "windres: ${name}.rc → ${name}_windows.syso"
                VERBATIM
            )
            set(_res_dep ${_syso})
        else()
            message( WARNING "windres not found - ${name}.rc is not compiled. "
                "${name} will show the wxWidgets 'no Common Controls v6 "
                "manifest' warning dialog at startup." )
        endif()
    endif()
    # wxWidgets static archives reachable through kwxFFI's PUBLIC link interface
    kwx_wx_archive_genexes(_wx_lib_genex)

    # C++ runtime libraries — compiler-dependent
    if(CMAKE_CXX_COMPILER_ID STREQUAL "GNU")
        set(_cxx_runtime "-lstdc++" "-lgcc_s" "-lgcc")
    elseif(WIN32 AND CMAKE_CXX_COMPILER_ID MATCHES "Clang")
        # clang64's -lstdc++ actually resolves to libc++, so a plain -lstdc++
        # leaves the .exe importing libc++.dll from C:\msys64\clang64\bin and it
        # then only runs with msys64 on PATH. Static libc++.a/libunwind.a keep
        # the executable self-contained (only Windows system DLLs imported).
        set(_cxx_runtime "-static" "-lc++" "-lunwind")
    else()
        set(_cxx_runtime "-lstdc++")
    endif()

    if(WIN32)
        # libkwxFFI.a holds only the FFI wrapper objects; the wxWidgets archives
        # (~70 MB) are linked separately via kwxFFI's PUBLIC link interface.
        # Only system libraries that wxWidgets calls into need to be listed.
        # --start-group/--end-group resolves circular deps between static archives.
        set(_cgo_ldflags
            "-L${PROJECT_SOURCE_DIR}/lib/Release"
            "-Wl,--start-group"
            "-lkwxGO" "-lkwxFFI"
            "$<JOIN:${_wx_lib_genex}, >"
            "-Wl,--end-group"
            "-lgdiplus" "-lmsimg32" "-lcomctl32" "-luuid" "-loleacc" "-lgdi32" "-luser32"
            "-lole32" "-loleaut32" "-luxtheme" "-lwinspool" "-lwinmm"
            "-lshell32" "-lcomdlg32" "-ladvapi32" "-lwsock32" "-lws2_32"
            "-lshlwapi" "-lversion" "-lopengl32"
            ${_cxx_runtime}
        )
    elseif(APPLE)
        set(_cgo_ldflags
            "-L${PROJECT_SOURCE_DIR}/lib"
            "-lkwxGO" "-lkwxFFI"
            "$<JOIN:${_wx_lib_genex}, >"
            "-framework OpenGL" "-framework Cocoa" "-framework IOKit" "-framework Carbon"
            ${_cxx_runtime}
        )
    else()  # Linux
        set(_cgo_ldflags
            "-L${PROJECT_SOURCE_DIR}/lib"
            "-lkwxGO" "-lkwxFFI"
            "$<JOIN:${_wx_lib_genex}, >"
            "-lGL" "-lGLU"
            ${_cxx_runtime}
        )
    endif()

    # Temporary: show all ld.lld errors instead of truncating
    list(APPEND _cgo_ldflags "-Wl,--error-limit=0")

    string(JOIN " " _cgo_ldflags_str ${_cgo_ldflags})

    add_custom_target(${name}
        COMMAND ${CMAKE_COMMAND} -E make_directory ${_exe_dir}
        COMMAND ${CMAKE_COMMAND} -E env
            CGO_ENABLED=1
            CC=${CMAKE_C_COMPILER}
            CXX=${CMAKE_CXX_COMPILER}
            "CGO_CFLAGS=-I${kwxFFI_SOURCE_DIR}/include"
            "CGO_LDFLAGS=${_cgo_ldflags_str}"
            go build "-gcflags=all=-N -l" -o ${_exe} .
        WORKING_DIRECTORY ${_dir}
        DEPENDS kwxGO generate_go_bindings ${_res_dep}
        COMMENT "go build: ${name}"
        VERBATIM
    )
endfunction()


# ###################### Go Example Helper ######################
#
# Registers an example under examples/<name>/ — shorthand for go_program().
# Usage:  go_example(<name>)
# Output: examples/<name>/<name>[.exe]
#
function(go_example name)
    go_program(examples/${name} ${name})
endfunction()
