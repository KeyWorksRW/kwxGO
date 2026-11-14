// samples/preferences — kwxGO port of wxSamples/src/preferences/preferences.cpp
//
// Demonstrates wxPreferencesEditor: a frame that displays two settings values
// and a File menu whose "Preferences" item opens the platform preferences
// dialog, built from a stock "General" page and a custom "Topics" page.
//
// Conversion notes
// ----------------
//   - MyApp::m_settings / m_prefEditor / m_frame are package-level state here,
//     because the page and panel virtual hooks below are called back from C++
//     long after the construction functions return.
//   - PrefsPageGeneral derives from wxStockPreferencesPage(Kind_General); the
//     binding has no Kind_General symbol, so the documented enum value is used
//     (kwxFFI/include/kwx_hook_exports.h: "0 = General, 1 = Advanced").
//   - PrefsPageTopics derives from wxPreferencesPage, so it is built with
//     wx.NewKwxPreferencesPage; GetName()/GetIcon()/CreateWindow() are supplied
//     through the SetOn* hooks. The FFI callback ABI cannot return values, so
//     each hook publishes its result with the matching SetOn*Result() call.
//   - PrefsPageGeneralPanel / PrefsPageTopicsPanel derive from wxPanel only to
//     override TransferDataToWindow()/TransferDataFromWindow(). They are created
//     as wx.KwxPanel objects with those two hooks installed, so the dialog
//     drives them exactly as it drives the C++ overrides:
//     wxWindowBase::OnInitDialog() calls TransferDataToWindow() and
//     wxDialogBase::AcceptAndClose() calls TransferDataFromWindow() on OK.
//   - #ifdef __WXGTK__ -> GTKAllowDiagnosticsControl() is not ported; this
//     build targets wxMSW, where the branch does not compile.
package main

import (
	"syscall"
	"unsafe"

	"github.com/KeyWorksRW/kwxGO/wx"
)

// stockPageKindGeneral is wxStockPreferencesPage::Kind_General. kwxFFI exposes
// the stock-page kinds as a plain int ("0 = General, 1 = Advanced").
const stockPageKindGeneral = 0

// ---------------------------------------------------------------------------
// MySettings
// ---------------------------------------------------------------------------

// mySettings combines the settings edited in the preferences dialog.
type mySettings struct {
	useMarkdown bool
	spellcheck  bool
}

// newMySettings is MySettings::MySettings(): values are hardcoded because this
// sample does not use wxConfig (a real application would load them there).
func newMySettings() mySettings {
	return mySettings{useMarkdown: true, spellcheck: false}
}

// appSettings is MyApp::m_settings.
var appSettings = newMySettings()

// ---------------------------------------------------------------------------
// Application state
// ---------------------------------------------------------------------------

var (
	prefEditor   *wx.PreferencesEditor    // MyApp::m_prefEditor
	generalPage  *wx.StockPreferencesPage // PrefsPageGeneral
	topicsPage   *wx.KwxPreferencesPage   // PrefsPageTopics
	generalPanel *generalPanelState       // the live PrefsPageGeneralPanel
	topicsPanel  *topicsPanelState        // the live PrefsPageTopicsPanel
	mainFrame    *frameState              // MyApp::m_frame
)

// frameState corresponds to MyFrame.
type frameState struct {
	frame              *wx.Frame
	textMarkdownSyntax *wx.StaticText
	textSpellcheck     *wx.StaticText
}

// generalPanelState corresponds to PrefsPageGeneralPanel.
type generalPanelState struct {
	panel       *wx.KwxPanel
	useMarkdown *wx.CheckBox
	spellcheck  *wx.CheckBox
	// Settings corresponding to the current values in this dialog
	// (PrefsPageGeneralPanel::m_settingsCurrent).
	current mySettings
}

// topicsPanelState corresponds to PrefsPageTopicsPanel.
type topicsPanelState struct {
	panel    *wx.KwxPanel
	fulltext *wx.CheckBox
}

