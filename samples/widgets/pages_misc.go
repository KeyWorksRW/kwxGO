// Pages demonstrating toggle buttons, radio boxes, gauges, sliders, activity
// indicators, search controls, hyperlinks and colour pickers
// (toggle.cpp, radiobox.cpp, gauge.cpp, slider.cpp, activityindicator.cpp,
// searchctrl.cpp, hyperlnk.cpp, clrpicker.cpp).
package main

import (
	"strconv"
	"unsafe"

	"github.com/KeyWorksRW/kwxGO/wx"
)

// --- ToggleButton page (toggle.cpp) ----------------------------------------

var (
	idToggleReset    = wx.ID_HIGHEST
	idToggleChgLabel = wx.ID_HIGHEST + 1
	idTogglePicker   = wx.ID_HIGHEST + 2
)

type togglePage struct {
	pageBase

	chkBitmapOnly *wx.CheckBox
	chkTextBitmap *wx.CheckBox
	chkMarkup     *wx.CheckBox
	chkFit        *wx.CheckBox
	chkDisable    *wx.CheckBox
	chkBitmapCls  *wx.CheckBox
	radioImagePos *wx.RadioBox
	radioHAlign   *wx.RadioBox
	radioVAlign   *wx.RadioBox
	textLabel     *wx.TextCtrl
	toggle        *wx.ToggleButton
	sizerToggle   *wx.BoxSizer
}

func init() {
	registerPage("ToggleButton", nativeCtrls, func(book unsafe.Pointer) *pageBase {
		pg := &togglePage{}
		pb := newPageWithIcon(book, "icons/toggle.svg", 32, 32)
		if pb == nil {
			return nil
		}
		pg.pageBase = *pb
		pg.content = pg.create
		pg.getWidgets = func() []unsafe.Pointer {
			if pg.toggle == nil {
				return nil
			}
			return []unsafe.Pointer{pg.toggle.Ptr()}
		}
		pg.recreate = func(f *widgetsFrame) { pg.createToggle(f) }
		return &pg.pageBase
	})
}

func (pg *togglePage) create(f *widgetsFrame) {
	top := wx.NewBoxSizer(wx.HORIZONTAL)
	if top == nil {
		return
	}

	boxL, left := addStaticBox(pg.win.Ptr(), "Styles", wx.VERTICAL)
	if left != nil {
		pg.chkBitmapOnly = pg.createCheckBoxAndAddToSizer(left, "&Bitmap only", wx.ID_ANY, boxL.Ptr())
		pg.chkTextBitmap = pg.createCheckBoxAndAddToSizer(left, "Text &and bitmap", wx.ID_ANY, boxL.Ptr())
		pg.chkMarkup = pg.createCheckBoxAndAddToSizer(left, "Interpret &markup", wx.ID_ANY, boxL.Ptr())
		pg.chkFit = pg.createCheckBoxAndAddToSizer(left, "&Fit exactly", wx.ID_ANY, boxL.Ptr())
		pg.chkFit.SetValue(1)
		pg.chkDisable = pg.createCheckBoxAndAddToSizer(left, "Disable", wx.ID_ANY, boxL.Ptr())
		pg.chkBitmapCls = pg.createCheckBoxAndAddToSizer(left, "Use wxBitmapToggleButton", wx.ID_ANY, boxL.Ptr())
		pg.chkBitmapCls.SetValue(1)
		addSpacerSizer(left, 10)

		pg.radioImagePos = newRadioBox(boxL.Ptr(), wx.ID_ANY, "Image &position",
			[]string{"left", "right", "top", "bottom"}, 0, 0)
		addExpandBorder(left, pg.radioImagePos.Ptr())
		addSpacerSizer(left, 15)

		pg.radioHAlign = newRadioBox(boxL.Ptr(), wx.ID_ANY, "&Horz alignment",
			[]string{"left", "centre", "right"}, 0, 0)
		if pg.radioHAlign != nil {
			pg.radioHAlign.SetSelection(1)
		}
		addExpandBorder(left, pg.radioHAlign.Ptr())
		pg.radioVAlign = newRadioBox(boxL.Ptr(), wx.ID_ANY, "&Vert alignment",
			[]string{"top", "centre", "bottom"}, 0, 0)
		if pg.radioVAlign != nil {
			pg.radioVAlign.SetSelection(1)
		}
		addExpandBorder(left, pg.radioVAlign.Ptr())
		addSpacerSizer(left, 5)

		reset := wx.NewButton(boxL.Ptr(), idToggleReset, "&Reset", -1, -1, -1, -1, 0)
		addWin(left, reset.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.CenterHorizontal()
			sf.TripleBorder(wx.ALL)
		})
		reset.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			pg.reset()
			pg.createToggle(f)
		})
	}

	boxM, middle := addStaticBox(pg.win.Ptr(), "&Operations", wx.VERTICAL)
	if middle != nil {
		row := pg.createSizerWithTextAndButton(idToggleChgLabel, "Change label", wx.ID_ANY, &pg.textLabel, boxM.Ptr())
		addSub(middle, row.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})
		if pg.textLabel != nil {
			pg.textLabel.SetValue("&Toggle me!")
		}
		pg.win.BindId(wx.EVT_BUTTON, idToggleChgLabel, func(evt *wx.Event) {
			if pg.toggle == nil || pg.textLabel == nil {
				return
			}
			if pg.chkMarkup != nil && pg.chkMarkup.GetValue() {
				pg.toggle.SetLabelMarkup(pg.textLabel.GetValue())
			} else {
				pg.toggle.SetLabel(pg.textLabel.GetValue())
			}
			if pg.chkBitmapOnly != nil && pg.chkBitmapOnly.GetValue() {
				pg.createToggle(f)
			}
		})
	}

	pg.sizerToggle = wx.NewBoxSizer(wx.HORIZONTAL)
	if pg.sizerToggle != nil {
		pg.sizerToggle.SetMinSize(150, 0)
	}

	pg.reset()
	pg.createToggle(f)

	addSub(top, left.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.LEFT)
	})
	addSub(top, middle.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL)
	})
	addSub(top, pg.sizerToggle.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.RIGHT)
	})

	pg.win.SetSizer(top.Ptr(), false)
	bindStyleCheckboxes(pg.win, func() { pg.createToggle(f) })
}

func (pg *togglePage) reset() {
	setChk(pg.chkBitmapOnly, 0)
	setChk(pg.chkTextBitmap, 0)
	setChk(pg.chkMarkup, 0)
	setChk(pg.chkFit, 1)
	setChk(pg.chkDisable, 0)
	setChk(pg.chkBitmapCls, 1)
	if pg.radioImagePos != nil {
		pg.radioImagePos.SetSelection(0)
	}
	if pg.radioHAlign != nil {
		pg.radioHAlign.SetSelection(1)
	}
	if pg.radioVAlign != nil {
		pg.radioVAlign.SetSelection(1)
	}
}

