package main

import (
	"fmt"
	"unsafe"

	"github.com/KeyWorksRW/kwxGO/wx"
)

// calendarPanel is MyPanel: a static text with the selected date plus the
// calendar control, both centred in a box sizer whose orientation depends on
// the screen aspect ratio.
type calendarPanel struct {
	owner *calendarFrame
	win   *wx.Panel

	date *wx.StaticText
	// cal is the shared wxCalendarCtrlBase API, valid for either the native
	// wxCalendarCtrl or a wxGenericCalendarCtrl.
	cal   *wx.CalendarCtrlBase
	sizer *wx.BoxSizer

	usingGeneric bool

	// dateLast is the function-static s_dateLast of OnCalendar(): clicking the
	// same date twice unmarks it.
	dateLast *wx.DateTime
}

// newCalendarPanel is MyPanel::MyPanel().
func newCalendarPanel(parent unsafe.Pointer, owner *calendarFrame) *calendarPanel {
	panel := wx.NewPanel(parent, wx.ID_ANY, -1, -1, -1, -1,
		wx.TAB_TRAVERSAL|wx.BORDER_NONE)
	if panel == nil {
		return nil
	}

	p := &calendarPanel{
		owner: owner,
		win:   panel,
		// wxHAS_NATIVE_CALENDARCTRL is defined on this platform, so the native
		// control is used first.
		usingGeneric: false,
	}

	today := wx.NewDateTime()
	if today == nil {
		return nil
	}
	today.Today()
	dateText := fmt.Sprintf("Selected date: %s", today.FormatISODate())
	today.Delete()

	p.date = wx.NewStaticText(panel.Ptr(), wx.ID_ANY, dateText, -1, -1, -1, -1, 0)

	invalid := wx.NewDateTime()
	if invalid == nil {
		return nil
	}
	p.cal = p.createCalendar(invalid, wx.CAL_SHOW_HOLIDAYS)
	invalid.Delete()
	if p.cal == nil {
		return nil
	}

	// adjust to vertical/horizontal display
	horizontal := wx.SystemSettings_GetMetric(wx.SYS_SCREEN_X) >
		wx.SystemSettings_GetMetric(wx.SYS_SCREEN_Y)
	orient := wx.VERTICAL
	if horizontal {
		orient = wx.HORIZONTAL
	}

	p.sizer = wx.NewBoxSizer(orient)
	if p.sizer == nil {
		return nil
	}

	flags := wx.NewSizerFlags(0)
	if flags != nil {
		flags.Center()
		flags.DoubleBorder(wx.ALL)
		if p.date != nil {
			p.sizer.AddWindowFlags(p.date.Ptr(), flags.Ptr())
		}
		flags.Delete()
	}

	flags = wx.NewSizerFlags(0)
	if flags != nil {
		flags.Center()
		p.sizer.AddWindowFlags(p.cal.Ptr(), flags.Ptr())
		flags.Delete()
	}

	panel.SetSizer(p.sizer.Ptr(), true)
	p.sizer.SetSizeHints(panel.Ptr())
	return p
}

// createCalendar is MyPanel::DoCreateCalendar(): it builds the native or the
// generic control, binds the right-click handler and the calendar events.
func (p *calendarPanel) createCalendar(dt *wx.DateTime, style int) *wx.CalendarCtrlBase {
	var cal *wx.CalendarCtrlBase

	if p.usingGeneric {
		generic := wx.NewGenericCalendarCtrl(p.win.Ptr(), idCalendarCtrl, dt.Ptr(),
			-1, -1, -1, -1, style)
		if generic == nil {
			return nil
		}
		// Borrow the shared wxCalendarCtrlBase view of the generic control.
		cal = wx.CalendarCtrlBaseFromPtr(generic.Ptr())
	} else {
		native := wx.NewCalendarCtrl(p.win.Ptr(), idCalendarCtrl, dt.Ptr(),
			-1, -1, -1, -1, style)
		if native == nil {
			return nil
		}
		// Borrow the shared wxCalendarCtrlBase view of the native control.
		cal = wx.CalendarCtrlBaseFromPtr(native.Ptr())
	}
	if cal == nil {
		return nil
	}

	// MyFrame::OnCalRClick() is connected to the calendar in the C++ code; the
	// handler only touches the calendar, so it is bound on the control itself.
	cal.Bind(wx.EVT_RIGHT_DOWN, p.owner.onCalRClick)
	p.bindCalendarEvents(cal)
	return cal
}

