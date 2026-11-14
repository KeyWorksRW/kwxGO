// Pages demonstrating the picker controls and the editable list box:
// datepick.cpp, timepick.cpp, filepicker.cpp, dirpicker.cpp and editlbox.cpp.
package main

import (
	"strconv"
	"unsafe"

	"github.com/KeyWorksRW/kwxGO/wx"
)

// addSubExpandBorder adds a nested sizer with Expand().Border().
func addSubExpandBorder(s sizerTarget, child unsafe.Pointer) {
	addSub(s, child, 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})
}

// --- DatePicker page (datepick.cpp) ----------------------------------------

var (
	idDpReset    = wx.ID_HIGHEST
	idDpSet      = wx.ID_HIGHEST + 1
	idDpSetRange = wx.ID_HIGHEST + 2
	idDpSetNull  = wx.ID_HIGHEST + 3
	idDpPicker   = wx.ID_HIGHEST + 4
)

type datePickPage struct {
	pageBase

	picker      *wx.DatePickerCtrl
	sizerPicker *wx.BoxSizer

	textCur  *wx.TextCtrl
	textMin  *wx.TextCtrl
	textMax  *wx.TextCtrl
	textNull *wx.TextCtrl

	radioKind  *wx.RadioBox
	chkCentury *wx.CheckBox
	chkNone    *wx.CheckBox
}

func init() {
	registerPage("DatePicker", nativeCtrls|pickerCtrls, func(book unsafe.Pointer) *pageBase {
		pg := &datePickPage{}
		pb := newPageWithIcon(book, "icons/datepick.svg", 32, 32)
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

func (pg *datePickPage) create(f *widgetsFrame) {
	top := wx.NewBoxSizer(wx.HORIZONTAL)
	if top == nil {
		return
	}

	left := wx.NewBoxSizer(wx.VERTICAL)
	if left == nil {
		return
	}
	pg.radioKind = newRadioBox(pg.win.Ptr(), wx.ID_ANY, "&Kind",
		[]string{"&Default", "&Spin", "Drop do&wn"}, 1, wx.RA_SPECIFY_COLS)
	addExpandBorder(left, pg.radioKind.Ptr())

	boxS, style := addStaticBox(pg.win.Ptr(), "&Style", wx.VERTICAL)
	if style != nil {
		pg.chkCentury = pg.createCheckBoxAndAddToSizer(style, "Show &century", wx.ID_ANY, boxS.Ptr())
		pg.chkNone = pg.createCheckBoxAndAddToSizer(style, "Allow &no value", wx.ID_ANY, boxS.Ptr())
		addSubExpandBorder(left, style.Ptr())
	}

	recreate := wx.NewButton(pg.win.Ptr(), idDpReset, "&Recreate", -1, -1, -1, -1, 0)
	addWin(left, recreate.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Center()
		sf.Border(wx.ALL)
	})
	recreate.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
		pg.reset()
		pg.createPicker(f)
	})

	middle := wx.NewBoxSizer(wx.VERTICAL)
	if middle == nil {
		return
	}
	row := pg.createSizerWithTextAndButton(idDpSet, "&Set date", wx.ID_ANY, &pg.textCur, nil)
	addSubExpandBorder(middle, row.Ptr())
	addSpacerSizer(middle, 10)
	row = pg.createSizerWithTextAndLabel("&Min date", wx.ID_ANY, &pg.textMin, nil)
	addSubExpandBorder(middle, row.Ptr())
	row = pg.createSizerWithTextAndLabel("Ma&x date", wx.ID_ANY, &pg.textMax, nil)
	addSubExpandBorder(middle, row.Ptr())
	setRange := wx.NewButton(pg.win.Ptr(), idDpSetRange, "Set &range", -1, -1, -1, -1, 0)
	addWin(middle, setRange.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Center()
		sf.Border(wx.ALL)
	})
	addSpacerSizer(middle, 10)
	row = pg.createSizerWithTextAndLabel("&Null text", wx.ID_ANY, &pg.textNull, nil)
	addSubExpandBorder(middle, row.Ptr())
	setNull := wx.NewButton(pg.win.Ptr(), idDpSetNull, "Set &null text", -1, -1, -1, -1, 0)
	addWin(middle, setNull.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Center()
		sf.Border(wx.ALL)
	})

	pg.sizerPicker = wx.NewBoxSizer(wx.HORIZONTAL)
	if pg.sizerPicker == nil {
		return
	}
	// the constructor dereferences its wxDateTime argument, so never pass nil
	initial := wx.NewDateTime()
	pg.picker = wx.NewDatePickerCtrl(pg.win.Ptr(), idDpPicker, initial.Ptr(), -1, -1, -1, -1, 0)
	pg.sizerPicker.AddStretchSpacer(0)
	addWin(pg.sizerPicker, pg.picker.Ptr(), 1, func(sf *wx.SizerFlags) { sf.Center() })
	pg.sizerPicker.AddStretchSpacer(0)

	addSub(top, left.Ptr(), 0, func(sf *wx.SizerFlags) { sf.DoubleBorder(wx.ALL & ^wx.LEFT) })
	addSub(top, middle.Ptr(), 0, func(sf *wx.SizerFlags) { sf.DoubleBorder(wx.TOP | wx.BOTTOM) })
	addSub(top, pg.sizerPicker.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.RIGHT)
	})

	setChk(pg.chkCentury, 1)
	pg.reset()

	pg.win.BindId(wx.EVT_BUTTON, idDpSet, func(evt *wx.Event) { pg.setDate() })
	pg.win.BindId(wx.EVT_BUTTON, idDpSetRange, func(evt *wx.Event) { pg.setRange(f) })
	pg.win.BindId(wx.EVT_BUTTON, idDpSetNull, func(evt *wx.Event) {
		if pg.picker != nil && pg.textNull != nil {
			pg.picker.SetNullText(pg.textNull.GetValue())
		}
	})
	pg.win.Bind(wx.EVT_DATE_CHANGED, func(evt *wx.Event) {
		if pg.picker == nil {
			return
		}
		dt := wx.NewDateTime()
		pg.picker.GetValue(dt.Ptr())
		f.logf("Date changed, now is %s (control value is %s).",
			dateString(dt), dateString(dt))
	})
	pg.win.Bind(wx.EVT_RADIOBOX, func(evt *wx.Event) { pg.createPicker(f) })
	pg.win.Bind(wx.EVT_CHECKBOX, func(evt *wx.Event) { pg.createPicker(f) })

	pg.win.SetSizer(top.Ptr(), false)
}