// createToggle is ToggleWidgetsPage::CreateToggle().
func (pg *togglePage) createToggle(f *widgetsFrame) {
	if pg.sizerToggle == nil {
		return
	}

	label := "&Toggle me!"
	value := false
	if pg.toggle != nil {
		label = pg.toggle.GetLabel()
		value = pg.toggle.GetValue()
		pg.toggle.Destroy()
		pg.toggle = nil
	}
	clearSizer(pg.sizerToggle)
	if label == "" && pg.textLabel != nil {
		label = pg.textLabel.GetValue()
	}

	flags := f.attrs.defaultFlags
	if pg.radioHAlign != nil {
		switch pg.radioHAlign.GetSelection() {
		case 1:
			// default alignment
		case 2:
			flags |= wx.BU_RIGHT
		default:
			flags |= wx.BU_LEFT
		}
	}
	if pg.radioVAlign != nil {
		switch pg.radioVAlign.GetSelection() {
		case 2:
			flags |= wx.BU_BOTTOM
		case 0:
			flags |= wx.BU_TOP
		}
	}
	if pg.chkFit != nil && pg.chkFit.GetValue() {
		flags |= wx.BU_EXACTFIT
	}

	if pg.chkBitmapOnly != nil && pg.chkBitmapOnly.GetValue() {
		var toggle *wx.ToggleButton
		if pg.chkBitmapCls != nil && pg.chkBitmapCls.GetValue() {
			bmp := wx.ArtProvider_GetBitmap(wx.ART_INFORMATION, wx.ART_BUTTON, -1, -1)
			bt := wx.NewBitmapToggleButton(pg.win.Ptr(), idTogglePicker, bmp.Ptr(), -1, -1, -1, -1, flags)
			if bt != nil {
				toggle = wx.ToggleButtonFromPtr(bt.Ptr())
			}
		} else {
			toggle = wx.NewToggleButton(pg.win.Ptr(), idTogglePicker, "", -1, -1, -1, -1, flags)
		}
		pg.toggle = toggle
	} else {
		pg.toggle = wx.NewToggleButton(pg.win.Ptr(), idTogglePicker, label, -1, -1, -1, -1, flags)
	}
	if pg.toggle == nil {
		return
	}
	pg.toggle.SetValue(value)
	if pg.chkDisable != nil && pg.chkDisable.GetValue() {
		wx.WindowFromPtr(pg.toggle.Ptr()).Enable(false)
	}
	pg.toggle.Bind(wx.EVT_TOGGLEBUTTON, func(evt *wx.Event) {
		cmd := wx.AsCommandEvent(evt)
		if cmd == nil || pg.toggle == nil {
			return
		}
		f.logf("Button toggled, currently %s (event) / %s (control)",
			boolWord(cmd.IsChecked()), boolWord(pg.toggle.GetValue()))
	})

	if pg.chkFit != nil && pg.chkFit.GetValue() {
		pg.sizerToggle.AddStretchSpacer(0)
		addWin(pg.sizerToggle, pg.toggle.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Center()
			sf.Border(wx.ALL)
		})
		pg.sizerToggle.AddStretchSpacer(0)
	} else {
		addWin(pg.sizerToggle, pg.toggle.Ptr(), 1, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})
	}
	pg.sizerToggle.Layout()
	pg.notifyWidgetRecreation(f, pg.toggle.Ptr())
}

func boolWord(b bool) string {
	if b {
		return "toggled"
	}
	return "not toggled"
}

func setChk(c *wx.CheckBox, v int) {
	if c != nil {
		c.SetValue(v)
	}
}

// --- Radio page (radiobox.cpp) ---------------------------------------------

const (
	radioDefaultEntries = 12
	radioDefaultDim     = 3
	radioTestButton     = 1
)

var (
	idRadioReset     = wx.ID_HIGHEST
	idRadioUpdate    = wx.ID_HIGHEST + 1
	idRadioLabelBtn  = wx.ID_HIGHEST + 2
	idRadioSelection = wx.ID_HIGHEST + 3
	idRadioLabel     = wx.ID_HIGHEST + 4
	idRadioEnableItm = wx.ID_HIGHEST + 5
	idRadioShowItm   = wx.ID_HIGHEST + 6
	idRadioCtrl      = wx.ID_HIGHEST + 7
)

type radioPage struct {
	pageBase

	chkSpecifyRows *wx.CheckBox
	textMajorDim   *wx.TextCtrl
	textNumBtns    *wx.TextCtrl
	textCurSel     *wx.TextCtrl
	textSel        *wx.TextCtrl
	textLabel      *wx.TextCtrl
	textLabelBtns  *wx.TextCtrl
	chkEnableItem  *wx.CheckBox
	chkShowItem    *wx.CheckBox
	radio          *wx.RadioBox
	sizerRadio     *wx.BoxSizer
}

func init() {
	registerPage("Radio", nativeCtrls|itemsCtrls, func(book unsafe.Pointer) *pageBase {
		pg := &radioPage{}
		pb := newPageWithIcon(book, "icons/radiobox.svg", 32, 32)
		if pb == nil {
			return nil
		}
		pg.pageBase = *pb
		pg.content = pg.create
		pg.getWidgets = func() []unsafe.Pointer {
			if pg.radio == nil {
				return nil
			}
			return []unsafe.Pointer{pg.radio.Ptr()}
		}
		pg.recreate = func(f *widgetsFrame) { pg.createRadio(f) }
		return &pg.pageBase
	})
}

func (pg *radioPage) create(f *widgetsFrame) {
	top := wx.NewBoxSizer(wx.HORIZONTAL)
	if top == nil {
		return
	}

	boxL, left := addStaticBox(pg.win.Ptr(), "&Set style", wx.VERTICAL)
	if left != nil {
		pg.chkSpecifyRows = pg.createCheckBoxAndAddToSizer(left, "Major specifies &rows count", wx.ID_ANY, boxL.Ptr())
		pg.textMajorDim = pg.rowStaticTextAndText(left, boxL.Ptr(), "&Major dimension:", "3")
		pg.textNumBtns = pg.rowStaticTextAndText(left, boxL.Ptr(), "&Number of buttons:", "12")

		update := wx.NewButton(boxL.Ptr(), idRadioUpdate, "&Update", -1, -1, -1, -1, 0)
		addWin(left, update.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.CenterHorizontal()
			sf.Border(wx.ALL)
		})
		addSpacerSizer(left, 5)
		reset := wx.NewButton(boxL.Ptr(), idRadioReset, "&Reset", -1, -1, -1, -1, 0)
		addWin(left, reset.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.CenterHorizontal()
			sf.TripleBorder(wx.ALL)
		})
		reset.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			pg.reset()
			pg.createRadio(f)
		})
		update.Bind(wx.EVT_BUTTON, func(evt *wx.Event) { pg.createRadio(f) })
	}

	boxM, middle := addStaticBox(pg.win.Ptr(), "&Change parameters", wx.VERTICAL)
	if middle != nil {
		pg.textCurSel = pg.rowStaticTextAndText(middle, boxM.Ptr(), "Current selection:", "")
		pg.textSel = pg.rowButtonAndText(middle, boxM.Ptr(), idRadioSelection, "&Change selection:", "")
		pg.textLabel = pg.rowButtonAndText(middle, boxM.Ptr(), idRadioLabel, "&Label for box:", "")
		pg.textLabel.SetValue("I'm a radiobox")
		pg.textLabelBtns = pg.rowButtonAndText(middle, boxM.Ptr(), idRadioLabelBtn, "&Label for buttons:", "item")

		pg.chkEnableItem = pg.createCheckBoxAndAddToSizer(middle, "Disable &2nd item", idRadioEnableItm, boxM.Ptr())
		pg.chkEnableItem.SetValue(1)
		pg.chkShowItem = pg.createCheckBoxAndAddToSizer(middle, "Hide 2nd &item", idRadioShowItm, boxM.Ptr())
		pg.chkShowItem.SetValue(1)

		pg.win.BindId(wx.EVT_BUTTON, idRadioLabelBtn, func(evt *wx.Event) { pg.createRadio(f) })
		pg.win.BindId(wx.EVT_BUTTON, idRadioSelection, func(evt *wx.Event) {
			if pg.radio == nil || pg.textSel == nil {
				return
			}
			if n, err := strconv.Atoi(pg.textSel.GetValue()); err == nil && n >= 0 && n < pg.radio.GetCount() {
				pg.radio.SetSelection(n)
			} else {
				f.logf("Invalid radio box selection.")
			}
		})
		pg.win.BindId(wx.EVT_BUTTON, idRadioLabel, func(evt *wx.Event) {
			if pg.radio != nil && pg.textLabel != nil {
				pg.radio.SetLabel(pg.textLabel.GetValue())
			}
		})
		pg.chkEnableItem.Bind(wx.EVT_CHECKBOX, func(evt *wx.Event) {
			if pg.radio != nil {
				cmd := wx.AsCommandEvent(evt)
				if cmd != nil {
					pg.radio.Enable(radioTestButton, cmd.IsChecked())
				}
			}
		})
		pg.chkShowItem.Bind(wx.EVT_CHECKBOX, func(evt *wx.Event) {
			if pg.radio != nil {
				cmd := wx.AsCommandEvent(evt)
				if cmd != nil {
					pg.radio.Show(radioTestButton, cmd.IsChecked())
				}
			}
		})
	}

	pg.sizerRadio = wx.NewBoxSizer(wx.HORIZONTAL)
	if pg.sizerRadio == nil {
		return
	}

	pg.reset()
	pg.createRadio(f)

	addSub(top, left.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.LEFT)
	})
	addSub(top, middle.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL)
	})
	addSub(top, pg.sizerRadio.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.RIGHT)
	})

	pg.win.SetSizer(top.Ptr(), false)
	bindStyleCheckboxes(pg.win, func() { pg.createRadio(f) })
}

func (pg *radioPage) reset() {
	setChk(pg.chkSpecifyRows, 0)
	if pg.textNumBtns != nil {
		pg.textNumBtns.SetValue(strconv.Itoa(radioDefaultEntries))
	}
	if pg.textMajorDim != nil {
		pg.textMajorDim.SetValue(strconv.Itoa(radioDefaultDim))
	}
}

