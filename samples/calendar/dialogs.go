package main

import (
	"unsafe"

	"github.com/KeyWorksRW/kwxGO/wx"
)

// dateDialog is MyDateDialog: a modal dialog holding a wxDatePickerCtrl, a
// static text showing the date in ISO format, and an OK/Cancel button sizer.
type dateDialog struct {
	dlg  *wx.Dialog
	pick *wx.DatePickerCtrl
	text *wx.StaticText
}

// newDateDialog is MyDateDialog::MyDateDialog(). The generic picker variant is
// compiled out on this platform (wxUSE_DATEPICKCTRL_GENERIC is 0), so only the
// native control is created.
func newDateDialog(parent unsafe.Pointer, dt *wx.DateTime, dtpStyle int) *dateDialog {
	dlg := wx.NewDialog(parent, wx.ID_ANY, "Calendar: Choose a date",
		-1, -1, -1, -1, wx.DEFAULT_DIALOG_STYLE)
	if dlg == nil {
		return nil
	}

	d := &dateDialog{dlg: dlg}

	// The picker Create() wrapper passes the supplied date to the constructor;
	// a nil date means wxDefaultDateTime (the "start with no date" option).
	var datePtr unsafe.Pointer
	if dt != nil {
		datePtr = dt.Ptr()
	}
	d.pick = wx.NewDatePickerCtrl(dlg.Ptr(), wx.ID_ANY, datePtr,
		-1, -1, -1, -1, dtpStyle)
	if d.pick == nil {
		return nil
	}

	// SetRange(wxDateTime(1, Jan, 1900), wxDefaultDateTime)
	lower := wx.NewDateTime()
	if lower != nil {
		lower.Set(1, wx.DateTime_Jan, 1900, 0, 0, 0, 0)
		d.pick.SetRange(lower.Ptr(), nil)
		lower.Delete()
	}

	initial := ""
	if dt != nil && dt.IsValid() {
		initial = dt.FormatISODate()
	}
	d.text = wx.NewStaticText(dlg.Ptr(), wx.ID_ANY, initial, -1, -1, -1, -1, 0)

	flags := wx.NewSizerFlags(0)
	if flags == nil {
		return nil
	}
	flags.Center()
	flags.Border(wx.ALL)

	sizerMain := wx.NewFlexGridSizer(0, 2, 0, 0)
	if sizerMain == nil {
		flags.Delete()
		return nil
	}

	label := wx.NewStaticText(dlg.Ptr(), wx.ID_ANY, "Enter &date:", -1, -1, -1, -1, 0)
	if label != nil {
		sizerMain.AddWindowFlags(label.Ptr(), flags.Ptr())
	}
	sizerMain.AddWindowFlags(d.pick.Ptr(), flags.Ptr())

	isoLabel := wx.NewStaticText(dlg.Ptr(), wx.ID_ANY, "Date in ISO format:",
		-1, -1, -1, -1, 0)
	if isoLabel != nil {
		sizerMain.AddWindowFlags(isoLabel.Ptr(), flags.Ptr())
	}
	if d.text != nil {
		sizerMain.AddWindowFlags(d.text.Ptr(), flags.Ptr())
	}

	sizerTop := wx.NewBoxSizer(wx.VERTICAL)
	if sizerTop == nil {
		flags.Delete()
		return nil
	}
	sizerTop.AddSizerFlags(sizerMain.Ptr(), flags.Ptr())

	if buttons := dlg.CreateStdDialogButtonSizer(wx.OK | wx.CANCEL); buttons != nil {
		sizerTop.AddSizerFlags(buttons, flags.Ptr())
	}

	dlg.SetSizerAndFit(sizerTop.Ptr(), true)
	flags.Delete()

	// EVT_DATE_CHANGED is handled by the dialog in the C++ event table, so the
	// binding stays on the dialog even though wxDatePickerCtrl now inherits Bind
	// from wxControl.
	dlg.Bind(wx.EVT_DATE_CHANGED, d.onDateChange)
	return d
}

// showModal is wxDialog::ShowModal().
func (d *dateDialog) showModal() int {
	if d == nil || d.dlg == nil {
		return wx.ID_CANCEL
	}
	return d.dlg.ShowModal()
}

// date is MyDateDialog::GetDate(): the current picker value, owned by the
// caller.
func (d *dateDialog) date() *wx.DateTime {
	if d == nil || d.pick == nil {
		return nil
	}
	dt := wx.NewDateTime()
	if dt == nil {
		return nil
	}
	d.pick.GetValue(dt.Ptr())
	return dt
}

// onDateChange is MyDateDialog::OnDateChange().
func (d *dateDialog) onDateChange(evt *wx.Event) {
	if evt == nil || d.text == nil {
		return
	}
	dateEvt := wx.AsDateEvent(evt)
	if dateEvt == nil {
		return
	}

	dt := wx.NewDateTime()
	if dt == nil {
		return
	}
	dateEvt.GetDate(dt.Ptr())

	if dt.IsValid() {
		d.text.SetLabel(dt.FormatISODate())
	} else {
		d.text.SetLabel("")
	}

	dt.Delete()
}