// bindCalendarEvents reproduces the MyPanel event table. The C++ table uses
// EVT_CALENDAR* with the calendar's ID; connecting on the control itself is
// equivalent and survives the control being recreated. EVT_CALENDAR is an
// alias for EVT_CALENDAR_DOUBLECLICKED (a day was double clicked).
func (p *calendarPanel) bindCalendarEvents(cal *wx.CalendarCtrlBase) {
	cal.Bind(wx.EVT_CALENDAR_DOUBLECLICKED, p.onCalendar)
	cal.Bind(wx.EVT_CALENDAR_PAGE_CHANGED, p.onCalMonthChange)
	cal.Bind(wx.EVT_CALENDAR_SEL_CHANGED, p.onCalendarChange)
	cal.Bind(wx.EVT_CALENDAR_WEEKDAY_CLICKED, p.onCalendarWeekDayClick)
	cal.Bind(wx.EVT_CALENDAR_WEEK_CLICKED, p.onCalendarWeekClick)
}

// recreateCalendar is MyPanel::RecreateCalendar().
func (p *calendarPanel) recreateCalendar(style int) {
	current := wx.NewDateTime()
	if current == nil {
		return
	}
	p.cal.GetDate(current.Ptr())

	next := p.createCalendar(current, style)
	current.Delete()
	if next == nil {
		return
	}

	old := p.cal
	if p.sizer.ReplaceWindow(old.Ptr(), next.Ptr(), false) {
		// The C++ code deletes the replaced control; Destroy() is the
		// wxWidgets-safe way to do that for a window removed from a sizer.
		old.Destroy()
	}
	p.cal = next

	p.sizer.Layout()
}

// toggleCalStyle is MyPanel::ToggleCalStyle().
func (p *calendarPanel) toggleCalStyle(on bool, flag int) {
	style := p.cal.GetWindowStyle()
	if on {
		style |= flag
	} else {
		style &^= flag
	}

	if flag == wx.CAL_SEQUENTIAL_MONTH_SELECTION || flag == wx.CAL_SHOW_WEEK_NUMBERS {
		// changing this style requires recreating the control
		p.recreateCalendar(style)
	} else { // just changing the style is enough
		p.cal.SetWindowStyle(style)
		p.cal.Refresh(true)
	}
}

// toggleUseGeneric is MyPanel::ToggleUseGeneric().
func (p *calendarPanel) toggleUseGeneric() {
	p.usingGeneric = !p.usingGeneric
	p.recreateCalendar(p.cal.GetWindowStyle())
}

// isUsingGeneric reports MyPanel::IsUsingGeneric().
func (p *calendarPanel) isUsingGeneric() bool { return p.usingGeneric }

// setDate is MyPanel::SetDate().
func (p *calendarPanel) setDate(dt *wx.DateTime) {
	if dt == nil {
		return
	}
	p.cal.SetDate(dt.Ptr())
}

// highlightSpecial is MyPanel::HighlightSpecial().
func (p *calendarPanel) highlightSpecial(on bool) {
	if on {
		// wxCalendarDateAttr(wxCAL_BORDER_ROUND, *wxRED)
		attrRedCircle := wx.NewCalendarDateAttrBorder(wx.CAL_BORDER_ROUND, wx.RED)
		// wxCalendarDateAttr(wxCAL_BORDER_SQUARE, *wxGREEN)
		attrGreenSquare := wx.NewCalendarDateAttrBorder(wx.CAL_BORDER_SQUARE, wx.GREEN)
		// wxCalendarDateAttr(*wxBLUE, *wxLIGHT_GREY)
		attrHeaderLike := wx.NewCalendarDateAttr(wx.BLUE, wx.LIGHT_GREY,
			wx.Null_Colour(), wx.Null_Font(), wx.CAL_BORDER_NONE)

		// SetAttr() takes ownership of the attribute.
		if attrRedCircle != nil {
			p.cal.SetAttr(17, attrRedCircle.Ptr())
		}
		if attrGreenSquare != nil {
			p.cal.SetAttr(29, attrGreenSquare.Ptr())
		}
		if attrHeaderLike != nil {
			p.cal.SetAttr(13, attrHeaderLike.Ptr())
		}
	} else { // off
		p.cal.ResetAttr(17)
		p.cal.ResetAttr(29)
		p.cal.ResetAttr(13)
	}

	p.cal.Refresh(true)
}

