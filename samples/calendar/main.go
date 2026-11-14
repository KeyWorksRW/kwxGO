// Command calendar is the kwxGO port of the wxWidgets "calendar" sample
// (samples/calendar/calendar.cpp).
//
// It exercises wxCalendarCtrl / wxGenericCalendarCtrl, wxDatePickerCtrl and
// wxTimePickerCtrl: the frame holds a splitter whose upper pane is a panel with
// a static text and a calendar control, and whose lower pane is a read-only
// text control that receives wxLogMessage() output. The menus toggle the
// calendar styles, the first weekday, holiday and special-date highlighting, a
// restricted date range and the resizable flag, and open the date/time picker
// dialogs.
package main

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/KeyWorksRW/kwxGO/wx"
)

// Command IDs, mirroring the anonymous enum in calendar.cpp. The File menu
// reuses the standard wxWidgets IDs; the calendar, date-picker and time-picker
// groups keep their C++ numeric values because the handlers are bound by ID.
const (
	idCalGeneric       = 200
	idCalAutoWeekday   = 201
	idCalSunday        = 202
	idCalMonday        = 203
	idCalHolidays      = 204
	idCalSpecial       = 205
	idCalMonth         = 206
	idCalLimitDates    = 207
	idCalSeqMonth      = 208
	idCalSurroundWeeks = 209
	idCalWeekNumbers   = 210
	idCalSetDate       = 211
	idCalToday         = 212
	idCalBeginDST      = 213
	idCalResizable     = 214

	idDatePickerAskDate       = 300
	idDatePickerShowCentury   = 301
	idDatePickerDropDown      = 302
	idDatePickerAllowNone     = 303
	idDatePickerStartWithNone = 304

	idTimePickerAskTime = 400

	// idCalendarCtrl is the ID shared by every calendar control the panel
	// creates, so the panel's event bindings keep matching after a recreate.
	idCalendarCtrl = 1000
)

// calendarFrame holds the state MyFrame keeps in member variables.
type calendarFrame struct {
	frame *wx.Frame
	panel *calendarPanel
	log   *wx.TextCtrl
}

func main() {
	wx.Run(func() bool {
		// Locale affects the language used in the calendar, and may affect the
		// first day of the week, so set it before creating the controls.
		wx.UILocale_UseDefault()

		f := newCalendarFrame()
		if f == nil {
			return false
		}
		f.frame.Show()
		wx.SetTopWindow(f.frame.Ptr())
		return true
	})
}

// newCalendarFrame builds the frame, its menu bar, the splitter with the
// calendar panel and the log pane, and installs the log target -- MyFrame's
// constructor.
func newCalendarFrame() *calendarFrame {
	frame := wx.NewFrame(nil, wx.ID_ANY, "Calendar wxWidgets sample",
		-1, -1, -1, -1, wx.DEFAULT_FRAME_STYLE)
	if frame == nil {
		return nil
	}

	f := &calendarFrame{frame: frame}

	if icon := wx.NewIconFromXPM(sampleXPM); icon != nil {
		frame.SetIcon(icon.Ptr())
	}

	// frame->SetSize(frame->FromDIP(wxSize(460, 340)))
	frame.SetSize(-1, -1, frame.FromDIP(460), frame.FromDIP(340), wx.SIZE_USE_EXISTING)

	f.buildMenus()

	splitter := wx.NewSplitterWindow(frame.Ptr(), wx.ID_ANY, -1, -1, -1, -1, wx.SP_NOBORDER)
	if splitter == nil {
		return nil
	}

	panel := newCalendarPanel(splitter.Ptr(), f)
	if panel == nil {
		return nil
	}
	f.panel = panel

	f.log = wx.NewTextCtrl(splitter.Ptr(), wx.ID_ANY, "",
		-1, -1, -1, -1, wx.TE_READONLY|wx.TE_MULTILINE)
	if f.log == nil {
		return nil
	}

	splitter.SplitHorizontally(panel.win.Ptr(), f.log.Ptr(), 0)
	splitter.SetMinimumPaneSize(20)

	if lt := wx.NewLogTextCtrl(f.log.Ptr()); lt != nil {
		wx.SetActiveLogTarget(lt.Ptr())
	}

	f.bindEvents()
	return f
}

