// Package main ports the wxWidgets "minimal" sample to kwxGO.
//
// Original: wxSamples/src/minimal/minimal.cpp
//
// MyApp               -> the wx.Run onInit closure in main.
// MyFrame             -> newMinimalFrame plus the two handlers below.
// MyFrame event table -> BindId calls on the frame.
package main

import (
	"fmt"

	"github.com/KeyWorksRW/kwxGO/wx"
)

// sampleIconSize is the pixel size of the original ../sample.xpm asset.
const sampleIconSize = 32

// newMinimalFrame is MyFrame::MyFrame(const wxString& title): it creates the
// frame, then sets the icon, the menu bar and the status bar in that order.
//
// The C++ source has no size or position argument (wxDefaultSize /
// wxDefaultPosition) and never calls Fit(), so -1 is passed for all four.
func newMinimalFrame(title string) *wx.Frame {
	frame := wx.NewFrame(nil, wx.ID_ANY, title, -1, -1, -1, -1, wx.DEFAULT_FRAME_STYLE)
	if frame == nil {
		return nil
	}

	// SetIcon(wxICON(sample)): the sample.xpm artwork, converted to sample.svg.
	if icon := loadSampleIcon(); icon != nil {
		frame.SetIcon(icon.Ptr())
		icon.Delete()
	}

	// wxUSE_MENUBAR branch of the original constructor. The #else branch (a
	// wxBoxSizer holding an "About..." wxButton, bound to the same handler) is
	// unreachable in a default wxWidgets build and is deliberately not ported.
	fileMenu := wx.NewMenu("", 0)
	helpMenu := wx.NewMenu("", 0)
	if fileMenu == nil || helpMenu == nil {
		return nil
	}

	// The "About" item belongs to the help menu and must keep wxID_ABOUT.
	helpMenu.Append(wx.ID_ABOUT, "&About\tF1", "Show about dialog", false)
	fileMenu.Append(wx.ID_EXIT, "E&xit\tAlt-X", "Quit this program", false)

	menuBar := wx.NewMenuBar(0)
	if menuBar == nil {
		return nil
	}
	menuBar.Append(fileMenu.Ptr(), "&File")
	menuBar.Append(helpMenu.Ptr(), "&Help")
	frame.SetMenuBar(menuBar.Ptr())

	// CreateStatusBar(2) — default style, default id 0 — then the initial text.
	frame.CreateStatusBar(2, wx.STB_DEFAULT_STYLE, 0)
	frame.SetStatusText("Welcome to wxWidgets!", 0)

	// wxBEGIN_EVENT_TABLE(MyFrame, wxFrame)
	//     EVT_MENU(Minimal_Quit,  MyFrame::OnQuit)
	//     EVT_MENU(Minimal_About, MyFrame::OnAbout)
	// wxEND_EVENT_TABLE()
	frame.BindId(wx.EVT_MENU, wx.ID_EXIT, func(evt *wx.Event) { onQuit(frame, evt) })
	frame.BindId(wx.EVT_MENU, wx.ID_ABOUT, func(evt *wx.Event) { onAbout(frame, evt) })

	return frame
}

// loadSampleIcon loads the SVG rendition of the sample's XPM icon and converts
// it to a wxIcon for the frame.
func loadSampleIcon() *wx.Icon {
	bundle := wx.BitmapBundle_FromSVGFile("sample.svg", sampleIconSize, sampleIconSize)
	if bundle == nil {
		return nil
	}
	defer bundle.Delete()

	bitmap := bundle.GetBitmap(sampleIconSize, sampleIconSize)
	if bitmap == nil {
		return nil
	}
	defer bitmap.Delete()

	icon := wx.NewIconDefault()
	if icon == nil {
		return nil
	}
	icon.CopyFromBitmap(bitmap.Ptr())
	if !icon.IsOk() {
		icon.Delete()
		return nil
	}
	return icon
}

// onQuit is MyFrame::OnQuit(): true forces the frame to close.
func onQuit(frame *wx.Frame, evt *wx.Event) {
	if evt == nil {
		return
	}
	frame.Close(true)
}

// onAbout is MyFrame::OnAbout(): the return value of the message box is ignored.
func onAbout(frame *wx.Frame, evt *wx.Event) {
	if evt == nil {
		return
	}
	wx.MessageBox(fmt.Sprintf(
		"Welcome to %s!\n"+
			"\n"+
			"This is the minimal wxWidgets sample\n"+
			"running under %s.",
		wx.VersionString(),
		wx.GetOsDescription()),
		"About wxWidgets minimal sample",
		wx.OK|wx.ICON_INFORMATION,
		frame.Ptr(), -1, -1)
}

// main is MyApp::OnInit().
func main() {
	wx.Run(func() bool {
		frame := newMinimalFrame("Minimal wxWidgets App")
		if frame == nil {
			return false
		}

		frame.Show()
		wx.SetTopWindow(frame.Ptr())
		return true
	})
}
