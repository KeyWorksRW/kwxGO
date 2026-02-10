// Package wx constants
//
// This file contains commonly-used wxWidgets constants exported from kwxFFI.
// Constants are called at package init time via the expwxXXX() functions.
//
// Note: These are declared as var (not const) because values come from C at runtime.
//
// Additional constants:
// - kwxFFI exports 2000+ constants from kwxFFI/src/defs.cpp
// - Search for "EXPORT int expwxXXXX()" to find more
// - Add new constants following the same pattern as below
// - For bulk imports, consider code generation from defs.cpp

package wx

/*
#cgo CFLAGS: -I${SRCDIR}/../build/_deps/kwxffi-src/include

#include "wxffi_def.h"

// Constants are exported as functions from kwxFFI
// These declarations allow us to call them from Go
extern int expwxID_ANY(void);
extern int expwxID_OK(void);
extern int expwxID_CANCEL(void);
extern int expwxID_YES(void);
extern int expwxID_NO(void);
extern int expwxID_APPLY(void);
extern int expwxID_CLOSE(void);
extern int expwxID_HELP(void);
extern int expwxID_NEW(void);
extern int expwxID_OPEN(void);
extern int expwxID_SAVE(void);
extern int expwxID_SAVEAS(void);
extern int expwxID_EXIT(void);
extern int expwxID_UNDO(void);
extern int expwxID_REDO(void);
extern int expwxID_PRINT(void);
extern int expwxID_ABOUT(void);
extern int expwxID_CUT(void);
extern int expwxID_COPY(void);
extern int expwxID_PASTE(void);
extern int expwxID_CLEAR(void);
extern int expwxID_FIND(void);
extern int expwxID_SELECTALL(void);

// Window styles
extern int expwxDEFAULT_FRAME_STYLE(void);
extern int expwxCAPTION(void);
extern int expwxRESIZE_BORDER(void);
extern int expwxSYSTEM_MENU(void);
extern int expwxMINIMIZE_BOX(void);
extern int expwxMAXIMIZE_BOX(void);
extern int expwxCLOSE_BOX(void);
extern int expwxSTAY_ON_TOP(void);
extern int expwxFRAME_NO_TASKBAR(void);
extern int expwxFRAME_TOOL_WINDOW(void);

// Sizer flags
extern int expwxEXPAND(void);
extern int expwxALL(void);
extern int expwxLEFT(void);
extern int expwxRIGHT(void);
extern int expwxTOP(void);
extern int expwxBOTTOM(void);

// Alignment
extern int expwxALIGN_LEFT(void);
extern int expwxALIGN_RIGHT(void);
extern int expwxALIGN_TOP(void);
extern int expwxALIGN_BOTTOM(void);
extern int expwxALIGN_CENTER(void);
extern int expwxALIGN_CENTRE(void);
extern int expwxALIGN_CENTER_HORIZONTAL(void);
extern int expwxALIGN_CENTRE_HORIZONTAL(void);
extern int expwxALIGN_CENTER_VERTICAL(void);
extern int expwxALIGN_CENTRE_VERTICAL(void);

// Orientation
extern int expwxHORIZONTAL(void);
extern int expwxVERTICAL(void);
extern int expwxBOTH(void);

// Text control styles
extern int expwxTE_MULTILINE(void);
extern int expwxTE_READONLY(void);
extern int expwxTE_PROCESS_ENTER(void);
extern int expwxTE_PROCESS_TAB(void);
extern int expwxTE_PASSWORD(void);
extern int expwxTE_RICH(void);
extern int expwxTE_RICH2(void);
extern int expwxTE_NO_VSCROLL(void);
extern int expwxTE_LEFT(void);
extern int expwxTE_RIGHT(void);
extern int expwxTE_CENTRE(void);
extern int expwxTE_WORDWRAP(void);

// Button styles
extern int expwxBU_LEFT(void);
extern int expwxBU_RIGHT(void);
extern int expwxBU_TOP(void);
extern int expwxBU_BOTTOM(void);
extern int expwxBU_EXACTFIT(void);

// Border styles
extern int expwxBORDER_DEFAULT(void);
extern int expwxBORDER_SIMPLE(void);
extern int expwxBORDER_SUNKEN(void);
extern int expwxBORDER_RAISED(void);
extern int expwxBORDER_STATIC(void);
extern int expwxBORDER_THEME(void);
extern int expwxBORDER_NONE(void);

// Event types (command events)
extern int expwxEVT_BUTTON(void);
extern int expwxEVT_CHECKBOX(void);
extern int expwxEVT_CHOICE(void);
extern int expwxEVT_LISTBOX(void);
extern int expwxEVT_LISTBOX_DCLICK(void);
extern int expwxEVT_TEXT(void);
extern int expwxEVT_TEXT_ENTER(void);
extern int expwxEVT_MENU(void);
extern int expwxEVT_SLIDER(void);
extern int expwxEVT_RADIOBOX(void);
extern int expwxEVT_RADIOBUTTON(void);

// Event types (window events)
extern int expwxEVT_CLOSE_WINDOW(void);
extern int expwxEVT_SIZE(void);
extern int expwxEVT_MOVE(void);
extern int expwxEVT_PAINT(void);
extern int expwxEVT_SHOW(void);
extern int expwxEVT_ACTIVATE(void);

// MessageBox styles
extern int expwxOK(void);
extern int expwxCANCEL(void);
extern int expwxYES_NO(void);
extern int expwxYES_DEFAULT(void);
extern int expwxNO_DEFAULT(void);
extern int expwxICON_ERROR(void);
extern int expwxICON_QUESTION(void);
extern int expwxICON_INFORMATION(void);
extern int expwxICON_WARNING(void);
*/
import "C"

