package main

import (
	"unsafe"

	"github.com/KeyWorksRW/kwxGO/wx"
)

// showCompareDialog is MyFrame::OnShowComparison(): it shows some
// wxOwnerDrawnComboBoxes side by side with native wxComboBoxes.
func showCompareDialog(parent unsafe.Pointer, items []string) {
	dlg := wx.NewDialog(parent, wx.ID_ANY,
		"Compare against wxComboBox",
		-1, -1, -1, -1,
		wx.DEFAULT_DIALOG_STYLE|wx.RESIZE_BORDER)
	if dlg == nil {
		return
	}

	colSizer := wx.NewBoxSizer(wx.VERTICAL)
	rowSizer := wx.NewBoxSizer(wx.HORIZONTAL)
	if colSizer == nil || rowSizer == nil {
		dlg.Destroy()
		return
	}

	// -- wxOwnerDrawnComboBox column -----------------------------------------
	if group := newOwnerDrawnGroup(dlg.Ptr(), items); group != nil {
		addSiz(rowSizer, group.Ptr(), 1, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})
	}

	// -- native wxComboBox column --------------------------------------------
	if group := newComboBoxGroup(dlg.Ptr(), items); group != nil {
		addSiz(rowSizer, group.Ptr(), 1, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})
	}

	addSiz(colSizer, rowSizer.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})

	dlg.SetSizer(colSizer.Ptr(), true)
	colSizer.SetSizeHints(dlg.Ptr())

	dlg.SetSize(-1, -1, 60, 240, wx.SIZE_USE_EXISTING)
	dlg.Centre(wx.BOTH)
	dlg.ShowModal()
	dlg.Destroy()
}

// newOwnerDrawnGroup builds the " wxOwnerDrawnComboBox " static box with the
// four owner-drawn combo variants.
func newOwnerDrawnGroup(parent unsafe.Pointer, items []string) *wx.StaticBoxSizer {
	box := wx.NewStaticBox(parent, wx.ID_ANY, " wxOwnerDrawnComboBox ",
		-1, -1, -1, -1, 0)
	if box == nil {
		return nil
	}
	group := wx.NewStaticBoxSizer(box.Ptr(), wx.VERTICAL)
	if group == nil {
		return nil
	}
	boxPtr := group.GetStaticBox()

	// Writable, with margins, sorted.
	addWin(group, label(boxPtr, "Writable, with margins, sorted:"), 0,
		func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.RIGHT)
		})

	odc := newOwnerDrawnComboBox(boxPtr, wx.CB_SORT, items)
	if odc != nil {
		odc.Append("H - Appended Item") // test sorting in append
		odc.SetValue("Dot Dash")
		odc.SetMargins(15, 10)
		addWin(group, odc.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Border(wx.ALL)
		})
	}
	group.AddStretchSpacer(1)

	// Read-only ODComboBox, big font.
	addWin(group, label(boxPtr, "Read-only, big font:"), 0,
		func(sf *wx.SizerFlags) {
			sf.Border(wx.RIGHT)
		})

	odc = newOwnerDrawnComboBox(boxPtr, wx.CB_SORT|wx.CB_READONLY, items)
	if odc != nil {
		// odc->SetFont(odc->GetFont().Scale(1.5))
		if font := scaledFont(odc.Ptr(), 1.5); font != nil {
			odc.SetFont(font.Ptr())
			font.Delete()
		}
		odc.SetValue("Dot Dash")
		odc.SetText("Dot Dash (Testing SetText)")
		addWin(group, odc.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Border(wx.ALL)
		})
	}
	group.AddStretchSpacer(1)

	// Disabled read-only ODComboBox.
	addWin(group, label(boxPtr, "Read-only disabled:"), 0,
		func(sf *wx.SizerFlags) {
			sf.Border(wx.RIGHT)
		})

	odc = newOwnerDrawnComboBox(boxPtr, wx.CB_READONLY, items)
	if odc != nil {
		odc.SetValue("Dot Dash")
		odc.Enable(false)
		addWin(group, odc.Ptr(), 3, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})
	}

	// Disabled ODComboBox.
	addWin(group, label(boxPtr, "Disabled:"), 0, func(sf *wx.SizerFlags) {
		sf.Border(wx.RIGHT)
	})

	odc = newOwnerDrawnComboBox(boxPtr, 0, items)
	if odc != nil {
		odc.SetValue("Dot Dash")
		odc.Enable(false)
		addWin(group, odc.Ptr(), 3, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})
	}

	return group
}

