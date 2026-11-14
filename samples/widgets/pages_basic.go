// Pages demonstrating the basic controls: button, checkbox, toggle button and
// radio box (button.cpp, checkbox.cpp, toggle.cpp, radiobox.cpp).
package main

import (
	"unsafe"

	"github.com/KeyWorksRW/kwxGO/wx"
)

// --- CheckBox page (checkbox.cpp) ------------------------------------------

var (
	idCbReset    = wx.ID_HIGHEST
	idCbChgLabel = wx.ID_HIGHEST + 1
	idCbCheck    = wx.ID_HIGHEST + 2
	idCbUncheck  = wx.ID_HIGHEST + 3
	idCbPart     = wx.ID_HIGHEST + 4
	idCbRight    = wx.ID_HIGHEST + 5
	idCbBox      = wx.ID_HIGHEST + 6
)

type checkboxPage struct {
	pageBase

	chkRight  *wx.CheckBox
	radioKind *wx.RadioBox
	textLabel *wx.TextCtrl
	checkbox  *wx.CheckBox
	sizerChk  *wx.BoxSizer
}

func init() {
	registerPage("CheckBox", nativeCtrls, func(book unsafe.Pointer) *pageBase {
		pg := &checkboxPage{}
		pb := newPageWithIcon(book, "icons/checkbox.svg", 32, 32)
		if pb == nil {
			return nil
		}
		pg.pageBase = *pb
		pg.content = pg.create
		pg.getWidgets = func() []unsafe.Pointer {
			if pg.checkbox == nil {
				return nil
			}
			return []unsafe.Pointer{pg.checkbox.Ptr()}
		}
		pg.recreate = func(f *widgetsFrame) { pg.createCheckbox(f) }
		return &pg.pageBase
	})
}

func (pg *checkboxPage) create(f *widgetsFrame) {
	top := wx.NewBoxSizer(wx.HORIZONTAL)
	if top == nil {
		return
	}

	// left pane: set style
	boxL, left := addStaticBox(pg.win.Ptr(), "&Set style", wx.VERTICAL)
	if left != nil {
		pg.chkRight = pg.createCheckBoxAndAddToSizer(left, "&Right aligned", idCbRight, boxL.Ptr())
		addSpacer(left, 15)

		pg.radioKind = newRadioBox(boxL.Ptr(), wx.ID_ANY, "&Kind",
			[]string{"usual &2-state checkbox", "&3rd state settable by program", "&user-settable 3rd state"},
			1, 0)
		addWin(left, pg.radioKind.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})

		reset := wx.NewButton(boxL.Ptr(), idCbReset, "&Reset", -1, -1, -1, -1, 0)
		addWin(left, reset.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.CenterHorizontal()
			sf.TripleBorder(wx.ALL)
		})
		reset.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			pg.reset()
			pg.createCheckbox(f)
		})
	}

	// middle pane: operations
	boxM, middle := addStaticBox(pg.win.Ptr(), "&Operations", wx.VERTICAL)
	if middle != nil {
		row := pg.createSizerWithTextAndButton(idCbChgLabel, "Change label", wx.ID_ANY, &pg.textLabel, boxM.Ptr())
		addSub(middle, row.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})

		check := wx.NewButton(boxM.Ptr(), idCbCheck, "&Check it", -1, -1, -1, -1, 0)
		addWin(middle, check.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})
		uncheck := wx.NewButton(boxM.Ptr(), idCbUncheck, "&Uncheck it", -1, -1, -1, -1, 0)
		addWin(middle, uncheck.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})
		part := wx.NewButton(boxM.Ptr(), idCbPart, "Put in &3rd state", -1, -1, -1, -1, 0)
		addWin(middle, part.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})

		check.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.checkbox != nil {
				pg.checkbox.SetValue(1)
			}
		})
		uncheck.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.checkbox != nil {
				pg.checkbox.SetValue(0)
			}
		})
		part.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.checkbox != nil {
				pg.checkbox.Set3StateValue(wx.CHK_UNDETERMINED)
			}
		})

		frameBindButton(pg.win, idCbChgLabel, func() {
			if pg.checkbox != nil && pg.textLabel != nil {
				pg.checkbox.SetLabel(pg.textLabel.GetValue())
			}
		})
	}

	// right pane: the checkbox itself
	pg.sizerChk = wx.NewBoxSizer(wx.HORIZONTAL)
	if pg.sizerChk == nil {
		return
	}
	pg.sizerChk.SetMinSize(150, 0)

	pg.checkbox = wx.NewCheckBox(pg.win.Ptr(), idCbBox, "&Check me!", -1, -1, -1, -1, 0)
	if pg.checkbox != nil {
		pg.checkbox.Bind(wx.EVT_CHECKBOX, func(evt *wx.Event) {
			cmd := wx.AsCommandEvent(evt)
			if cmd == nil || pg.checkbox == nil {
				return
			}
			state := "not "
			if cmd.IsChecked() {
				state = ""
			}
			f.logf("Test checkbox %schecked (value = %d)", state, pg.checkbox.Get3StateValue())
		})
	}
	addWin(pg.sizerChk, pg.checkbox.Ptr(), 1, func(sf *wx.SizerFlags) { sf.Center() })

	addSub(top, left.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.LEFT)
	})
	addSub(top, middle.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL)
	})
	addSub(top, pg.sizerChk.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.RIGHT)
	})

	pg.reset()
	pg.win.SetSizer(top.Ptr(), false)

	// style changes recreate the control
	bindStyleCheckboxes(pg.win, func() { pg.createCheckbox(f) })
}