// Standard IDs for buttons and menu items
var (
	ID_ANY       = int(C.expwxID_ANY())
	ID_OK        = int(C.expwxID_OK())
	ID_CANCEL    = int(C.expwxID_CANCEL())
	ID_YES       = int(C.expwxID_YES())
	ID_NO        = int(C.expwxID_NO())
	ID_APPLY     = int(C.expwxID_APPLY())
	ID_CLOSE     = int(C.expwxID_CLOSE())
	ID_HELP      = int(C.expwxID_HELP())
	ID_NEW       = int(C.expwxID_NEW())
	ID_OPEN      = int(C.expwxID_OPEN())
	ID_SAVE      = int(C.expwxID_SAVE())
	ID_SAVEAS    = int(C.expwxID_SAVEAS())
	ID_EXIT      = int(C.expwxID_EXIT())
	ID_UNDO      = int(C.expwxID_UNDO())
	ID_REDO      = int(C.expwxID_REDO())
	ID_PRINT     = int(C.expwxID_PRINT())
	ID_ABOUT     = int(C.expwxID_ABOUT())
	ID_CUT       = int(C.expwxID_CUT())
	ID_COPY      = int(C.expwxID_COPY())
	ID_PASTE     = int(C.expwxID_PASTE())
	ID_CLEAR     = int(C.expwxID_CLEAR())
	ID_FIND      = int(C.expwxID_FIND())
	ID_SELECTALL = int(C.expwxID_SELECTALL())
)

// Window styles for frames and dialogs
var (
	DEFAULT_FRAME_STYLE = int(C.expwxDEFAULT_FRAME_STYLE())
	CAPTION             = int(C.expwxCAPTION())
	RESIZE_BORDER       = int(C.expwxRESIZE_BORDER())
	SYSTEM_MENU         = int(C.expwxSYSTEM_MENU())
	MINIMIZE_BOX        = int(C.expwxMINIMIZE_BOX())
	MAXIMIZE_BOX        = int(C.expwxMAXIMIZE_BOX())
	CLOSE_BOX           = int(C.expwxCLOSE_BOX())
	STAY_ON_TOP         = int(C.expwxSTAY_ON_TOP())
	FRAME_NO_TASKBAR    = int(C.expwxFRAME_NO_TASKBAR())
	FRAME_TOOL_WINDOW   = int(C.expwxFRAME_TOOL_WINDOW())
)

// Sizer flags
var (
	EXPAND = int(C.expwxEXPAND())
	ALL    = int(C.expwxALL())
	LEFT   = int(C.expwxLEFT())
	RIGHT  = int(C.expwxRIGHT())
	TOP    = int(C.expwxTOP())
	BOTTOM = int(C.expwxBOTTOM())
)

// Alignment flags
var (
	ALIGN_LEFT              = int(C.expwxALIGN_LEFT())
	ALIGN_RIGHT             = int(C.expwxALIGN_RIGHT())
	ALIGN_TOP               = int(C.expwxALIGN_TOP())
	ALIGN_BOTTOM            = int(C.expwxALIGN_BOTTOM())
	ALIGN_CENTER            = int(C.expwxALIGN_CENTER())
	ALIGN_CENTRE            = int(C.expwxALIGN_CENTRE())
	ALIGN_CENTER_HORIZONTAL = int(C.expwxALIGN_CENTER_HORIZONTAL())
	ALIGN_CENTRE_HORIZONTAL = int(C.expwxALIGN_CENTRE_HORIZONTAL())
	ALIGN_CENTER_VERTICAL   = int(C.expwxALIGN_CENTER_VERTICAL())
	ALIGN_CENTRE_VERTICAL   = int(C.expwxALIGN_CENTRE_VERTICAL())
)

// Orientation
var (
	HORIZONTAL = int(C.expwxHORIZONTAL())
	VERTICAL   = int(C.expwxVERTICAL())
	BOTH       = int(C.expwxBOTH())
)