func dateString(dt *wx.DateTime) string {
	if dt == nil || !dt.IsValid() {
		return "[none]"
	}
	return dt.FormatISODate()
}

func (pg *datePickPage) reset() {
	today := wx.NewDateTime()
	today.SetToCurrent()
	if pg.picker != nil {
		pg.picker.SetValue(today.Ptr())
	}
	if pg.textCur != nil {
		pg.textCur.SetValue(today.FormatISODate())
	}
}

// createPicker is DatePickerWidgetsPage::CreateDatePicker().
func (pg *datePickPage) createPicker(f *widgetsFrame) {
	if pg.sizerPicker == nil {
		return
	}

	value := wx.NewDateTime()
	old := pg.picker
	if pg.picker != nil {
		pg.picker.GetValue(value.Ptr())
	}
	clearSizer(pg.sizerPicker)
	destroyWin(old)
	pg.picker = nil

	style := f.attrs.defaultFlags
	switch selection(pg.radioKind) {
	case 1:
		style = wx.DP_SPIN
	case 2:
		style = wx.DP_DROPDOWN
	default:
		style = wx.DP_DEFAULT
	}
	if checked(pg.chkCentury) {
		style |= wx.DP_SHOWCENTURY
	}
	if checked(pg.chkNone) {
		style |= wx.DP_ALLOWNONE
	}

	pg.picker = wx.NewDatePickerCtrl(pg.win.Ptr(), idDpPicker, value.Ptr(), -1, -1, -1, -1, style)
	if pg.picker == nil {
		return
	}
	pg.notifyWidgetRecreation(f, pg.picker.Ptr())
	pg.sizerPicker.AddStretchSpacer(0)
	addWin(pg.sizerPicker, pg.picker.Ptr(), 1, func(sf *wx.SizerFlags) { sf.Center() })
	pg.sizerPicker.AddStretchSpacer(0)
	pg.sizerPicker.Layout()
}

func (pg *datePickPage) setDate() {
	if pg.picker == nil || pg.textCur == nil {
		return
	}
	text := pg.textCur.GetValue()
	if text == "" {
		return
	}
	dt := wx.NewDateTime()
	ws := wx.NewWxString(text)
	dt.ParseDate(ws.Ptr())
	ws.Free()
	pg.picker.SetValue(dt.Ptr())
}

