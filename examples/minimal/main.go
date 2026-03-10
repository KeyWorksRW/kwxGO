// Minimal wxGO example application
//
// This example demonstrates the basic structure of a wxGO application.
// To build and run:
//
//  1. Build kwxFFI first:
//     cmake -B build -G Ninja
//     ninja -C build -f build-Debug.ninja
//
//  2. Build this example:
//     go build -o minimal.exe ./examples/minimal/
//
//  3. Run:
//     ./minimal.exe
package main

import (
	"github.com/KeyWorksRW/kwxGO/wx"
)

func main() {
	wx.Run(func() bool {
		// Create main frame
		frame := wx.NewFrame(nil, wx.ID_ANY, "Hello wxGO",
			-1, -1, 400, 300, wx.DEFAULT_FRAME_STYLE)

		// Create panel to hold controls
		panel := wx.NewPanel(frame.Ptr(), wx.ID_ANY,
			-1, -1, -1, -1, 0)

		// Create a button
		button := wx.NewButton(panel.Ptr(), wx.ID_ANY, "Click Me!",
			50, 50, -1, -1, 0)

		// Bind button click event
		button.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			wx.MessageBox("Hello from Go!", "wxGO",
				wx.OK|wx.ICON_INFORMATION, frame.Ptr(), -1, -1)
		})

		// Show the frame
		frame.Show()

		// Set as top window so MainLoop knows the app is active
		wx.SetTopWindow(frame.Ptr())

		return true
	})
}