// createRadio is RadioWidgetsPage::CreateRadio().
func (pg *radioPage) createRadio(f *widgetsFrame) {
	if pg.sizerRadio == nil {
		return
	}

	sel := -1
	if pg.radio != nil {
		sel = pg.radio.GetSelection()
		pg.radio.Destroy()
		pg.radio = nil
	}
	clearSizer(pg.sizerRadio)

	count := radioDefaultEntries
	if pg.textNumBtns != nil {
		if n, err := strconv.Atoi(pg.textNumBtns.GetValue()); err == nil {
			count = n
		}
	}
	majorDim := radioDefaultDim
	if pg.textMajorDim != nil {
		if n, err := strconv.Atoi(pg.textMajorDim.GetValue()); err == nil {
			majorDim = n
		}
	}

	prefix := "item"
	if pg.textLabelBtns != nil {
		prefix = pg.textLabelBtns.GetValue()
	}
	items := make([]string, 0, count)
	for n := 0; n < count; n++ {
		items = append(items, prefix+" "+strconv.Itoa(n+1))
	}

	flags := f.attrs.defaultFlags
	if pg.chkSpecifyRows != nil && pg.chkSpecifyRows.GetValue() {
		flags |= wx.RA_SPECIFY_ROWS
	} else {
		flags |= wx.RA_SPECIFY_COLS
	}

	label := "I'm a radiobox"
	if pg.textLabel != nil {
		label = pg.textLabel.GetValue()
	}

	pg.radio = newRadioBox(pg.win.Ptr(), idRadioCtrl, label, items, majorDim, flags)
	if pg.radio == nil {
		return
	}
	if sel >= 0 && sel < count {
		pg.radio.SetSelection(sel)
	}
	if count > 1 {
		pg.radio.Enable(radioTestButton, pg.chkEnableItem.GetValue())
		pg.radio.Show(radioTestButton, pg.chkShowItem.GetValue())
	}
	pg.radio.Bind(wx.EVT_RADIOBOX, func(evt *wx.Event) {
		cmd := wx.AsCommandEvent(evt)
		if cmd == nil || pg.radio == nil {
			return
		}
		f.logf("Radio item %d selected", cmd.GetSelection())
		if pg.textCurSel != nil {
			pg.textCurSel.SetValue(strconv.Itoa(pg.radio.GetSelection()))
		}
	})

	addWin(pg.sizerRadio, pg.radio.Ptr(), 1, func(sf *wx.SizerFlags) { sf.Expand() })
	pg.sizerRadio.Layout()
	pg.notifyWidgetRecreation(f, pg.radio.Ptr())
}

// --- Gauge page (gauge.cpp) -------------------------------------------------

var (
	idGaugeReset      = wx.ID_HIGHEST
	idGaugeProgress   = wx.ID_HIGHEST + 1
	idGaugeIndet      = wx.ID_HIGHEST + 2
	idGaugeClear      = wx.ID_HIGHEST + 3
	idGaugeSetValue   = wx.ID_HIGHEST + 4
	idGaugeSetRange   = wx.ID_HIGHEST + 5
	idGaugeCurText    = wx.ID_HIGHEST + 6
	idGaugeValueText  = wx.ID_HIGHEST + 7
	idGaugeRangeText  = wx.ID_HIGHEST + 8
	idGaugeTimer      = wx.ID_HIGHEST + 9
	idGaugeIndetTimer = wx.ID_HIGHEST + 10
	idGaugeCtrl       = wx.ID_HIGHEST + 11
)

type gaugePage struct {
	pageBase

	chkVert   *wx.CheckBox
	chkSmooth *wx.CheckBox
	chkProg   *wx.CheckBox
	btnProg   *wx.Button
	btnIndet  *wx.Button
	textValue *wx.TextCtrl
	textRange *wx.TextCtrl
	textCur   *wx.TextCtrl
	gauge     *wx.Gauge
	sizerG    *wx.BoxSizer
	timer     *wx.Timer
	rng       int
}

func init() {
	registerPage("Gauge", nativeCtrls, func(book unsafe.Pointer) *pageBase {
		pg := &gaugePage{rng: 100}
		pb := newPageWithIcon(book, "icons/gauge.svg", 32, 32)
		if pb == nil {
			return nil
		}
		pg.pageBase = *pb
		pg.content = pg.create
		pg.getWidgets = func() []unsafe.Pointer {
			if pg.gauge == nil {
				return nil
			}
			return []unsafe.Pointer{pg.gauge.Ptr()}
		}
		pg.recreate = func(f *widgetsFrame) { pg.createGauge(f) }
		return &pg.pageBase
	})
}

func (pg *gaugePage) create(f *widgetsFrame) {
	top := wx.NewBoxSizer(wx.HORIZONTAL)
	if top == nil {
		return
	}

	boxL, left := addStaticBox(pg.win.Ptr(), "&Set style", wx.VERTICAL)
	if left != nil {
		pg.chkVert = pg.createCheckBoxAndAddToSizer(left, "&Vertical", wx.ID_ANY, boxL.Ptr())
		pg.chkSmooth = pg.createCheckBoxAndAddToSizer(left, "&Smooth", wx.ID_ANY, boxL.Ptr())
		pg.chkProg = pg.createCheckBoxAndAddToSizer(left, "&Progress", wx.ID_ANY, boxL.Ptr())
		addSpacerSizer(left, 15)
		reset := wx.NewButton(boxL.Ptr(), idGaugeReset, "&Reset", -1, -1, -1, -1, 0)
		addWin(left, reset.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.CenterHorizontal()
			sf.TripleBorder(wx.ALL)
		})
		reset.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			pg.reset()
			pg.createGauge(f)
		})
	}

	boxM, middle := addStaticBox(pg.win.Ptr(), "&Change gauge value", wx.VERTICAL)
	if middle != nil {
		cur := pg.rowStaticTextAndText(middle, boxM.Ptr(), "Current value", "")
		cur.SetEditable(false)
		pg.textCur = cur
		// The C++ sample refreshes this field from an EVT_UPDATE_UI handler;
		// wxUpdateUIEvent::SetText through the binding crashes the update pass,
		// so the field is refreshed from the gauge's own value changes instead.

		pg.textValue = pg.rowButtonAndText(middle, boxM.Ptr(), idGaugeSetValue, "Set &value", "")
		pg.textRange = pg.rowButtonAndText(middle, boxM.Ptr(), idGaugeSetRange, "Set &range", "100")

		pg.btnProg = wx.NewButton(boxM.Ptr(), idGaugeProgress, "Simulate &progress", -1, -1, -1, -1, 0)
		addExpandBorder(middle, pg.btnProg.Ptr())
		pg.btnIndet = wx.NewButton(boxM.Ptr(), idGaugeIndet, "Simulate &indeterminate job", -1, -1, -1, -1, 0)
		addExpandBorder(middle, pg.btnIndet.Ptr())
		clear := wx.NewButton(boxM.Ptr(), idGaugeClear, "&Clear", -1, -1, -1, -1, 0)
		addExpandBorder(middle, clear.Ptr())

		pg.btnProg.Bind(wx.EVT_BUTTON, func(evt *wx.Event) { pg.toggleTimer(f, pg.btnProg, idGaugeTimer) })
		pg.btnIndet.Bind(wx.EVT_BUTTON, func(evt *wx.Event) { pg.toggleTimer(f, pg.btnIndet, idGaugeIndetTimer) })
		clear.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.gauge != nil {
				pg.gauge.SetValue(0)
				pg.updateCur()
			}
		})
		pg.win.BindId(wx.EVT_BUTTON, idGaugeSetValue, func(evt *wx.Event) {
			if pg.gauge == nil || pg.textValue == nil {
				return
			}
			if n, err := strconv.Atoi(pg.textValue.GetValue()); err == nil {
				pg.gauge.SetValue(n)
			}
		})
		pg.win.BindId(wx.EVT_BUTTON, idGaugeSetRange, func(evt *wx.Event) {
			if pg.gauge == nil || pg.textRange == nil {
				return
			}
			if n, err := strconv.Atoi(pg.textRange.GetValue()); err == nil {
				pg.rng = n
				pg.gauge.SetRange(n)
			}
		})
	}

	pg.sizerG = wx.NewBoxSizer(wx.HORIZONTAL)
	if pg.sizerG != nil {
		pg.sizerG.SetMinSize(150, 0)
	}

	addSub(top, left.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.LEFT)
	})
	addSub(top, middle.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL)
	})
	addSub(top, pg.sizerG.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.RIGHT)
	})

	pg.reset()
	pg.win.SetSizer(top.Ptr(), false)
	bindStyleCheckboxes(pg.win, func() { pg.createGauge(f) })
}