// buildMenus recreates the File / Calendar / Date picker / Time picker menus
// verbatim. The generic date- and time-picker menu items from the C++ sample
// are compiled out on Windows (wxUSE_DATEPICKCTRL_GENERIC and
// wxUSE_TIMEPICKCTRL_GENERIC are 0), so they are not ported.
func (f *calendarFrame) buildMenus() {
	menuFile := wx.NewMenu("", 0)
	if menuFile == nil {
		return
	}
	menuFile.Append(wx.ID_ABOUT, "&About\tCtrl-A", "Show about dialog", false)
	menuFile.AppendSeparator()
	menuFile.Append(wx.ID_CLEAR, "&Clear log\tCtrl-L", "", false)
	menuFile.AppendSeparator()
	menuFile.Append(wx.ID_EXIT, "E&xit\tAlt-X", "Quit this program", false)

	menuCal := wx.NewMenu("", 0)
	if menuCal == nil {
		return
	}
	menuCal.AppendCheckItem(idCalGeneric, "Use &generic version\tCtrl-G",
		"Toggle between native and generic control")
	menuCal.AppendSeparator()
	menuCal.AppendRadioItem(idCalAutoWeekday, "Automatic &first weekday\tCtrl-V", "")
	menuCal.AppendRadioItem(idCalSunday, "Sunday &first weekday\tCtrl-Z", "")
	menuCal.AppendRadioItem(idCalMonday, "Monday &first weekday\tCtrl-F", "")
	menuCal.Append(idCalHolidays, "Show &holidays\tCtrl-H",
		"Toggle highlighting the holidays", true)
	menuCal.Append(idCalSpecial, "Highlight &special dates\tCtrl-S",
		"Test custom highlighting", true)
	menuCal.Append(idCalSurroundWeeks, "Show s&urrounding weeks\tCtrl-W",
		"Show the neighbouring weeks in the prev/next month", true)
	menuCal.Append(idCalWeekNumbers, "Show &week numbers",
		"Toggle week numbers", true)
	menuCal.AppendSeparator()
	menuCal.Append(idCalSeqMonth, "Toggle month selector st&yle\tCtrl-Y",
		"Use another style for the calendar controls", true)
	menuCal.Append(idCalMonth, "&Month can be changed\tCtrl-M",
		"Allow changing the month in the calendar", true)
	menuCal.AppendCheckItem(idCalLimitDates, "Toggle date ra&nge\tCtrl-N",
		"Limit the valid dates")
	menuCal.AppendSeparator()
	menuCal.Append(idCalSetDate, "Call &SetDate(2005-12-24)", "Set date to 2005-12-24.", false)
	menuCal.Append(idCalToday, "Call &Today()", "Set to the current date.", false)
	menuCal.Append(idCalBeginDST, "Call SetDate(GetBeginDST())", "", false)
	menuCal.AppendSeparator()
	menuCal.AppendCheckItem(idCalResizable, "Make &resizable\tCtrl-R", "")

	menuDate := wx.NewMenu("", 0)
	if menuDate == nil {
		return
	}
	menuDate.AppendCheckItem(idDatePickerShowCentury, "Al&ways show century", "")
	menuDate.AppendCheckItem(idDatePickerDropDown, "Use &drop down control", "")
	menuDate.AppendCheckItem(idDatePickerAllowNone, "Allow &no date", "")
	menuDate.AppendCheckItem(idDatePickerStartWithNone, "Start &with no date", "")
	menuDate.AppendSeparator()
	menuDate.Append(idDatePickerAskDate, "&Choose date...\tCtrl-D",
		"Show dialog with wxDatePickerCtrl", false)

	menuTime := wx.NewMenu("", 0)
	if menuTime == nil {
		return
	}
	menuTime.Append(idTimePickerAskTime, "&Choose time...\tCtrl-T",
		"Show dialog with wxTimePickerCtrl", false)

	menuBar := wx.NewMenuBar(0)
	if menuBar == nil {
		return
	}
	menuBar.Append(menuFile.Ptr(), "&File")
	menuBar.Append(menuCal.Ptr(), "&Calendar")
	menuBar.Append(menuDate.Ptr(), "&Date picker")
	menuBar.Append(menuTime.Ptr(), "&Time picker")
	f.frame.SetMenuBar(menuBar.Ptr())

	menuBar.Check(idCalAutoWeekday, true)
	menuBar.Check(idCalSunday, false)
	menuBar.Check(idCalMonday, false)
	menuBar.Check(idCalHolidays, true)
	menuBar.Check(idCalMonth, true)
	menuBar.Check(idCalLimitDates, false)
	menuBar.Check(idDatePickerShowCentury, true)
}