// destroy releases the dialog's native window.
func (d *dateDialog) destroy() {
	if d == nil || d.dlg == nil {
		return
	}
	d.dlg.Destroy()
	d.dlg = nil
}

// timeDialog is MyTimeDialog: a modal dialog holding a wxTimePickerCtrl, a
// static text showing the time in ISO format, and an OK/Cancel button sizer.
type timeDialog struct {
	dlg  *wx.Dialog
	pick *wx.TimePickerCtrl
	text *wx.StaticText
}

// newTimeDialog is MyTimeDialog::MyTimeDialog(). The generic picker variant is
// compiled out on this platform (wxUSE_TIMEPICKCTRL_GENERIC is 0).
func newTimeDialog(parent unsafe.Pointer) *timeDialog {
	dlg := wx.NewDialog(parent, wx.ID_ANY, "Calendar: Choose time",
		-1, -1, -1, -1, wx.DEFAULT_DIALOG_STYLE)
	if dlg == nil {
		return nil
	}

	d := &timeDialog{dlg: dlg}

	d.pick = wx.NewTimePickerCtrl(dlg.Ptr(), wx.ID_ANY, nil, -1, -1, -1, -1, 0)
	if d.pick == nil {
		return nil
	}

	d.text = wx.NewStaticText(dlg.Ptr(), wx.ID_ANY, d.timeISO(), -1, -1, -1, -1, 0)

	flags := wx.NewSizerFlags(0)
	if flags == nil {
		return nil
	}
	flags.Center()
	flags.Border(wx.ALL)

	sizerMain := wx.NewFlexGridSizer(0, 2, 0, 0)
	if sizerMain == nil {
		flags.Delete()
		return nil
	}

	label := wx.NewStaticText(dlg.Ptr(), wx.ID_ANY, "Enter &time:", -1, -1, -1, -1, 0)
	if label != nil {
		sizerMain.AddWindowFlags(label.Ptr(), flags.Ptr())
	}
	sizerMain.AddWindowFlags(d.pick.Ptr(), flags.Ptr())

	isoLabel := wx.NewStaticText(dlg.Ptr(), wx.ID_ANY, "Time in ISO format:",
		-1, -1, -1, -1, 0)
	if isoLabel != nil {
		sizerMain.AddWindowFlags(isoLabel.Ptr(), flags.Ptr())
	}
	if d.text != nil {
		sizerMain.AddWindowFlags(d.text.Ptr(), flags.Ptr())
	}

	sizerTop := wx.NewBoxSizer(wx.VERTICAL)
	if sizerTop == nil {
		flags.Delete()
		return nil
	}
	sizerTop.AddSizerFlags(sizerMain.Ptr(), flags.Ptr())

	if buttons := dlg.CreateStdDialogButtonSizer(wx.OK | wx.CANCEL); buttons != nil {
		sizerTop.AddSizerFlags(buttons, flags.Ptr())
	}

	dlg.SetSizerAndFit(sizerTop.Ptr(), true)
	flags.Delete()

	// EVT_TIME_CHANGED is handled by the dialog in the C++ event table, so the
	// binding stays on the dialog even though wxTimePickerCtrl now inherits Bind
	// from wxControl.
	dlg.Bind(wx.EVT_TIME_CHANGED, d.onTimeChange)
	return d
}

// timeISO is GetTime().FormatISOTime().
func (d *timeDialog) timeISO() string {
	if d == nil || d.pick == nil {
		return ""
	}
	dt := wx.NewDateTime()
	if dt == nil {
		return ""
	}
	defer dt.Delete()
	d.pick.GetValue(dt.Ptr())
	return dt.FormatISOTime()
}

// showModal is wxDialog::ShowModal().
func (d *timeDialog) showModal() int {
	if d == nil || d.dlg == nil {
		return wx.ID_CANCEL
	}
	return d.dlg.ShowModal()
}

// time is MyTimeDialog::GetTime(): the current picker value, owned by the
// caller.
func (d *timeDialog) time() *wx.DateTime {
	if d == nil || d.pick == nil {
		return nil
	}
	dt := wx.NewDateTime()
	if dt == nil {
		return nil
	}
	d.pick.GetValue(dt.Ptr())
	return dt
}

// onTimeChange is MyTimeDialog::OnTimeChange().
func (d *timeDialog) onTimeChange(evt *wx.Event) {
	if evt == nil || d.text == nil {
		return
	}
	dateEvt := wx.AsDateEvent(evt)
	if dateEvt == nil {
		return
	}

	dt := wx.NewDateTime()
	if dt == nil {
		return
	}
	dateEvt.GetDate(dt.Ptr())
	d.text.SetLabel(dt.FormatISOTime())
	dt.Delete()
}

// destroy releases the dialog's native window.
func (d *timeDialog) destroy() {
	if d == nil || d.dlg == nil {
		return
	}
	d.dlg.Destroy()
	d.dlg = nil
}