// newComboBoxGroup builds the " wxComboBox " static box with the four native
// combo box variants that mirror the owner-drawn column.
func newComboBoxGroup(parent unsafe.Pointer, items []string) *wx.StaticBoxSizer {
	box := wx.NewStaticBox(parent, wx.ID_ANY, " wxComboBox ",
		-1, -1, -1, -1, 0)
	if box == nil {
		return nil
	}
	group := wx.NewStaticBoxSizer(box.Ptr(), wx.VERTICAL)
	if group == nil {
		return nil
	}
	boxPtr := group.GetStaticBox()

	// Writable, with margins, sorted.
	addWin(group, label(boxPtr, "Writable, with margins, sorted:"), 0,
		func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.RIGHT)
		})

	cb := newComboBox(boxPtr, items, wx.CB_SORT)
	if cb != nil {
		cb.Append("H - Appended Item") // test sorting in append
		cb.SetValue("Dot Dash")
		// wxComboBox inherits wxTextEntry::SetMargins(), but kwxgen exposes the
		// text-entry methods only on the TextEntry wrapper. Reinterpreting the
		// same native object as a wxTextEntry is the documented borrow path and
		// reaches exactly the same C++ member.
		if te := wx.TextEntryFromPtr(cb.Ptr()); te != nil {
			te.SetMargins(15, 10)
		}
		addWin(group, cb.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Border(wx.ALL)
		})
	}
	group.AddStretchSpacer(1)

	// Read-only wxComboBox, big font.
	addWin(group, label(boxPtr, "Read-only, big font:"), 0,
		func(sf *wx.SizerFlags) {
			sf.Border(wx.RIGHT)
		})

	cb = newComboBox(boxPtr, items, wx.CB_SORT|wx.CB_READONLY)
	if cb != nil {
		// cb->SetFont(cb->GetFont().Scale(1.5))
		if font := scaledFont(cb.Ptr(), 1.5); font != nil {
			cb.SetFont(font.Ptr())
			font.Delete()
		}
		cb.SetValue("Dot Dash")
		addWin(group, cb.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Border(wx.ALL)
		})
	}
	group.AddStretchSpacer(1)

	// Disabled read-only wxComboBox.
	addWin(group, label(boxPtr, "Read-only disabled:"), 0,
		func(sf *wx.SizerFlags) {
			sf.Border(wx.RIGHT)
		})

	cb = newComboBox(boxPtr, items, wx.CB_READONLY)
	if cb != nil {
		cb.SetValue("Dot Dash")
		cb.Enable(false)
		addWin(group, cb.Ptr(), 3, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})
	}

	// Disabled wxComboBox.
	addWin(group, label(boxPtr, "Disabled:"), 0, func(sf *wx.SizerFlags) {
		sf.Border(wx.RIGHT)
	})

	cb = newComboBox(boxPtr, items, 0)
	if cb != nil {
		cb.SetValue("Dot Dash")
		cb.Enable(false)
		addWin(group, cb.Ptr(), 3, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})
	}

	return group
}

// newComboBox builds a wxComboBox from the shared item list. The constructor
// takes the choices as a Go []string and packs them for C itself.
func newComboBox(parent unsafe.Pointer, items []string, style int) *wx.ComboBox {
	if len(items) == 0 {
		return nil
	}

	return wx.NewComboBox(parent, wx.ID_ANY, "", -1, -1, -1, -1, items, style)
}

// scaledFont is `window->GetFont().Scale(factor)`: it copies the window's font
// into a new wxFont and scales it in place. The caller owns the result.
func scaledFont(window unsafe.Pointer, factor float32) *wx.Font {
	if window == nil {
		return nil
	}
	win := wx.WindowFromPtr(window)
	if win == nil {
		return nil
	}

	font := wx.NewFontDefault()
	if font == nil {
		return nil
	}
	win.GetFont(font.Ptr())
	font.Scale(factor)
	return font
}