// bindEvents wires the C++ event table onto the frame.
func (f *calendarFrame) bindEvents() {
	// File menu
	f.frame.BindId(wx.EVT_MENU, wx.ID_ABOUT, f.onAbout)
	f.frame.BindId(wx.EVT_MENU, wx.ID_CLEAR, f.onClearLog)
	f.frame.BindId(wx.EVT_MENU, wx.ID_EXIT, f.onQuit)

	// Date picker menu
	f.frame.BindId(wx.EVT_MENU, idDatePickerAskDate, f.onAskDate)
	f.frame.BindId(wx.EVT_UPDATE_UI, idDatePickerStartWithNone, f.onUpdateUIStartWithNone)

	// Time picker menu
	f.frame.BindId(wx.EVT_MENU, idTimePickerAskTime, f.onAskTime)

	// Calendar menu
	f.frame.BindId(wx.EVT_MENU, idCalGeneric, f.onCalGeneric)
	f.frame.BindId(wx.EVT_MENU, idCalAutoWeekday, f.onCalAutoWeekday)
	f.frame.BindId(wx.EVT_MENU, idCalSunday, f.onCalSunday)
	f.frame.BindId(wx.EVT_MENU, idCalMonday, f.onCalMonday)
	f.frame.BindId(wx.EVT_MENU, idCalHolidays, f.onCalHolidays)
	f.frame.BindId(wx.EVT_MENU, idCalSpecial, f.onCalSpecial)
	f.frame.BindId(wx.EVT_MENU, idCalMonth, f.onCalAllowMonth)
	f.frame.BindId(wx.EVT_MENU, idCalLimitDates, f.onCalLimitDates)
	f.frame.BindId(wx.EVT_MENU, idCalSeqMonth, f.onCalSeqMonth)
	f.frame.BindId(wx.EVT_MENU, idCalSurroundWeeks, f.onCalShowSurroundingWeeks)
	f.frame.BindId(wx.EVT_MENU, idCalWeekNumbers, f.onCalShowWeekNumbers)
	f.frame.BindId(wx.EVT_MENU, idCalSetDate, f.onSetDate)
	f.frame.BindId(wx.EVT_MENU, idCalToday, f.onToday)
	f.frame.BindId(wx.EVT_MENU, idCalBeginDST, f.onBeginDST)
	f.frame.BindId(wx.EVT_MENU, idCalResizable, f.onCalToggleResizable)

	// EVT_UPDATE_UI(Calendar_Cal_SeqMonth, OnUpdateUIGenericOnly) and, on this
	// platform, the same for Calendar_Cal_Special and
	// Calendar_Cal_SurroundWeeks.
	f.frame.BindId(wx.EVT_UPDATE_UI, idCalSeqMonth, f.onUpdateUIGenericOnly)
	f.frame.BindId(wx.EVT_UPDATE_UI, idCalSpecial, f.onUpdateUIGenericOnly)
	f.frame.BindId(wx.EVT_UPDATE_UI, idCalSurroundWeeks, f.onUpdateUIGenericOnly)

	// The C++ event table additionally disables the first-weekday radio items
	// under #ifdef __WXGTK__ while the native control is in use. Mirror that
	// branch; it is inactive on the MSW build this port targets.
	if runtime.GOOS == "linux" {
		f.frame.BindId(wx.EVT_UPDATE_UI, idCalAutoWeekday, f.onUpdateUIGenericOnly)
		f.frame.BindId(wx.EVT_UPDATE_UI, idCalSunday, f.onUpdateUIGenericOnly)
		f.frame.BindId(wx.EVT_UPDATE_UI, idCalMonday, f.onUpdateUIGenericOnly)
		f.frame.BindId(wx.EVT_UPDATE_UI, idCalHolidays, f.onUpdateUIGenericOnly)
	}
}