func (pg *gaugePage) reset() {
	setChk(pg.chkVert, 0)
	setChk(pg.chkSmooth, 0)
	setChk(pg.chkProg, 0)
}

// createGauge is GaugeWidgetsPage::CreateGauge().
func (pg *gaugePage) createGauge(f *widgetsFrame) {
	if pg.sizerG == nil {
		return
	}

	flags := f.attrs.defaultFlags
	flags |= wx.GA_HORIZONTAL
	if pg.chkSmooth != nil && pg.chkSmooth.GetValue() {
		flags |= wx.GA_SMOOTH
	}
	if pg.chkProg != nil && pg.chkProg.GetValue() {
		flags |= wx.GA_PROGRESS
	}

	val := 0
	if pg.gauge != nil {
		val = pg.gauge.GetValue()
		pg.gauge.Destroy()
		pg.gauge = nil
	}
	clearSizer(pg.sizerG)

	pg.gauge = wx.NewGauge(pg.win.Ptr(), idGaugeCtrl, pg.rng, -1, -1, -1, -1, flags)
	if pg.gauge == nil {
		return
	}
	pg.gauge.SetValue(val)
	addWin(pg.sizerG, pg.gauge.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Center()
		sf.Border(wx.ALL)
	})
	pg.sizerG.Layout()
	pg.updateCur()
	pg.notifyWidgetRecreation(f, pg.gauge.Ptr())
}

// updateCur refreshes the "Current value" field (C++ does it in EVT_UPDATE_UI).
func (pg *gaugePage) updateCur() {
	if pg.textCur != nil && pg.gauge != nil {
		pg.textCur.SetValue(strconv.Itoa(pg.gauge.GetValue()))
	}
}

// toggleTimer is StartTimer()/StopTimer() of gauge.cpp.
func (pg *gaugePage) toggleTimer(f *widgetsFrame, btn *wx.Button, id int) {
	if pg.gauge == nil {
		return
	}

	if pg.timer == nil {
		pg.timer = wx.NewTimer(pg.win.Ptr(), id)
		if pg.timer == nil {
			return
		}
		progress := id == idGaugeTimer
		pg.timer.Bind(wx.EVT_TIMER, func(evt *wx.Event) {
			if pg.gauge == nil {
				return
			}
			if progress {
				if v := pg.gauge.GetValue(); v < pg.gauge.GetRange() {
					pg.gauge.SetValue(v + 1)
				} else {
					pg.stopTimer(f, btn, id)
				}
			} else {
				pg.gauge.Pulse()
			}
		})
		pg.timer.Start(300, false)
		btn.SetLabel("&Stop timer")
		if progress {
			wx.WindowFromPtr(pg.btnIndet.Ptr()).Enable(false)
		} else {
			wx.WindowFromPtr(pg.btnProg.Ptr()).Enable(false)
		}
		f.logf("Starting the timer.")
		return
	}

	pg.stopTimer(f, btn, id)
}

func (pg *gaugePage) stopTimer(f *widgetsFrame, btn *wx.Button, id int) {
	if pg.timer != nil {
		pg.timer.Stop()
		pg.timer.Delete()
		pg.timer = nil
	}
	if id == idGaugeTimer {
		btn.SetLabel("Simulate &progress")
		wx.WindowFromPtr(pg.btnIndet.Ptr()).Enable(true)
	} else {
		btn.SetLabel("Simulate &indeterminate job")
		wx.WindowFromPtr(pg.btnProg.Ptr()).Enable(true)
		if pg.gauge != nil {
			pg.gauge.SetValue(0)
		}
	}
	f.logf("Progress finished.")
}

// --- Slider page (slider.cpp) ----------------------------------------------

var (
	idSliderReset     = wx.ID_HIGHEST
	idSliderSetValue  = wx.ID_HIGHEST + 1
	idSliderSetMinMax = wx.ID_HIGHEST + 2
	idSliderSetRange  = wx.ID_HIGHEST + 3
	idSliderSetLine   = wx.ID_HIGHEST + 4
	idSliderSetPage   = wx.ID_HIGHEST + 5
	idSliderSetTick   = wx.ID_HIGHEST + 6
	idSliderSetThumb  = wx.ID_HIGHEST + 7
	idSliderCurText   = wx.ID_HIGHEST + 8
	idSliderValueText = wx.ID_HIGHEST + 9
	idSliderMinText   = wx.ID_HIGHEST + 10
	idSliderMaxText   = wx.ID_HIGHEST + 11
	idSliderRngMin    = wx.ID_HIGHEST + 12
	idSliderRngMax    = wx.ID_HIGHEST + 13
	idSliderLineText  = wx.ID_HIGHEST + 14
	idSliderPageText  = wx.ID_HIGHEST + 15
	idSliderTickText  = wx.ID_HIGHEST + 16
	idSliderThumbText = wx.ID_HIGHEST + 17
	idSliderRadioSide = wx.ID_HIGHEST + 18
	idSliderBoth      = wx.ID_HIGHEST + 19
	idSliderSelRange  = wx.ID_HIGHEST + 20
	idSliderCtrl      = wx.ID_HIGHEST + 21
)

type sliderPage struct {
	pageBase

	chkInverse   *wx.CheckBox
	chkTicks     *wx.CheckBox
	chkMinMax    *wx.CheckBox
	chkValue     *wx.CheckBox
	radioSides   *wx.RadioBox
	chkBoth      *wx.CheckBox
	chkSelRange  *wx.CheckBox
	textValue    *wx.TextCtrl
	textMin      *wx.TextCtrl
	textMax      *wx.TextCtrl
	textRangeMin *wx.TextCtrl
	textRangeMax *wx.TextCtrl
	textLine     *wx.TextCtrl
	textPage     *wx.TextCtrl
	textTick     *wx.TextCtrl
	textThumb    *wx.TextCtrl
	textCur      *wx.TextCtrl
	slider       *wx.Slider
	sizerSl      *wx.BoxSizer

	min, max int
}

func init() {
	registerPage("Slider", nativeCtrls, func(book unsafe.Pointer) *pageBase {
		pg := &sliderPage{max: 100}
		pb := newPageWithIcon(book, "icons/slider.svg", 32, 32)
		if pb == nil {
			return nil
		}
		pg.pageBase = *pb
		pg.content = pg.create
		pg.getWidgets = func() []unsafe.Pointer {
			if pg.slider == nil {
				return nil
			}
			return []unsafe.Pointer{pg.slider.Ptr()}
		}
		pg.recreate = func(f *widgetsFrame) { pg.createSlider(f) }
		return &pg.pageBase
	})
}

