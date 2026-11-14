// Command collpane is the kwxGO port of the wxWidgets "collpane" sample
// (samples/collpane/collpane.cpp).
//
// The frame holds a wxCollapsiblePane whose pane contains static texts, a text
// control, a button and a nested (empty) wxCollapsiblePane, plus a
// wxCollapsibleHeaderCtrl below it. The Pane menu collapses, expands and
// relabels the pane, opens a test dialog and shows the About box; the
// update-UI handlers keep Collapse/Expand enabled only in the direction that
// is currently possible. The test dialog demonstrates
// wxCollapsiblePane/wxCollapsibleHeaderCtrl events and re-aligns one sizer item
// at runtime.
package main

import (
	"fmt"
	"unsafe"

	"github.com/KeyWorksRW/kwxGO/wx"
)

// Command IDs, mirroring the anonymous enum in collpane.cpp. The two standard
// IDs reuse the wxWidgets values, so the menu labels come from wxWidgets'
// stock labels exactly as in the C++ sample.
const (
	idPaneCollapse = 100
	idPaneExpand   = 101
	idPaneSetLabel = 102
	idPaneShowDlg  = 103
	idPaneButton   = 104
	idPaneTextCtrl = 105
)

// idPaneToggleStatus mirrors PANEDLG_TOGGLESTATUS_BTN = wxID_HIGHEST. It is a
// var because wx.ID_HIGHEST is resolved through cgo at package init.
var idPaneToggleStatus = wx.ID_HIGHEST

// Dimensions of the original sample.xpm icon, converted to sample.svg.
const (
	sampleIconWidth  = 32
	sampleIconHeight = 32
)

// collpaneFrame holds the state MyFrame keeps in member variables.
type collpaneFrame struct {
	frame    *wx.Frame
	collPane *wx.CollapsiblePane
	header   *wx.CollapsibleHeaderCtrl
}

func main() {
	wx.Run(func() bool {
		f := newCollpaneFrame()
		if f == nil {
			return false
		}
		f.frame.Show()
		wx.SetTopWindow(f.frame.Ptr())
		return true
	})
}

// newCollpaneFrame builds the frame, its icon, status bar, menu bar and the
// contents of the collapsible pane -- MyFrame::MyFrame().
func newCollpaneFrame() *collpaneFrame {
	frame := wx.NewFrame(nil, wx.ID_ANY, "wxCollapsiblePane sample",
		-1, -1, 420, 300, wx.DEFAULT_FRAME_STYLE)
	if frame == nil {
		return nil
	}

	f := &collpaneFrame{frame: frame}

	// SetIcon(wxICON(sample)): the XPM icon has been converted to sample.svg.
	if icon := loadSampleIcon(); icon != nil {
		frame.SetIcon(icon.Ptr())
		icon.Delete()
	}

	// CreateStatusBar(2), with the default style and id.
	frame.CreateStatusBar(2, wx.STB_DEFAULT_STYLE, 0)

	// m_collPane = new wxCollapsiblePane(this, -1, "...")
	// The constructor's default style is wxCP_DEFAULT_STYLE.
	f.collPane = wx.NewCollapsiblePane(frame.Ptr(), -1, "This is a wxCollapsiblePane",
		-1, -1, -1, -1, wx.CP_DEFAULT_STYLE)
	if f.collPane == nil {
		return nil
	}

	paneWindow := wx.WindowFromPtr(f.collPane.GetPane())
	if paneWindow == nil {
		return nil
	}

	paneSizer := wx.NewBoxSizer(wx.HORIZONTAL)
	paneSubSizer := wx.NewBoxSizer(wx.VERTICAL)
	if paneSizer == nil || paneSubSizer == nil {
		return nil
	}

	paneSizer.AddSpacer(20)
	addSizer(paneSizer, paneSubSizer.Ptr(), 1, nil)

	staticText1 := wx.NewStaticText(paneWindow.Ptr(), -1, "Static text",
		-1, -1, -1, -1, 0)
	if staticText1 == nil {
		return nil
	}
	addWindow(paneSubSizer, staticText1.Ptr(), 0, leftBorder)

	staticText2 := wx.NewStaticText(paneWindow.Ptr(), -1, "Yet another one!",
		-1, -1, -1, -1, 0)
	if staticText2 == nil {
		return nil
	}
	addWindow(paneSubSizer, staticText2.Ptr(), 0, leftBorder)

	textCtrl := wx.NewTextCtrl(paneWindow.Ptr(), idPaneTextCtrl, "Text control",
		-1, -1, -1, -1, 0)
	if textCtrl == nil {
		return nil
	}
	addWindow(paneSubSizer, textCtrl.Ptr(), 0, leftBorder)

	alignButton := wx.NewButton(paneWindow.Ptr(), idPaneButton, "Press to align right",
		-1, -1, -1, -1, 0)
	if alignButton == nil {
		return nil
	}
	addWindow(paneSubSizer, alignButton.Ptr(), 0, leftBorder)

	nestedPane := wx.NewCollapsiblePane(paneWindow.Ptr(), -1, "Nested empty wxCollapsiblePane",
		-1, -1, -1, -1, wx.CP_DEFAULT_STYLE)
	if nestedPane == nil {
		return nil
	}
	addWindow(paneSubSizer, nestedPane.Ptr(), 0, expandBorder)

	paneWindow.SetSizer(paneSizer.Ptr(), true)

	// m_collHeaderCtrl = new wxCollapsibleHeaderCtrl(this, wxID_ANY, "...")
	// The C++ constructor's default style is wxBORDER_NONE.
	f.header = wx.NewCollapsibleHeaderCtrl(frame.Ptr(), wx.ID_ANY,
		"Collapsed wxCollapsibleHeaderCtrl", -1, -1, -1, -1, wx.BORDER_NONE)
	if f.header == nil {
		return nil
	}

	sizerTop := wx.NewBoxSizer(wx.VERTICAL)
	if sizerTop == nil {
		return nil
	}
	addWindow(sizerTop, f.collPane.Ptr(), 1, expandOnly)
	addWindow(sizerTop, f.header.Ptr(), 0, nil)
	frame.SetSizer(sizerTop.Ptr(), true)

	// The C++ sample installs its event table statically; installing the
	// handlers once the widgets exist is the equivalent Go ordering.
	f.buildMenus()

	return f
}