func (pg *checkboxPage) reset() {
	if pg.chkRight != nil {
		pg.chkRight.SetValue(0)
	}
	if pg.radioKind != nil {
		pg.radioKind.SetSelection(0)
	}
}

// createCheckbox is CheckBoxWidgetsPage::CreateCheckbox().
func (pg *checkboxPage) createCheckbox(f *widgetsFrame) {
	if pg.sizerChk == nil {
		return
	}

	label := "&Check me!"
	if pg.checkbox != nil {
		label = pg.checkbox.GetLabel()
		pg.checkbox.Destroy()
		pg.checkbox = nil
	}
	clearSizer(pg.sizerChk)

	flags := f.attrs.defaultFlags
	if pg.chkRight != nil && pg.chkRight.GetValue() {
		flags |= wx.ALIGN_RIGHT
	}
	if pg.radioKind != nil {
		switch pg.radioKind.GetSelection() {
		case 1:
			flags |= wx.CHK_3STATE
		case 2:
			flags |= wx.CHK_3STATE | wx.CHK_ALLOW_3RD_STATE_FOR_USER
		default:
			flags |= wx.CHK_2STATE
		}
	}

	pg.checkbox = wx.NewCheckBox(pg.win.Ptr(), idCbBox, label, -1, -1, -1, -1, flags)
	if pg.checkbox == nil {
		return
	}
	pg.checkbox.Bind(wx.EVT_CHECKBOX, func(evt *wx.Event) {
		cmd := wx.AsCommandEvent(evt)
		if cmd == nil || pg.checkbox == nil {
			return
		}
		state := "not "
		if cmd.IsChecked() {
			state = ""
		}
		f.logf("Test checkbox %schecked (value = %d)", state, pg.checkbox.Get3StateValue())
	})

	addWin(pg.sizerChk, pg.checkbox.Ptr(), 1, func(sf *wx.SizerFlags) { sf.Center() })
	pg.sizerChk.Layout()

	pg.notifyWidgetRecreation(f, pg.checkbox.Ptr())
}

// bindStyleCheckboxes connects every checkbox and radio box on the page to the
// sample's OnCheckOrRadioBox handler.
func bindStyleCheckboxes(page *wx.ScrolledWindow, fn func()) {
	page.Bind(wx.EVT_CHECKBOX, func(evt *wx.Event) { fn() })
	page.Bind(wx.EVT_RADIOBOX, func(evt *wx.Event) { fn() })
}

// frameBindButton binds a button on the page to fn by id, walking the page's
// children through the event bubbling of the page window.
func frameBindButton(page *wx.ScrolledWindow, id int, fn func()) {
	page.BindId(wx.EVT_BUTTON, id, func(evt *wx.Event) { fn() })
}