func (pg *datePickPage) setRange(f *widgetsFrame) {
	if pg.picker == nil {
		return
	}
	dt1, dt2 := wx.NewDateTime(), wx.NewDateTime()
	dt1.SetToCurrent()
	dt2.SetToCurrent()
	if pg.textMin != nil && pg.textMin.GetValue() != "" {
		ws := wx.NewWxString(pg.textMin.GetValue())
		dt1.ParseDate(ws.Ptr())
		ws.Free()
	}
	if pg.textMax != nil && pg.textMax.GetValue() != "" {
		ws := wx.NewWxString(pg.textMax.GetValue())
		dt2.ParseDate(ws.Ptr())
		ws.Free()
	}
	pg.picker.SetRange(dt1.Ptr(), dt2.Ptr())

	out1, out2 := wx.NewDateTime(), wx.NewDateTime()
	if !pg.picker.GetRange(out1.Ptr(), out2.Ptr()) {
		f.logf("No range set")
		return
	}
	if pg.textMin != nil {
		pg.textMin.SetValue(dateString(out1))
	}
	if pg.textMax != nil {
		pg.textMax.SetValue(dateString(out2))
	}
	f.logf("Date picker range updated")
}

// --- TimePicker page (timepick.cpp) ----------------------------------------

var (
	idTpReset  = wx.ID_HIGHEST
	idTpSet    = wx.ID_HIGHEST + 1
	idTpPicker = wx.ID_HIGHEST + 2
)

type timePickPage struct {
	pageBase

	picker      *wx.TimePickerCtrl
	sizerPicker *wx.BoxSizer
	textCur     *wx.TextCtrl
}

func init() {
	registerPage("TimePicker", nativeCtrls|pickerCtrls, func(book unsafe.Pointer) *pageBase {
		pg := &timePickPage{}
		pb := newPageWithIcon(book, "icons/timepick.svg", 32, 32)
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

func (pg *timePickPage) create(f *widgetsFrame) {
	top := wx.NewBoxSizer(wx.HORIZONTAL)
	if top == nil {
		return
	}

	left := wx.NewBoxSizer(wx.VERTICAL)
	if left == nil {
		return
	}
	reset := wx.NewButton(pg.win.Ptr(), idTpReset, "&Reset", -1, -1, -1, -1, 0)
	addWin(left, reset.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Center()
		sf.Border(wx.ALL)
	})
	reset.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
		pg.reset()
		pg.createPicker(f)
	})

	middle := wx.NewBoxSizer(wx.VERTICAL)
	if middle == nil {
		return
	}
	row := pg.createSizerWithTextAndButton(idTpSet, "&Set time", wx.ID_ANY, &pg.textCur, nil)
	addSubExpandBorder(middle, row.Ptr())

	pg.sizerPicker = wx.NewBoxSizer(wx.HORIZONTAL)
	if pg.sizerPicker == nil {
		return
	}
	// the constructor dereferences its wxDateTime argument, so never pass nil
	now := wx.NewDateTime()
	now.SetToCurrent()
	pg.picker = wx.NewTimePickerCtrl(pg.win.Ptr(), idTpPicker, now.Ptr(), -1, -1, -1, -1, 0)
	pg.sizerPicker.AddStretchSpacer(0)
	addWin(pg.sizerPicker, pg.picker.Ptr(), 1, func(sf *wx.SizerFlags) { sf.Center() })
	pg.sizerPicker.AddStretchSpacer(0)

	addSub(top, left.Ptr(), 0, func(sf *wx.SizerFlags) { sf.DoubleBorder(wx.ALL & ^wx.LEFT) })
	addSub(top, middle.Ptr(), 0, func(sf *wx.SizerFlags) { sf.DoubleBorder(wx.TOP | wx.BOTTOM) })
	addSub(top, pg.sizerPicker.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.RIGHT)
	})

	pg.reset()
	pg.win.SetSizer(top.Ptr(), false)

	pg.win.BindId(wx.EVT_BUTTON, idTpSet, func(evt *wx.Event) { pg.setTime(f) })
	pg.win.Bind(wx.EVT_TIME_CHANGED, func(evt *wx.Event) {
		if pg.picker == nil {
			return
		}
		var h, m, s int
		pg.picker.GetTime(unsafe.Pointer(&h), unsafe.Pointer(&m), unsafe.Pointer(&s))
		f.logf("Time changed, now is %02d:%02d:%02d (control value is %02d:%02d:%02d).", h, m, s, h, m, s)
	})
}

