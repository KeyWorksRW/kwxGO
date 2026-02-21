For kwxFFI questions, first read `.shared/kwxffi-architecture.md`.

The following are the languages that will be binding to kwxFFI:

- ../kwxFortran-dev
- ../kwxGO-dev
- ../kwxJulia-dev
- ../kwxLuaJit-dev
- ../kwxPerl-dev
- ../kwxRust-dev

## kwxFFI Architecture Quick Reference

Full details in `.shared/kwxffi-architecture.md`. Key facts:

### Headers
| Header | Purpose |
|--------|---------|
| `kwx_glue.h` | All C wrapper function declarations (~7000 lines) — the primary API header |
| `kwx_wrapper.h` | Precompiled header — all wxWidgets includes + helper class declarations |
| `kwx_def.h` | Export macros (`EXPORT`, `WXFFI_FUNC`, `WXFFI_CONSTANT_INT`, etc.) |
| `kwx_types.h` | Type macros (`TClass`, `TSelf`, `TString`, `TBool`, `TPoint(x,y)`, etc.) |
| `kwx_grid.h` | Callback types and `kwxGridTable` for custom grid data sources |

### Function Naming
- Constructor: `wxButton_Create(parent, id, text, x, y, w, h, style)`
- Method: `wxButton_SetDefault(self)`
- Constant: `expwxDEFAULT_FRAME_STYLE()` → returns `int`
- Event type: `expEVT_BUTTON()` → returns event ID `int`
- String ops: `wxString_Create`, `wxString_Delete`, `wxString_GetString`, `kwxUtf8Buffer_*`

### String Handling
```cpp
// Input: pass wxString*, function dereferences with *str
wxString* s = wxString_Create("Hello");   // or wxString_CreateUTF8
wxTextCtrl_SetValue(ctrl, s);
wxString_Delete(s);                        // caller always deletes input strings

// Output: function allocates, caller must delete
wxString* result = wxTextCtrl_GetValue(ctrl);
// extract via wxString_GetString(result, buf) or kwxUtf8Buffer_*
wxString_Delete(result);
```

### Type Macros (kwx_types.h)
`TClass(tp)` = `void*`, `TSelf(tp)` = `void*`, `TString` = `char*`, `TBool` = `int`,
`TPoint(x,y)` = `int x, int y`, `TSize(w,h)` = `int w, int h`, `TRect(x,y,w,h)` = four ints.

### App Lifecycle (from `examples/CApp/kwxApp.h`)
```c
kwxApp_Initialize(argc, argv);  // call once before any wx functions
kwxApp_MainLoop();               // blocks until exit
kwxApp_ExitMainLoop();           // call from event handler to quit
```

### Event System
```cpp
// Connect: closure wraps a foreign function pointer + data
wxEvtHandler_Connect(obj, first, last, type, closure);
// ClosureFun signature: void(void* closureFun, void* data, void* event)
```

### Library Output (MinGW/static build)
- Windows: `lib/<Config>/libkwxFFI.a`
- Linux: `lib/libkwxFFI.a`
- macOS: `lib/libkwxFFI.a`

### wxWidgets Platform Setup Header Location
After building via FetchContent:
- Windows: `build/_deps/wxwidgets-build/lib/gcc_x64_lib/mswu` (MinGW GCC)
- macOS: `build/_deps/wxwidgets-build/lib/osx_cocoa-unicode-static`
- Linux: `build/_deps/wxwidgets-build/lib/gtk3_unixu-3.3/include`

---

## Running Created Programs

**Do not run programs we create from `key_term` or `run_in_terminal`.** If the program hangs (e.g., GUI app waiting for input, dialog box, crash), the terminal will also hang and become unresponsive.

Instead:
1. **Preferred:** Create a `.vscode/launch.json` configuration so the program can be launched from the VS Code debugger (F5). This allows proper debugging and the user can terminate it independently.
2. **Alternative:** Provide command-line instructions for the user to run manually in their own terminal.

## kwxFFI Issues

**Do NOT create workarounds for missing features or bugs in kwxFFI.** The user owns kwxFFI and wants issues fixed at the source.

When you encounter a missing export, missing function, bug, or limitation in kwxFFI:
1. **Stop** - Do not implement a workaround
2. **Explain** the problem clearly (what's missing, what you tried, what the expected behavior should be)
3. **Ask** if the user wants you to create an issue in kwxFFI
4. If yes, use `gh issue create --repo KeyWorksRW/kwxFFI --assignee Randalphwa --title "..." --body "..."` to create the issue
5. **Wait** for the issue to be fixed before continuing with that feature