// limitDateRange is MyPanel::LimitDateRange(): toggle a restricted date range
// to the six months centred on today's date.
func (p *calendarPanel) limitDateRange(on bool) {
	if on {
		// limit the choice of date to 3 months around today
		diff := wx.NewDateSpanFromMonths(3)
		today := wx.NewDateTime()
		lower := wx.NewDateTime()
		upper := wx.NewDateTime()
		if diff == nil || today == nil || lower == nil || upper == nil {
			if diff != nil {
				diff.Delete()
			}
			if today != nil {
				today.Delete()
			}
			if lower != nil {
				lower.Delete()
			}
			if upper != nil {
				upper.Delete()
			}
			return
		}

		today.Today()

		// today - wxDateSpan::Months(3) and today + wxDateSpan::Months(3)
		today.SubtractDate(diff.Ptr(), lower.Ptr())
		today.AddDate(diff.Ptr(), upper.Ptr())
		diff.Delete()
		today.Delete()

		// Set the restricted date range.
		if p.cal.SetDateRange(lower.Ptr(), upper.Ptr()) {
			wx.LogStatus("Date range limited to 3 months around today.")

			firstValidDate := wx.NewDateTime()
			lastValidDate := wx.NewDateTime()
			if firstValidDate != nil && lastValidDate != nil {
				if p.cal.GetDateRange(firstValidDate.Ptr(), lastValidDate.Ptr()) {
					wx.LogMessage(fmt.Sprintf("First valid date: %s, last valid date: %s",
						firstValidDate.FormatISODate(), lastValidDate.FormatISODate()))
				} else {
					wx.LogWarning("Failed to get back the valid dates range.")
				}
			}
			if firstValidDate != nil {
				firstValidDate.Delete()
			}
			if lastValidDate != nil {
				lastValidDate.Delete()
			}
		} else {
			wx.LogWarning("Date range not supported.")
		}

		lower.Delete()
		upper.Delete()
	} else { // off
		// Remove the date restrictions.
		if p.cal.SetDateRange(nil, nil) {
			wx.LogStatus("Date choice is unlimited now.")
		} else {
			wx.LogWarning("Date range not supported.")
		}
	}

	p.cal.Refresh(true)
}

// ---------------------------------------------------------------------------
// Calendar event handlers
// ---------------------------------------------------------------------------

// onCalendar is MyPanel::OnCalendar(): clicking the same date twice unmarks it.
func (p *calendarPanel) onCalendar(evt *wx.Event) {
	if evt == nil {
		return
	}
	calEvt := wx.AsCalendarEvent(evt)
	if calEvt == nil {
		return
	}

	date := wx.NewDateTime()
	if date == nil {
		return
	}
	calEvt.GetDate(date.Ptr())

	mark := p.dateLast == nil || !p.dateLast.IsValid() || !date.IsEqualTo(p.dateLast.Ptr())

	if p.dateLast != nil {
		p.dateLast.Delete()
	}
	// Copy the event date into a persistent wxDateTime, as s_dateLast does.
	p.dateLast = wx.NewDateTimeFromTimeT(int64(date.GetTicks()))

	p.cal.Mark(date.GetDay(wx.DateTime_Local), mark)

	markText := ""
	if !mark {
		markText = "un"
	}
	wx.LogMessage(fmt.Sprintf("Selected (and %smarked) %s from calendar.",
		markText, date.FormatISODate()))

	date.Delete()
}

// onCalendarChange is MyPanel::OnCalendarChange().
func (p *calendarPanel) onCalendarChange(evt *wx.Event) {
	if evt == nil {
		return
	}
	calEvt := wx.AsCalendarEvent(evt)
	if calEvt == nil {
		return
	}

	date := wx.NewDateTime()
	if date == nil {
		return
	}
	calEvt.GetDate(date.Ptr())

	s := fmt.Sprintf("Selected date: %s", date.FormatISODate())
	if p.date != nil {
		p.date.SetLabel(s)
	}
	wx.LogStatus(s)

	date.Delete()
}

// onCalMonthChange is MyPanel::OnCalMonthChange().
func (p *calendarPanel) onCalMonthChange(evt *wx.Event) {
	if evt == nil {
		return
	}
	calEvt := wx.AsCalendarEvent(evt)
	if calEvt == nil {
		return
	}

	date := wx.NewDateTime()
	if date == nil {
		return
	}
	calEvt.GetDate(date.Ptr())

	wx.LogStatus(fmt.Sprintf("Calendar month changed to %s %d",
		wx.DateTime_GetMonthName(date.GetMonth(wx.DateTime_Local), wx.DateTime_Name_Full),
		date.GetYear(wx.DateTime_Local)))

	date.Delete()
}

// onCalendarWeekDayClick is MyPanel::OnCalendarWeekDayClick().
func (p *calendarPanel) onCalendarWeekDayClick(evt *wx.Event) {
	if evt == nil {
		return
	}
	calEvt := wx.AsCalendarEvent(evt)
	if calEvt == nil {
		return
	}

	wx.LogMessage(fmt.Sprintf("Clicked on %s",
		wx.DateTime_GetWeekDayName(calEvt.GetWeekDay(), wx.DateTime_Name_Full)))
}

// onCalendarWeekClick is MyPanel::OnCalendarWeekClick().
func (p *calendarPanel) onCalendarWeekClick(evt *wx.Event) {
	if evt == nil {
		return
	}
	calEvt := wx.AsCalendarEvent(evt)
	if calEvt == nil {
		return
	}

	date := wx.NewDateTime()
	if date == nil {
		return
	}
	calEvt.GetDate(date.Ptr())

	wx.LogMessage(fmt.Sprintf("Clicked on week %d",
		date.GetWeekOfYear(wx.DateTime_Monday_First, wx.DateTime_Local)))

	date.Delete()
}
