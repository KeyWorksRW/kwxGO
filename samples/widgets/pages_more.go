// Pages demonstrating the remaining native controls: button (button.cpp),
// static controls (static.cpp) and static bitmap (statbmp.cpp).
package main

import (
	"strconv"
	"unsafe"

	"github.com/KeyWorksRW/kwxGO/wx"
)

// --- Button page (button.cpp) ----------------------------------------------

var (
	idButtonReset       = wx.ID_HIGHEST
	idButtonChangeLabel = wx.ID_HIGHEST + 1
	idButtonChangeNote  = wx.ID_HIGHEST + 2
	idButtonChangeMarg  = wx.ID_HIGHEST + 3
	idButtonCtrl        = wx.ID_HIGHEST + 4
)

type buttonPage struct {
	pageBase

	chkBitmapOnly  *wx.CheckBox
	chkTextBitmap  *wx.CheckBox
	chkFit         *wx.CheckBox
	chkAuthNeeded  *wx.CheckBox
	chkCommandLink *wx.CheckBox
	chkMarkup      *wx.CheckBox
	chkDefault     *wx.CheckBox
	chkBitmapClass *wx.CheckBox

	chkUsePressed  *wx.CheckBox
	chkUseFocused  *wx.CheckBox
	chkUseCurrent  *wx.CheckBox
	chkUseDisabled *wx.CheckBox

	radioImagePos *wx.RadioBox
	radioHAlign   *wx.RadioBox
	radioVAlign   *wx.RadioBox

	textLabel  *wx.TextCtrl
	textNote   *wx.TextCtrl
	textMargH  *wx.TextCtrl
	textMargV  *wx.TextCtrl
	sizerNote  *wx.BoxSizer
	button     *wx.Button
	bitmapBtn  *wx.BitmapButton
	anyBtn     *wx.AnyButton
	sizerBtn   *wx.BoxSizer
	margH      int
	margV      int
	cmdLnkBtns []*wx.CommandLinkButton
}

func init() {
	registerPage("Button", nativeCtrls, func(book unsafe.Pointer) *pageBase {
		pg := &buttonPage{}
		pb := newPageWithIcon(book, "icons/button.svg", 32, 32)
		if pb == nil {
			return nil
		}
		pg.pageBase = *pb
		pg.content = pg.create
		pg.getWidgets = func() []unsafe.Pointer {
			if pg.button == nil {
				return nil
			}
			return []unsafe.Pointer{pg.button.Ptr()}
		}
		pg.recreate = func(f *widgetsFrame) { pg.createButton(f) }
		return &pg.pageBase
	})
}

// artBitmap returns an ArtProvider bitmap; the caller owns it and must Delete it.
func artBitmap(id string) *wx.Bitmap {
	return wx.ArtProvider_GetBitmap(id, wx.ART_BUTTON, -1, -1)
}

// artBundle returns an ArtProvider bitmap bundle for the wxAnyButton setters,
// which take a wxBitmapBundle in wx 3.3. The caller owns it and must Delete it.
func artBundle(id string) *wx.BitmapBundle {
	return wx.ArtProvider_GetBitmapBundle(id, wx.ART_BUTTON, -1, -1)
}

