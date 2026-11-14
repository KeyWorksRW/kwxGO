// Command combo is the kwxGO port of the wxWidgets "combo" sample
// (samples/combo/combo.cpp).
//
// The frame demonstrates wxComboCtrl and wxOwnerDrawnComboBox variants: two
// owner-drawn combo boxes populated with the pen-style item list, a wxComboCtrl
// with a wxListView popup, a wxGenericComboCtrl with a wxTreeCtrl popup, two
// owner-drawn combos with custom drop-button bitmaps, and a combo with a custom
// button plus a wxCheckBox as its main control. A multiline wxTextCtrl at the
// right logs the wxEVT_TEXT / wxEVT_TEXT_ENTER / wxEVT_COMBOBOX notifications
// through wxLogTextCtrl, and the File menu opens a dialog comparing the
// owner-drawn combos with native wxComboBox controls.
package main

import (
	"fmt"
	"unsafe"

	"github.com/KeyWorksRW/kwxGO/wx"
)

// Dimensions of the original sample.xpm icon, now sample.svg.
const (
	sampleIconWidth  = 32
	sampleIconHeight = 32
)

// idLogWindow is the explicit ID MyFrame gives the log wxTextCtrl.
const idLogWindow = 105

// comboFrame holds the state MyFrame keeps in member variables.
type comboFrame struct {
	frame  *wx.Frame
	logWin *wx.TextCtrl

	// logger is the wxLogTextCtrl installed as the active log target and
	// oldLogTarget is the target it replaced. MyFrame::~MyFrame() does
	// `delete wxLog::SetActiveTarget(m_logOld)`; onClose performs the same
	// cleanup before the frame closes.
	logger       *wx.LogTextCtrl
	oldLogTarget unsafe.Pointer

	// idCompare is ComboCtrl_Compare, the first custom command ID. It is
	// wxID_HIGHEST, which the binding exposes as a runtime value rather than a
	// Go constant, so it is resolved when the frame is built.
	idCompare int

	// items is the shared item list every combo in the sample is filled from.
	items []string
}

func main() {
	wx.Run(func() bool {
		f := newComboFrame()
		if f == nil {
			return false
		}
		f.frame.Show()
		wx.SetTopWindow(f.frame.Ptr())
		return true
	})
}

