package main

import (
	"unsafe"

	"github.com/KeyWorksRW/kwxGO/wx"
)

// createUserCreatedPage is CreateUserCreatedPage(): a panel with one button.
func createUserCreatedPage(parent unsafe.Pointer) unsafe.Pointer {
	panel := wx.NewPanel(parent, wx.ID_ANY, -1, -1, -1, -1, 0)
	if panel == nil {
		return nil
	}

	panel.SetHelpText("Panel with a Button")

	wx.NewButton(panel.Ptr(), wx.ID_ANY, "Button", 10, 10, -1, -1, 0)

	return panel.Ptr()
}

// createRadioButtonsPage is CreateRadioButtonsPage(): two radio boxes in a
// vertical sizer, with proportions 2 and 1.
func createRadioButtonsPage(parent unsafe.Pointer) unsafe.Pointer {
	panel := wx.NewPanel(parent, wx.ID_ANY, -1, -1, -1, -1, 0)
	if panel == nil {
		return nil
	}

	panel.SetHelpText("Panel with some Radio Buttons")

	animals := []string{"Fox", "Hare", "Rabbit", "Sabre-toothed tiger", "T Rex"}
	radiobox1 := newRadioBox(panel.Ptr(), wx.ID_ANY, "Choose one", animals, 2, wx.RA_SPECIFY_ROWS)

	computers := []string{"Amiga", "Commodore 64", "PET", "Another"}
	radiobox2 := newRadioBox(panel.Ptr(), wx.ID_ANY, "Choose your favourite", computers, 0, wx.RA_SPECIFY_COLS)

	sizerPanel := wx.NewBoxSizer(wx.VERTICAL)
	if sizerPanel != nil {
		if radiobox1 != nil {
			wxf := wx.NewSizerFlags(2)
			wxf.Expand()
			sizerPanel.AddWindowFlags(radiobox1.Ptr(), wxf.Ptr())
			wxf.Delete()
		}
		if radiobox2 != nil {
			wxf := wx.NewSizerFlags(1)
			wxf.Expand()
			sizerPanel.AddWindowFlags(radiobox2.Ptr(), wxf.Ptr())
			wxf.Delete()
		}
		panel.SetSizer(sizerPanel.Ptr(), true)
	}

	return panel.Ptr()
}

// createVetoPage is CreateVetoPage(): an otherwise empty panel with a label.
func createVetoPage(parent unsafe.Pointer) unsafe.Pointer {
	panel := wx.NewPanel(parent, wx.ID_ANY, -1, -1, -1, -1, 0)
	if panel == nil {
		return nil
	}

	panel.SetHelpText("An empty panel")

	s := wx.NewBoxSizer(wx.VERTICAL)
	if s != nil {
		if st := wx.NewStaticText(panel.Ptr(), wx.ID_ANY,
			"This page intentionally left blank", -1, -1, -1, -1, 0); st != nil {
			wxf := wx.NewSizerFlags(0)
			wxf.Expand()
			wxf.DoubleBorder(wx.ALL)
			s.AddWindowFlags(st.Ptr(), wxf.Ptr())
			wxf.Delete()
		}
		panel.SetSizer(s.Ptr(), true)
	}

	return panel.Ptr()
}

// createFullPageText is CreateFullPageText(): the page is the text control
// itself, not a panel.
func createFullPageText(parent unsafe.Pointer) unsafe.Pointer {
	text := wx.NewTextCtrl(parent, wx.ID_ANY, "Full page text",
		-1, -1, -1, -1, wx.TE_MULTILINE)
	if text == nil {
		return nil
	}

	text.SetHelpText("Page consisting of just a text control")

	return text.Ptr()
}

// createInsertPage is CreateInsertPage(): a maroon panel with a label.
func createInsertPage(parent unsafe.Pointer) unsafe.Pointer {
	panel := wx.NewPanel(parent, wx.ID_ANY, -1, -1, -1, -1, 0)
	if panel == nil {
		return nil
	}

	panel.SetHelpText("Maroon panel")

	if colour := wx.NewColourByName("MAROON"); colour != nil {
		panel.SetBackgroundColour(colour.Ptr())
		// SetBackgroundColour copies the colour; release the temporary.
		colour.Delete()
	}

	s := wx.NewBoxSizer(wx.VERTICAL)
	if s != nil {
		if st := wx.NewStaticText(panel.Ptr(), wx.ID_ANY,
			"This page has been inserted, not added.", -1, -1, -1, -1, 0); st != nil {
			wxf := wx.NewSizerFlags(0)
			wxf.Expand()
			wxf.DoubleBorder(wx.ALL)
			s.AddWindowFlags(st.Ptr(), wxf.Ptr())
			wxf.Delete()
		}
		panel.SetSizer(s.Ptr(), true)
	}

	return panel.Ptr()
}

// createNewPage is MyFrame::CreateNewPage(): a panel with two buttons in a
// horizontal sizer.
func createNewPage(parent unsafe.Pointer) unsafe.Pointer {
	panel := wx.NewPanel(parent, wx.ID_ANY, -1, -1, -1, -1, 0)
	if panel == nil {
		return nil
	}

	panel.SetHelpText("Panel with \"First\" and \"Second\" buttons")

	sizer := wx.NewBoxSizer(wx.HORIZONTAL)
	if sizer == nil {
		return panel.Ptr()
	}
	panel.SetSizer(sizer.Ptr(), true)

	if btn := wx.NewButton(panel.Ptr(), wx.ID_ANY, "First button", -1, -1, -1, -1, 0); btn != nil {
		wxf := wx.NewSizerFlags(0)
		wxf.Border(wx.ALL)
		sizer.AddWindowFlags(btn.Ptr(), wxf.Ptr())
		wxf.Delete()
	}
	if btn := wx.NewButton(panel.Ptr(), wx.ID_ANY, "Second button", -1, -1, -1, -1, 0); btn != nil {
		wxf := wx.NewSizerFlags(0)
		wxf.Border(wx.ALL)
		sizer.AddWindowFlags(btn.Ptr(), wxf.Ptr())
		wxf.Delete()
	}

	return panel.Ptr()
}

// newRadioBox builds a wxRadioBox from a Go string slice. The constructor takes
// the choices as a Go []string and packs them for C itself.
func newRadioBox(parent unsafe.Pointer, id int, label string, choices []string,
	majorDim int, style int) *wx.RadioBox {
	if len(choices) == 0 {
		return nil
	}

	return wx.NewRadioBox(parent, id, label, -1, -1, -1, -1, choices, majorDim, style)
}