func (pg *buttonPage) create(f *widgetsFrame) {
	top := wx.NewBoxSizer(wx.HORIZONTAL)
	if top == nil {
		return
	}

	boxL, left := addStaticBox(pg.win.Ptr(), "&Set style", wx.VERTICAL)
	if left != nil {
		pg.chkBitmapOnly = pg.createCheckBoxAndAddToSizer(left, "&Bitmap only", wx.ID_ANY, boxL.Ptr())
		pg.chkTextBitmap = pg.createCheckBoxAndAddToSizer(left, "Text &and bitmap", wx.ID_ANY, boxL.Ptr())
		pg.chkFit = pg.createCheckBoxAndAddToSizer(left, "&Fit exactly", wx.ID_ANY, boxL.Ptr())
		pg.chkAuthNeeded = pg.createCheckBoxAndAddToSizer(left, "Require a&uth", wx.ID_ANY, boxL.Ptr())
		pg.chkCommandLink = pg.createCheckBoxAndAddToSizer(left, "Use command &link button", wx.ID_ANY, boxL.Ptr())
		pg.chkMarkup = pg.createCheckBoxAndAddToSizer(left, "Interpret &markup", wx.ID_ANY, boxL.Ptr())
		pg.chkDefault = pg.createCheckBoxAndAddToSizer(left, "&Default", wx.ID_ANY, boxL.Ptr())
		pg.chkBitmapClass = pg.createCheckBoxAndAddToSizer(left, "Use wxBitmapButton", wx.ID_ANY, boxL.Ptr())
		pg.chkBitmapClass.SetValue(1)

		addSpacerSizer(left, 5)

		boxU, useLabels := addStaticBox(boxL.Ptr(), "&Use the following bitmaps in addition to the normal one?", wx.VERTICAL)
		if useLabels != nil {
			pg.chkUsePressed = pg.createCheckBoxAndAddToSizer(useLabels, "&Pressed (small help icon)", wx.ID_ANY, boxU.Ptr())
			pg.chkUseFocused = pg.createCheckBoxAndAddToSizer(useLabels, "&Focused (small error icon)", wx.ID_ANY, boxU.Ptr())
			pg.chkUseCurrent = pg.createCheckBoxAndAddToSizer(useLabels, "&Current (small warning icon)", wx.ID_ANY, boxU.Ptr())
			pg.chkUseDisabled = pg.createCheckBoxAndAddToSizer(useLabels, "&Disabled (broken image icon)", wx.ID_ANY, boxU.Ptr())
			addSub(left, useLabels.Ptr(), 0, func(sf *wx.SizerFlags) {
				sf.Expand()
				sf.Border(wx.ALL)
			})
		}

		addSpacerSizer(left, 10)

		pg.radioImagePos = newRadioBox(boxL.Ptr(), wx.ID_ANY, "Image &position",
			[]string{"left", "right", "top", "bottom"}, 0, wx.RA_SPECIFY_COLS)
		addExpandBorder(left, pg.radioImagePos.Ptr())

		boxM2, margins := addStaticBox(boxL.Ptr(), "Image margins", wx.VERTICAL)
		if margins != nil {
			row := pg.createSizerWithTextAndButton(idButtonChangeMarg, "Horizontal and vertical",
				wx.ID_ANY, &pg.textMargH, boxM2.Ptr())
			addSub(margins, row.Ptr(), 0, func(sf *wx.SizerFlags) {
				sf.Border(wx.ALL)
				sf.Center()
			})
			pg.textMargV = wx.NewTextCtrl(boxM2.Ptr(), wx.ID_ANY, "", -1, -1, -1, -1, 0)
			addWin(row, pg.textMargV.Ptr(), 1, func(sf *wx.SizerFlags) {
				sf.CenterVertical()
				sf.Border(wx.LEFT)
			})
			addSub(left, margins.Ptr(), 0, func(sf *wx.SizerFlags) {
				sf.Expand()
				sf.Border(wx.ALL)
			})
		}

		addSpacerSizer(left, 15)

		pg.radioHAlign = newRadioBox(boxL.Ptr(), wx.ID_ANY, "&Horz alignment",
			[]string{"left", "centre", "right"}, 0, wx.RA_SPECIFY_COLS)
		pg.radioVAlign = newRadioBox(boxL.Ptr(), wx.ID_ANY, "&Vert alignment",
			[]string{"top", "centre", "bottom"}, 0, wx.RA_SPECIFY_COLS)
		addExpandBorder(left, pg.radioHAlign.Ptr())
		addExpandBorder(left, pg.radioVAlign.Ptr())

		addSpacerSizer(left, 5)

		reset := wx.NewButton(boxL.Ptr(), idButtonReset, "&Reset", -1, -1, -1, -1, 0)
		addWin(left, reset.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.CenterHorizontal()
			sf.TripleBorder(wx.ALL)
		})
		reset.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			pg.reset()
			pg.createButton(f)
		})
	}

	boxM, middle := addStaticBox(pg.win.Ptr(), "&Operations", wx.VERTICAL)
	if middle != nil {
		row := pg.createSizerWithTextAndButton(idButtonChangeLabel, "Change label", wx.ID_ANY, &pg.textLabel, boxM.Ptr())
		addSub(middle, row.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})
		if pg.textLabel != nil {
			pg.textLabel.SetValue("&Press me!")
		}

		pg.sizerNote = pg.createSizerWithTextAndButton(idButtonChangeNote, "Change note", wx.ID_ANY, &pg.textNote, boxM.Ptr())
		addSub(middle, pg.sizerNote.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})
		if pg.textNote != nil {
			pg.textNote.SetValue("Writes down button clicks in the log.")
		}

		pg.win.BindId(wx.EVT_BUTTON, idButtonChangeLabel, func(evt *wx.Event) { pg.changeLabel(f) })
		pg.win.BindId(wx.EVT_BUTTON, idButtonChangeNote, func(evt *wx.Event) {
			if len(pg.cmdLnkBtns) > 0 && pg.textNote != nil && pg.cmdLnkBtns[0] != nil {
				pg.cmdLnkBtns[0].SetNote(pg.textNote.GetValue())
				pg.sizerBtn.Layout()
			}
		})
		pg.win.BindId(wx.EVT_BUTTON, idButtonChangeMarg, func(evt *wx.Event) { pg.changeMargins(f) })
	}

	pg.sizerBtn = wx.NewBoxSizer(wx.VERTICAL)
	if pg.sizerBtn != nil {
		pg.sizerBtn.SetMinSize(150, 0)
	}

	addSub(top, left.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.LEFT)
	})
	addSub(top, middle.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL)
	})
	addSub(top, pg.sizerBtn.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.RIGHT)
	})

	pg.reset()
	pg.createButton(f)
	pg.win.SetSizer(top.Ptr(), false)
	pg.bindStyleChanges(f)

	if pg.textMargH != nil {
		pg.textMargH.SetValue(strconv.Itoa(pg.margH))
	}
	if pg.textMargV != nil {
		pg.textMargV.SetValue(strconv.Itoa(pg.margV))
	}
}