// ---------------------------------------------------------------------------
// C++ virtual hooks
//
// wxPreferencesPage::GetName/GetIcon/CreateWindow and the two wxPanel transfer
// methods are virtuals, so kwxFFI exposes them as SetOn* hooks that take a C
// function pointer. syscall.NewCallback is the standard library's supported
// conversion from a Go function to a Win32 callback pointer; every hook body
// below still calls only generated wx wrappers.
// ---------------------------------------------------------------------------

// cCallback builds the C function pointer a SetOn* hook expects. The
// uintptr -> unsafe.Pointer conversion is the same one wx/event.go uses for its
// handler handles.
func cCallback(fn any) unsafe.Pointer {
	return unsafe.Pointer(uintptr(syscall.NewCallback(fn)))
}

// The page panels hand an opaque context pointer back to their transfer hooks.
// A package-level byte is stable for the life of the process, holds no Go
// pointer, and distinguishes the two panels, which is all the hooks need.
var (
	generalPanelTag byte
	topicsPanelTag  byte
)

var (
	// wxStockPreferencesPage::CreateWindow(parent) for the General page.
	hookCreateGeneralWindow = cCallback(func(_ unsafe.Pointer, parent unsafe.Pointer) uintptr {
		if generalPage != nil && parent != nil {
			if panel := newGeneralPanel(parent); panel != nil {
				generalPage.SetOnCreateWindowResult(panel.Ptr())
			}
		}
		return 0
	})

	// wxPreferencesPage::CreateWindow(parent) for the Topics page.
	hookCreateTopicsWindow = cCallback(func(_ unsafe.Pointer, parent unsafe.Pointer) uintptr {
		if topicsPage != nil && parent != nil {
			if panel := newTopicsPanel(parent); panel != nil {
				topicsPage.SetOnCreateWindowResult(panel.Ptr())
			}
		}
		return 0
	})

	// PrefsPageTopics::GetName()
	hookGetTopicsName = cCallback(func(_ unsafe.Pointer) uintptr {
		if topicsPage != nil {
			topicsPage.SetOnGetNameResult("Topics")
		}
		return 0
	})

	// PrefsPageTopics::GetIcon() -> wxArtProvider::GetBitmapBundle(wxART_HELP,
	// wxART_TOOLBAR). The helper copies the bundle and deletes the original, so
	// ownership passes over and the wrapper must not be Delete()d.
	hookGetTopicsIcon = cCallback(func(_ unsafe.Pointer) uintptr {
		if topicsPage == nil {
			return 0
		}
		if bundle := wx.ArtProvider_GetBitmapBundle(wx.ART_HELP, wx.ART_TOOLBAR, -1, -1); bundle != nil {
			topicsPage.SetOnGetIconResult(bundle.Ptr())
		}
		return 0
	})

	// wxPanel::TransferDataToWindow() on either page panel.
	hookTransferToWindow = cCallback(func(pObject unsafe.Pointer) uintptr {
		switch pObject {
		case unsafe.Pointer(&generalPanelTag):
			return transferToGeneralWindow()
		case unsafe.Pointer(&topicsPanelTag):
			return transferToTopicsWindow()
		}
		return 1
	})

	// wxPanel::TransferDataFromWindow() on either page panel.
	hookTransferFromWindow = cCallback(func(pObject unsafe.Pointer) uintptr {
		switch pObject {
		case unsafe.Pointer(&generalPanelTag):
			return transferFromGeneralWindow()
		case unsafe.Pointer(&topicsPanelTag):
			return transferFromTopicsWindow()
		}
		return 1
	})
)

// ---------------------------------------------------------------------------
// Sizer helpers
// ---------------------------------------------------------------------------

// boxSizer / flexGridSizer hand out the embedded base sizer, which is where the
// wxSizerFlags overloads live.
func boxSizer(o *wx.BoxSizer) *wx.Sizer           { return &o.Sizer }
func flexGridSizer(o *wx.FlexGridSizer) *wx.Sizer { return &o.GridSizer.Sizer }