// newComboFrame builds the frame, its menu bar, the panel with every combo
// control and the log pane -- MyFrame::MyFrame().
func newComboFrame() *comboFrame {
	frame := wx.NewFrame(nil, wx.ID_ANY,
		"wxComboCtrl and wxOwnerDrawnComboBox Sample",
		-1, -1, -1, -1, wx.DEFAULT_FRAME_STYLE)
	if frame == nil {
		return nil
	}

	f := &comboFrame{
		frame:     frame,
		idCompare: wx.ID_HIGHEST,
		items: []string{
			"Solid",
			"Transparent",
			"Dot",
			"Long Dash",
			"Short Dash",
			"Dot Dash",
			"Backward Diagonal Hatch",
			"Cross-diagonal Hatch",
			"Forward Diagonal Hatch",
			"Cross Hatch",
			"Horizontal Hatch",
			"Vertical Hatch",
		},
	}

	// SetIcon(wxICON(sample)): the XPM icon has been converted to sample.svg.
	if icon := loadSampleIcon(); icon != nil {
		frame.SetIcon(icon.Ptr())
		icon.Delete()
	}

	f.buildMenus()

	panel := wx.NewPanel(frame.Ptr(), wx.ID_ANY, -1, -1, -1, -1, 0)
	if panel == nil {
		return nil
	}

	// Prepare the log window right away since it shows EVT_TEXTs.
	f.logWin = wx.NewTextCtrl(panel.Ptr(), idLogWindow, "",
		-1, -1, -1, 125, wx.TE_MULTILINE)
	if f.logWin == nil {
		return nil
	}
	if logger := wx.NewLogTextCtrl(f.logWin.Ptr()); logger != nil {
		f.logger = logger
		f.oldLogTarget = wx.SetActiveLogTarget(logger.Ptr())
		// logger->DisableTimestamp().
		logger.DisableTimestamp()
	}

	topSizer := wx.NewBoxSizer(wx.VERTICAL)
	topRowSizer := wx.NewBoxSizer(wx.HORIZONTAL)
	colSizer := wx.NewBoxSizer(wx.VERTICAL)
	if topSizer == nil || topRowSizer == nil || colSizer == nil {
		return nil
	}

	// -- OwnerDrawnComboBox with owner-drawn items ---------------------------
	// The C++ sample uses wxPenStyleComboBox, a wxOwnerDrawnComboBox subclass that
	// overrides OnDrawItem (pen style per item), OnDrawBackground (alternating
	// background) and OnMeasureItem (variable 36/24px rows). Those are protected
	// virtuals: the binding only exposes them through KwxOwnerDrawnComboBox's
	// SetOnDrawItem/SetOnDrawBackground/SetOnMeasureItem, which take a raw C
	// function pointer that Go cannot produce without a cgo shim in the sample.
	// The controls therefore render with the stock wxOwnerDrawnComboBox drawing.
	rowSizer := wx.NewBoxSizer(wx.HORIZONTAL)
	if rowSizer == nil {
		return nil
	}
	addWin(rowSizer, label(panel.Ptr(),
		"OwnerDrawnComboBox with owner-drawn items:"), 1, func(sf *wx.SizerFlags) {
		sf.CenterVertical()
		sf.Border(wx.RIGHT)
	})
	addSiz(colSizer, rowSizer.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})

	rowSizer = wx.NewBoxSizer(wx.HORIZONTAL)
	if rowSizer == nil {
		return nil
	}
	odc := newOwnerDrawnComboBox(panel.Ptr(), wx.CB_READONLY, f.items)
	if odc == nil {
		return nil
	}
	odc.SetSelection(0)
	addWin(rowSizer, odc.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.CenterVertical()
		sf.Border(wx.ALL)
	})
	rowSizer.AddStretchSpacer(1)
	addSiz(colSizer, rowSizer.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})

	// -- Same, but with the button on the left -------------------------------
	rowSizer = wx.NewBoxSizer(wx.HORIZONTAL)
	if rowSizer == nil {
		return nil
	}
	addWin(rowSizer, label(panel.Ptr(),
		"OwnerDrawnComboBox with owner-drawn items and button on the left:"), 1, func(sf *wx.SizerFlags) {
		sf.CenterVertical()
		sf.Border(wx.RIGHT)
	})
	addSiz(colSizer, rowSizer.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})

	rowSizer = wx.NewBoxSizer(wx.HORIZONTAL)
	if rowSizer == nil {
		return nil
	}
	odc = newOwnerDrawnComboBox(panel.Ptr(), wx.CB_READONLY, f.items)
	if odc == nil {
		return nil
	}
	odc.SetSelection(0)

	// Use a button size that is slightly smaller than the default.
	butSize := odc.GetButtonSize()
	odc.SetButtonPosition(butSize.Width-2, butSize.Height-6, wx.LEFT, 2)

	addWin(rowSizer, odc.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.CenterVertical()
		sf.Border(wx.ALL)
	})
	rowSizer.AddStretchSpacer(1)
	addSiz(colSizer, rowSizer.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})

	// -- List view / tree control combo labels -------------------------------
	rowSizer = wx.NewBoxSizer(wx.HORIZONTAL)
	if rowSizer == nil {
		return nil
	}
	addWin(rowSizer, label(panel.Ptr(),
		"List View wxComboCtrl (custom animation):"), 1, func(sf *wx.SizerFlags) {
		sf.CenterVertical()
		sf.Border(wx.RIGHT)
	})
	addWin(rowSizer, label(panel.Ptr(),
		"Tree Ctrl wxComboCtrl:"), 1, func(sf *wx.SizerFlags) {
		sf.CenterVertical()
		sf.Border(wx.RIGHT)
	})
	addSiz(colSizer, rowSizer.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})

	rowSizer = wx.NewBoxSizer(wx.HORIZONTAL)
	if rowSizer == nil {
		return nil
	}

	// Two-step creation so the contained wxTextCtrl gets wxTE_RIGHT. The C++
	// class wxComboCtrlWithCustomPopupAnim also overrides AnimateShow() to fade
	// the popup in with ShowWithEffect(wxSHOW_EFFECT_BLEND); that protected
	// virtual is only reachable through SetOnAnimateShow(), which takes a raw C
	// function pointer Go cannot produce, so the popup uses the default animation.
	cc := wx.NewKwxComboCtrlDefault(nil)
	if cc == nil {
		return nil
	}
	cc.SetTextCtrlStyle(wx.TE_RIGHT)
	if !cc.Create(panel.Ptr(), wx.ID_ANY, "", -1, -1, -1, -1, 0) {
		return nil
	}
	// Make sure we use a popup that allows focusing the listview.
	cc.UseAltPopupWindow(true)
	cc.SetPopupMinWidth(300)

	// The popup must be installed before AddString(): the C helper creates its
	// wxListView from Create(), which wxComboCtrl calls from SetPopupControl().
	iface := wx.NewKwxListViewComboPopup()
	if iface == nil {
		return nil
	}
	cc.SetPopupControl(iface.Ptr())
	for i := 0; i < 100; i++ {
		iface.AddString(fmt.Sprintf("Item %02d", i))
	}

	addWin(rowSizer, cc.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.CenterVertical()
		sf.Border(wx.ALL)
	})

	// -- Tree control combo --------------------------------------------------
	gcc := wx.NewGenericComboCtrl(panel.Ptr(), wx.ID_ANY, "",
		-1, -1, -1, -1, 0)
	if gcc == nil {
		return nil
	}
	// Make sure we use a popup that allows focusing the treectrl.
	gcc.UseAltPopupWindow(true)

	tcPopup := wx.NewKwxTreeCtrlComboPopup()
	if tcPopup == nil {
		return nil
	}
	gcc.SetPopupControl(tcPopup.Ptr())

	// Add items using wxTreeCtrl methods directly.
	rootID := wx.NewTreeItemId()
	if rootID == nil {
		return nil
	}
	defer rootID.Delete()
	tcPopup.AddRoot("<hidden_root>", rootID.Ptr())

	for i := 0; i < 4; i++ {
		groupID := wx.NewTreeItemId()
		if groupID == nil {
			return nil
		}
		tcPopup.AppendItem(rootID.Ptr(), fmt.Sprintf("Branch %02d", i), groupID.Ptr())
		for n := 0; n < 25; n++ {
			subID := wx.NewTreeItemId()
			if subID == nil {
				groupID.Delete()
				return nil
			}
			tcPopup.AppendItem(groupID.Ptr(),
				fmt.Sprintf("Subitem %02d", i*25+n), subID.Ptr())
			subID.Delete()
		}
		groupID.Delete()
	}

	gcc.SetValue("Subitem 05")

	// Move the button to the left - it makes more sense for a tree ctrl.
	gcc.SetButtonPosition(-1, -1, wx.LEFT, 0)

	addWin(rowSizer, gcc.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.CenterVertical()
		sf.Border(wx.ALL)
	})
	addSiz(colSizer, rowSizer.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})

	// -- Custom drop-button bitmaps ------------------------------------------
	wx.InitAllImageHandlers()

	rowSizer = wx.NewBoxSizer(wx.HORIZONTAL)
	if rowSizer == nil {
		return nil
	}
	addWin(rowSizer, label(panel.Ptr(),
		"OwnerDrawnComboBox with simple dropbutton graphics:"), 1, func(sf *wx.SizerFlags) {
		sf.CenterVertical()
		sf.Border(wx.RIGHT)
	})
	addSiz(colSizer, rowSizer.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})

	rowSizer = wx.NewBoxSizer(wx.HORIZONTAL)
	if rowSizer == nil {
		return nil
	}
	odc1 := newOwnerDrawnComboBox(panel.Ptr(), 0, f.items)
	odc2 := newOwnerDrawnComboBox(panel.Ptr(), 0, f.items)
	if odc1 == nil || odc2 == nil {
		return nil
	}

	// Load images from disk. The second one uses a blank button background.
	bmpNormal := loadSVGBundle("dropbutn.svg", 16, 16)
	bmpPressed := loadSVGBundle("dropbutp.svg", 16, 16)
	bmpHover := loadSVGBundle("dropbuth.svg", 16, 16)

	if bmpNormal != nil && bmpNormal.IsOk() &&
		bmpPressed != nil && bmpPressed.IsOk() &&
		bmpHover != nil && bmpHover.IsOk() {
		odc1.SetButtonBitmaps(bmpNormal.Ptr(), false,
			bmpPressed.Ptr(), bmpHover.Ptr(), nil)
		odc2.SetButtonBitmaps(bmpNormal.Ptr(), true,
			bmpPressed.Ptr(), bmpHover.Ptr(), nil)
	} else {
		wx.LogError("Dropbutton images not found")
	}

	addWin(rowSizer, odc1.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.CenterVertical()
		sf.Border(wx.ALL)
	})
	addWin(rowSizer, odc2.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.CenterVertical()
		sf.Border(wx.ALL)
	})
	addSiz(colSizer, rowSizer.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})

	// -- Custom button and custom main control -------------------------------
	rowSizer = wx.NewBoxSizer(wx.HORIZONTAL)
	if rowSizer == nil {
		return nil
	}
	addWin(rowSizer, label(panel.Ptr(),
		"wxComboCtrl with custom button and custom main control:"), 1, func(sf *wx.SizerFlags) {
		sf.CenterVertical()
		sf.Border(wx.RIGHT)
	})
	addSiz(colSizer, rowSizer.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})

	rowSizer = wx.NewBoxSizer(wx.HORIZONTAL)
	if rowSizer == nil {
		return nil
	}

	// wxFileSelectorCombo: wxCC_STD_BUTTON plus the three-dots button graphics.
	// Its OnButtonClick() opened a wxFileDialog in C++, and DoSetPopupControl()
	// was emptied to suppress the assertion. Both are protected virtuals reached
	// only through SetOnButtonClick/SetOnDoSetPopupControl, which take a raw C
	// function pointer that Go cannot produce. A list popup is installed instead
	// so the drop button still does something.
	fsc := wx.NewKwxComboCtrlWith(nil, panel.Ptr(), wx.ID_ANY, "",
		-1, -1, -1, -1, wx.CC_STD_BUTTON)
	if fsc == nil {
		return nil
	}
	if dots := loadSVGBundle("three-dots.svg", 16, 16); dots != nil && dots.IsOk() {
		fsc.SetButtonBitmaps(dots.Ptr(), true, nil, nil, nil)
	}
	if fscPopup := wx.NewKwxListViewComboPopup(); fscPopup != nil {
		fsc.SetPopupControl(fscPopup.Ptr())
	}

	// A perfectly useless control, as the popup and main control don't interact,
	// but it shows that something other than wxTextCtrl can be used for the main
	// part of wxComboCtrl too.
	//
	// The controls must be created in this order so that the custom main control
	// exists before the wxComboCtrl window is actually created: otherwise the
	// combo would needlessly create a default wxTextCtrl first.
	comboCustom := wx.NewKwxComboCtrlDefault(nil)
	if comboCustom == nil {
		return nil
	}
	cbox := wx.NewCheckBoxDefault()
	if cbox == nil {
		return nil
	}
	comboCustom.SetMainControl(cbox.Ptr())
	if !comboCustom.Create(panel.Ptr(), wx.ID_ANY, "", -1, -1, -1, -1, 0) {
		return nil
	}
	if !cbox.Create2(comboCustom.Ptr(), wx.ID_ANY,
		"Checkbox as main control", -1, -1, -1, -1, 0) {
		return nil
	}
	if windowCol := sysColour(wx.SYS_COLOUR_WINDOW); windowCol != nil {
		cbox.SetBackgroundColour(windowCol.Ptr())
		windowCol.Delete()
	}
	if customPopup := wx.NewKwxListViewComboPopup(); customPopup != nil {
		comboCustom.SetPopupControl(customPopup.Ptr())
	}

	addWin(rowSizer, fsc.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.CenterVertical()
		sf.Border(wx.ALL)
	})
	addWin(rowSizer, comboCustom.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.CenterVertical()
		sf.Border(wx.ALL)
	})
	addSiz(colSizer, rowSizer.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})

	// Make sure GetFeatures is implemented.
	wx.ComboCtrl_GetFeatures()

	// -- Right-hand column: the log messages ---------------------------------
	addSiz(topRowSizer, colSizer.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.BorderSize(wx.ALL, frame.FromDIP(2))
	})

	colSizer = wx.NewBoxSizer(wx.VERTICAL)
	if colSizer == nil {
		return nil
	}
	colSizer.AddSpacer(8)
	addWin(colSizer, label(panel.Ptr(), "Log Messages:"), 0, func(sf *wx.SizerFlags) {
		sf.BorderSize(wx.TOP|wx.LEFT, frame.FromDIP(3))
	})
	addWin(colSizer, f.logWin.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.BorderSize(wx.ALL, frame.FromDIP(3))
	})

	addSiz(topRowSizer, colSizer.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.BorderSize(wx.ALL, frame.FromDIP(2))
	})
	addSiz(topSizer, topRowSizer.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
	})

	panel.SetSizer(topSizer.Ptr(), true)
	topSizer.SetSizeHints(panel.Ptr())

	frame.Fit()
	frame.Centre(wx.BOTH)

	f.bindEvents()
	return f
}