// bindStyleChanges connects every check/radio box to ButtonWidgetsPage's
// OnCheckOrRadioBox, which recreates the button.
func (pg *buttonPage) bindStyleChanges(f *widgetsFrame) {
	pg.win.Bind(wx.EVT_CHECKBOX, func(evt *wx.Event) {
		pg.createButton(f)
		pg.win.Layout()
	})
	pg.win.Bind(wx.EVT_RADIOBOX, func(evt *wx.Event) {
		pg.createButton(f)
		pg.win.Layout()
	})
}

func (pg *buttonPage) reset() {
	setChk(pg.chkBitmapOnly, 0)
	setChk(pg.chkFit, 0)
	setChk(pg.chkAuthNeeded, 0)
	setChk(pg.chkTextBitmap, 0)
	setChk(pg.chkDefault, 0)
	setChk(pg.chkCommandLink, 0)
	setChk(pg.chkMarkup, 0)
	setChk(pg.chkBitmapClass, 1)
	setChk(pg.chkUsePressed, 1)
	setChk(pg.chkUseFocused, 1)
	setChk(pg.chkUseCurrent, 1)
	setChk(pg.chkUseDisabled, 1)

	if pg.radioImagePos != nil {
		pg.radioImagePos.SetSelection(0)
	}
	if pg.radioHAlign != nil {
		pg.radioHAlign.SetSelection(1)
	}
	if pg.radioVAlign != nil {
		pg.radioVAlign.SetSelection(1)
	}

	pg.margH, pg.margV = 0, 0
	if pg.textMargH != nil {
		pg.textMargH.SetValue("0")
	}
	if pg.textMargV != nil {
		pg.textMargV.SetValue("0")
	}
}

func checked(c *wx.CheckBox) bool {
	return c != nil && c.GetValue()
}

func selection(r *wx.RadioBox) int {
	if r == nil {
		return 0
	}
	return r.GetSelection()
}

// createButton is ButtonWidgetsPage::CreateButton().
func (pg *buttonPage) createButton(f *widgetsFrame) {
	if pg.sizerBtn == nil {
		return
	}

	label := ""
	if pg.button != nil {
		label = wx.WindowFromPtr(pg.button.Ptr()).GetLabel()
		if len(pg.cmdLnkBtns) > 0 && pg.cmdLnkBtns[0] != nil {
			label = pg.cmdLnkBtns[0].GetMainLabel()
		}
		clearSizer(pg.sizerBtn)
		destroyWin(pg.button)
		pg.cmdLnkBtns = nil
		// the old button is destroyed above; drop our reference so it is never
		// used again.
		pg.button = nil
	}
	if label == "" && pg.textLabel != nil {
		label = pg.textLabel.GetValue()
	}

	flags := f.attrs.defaultFlags
	switch selection(pg.radioHAlign) {
	case 0:
		flags |= wx.BU_LEFT
	case 2:
		flags |= wx.BU_RIGHT
	}
	switch selection(pg.radioVAlign) {
	case 0:
		flags |= wx.BU_TOP
	case 2:
		flags |= wx.BU_BOTTOM
	}
	if checked(pg.chkFit) {
		flags |= wx.BU_EXACTFIT
	}

	showsBitmap := false
	pg.bitmapBtn, pg.anyBtn = nil, nil
	if checked(pg.chkBitmapOnly) {
		showsBitmap = true
		if checked(pg.chkBitmapClass) {
			bmp := artBitmap(wx.ART_INFORMATION)
			if bmp != nil {
				defer bmp.Delete()
				if bb := wx.NewBitmapButton(pg.win.Ptr(), idButtonCtrl, bmp.Ptr(), -1, -1, -1, -1, flags); bb != nil {
					pg.bitmapBtn = bb
					pg.button = wx.ButtonFromPtr(bb.Ptr())
				}
			}
		} else {
			pg.button = wx.NewButton(pg.win.Ptr(), idButtonCtrl, "", -1, -1, -1, -1, flags)
			if pg.button != nil {
				bmp := artBundle(wx.ART_INFORMATION)
				if bmp != nil {
					defer bmp.Delete()
					if ab := wx.AnyButtonFromPtr(pg.button.Ptr()); ab != nil {
						ab.SetBitmapLabel(bmp.Ptr())
					}
				}
			}
		}
		if pg.button != nil {
			pg.anyBtn = wx.AnyButtonFromPtr(pg.button.Ptr())
			pg.anyBtn.SetBitmapMargins(pg.margH, pg.margV)
			pg.setExtraBitmaps(pg.anyBtn)
		}
	} else if checked(pg.chkCommandLink) {
		note := ""
		if pg.textNote != nil {
			note = pg.textNote.GetValue()
		}
		clb := wx.NewCommandLinkButton(pg.win.Ptr(), idButtonCtrl, label, note, -1, -1, -1, -1, flags)
		if clb != nil {
			pg.cmdLnkBtns = append(pg.cmdLnkBtns, clb)
			pg.button = wx.ButtonFromPtr(clb.Ptr())
			pg.anyBtn = wx.AnyButtonFromPtr(clb.Ptr())
		}
	} else {
		pg.button = wx.NewButton(pg.win.Ptr(), idButtonCtrl, label, -1, -1, -1, -1, flags)
		if pg.button != nil {
			pg.anyBtn = wx.AnyButtonFromPtr(pg.button.Ptr())
		}
	}

	if pg.sizerNote != nil {
		pg.sizerNote.ShowItems(checked(pg.chkCommandLink))
	}

	if pg.button == nil {
		return
	}
	pg.notifyWidgetRecreation(f, pg.button.Ptr())

	if !showsBitmap && checked(pg.chkTextBitmap) {
		showsBitmap = true
		dirs := []int{wx.LEFT, wx.RIGHT, wx.TOP, wx.BOTTOM}
		if bmp := artBundle(wx.ART_INFORMATION); bmp != nil {
			if pg.anyBtn == nil {
				pg.anyBtn = wx.AnyButtonFromPtr(pg.button.Ptr())
			}
			pg.anyBtn.SetBitmap(bmp.Ptr(), dirs[selection(pg.radioImagePos)])
			bmp.Delete()
		}
		if pg.anyBtn != nil {
			pg.anyBtn.SetBitmapMargins(pg.margH, pg.margV)
			pg.setExtraBitmaps(pg.anyBtn)
		}
	}

	// enable/disable the controls that only make sense with a bitmap
	enableWin(pg.chkTextBitmap, !checked(pg.chkBitmapOnly))
	enableWin(pg.chkBitmapOnly, !checked(pg.chkTextBitmap))
	enableWin(pg.chkCommandLink, !checked(pg.chkBitmapOnly))
	enableWin(pg.chkBitmapClass, showsBitmap)
	enableWin(pg.chkUsePressed, showsBitmap)
	enableWin(pg.chkUseFocused, showsBitmap)
	enableWin(pg.chkUseCurrent, showsBitmap)
	enableWin(pg.chkUseDisabled, showsBitmap)
	enableWin(pg.radioImagePos, checked(pg.chkTextBitmap))
	enableWin(pg.textMargH, showsBitmap)
	enableWin(pg.textMargV, showsBitmap)
	if w := wx.WindowFromPtr(wx.Window_FindWindowById(idButtonChangeMarg, pg.win.Ptr())); w != nil {
		w.Enable(showsBitmap)
	}

	if checked(pg.chkAuthNeeded) {
		pg.button.SetAuthNeeded(true)
	}
	if checked(pg.chkDefault) {
		pg.button.SetDefault()
	}

	pg.sizerBtn.AddStretchSpacer(0)
	addWin(pg.sizerBtn, pg.button.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Center()
		sf.Border(wx.ALL)
	})
	pg.sizerBtn.AddStretchSpacer(2)
	pg.sizerBtn.Layout()
}