func (pg *timePickPage) reset() {
	today := wx.NewDateTime()
	today.SetToCurrent()
	if pg.picker != nil {
		pg.picker.SetValue(today.Ptr())
	}
	if pg.textCur != nil {
		pg.textCur.SetValue(today.FormatISOTime())
	}
}

// createPicker is TimePickerWidgetsPage::CreateTimePicker().
func (pg *timePickPage) createPicker(f *widgetsFrame) {
	if pg.sizerPicker == nil {
		return
	}
	value := wx.NewDateTime()
	old := pg.picker
	if pg.picker != nil {
		pg.picker.GetValue(value.Ptr())
	}
	clearSizer(pg.sizerPicker)
	destroyWin(old)
	pg.picker = nil

	pg.picker = wx.NewTimePickerCtrl(pg.win.Ptr(), idTpPicker, value.Ptr(), -1, -1, -1, -1,
		f.attrs.defaultFlags)
	if pg.picker == nil {
		return
	}
	pg.notifyWidgetRecreation(f, pg.picker.Ptr())
	pg.sizerPicker.AddStretchSpacer(0)
	addWin(pg.sizerPicker, pg.picker.Ptr(), 1, func(sf *wx.SizerFlags) { sf.Center() })
	pg.sizerPicker.AddStretchSpacer(0)
	pg.sizerPicker.Layout()
}

func (pg *timePickPage) setTime(f *widgetsFrame) {
	if pg.picker == nil || pg.textCur == nil {
		return
	}
	text := pg.textCur.GetValue()
	var h, m, s int
	if n, err := sscanfColon(text, &h, &m, &s); err != nil || n != 3 {
		f.logf("Invalid time, please use HH:MM:SS format.")
		return
	}
	pg.picker.SetTime(h, m, s)
}

// sscanfColon parses "HH:MM:SS" (the format the sample documents).
func sscanfColon(text string, h, m, s *int) (int, error) {
	parts := [3]int{}
	str := text
	n := 0
	start := 0
	for n < 3 {
		end := start
		for end < len(str) && str[end] != ':' {
			end++
		}
		if end == start {
			return n, strconv.ErrSyntax
		}
		v, err := strconv.Atoi(str[start:end])
		if err != nil {
			return n, err
		}
		parts[n] = v
		n++
		if end == len(str) {
			break
		}
		start = end + 1
	}
	*h, *m, *s = parts[0], parts[1], parts[2]
	return n, nil
}

// --- FilePicker page (filepicker.cpp) --------------------------------------

var (
	idFpReset   = wx.ID_HIGHEST
	idFpFile    = wx.ID_HIGHEST + 1
	idFpSetDir  = wx.ID_HIGHEST + 2
	idFpCurPath = wx.ID_HIGHEST + 3
)

type filePickPage struct {
	pageBase

	picker *wx.FilePickerCtrl

	chkTextCtrl  *wx.CheckBox
	chkOverwrite *wx.CheckBox
	chkMustExist *wx.CheckBox
	chkChangeDir *wx.CheckBox
	chkSmall     *wx.CheckBox

	radioMode *wx.RadioBox
	labelPath *wx.StaticText
	textDir   *wx.TextCtrl
	sizer     *wx.BoxSizer
}

func init() {
	registerPage("FilePicker", pickerCtrls|genericCtrls, func(book unsafe.Pointer) *pageBase {
		pg := &filePickPage{}
		pb := newPageWithIcon(book, "icons/filepicker.svg", 32, 32)
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
		pg.recreate = func(f *widgetsFrame) { pg.recreatePicker(f) }
		return &pg.pageBase
	})
}

