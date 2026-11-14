// Test dialog of the collpane sample -- MyDialog in collpane.cpp.
package main

import (
	"fmt"
	"unsafe"

	"github.com/KeyWorksRW/kwxGO/wx"
)

// collpaneDialog holds the state MyDialog keeps in member variables.
type collpaneDialog struct {
	dlg   *wx.Dialog
	pane  *wx.CollapsiblePane
	sizer *wx.GridSizer
}

// showTestDialog runs the modal "Test dialog" and returns its modal result --
// the `MyDialog dlg(this); dlg.ShowModal();` pair of MyFrame::OnShowDialog.
func showTestDialog(parent unsafe.Pointer) int {
	dlg := wx.NewDialog(parent, wx.ID_ANY, "Test dialog",
		-1, -1, -1, -1, wx.RESIZE_BORDER|wx.DEFAULT_DIALOG_STYLE)
	if dlg == nil {
		return 0
	}
	defer dlg.Destroy()

	d := &collpaneDialog{dlg: dlg}

	sz := wx.NewBoxSizer(wx.VERTICAL)
	if sz == nil {
		return 0
	}

	intro := wx.NewStaticText(dlg.Ptr(), -1,
		"This dialog allows you to test the wxCollapsiblePane control",
		-1, -1, -1, -1, 0)
	if intro == nil {
		return 0
	}
	addWindow(sz, intro.Ptr(), 0, borderOnly)

	toggleButton := wx.NewButton(dlg.Ptr(), idPaneToggleStatus, "Change status",
		-1, -1, -1, -1, 0)
	if toggleButton == nil {
		return 0
	}
	toggleButton.Bind(wx.EVT_BUTTON, d.onToggleStatus)
	addWindow(sz, toggleButton.Ptr(), 1, growExpandBorder)

	d.pane = wx.NewCollapsiblePane(dlg.Ptr(), -1, "Click here for a surprise",
		-1, -1, -1, -1, wx.CP_DEFAULT_STYLE)
	if d.pane == nil {
		return 0
	}
	d.pane.Bind(wx.EVT_COLLAPSIBLEPANE_CHANGED, d.onPaneChanged)
	addWindow(sz, d.pane.Ptr(), 0, expandBorder)

	testCtrl := wx.NewTextCtrl(dlg.Ptr(), -1, "just a test",
		-1, -1, -1, -1, 0)
	if testCtrl == nil {
		return 0
	}
	addWindow(sz, testCtrl.Ptr(), 0, expandBorder)

	sz.AddSpacer(10)

	// new wxButton(this, wxID_OK) -- the empty label becomes the stock "OK".
	okButton := wx.NewButton(dlg.Ptr(), wx.ID_OK, "", -1, -1, -1, -1, 0)
	if okButton == nil {
		return 0
	}
	addWindow(sz, okButton.Ptr(), 0, rightBorder)

	// Test controls inside the collapsible pane: wxGridSizer(4, 1, 5, 5).
	paneWindow := wx.WindowFromPtr(d.pane.GetPane())
	if paneWindow == nil {
		return 0
	}

	d.sizer = wx.NewGridSizer(4, 1, 5, 5)
	if d.sizer == nil {
		return 0
	}

	paneText1 := wx.NewStaticText(paneWindow.Ptr(), -1, "Static text",
		-1, -1, -1, -1, 0)
	if paneText1 == nil {
		return 0
	}
	addWindow(d.sizer, paneText1.Ptr(), 0, leftOnly)

	paneText2 := wx.NewStaticText(paneWindow.Ptr(), -1, "Yet another one!",
		-1, -1, -1, -1, 0)
	if paneText2 == nil {
		return 0
	}
	addWindow(d.sizer, paneText2.Ptr(), 0, leftOnly)

	paneTextCtrl := wx.NewTextCtrl(paneWindow.Ptr(), idPaneTextCtrl, "Text control",
		-1, -1, -1, -1, 0)
	if paneTextCtrl == nil {
		return 0
	}
	addWindow(d.sizer, paneTextCtrl.Ptr(), 0, leftOnly)

	paneButton := wx.NewButton(paneWindow.Ptr(), idPaneButton, "Press to align right",
		-1, -1, -1, -1, 0)
	if paneButton == nil {
		return 0
	}
	paneButton.Bind(wx.EVT_BUTTON, d.onAlignButton)
	addWindow(d.sizer, paneButton.Ptr(), 0, leftOnly)

	paneWindow.SetSizer(d.sizer.Ptr(), true)
	d.sizer.SetSizeHints(paneWindow.Ptr())

	dlg.SetSizer(sz.Ptr(), true)
	sz.SetSizeHints(dlg.Ptr())

	return dlg.ShowModal()
}

// onToggleStatus is MyDialog::OnToggleStatus.
func (d *collpaneDialog) onToggleStatus(evt *wx.Event) {
	if evt == nil {
		return
	}
	d.pane.Collapse(!d.pane.IsCollapsed())
}

// onAlignButton is MyDialog::OnAlignButton: it changes the flags of the pane's
// text-control item to wxALIGN_RIGHT and re-lays out the dialog.
func (d *collpaneDialog) onAlignButton(evt *wx.Event) {
	if evt == nil {
		return
	}

	// m_paneSizer->GetItem(FindWindow(PANE_TEXTCTRL), true)
	item := wx.SizerItemFromPtr(d.sizer.GetItemWindow(d.dlg.FindWindowByIdSelf(idPaneTextCtrl), true))
	if item == nil {
		return
	}
	item.SetFlag(wx.ALIGN_RIGHT)

	d.dlg.Layout()
}

// onPaneChanged is MyDialog::OnPaneChanged; the C++ event table entry is on the
// dialog with wxID_ANY, so binding it on the collapsible pane itself is
// equivalent.
func (d *collpaneDialog) onPaneChanged(evt *wx.Event) {
	if evt == nil {
		return
	}
	paneEvent := wx.AsCollapsiblePaneEvent(evt)
	if paneEvent == nil {
		return
	}

	state := "expanded"
	if paneEvent.GetCollapsed() {
		state = "collapsed"
	}
	wx.LogMessage(fmt.Sprintf("The pane has just been %s by the user", state))
}

// wxSizerFlags().Left()
func leftOnly(sf *wx.SizerFlags) {
	sf.Left()
}