// setExtraBitmaps applies the pressed/focused/current/disabled bitmaps.
func (pg *buttonPage) setExtraBitmaps(btn *wx.AnyButton) {
	if btn == nil {
		return
	}
	if checked(pg.chkUsePressed) {
		if bmp := artBundle(wx.ART_HELP); bmp != nil {
			btn.SetBitmapPressed(bmp.Ptr())
			bmp.Delete()
		}
	}
	if checked(pg.chkUseFocused) {
		if bmp := artBundle(wx.ART_ERROR); bmp != nil {
			btn.SetBitmapFocus(bmp.Ptr())
			bmp.Delete()
		}
	}
	if checked(pg.chkUseCurrent) {
		if bmp := artBundle(wx.ART_WARNING); bmp != nil {
			btn.SetBitmapCurrent(bmp.Ptr())
			bmp.Delete()
		}
	}
	if checked(pg.chkUseDisabled) {
		if bmp := artBundle(wx.ART_MISSING_IMAGE); bmp != nil {
			btn.SetBitmapDisabled(bmp.Ptr())
			bmp.Delete()
		}
	}
}

func (pg *buttonPage) changeLabel(f *widgetsFrame) {
	if pg.textLabel == nil {
		return
	}
	label := pg.textLabel.GetValue()
	if len(pg.cmdLnkBtns) > 0 && pg.cmdLnkBtns[0] != nil {
		pg.cmdLnkBtns[0].SetMainLabel(label)
	} else if pg.button != nil {
		if checked(pg.chkMarkup) {
			pg.button.SetLabelMarkup(label)
		} else {
			pg.button.SetLabel(label)
		}
	}
	if checked(pg.chkBitmapOnly) {
		pg.createButton(f)
	}
	if pg.sizerBtn != nil {
		pg.sizerBtn.Layout()
	}
}

func (pg *buttonPage) changeMargins(f *widgetsFrame) {
	if pg.textMargH == nil || pg.textMargV == nil {
		return
	}
	mh, errH := strconv.Atoi(pg.textMargH.GetValue())
	mv, errV := strconv.Atoi(pg.textMargV.GetValue())
	if errH != nil || errV != nil || mh < 0 || mv < 0 {
		f.logf("Invalid margin values for bitmap.")
		return
	}
	pg.margH, pg.margV = mh, mv
	if pg.anyBtn != nil {
		pg.anyBtn.SetBitmapMargins(pg.margH, pg.margV)
	}
	if pg.button != nil {
		wx.WindowFromPtr(pg.button.Ptr()).Refresh(true)
	}
	if pg.sizerBtn != nil {
		pg.sizerBtn.Layout()
	}
}

// enableWin enables or disables a window given by wrapper.
func enableWin(w interface {
	Ptr() unsafe.Pointer
}, enable bool) {
	if ptr := wrapperPtr(w); ptr != nil {
		if win := wx.WindowFromPtr(ptr); win != nil {
			win.Enable(enable)
		}
	}
}

