package wx

/*
#cgo CFLAGS: -I${SRCDIR}/../build/_deps/kwxffi-src/include -I${SRCDIR}/../build/_deps/kwxffi-src/examples/CApp
#cgo windows LDFLAGS: ${SRCDIR}/../bin/Release/kwxGO.dll ${SRCDIR}/../bin/Release/kwxFFI.dll
#cgo linux LDFLAGS: -L${SRCDIR}/../bin/Release -lkwxGO -lkwxFFI -Wl,-rpath,${SRCDIR}/../bin/Release
#cgo darwin LDFLAGS: -L${SRCDIR}/../bin/Release -lkwxGO -lkwxFFI -Wl,-rpath,${SRCDIR}/../bin/Release

#include <stdlib.h>
#include "wxffi_def.h"
#include "wxffi_types.h"
#include "kwxApp.h"

// =============================================================================
// String Functions
// =============================================================================

// Create a wxString from a UTF-8 C string
extern void* wxString_CreateUTF8(char* buffer);

// Create a wxString from a C string with specified length
extern void* wxString_CreateLen(char* buffer, int len);

// Delete a wxString - must be called to free memory
extern void wxString_Delete(void* s);

// Get the UTF-8 content of a wxString into buffer, returns length
extern int wxString_GetString(void* s, char* buffer);

// Get the length of a wxString
extern size_t wxString_Length(void* s);

// =============================================================================
// Closure/Callback Functions
// =============================================================================

// Closure function type: void (*)(void* fun, void* data, void* event)
typedef void (*ClosureFun)(void*, void*, void*);

// Create a closure for event callbacks
extern void* wxClosure_Create(ClosureFun fun, void* data);

// Get the data associated with a closure
extern void* wxClosure_GetData(void* closure);

// =============================================================================
// Event Handler Functions
// =============================================================================

// Create a new event handler
extern void* wxEvtHandler_Create(void);

// Delete an event handler
extern void wxEvtHandler_Delete(void* self);

// Connect an event handler with a closure
// Implemented in src/go_events.cpp as kwxGO_Connect (see wx/events.go)

// =============================================================================
// Window Base Functions
// =============================================================================

// Show a window
extern int wxWindow_Show(void* self);

// Hide a window
extern int wxWindow_Hide(void* self);

// Close a window
extern int wxWindow_Close(void* self, int force);

// Destroy a window
extern int wxWindow_Destroy(void* self);

// Set window size
extern void wxWindow_SetSize(void* self, int x, int y, int width, int height, int sizeFlags);

// Get window size - returns wxSize* that must be deleted
extern void* wxWindow_GetSize(void* self);

// Set window position
extern void wxWindow_Move(void* self, int x, int y);

// Get window position - returns wxPoint* that must be deleted
extern void* wxWindow_GetPosition(void* self);

// Set window label/title
extern void wxWindow_SetLabel(void* self, void* label);

// Enable a window
extern int wxWindow_Enable(void* self);

// Disable a window
extern int wxWindow_Disable(void* self);

// Check if window is enabled
extern int wxWindow_IsEnabled(void* self);

// Check if window is shown
extern int wxWindow_IsShown(void* self);

// Set focus to window
extern void wxWindow_SetFocus(void* self);

// Refresh window
extern void wxWindow_Refresh(void* self, int eraseBackground);

// Update window
extern void wxWindow_Update(void* self);

// Get window ID
extern int wxWindow_GetId(void* self);

// Set window ID
extern void wxWindow_SetId(void* self, int id);

// Get parent window
extern void* wxWindow_GetParent(void* self);

// =============================================================================
// Frame Functions
// =============================================================================

// Create a new frame window
extern void* wxFrame_Create(void* parent, int id, void* title,
                            int x, int y, int w, int h, int style);

// Set frame status text
extern void wxFrame_SetStatusText(void* self, void* text, int number);

// Create status bar
extern void* wxFrame_CreateStatusBar(void* self, int number, int style);

// Create toolbar
extern void* wxFrame_CreateToolBar(void* self, long style);

// Get menu bar
extern void* wxFrame_GetMenuBar(void* self);

// Set menu bar
extern void wxFrame_SetMenuBar(void* self, void* menuBar);

// =============================================================================
// Panel Functions
// =============================================================================

// Create a new panel
extern void* wxPanel_Create(void* parent, int id,
                            int x, int y, int w, int h, int style);

// =============================================================================
// Button Functions
// =============================================================================

// Create a new button
extern void* wxButton_Create(void* parent, int id, void* label,
                             int x, int y, int w, int h, int style);

// Set button as default
extern void wxButton_SetDefault(void* self);

// =============================================================================
// Static Text Functions
// =============================================================================

// Create a new static text (label)
extern void* wxStaticText_Create(void* parent, int id, void* label,
                                 int x, int y, int w, int h, int style);

// =============================================================================
// Text Control Functions
// =============================================================================

// Create a new text control
extern void* wxTextCtrl_Create(void* parent, int id, void* value,
                               int x, int y, int w, int h, long style);

// Get text control value - returns new wxString that must be deleted
extern void* wxTextCtrl_GetValue(void* self);

// Set text control value
extern void wxTextCtrl_SetValue(void* self, void* value);

// Append text to control
extern void wxTextCtrl_AppendText(void* self, void* text);

// Clear text control
extern void wxTextCtrl_Clear(void* self);

// =============================================================================
// CheckBox Functions
// =============================================================================

// Create a new checkbox
extern void* wxCheckBox_Create(void* parent, int id, void* label,
                               int x, int y, int w, int h, int style);

// Get checkbox value
extern int wxCheckBox_GetValue(void* self);

// Set checkbox value
extern void wxCheckBox_SetValue(void* self, int state);

// Check if checkbox is 3-state
extern int wxCheckBox_Is3State(void* self);

// Get 3-state value (0=unchecked, 1=checked, 2=undetermined)
extern int wxCheckBox_Get3StateValue(void* self);

// Set 3-state value
extern void wxCheckBox_Set3StateValue(void* self, int state);

// =============================================================================
// Sizer Functions
// =============================================================================

// Create a box sizer
extern void* wxBoxSizer_Create(int orient);

// Create a flex grid sizer
extern void* wxFlexGridSizer_Create(int rows, int cols, int vgap, int hgap);

// Add flexible growable row
extern void wxFlexGridSizer_AddGrowableRow(void* self, size_t idx);

// Add flexible growable column
extern void wxFlexGridSizer_AddGrowableCol(void* self, size_t idx);

// Create a grid bag sizer
extern void* wxGridBagSizer_Create(int vgap, int hgap);

// Add window to grid bag sizer at position
extern void wxGridBagSizer_AddWindow(void* self, void* window, int row, int col,
                                     int rowspan, int colspan, int flag, int border, void* userData);

// Add sizer to grid bag sizer at position
extern void wxGridBagSizer_AddSizer(void* self, void* sizer, int row, int col,
                                    int rowspan, int colspan, int flag, int border);

// Add spacer to grid bag sizer at position
extern void wxGridBagSizer_AddSpacer(void* self, int width, int height, int row, int col,
                                     int rowspan, int colspan, int flag, int border);

// Create a static box
extern void* wxStaticBox_Create(void* parent, int id, void* text,
                                int x, int y, int w, int h, int style);

// Create a static box sizer
extern void* wxStaticBoxSizer_Create(void* box, int orient);

// Get the static box from a static box sizer
extern void* wxStaticBoxSizer_GetStaticBox(void* self);

// Add a window to a sizer
extern void wxSizer_AddWindow(void* self, void* window, int proportion, int flag, int border, void* userData);

// Add a sizer to a sizer (for nested layouts)
extern void wxSizer_AddSizer(void* self, void* sizer, int proportion, int flag, int border, void* userData);

// Add a spacer to a sizer
extern void wxSizer_AddSpacer(void* self, int size);

// Add stretchable space
extern void wxSizer_AddStretchSpacer(void* self, int prop);

// Layout the sizer
extern void wxSizer_Layout(void* self);

// Set minimum size for sizer
extern void wxSizer_SetMinSize(void* self, int width, int height);

// Fit window to sizer
extern void wxSizer_Fit(void* self, void* window);

// Set sizer and resize window
extern void wxWindow_SetSizerAndFit(void* self, void* sizer, int deleteOld);

// Set sizer without resizing
extern void wxWindow_SetSizer(void* self, void* sizer);

// =============================================================================
// Message Box
// =============================================================================

// Show a message box
extern int wxcMessageBox(void* message, void* caption, int style, void* parent,
                        int x, int y);

// =============================================================================
// Dialog Functions
// =============================================================================

// Dialog base functions (via inheritance)
extern int wxDialog_ShowModal(void* self);

// File dialog
extern void* wxFileDialog_Create(void* parent, void* message, void* defaultDir,
                                 void* defaultFile, void* wildcard, int x, int y, int style);
extern void* wxFileDialog_GetPath(void* self);

// Directory dialog
extern void* wxDirDialog_Create(void* parent, void* message, void* defaultPath, int x, int y, int style);
extern void* wxDirDialog_GetPath(void* self);

// Text entry dialog
extern void* wxTextEntryDialog_Create(void* parent, void* message, void* caption,
                                      void* defaultValue, long style, int x, int y);
extern void* wxTextEntryDialog_GetValue(void* self);

// Color dialog
extern void* wxColourDialog_Create(void* parent, void* data);
extern void* wxColourDialog_GetColourData(void* self);

// Font dialog
extern void* wxFontDialog_Create(void* parent, void* data);
extern void* wxFontDialog_GetFontData(void* self);

// =============================================================================
// Application Functions (kwxApp)
// =============================================================================

// kwxApp functions are declared in kwxApp.h (included above)
// No extern declarations needed - avoiding duplicate declarations

*/
import "C"

// This file contains the raw cgo declarations for wxFFI functions.
// Higher-level Go wrappers are in other files (types.go, string.go, etc.)

// Note: The actual function implementations are in the wxFFI shared library.
// This file just declares the C function signatures for cgo to use.