func (pg *sliderPage) create(f *widgetsFrame) {
	top := wx.NewBoxSizer(wx.HORIZONTAL)
	if top == nil {
		return
	}

	boxL, left := addStaticBox(pg.win.Ptr(), "&Set style", wx.VERTICAL)
	if left != nil {
		pg.chkInverse = pg.createCheckBoxAndAddToSizer(left, "&Inverse", wx.ID_ANY, boxL.Ptr())
		pg.chkTicks = pg.createCheckBoxAndAddToSizer(left, "Show &ticks", wx.ID_ANY, boxL.Ptr())
		pg.chkMinMax = pg.createCheckBoxAndAddToSizer(left, "Show min/max &labels", wx.ID_ANY, boxL.Ptr())
		pg.chkValue = pg.createCheckBoxAndAddToSizer(left, "Show &value label", wx.ID_ANY, boxL.Ptr())

		pg.radioSides = newRadioBox(boxL.Ptr(), idSliderRadioSide, "&Label position",
			[]string{"default", "top", "bottom", "left", "right"}, 0, wx.RA_SPECIFY_COLS)
		addExpandBorder(left, pg.radioSides.Ptr())
		pg.chkBoth = pg.createCheckBoxAndAddToSizer(left, "&Both sides", idSliderBoth, boxL.Ptr())
		pg.chkSelRange = pg.createCheckBoxAndAddToSizer(left, "&Selection range", idSliderSelRange, boxL.Ptr())
		addSpacerSizer(left, 5)

		reset := wx.NewButton(boxL.Ptr(), idSliderReset, "&Reset", -1, -1, -1, -1, 0)
		addWin(left, reset.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.CenterHorizontal()
			sf.TripleBorder(wx.ALL)
		})
		reset.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			pg.reset()
			pg.createSlider(f)
		})
	}

	boxM, middle := addStaticBox(pg.win.Ptr(), "&Change slider value", wx.VERTICAL)
	if middle != nil {
		cur := pg.rowStaticTextAndText(middle, boxM.Ptr(), "Current value", "")
		cur.SetEditable(false)
		pg.textCur = cur
		// See the gauge page: the update-UI based refresh is replaced by a
		// direct update from the slider's own value changes.
		pg.textValue = pg.rowButtonAndText(middle, boxM.Ptr(), idSliderSetValue, "Set &value", "")

		pg.textMin = pg.rowButtonAndText(middle, boxM.Ptr(), idSliderSetMinMax, "&Min and max", "0")
		pg.textMax = wx.NewTextCtrl(boxM.Ptr(), idSliderMaxText, "100", -1, -1, -1, -1, 0)
		pg.addToLastRow(pg.textMin, pg.textMax)

		pg.textRangeMin = pg.rowButtonAndText(middle, boxM.Ptr(), idSliderSetRange, "&Selection", "20")
		pg.textRangeMax = wx.NewTextCtrl(boxM.Ptr(), idSliderRngMax, "80", -1, -1, -1, -1, 0)
		pg.addToLastRow(pg.textRangeMin, pg.textRangeMax)

		pg.textLine = pg.rowButtonAndText(middle, boxM.Ptr(), idSliderSetLine, "Li&ne size", "")
		pg.textPage = pg.rowButtonAndText(middle, boxM.Ptr(), idSliderSetPage, "P&age size", "")
		pg.textTick = pg.rowButtonAndText(middle, boxM.Ptr(), idSliderSetTick, "Tick &frequency", "10")
		pg.textThumb = pg.rowButtonAndText(middle, boxM.Ptr(), idSliderSetThumb, "Thumb &length", "")

		pg.win.BindId(wx.EVT_BUTTON, idSliderSetValue, func(evt *wx.Event) {
			if pg.slider == nil || pg.textValue == nil {
				return
			}
			if n, err := strconv.Atoi(pg.textValue.GetValue()); err == nil && pg.validValue(n) {
				pg.slider.SetValue(n)
			} else {
				f.logf("Invalid slider value.")
			}
		})
		pg.win.BindId(wx.EVT_BUTTON, idSliderSetMinMax, func(evt *wx.Event) {
			if pg.slider == nil || pg.textMin == nil || pg.textMax == nil {
				return
			}
			mn, e1 := strconv.Atoi(pg.textMin.GetValue())
			mx, e2 := strconv.Atoi(pg.textMax.GetValue())
			if e1 != nil || e2 != nil || mn >= mx {
				f.logf("Invalid min/max values.")
				return
			}
			pg.min, pg.max = mn, mx
			pg.slider.SetRange(mn, mx)
		})
		pg.win.BindId(wx.EVT_BUTTON, idSliderSetRange, func(evt *wx.Event) {
			pg.doSetSelectionRange()
		})
		pg.win.BindId(wx.EVT_BUTTON, idSliderSetLine, func(evt *wx.Event) {
			if pg.slider == nil || pg.textLine == nil {
				return
			}
			if n, err := strconv.Atoi(pg.textLine.GetValue()); err == nil {
				pg.slider.SetLineSize(n)
			}
		})
		pg.win.BindId(wx.EVT_BUTTON, idSliderSetPage, func(evt *wx.Event) {
			if pg.slider == nil || pg.textPage == nil {
				return
			}
			if n, err := strconv.Atoi(pg.textPage.GetValue()); err == nil {
				pg.slider.SetPageSize(n)
			}
		})
		pg.win.BindId(wx.EVT_BUTTON, idSliderSetTick, func(evt *wx.Event) {
			if pg.slider == nil || pg.textTick == nil {
				return
			}
			if n, err := strconv.Atoi(pg.textTick.GetValue()); err == nil {
				pg.slider.SetTickFreq(n)
			}
		})
		pg.win.BindId(wx.EVT_BUTTON, idSliderSetThumb, func(evt *wx.Event) {
			if pg.slider == nil || pg.textThumb == nil {
				return
			}
			if n, err := strconv.Atoi(pg.textThumb.GetValue()); err == nil {
				pg.slider.SetThumbLength(n)
			}
		})
	}

	pg.sizerSl = wx.NewBoxSizer(wx.HORIZONTAL)
	if pg.sizerSl == nil {
		return
	}

	pg.reset()
	pg.createSlider(f)

	addSub(top, left.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.LEFT)
	})
	addSub(top, middle.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL)
	})
	addSub(top, pg.sizerSl.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.RIGHT)
	})

	pg.win.SetSizer(top.Ptr(), false)
	bindStyleCheckboxes(pg.win, func() { pg.createSlider(f) })
}

func (pg *sliderPage) validValue(v int) bool { return v >= pg.min && v <= pg.max }

func (pg *sliderPage) reset() {
	setChk(pg.chkInverse, 0)
	setChk(pg.chkTicks, 1)
	setChk(pg.chkValue, 1)
	setChk(pg.chkMinMax, 1)
	setChk(pg.chkBoth, 0)
	setChk(pg.chkSelRange, 0)
	if pg.radioSides != nil {
		pg.radioSides.SetSelection(0)
	}
}

func (pg *sliderPage) doSetSelectionRange() {
	if pg.slider == nil || pg.textRangeMin == nil || pg.textRangeMax == nil {
		return
	}
	mn, e1 := strconv.Atoi(pg.textRangeMin.GetValue())
	mx, e2 := strconv.Atoi(pg.textRangeMax.GetValue())
	if e1 != nil || e2 != nil {
		return
	}
	pg.slider.SetSelection(mn, mx)
}

// createSlider is SliderWidgetsPage::CreateSlider().
func (pg *sliderPage) createSlider(f *widgetsFrame) {
	if pg.sizerSl == nil {
		return
	}

	flags := f.attrs.defaultFlags
	if pg.chkInverse != nil && pg.chkInverse.GetValue() {
		flags |= wx.SL_INVERSE
	}
	if pg.chkTicks != nil && pg.chkTicks.GetValue() {
		flags |= wx.SL_AUTOTICKS
	}
	if pg.chkMinMax != nil && pg.chkMinMax.GetValue() {
		flags |= wx.SL_MIN_MAX_LABELS
	}
	if pg.chkValue != nil && pg.chkValue.GetValue() {
		flags |= wx.SL_VALUE_LABEL
	}
	if pg.chkBoth != nil && pg.chkBoth.GetValue() {
		flags |= wx.SL_BOTH
	}
	if pg.chkSelRange != nil && pg.chkSelRange.GetValue() {
		flags |= wx.SL_SELRANGE
	}
	if pg.radioSides != nil {
		switch pg.radioSides.GetSelection() {
		case 1:
			flags |= wx.SL_BOTTOM
		case 2:
			flags |= wx.SL_TOP
		case 3:
			flags |= wx.SL_RIGHT
		case 4:
			flags |= wx.SL_LEFT
		}
	}

	val := pg.min
	if pg.slider != nil {
		val = pg.slider.GetValue()
		pg.slider.Destroy()
		pg.slider = nil
	}
	clearSizer(pg.sizerSl)
	if !pg.validValue(val) {
		val = pg.min
	}

	pg.slider = wx.NewSlider(pg.win.Ptr(), idSliderCtrl, val, pg.min, pg.max, -1, -1, -1, -1, flags)
	if pg.slider == nil {
		return
	}
	addWin(pg.sizerSl, pg.slider.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Center()
		sf.Border(wx.ALL)
	})
	if pg.textLine != nil {
		pg.textLine.SetValue(strconv.Itoa(pg.slider.GetLineSize()))
	}
	if pg.textPage != nil {
		pg.textPage.SetValue(strconv.Itoa(pg.slider.GetPageSize()))
	}
	if pg.textThumb != nil {
		pg.textThumb.SetValue(strconv.Itoa(pg.slider.GetThumbLength()))
	}
	pg.slider.Bind(wx.EVT_SLIDER, func(evt *wx.Event) {
		cmd := wx.AsCommandEvent(evt)
		if cmd == nil || pg.slider == nil {
			return
		}
		f.logf("Slider value changed to %d (event %d)", pg.slider.GetValue(), cmd.GetInt())
		pg.updateCur()
	})
	pg.sizerSl.Layout()
	pg.updateCur()
	pg.notifyWidgetRecreation(f, pg.slider.Ptr())
}