// --- Static page (static.cpp) ----------------------------------------------

type staticPage struct {
	pageBase

	chkGeneric   *wx.CheckBox
	chkVert      *wx.CheckBox
	chkAutoSize  *wx.CheckBox
	chkWrap      *wx.CheckBox
	chkEllipsize *wx.CheckBox

	chkMarkup *wx.CheckBox
	chkGreen  *wx.CheckBox

	radioHAlign    *wx.RadioBox
	radioVAlign    *wx.RadioBox
	radioEllipsize *wx.RadioBox

	box         *wx.StaticBox
	sizerBox    *wx.StaticBoxSizer
	statText    *wx.StaticText
	statMarkup  *wx.StaticText
	statLine    *wx.StaticLine
	sizerStatic *wx.BoxSizer

	textBox         *wx.TextCtrl
	textLabel       *wx.TextCtrl
	textLabelMarkup *wx.TextCtrl
}

func init() {
	registerPage("Static", nativeCtrls, func(book unsafe.Pointer) *pageBase {
		pg := &staticPage{}
		pb := newPageWithIcon(book, "icons/statbox.svg", 32, 32)
		if pb == nil {
			return nil
		}
		pg.pageBase = *pb
		pg.content = pg.create
		pg.getWidgets = func() []unsafe.Pointer {
			var widgets []unsafe.Pointer
			if pg.box != nil {
				widgets = append(widgets, pg.box.Ptr())
			}
			if pg.statText != nil {
				widgets = append(widgets, pg.statText.Ptr())
			}
			if pg.statMarkup != nil {
				widgets = append(widgets, pg.statMarkup.Ptr())
			}
			if pg.statLine != nil {
				widgets = append(widgets, pg.statLine.Ptr())
			}
			return widgets
		}
		pg.recreate = func(f *widgetsFrame) { pg.createStatic(f) }
		return &pg.pageBase
	})
}