// ---------------------------------------------------------------------------
// Event handlers
// ---------------------------------------------------------------------------

func (f *calendarFrame) onQuit(evt *wx.Event) {
	if evt == nil {
		return
	}
	// true is to force the frame to close
	f.frame.Close(true)
}

func (f *calendarFrame) onAbout(evt *wx.Event) {
	if evt == nil {
		return
	}
	wx.MessageBox("wxCalendarCtrl sample\n(c) 2000--2008 Vadim Zeitlin",
		"About Calendar", wx.OK|wx.ICON_INFORMATION, f.frame.Ptr(), -1, -1)
}

func (f *calendarFrame) onClearLog(evt *wx.Event) {
	if evt == nil || f.log == nil {
		return
	}
	f.log.Clear()
}

func (f *calendarFrame) onCalGeneric(evt *wx.Event) {
	if evt == nil {
		return
	}
	f.panel.toggleUseGeneric()
}

func (f *calendarFrame) onCalAutoWeekday(evt *wx.Event) {
	if evt == nil {
		return
	}
	f.panel.toggleCalStyle(false, wx.CAL_SUNDAY_FIRST)
	f.panel.toggleCalStyle(false, wx.CAL_MONDAY_FIRST)
}

func (f *calendarFrame) onCalSunday(evt *wx.Event) {
	if evt == nil {
		return
	}
	cmd := wx.AsCommandEvent(evt)
	if cmd == nil {
		return
	}
	f.panel.toggleCalStyle(false, wx.CAL_MONDAY_FIRST)
	f.panel.toggleCalStyle(cmd.IsChecked(), wx.CAL_SUNDAY_FIRST)
}

func (f *calendarFrame) onCalMonday(evt *wx.Event) {
	if evt == nil {
		return
	}
	cmd := wx.AsCommandEvent(evt)
	if cmd == nil {
		return
	}
	f.panel.toggleCalStyle(false, wx.CAL_SUNDAY_FIRST)
	f.panel.toggleCalStyle(cmd.IsChecked(), wx.CAL_MONDAY_FIRST)
}

func (f *calendarFrame) onCalHolidays(evt *wx.Event) {
	if evt == nil {
		return
	}
	cmd := wx.AsCommandEvent(evt)
	if cmd == nil {
		return
	}
	display := 0
	if cmd.IsChecked() {
		display = 1
	}
	f.panel.cal.EnableHolidayDisplay(display)
}

func (f *calendarFrame) onCalSpecial(evt *wx.Event) {
	if evt == nil {
		return
	}
	f.panel.highlightSpecial(f.menuChecked(evt.GetId()))
}

func (f *calendarFrame) onCalLimitDates(evt *wx.Event) {
	if evt == nil {
		return
	}
	f.panel.limitDateRange(f.menuChecked(evt.GetId()))
}

func (f *calendarFrame) onCalAllowMonth(evt *wx.Event) {
	if evt == nil {
		return
	}
	cmd := wx.AsCommandEvent(evt)
	if cmd == nil {
		return
	}
	f.panel.cal.EnableMonthChange(cmd.IsChecked())
}