// addToSizer is sizer->Add(window, wxSizerFlags(proportion).<setup>()).
func addToSizer(sizer *wx.Sizer, window unsafe.Pointer, proportion int, setup func(*wx.SizerFlags)) {
	flags := wx.NewSizerFlags(proportion)
	if flags == nil {
		sizer.AddWindow(window, proportion, 0, 0, nil)
		return
	}
	defer flags.Delete()
	if setup != nil {
		setup(flags)
	}
	sizer.AddWindowFlags(window, flags.Ptr())
}

func checkValue(v bool) int {
	if v {
		return 1
	}
	return 0
}

// ---------------------------------------------------------------------------
// Settings display (MyFrame::UpdateSettings and MyApp::UpdateSettings)
// ---------------------------------------------------------------------------

// displaySettings refreshes the values shown in the frame.
func displaySettings(state *frameState) {
	if state == nil {
		return
	}
	if appSettings.useMarkdown {
		state.textMarkdownSyntax.SetLabel("yes")
	} else {
		state.textMarkdownSyntax.SetLabel("no")
	}
	if appSettings.spellcheck {
		state.textSpellcheck.SetLabel("on")
	} else {
		state.textSpellcheck.SetLabel("off")
	}
}

// updateSettings is MyApp::UpdateSettings(): store the values, then refresh the
// frame.
func updateSettings(settings mySettings) {
	appSettings = settings
	displaySettings(mainFrame)
}

// updateSettingsIfNecessary is PrefsPageGeneralPanel::UpdateSettingsIfNecessary():
// changes are applied immediately only on platforms whose preferences dialog
// works that way (macOS, GTK+); on wxMSW they wait for the OK button.
func updateSettingsIfNecessary() {
	if generalPanel != nil && wx.PreferencesEditor_ShouldApplyChangesImmediately() {
		updateSettings(generalPanel.current)
	}
}

// ---------------------------------------------------------------------------
// PrefsPageGeneralPanel
// ---------------------------------------------------------------------------

// newGeneralPanel is PrefsPageGeneralPanel::PrefsPageGeneralPanel().
func newGeneralPanel(parent unsafe.Pointer) *wx.KwxPanel {
	panel := wx.NewKwxPanel(unsafe.Pointer(&generalPanelTag), parent, wx.ID_ANY,
		-1, -1, -1, -1, wx.TAB_TRAVERSAL|wx.BORDER_NONE)
	if panel == nil {
		return nil
	}

	useMarkdown := wx.NewCheckBox(panel.Ptr(), wx.ID_ANY, "Use Markdown syntax", -1, -1, -1, -1, 0)
	spellcheck := wx.NewCheckBox(panel.Ptr(), wx.ID_ANY, "Check spelling", -1, -1, -1, -1, 0)
	if useMarkdown == nil || spellcheck == nil {
		return nil
	}

	sizer := wx.NewBoxSizer(wx.VERTICAL)
	if sizer == nil {
		return nil
	}
	base := boxSizer(sizer)
	addToSizer(base, useMarkdown.Ptr(), 0, func(f *wx.SizerFlags) { f.Border(wx.ALL) })
	addToSizer(base, spellcheck.Ptr(), 0, func(f *wx.SizerFlags) { f.Border(wx.ALL) })
	panel.SetSizerAndFit(sizer.Ptr(), true)

	// The two virtual overrides. The dialog calls them; see the notes at the
	// top of the file.
	panel.SetOnTransferDataToWindow(hookTransferToWindow)
	panel.SetOnTransferDataFromWindow(hookTransferFromWindow)

	state := &generalPanelState{
		panel:       panel,
		useMarkdown: useMarkdown,
		spellcheck:  spellcheck,
	}
	generalPanel = state

	useMarkdown.Bind(wx.EVT_CHECKBOX, func(evt *wx.Event) {
		if evt == nil {
			return
		}
		cmd := wx.CommandEventFromPtr(evt.Ptr())
		if cmd == nil || generalPanel == nil {
			return
		}
		generalPanel.current.useMarkdown = cmd.IsChecked()
		updateSettingsIfNecessary()
	})

	spellcheck.Bind(wx.EVT_CHECKBOX, func(evt *wx.Event) {
		if evt == nil {
			return
		}
		cmd := wx.CommandEventFromPtr(evt.Ptr())
		if cmd == nil || generalPanel == nil {
			return
		}
		generalPanel.current.spellcheck = cmd.IsChecked()
		updateSettingsIfNecessary()
	})

	return panel
}