func (pg *filePickPage) create(f *widgetsFrame) {
	left := wx.NewBoxSizer(wx.VERTICAL)
	if left == nil {
		return
	}

	pg.radioMode = newRadioBox(pg.win.Ptr(), wx.ID_ANY, "wxFilePicker mode",
		[]string{"open", "save"}, 0, wx.RA_SPECIFY_COLS)
	addExpandBorder(left, pg.radioMode.Ptr())

	boxS, style := addStaticBox(pg.win.Ptr(), "&FilePicker style", wx.VERTICAL)
	if style != nil {
		pg.chkTextCtrl = pg.createCheckBoxAndAddToSizer(style, "With textctrl", wx.ID_ANY, boxS.Ptr())
		pg.chkOverwrite = pg.createCheckBoxAndAddToSizer(style, "Overwrite prompt", wx.ID_ANY, boxS.Ptr())
		pg.chkMustExist = pg.createCheckBoxAndAddToSizer(style, "File must exist", wx.ID_ANY, boxS.Ptr())
		pg.chkChangeDir = pg.createCheckBoxAndAddToSizer(style, "Change working dir", wx.ID_ANY, boxS.Ptr())
		pg.chkSmall = pg.createCheckBoxAndAddToSizer(style, "&Small version", wx.ID_ANY, boxS.Ptr())
		addSubExpandBorder(left, style.Ptr())
	}

	row := pg.createSizerWithTextAndButton(idFpSetDir, "&Initial directory", wx.ID_ANY, &pg.textDir, nil)
	addSubExpandBorder(left, row.Ptr())
	addSpacerSizer(left, 10)

	reset := wx.NewButton(pg.win.Ptr(), idFpReset, "&Reset", -1, -1, -1, -1, 0)
	addWin(left, reset.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.CenterHorizontal()
		sf.TripleBorder(wx.ALL)
	})
	reset.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
		pg.reset()
		pg.recreatePicker(f)
	})

	pg.reset()

	pg.sizer = wx.NewBoxSizer(wx.VERTICAL)
	if pg.sizer == nil {
		return
	}
	pg.sizer.AddStretchSpacer(0)
	pg.labelPath = wx.NewStaticText(pg.win.Ptr(), idFpCurPath, "", -1, -1, -1, -1, 0)
	pg.createPicker(f)
	pg.sizer.AddStretchSpacer(0)
	addWin(pg.sizer, pg.labelPath.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})
	pg.sizer.AddStretchSpacer(0)

	sz := wx.NewBoxSizer(wx.HORIZONTAL)
	if sz == nil {
		return
	}
	addSub(sz, left.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})
	addSub(sz, pg.sizer.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})
	pg.win.SetSizer(sz.Ptr(), false)

	pg.win.BindId(wx.EVT_BUTTON, idFpSetDir, func(evt *wx.Event) {
		if pg.picker == nil || pg.textDir == nil {
			return
		}
		dir := pg.textDir.GetValue()
		pg.picker.SetInitialDirectory(dir)
		f.logf("Initial directory set to \"%s\"", dir)
	})
	pg.win.Bind(wx.EVT_FILEPICKER_CHANGED, func(evt *wx.Event) {
		if pg.picker == nil {
			return
		}
		f.logf("The file changed to '%s' (the control has '%s').", pg.picker.GetPath(), pg.picker.GetPath())
	})
	pg.win.Bind(wx.EVT_CHECKBOX, func(evt *wx.Event) { pg.recreatePicker(f) })
	pg.win.Bind(wx.EVT_RADIOBOX, func(evt *wx.Event) {
		pg.updateMode()
		pg.recreatePicker(f)
	})
	pg.win.BindId(wx.EVT_UPDATE_UI, idFpCurPath, func(evt *wx.Event) {
		ui := wx.AsUpdateUIEvent(evt)
		if ui == nil || pg.picker == nil {
			return
		}
		ui.SetText("Current path: " + pg.picker.GetPath())
	})
}

func (pg *filePickPage) reset() {
	def := wx.FLP_DEFAULT_STYLE
	if pg.radioMode != nil {
		if def&wx.FLP_OPEN != 0 {
			pg.radioMode.SetSelection(0)
		} else {
			pg.radioMode.SetSelection(1)
		}
	}
	setChk(pg.chkTextCtrl, boolToInt(def&wx.FLP_USE_TEXTCTRL != 0))
	setChk(pg.chkOverwrite, boolToInt(def&wx.FLP_OVERWRITE_PROMPT != 0))
	setChk(pg.chkMustExist, boolToInt(def&wx.FLP_FILE_MUST_EXIST != 0))
	setChk(pg.chkChangeDir, boolToInt(def&wx.FLP_CHANGE_DIR != 0))
	setChk(pg.chkSmall, boolToInt(def&wx.FLP_SMALL != 0))
	pg.updateMode()
}

