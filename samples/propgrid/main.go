// Command propgrid is the kwxGO port of the wxWidgets "propgrid" sample
// (samples/propgrid/propgrid.cpp, propgrid.h, sampleprops.cpp).
//
// The frame hosts a wxPropertyGridManager with three pages ("Standard Items",
// "wxWidgets Library Config" and "Examples"), a status bar, an extensive menu
// bar, two buttons below the grid and a wxLogWindow beside the frame.
package main

import (
	"fmt"
	"unsafe"

	"github.com/KeyWorksRW/kwxGO/wx"
)

// Command IDs, mirroring the anonymous enum in propgrid.cpp. The numeric values
// matter: the frame binds its handlers by ID and the grid manager is created
// with idPGID so that its command events reach the frame.
const (
	idPGID               = 1
	idAbout              = 2
	idQuit               = 3
	idAppendProp         = 4
	idAppendCat          = 5
	idInsertProp         = 6
	idInsertCat          = 7
	idEnable             = 8
	idSetReadOnly        = 9
	idHide               = 10
	idBoolCheckbox       = 11
	idDelete             = 12
	idDeleteR            = 13
	idDeleteAll          = 14
	idUnspecify          = 15
	idIterate1           = 16
	idIterate2           = 17
	idIterate3           = 18
	idIterate4           = 19
	idClearModif         = 20
	idFreeze             = 21
	idDumpList           = 22
	idColourScheme1      = 23
	idColourScheme2      = 24
	idColourScheme3      = 25
	idCatColours         = 26
	idSetBgColour        = 27
	idSetBgColourRecur   = 28
	idStaticLayout       = 29
	idPopulate1          = 30
	idPopulate2          = 31
	idCollapse           = 32
	idCollapseAll        = 33
	idGetValues          = 34
	idSetValues          = 35
	idSetValues2         = 36
	idFitColumns         = 37
	idChangeFlagsItems   = 38
	idTestInsertChoice   = 39
	idTestDeleteChoice   = 40
	idInsertPage         = 41
	idRemovePage         = 42
	idSetSpinCtrlEditor  = 43
	idSetPropertyValue   = 44
	idTestReplace        = 45
	idSetColumns         = 46
	idSetVirtWidth       = 47
	idSetPGDisabled      = 48
	idTestXRC            = 49
	idEnableCommonValues = 50
	idSelectStyle        = 51
	idSaveState          = 52
	idRestoreState       = 53
	idRunMinimal         = 54
	idEnableLabelEditing = 55
	idVetoColDrag        = 56
	idShowHeader         = 57
	idOnExtendedKeyNav   = 58
	idShowPopup          = 59
	idPopupGrid          = 60

	// Declared in propgrid.h.
	idColourScheme4 = 100

	// Values from wxPGPropertyValuesFlags (wx/propgrid/propgriddefs.h). The FFI
	// takes these as plain ints because the enum is not exported as Go
	// constants. wxPGVFBFlags come from the wx package (wx.PG_VFB_*).
	pgDontRecurse   = 0x00
	pgKeepStructure = 0x10
	pgRecurse       = 0x20
	pgIncAttributes = 0x40
	pgAllStates     = 0x3F
	pgNextProperty  = 1

	// wxPG_PREFIX_NONE and wxPG_BASE_HEX.
	pgPrefixNone = 0
	pgBaseHex    = 16
)

// mainFrame holds the state FormMain keeps in member variables.
type mainFrame struct {
	frame      *wx.Frame
	panel      *wx.Panel
	topSizer   *wx.BoxSizer
	manager    *wx.PropertyGridManager
	grid       *wx.PropertyGrid
	iconBundle *wx.BitmapBundle

	itemEnable     *wx.MenuItem
	itemFreeze     *wx.MenuItem
	itemCatColours *wx.MenuItem
	itemVetoDrag   *wx.MenuItem

	sampleEditor unsafe.Pointer

	storedValues *wx.Variant
	savedState   string

	hasHeader           bool
	labelEditingEnabled bool
	xrcHandlersAdded    bool

	logWindow *wx.LogWindow

	popup *gridPopup
}

func main() {
	wx.Run(func() bool {
		f := newMainFrame("wxPropertyGrid Sample")
		if f == nil {
			return false
		}
		f.frame.Show()
		wx.SetTopWindow(f.frame.Ptr())
		return true
	})
}