// buildMenus recreates the Pane and Help menus and installs every handler of
// MyFrame's event table.
func (f *collpaneFrame) buildMenus() {
	paneMenu := wx.NewMenu("", 0)
	if paneMenu == nil {
		return
	}
	paneMenu.Append(idPaneCollapse, "Collapse\tCtrl-C", "", false)
	paneMenu.Append(idPaneExpand, "Expand\tCtrl-E", "", false)
	paneMenu.AppendSeparator()
	paneMenu.Append(idPaneSetLabel, "Set label...\tCtrl-L", "", false)
	paneMenu.AppendSeparator()
	paneMenu.Append(idPaneShowDlg, "Show dialog...\tCtrl-S", "", false)
	paneMenu.AppendSeparator()
	// Append(PANE_QUIT) uses the stock label and help string for wxID_EXIT.
	paneMenu.Append(wx.ID_EXIT, "", "", false)

	helpMenu := wx.NewMenu("", 0)
	if helpMenu == nil {
		return
	}
	helpMenu.Append(wx.ID_ABOUT, "", "", false)

	menuBar := wx.NewMenuBar(0)
	if menuBar == nil {
		return
	}
	menuBar.Append(paneMenu.Ptr(), "&Pane")
	menuBar.Append(helpMenu.Ptr(), "&Help")
	f.frame.SetMenuBar(menuBar.Ptr())

	f.frame.BindId(wx.EVT_MENU, idPaneCollapse, f.onCollapse)
	f.frame.BindId(wx.EVT_MENU, idPaneExpand, f.onExpand)
	f.frame.BindId(wx.EVT_MENU, idPaneSetLabel, f.onSetLabel)
	f.frame.BindId(wx.EVT_MENU, idPaneShowDlg, f.onShowDialog)
	f.frame.BindId(wx.EVT_MENU, wx.ID_ABOUT, f.onAbout)
	f.frame.BindId(wx.EVT_MENU, wx.ID_EXIT, f.quit)

	f.frame.BindId(wx.EVT_UPDATE_UI, idPaneCollapse, f.onCollapseUpdateUI)
	f.frame.BindId(wx.EVT_UPDATE_UI, idPaneExpand, f.onExpandUpdateUI)

	f.header.Bind(wx.EVT_COLLAPSIBLEHEADER_CHANGED, f.onCollapsibleHeaderChanged)
}

func (f *collpaneFrame) quit(evt *wx.Event) {
	if evt == nil {
		return
	}
	f.frame.Close(true)
}

func (f *collpaneFrame) onCollapse(evt *wx.Event) {
	if evt == nil {
		return
	}
	f.collPane.Collapse(true)
}

func (f *collpaneFrame) onExpand(evt *wx.Event) {
	if evt == nil {
		return
	}
	f.collPane.Expand()
}

// onSetLabel implements MyFrame::OnSetLabel: it writes the text entered in the
// wxGetTextFromUser dialog, which is empty when the dialog is cancelled, onto
// the collapsible pane unconditionally.
func (f *collpaneFrame) onSetLabel(evt *wx.Event) {
	if evt == nil {
		return
	}

	// wxGetTextFromUser("Enter new label", wxGetTextFromUserPromptStr, GetLabel()):
	// no parent, default x/y and wxCENTRE.
	text := wx.GetTextFromUserEx("Enter new label", "Input Text",
		f.collPane.GetLabel(), nil, -1, -1, wx.CENTRE)

	f.collPane.SetLabel(text)
}