func (pg *filePickPage) updateMode() {
	if selection(pg.radioMode) == 0 {
		setChk(pg.chkOverwrite, 0)
		enableWin(pg.chkOverwrite, false)
		enableWin(pg.chkMustExist, true)
	} else {
		setChk(pg.chkMustExist, 0)
		enableWin(pg.chkMustExist, false)
		enableWin(pg.chkOverwrite, true)
	}
}

// createPicker is FilePickerWidgetsPage::CreatePicker().
func (pg *filePickPage) createPicker(f *widgetsFrame) {
	destroyWin(pg.picker)
	pg.picker = nil

	style := f.attrs.defaultFlags
	if checked(pg.chkTextCtrl) {
		style |= wx.FLP_USE_TEXTCTRL
	}
	if checked(pg.chkOverwrite) {
		style |= wx.FLP_OVERWRITE_PROMPT
	}
	if checked(pg.chkMustExist) {
		style |= wx.FLP_FILE_MUST_EXIST
	}
	if checked(pg.chkChangeDir) {
		style |= wx.FLP_CHANGE_DIR
	}
	if checked(pg.chkSmall) {
		style |= wx.FLP_SMALL
	}
	if selection(pg.radioMode) == 0 {
		style |= wx.FLP_OPEN
	} else {
		style |= wx.FLP_SAVE
	}

	pg.picker = wx.NewFilePickerCtrl(pg.win.Ptr(), idFpFile, -1, -1, -1, -1, "", "Hello!", "*", style)
	if pg.picker == nil {
		return
	}
	pg.notifyWidgetRecreation(f, pg.picker.Ptr())
	addWin(pg.sizer, pg.picker.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})
}

// recreatePicker is FilePickerWidgetsPage::RecreatePicker().
func (pg *filePickPage) recreatePicker(f *widgetsFrame) {
	if pg.sizer == nil {
		return
	}
	pg.sizer.Remove(1)
	clearSizerIndex(pg.sizer, 1)
	pg.createPicker(f)
	pg.sizer.Layout()
}

// clearSizerIndex removes the item that createPicker may have appended.
func clearSizerIndex(s *wx.BoxSizer, from int) {
	for s.GetItemCount() > from {
		if !s.Remove(from) {
			return
		}
	}
}

// --- DirPicker page (dirpicker.cpp) ----------------------------------------

var (
	idDpDirReset = wx.ID_HIGHEST
	idDpDir      = wx.ID_HIGHEST + 1
	idDpDirSet   = wx.ID_HIGHEST + 2
)

type dirPickPage struct {
	pageBase

	picker       *wx.DirPickerCtrl
	chkTextCtrl  *wx.CheckBox
	chkChange    *wx.CheckBox
	chkMustExist *wx.CheckBox
	chkSmall     *wx.CheckBox
	textDir      *wx.TextCtrl
	sizer        *wx.BoxSizer
}

func init() {
	registerPage("DirPicker", pickerCtrls|genericCtrls, func(book unsafe.Pointer) *pageBase {
		pg := &dirPickPage{}
		pb := newPageWithIcon(book, "icons/dirpicker.svg", 32, 32)
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
		pg.recreate = func(f *widgetsFrame) { pg.recreatePicker(f) }
		return &pg.pageBase
	})
}

