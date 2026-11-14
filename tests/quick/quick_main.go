// Command quick is the hand-written entry point for the kwxGO quick test app.
//
// quick.wxui generates quick_mainframe.go and quick_dlg.go into this folder; the
// generated files only build the (hidden) forms. This file owns the application
// bootstrap: create the frame, show it, register it as the top window, run the
// modal dialog and then enter the event loop.
//
// Build and run from the repository root with "build_quick.cmd" (build) or
// "build_quick.cmd run" (build, then launch).
package main

import (
	"github.com/KeyWorksRW/kwxGO/wx"
)

func main() {
	// wx.Run returns no exit code, so a clean exit is the only success signal.
	wx.Run(func() bool {
		// Register the image handlers before any form decodes an embedded bitmap.
		wx.InitAllImageHandlers()
		frame := newQuickMainFrame(nil)
		if frame == nil {
			return false
		}
		frame.Show()
		wx.SetTopWindow(frame.Ptr())

		dialog := newQuickDialog(frame.Ptr())
		if dialog == nil || dialog.dialog == nil {
			return false
		}

		// wxWidgets vetoes a frame close while an app-modal dialog is open
		// (wxTopLevelWindowBase::OnCloseWindow). A close request on the frame still has to end
		// the app, so record it, stop the handler chain (the default handler would veto) and
		// end the modal dialog; the frame is destroyed once the modal loop has unwound.
		closeRequested := false
		frame.Bind(wx.EVT_CLOSE_WINDOW, func(evt *wx.Event) {
			closeRequested = true
			evt.Skip(false)
			dialog.dialog.EndModal(wx.ID_CANCEL)
		})

		dialog.dialog.ShowModal()

		if closeRequested {
			// Destroy the frame and leave without entering the main loop, so the process exits
			// with code 0.
			frame.Destroy()
			return false
		}

		return true
	})
}