func (f *collpaneFrame) onShowDialog(evt *wx.Event) {
	if evt == nil {
		return
	}
	showTestDialog(f.frame.Ptr())
}

func (f *collpaneFrame) onAbout(evt *wx.Event) {
	if evt == nil {
		return
	}

	info := wx.NewAboutDialogInfo()
	if info == nil {
		return
	}
	defer info.Delete()

	info.SetName("wxCollapsiblePane sample")
	info.SetDescription("This sample program demonstrates usage of wxCollapsiblePane")
	info.SetCopyright("(C) 2006 Francesco Montorsi <frm@users.sourceforge.net>")

	wx.AboutBox(info.Ptr(), f.frame.Ptr())
}

func (f *collpaneFrame) onCollapseUpdateUI(evt *wx.Event) {
	if evt == nil {
		return
	}
	update := wx.AsUpdateUIEvent(evt)
	if update == nil {
		return
	}
	update.Enable(!f.collPane.IsCollapsed())
}

func (f *collpaneFrame) onExpandUpdateUI(evt *wx.Event) {
	if evt == nil {
		return
	}
	update := wx.AsUpdateUIEvent(evt)
	if update == nil {
		return
	}
	update.Enable(f.collPane.IsCollapsed())
}

// onCollapsibleHeaderChanged is MyFrame::OnCollapsibleHeaderChanged; in the C++
// sample the event table entry is on the frame and matches wxID_ANY, so binding
// it on the header control itself is equivalent.
func (f *collpaneFrame) onCollapsibleHeaderChanged(evt *wx.Event) {
	if evt == nil {
		return
	}

	state := "Expanded"
	if f.header.IsCollapsed() {
		state = "Collapsed"
	}
	f.header.SetLabel(fmt.Sprintf("%s wxCollapsibleHeaderCtrl", state))
}

// sizerTarget is implemented by every sizer wrapper through its embedded
// wx.Sizer, letting the helpers below fill box and grid sizers alike.
type sizerTarget interface {
	AddWindowFlags(window unsafe.Pointer, flags unsafe.Pointer)
	AddSizerFlags(sizer unsafe.Pointer, flags unsafe.Pointer)
}

// addWindow adds a window with wxSizerFlags(proportion) plus the modifiers
// applied by flagsFn (nil for plain default flags).
func addWindow(s sizerTarget, window unsafe.Pointer, proportion int, flagsFn func(*wx.SizerFlags)) {
	flags := wx.NewSizerFlags(proportion)
	if flags == nil {
		return
	}
	if flagsFn != nil {
		flagsFn(flags)
	}
	s.AddWindowFlags(window, flags.Ptr())
	flags.Delete()
}

// addSizer adds a nested sizer with wxSizerFlags(proportion) plus modifiers.
func addSizer(s sizerTarget, child unsafe.Pointer, proportion int, flagsFn func(*wx.SizerFlags)) {
	flags := wx.NewSizerFlags(proportion)
	if flags == nil {
		return
	}
	if flagsFn != nil {
		flagsFn(flags)
	}
	s.AddSizerFlags(child, flags.Ptr())
	flags.Delete()
}

// wxSizerFlags().Left().Border()
func leftBorder(sf *wx.SizerFlags) {
	sf.Left()
	sf.Border(wx.ALL)
}

// wxSizerFlags().Expand().Border()
func expandBorder(sf *wx.SizerFlags) {
	sf.Expand()
	sf.Border(wx.ALL)
}

// wxSizerFlags().Expand()
func expandOnly(sf *wx.SizerFlags) {
	sf.Expand()
}

// wxSizerFlags().Border()
func borderOnly(sf *wx.SizerFlags) {
	sf.Border(wx.ALL)
}

// wxSizerFlags().Right().Border()
func rightBorder(sf *wx.SizerFlags) {
	sf.Right()
	sf.Border(wx.ALL)
}

// wxSizerFlags(1).Expand().Border()
func growExpandBorder(sf *wx.SizerFlags) {
	sf.Proportion(1)
	sf.Expand()
	sf.Border(wx.ALL)
}

// loadSampleIcon loads the SVG rendition of the sample's XPM icon and converts
// it to a wxIcon for the frame.
func loadSampleIcon() *wx.Icon {
	bundle := wx.BitmapBundle_FromSVGFile("sample.svg", sampleIconWidth, sampleIconHeight)
	if bundle == nil {
		return nil
	}
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