func (pg *dirPickPage) create(f *widgetsFrame) {
	left := wx.NewBoxSizer(wx.VERTICAL)
	if left == nil {
		return
	}

	boxS, style := addStaticBox(pg.win.Ptr(), "&DirPicker style", wx.VERTICAL)
	if style != nil {
		pg.chkTextCtrl = pg.createCheckBoxAndAddToSizer(style, "With textctrl", wx.ID_ANY, boxS.Ptr())
		pg.chkMustExist = pg.createCheckBoxAndAddToSizer(style, "Dir must exist", wx.ID_ANY, boxS.Ptr())
		pg.chkChange = pg.createCheckBoxAndAddToSizer(style, "Change working dir", wx.ID_ANY, boxS.Ptr())
		pg.chkSmall = pg.createCheckBoxAndAddToSizer(style, "&Small version", wx.ID_ANY, boxS.Ptr())
		addSubExpandBorder(left, style.Ptr())
	}

	row := pg.createSizerWithTextAndButton(idDpDirSet, "&Initial directory", wx.ID_ANY, &pg.textDir, nil)
	addSubExpandBorder(left, row.Ptr())
	addSpacerSizer(left, 10)

	reset := wx.NewButton(pg.win.Ptr(), idDpDirReset, "&Reset", -1, -1, -1, -1, 0)
	addWin(left, reset.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.CenterHorizontal()
		sf.TripleBorder(wx.ALL)
	})
	reset.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
		pg.reset()
		pg.recreatePicker(f)
	})

	pg.reset()

	pg.sizer = wx.NewBoxSizer(wx.VERTICAL)
	if pg.sizer == nil {
		return
	}
	pg.sizer.AddStretchSpacer(0)
	pg.createPicker(f)
	pg.sizer.AddStretchSpacer(0)

	sz := wx.NewBoxSizer(wx.HORIZONTAL)
	if sz == nil {
		return
	}
	addSub(sz, left.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})
	addSub(sz, pg.sizer.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})
	pg.win.SetSizer(sz.Ptr(), false)

	pg.win.BindId(wx.EVT_BUTTON, idDpDirSet, func(evt *wx.Event) {
		if pg.picker != nil && pg.textDir != nil {
			pg.picker.SetInitialDirectory(pg.textDir.GetValue())
		}
	})
	pg.win.Bind(wx.EVT_DIRPICKER_CHANGED, func(evt *wx.Event) {
		if pg.picker != nil {
			f.logf("The directory changed to '%s' !", pg.picker.GetPath())
		}
	})
	pg.win.Bind(wx.EVT_CHECKBOX, func(evt *wx.Event) { pg.recreatePicker(f) })
}

func (pg *dirPickPage) reset() {
	def := wx.DIRP_DEFAULT_STYLE
	setChk(pg.chkTextCtrl, boolToInt(def&wx.DIRP_USE_TEXTCTRL != 0))
	setChk(pg.chkMustExist, boolToInt(def&wx.DIRP_DIR_MUST_EXIST != 0))
	setChk(pg.chkChange, boolToInt(def&wx.DIRP_CHANGE_DIR != 0))
	setChk(pg.chkSmall, boolToInt(def&wx.DIRP_SMALL != 0))
}

// createPicker is DirPickerWidgetsPage::CreatePicker().
func (pg *dirPickPage) createPicker(f *widgetsFrame) {
	destroyWin(pg.picker)
	pg.picker = nil
	style := f.attrs.defaultFlags
	if checked(pg.chkTextCtrl) {
		style |= wx.DIRP_USE_TEXTCTRL
	}
	if checked(pg.chkMustExist) {
		style |= wx.DIRP_DIR_MUST_EXIST
	}
	if checked(pg.chkChange) {
		style |= wx.DIRP_CHANGE_DIR
	}
	if checked(pg.chkSmall) {
		style |= wx.DIRP_SMALL
	}

	pg.picker = wx.NewDirPickerCtrl(pg.win.Ptr(), idDpDir, -1, -1, -1, -1, "", "Hello!", style)
	if pg.picker == nil {
		return
	}
	pg.notifyWidgetRecreation(f, pg.picker.Ptr())
	addWin(pg.sizer, pg.picker.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})
}

// recreatePicker is DirPickerWidgetsPage::RecreatePicker().
func (pg *dirPickPage) recreatePicker(f *widgetsFrame) {
	if pg.sizer == nil {
		return
	}
	pg.sizer.Remove(1)
	clearSizerIndex(pg.sizer, 1)
	pg.createPicker(f)
	pg.sizer.Layout()
}

// --- EditableListbox page (editlbox.cpp) -----------------------------------

var (
	idElReset   = wx.ID_HIGHEST
	idElListbox = wx.ID_HIGHEST + 1
)

type editableLboxPage struct {
	pageBase

	chkAllowNew    *wx.CheckBox
	chkAllowEdit   *wx.CheckBox
	chkAllowDelete *wx.CheckBox
	chkNoReorder   *wx.CheckBox
	lbox           *wx.EditableListBox
	sizerLbox      *wx.BoxSizer
}

