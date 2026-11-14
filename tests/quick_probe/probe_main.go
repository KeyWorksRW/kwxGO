// Hand-written probe entry for the kwxGO probe lane (tests/quick_probe).
//
// probe.wxui (../quick/tests/quick_probe/probe.wxui) generates
// probe_mainframe.go and probe_dlg.go into gen/; this package is built from
// gen/ because Go requires one directory per package, so build_probe.cmd copies
// this file there before invoking the compiler.
//
// Probe output never goes to tests/quick.
package main

import (
	"github.com/KeyWorksRW/kwxGO/wx"
)

func main() {
	wx.Run(func() bool {
		// Register the image handlers before any form decodes an embedded bitmap; a form
		// handing an unregistered PNG format to wxImage produces an invalid bitmap.
		wx.InitAllImageHandlers()
		frame := newProbeMainFrame(nil)
		if frame == nil {
			return false
		}
		frame.Show()
		wx.SetTopWindow(frame.Ptr())

		dialog := newProbeDialog(frame.Ptr())
		if dialog == nil || dialog.dialog == nil {
			return false
		}
		dialog.dialog.ShowModal()

		return true
	})
}