// updateCur refreshes the "Current value" field (C++ does it in EVT_UPDATE_UI).
func (pg *sliderPage) updateCur() {
	if pg.textCur != nil && pg.slider != nil {
		pg.textCur.SetValue(strconv.Itoa(pg.slider.GetValue()))
	}
}

// --- ActivityIndicator page ------------------------------------------------

var (
	idActivityStart  = wx.ID_HIGHEST
	idActivityStop   = wx.ID_HIGHEST + 1
	idActivityStatus = wx.ID_HIGHEST + 2
)

type activityPage struct {
	pageBase

	indicator *wx.ActivityIndicator
	sizerInd  *wx.StaticBoxSizer
	box       *wx.StaticBox
}

func init() {
	registerPage("ActivityIndicator", nativeCtrls, func(book unsafe.Pointer) *pageBase {
		pg := &activityPage{}
		pb := newPageWithIcon(book, "icons/activityindicator.svg", 16, 16)
		if pb == nil {
			return nil
		}
		pg.pageBase = *pb
		pg.content = pg.create
		pg.getWidgets = func() []unsafe.Pointer {
			if pg.indicator == nil {
				return nil
			}
			return []unsafe.Pointer{pg.indicator.Ptr()}
		}
		pg.recreate = func(f *widgetsFrame) { pg.recreateIndicator(f) }
		return &pg.pageBase
	})
}

func (pg *activityPage) create(f *widgetsFrame) {
	top := wx.NewBoxSizer(wx.HORIZONTAL)
	if top == nil {
		return
	}

	boxO, oper := addStaticBox(pg.win.Ptr(), "&Operations", wx.VERTICAL)
	if oper != nil {
		start := wx.NewButton(boxO.Ptr(), idActivityStart, "&Start", -1, -1, -1, -1, 0)
		addExpandBorder(oper, start.Ptr())
		stop := wx.NewButton(boxO.Ptr(), idActivityStop, "&Stop", -1, -1, -1, -1, 0)
		addExpandBorder(oper, stop.Ptr())
		status := wx.NewStaticText(boxO.Ptr(), idActivityStatus, "Indicator is initializing...", -1, -1, -1, -1, 0)
		addExpandBorder(oper, status.Ptr())
		status.Bind(wx.EVT_UPDATE_UI, func(evt *wx.Event) {
			ui := wx.AsUpdateUIEvent(evt)
			if ui == nil {
				return
			}
			if pg.indicator != nil && pg.indicator.IsRunning() {
				ui.SetText("Indicator is running")
			} else {
				ui.SetText("Indicator is stopped")
			}
		})

		start.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.indicator != nil {
				pg.indicator.Start()
			}
		})
		stop.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.indicator != nil {
				pg.indicator.Stop()
			}
		})
	}

	pg.box, pg.sizerInd = addStaticBox(pg.win.Ptr(), "Activity Indicator", wx.HORIZONTAL)
	pg.recreateIndicator(f)

	addSub(top, oper.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL)
	})
	addSub(top, pg.sizerInd.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL)
	})

	pg.win.SetSizer(top.Ptr(), false)
}

func (pg *activityPage) recreateIndicator(f *widgetsFrame) {
	if pg.sizerInd == nil {
		return
	}
	destroyWin(pg.indicator)
	clearSizer(pg.sizerInd)
	pg.indicator = wx.NewActivityIndicator(pg.box.Ptr(), wx.ID_ANY, -1, -1, -1, -1, f.attrs.defaultFlags)
	if pg.indicator == nil {
		return
	}
	pg.sizerInd.AddStretchSpacer(0)
	addWin(pg.sizerInd, pg.indicator.Ptr(), 0, func(sf *wx.SizerFlags) { sf.Center() })
	pg.sizerInd.AddStretchSpacer(0)
	pg.sizerInd.Layout()
	pg.notifyWidgetRecreation(f, pg.indicator.Ptr())
}

// --- SearchCtrl page (searchctrl.cpp) --------------------------------------

var (
	idSearchCb      = wx.ID_HIGHEST
	idCancelCb      = wx.ID_HIGHEST + 1
	idMenuCb        = wx.ID_HIGHEST + 2
	idSearchMenu    = wx.ID_HIGHEST + 3
	idSearchMenuMax = idSearchMenu + 5
)

type searchPage struct {
	pageBase

	chkSearch *wx.CheckBox
	chkCancel *wx.CheckBox
	chkMenu   *wx.CheckBox
	ctrl      *wx.SearchCtrl
	sizer     *wx.BoxSizer
}

func init() {
	registerPage("SearchCtrl", genericCtrls|editableCtrls, func(book unsafe.Pointer) *pageBase {
		pg := &searchPage{}
		pb := newPageWithIcon(book, "icons/text.svg", 32, 32)
		if pb == nil {
			return nil
		}
		pg.pageBase = *pb
		pg.content = pg.create
		pg.getWidgets = func() []unsafe.Pointer {
			if pg.ctrl == nil {
				return nil
			}
			return []unsafe.Pointer{pg.ctrl.Ptr()}
		}
		pg.getEntry = func() unsafe.Pointer {
			if pg.ctrl == nil {
				return nil
			}
			return pg.ctrl.Ptr()
		}
		pg.recreate = func(f *widgetsFrame) { pg.createControl(f) }
		return &pg.pageBase
	})
}

func (pg *searchPage) create(f *widgetsFrame) {
	pg.createControl(f)

	box, options := addStaticBox(pg.win.Ptr(), "Options", wx.VERTICAL)
	pg.chkSearch = wx.NewCheckBox(box.Ptr(), idSearchCb, "Search button", -1, -1, -1, -1, 0)
	pg.chkCancel = wx.NewCheckBox(box.Ptr(), idCancelCb, "Cancel button", -1, -1, -1, -1, 0)
	pg.chkMenu = wx.NewCheckBox(box.Ptr(), idMenuCb, "Search menu", -1, -1, -1, -1, 0)
	pg.chkSearch.SetValue(1)
	addWin(options, pg.chkSearch.Ptr(), 0, func(sf *wx.SizerFlags) { sf.Border(wx.ALL) })
	addWin(options, pg.chkCancel.Ptr(), 0, func(sf *wx.SizerFlags) { sf.Border(wx.ALL) })
	addWin(options, pg.chkMenu.Ptr(), 0, func(sf *wx.SizerFlags) { sf.Border(wx.ALL) })

	pg.chkSearch.Bind(wx.EVT_CHECKBOX, func(evt *wx.Event) {
		if pg.ctrl != nil {
			pg.ctrl.ShowSearchButton(pg.chkSearch.GetValue())
		}
	})
	pg.chkCancel.Bind(wx.EVT_CHECKBOX, func(evt *wx.Event) {
		if pg.ctrl != nil {
			pg.ctrl.ShowCancelButton(pg.chkCancel.GetValue())
		}
	})
	pg.chkMenu.Bind(wx.EVT_CHECKBOX, func(evt *wx.Event) {
		if pg.ctrl == nil {
			return
		}
		if pg.chkMenu.GetValue() {
			menu := pg.createTestMenu()
			if menu != nil {
				pg.ctrl.SetMenu(menu.Ptr())
			}
		} else {
			pg.ctrl.SetMenu(nil)
		}
	})

	sizer := wx.NewBoxSizer(wx.HORIZONTAL)
	if sizer == nil {
		return
	}
	pg.sizer = sizer
	addSub(sizer, options.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.TripleBorder(wx.ALL)
	})
	addWin(sizer, pg.ctrl.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Center()
		sf.TripleBorder(wx.ALL)
	})
	pg.win.SetSizer(sizer.Ptr(), false)
}