func init() {
	registerPage("EditableListbox", genericCtrls, func(book unsafe.Pointer) *pageBase {
		pg := &editableLboxPage{}
		pb := newPageWithIcon(book, "icons/listbox.svg", 32, 32)
		if pb == nil {
			return nil
		}
		pg.pageBase = *pb
		pg.content = pg.create
		pg.getWidgets = func() []unsafe.Pointer {
			if pg.lbox == nil {
				return nil
			}
			return []unsafe.Pointer{pg.lbox.GetListCtrl()}
		}
		pg.recreate = func(f *widgetsFrame) { pg.createLbox(f) }
		return &pg.pageBase
	})
}

func (pg *editableLboxPage) create(f *widgetsFrame) {
	top := wx.NewBoxSizer(wx.HORIZONTAL)
	if top == nil {
		return
	}

	boxL, left := addStaticBox(pg.win.Ptr(), "&Set listbox parameters", wx.VERTICAL)
	if left != nil {
		pg.chkAllowNew = pg.createCheckBoxAndAddToSizer(left, "Allow new items", wx.ID_ANY, boxL.Ptr())
		pg.chkAllowEdit = pg.createCheckBoxAndAddToSizer(left, "Allow editing items", wx.ID_ANY, boxL.Ptr())
		pg.chkAllowDelete = pg.createCheckBoxAndAddToSizer(left, "Allow deleting items", wx.ID_ANY, boxL.Ptr())
		pg.chkNoReorder = pg.createCheckBoxAndAddToSizer(left, "Block user reordering", wx.ID_ANY, boxL.Ptr())

		reset := wx.NewButton(boxL.Ptr(), idElReset, "&Reset", -1, -1, -1, -1, 0)
		addWin(left, reset.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.CenterHorizontal()
			sf.TripleBorder(wx.ALL)
		})
		reset.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			pg.reset()
			pg.createLbox(f)
		})
	}

	pg.sizerLbox = wx.NewBoxSizer(wx.VERTICAL)
	if pg.sizerLbox == nil {
		return
	}
	pg.sizerLbox.SetMinSize(150, 0)

	pg.lbox = wx.NewEditableListBox(pg.win.Ptr(), idElListbox, "Match these wildcards:", -1, -1, -1, -1, 0)
	if pg.lbox != nil {
		addWin(pg.sizerLbox, pg.lbox.Ptr(), 1, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})
	}

	pg.reset()

	addSub(top, left.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.LEFT)
	})
	addSub(top, pg.sizerLbox.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.RIGHT)
	})
	pg.win.SetSizer(top.Ptr(), false)

	pg.win.Bind(wx.EVT_CHECKBOX, func(evt *wx.Event) { pg.createLbox(f) })
}

func (pg *editableLboxPage) reset() {
	setChk(pg.chkAllowNew, 0)
	setChk(pg.chkAllowEdit, 0)
	setChk(pg.chkAllowDelete, 0)
	setChk(pg.chkNoReorder, 0)
}

// createLbox is EditableListboxWidgetsPage::CreateLbox().
func (pg *editableLboxPage) createLbox(f *widgetsFrame) {
	if pg.sizerLbox == nil {
		return
	}
	flags := f.attrs.defaultFlags
	if checked(pg.chkAllowNew) {
		flags |= wx.EL_ALLOW_NEW
	}
	if checked(pg.chkAllowEdit) {
		flags |= wx.EL_ALLOW_EDIT
	}
	if checked(pg.chkAllowDelete) {
		flags |= wx.EL_ALLOW_DELETE
	}
	if checked(pg.chkNoReorder) {
		flags |= wx.EL_NO_REORDER
	}

	// carry the current items over to the new control
	var items []string
	old := pg.lbox
	if pg.lbox != nil {
		if lc := wx.ListCtrlFromPtr(pg.lbox.GetListCtrl()); lc != nil {
			for i := 0; i < lc.GetItemCount(); i++ {
				items = append(items, lc.GetItemText(i, 0))
			}
		}
		clearSizer(pg.sizerLbox)
		destroyWin(old)
		pg.lbox = nil
	}

	pg.lbox = wx.NewEditableListBox(pg.win.Ptr(), idElListbox, "Match these wildcards:", -1, -1, -1, -1, flags)
	if pg.lbox == nil {
		return
	}
	pg.notifyWidgetRecreation(f, pg.lbox.Ptr())

	if len(items) > 0 {
		if arr := stringsArray(items); arr != nil {
			pg.lbox.SetStrings(arr.Ptr())
			arr.Delete()
		}
	}
	addWin(pg.sizerLbox, pg.lbox.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})
	pg.sizerLbox.Layout()
}