// buildMenus recreates the File and Help menus verbatim.
func (f *comboFrame) buildMenus() {
	fileMenu := wx.NewMenu("", 0)
	if fileMenu == nil {
		return
	}
	fileMenu.Append(f.idCompare, "&Compare against wxComboBox...",
		"Show some wxOwnerDrawnComboBoxes side-by-side with native wxComboBoxes.",
		false)
	fileMenu.AppendSeparator()
	fileMenu.Append(wx.ID_EXIT, "E&xit\tAlt-X", "Quit this program", false)

	// The "About" item belongs in the help menu.
	helpMenu := wx.NewMenu("", 0)
	if helpMenu == nil {
		return
	}
	helpMenu.Append(wx.ID_ABOUT, "&About\tF1", "Show about dialog", false)

	menuBar := wx.NewMenuBar(0)
	if menuBar == nil {
		return
	}
	menuBar.Append(fileMenu.Ptr(), "&File")
	menuBar.Append(helpMenu.Ptr(), "&Help")
	f.frame.SetMenuBar(menuBar.Ptr())
}

// bindEvents wires the C++ event table onto the frame.
func (f *comboFrame) bindEvents() {
	// EVT_TEXT / EVT_TEXT_ENTER / EVT_COMBOBOX with wxID_ANY.
	f.frame.Bind(wx.EVT_TEXT, f.onComboBoxUpdate)
	f.frame.Bind(wx.EVT_TEXT_ENTER, f.onComboBoxUpdate)
	f.frame.Bind(wx.EVT_COMBOBOX, f.onComboBoxUpdate)

	f.frame.BindId(wx.EVT_MENU, f.idCompare, f.onShowComparison)
	f.frame.BindId(wx.EVT_MENU, wx.ID_EXIT, f.onQuit)
	f.frame.BindId(wx.EVT_MENU, wx.ID_ABOUT, f.onAbout)

	f.frame.Bind(wx.EVT_IDLE, f.onIdle)
	f.frame.Bind(wx.EVT_CLOSE_WINDOW, f.onClose)
}