// newMainFrame is FormMain::FormMain().
func newMainFrame(title string) *mainFrame {
	style := wx.MINIMIZE_BOX | wx.MAXIMIZE_BOX | wx.RESIZE_BORDER | wx.SYSTEM_MENU |
		wx.CAPTION | wx.TAB_TRAVERSAL | wx.CLOSE_BOX
	frame := wx.NewFrame(nil, -1, title, -1, -1, -1, -1, style)
	if frame == nil {
		return nil
	}

	f := &mainFrame{frame: frame}

	// SetIcon(wxICON(sample)): the sample.xpm artwork is shipped as sample.svg.
	if bundle := wx.BitmapBundle_FromSVGFile("sample.svg", 32, 32); bundle != nil {
		f.iconBundle = bundle
		if bmp := f.iconBundle.GetBitmap(32, 32); bmp != nil {
			if icon := wx.NewIconDefault(); icon != nil {
				icon.CopyFromBitmap(bmp.Ptr())
				frame.SetIcon(icon.Ptr())
			}
		}
	}

	frame.SetClientSize(frame.FromDIP(600), frame.FromDIP(800))
	frame.Centre(wx.BOTH)

	// This is here to really test the wxImageFileProperty.
	wx.InitAllImageHandlers()

	f.buildMenus()

	frame.CreateStatusBar(1, 0, wx.ID_ANY)
	frame.SetStatusText("", 0)

	// Register all editors (SpinCtrl etc.).
	wx.PropertyGridInterface_RegisterAdditionalEditors()

	// Register our sample custom editor. The C++ version adds three buttons to
	// the text ctrl; the callback ABI for kwxPGTextCtrlEditor needs C function
	// pointers, so a text-ctrl-only editor is registered instead.
	if editor := wx.NewKwxPGTextCtrlEditor(nil); editor != nil {
		f.sampleEditor = wx.PropertyGrid_RegisterEditorClass(editor.Ptr(), false)
	}

	f.panel = wx.NewPanel(frame.Ptr(), wx.ID_ANY, -1, -1, -1, -1, wx.TAB_TRAVERSAL)
	if f.panel == nil {
		return nil
	}

	f.createGrid(
		wx.PG_BOLD_MODIFIED|wx.PG_SPLITTER_AUTO_CENTER|wx.PG_AUTO_SORT|
			wx.PG_TOOLBAR|wx.PG_DESCRIPTION,
		wx.PG_EX_NATIVE_DOUBLE_BUFFERING|wx.PG_EX_MODE_BUTTONS|
			wx.PG_EX_MULTIPLE_SELECTION)
	if f.manager == nil {
		return nil
	}

	f.topSizer = wx.NewBoxSizer(wx.VERTICAL)
	if f.topSizer == nil {
		return nil
	}

	gridFlags := wx.NewSizerFlags(1)
	gridFlags.Expand()
	f.topSizer.AddWindowFlags(f.manager.Ptr(), gridFlags.Ptr())
	gridFlags.Delete()

	// Button for tab traversal testing.
	btnSizer := wx.NewBoxSizer(wx.HORIZONTAL)
	if b := wx.NewButton(f.panel.Ptr(), wx.ID_ANY,
		"Should be able to move here with Tab", -1, -1, -1, -1, 0); b != nil {
		flags := wx.NewSizerFlags(1)
		flags.DoubleBorder(wx.ALL)
		btnSizer.AddWindowFlags(b.Ptr(), flags.Ptr())
		flags.Delete()
	}
	if b := wx.NewButton(f.panel.Ptr(), idShowPopup, "Show Popup", -1, -1, -1, -1, 0); b != nil {
		flags := wx.NewSizerFlags(1)
		flags.DoubleBorder(wx.ALL)
		btnSizer.AddWindowFlags(b.Ptr(), flags.Ptr())
		flags.Delete()
	}
	btnFlags := wx.NewSizerFlags(0)
	btnFlags.Border(wx.ALL)
	btnFlags.Expand()
	f.topSizer.AddSizerFlags(btnSizer.Ptr(), btnFlags.Ptr())
	btnFlags.Delete()

	f.panel.SetSizer(f.topSizer.Ptr(), true)
	f.topSizer.SetSizeHints(f.panel.Ptr())
	f.panel.Layout()

	// Create log window.
	f.logWindow = wx.NewLogWindow(frame.Ptr(), "Log Messages", false, false)
	if f.logWindow != nil {
		if lf := wx.FrameFromPtr(f.logWindow.GetFrame()); lf != nil {
			pos := frame.GetPosition()
			size := frame.GetSize()
			lf.Move(pos.X+size.Width+10, pos.Y, wx.SIZE_USE_EXISTING)
		}
		f.logWindow.Show(true)
	}

	f.bindEvents()
	return f
}

