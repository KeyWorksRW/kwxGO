// Command about is the kwxGO port of the wxWidgets "about" sample
// (samples/html/about/about.cpp).
//
// The frame has a File menu whose About item opens a modal dialog holding a
// wxHtmlWindow. That window renders data/about.htm -- which itself embeds
// data/logo.png -- over the tiled SVG background data/bg.svg, sized to the
// parsed page instead of scrolling.
package main

import (
	"github.com/KeyWorksRW/kwxGO/wx"
)

func main() {
	wx.Run(func() bool {
		// we use a PNG image in our HTML page
		if handler := wx.NewPNGHandler(); handler != nil {
			// wxImage::AddHandler() takes ownership of the handler.
			wx.Image_AddHandler(handler.Ptr())
		}

		frame := newAboutFrame()
		if frame == nil {
			return false
		}

		// create and show the main application window
		frame.Show()
		wx.SetTopWindow(frame.Ptr())
		return true
	})
}

// newAboutFrame is MyFrame::MyFrame(): it creates the frame, its icon and the
// File menu, and connects the two menu handlers.
func newAboutFrame() *wx.Frame {
	frame := wx.NewFrame(nil, wx.ID_ANY, "wxHtmlWindow testing application",
		-1, -1, -1, -1, wx.DEFAULT_FRAME_STYLE)
	if frame == nil {
		return nil
	}

	if icon := wx.NewIconFromXPM(sampleXPM); icon != nil {
		frame.SetIcon(icon.Ptr())
	}

	// create a menu bar
	menuFile := wx.NewMenu("", 0)
	if menuFile == nil {
		return nil
	}

	// An empty label makes wxWidgets use the stock label for the id, exactly
	// like the C++ Append(wxID_ABOUT) / Append(wxID_EXIT) overloads.
	menuFile.Append(wx.ID_ABOUT, "", "", false)
	menuFile.Append(wx.ID_EXIT, "", "", false)

	menuBar := wx.NewMenuBar(0)
	if menuBar == nil {
		return nil
	}
	menuBar.Append(menuFile.Ptr(), "&File")
	frame.SetMenuBar(menuBar.Ptr())

	frame.BindId(wx.EVT_MENU, wx.ID_ABOUT, func(evt *wx.Event) {
		if evt == nil {
			return
		}
		showAboutDialog(frame)
	})

	frame.BindId(wx.EVT_MENU, wx.ID_EXIT, func(evt *wx.Event) {
		if evt == nil {
			return
		}
		// true is to force the frame to close
		frame.Close(true)
	})

	return frame
}

// showAboutDialog is MyFrame::OnAbout(): it builds and shows the modal dialog
// with the HTML about page.
func showAboutDialog(parent *wx.Frame) {
	dlg := wx.NewDialog(parent.Ptr(), wx.ID_ANY, "About",
		-1, -1, -1, -1, wx.DEFAULT_DIALOG_STYLE)
	if dlg == nil {
		return
	}

	topsizer := wx.NewBoxSizer(wx.VERTICAL)
	if topsizer == nil {
		dlg.Delete()
		return
	}

	html := wx.NewHtmlWindow(dlg.Ptr(), wx.ID_ANY, -1, -1, 380, 160,
		wx.HW_SCROLLBAR_NEVER, "htmlWindow")
	if html == nil {
		dlg.Delete()
		return
	}

	html.SetBorders(0)
	html.LoadPage("data/about.htm")
	// SetHTMLBackgroundImage() and SetBackgroundImage() are the same wx
	// method in this wxWidgets version; both take a wxBitmapBundle.
	if bg := wx.BitmapBundle_FromSVGFile("data/bg.svg", 65, 45); bg != nil {
		html.SetBackgroundImage(bg.Ptr())
		bg.Delete()
	}

	// Shrink the window to the size the page needs.
	if rep := wx.HtmlCellFromPtr(html.GetInternalRepresentation()); rep != nil {
		html.SetInitialSize(rep.GetWidth(), rep.GetHeight())
	}

	htmlFlags := wx.NewSizerFlags(1)
	htmlFlags.DoubleBorder(wx.ALL)
	topsizer.AddWindowFlags(html.Ptr(), htmlFlags.Ptr())
	htmlFlags.Delete()

	if line := wx.NewStaticLine(dlg.Ptr(), wx.ID_ANY, -1, -1, -1, -1, 0); line != nil {
		lineFlags := wx.NewSizerFlags(0)
		lineFlags.Expand()
		lineFlags.DoubleHorzBorder()
		topsizer.AddWindowFlags(line.Ptr(), lineFlags.Ptr())
		lineFlags.Delete()
	}

	if ok := wx.NewButton(dlg.Ptr(), wx.ID_OK, "OK", -1, -1, -1, -1, 0); ok != nil {
		ok.SetDefault()

		okFlags := wx.NewSizerFlags(0)
		okFlags.Right()
		okFlags.TripleBorder(wx.ALL)
		topsizer.AddWindowFlags(ok.Ptr(), okFlags.Ptr())
		okFlags.Delete()
	}

	dlg.SetSizer(topsizer.Ptr(), true)
	topsizer.Fit(dlg.Ptr())

	dlg.ShowModal()
	dlg.Delete()
}

// sampleXPM is the string array of the wxWidgets sample icon
// (samples/sample.xpm), used as the frame icon. It starts at the size line
// exactly like the C++ sample_xpm[] array: the XPM reader expects element 0 to
// be "width height colors chars_per_pixel".
var sampleXPM = []string{
	"32 32 10 1",
	" \tc None",
	".\tc #2020FF",
	"+\tc #8080FF",
	"@\tc #C0C0FF",
	"#\tc #FFFF00",
	"$\tc #FFFF80",
	"%\tc #FF2020",
	"&\tc #FF4848",
	"*\tc #FFFFDF",
	"=\tc #FFC0C0",
	"                                ",
	"                                ",
	"                                ",
	"                                ",
	"                                ",
	"        ..............          ",
	"        .++++++++++++.          ",
	"        .++++++++++++.          ",
	"        .+@@+++++++++.          ",
	"        .+@@+++++++++.          ",
	"        .+@@+++++++++.          ",
	"        .+@@++++++##############",
	"        .+@@++++++#$$$$$$$$$$$$#",
	"%%%%%%%%.+@@++++++#$$$$$$$$$$$$#",
	"%&&&&&&&.+@@++++++#$**$$$$$$$$$#",
	"%&&&&&&&.+@@++++++#$**$$$$$$$$$#",
	"%&==&&&&.+++++++++#$**$$$$$$$$$#",
	"%&==&&&&.+++++++++#$**$$$$$$$$$#",
	"%&==&&&&..........#$**$$$$$$$$$#",
	"%&==&&&&&&&&&%    #$**$$$$$$$$$#",
	"%&==&&&&&&&&&%    #$**$$$$$$$$$#",
	"%&==&&&&&&&&&%    #$**$$$$$$$$$#",
	"%&==&&&&&&&&&%    #$$$$$$$$$$$$#",
	"%&==&&&&&&&&&%    #$$$$$$$$$$$$#",
	"%&&&&&&&&&&&&%    ##############",
	"%&&&&&&&&&&&&%                  ",
	"%%%%%%%%%%%%%%                  ",
	"                                ",
	"                                ",
	"                                ",
	"                                ",
	"                                ",
}
