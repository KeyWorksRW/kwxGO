// kwxGO Demo application
//
// Demonstrates menus, status bar, sizers, text input, and event handling.
// Implements the standard kwxFFI demo UI specification.
package main

import (
	"github.com/KeyWorksRW/kwxGO/wx"
)

func main() {
	wx.Run(func() bool {
		wx.SetAppName("kwxGO Demo")

		// Create main frame (480×320, centered)
		frame := wx.NewFrame(nil, wx.ID_ANY, "kwxGO Demo",
			-1, -1, 480, 320, wx.DEFAULT_FRAME_STYLE)

		// ---- Menu Bar ----
		buildMenuBar(frame)

		// ---- Status Bar ----
		frame.CreateStatusBar(1, 0)
		frame.SetStatusText("Ready", 0)

		// ---- Client Area ----
		panel := wx.NewPanel(frame.Ptr(), wx.ID_ANY,
			-1, -1, -1, -1, 0)

		// Vertical box sizer (top-level)
		vbox := wx.NewBoxSizer(wx.VERTICAL)

		// Horizontal box sizer with label + text control
		hbox := wx.NewBoxSizer(wx.HORIZONTAL)

		label := wx.NewStaticText(panel.Ptr(), wx.ID_ANY, "Text:",
			-1, -1, -1, -1, 0)
		hbox.AddWindow(label.Ptr(), 0, wx.ALL|wx.ALIGN_CENTER_VERTICAL, 5, nil)

		textCtrl := wx.NewTextCtrl(panel.Ptr(), wx.ID_ANY, "",
			-1, -1, -1, -1, 0)
		hbox.AddWindow(textCtrl.Ptr(), 1, wx.ALL|wx.EXPAND, 5, nil)

		vbox.AddSizer(hbox.Ptr(), 0, wx.EXPAND|wx.ALL, 5, nil)

		// "Click Me" button
		button := wx.NewButton(panel.Ptr(), wx.ID_ANY, "Click Me",
			-1, -1, -1, -1, 0)
		vbox.AddWindow(button.Ptr(), 0, wx.ALL, 10, nil)

		panel.SetSizer(vbox.Ptr(), false)

		// ---- Event Handlers ----

		// Button click
		button.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			text := textCtrl.GetValue()
			if text == "" {
				wx.MessageBox("(no text entered)", "Demo",
					wx.OK|wx.ICON_INFORMATION, frame.Ptr(), -1, -1)
			} else {
				wx.MessageBox(text, "You typed",
					wx.OK|wx.ICON_INFORMATION, frame.Ptr(), -1, -1)
			}
		})

		// File > Exit
		frame.BindId(wx.EVT_MENU, wx.ID_EXIT, func(evt *wx.Event) {
			frame.Close(true)
		})

		// Window close (X button)
		frame.Bind(wx.EVT_CLOSE_WINDOW, func(evt *wx.Event) {
			frame.Destroy()
		})

		// Help > About
		frame.BindId(wx.EVT_MENU, wx.ID_ABOUT, func(evt *wx.Event) {
			wx.MessageBox("kwxGO Demo\nGo bindings for wxWidgets via kwxFFI",
				"About Demo", wx.OK|wx.ICON_INFORMATION, frame.Ptr(), -1, -1)
		})

		// Show
		frame.Center(wx.BOTH)
		wx.SetTopWindow(frame.Ptr())
		frame.Show()

		return true
	})
}

// buildMenuBar creates and attaches the menu bar to the frame.
func buildMenuBar(frame *wx.Frame) {
	// File menu
	fileMenu := wx.NewMenu("", 0)
	fileMenu.Append(wx.ID_NEW, "&New", "Create new", false)
	fileMenu.Append(wx.ID_OPEN, "&Open", "Open file", false)
	fileMenu.Append(wx.ID_SAVE, "&Save", "Save file", false)
	fileMenu.AppendSeparator()
	fileMenu.Append(wx.ID_EXIT, "E&xit", "Exit application", false)

	// Edit menu
	editMenu := wx.NewMenu("", 0)
	editMenu.Append(wx.ID_UNDO, "&Undo", "", false)
	editMenu.Append(wx.ID_REDO, "&Redo", "", false)
	editMenu.AppendSeparator()
	editMenu.Append(wx.ID_CUT, "Cu&t", "", false)
	editMenu.Append(wx.ID_COPY, "&Copy", "", false)
	editMenu.Append(wx.ID_PASTE, "&Paste", "", false)

	// Help menu
	helpMenu := wx.NewMenu("", 0)
	helpMenu.Append(wx.ID_ABOUT, "&About", "About this demo", false)

	// Menu bar
	menuBar := wx.NewMenuBar(0)
	menuBar.Append(fileMenu.Ptr(), "&File")
	menuBar.Append(editMenu.Ptr(), "&Edit")
	menuBar.Append(helpMenu.Ptr(), "&Help")

	frame.SetMenuBar(menuBar.Ptr())
}