// transferToGeneralWindow is PrefsPageGeneralPanel::TransferDataToWindow().
func transferToGeneralWindow() uintptr {
	state := generalPanel
	if state == nil {
		return 1
	}
	state.current = appSettings
	state.useMarkdown.SetValue(checkValue(state.current.useMarkdown))
	state.spellcheck.SetValue(checkValue(state.current.spellcheck))
	return 1
}

// transferFromGeneralWindow is
// PrefsPageGeneralPanel::TransferDataFromWindow(): called on platforms with a
// modal preferences dialog to save and apply the changes.
func transferFromGeneralWindow() uintptr {
	if state := generalPanel; state != nil {
		updateSettings(state.current)
	}
	return 1
}

// ---------------------------------------------------------------------------
// PrefsPageGeneral : wxStockPreferencesPage(Kind_General)
// ---------------------------------------------------------------------------

func newGeneralPage() *wx.StockPreferencesPage {
	page := wx.NewStockPreferencesPage(stockPageKindGeneral)
	if page == nil {
		return nil
	}
	// PrefsPageGeneral::CreateWindow() override.
	page.SetOnCreateWindow(hookCreateGeneralWindow)
	return page
}

// ---------------------------------------------------------------------------
// PrefsPageTopicsPanel
// ---------------------------------------------------------------------------

// newTopicsPanel is PrefsPageTopicsPanel::PrefsPageTopicsPanel().
func newTopicsPanel(parent unsafe.Pointer) *wx.KwxPanel {
	panel := wx.NewKwxPanel(unsafe.Pointer(&topicsPanelTag), parent, wx.ID_ANY,
		-1, -1, -1, -1, wx.TAB_TRAVERSAL|wx.BORDER_NONE)
	if panel == nil {
		return nil
	}

	sizer := wx.NewBoxSizer(wx.VERTICAL)
	if sizer == nil {
		return nil
	}
	base := boxSizer(sizer)

	label := wx.NewStaticText(panel.Ptr(), wx.ID_ANY, "Search in these topics:", -1, -1, -1, -1, 0)
	if label == nil {
		return nil
	}
	addToSizer(base, label.Ptr(), 0, func(f *wx.SizerFlags) { f.Border(wx.ALL) })

	// wxListBox(this, wxID_ANY) has no entries yet (the TODO in the C++).
	box := wx.NewListBox(panel.Ptr(), wx.ID_ANY, -1, -1, -1, -1, nil, wx.ExpLB_DEFAULT())
	if box == nil {
		return nil
	}
	box.SetMinSize(400, 300)
	addToSizer(base, box.Ptr(), 1, func(f *wx.SizerFlags) {
		f.Border(wx.ALL)
		f.Expand()
	})

	fulltext := wx.NewCheckBox(panel.Ptr(), wx.ID_ANY, "Automatically build fulltext index", -1, -1, -1, -1, 0)
	if fulltext == nil {
		return nil
	}
	addToSizer(base, fulltext.Ptr(), 0, func(f *wx.SizerFlags) { f.Border(wx.ALL) })

	panel.SetSizerAndFit(sizer.Ptr(), true)

	topicsPanel = &topicsPanelState{panel: panel, fulltext: fulltext}

	panel.SetOnTransferDataToWindow(hookTransferToWindow)
	panel.SetOnTransferDataFromWindow(hookTransferFromWindow)

	// ChangedFulltext() is bound only on platforms that apply changes
	// immediately; its C++ body holds a comment only.
	if wx.PreferencesEditor_ShouldApplyChangesImmediately() {
		fulltext.Bind(wx.EVT_CHECKBOX, func(evt *wx.Event) {
			if evt == nil {
				return
			}
			// PrefsPageTopicsPanel::ChangedFulltext(): save the new value and
			// apply it to the application here.
		})
	}

	return panel
}