// onComboBoxUpdate is MyFrame::OnComboBoxUpdate(): it logs the combo control
// notifications and ignores everything else (in particular the log window's own
// EVT_TEXT, which would recurse).
func (f *comboFrame) onComboBoxUpdate(evt *wx.Event) {
	if evt == nil {
		return
	}

	// Don't show messages for the log output window (it'll crash).
	if !isKindOf(evt.GetEventObject(), "wxComboCtrl") {
		return
	}

	cmd := wx.AsCommandEvent(evt)
	if cmd == nil {
		return
	}

	switch evt.GetEventType() {
	case wx.EVT_COMBOBOX:
		wx.LogMessage(fmt.Sprintf("EVT_COMBOBOX(id=%d,selection=%d)",
			cmd.GetId(), cmd.GetSelection()))
	case wx.EVT_TEXT:
		wx.LogMessage(fmt.Sprintf("EVT_TEXT(id=%d,string=\"%s\")",
			cmd.GetId(), cmd.GetString()))
	case wx.EVT_TEXT_ENTER:
		wx.LogMessage(fmt.Sprintf("EVT_TEXT_ENTER(id=%d,string=\"%s\")",
			cmd.GetId(), cmd.GetString()))
	}
}

// onShowComparison is MyFrame::OnShowComparison().
func (f *comboFrame) onShowComparison(evt *wx.Event) {
	if evt == nil {
		return
	}
	showCompareDialog(f.frame.Ptr(), f.items)
}