func (f *calendarFrame) onCalSeqMonth(evt *wx.Event) {
	if evt == nil {
		return
	}
	cmd := wx.AsCommandEvent(evt)
	if cmd == nil {
		return
	}
	f.panel.toggleCalStyle(cmd.IsChecked(), wx.CAL_SEQUENTIAL_MONTH_SELECTION)
}

func (f *calendarFrame) onCalShowSurroundingWeeks(evt *wx.Event) {
	if evt == nil {
		return
	}
	cmd := wx.AsCommandEvent(evt)
	if cmd == nil {
		return
	}
	f.panel.toggleCalStyle(cmd.IsChecked(), wx.CAL_SHOW_SURROUNDING_WEEKS)
}

func (f *calendarFrame) onCalShowWeekNumbers(evt *wx.Event) {
	if evt == nil {
		return
	}
	cmd := wx.AsCommandEvent(evt)
	if cmd == nil {
		return
	}
	f.panel.toggleCalStyle(cmd.IsChecked(), wx.CAL_SHOW_WEEK_NUMBERS)
}

func (f *calendarFrame) onSetDate(evt *wx.Event) {
	if evt == nil {
		return
	}
	dt := wx.NewDateTime()
	if dt == nil {
		return
	}
	defer dt.Delete()
	dt.Set(24, wx.DateTime_Dec, 2005, 22, 0, 0, 0)
	f.panel.setDate(dt)
}

func (f *calendarFrame) onToday(evt *wx.Event) {
	if evt == nil {
		return
	}
	dt := wx.NewDateTime()
	if dt == nil {
		return
	}
	defer dt.Delete()
	dt.Today()
	f.panel.setDate(dt)
}

func (f *calendarFrame) onBeginDST(evt *wx.Event) {
	if evt == nil {
		return
	}
	current := wx.NewDateTime()
	if current == nil {
		return
	}
	defer current.Delete()
	f.panel.cal.GetDate(current.Ptr())

	dst := wx.NewDateTime()
	if dst == nil {
		return
	}
	defer dst.Delete()
	wx.DateTime_GetBeginDST(current.GetYear(wx.DateTime_Local), wx.DateTime_Country_Default, dst.Ptr())
	f.panel.setDate(dst)
}

func (f *calendarFrame) onCalToggleResizable(evt *wx.Event) {
	if evt == nil {
		return
	}
	cmd := wx.AsCommandEvent(evt)
	if cmd == nil {
		return
	}

	sizer := wx.SizerFromPtr(f.panel.win.GetSizer())
	if sizer == nil {
		return
	}
	item := wx.SizerItemFromPtr(sizer.GetItemWindow(f.panel.cal.Ptr(), false))
	if item == nil {
		return
	}

	if cmd.IsChecked() {
		item.SetProportion(1)
		item.SetFlag(wx.EXPAND)
	} else { // not resizable
		item.SetProportion(0)
		item.SetFlag(wx.ALIGN_CENTER)
	}

	sizer.Layout()
}

// onCalRClick is bound to wxEVT_RIGHT_DOWN on the calendar control itself, with
// the frame as the logical handler -- MyFrame::OnCalRClick().
func (f *calendarFrame) onCalRClick(evt *wx.Event) {
	if evt == nil {
		return
	}
	mouse := wx.AsMouseEvent(evt)
	if mouse == nil {
		return
	}

	pt := mouse.GetPosition()

	dt := wx.NewDateTime()
	if dt == nil {
		return
	}
	defer dt.Delete()

	// kwxFFI declares the weekday out-parameter as int*, so a plain writable
	// local is passed directly; HitTest() writes it only for the header hit.
	var weekday int32
	hit := f.panel.cal.HitTest(pt.X, pt.Y, dt.Ptr(), unsafe.Pointer(&weekday))

	msg := fmt.Sprintf("Point (%d, %d) is ", pt.X, pt.Y)

	// The wxCAL_HITTEST_* values are runtime vars rather than constants in this
	// binding, but Go expression switches accept non-constant case values.
	switch hit {
	case wx.CAL_HITTEST_HEADER:
		msg += fmt.Sprintf("over %s",
			wx.DateTime_GetWeekDayName(int(weekday), wx.DateTime_Name_Full))
	case wx.CAL_HITTEST_DAY:
		msg += fmt.Sprintf("over %s", dt.FormatISODate())
	case wx.CAL_HITTEST_INCMONTH:
		msg += "over next month button"
	case wx.CAL_HITTEST_DECMONTH:
		msg += "over previous month button"
	case wx.CAL_HITTEST_SURROUNDING_WEEK:
		msg += "over a day from another month"
	default:
		// includes wxCAL_HITTEST_NOWHERE and any unexpected result
		msg += "nowhere"
	}

	wx.LogMessage(msg)
}