// transferToTopicsWindow is PrefsPageTopicsPanel::TransferDataToWindow(): the
// values are hardcoded because this sample does not use wxConfig.
func transferToTopicsWindow() uintptr {
	if topicsPanel != nil {
		topicsPanel.fulltext.SetValue(1)
	}
	return 1
}

// transferFromTopicsWindow is PrefsPageTopicsPanel::TransferDataFromWindow().
// ChangedFulltext(dummy) is a comment-only body in the C++, and the listbox is
// still a TODO, so nothing is transferred back.
func transferFromTopicsWindow() uintptr {
	return 1
}

// ---------------------------------------------------------------------------
// PrefsPageTopics : wxPreferencesPage
// ---------------------------------------------------------------------------

func newTopicsPage() *wx.KwxPreferencesPage {
	page := wx.NewKwxPreferencesPage(nil)
	if page == nil {
		return nil
	}
	// The editor takes ownership of the page, so this wrapper is never Delete()d.
	page.SetOnGetName(hookGetTopicsName)
	page.SetOnGetIcon(hookGetTopicsIcon)
	page.SetOnCreateWindow(hookCreateTopicsWindow)
	return page
}

// ---------------------------------------------------------------------------
// MyApp::ShowPreferencesEditor / MyApp::DismissPreferencesEditor
// ---------------------------------------------------------------------------

func showPreferencesEditor(parent unsafe.Pointer) {
	if prefEditor == nil {
		editor := wx.NewPreferencesEditor("")
		if editor == nil {
			return
		}

		// AddPage(new PrefsPageGeneral());
		if page := newGeneralPage(); page != nil {
			generalPage = page
			editor.AddStockPage(page.Ptr())
		}

		// AddPage(new PrefsPageTopics());
		if page := newTopicsPage(); page != nil {
			topicsPage = page
			editor.AddPage(page.Ptr())
		}

		prefEditor = editor
	}

	// On wxMSW Show() creates the dialog, runs it modally and destroys it again.
	prefEditor.Show(parent)
}

func dismissPreferencesEditor() {
	if prefEditor != nil {
		prefEditor.Dismiss()
	}
}

// ---------------------------------------------------------------------------
// MyFrame
// ---------------------------------------------------------------------------