// onQuit is MyFrame::OnQuit(). The menu handler is not a close handler, so
// closing the frame here is what the C++ code does.
func (f *comboFrame) onQuit(evt *wx.Event) {
	if evt == nil {
		return
	}
	f.frame.Close(true)
}

// onAbout is MyFrame::OnAbout().
func (f *comboFrame) onAbout(evt *wx.Event) {
	if evt == nil {
		return
	}
	wx.MessageBox(fmt.Sprintf(
		"Welcome to %s!\n"+
			"\n"+
			"This is the wxWidgets wxComboCtrl and wxOwnerDrawnComboBox sample\n"+
			"running under %s.",
		wx.VersionString(),
		wx.GetOsDescription()),
		"About wxComboCtrl sample",
		wx.OK|wx.ICON_INFORMATION,
		f.frame.Ptr(), -1, -1)
}

// onClose performs the teardown MyFrame::~MyFrame() does: it restores the log
// target that was active before m_logWin was installed and deletes the logger.
// Go has no frame destructor, so this runs on the close event; skipping the
// event afterwards lets the frame close and be destroyed as usual.
func (f *comboFrame) onClose(evt *wx.Event) {
	if evt == nil {
		return
	}
	if f.logger != nil {
		// delete wxLog::SetActiveTarget(m_logOld)
		wx.SetActiveLogTarget(f.oldLogTarget)
		f.logger.Delete()
		f.logger = nil
	}
	evt.Skip(true)
}

