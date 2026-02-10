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
			wx.DefaultPosition, wx.Size{W: 400, H: 300},
			wx.DEFAULT_FRAME_STYLE)

		// Create panel to hold controls
		panel := wx.NewPanel(frame, wx.ID_ANY,
			wx.DefaultPosition, wx.DefaultSize, 0)

		// Create a button
		button := wx.NewButton(panel, wx.ID_ANY, "Click Me!",
			wx.Point{X: 50, Y: 50}, wx.DefaultSize, 0)

		// Bind button click event
		button.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			wx.MessageBox("Hello from Go!", "wxGO",
				wx.OK|wx.ICON_INFORMATION, frame)
		})

		// Show the frame
		frame.Show(true)

		// Set as top window so MainLoop knows the app is active
		wx.SetTopWindow(frame)

		return true
	})
}