// newMyFrame is MyFrame::MyFrame().
func newMyFrame() *frameState {
	frame := wx.NewFrame(nil, wx.ID_ANY, "Preferences sample", -1, -1, -1, -1, wx.DEFAULT_FRAME_STYLE)
	if frame == nil {
		return nil
	}

	fileMenu := wx.NewMenu("", 0)
	if fileMenu == nil {
		return nil
	}
	// fileMenu->Append(wxID_PREFERENCES) / Append(wxID_EXIT): the one-argument
	// overload, which keeps the stock label, accelerator and help text.
	fileMenu.AppendById(wx.ID_PREFERENCES)
	fileMenu.AppendById(wx.ID_EXIT)

	menuBar := wx.NewMenuBar(0)
	if menuBar == nil {
		return nil
	}
	menuBar.Append(fileMenu.Ptr(), "&File")
	frame.SetMenuBar(menuBar.Ptr())

	frame.BindId(wx.EVT_MENU, wx.ID_PREFERENCES, func(evt *wx.Event) {
		if evt == nil {
			return
		}
		showPreferencesEditor(frame.Ptr())
	})

	frame.BindId(wx.EVT_MENU, wx.ID_EXIT, func(evt *wx.Event) {
		if evt == nil {
			return
		}
		frame.Close(false)
	})

	frame.Bind(wx.EVT_CLOSE_WINDOW, func(evt *wx.Event) {
		if evt == nil {
			return
		}
		dismissPreferencesEditor()
		// MyFrame::OnClose() = wxGetApp().DismissPreferencesEditor(); e.Skip();
		// The skip lets wx's default close handler run, which closes the frame
		// and (exit-on-frame-delete) ends the main loop. wxWindow::Destroy()
		// defers the real deletion to idle time, so the frame is not destroyed
		// inside this handler.
		evt.Skip(true)
	})

	panel := wx.NewPanel(frame.Ptr(), wx.ID_ANY, -1, -1, -1, -1, wx.TAB_TRAVERSAL|wx.BORDER_NONE)
	if panel == nil {
		return nil
	}

	textMarkdownSyntax := wx.NewStaticText(panel.Ptr(), wx.ID_ANY, "", -1, -1, -1, -1, 0)
	textSpellcheck := wx.NewStaticText(panel.Ptr(), wx.ID_ANY, "", -1, -1, -1, -1, 0)
	if textMarkdownSyntax == nil || textSpellcheck == nil {
		return nil
	}

	// new wxFlexGridSizer(2, wxSize(5, 5)): two columns, 5px gaps.
	sizer := wx.NewFlexGridSizer(0, 2, 5, 5)
	if sizer == nil {
		return nil
	}
	base := flexGridSizer(sizer)

	labelMarkdown := wx.NewStaticText(panel.Ptr(), wx.ID_ANY, "Markdown syntax:", -1, -1, -1, -1, 0)
	if labelMarkdown == nil {
		return nil
	}
	// sizer->Add(label, wxSizerFlags().Center().Right())
	addToSizer(base, labelMarkdown.Ptr(), 0, func(f *wx.SizerFlags) {
		f.Center()
		f.Right()
	})
	// sizer->Add(m_textMarkdownSyntax, wxSizerFlags().Center())
	addToSizer(base, textMarkdownSyntax.Ptr(), 0, func(f *wx.SizerFlags) { f.Center() })

	labelSpellcheck := wx.NewStaticText(panel.Ptr(), wx.ID_ANY, "Spell checking:", -1, -1, -1, -1, 0)
	if labelSpellcheck == nil {
		return nil
	}
	addToSizer(base, labelSpellcheck.Ptr(), 0, func(f *wx.SizerFlags) {
		f.Center()
		f.Right()
	})
	addToSizer(base, textSpellcheck.Ptr(), 0, func(f *wx.SizerFlags) { f.Center() })

	panel.SetSizer(sizer.Ptr(), true)

	state := &frameState{
		frame:              frame,
		textMarkdownSyntax: textMarkdownSyntax,
		textSpellcheck:     textSpellcheck,
	}

	// Show the initial values.
	displaySettings(state)

	return state
}

// ---------------------------------------------------------------------------
// MyApp::OnInit
// ---------------------------------------------------------------------------

func main() {
	wx.Run(func() bool {
		// SetAppDisplayName("wxWidgets Sample") is used in the preferences
		// dialog title, so it must not stay the default "Preferences".
		if app := wx.AppFromPtr(wx.App_GetApp()); app != nil {
			app.SetAppDisplayName("wxWidgets Sample")
		}

		state := newMyFrame()
		if state == nil {
			return false
		}
		mainFrame = state

		state.frame.Show()
		wx.SetTopWindow(state.frame.Ptr())
		return true
	})

	// MyApp owns the editor and the pages; wx.Run() already destroyed the
	// application, so the natives are released with the process.
}
