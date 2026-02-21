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
			wx.DefaultPosition, wx.Size{W: 480, H: 320},
			wx.DEFAULT_FRAME_STYLE)

		// ---- Menu Bar ----
		buildMenuBar(frame)

		// ---- Status Bar ----
		frame.CreateStatusBar(1, 0)
		frame.SetStatusText("Ready", 0)

		// ---- Client Area ----
		panel := wx.NewPanel(frame, wx.ID_ANY,
			wx.DefaultPosition, wx.DefaultSize, 0)

		// Vertical box sizer (top-level)
		vbox := wx.NewBoxSizer(wx.VERTICAL)

		// Horizontal box sizer with label + text control
		hbox := wx.NewBoxSizer(wx.HORIZONTAL)

		label := wx.NewStaticText(panel, wx.ID_ANY, "Text:",
			wx.DefaultPosition, wx.DefaultSize, 0)
		hbox.Add(label, 0, wx.ALL|wx.ALIGN_CENTER_VERTICAL, 5)

		textCtrl := wx.NewTextCtrl(panel, wx.ID_ANY, "",
			wx.DefaultPosition, wx.DefaultSize, 0)
		textCtrl.SetHint("hint text")
		hbox.Add(textCtrl, 1, wx.ALL|wx.EXPAND, 5)

		vbox.AddSizer(hbox, 0, wx.EXPAND|wx.ALL, 5)

		// "Click Me" button
		button := wx.NewButton(panel, wx.ID_ANY, "Click Me",
			wx.DefaultPosition, wx.DefaultSize, 0)
		vbox.Add(button, 0, wx.ALL, 10)

		panel.SetSizerInterface(vbox)

		// ---- Event Handlers ----

		// Button click
		button.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			text := textCtrl.GetValue()
			if text == "" {
				wx.MessageBox("(no text entered)", "Demo",
					wx.OK|wx.ICON_INFORMATION, frame)
			} else {
				wx.MessageBox(text, "You typed",
					wx.OK|wx.ICON_INFORMATION, frame)
			}
		})

		// File > Exit
		frame.BindId(wx.EVT_MENU, wx.ID_EXIT, func(evt *wx.Event) {
			frame.Close(true)
		})

		// Help > About
		frame.BindId(wx.EVT_MENU, wx.ID_ABOUT, func(evt *wx.Event) {
			wx.MessageBox("kwxGO Demo\nGo bindings for wxWidgets via kwxFFI",
				"About Demo", wx.OK|wx.ICON_INFORMATION, frame)
		})

		// Window close (X button)
		frame.Bind(wx.EVT_CLOSE_WINDOW, func(evt *wx.Event) {
			wx.ExitMainLoop()
		})

		// Show
		frame.Center()
		wx.SetTopWindow(frame)
		frame.Show(true)

		return true
	})
}

// buildMenuBar creates and attaches the menu bar to the frame.
func buildMenuBar(frame *wx.Frame) {
	// File menu
	fileMenu := wx.NewMenu("")
	fileMenu.Append(wx.ID_NEW, "&New", "Create new")
	fileMenu.Append(wx.ID_OPEN, "&Open", "Open file")
	fileMenu.Append(wx.ID_SAVE, "&Save", "Save file")
	fileMenu.AppendSeparator()
	fileMenu.Append(wx.ID_EXIT, "E&xit", "Exit application")

	// Edit menu
	editMenu := wx.NewMenu("")
	editMenu.Append(wx.ID_UNDO, "&Undo", "")
	editMenu.Append(wx.ID_REDO, "&Redo", "")
	editMenu.AppendSeparator()
	editMenu.Append(wx.ID_CUT, "Cu&t", "")
	editMenu.Append(wx.ID_COPY, "&Copy", "")
	editMenu.Append(wx.ID_PASTE, "&Paste", "")

	// Help menu
	helpMenu := wx.NewMenu("")
	helpMenu.Append(wx.ID_ABOUT, "&About", "About this demo")

	// Menu bar
	menuBar := wx.NewMenuBar()
	menuBar.Append(fileMenu, "&File")
	menuBar.Append(editMenu, "&Edit")
	menuBar.Append(helpMenu, "&Help")

	frame.SetMenuBar(menuBar.Ptr())
}