// buildMenus is the menu bar construction in FormMain::FormMain().
func (f *mainFrame) buildMenus() {
	menuFile := wx.NewMenu("", wx.MENU_TEAROFF)
	menuTry := wx.NewMenu("", 0)
	menuTools1 := wx.NewMenu("", 0)
	menuTools2 := wx.NewMenu("", 0)
	menuHelp := wx.NewMenu("", 0)
	if menuFile == nil || menuTry == nil || menuTools1 == nil ||
		menuTools2 == nil || menuHelp == nil {
		return
	}

	menuHelp.Append(idAbout, "&About", "Show about dialog", false)

	menuTools1.Append(idAppendProp, "Append New Property", "", false)
	menuTools1.Append(idAppendCat, "Append New Category\tCtrl-S", "", false)
	menuTools1.AppendSeparator()
	menuTools1.Append(idInsertProp, "Insert New Property\tCtrl-I", "", false)
	menuTools1.Append(idInsertCat, "Insert New Category\tCtrl-W", "", false)
	menuTools1.AppendSeparator()
	menuTools1.Append(idDelete, "Delete Selected", "", false)
	menuTools1.Append(idDeleteR, "Delete Random", "", false)
	menuTools1.Append(idDeleteAll, "Delete All", "", false)
	menuTools1.AppendSeparator()
	menuTools1.Append(idSetBgColour, "Set Bg Colour", "", false)
	menuTools1.Append(idSetBgColourRecur, "Set Bg Colour (Recursively)", "", false)
	menuTools1.Append(idUnspecify, "Set Value to Unspecified", "", false)
	menuTools1.AppendSeparator()
	f.itemEnable = wx.NewMenuItemEx(idEnable, "Enable",
		"Toggles item's enabled state.", wx.ITEM_NORMAL, nil)
	if f.itemEnable != nil {
		menuTools1.AppendItem(f.itemEnable.Ptr())
		f.itemEnable.Enable(false)
	}
	menuTools1.Append(idHide, "Hide", "Hides a property", false)
	menuTools1.Append(idSetReadOnly, "Set as Read-Only", "Set property as read-only", false)

	menuTools2.Append(idIterate1, "Iterate Over Properties", "", false)
	menuTools2.Append(idIterate2, "Iterate Over Visible Items", "", false)
	menuTools2.Append(idIterate3, "Reverse Iterate Over Properties", "", false)
	menuTools2.Append(idIterate4, "Iterate Over Categories", "", false)
	menuTools2.AppendSeparator()
	menuTools2.Append(idOnExtendedKeyNav, "Extend Keyboard Navigation",
		"This will set Enter to navigate to next property, "+
			"and allows arrow keys to navigate even when in editor control.", false)
	menuTools2.AppendSeparator()
	menuTools2.Append(idSetPropertyValue, "Set Property Value", "", false)
	menuTools2.Append(idClearModif, "Clear Modified Status",
		"Clears wxPG_MODIFIED flag from all properties.", false)
	menuTools2.AppendSeparator()
	f.itemFreeze = wx.NewMenuItemEx(idFreeze, "Freeze",
		"Disables painting, auto-sorting, etc.", wx.ITEM_CHECK, nil)
	if f.itemFreeze != nil {
		menuTools2.AppendItem(f.itemFreeze.Ptr())
	}
	menuTools2.AppendSeparator()
	menuTools2.Append(idDumpList, "Display Values as wxVariant List",
		"Tests GetAllValues method and wxVariant conversion.", false)
	menuTools2.AppendSeparator()
	menuTools2.Append(idGetValues, "Get Property Values", "Stores all property values.", false)
	menuTools2.Append(idSetValues, "Set Property Values",
		"Reverts property values to those last stored.", false)
	menuTools2.Append(idSetValues2, "Set Property Values 2",
		"Adds property values that should not initially be as items (so new items are created).", false)
	menuTools2.AppendSeparator()
	menuTools2.Append(idSaveState, "Save Editable State", "", false)
	menuTools2.Append(idRestoreState, "Restore Editable State", "", false)
	menuTools2.AppendSeparator()
	menuTools2.Append(idEnableCommonValues, "Enable Common Value",
		"Enable values that are common to all properties, for selected property.", false)
	menuTools2.AppendSeparator()
	menuTools2.Append(idCollapse, "Collapse Selected", "", false)
	menuTools2.Append(idCollapseAll, "Collapse All", "", false)
	menuTools2.AppendSeparator()
	menuTools2.Append(idInsertPage, "Add Page", "", false)
	menuTools2.Append(idRemovePage, "Remove Page", "", false)
	menuTools2.AppendSeparator()
	menuTools2.Append(idFitColumns, "Fit Columns", "", false)
	f.itemVetoDrag = wx.NewMenuItemEx(idVetoColDrag, "Veto Column Dragging", "",
		wx.ITEM_CHECK, nil)
	if f.itemVetoDrag != nil {
		menuTools2.AppendItem(f.itemVetoDrag.Ptr())
	}
	menuTools2.AppendSeparator()
	menuTools2.Append(idChangeFlagsItems, "Change Children of FlagsProp", "", false)
	menuTools2.AppendSeparator()
	menuTools2.Append(idTestInsertChoice, "Test InsertPropertyChoice", "", false)
	menuTools2.Append(idTestDeleteChoice, "Test DeletePropertyChoice", "", false)
	menuTools2.AppendSeparator()
	menuTools2.Append(idSetSpinCtrlEditor, "Use SpinCtrl Editor", "", false)
	menuTools2.Append(idTestReplace, "Test ReplaceProperty", "", false)

	menuTry.Append(idSelectStyle, "Set Window Style",
		"Select window style flags used by the grid.", false)
	menuTry.AppendCheckItem(idEnableLabelEditing, "Enable label editing",
		"This calls wxPropertyGrid::MakeColumnEditable(0)")
	menuTry.Check(idEnableLabelEditing, f.labelEditingEnabled)
	menuTry.AppendCheckItem(idShowHeader, "Enable header",
		"This calls wxPropertyGridManager::ShowHeader()")
	menuTry.Check(idShowHeader, f.hasHeader)
	menuTry.AppendSeparator()
	menuTry.AppendRadioItem(idColourScheme1, "Standard Colour Scheme", "")
	menuTry.AppendRadioItem(idColourScheme2, "White Colour Scheme", "")
	menuTry.AppendRadioItem(idColourScheme3, ".NET Colour Scheme", "")
	menuTry.AppendRadioItem(idColourScheme4, "Cream Colour Scheme", "")
	menuTry.AppendSeparator()
	f.itemCatColours = wx.NewMenuItemEx(idCatColours, "Category Specific Colours",
		"Switches between category-specific cell colours and default scheme "+
			"(actually done using SetPropertyTextColour and SetPropertyBackgroundColour).",
		wx.ITEM_CHECK, nil)
	if f.itemCatColours != nil {
		menuTry.AppendItem(f.itemCatColours.Ptr())
	}
	menuTry.AppendSeparator()
	menuTry.AppendCheckItem(idStaticLayout, "Static Layout",
		"Switches between user-modifiable and static layouts.")
	menuTry.AppendCheckItem(idBoolCheckbox, "Render Boolean values as checkboxes",
		"Renders Boolean values as checkboxes")
	menuTry.Append(idSetColumns, "Set Number of Columns", "", false)
	menuTry.Append(idSetVirtWidth, "Set Virtual Width", "", false)
	menuTry.AppendCheckItem(idSetPGDisabled, "Disable Grid", "")

	menuFile.Append(idRunMinimal, "Run Minimal Sample", "", false)
	menuFile.Append(idTestXRC, "Load from XRC\tCtrl-R", "", false)
	menuFile.AppendSeparator()
	menuFile.Append(idQuit, "E&xit\tAlt-X", "Quit this program", false)

	menuBar := wx.NewMenuBar(0)
	if menuBar == nil {
		return
	}
	menuBar.Append(menuFile.Ptr(), "&File")
	menuBar.Append(menuTry.Ptr(), "&Try These!")
	menuBar.Append(menuTools1.Ptr(), "&Basic")
	menuBar.Append(menuTools2.Ptr(), "&Advanced")
	menuBar.Append(menuHelp.Ptr(), "&Help")
	f.frame.SetMenuBar(menuBar.Ptr())
}

// setStatusText mirrors GetStatusBar()->SetStatusText(text).
func (f *mainFrame) setStatusText(text string) {
	f.frame.SetStatusText(text, 0)
}

// onAbout is FormMain::OnAbout().
func (f *mainFrame) onAbout(evt *wx.Event) {
	if evt == nil {
		return
	}
	pi := wx.PlatformInfoFromPtr(wx.PlatformInfo_Get())
	toolkit := ""
	if pi != nil {
		toolkit = fmt.Sprintf("%s %d.%d.%d", pi.GetPortIdName(),
			pi.GetToolkitMajorVersion(), pi.GetToolkitMinorVersion(),
			pi.GetToolkitMicroVersion())
	}
	msg := fmt.Sprintf("wxPropertyGrid Sample <utf-8> <debug>\n\n"+
		"Programmed by %s\n\nUsing %s (%s)\n\n",
		"Jaakko Salli", wx.VersionString(), toolkit)
	wx.MessageBox(msg, "About", wx.OK|wx.ICON_INFORMATION, f.frame.Ptr(), -1, -1)
}