// onIdle is MyFrame::OnIdle(), which only skips the event.
func (f *comboFrame) onIdle(evt *wx.Event) {
	if evt == nil {
		return
	}
	evt.Skip(true)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// sizerTarget is implemented by every sizer wrapper through its embedded
// wx.Sizer, so box and static-box sizers can both be filled by the helpers
// below.
type sizerTarget interface {
	AddWindowFlags(window unsafe.Pointer, flags unsafe.Pointer)
	AddSizerFlags(sizer unsafe.Pointer, flags unsafe.Pointer)
}

// addWin adds a window to the sizer with the given wxSizerFlags modifiers.
func addWin(s sizerTarget, window unsafe.Pointer, proportion int, apply func(*wx.SizerFlags)) {
	if s == nil || window == nil {
		return
	}
	flags := wx.NewSizerFlags(proportion)
	if flags == nil {
		return
	}
	apply(flags)
	s.AddWindowFlags(window, flags.Ptr())
	flags.Delete()
}

// addSiz adds a nested sizer with the given wxSizerFlags modifiers.
func addSiz(s sizerTarget, child unsafe.Pointer, proportion int, apply func(*wx.SizerFlags)) {
	if s == nil || child == nil {
		return
	}
	flags := wx.NewSizerFlags(proportion)
	if flags == nil {
		return
	}
	apply(flags)
	s.AddSizerFlags(child, flags.Ptr())
	flags.Delete()
}

// label creates a wxStaticText and returns its raw pointer, or nil when the
// control could not be created. Callers pass the result straight to addWin.
func label(parent unsafe.Pointer, text string) unsafe.Pointer {
	st := wx.NewStaticText(parent, wx.ID_ANY, text, -1, -1, -1, -1, 0)
	if st == nil {
		return nil
	}
	return st.Ptr()
}

// isKindOf reports whether the native object behind ptr is a className or a
// derived class -- the C++ `event.GetEventObject()->IsKindOf(CLASSINFO(...))`.
func isKindOf(ptr unsafe.Pointer, className string) bool {
	if ptr == nil {
		return false
	}
	obj := wx.ObjectFromPtr(ptr)
	if obj == nil {
		return false
	}
	info := wx.ClassInfo_FindClass(className)
	if info == nil {
		return false
	}
	return obj.IsKindOf(info)
}

// newOwnerDrawnComboBox creates a wxOwnerDrawnComboBox and fills it from items.
// The C++ constructor takes the same wxArrayString, so the resulting item order
// is identical (wxCB_SORT sorts in both paths).
func newOwnerDrawnComboBox(parent unsafe.Pointer, style int, items []string) *wx.OwnerDrawnComboBox {
	odc := wx.NewOwnerDrawnComboBox(parent, wx.ID_ANY, "", -1, -1, -1, -1, style)
	if odc == nil {
		return nil
	}
	for _, item := range items {
		odc.Append(item)
	}
	return odc
}

// sysColour is wxSystemSettings::GetColour(index); the caller owns the result.
func sysColour(index int) *wx.Colour {
	col := wx.NewColourEmpty()
	if col == nil {
		return nil
	}
	wx.SystemSettings_GetColour(index, col.Ptr())
	return col
}

// loadSVGBundle loads an SVG drop-button graphic as a wxBitmapBundle. The
// returned wrapper is owned by the caller and kept alive for the event loop.
func loadSVGBundle(name string, width int, height int) *wx.BitmapBundle {
	return wx.BitmapBundle_FromSVGFile(name, width, height)
}

// loadSampleIcon loads the SVG rendition of the sample's XPM icon and converts
// it to a wxIcon for the frame.
func loadSampleIcon() *wx.Icon {
	bundle := loadSVGBundle("sample.svg", sampleIconWidth, sampleIconHeight)
	if bundle == nil || !bundle.IsOk() {
		return nil
	}

	bitmap := bundle.GetBitmap(sampleIconWidth, sampleIconHeight)
	if bitmap == nil {
		return nil
	}
	defer bitmap.Delete()

	icon := wx.NewIconDefault()
	if icon == nil {
		return nil
	}
	icon.CopyFromBitmap(bitmap.Ptr())
	if !icon.IsOk() {
		icon.Delete()
		return nil
	}
	return icon
}
