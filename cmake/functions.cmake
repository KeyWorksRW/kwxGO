# ###################### Go Example Helper ######################
#
# Registers a Go example as a CMake build target.
# Usage:  go_example(<name>)
#
#   Expects sources under examples/<name>/
#   On Windows, compiles examples/<name>/<name>.rc → <name>.syso if the RC exists.
#   Depends on the kwxGO static library (built in the same $<CONFIG>).
#   Output: examples/<name>/<name>[.exe]
#
#   CGO_LDFLAGS prepends $<CONFIG>-specific library search paths so that the
#   correct Debug/Release libkwxGO.a and libkwxFFI.a are found first, regardless
#   of the hardcoded "Release" paths in wx/ffi.go (which remain as a fallback
#   for developers running 'go build' directly).
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

    add_custom_target(${name}
        COMMAND ${CMAKE_COMMAND} -E env
            CGO_ENABLED=1
            CC=${CMAKE_C_COMPILER}
            CXX=${CMAKE_CXX_COMPILER}
            "CGO_LDFLAGS=-L${PROJECT_SOURCE_DIR}/bin/$<CONFIG> -L${PROJECT_SOURCE_DIR}/lib/$<CONFIG>"
            go build "-gcflags=all=-N -l" -o ${_exe} .
        WORKING_DIRECTORY ${_dir}
        DEPENDS kwxGO ${_res_dep}
        COMMENT "go build: ${name}"
        VERBATIM
    )
endfunction()
