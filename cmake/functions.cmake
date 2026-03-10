# ###################### Go Example Helper ######################
#
# Registers a Go example as a CMake build target.
# Usage:  go_example(<name>)
#
#   Expects sources under examples/<name>/
#   On Windows, compiles examples/<name>/<name>.rc → <name>.syso if the RC exists.
#   Depends on the kwxGO static library.
#   Output: examples/<name>/<name>[.exe]
#
#   CGO_CFLAGS points to the kwxFFI include directory so CGo can find kwx_classes.h
#   and other generated headers.
#   CGO_LDFLAGS links kwxGO and kwxFFI (which already bundles all wxWidgets symbols)
#   plus the required platform system libraries.
#
function(go_example name)
    set(_dir ${PROJECT_SOURCE_DIR}/examples/${name})

    if(WIN32)
        set(_exe ${_dir}/${name}.exe)
    else()
        set(_exe ${_dir}/${name})
    endif()

    # Windows resources: .rc → .syso  (Go auto-links .syso files in the package dir)
    if(WIN32 AND EXISTS ${_dir}/${name}.rc)
        set(_syso ${_dir}/${name}.syso)
        add_custom_command(
            OUTPUT  ${_syso}
            COMMAND windres -i ${name}.rc -o ${name}.syso -O coff
            WORKING_DIRECTORY ${_dir}
            DEPENDS ${_dir}/${name}.rc
            COMMENT "windres: ${name}.rc → ${name}.syso"
            VERBATIM
        )
        set(_res_dep ${_syso})
    else()
        set(_res_dep "")
    endif()
    # Extract wxWidgets library paths from kwxFFI's PUBLIC dependencies
    get_target_property(_kwxffi_deps kwxFFI INTERFACE_LINK_LIBRARIES)
    set(_wx_lib_genex "")
    foreach(_dep ${_kwxffi_deps})
        if(TARGET ${_dep})
            list(APPEND _wx_lib_genex "$<TARGET_LINKER_FILE:${_dep}>")
        endif()
    endforeach()

    # C++ runtime libraries — compiler-dependent
    if(CMAKE_CXX_COMPILER_ID STREQUAL "GNU")
        set(_cxx_runtime "-lstdc++" "-lgcc_s" "-lgcc")
    elseif(CMAKE_CXX_COMPILER_ID MATCHES "Clang")
        set(_cxx_runtime "-lstdc++")
    else()
        set(_cxx_runtime "-lstdc++")
    endif()

    if(WIN32)
        # libkwxFFI.a already bundles all wxWidgets object code.
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
            "-lshlwapi" "-lversion"
            ${_cxx_runtime}
        )
    elseif(APPLE)
        set(_cgo_ldflags
            "-L${PROJECT_SOURCE_DIR}/lib"
            "-lkwxGO" "-lkwxFFI"
            "$<JOIN:${_wx_lib_genex}, >"
            "-framework Cocoa" "-framework IOKit" "-framework Carbon"
            ${_cxx_runtime}
        )
    else()  # Linux
        set(_cgo_ldflags
            "-L${PROJECT_SOURCE_DIR}/lib"
            "-lkwxGO" "-lkwxFFI"
            "$<JOIN:${_wx_lib_genex}, >"
            ${_cxx_runtime}
        )
    endif()

    # Temporary: show all ld.lld errors instead of truncating
    list(APPEND _cgo_ldflags "-Wl,--error-limit=0")

    string(JOIN " " _cgo_ldflags_str ${_cgo_ldflags})

    add_custom_target(${name}
        COMMAND ${CMAKE_COMMAND} -E env
            CGO_ENABLED=1
            CC=${CMAKE_C_COMPILER}
            CXX=${CMAKE_CXX_COMPILER}
            "CGO_CFLAGS=-I${kwxFFI_SOURCE_DIR}/include"
            "CGO_LDFLAGS=${_cgo_ldflags_str}"
            go build "-gcflags=all=-N -l" -o ${_exe} .
        WORKING_DIRECTORY ${_dir}
        DEPENDS kwxGO ${_res_dep}
        COMMENT "go build: ${name}"
        VERBATIM
    )
endfunction()