// Text control styles
var (
	TE_MULTILINE     = int(C.expwxTE_MULTILINE())
	TE_READONLY      = int(C.expwxTE_READONLY())
	TE_PROCESS_ENTER = int(C.expwxTE_PROCESS_ENTER())
	TE_PROCESS_TAB   = int(C.expwxTE_PROCESS_TAB())
	TE_PASSWORD      = int(C.expwxTE_PASSWORD())
	TE_RICH          = int(C.expwxTE_RICH())
	TE_RICH2         = int(C.expwxTE_RICH2())
	TE_NO_VSCROLL    = int(C.expwxTE_NO_VSCROLL())
	TE_LEFT          = int(C.expwxTE_LEFT())
	TE_RIGHT         = int(C.expwxTE_RIGHT())
	TE_CENTRE        = int(C.expwxTE_CENTRE())
	TE_WORDWRAP      = int(C.expwxTE_WORDWRAP())
)

// Button styles
var (
	BU_LEFT     = int(C.expwxBU_LEFT())
	BU_RIGHT    = int(C.expwxBU_RIGHT())
	BU_TOP      = int(C.expwxBU_TOP())
	BU_BOTTOM   = int(C.expwxBU_BOTTOM())
	BU_EXACTFIT = int(C.expwxBU_EXACTFIT())
)

// Border styles
var (
	BORDER_DEFAULT = int(C.expwxBORDER_DEFAULT())
	BORDER_SIMPLE  = int(C.expwxBORDER_SIMPLE())
	BORDER_SUNKEN  = int(C.expwxBORDER_SUNKEN())
	BORDER_RAISED  = int(C.expwxBORDER_RAISED())
	BORDER_STATIC  = int(C.expwxBORDER_STATIC())
	BORDER_THEME   = int(C.expwxBORDER_THEME())
	BORDER_NONE    = int(C.expwxBORDER_NONE())
)

// Event types - Command events
var (
	EVT_COMMAND_BUTTON_CLICKED        int = int(C.expwxEVT_BUTTON())
	EVT_COMMAND_CHECKBOX_CLICKED      int = int(C.expwxEVT_CHECKBOX())
	EVT_COMMAND_CHOICE_SELECTED       int = int(C.expwxEVT_CHOICE())
	EVT_COMMAND_LISTBOX_SELECTED      int = int(C.expwxEVT_LISTBOX())
	EVT_COMMAND_LISTBOX_DOUBLECLICKED int = int(C.expwxEVT_LISTBOX_DCLICK())
	EVT_COMMAND_TEXT_UPDATED          int = int(C.expwxEVT_TEXT())
	EVT_COMMAND_TEXT_ENTER            int = int(C.expwxEVT_TEXT_ENTER())
	EVT_COMMAND_MENU_SELECTED         int = int(C.expwxEVT_MENU())
	EVT_COMMAND_SLIDER_UPDATED        int = int(C.expwxEVT_SLIDER())
	EVT_COMMAND_RADIOBOX_SELECTED     int = int(C.expwxEVT_RADIOBOX())
	EVT_COMMAND_RADIOBUTTON_SELECTED  int = int(C.expwxEVT_RADIOBUTTON())
)

// Event types - Window events
var (
	EVT_CLOSE_WINDOW int = int(C.expwxEVT_CLOSE_WINDOW())
	EVT_SIZE         int = int(C.expwxEVT_SIZE())
	EVT_MOVE         int = int(C.expwxEVT_MOVE())
	EVT_PAINT        int = int(C.expwxEVT_PAINT())
	EVT_SHOW         int = int(C.expwxEVT_SHOW())
	EVT_ACTIVATE     int = int(C.expwxEVT_ACTIVATE())
)

// MessageBox styles
var (
	OK               = int(C.expwxOK())
	CANCEL           = int(C.expwxCANCEL())
	YES_NO           = int(C.expwxYES_NO())
	YES_DEFAULT      = int(C.expwxYES_DEFAULT())
	NO_DEFAULT       = int(C.expwxNO_DEFAULT())
	ICON_ERROR       = int(C.expwxICON_ERROR())
	ICON_QUESTION    = int(C.expwxICON_QUESTION())
	ICON_INFORMATION = int(C.expwxICON_INFORMATION())
	ICON_WARNING     = int(C.expwxICON_WARNING())

	// Aliases (wxICON_EXCLAMATION == wxICON_WARNING, wxICON_HAND == wxICON_ERROR)
	ICON_EXCLAMATION = int(C.expwxICON_WARNING())
	ICON_HAND        = int(C.expwxICON_ERROR())
)

// Commonly used event type aliases (shorter names)
var (
	EVT_BUTTON int = EVT_COMMAND_BUTTON_CLICKED
	EVT_TEXT   int = EVT_COMMAND_TEXT_UPDATED
	EVT_MENU   int = EVT_COMMAND_MENU_SELECTED
)