func (pg *searchPage) createControl(f *widgetsFrame) {
	old := pg.ctrl
	pg.ctrl = wx.NewSearchCtrl(pg.win.Ptr(), -1, "", -1, -1, 150, -1, f.attrs.defaultFlags)
	if pg.ctrl == nil {
		pg.ctrl = old
		return
	}
	// keep the sizer item pointing at the surviving control before the old one
	// is destroyed
	if old != nil {
		if pg.sizer != nil {
			pg.sizer.ReplaceWindow(old.Ptr(), pg.ctrl.Ptr(), false)
		}
		// SearchWidgetsPage::CreateControl() destroys the previous control
		destroyWin(old)
	}
	pg.ctrl.Bind(wx.EVT_TEXT, func(evt *wx.Event) {
		cmd := wx.AsCommandEvent(evt)
		if cmd != nil {
			f.logf("Text changed to \"%s\"", cmd.GetString())
		}
	})
	pg.ctrl.Bind(wx.EVT_TEXT_ENTER, func(evt *wx.Event) {
		cmd := wx.AsCommandEvent(evt)
		if cmd != nil {
			f.logf("Search text entered: \"%s\"", cmd.GetString())
		}
	})
	pg.ctrl.Bind(wx.EVT_SEARCHCTRL_SEARCH_BTN, func(evt *wx.Event) {
		f.logf("Search button pressed.")
	})
	pg.ctrl.Bind(wx.EVT_SEARCHCTRL_CANCEL_BTN, func(evt *wx.Event) {
		f.logf("Cancel button pressed.")
	})
	pg.notifyWidgetRecreation(f, pg.ctrl.Ptr())
}

func (pg *searchPage) createTestMenu() *wx.Menu {
	menu := wx.NewMenu("", 0)
	if menu == nil {
		return nil
	}
	menu.Append(wx.ID_ANY, "Recent Searches", "", false)
	menu.Enable(wx.ID_ANY, false)
	for i := 0; i < 5; i++ {
		menu.AppendCheckItem(idSearchMenu+i, "item "+strconv.Itoa(i), "tip "+strconv.Itoa(i))
	}
	return menu
}

// --- Hyperlink page (hyperlnk.cpp) -----------------------------------------

var (
	idHyperReset    = wx.ID_HIGHEST
	idHyperSetLabel = wx.ID_HIGHEST + 1
	idHyperSetURL   = wx.ID_HIGHEST + 2
	idHyperCtrl     = wx.ID_HIGHEST + 3
)

type hyperlinkPage struct {
	pageBase

	radioAlign *wx.RadioBox
	chkGeneric *wx.CheckBox
	textLabel  *wx.TextCtrl
	textURL    *wx.TextCtrl
	visit      *wx.StaticText
	fun        *wx.StaticText
	hyperlink  *wx.HyperlinkCtrl
	longLink   *wx.HyperlinkCtrl
	sizerLong  *wx.BoxSizer
	shortRow   *wx.BoxSizer
}

func init() {
	registerPage("Hyperlink", genericCtrls, func(book unsafe.Pointer) *pageBase {
		pg := &hyperlinkPage{}
		pb := newPageWithIcon(book, "icons/hyperlnk.svg", 32, 32)
		if pb == nil {
			return nil
		}
		pg.pageBase = *pb
		pg.content = pg.create
		pg.getWidgets = func() []unsafe.Pointer {
			if pg.hyperlink == nil {
				return nil
			}
			return []unsafe.Pointer{pg.hyperlink.Ptr(), ptrOf(pg.longLink)}
		}
		pg.recreate = func(f *widgetsFrame) { pg.createHyperlink(f) }
		return &pg.pageBase
	})
}

func (pg *hyperlinkPage) create(f *widgetsFrame) {
	top := wx.NewBoxSizer(wx.HORIZONTAL)
	if top == nil {
		return
	}

	boxL, left := addStaticBox(pg.win.Ptr(), "Hyperlink details", wx.VERTICAL)
	if left != nil {
		pg.textLabel = pg.rowButtonAndText(left, boxL.Ptr(), idHyperSetLabel, "Set &Label", "wxWidgets website")
		pg.textURL = pg.rowButtonAndText(left, boxL.Ptr(), idHyperSetURL, "Set &URL", "www.wxwidgets.org")

		pg.radioAlign = newRadioBox(boxL.Ptr(), wx.ID_ANY, "alignment",
			[]string{"&left", "&centre", "&right"}, 0, wx.RA_SPECIFY_COLS)
		if pg.radioAlign != nil {
			pg.radioAlign.SetSelection(1)
		}
		addExpandBorder(left, pg.radioAlign.Ptr())

		pg.chkGeneric = wx.NewCheckBox(boxL.Ptr(), wx.ID_ANY, "Use generic version", -1, -1, -1, -1, 0)
		addExpandBorder(left, pg.chkGeneric.Ptr())
	}

	pg.sizerLong = wx.NewBoxSizer(wx.VERTICAL)
	if pg.sizerLong != nil {
		pg.longLink = wx.NewHyperlinkCtrl(pg.win.Ptr(), wx.ID_ANY, "This is a long hyperlink",
			"www.wxwidgets.org", -1, -1, -1, -1, f.attrs.defaultFlags)
	}

	shortRow := wx.NewBoxSizer(wx.HORIZONTAL)
	if shortRow != nil {
		pg.shortRow = shortRow
		pg.visit = wx.NewStaticText(pg.win.Ptr(), wx.ID_ANY, "Visit ", -1, -1, -1, -1, 0)
		pg.hyperlink = wx.NewHyperlinkCtrl(pg.win.Ptr(), idHyperCtrl, "wxWidgets website",
			"www.wxwidgets.org", -1, -1, -1, -1, f.attrs.defaultFlags)
		pg.fun = wx.NewStaticText(pg.win.Ptr(), wx.ID_ANY, " for fun!", -1, -1, -1, -1, 0)

		shortRow.SetMinSize(150, 0)
		shortRow.AddStretchSpacer(0)
		addWin(shortRow, pg.visit.Ptr(), 0, func(sf *wx.SizerFlags) { sf.Center() })
		addWin(shortRow, pg.hyperlink.Ptr(), 0, func(sf *wx.SizerFlags) { sf.Center() })
		addWin(shortRow, pg.fun.Ptr(), 0, func(sf *wx.SizerFlags) { sf.Center() })
		shortRow.AddStretchSpacer(0)
	}

	if pg.sizerLong != nil {
		pg.sizerLong.AddStretchSpacer(0)
		addSub(pg.sizerLong, shortRow.Ptr(), 1, func(sf *wx.SizerFlags) { sf.Center() })
		pg.sizerLong.AddStretchSpacer(0)
		addWin(pg.sizerLong, pg.longLink.Ptr(), 3, func(sf *wx.SizerFlags) { sf.Expand() })
		pg.sizerLong.AddStretchSpacer(0)
	}

	if pg.radioAlign != nil {
		pg.radioAlign.Bind(wx.EVT_RADIOBOX, func(evt *wx.Event) { pg.createLong(f) })
	}
	pg.chkGeneric.Bind(wx.EVT_CHECKBOX, func(evt *wx.Event) {
		pg.createHyperlink(f)
		pg.createLong(f)
	})
	pg.win.BindId(wx.EVT_BUTTON, idHyperSetLabel, func(evt *wx.Event) {
		if pg.hyperlink != nil && pg.textLabel != nil {
			pg.hyperlink.SetLabel(pg.textLabel.GetValue())
		}
	})
	pg.win.BindId(wx.EVT_BUTTON, idHyperSetURL, func(evt *wx.Event) {
		if pg.hyperlink != nil && pg.textURL != nil {
			pg.hyperlink.SetURL(pg.textURL.GetValue())
		}
	})

	addSub(top, left.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.LEFT)
	})
	addSub(top, pg.sizerLong.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.RIGHT)
	})

	pg.win.SetSizer(top.Ptr(), false)
	pg.reset()
}

func (pg *hyperlinkPage) reset() {
	if pg.textLabel != nil && pg.hyperlink != nil {
		pg.textLabel.SetValue(wx.WindowFromPtr(pg.hyperlink.Ptr()).GetLabel())
	}
	if pg.textURL != nil && pg.hyperlink != nil {
		pg.textURL.SetValue(pg.hyperlink.GetURL())
	}
}

func (pg *hyperlinkPage) createHyperlink(f *widgetsFrame) {
	if pg.hyperlink == nil {
		return
	}
	label := wx.WindowFromPtr(pg.hyperlink.Ptr()).GetLabel()
	url := pg.hyperlink.GetURL()
	old := pg.hyperlink
	pg.hyperlink = wx.NewHyperlinkCtrl(pg.win.Ptr(), idHyperCtrl, label, url, -1, -1, -1, -1,
		f.attrs.defaultFlags|(wx.HL_DEFAULT_STYLE&^wx.BORDER_MASK))
	if pg.hyperlink == nil {
		pg.hyperlink = old
		return
	}
	if pg.shortRow != nil {
		pg.shortRow.ReplaceWindow(old.Ptr(), pg.hyperlink.Ptr(), false)
	}
	destroyWin(old)
	pg.notifyWidgetRecreation(f, pg.hyperlink.Ptr())
	pg.win.Layout()
}