func (pg *staticPage) create(f *widgetsFrame) {
	top := wx.NewBoxSizer(wx.HORIZONTAL)
	if top == nil {
		return
	}

	boxL, left := addStaticBox(pg.win.Ptr(), "&Set style", wx.VERTICAL)
	if left != nil {
		pg.chkGeneric = pg.createCheckBoxAndAddToSizer(left, "&Generic wxStaticText", wx.ID_ANY, boxL.Ptr())
		pg.chkVert = pg.createCheckBoxAndAddToSizer(left, "&Vertical line", wx.ID_ANY, boxL.Ptr())
		pg.chkAutoSize = pg.createCheckBoxAndAddToSizer(left, "&Fit to text", wx.ID_ANY, boxL.Ptr())
		pg.chkWrap = pg.createCheckBoxAndAddToSizer(left, "&Wrap", wx.ID_ANY, boxL.Ptr())

		addSpacerSizer(left, 15)

		pg.radioHAlign = newRadioBox(boxL.Ptr(), wx.ID_ANY, "&Horz alignment",
			[]string{"left", "centre", "right"}, 3, wx.RA_SPECIFY_COLS)
		pg.radioVAlign = newRadioBox(boxL.Ptr(), wx.ID_ANY, "&Vert alignment",
			[]string{"top", "centre", "bottom"}, 3, wx.RA_SPECIFY_COLS)
		if pg.radioVAlign != nil {
			wx.WindowFromPtr(pg.radioVAlign.Ptr()).SetToolTip("Relevant for Generic wxStaticText only")
		}
		addExpandBorder(left, pg.radioHAlign.Ptr())
		addExpandBorder(left, pg.radioVAlign.Ptr())

		addSpacerSizer(left, 15)

		pg.chkEllipsize = pg.createCheckBoxAndAddToSizer(left, "&Ellipsize", wx.ID_ANY, boxL.Ptr())
		pg.radioEllipsize = newRadioBox(boxL.Ptr(), wx.ID_ANY, "&Ellipsize mode",
			[]string{"&start", "&middle", "&end"}, 3, wx.RA_SPECIFY_COLS)
		addExpandBorder(left, pg.radioEllipsize.Ptr())
		pg.chkEllipsize.Bind(wx.EVT_CHECKBOX, func(evt *wx.Event) {
			cmd := wx.AsCommandEvent(evt)
			if pg.radioEllipsize != nil && cmd != nil {
				wx.WindowFromPtr(pg.radioEllipsize.Ptr()).Enable(cmd.IsChecked())
			}
			pg.createStatic(f)
		})
		pg.win.Bind(wx.EVT_RADIOBOX, func(evt *wx.Event) { pg.createStatic(f) })
		pg.win.Bind(wx.EVT_CHECKBOX, func(evt *wx.Event) { pg.createStatic(f) })

		reset := wx.NewButton(boxL.Ptr(), wx.ID_ANY, "&Reset", -1, -1, -1, -1, 0)
		addWin(left, reset.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.CenterHorizontal()
			sf.TripleBorder(wx.ALL)
		})
		reset.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			pg.reset()
			pg.createStatic(f)
		})
	}

	boxM, middle := addStaticBox(pg.win.Ptr(), "&Change labels", wx.VERTICAL)
	if middle != nil {
		pg.textBox = wx.NewTextCtrl(boxM.Ptr(), wx.ID_ANY, "", -1, -1, -1, -1, 0)
		addWin(middle, pg.textBox.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})
		b1 := wx.NewButton(boxM.Ptr(), wx.ID_ANY, "Change &box label", -1, -1, -1, -1, 0)
		addWin(middle, b1.Ptr(), 0, func(sf *wx.SizerFlags) { sf.Border(wx.LEFT | wx.BOTTOM) })
		b1.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.sizerBox != nil && pg.textBox != nil {
				if st := wx.StaticBoxFromPtr(pg.sizerBox.GetStaticBox()); st != nil {
					wx.WindowFromPtr(st.Ptr()).SetLabel(pg.textBox.GetValue())
				}
			}
		})

		pg.textLabel = wx.NewTextCtrl(boxM.Ptr(), wx.ID_ANY, "", -1, -1, -1, -1, wx.TE_MULTILINE|wx.HSCROLL)
		addWin(middle, pg.textLabel.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})
		b2 := wx.NewButton(boxM.Ptr(), wx.ID_ANY, "Change &text label", -1, -1, -1, -1, 0)
		addWin(middle, b2.Ptr(), 0, func(sf *wx.SizerFlags) { sf.Border(wx.LEFT | wx.BOTTOM) })
		b2.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.statText != nil && pg.textLabel != nil {
				pg.statText.SetLabel(pg.textLabel.GetValue())
				f.logf("The original label should be '%s'", pg.statText.GetLabel())
				f.logf("The label text is '%s'", pg.statText.GetLabelText())
			}
		})

		pg.textLabelMarkup = wx.NewTextCtrl(boxM.Ptr(), wx.ID_ANY, "", -1, -1, -1, -1, wx.TE_MULTILINE|wx.HSCROLL)
		addWin(middle, pg.textLabelMarkup.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})
		b3 := wx.NewButton(boxM.Ptr(), wx.ID_ANY, "Change decorated text label", -1, -1, -1, -1, 0)
		addWin(middle, b3.Ptr(), 0, func(sf *wx.SizerFlags) { sf.Border(wx.LEFT | wx.BOTTOM) })
		b3.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.statMarkup != nil && pg.textLabelMarkup != nil {
				pg.statMarkup.SetLabelMarkup(pg.textLabelMarkup.GetValue())
				f.logf("The original label should be '%s'", pg.statMarkup.GetLabel())
				f.logf("The label text is '%s'", pg.statMarkup.GetLabelText())
			}
		})
		pg.chkGreen = pg.createCheckBoxAndAddToSizer(middle, "Decorated label on g&reen", wx.ID_ANY, boxM.Ptr())
	}

	pg.reset()

	if pg.textBox != nil {
		pg.textBox.SetValue("This is a &box")
	}
	if pg.textLabel != nil {
		pg.textLabel.SetValue("And this is a\n\tlabel inside the box with a &mnemonic.\n" +
			"Only this text is affected by the ellipsize settings.")
	}
	if pg.textLabelMarkup != nil {
		pg.textLabelMarkup.SetValue("Another label, this time <b>decorated</b>\n" +
			"with <u>markup</u>; here you need entities\n" +
			"for the symbols: &lt; &gt; &amp;&amp; &apos; &quot;\n" +
			"but you can still use &mnemonics too")
	}

	pg.sizerStatic = wx.NewBoxSizer(wx.HORIZONTAL)
	if pg.sizerStatic == nil {
		return
	}
	pg.sizerStatic.SetMinSize(150, 0)

	pg.createStatic(f)

	addSub(top, left.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.LEFT)
	})
	addSub(top, middle.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL)
	})
	addSub(top, pg.sizerStatic.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.RIGHT)
	})
	pg.win.SetSizer(top.Ptr(), false)
}

func (pg *staticPage) reset() {
	setChk(pg.chkGeneric, 0)
	setChk(pg.chkVert, 0)
	setChk(pg.chkAutoSize, 1)
	setChk(pg.chkWrap, 0)
	setChk(pg.chkEllipsize, 0)
	setChk(pg.chkGreen, 0)
	if pg.radioEllipsize != nil {
		wx.WindowFromPtr(pg.radioEllipsize.Ptr()).Enable(false)
	}
	if pg.radioHAlign != nil {
		pg.radioHAlign.SetSelection(0)
	}
	if pg.radioVAlign != nil {
		pg.radioVAlign.SetSelection(0)
	}
	if pg.radioEllipsize != nil {
		pg.radioEllipsize.SetSelection(0)
	}
}