func (f *calendarFrame) onUpdateUIStartWithNone(evt *wx.Event) {
	update := wx.AsUpdateUIEvent(evt)
	if update == nil {
		return
	}
	// it only makes sense to start with invalid date if we can have no date
	update.Enable(f.menuChecked(idDatePickerAllowNone))
}

func (f *calendarFrame) onUpdateUIGenericOnly(evt *wx.Event) {
	update := wx.AsUpdateUIEvent(evt)
	if update == nil {
		return
	}
	update.Enable(f.panel.isUsingGeneric())
}

// menuChecked wraps GetMenuBar()->IsChecked(id).
func (f *calendarFrame) menuChecked(id int) bool {
	mb := wx.MenuBarFromPtr(f.frame.GetMenuBar())
	if mb == nil {
		return false
	}
	return mb.IsChecked(id)
}

// onAskDate is MyFrame::OnAskDate().
func (f *calendarFrame) onAskDate(evt *wx.Event) {
	if evt == nil {
		return
	}

	dt := wx.NewDateTime()
	if dt == nil {
		return
	}
	f.panel.cal.GetDate(dt.Ptr())

	style := wx.DP_DEFAULT
	if f.menuChecked(idDatePickerShowCentury) {
		style |= wx.DP_SHOWCENTURY
	}
	if f.menuChecked(idDatePickerDropDown) {
		style |= wx.DP_DROPDOWN
	}
	if f.menuChecked(idDatePickerAllowNone) {
		style |= wx.DP_ALLOWNONE
		if f.menuChecked(idDatePickerStartWithNone) {
			// dt = wxDefaultDateTime
			dt.Delete()
			dt = wx.NewDateTime()
			if dt == nil {
				return
			}
		}
	}

	dlg := newDateDialog(f.frame.Ptr(), dt, style)
	if dlg != nil {
		if dlg.showModal() == wx.ID_OK {
			chosen := dlg.date()
			if chosen != nil {
				if chosen.IsValid() {
					today := wx.NewDateTime()
					if today != nil {
						today.Today()
						if chosen.GetDay(wx.DateTime_Local) == today.GetDay(wx.DateTime_Local) &&
							chosen.GetMonth(wx.DateTime_Local) == today.GetMonth(wx.DateTime_Local) {
							// wxMessageBox() default style is wxOK | wxCENTRE
							wx.MessageBox("Happy birthday!", "Calendar Sample",
								wx.OK|wx.CENTRE, nil, -1, -1)
						}
						today.Delete()
					}

					f.panel.setDate(chosen)
					wx.LogStatus("Changed the date to your input")
				} else {
					wx.LogStatus("No date entered")
				}
				chosen.Delete()
			}
		}
		dlg.destroy()
	}

	dt.Delete()
}

// onAskTime is MyFrame::OnAskTime().
func (f *calendarFrame) onAskTime(evt *wx.Event) {
	if evt == nil {
		return
	}

	dlg := newTimeDialog(f.frame.Ptr())
	if dlg == nil {
		return
	}
	if dlg.showModal() == wx.ID_OK {
		if t := dlg.time(); t != nil {
			wx.LogMessage(fmt.Sprintf("You entered %s", t.FormatISOTime()))
			t.Delete()
		}
	}
	dlg.destroy()
}