func (pg *hyperlinkPage) createLong(f *widgetsFrame) {
	if pg.sizerLong == nil || pg.longLink == nil {
		return
	}
	align := wx.HL_ALIGN_LEFT
	if pg.radioAlign != nil {
		switch pg.radioAlign.GetSelection() {
		case 1:
			align = wx.HL_ALIGN_CENTRE
		case 2:
			align = wx.HL_ALIGN_RIGHT
		}
	}
	old := pg.longLink
	pg.longLink = wx.NewHyperlinkCtrl(pg.win.Ptr(), wx.ID_ANY, "This is a long hyperlink",
		"www.wxwidgets.org", -1, -1, -1, -1,
		f.attrs.defaultFlags|align|(wx.HL_DEFAULT_STYLE&^(wx.HL_ALIGN_CENTRE|wx.BORDER_MASK)))
	if pg.longLink == nil {
		pg.longLink = old
		return
	}
	pg.sizerLong.ReplaceWindow(old.Ptr(), pg.longLink.Ptr(), false)
	destroyWin(old)
	pg.win.Layout()
}

// --- ColourPicker page (clrpicker.cpp) -------------------------------------

var (
	idClrReset  = wx.ID_HIGHEST
	idClrPicker = wx.ID_HIGHEST + 1
)

type colourPage struct {
	pageBase

	chkText  *wx.CheckBox
	chkLabel *wx.CheckBox
	chkAlpha *wx.CheckBox
	picker   *wx.ColorPickerCtrl
	sizerP   *wx.BoxSizer
}

func init() {
	registerPage("ColourPicker", pickerCtrls|genericCtrls, func(book unsafe.Pointer) *pageBase {
		pg := &colourPage{}
		pb := newPageWithIcon(book, "icons/clrpicker.svg", 32, 32)
		if pb == nil {
			return nil
		}
		pg.pageBase = *pb
		pg.content = pg.create
		pg.getWidgets = func() []unsafe.Pointer {
			if pg.picker == nil {
				return nil
			}
			return []unsafe.Pointer{pg.picker.Ptr()}
		}
		pg.recreate = func(f *widgetsFrame) { pg.createPicker(f) }
		return &pg.pageBase
	})
}

func (pg *colourPage) create(f *widgetsFrame) {
	left := wx.NewBoxSizer(wx.VERTICAL)
	box, style := addStaticBox(pg.win.Ptr(), "&ColourPicker style", wx.VERTICAL)
	if style != nil {
		pg.chkText = pg.createCheckBoxAndAddToSizer(style, "With textctrl", wx.ID_ANY, box.Ptr())
		pg.chkLabel = pg.createCheckBoxAndAddToSizer(style, "With label", wx.ID_ANY, box.Ptr())
		pg.chkAlpha = pg.createCheckBoxAndAddToSizer(style, "With opacity", wx.ID_ANY, box.Ptr())
		addSub(left, style.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})
	}

	reset := wx.NewButton(pg.win.Ptr(), idClrReset, "&Reset", -1, -1, -1, -1, 0)
	addWin(left, reset.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.CenterHorizontal()
		sf.TripleBorder(wx.ALL)
	})
	reset.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
		pg.reset()
		pg.createPicker(f)
	})

	pg.sizerP = wx.NewBoxSizer(wx.VERTICAL)
	if pg.sizerP != nil {
		pg.sizerP.AddStretchSpacer(0)
	}

	pg.reset()
	pg.createPicker(f)

	sz := wx.NewBoxSizer(wx.HORIZONTAL)
	if sz == nil {
		return
	}
	addSub(sz, left.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})
	addSub(sz, pg.sizerP.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})
	pg.win.SetSizer(sz.Ptr(), false)

	bindStyleCheckboxes(pg.win, func() { pg.createPicker(f) })
}

func (pg *colourPage) reset() {
	def := wx.CLRP_DEFAULT_STYLE
	setChk(pg.chkText, boolToInt(def&wx.CLRP_USE_TEXTCTRL != 0))
	setChk(pg.chkLabel, boolToInt(def&wx.CLRP_SHOW_LABEL != 0))
	setChk(pg.chkAlpha, boolToInt(def&wx.CLRP_SHOW_ALPHA != 0))
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (pg *colourPage) createPicker(f *widgetsFrame) {
	if pg.sizerP == nil {
		return
	}
	clearSizer(pg.sizerP)
	if pg.picker != nil {
		pg.picker.Destroy()
		pg.picker = nil
	}

	flags := f.attrs.defaultFlags
	if pg.chkText != nil && pg.chkText.GetValue() {
		flags |= wx.CLRP_USE_TEXTCTRL
	}
	if pg.chkLabel != nil && pg.chkLabel.GetValue() {
		flags |= wx.CLRP_SHOW_LABEL
	}
	if pg.chkAlpha != nil && pg.chkAlpha.GetValue() {
		flags |= wx.CLRP_SHOW_ALPHA
	}
	red := wx.NewColourRGB(255, 0, 0, 255)
	pg.picker = wx.NewColorPickerCtrl(pg.win.Ptr(), idClrPicker, red.Ptr(), -1, -1, -1, -1, flags)
	if pg.picker == nil {
		return
	}
	pg.picker.Bind(wx.EVT_COLOURPICKER_CHANGED, func(evt *wx.Event) {
		col := wx.NewColourEmpty()
		pg.picker.GetColour(col.Ptr())
		f.logf("The colour changed to '%s' !", colourString(col))
	})
	addWin(pg.sizerP, pg.picker.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Center()
		sf.Border(wx.ALL)
	})
	pg.sizerP.AddStretchSpacer(0)
	pg.sizerP.Layout()
	pg.notifyWidgetRecreation(f, pg.picker.Ptr())
}

func colourString(c *wx.Colour) string {
	if c == nil {
		return ""
	}
	return c.GetAsString(wx.C2S_CSS_SYNTAX)
}

// --- shared row helpers ----------------------------------------------------

// lastRow is set by rowButtonAndText/rowStaticTextAndText to the row it built.
// It exists so the callers can append a second text control to that row.
var lastRowFor = map[*wx.TextCtrl]*wx.BoxSizer{}

// rowStaticTextAndText adds a "label + text" row to the sizer.
func (p *pageBase) rowStaticTextAndText(sizer sizerTarget, parent unsafe.Pointer, label, value string) *wx.TextCtrl {
	st := wx.NewStaticText(parent, wx.ID_ANY, label, -1, -1, -1, -1, 0)
	text := wx.NewTextCtrl(parent, wx.ID_ANY, "", -1, -1, -1, -1, wx.TE_PROCESS_ENTER)
	row := wx.NewBoxSizer(wx.HORIZONTAL)
	if st == nil || text == nil || row == nil {
		return nil
	}
	if value != "" {
		text.SetValue(value)
	}
	addWin(row, st.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.CenterVertical()
		sf.Border(wx.RIGHT)
	})
	addWin(row, text.Ptr(), 1, func(sf *wx.SizerFlags) { sf.CenterVertical() })
	addSub(sizer, row.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})
	lastRowFor[text] = row
	p.lastRow = row
	return text
}

// rowButtonAndText adds a "button + text" row to the sizer.
func (p *pageBase) rowButtonAndText(sizer sizerTarget, parent unsafe.Pointer, idBtn int, label, value string) *wx.TextCtrl {
	btn := wx.NewButton(parent, idBtn, label, -1, -1, -1, -1, 0)
	text := wx.NewTextCtrl(parent, wx.ID_ANY, "", -1, -1, -1, -1, wx.TE_PROCESS_ENTER)
	row := wx.NewBoxSizer(wx.HORIZONTAL)
	if btn == nil || text == nil || row == nil {
		return nil
	}
	if value != "" {
		text.SetValue(value)
	}
	addWin(row, btn.Ptr(), 0, func(sf *wx.SizerFlags) { sf.Border(wx.RIGHT) })
	addWin(row, text.Ptr(), 1, func(sf *wx.SizerFlags) { sf.Border(wx.LEFT) })
	addSub(sizer, row.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})
	lastRowFor[text] = row
	p.lastRow = row
	return text
}

// addToLastRow appends a text control to the row holding anchor.
func (p *pageBase) addToLastRow(anchor *wx.TextCtrl, text *wx.TextCtrl) {
	row := lastRowFor[anchor]
	if row == nil || text == nil {
		return
	}
	addWin(row, text.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.CenterVertical()
		sf.Border(wx.LEFT)
	})
}