// createStatic is StaticWidgetsPage::CreateStatic().
func (pg *staticPage) createStatic(f *widgetsFrame) {
	if pg.sizerStatic == nil {
		return
	}

	isVert := checked(pg.chkVert)

	if pg.sizerBox != nil {
		pg.sizerStatic.Remove(0)
		// wxSizer::Remove destroys the sizer item, which destroys the child
		// sizer and the static box (and therefore its children) as well.
		pg.sizerBox, pg.box = nil, nil
		pg.statText, pg.statMarkup, pg.statLine = nil, nil, nil
	}

	flagsBox := 0
	flagsText := f.attrs.defaultFlags
	flagsMarkup := f.attrs.defaultFlags

	if !checked(pg.chkAutoSize) {
		flagsText |= wx.ST_NO_AUTORESIZE
		flagsMarkup |= wx.ST_NO_AUTORESIZE
	}
	if checked(pg.chkWrap) {
		flagsText |= wx.ST_WRAP
		flagsMarkup |= wx.ST_WRAP
	}

	align := 0
	switch selection(pg.radioHAlign) {
	case 1:
		align |= wx.ALIGN_CENTRE_HORIZONTAL
	case 2:
		align |= wx.ALIGN_RIGHT
	default:
		align |= wx.ALIGN_LEFT
	}
	switch selection(pg.radioVAlign) {
	case 1:
		align |= wx.ALIGN_CENTRE_VERTICAL
	case 2:
		align |= wx.ALIGN_BOTTOM
	default:
		align |= wx.ALIGN_TOP
	}

	if checked(pg.chkEllipsize) {
		switch selection(pg.radioEllipsize) {
		case 1:
			flagsText |= wx.ST_ELLIPSIZE_MIDDLE
		case 2:
			flagsText |= wx.ST_ELLIPSIZE_END
		default:
			flagsText |= wx.ST_ELLIPSIZE_START
		}
	}

	flagsText |= align
	flagsMarkup |= align
	flagsBox |= align

	boxLabel := ""
	if pg.textBox != nil {
		boxLabel = pg.textBox.GetValue()
	}
	box := wx.NewStaticBox(pg.win.Ptr(), wx.ID_ANY, boxLabel, -1, -1, -1, -1, flagsBox)
	if box == nil {
		return
	}
	orient := wx.VERTICAL
	if isVert {
		orient = wx.HORIZONTAL
	}
	pg.box = box
	pg.sizerBox = wx.NewStaticBoxSizer(box.Ptr(), orient)
	if pg.sizerBox == nil {
		return
	}

	label := ""
	if pg.textLabel != nil {
		label = pg.textLabel.GetValue()
	}
	if checked(pg.chkGeneric) {
		if st := wx.NewGenericStaticText(box.Ptr(), wx.ID_ANY, label, -1, -1, -1, -1, flagsText); st != nil {
			pg.statText = wx.StaticTextFromPtr(st.Ptr())
		}
		if st := wx.NewGenericStaticText(box.Ptr(), wx.ID_ANY, "", -1, -1, -1, -1, flagsMarkup); st != nil {
			pg.statMarkup = wx.StaticTextFromPtr(st.Ptr())
		}
	} else {
		pg.statText = wx.NewStaticText(box.Ptr(), wx.ID_ANY, label, -1, -1, -1, -1, flagsText)
		pg.statMarkup = wx.NewStaticText(box.Ptr(), wx.ID_ANY, "", -1, -1, -1, -1, flagsMarkup)
	}

	if pg.statText != nil {
		wx.WindowFromPtr(pg.statText.Ptr()).SetToolTip("Tooltip for a label inside the box")
	}
	if pg.statMarkup != nil {
		if pg.textLabelMarkup != nil {
			pg.statMarkup.SetLabelMarkup(pg.textLabelMarkup.GetValue())
		}
		if checked(pg.chkGreen) {
			if green := wx.NewColourRGB(0, 255, 0, 255); green != nil {
				wx.WindowFromPtr(pg.statMarkup.Ptr()).SetBackgroundColour(green.Ptr())
			}
		}
	}

	lineFlags := wx.LI_HORIZONTAL
	if isVert {
		lineFlags = wx.LI_VERTICAL
	}
	pg.statLine = wx.NewStaticLine(box.Ptr(), wx.ID_ANY, -1, -1, -1, -1, lineFlags)

	addWin(pg.sizerBox, pg.statText.Ptr(), 0, func(sf *wx.SizerFlags) { sf.Expand() })
	if pg.statLine != nil {
		addWin(pg.sizerBox, pg.statLine.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.DoubleBorder(wx.TOP | wx.BOTTOM)
		})
	}
	if pg.statMarkup != nil {
		addWin(pg.sizerBox, pg.statMarkup.Ptr(), 0, nil)
	}

	pg.notifyWidgetRecreation(f, pg.statText.Ptr())
	if pg.statMarkup != nil {
		pg.notifyWidgetRecreation(f, pg.statMarkup.Ptr())
	}
	if pg.statLine != nil {
		pg.notifyWidgetRecreation(f, pg.statLine.Ptr())
	}

	addSub(pg.sizerStatic, pg.sizerBox.Ptr(), 1, func(sf *wx.SizerFlags) { sf.Expand() })
	pg.sizerStatic.Layout()

	pg.statText.Bind(wx.EVT_LEFT_UP, func(evt *wx.Event) { f.logf("Clicked on static text") })
	wx.WindowFromPtr(box.Ptr()).Bind(wx.EVT_LEFT_UP, func(evt *wx.Event) { f.logf("Clicked on static box") })

	f.setUpWidget(&pg.pageBase)
}

// --- StaticBitmap page (statbmp.cpp) ---------------------------------------

type statBmpPage struct {
	pageBase

	radioImpl  *wx.RadioBox
	radioScale *wx.RadioBox
	filePicker *wx.FilePickerCtrl
	box        *wx.StaticBox
	sizerBox   *wx.StaticBoxSizer
	statbmp    *wx.StaticBitmap
	genericBmp *wx.GenericStaticBitmap
}

