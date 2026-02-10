For wxFFI questions, first read `.github/wxffi-architecture.md`

The following are the languages that will be binding to wxFFI:

- ../kwxFortran-dev
- ../kwxGO-dev
- ../kwxJulia-dev
- ../kwxLuaJit-dev
- ../kwxPerl-dev
- ../kwxRust-dev

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