func init() {
	registerPage("StaticBitmap", allCtrls, func(book unsafe.Pointer) *pageBase {
		pg := &statBmpPage{}
		pb := newPageWithIcon(book, "icons/statbmp.svg", 32, 32)
		if pb == nil {
			return nil
		}
		pg.pageBase = *pb
		pg.content = pg.create
		pg.getWidgets = func() []unsafe.Pointer {
			if pg.statbmp != nil {
				return []unsafe.Pointer{pg.statbmp.Ptr()}
			}
			if pg.genericBmp != nil {
				return []unsafe.Pointer{pg.genericBmp.Ptr()}
			}
			return nil
		}
		pg.recreate = func(f *widgetsFrame) { pg.recreateBitmap(f) }
		return &pg.pageBase
	})
}

func (pg *statBmpPage) create(f *widgetsFrame) {
	pg.radioImpl = newRadioBox(pg.win.Ptr(), wx.ID_ANY, "Implementation",
		[]string{"native", "generic"}, 0, wx.RA_SPECIFY_COLS)
	pg.radioScale = newRadioBox(pg.win.Ptr(), wx.ID_ANY, "Scale Mode",
		[]string{"None", "Fill", "Aspect Fit", "Aspect Fill"}, 0, wx.RA_SPECIFY_COLS)

	pg.filePicker = wx.NewFilePickerCtrl(pg.win.Ptr(), wx.ID_ANY, -1, -1, -1, -1,
		"", "Choose an image file", "All files (*.*)|*.*", 0)

	pg.box, pg.sizerBox = addStaticBox(pg.win.Ptr(), "wxStaticBitmap inside", wx.VERTICAL)

	left := wx.NewBoxSizer(wx.VERTICAL)
	if left != nil {
		addExpandBorder(left, pg.radioImpl.Ptr())
		addExpandBorder(left, pg.radioScale.Ptr())
		addExpandBorder(left, pg.filePicker.Ptr())
	}

	sz := wx.NewBoxSizer(wx.HORIZONTAL)
	if sz == nil {
		return
	}
	addSub(sz, left.Ptr(), 0, func(sf *wx.SizerFlags) { sf.Border(wx.ALL) })
	if pg.sizerBox != nil {
		addSub(sz, pg.sizerBox.Ptr(), 0, func(sf *wx.SizerFlags) { sf.Center() })
	}
	pg.win.SetSizer(sz.Ptr(), false)

	pg.win.Bind(wx.EVT_FILEPICKER_CHANGED, func(evt *wx.Event) { pg.recreateBitmap(f) })
	pg.win.Bind(wx.EVT_RADIOBOX, func(evt *wx.Event) { pg.recreateBitmap(f) })

	pg.recreateBitmap(f)
}

// recreateBitmap is StatBmpWidgetsPage::RecreateWidget().
func (pg *statBmpPage) recreateBitmap(f *widgetsFrame) {
	if pg.sizerBox == nil {
		return
	}

	destroyWidgets(ptrOf(pg.statbmp), ptrOf(pg.genericBmp))
	pg.statbmp = nil
	pg.genericBmp = nil
	clearSizer(pg.sizerBox)

	path := ""
	if pg.filePicker != nil {
		path = pg.filePicker.GetPath()
	}

	var bmp *wx.Bitmap
	if path != "" {
		if img := wx.NewImageFromFile(path); img != nil && img.IsOk() {
			bmp = wx.NewBitmapFromImage(img.Ptr(), -1)
		} else {
			f.logf("Reading image from file '%s' failed.", path)
		}
	}
	if bmp == nil || !bmp.IsOk() {
		bmp = artBitmap(wx.ART_MISSING_IMAGE)
	}
	if bmp == nil {
		return
	}
	defer bmp.Delete()

	style := f.attrs.defaultFlags
	if selection(pg.radioImpl) == 0 {
		pg.statbmp = wx.NewStaticBitmap(pg.box.Ptr(), wx.ID_ANY, bmp.Ptr(), -1, -1, -1, -1, style)
	} else {
		pg.genericBmp = wx.NewGenericStaticBitmap(pg.box.Ptr(), wx.ID_ANY, bmp.Ptr(), -1, -1, -1, -1, style)
	}

	scale := selection(pg.radioScale)
	var handle unsafe.Pointer
	if pg.statbmp != nil {
		handle = pg.statbmp.Ptr()
		pg.statbmp.SetScaleMode(scale)
		if pg.statbmp.GetScaleMode() != scale {
			f.logf("Scale mode not supported by current implementation")
		}
	} else if pg.genericBmp != nil {
		handle = pg.genericBmp.Ptr()
		pg.genericBmp.SetScaleMode(scale)
		if pg.genericBmp.GetScaleMode() != scale {
			f.logf("Scale mode not supported by current implementation")
		}
	} else {
		return
	}

	pg.notifyWidgetRecreation(f, handle)
	addWin(pg.sizerBox, handle, 1, func(sf *wx.SizerFlags) { sf.Expand() })
	pg.win.Layout()
	wx.WindowFromPtr(handle).Bind(wx.EVT_LEFT_DOWN, func(evt *wx.Event) { f.logf("wxStaticBitmap clicked.") })
	wx.WindowFromPtr(handle).Refresh(true)
}
