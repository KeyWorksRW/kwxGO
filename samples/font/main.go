// Command font is the kwxGO port of the wxWidgets "font" sample
// (samples/font/font.cpp).
//
// The frame holds a splitter: the upper pane is a FontWindow with a FontPanel
// (face name, family, point size, style, weight, underline, strikethrough and
// fixed-width controls) and a FontCanvas that paints characters 32..255 of the
// current font in a 7x32 grid; the lower pane is a multiline text control that
// is shown with the selected font. The three menus change the font through
// standard wxWidgets font objects, enumerate font families and encodings, and
// exercise the native font descriptor strings.
package main

import (
	"fmt"
	"sort"
	"strings"
	"unsafe"

	"github.com/KeyWorksRW/kwxGO/wx"
)

// Dimensions of the original sample.xpm icon, now sample.svg.
const (
	sampleIconWidth  = 32
	sampleIconHeight = 32
)

const sampleTitle = "wxWidgets Font Sample"

// Menu and command IDs, mirroring the anonymous enum in font.cpp. The two
// standard IDs are runtime values in the binding, so the whole block is a var
// declaration and keeps the C++ ordering relative to wxID_HIGHEST.
var (
	idFontQuit  = wx.ID_EXIT
	idFontAbout = wx.ID_ABOUT

	idFontViewMsg        = wx.ID_HIGHEST
	idFontTestTextValue  = idFontViewMsg + 1
	idFontIncSize        = idFontViewMsg + 2
	idFontDecSize        = idFontViewMsg + 3
	idFontGetBaseFont    = idFontViewMsg + 4
	idFontBold           = idFontViewMsg + 5
	idFontLight          = idFontViewMsg + 6
	idFontItalic         = idFontViewMsg + 7
	idFontSlant          = idFontViewMsg + 8
	idFontUnderlined     = idFontViewMsg + 9
	idFontStrikethrough  = idFontViewMsg + 10
	idFontNormal         = idFontViewMsg + 11
	idFontSmall          = idFontViewMsg + 12
	idFontItalicStock    = idFontViewMsg + 13
	idFontSwiss          = idFontViewMsg + 14
	idFontDefault        = idFontViewMsg + 15
	idFontStandard       = idFontViewMsg + 16
	idFontSysOEMFixed    = idFontViewMsg + 17
	idFontSysANSIFixed   = idFontViewMsg + 18
	idFontSysANSIVar     = idFontViewMsg + 19
	idFontSysSystem      = idFontViewMsg + 20
	idFontSysDevice      = idFontViewMsg + 21
	idFontSysDefaultGUI  = idFontViewMsg + 22
	idFontSystemSettings = idFontViewMsg + 23
)

const (
	idFontChoose                  = 100
	idFontEnumFamiliesForEncoding = 101
	idFontEnumFamilies            = 102
	idFontEnumFixedFamilies       = 103
	idFontEnumEncodings           = 104
	idFontSetNativeDesc           = 105
	idFontSetNativeUserDesc       = 106
	idFontSetFamily               = 107
	idFontSetFaceName             = 108
	idFontSetEncoding             = 109
	idFontPrivate                 = 110
	idFontMax                     = 111
)

// familyNames/families and friends mirror the Family, Style and Weight enums.
var (
	familyNames = []string{"Default", "Decorative", "Roman", "Script", "Swiss", "Modern", "Teletype"}
	styleNames  = []string{"Normal", "Italic", "Slant"}
	weightNames = []string{
		"Thin", "Extra light", "Light", "Normal", "Medium",
		"Semi-bold", "Bold", "Extra bold", "Heavy", "Extra heavy",
	}
)

func families() []int {
	return []int{
		wx.FONTFAMILY_DEFAULT, wx.FONTFAMILY_DECORATIVE, wx.FONTFAMILY_ROMAN,
		wx.FONTFAMILY_SCRIPT, wx.FONTFAMILY_SWISS, wx.FONTFAMILY_MODERN,
		wx.FONTFAMILY_TELETYPE,
	}
}

func weights() []int {
	return []int{
		wx.FONTWEIGHT_THIN, wx.FONTWEIGHT_EXTRALIGHT, wx.FONTWEIGHT_LIGHT,
		wx.FONTWEIGHT_NORMAL, wx.FONTWEIGHT_MEDIUM, wx.FONTWEIGHT_SEMIBOLD,
		wx.FONTWEIGHT_BOLD, wx.FONTWEIGHT_EXTRABOLD, wx.FONTWEIGHT_HEAVY,
		wx.FONTWEIGHT_EXTRAHEAVY,
	}
}

// fontFrame holds the state MyFrame keeps in member variables.
type fontFrame struct {
	frame      *wx.Frame
	icon       *wx.Icon
	fontWindow *fontWindow
	textCtrl   *wx.TextCtrl

	// Static state of MyFrame::OnViewMsg.
	lastDir  string
	lastFile string
}

// fontWindow is FontWindow: the FontPanel, a separator and the FontCanvas.
type fontWindow struct {
	win    *wx.Window
	panel  *fontPanel
	canvas *fontCanvas
}

// fontPanel is FontPanel: the controls that describe a font.
type fontPanel struct {
	win *wx.Panel

	textFaceName       *wx.TextCtrl
	choiceFamily       *wx.Choice
	spinPointSize      *wx.SpinCtrlDouble
	choiceStyle        *wx.Choice
	choiceWeight       *wx.Choice
	spinWeight         *wx.SpinCtrl
	checkUnderlined    *wx.CheckBox
	checkStrikethrough *wx.CheckBox
	checkFixedWidth    *wx.CheckBox
	labelInfo          *wx.StaticText

	font             *wx.Font
	useFamily        bool
	useNumericWeight bool
}

// fontCanvas is FontCanvas: it paints the font's character map.
type fontCanvas struct {
	win    *wx.Window
	font   *wx.Font
	colour *wx.Colour
}

func main() {
	wx.Run(func() bool {
		f := newFontFrame()
		if f == nil {
			return false
		}
		f.frame.Show()
		wx.SetTopWindow(f.frame.Ptr())
		return true
	})
}

// cloneFont makes an owned deep copy of a font. The returned font must be
// deleted by the caller.
func cloneFont(font *wx.Font) *wx.Font {
	if font == nil {
		return nil
	}
	return wx.NewFontFromFont(font.Ptr())
}

// copyColour makes an owned copy of a colour from an opaque wxColour pointer.
func copyColour(ptr unsafe.Pointer) *wx.Colour {
	if ptr == nil {
		return nil
	}
	col := wx.NewColourEmpty()
	if col == nil {
		return nil
	}
	col.Copy(ptr)
	return col
}

// stockFont clones one of the global wxWidgets font objects.
func stockFont(ptr unsafe.Pointer) *wx.Font {
	return wx.NewFontFromFont(ptr)
}

// windowTextExtent measures text with the window's current font, like
// wxWindow::GetTextExtent().
func windowTextExtent(winPtr unsafe.Pointer, text string) wx.Size {
	dc := wx.NewClientDC(winPtr)
	if dc == nil {
		return wx.Size{}
	}
	defer dc.Delete()

	font := wx.NewFontDefault()
	if font != nil {
		if win := wx.WindowFromPtr(winPtr); win != nil {
			win.GetFont(font.Ptr())
			dc.SetFont(font.Ptr())
		}
		font.Delete()
	}
	return dc.GetTextExtentSize(text)
}

// loadSampleIcon loads the SVG rendition of the sample's XPM icon and converts
// it to a wxIcon for the frame.
func loadSampleIcon() *wx.Icon {
	bundle := wx.BitmapBundle_FromSVGFile("sample.svg", sampleIconWidth, sampleIconHeight)
	if bundle == nil {
		return nil
	}
	if bundle == nil || !bundle.IsOk() {
		return nil
	}

	bitmap := bundle.GetBitmap(sampleIconWidth, sampleIconHeight)
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

// newFontFrame is MyFrame::MyFrame().
func newFontFrame() *fontFrame {
	frame := wx.NewFrame(nil, wx.ID_ANY, "wxWidgets font sample",
		-1, -1, -1, -1, wx.DEFAULT_FRAME_STYLE)
	if frame == nil {
		return nil
	}

	f := &fontFrame{frame: frame}

	if icon := loadSampleIcon(); icon != nil {
		f.icon = icon
		frame.SetIcon(icon.Ptr())
	}

	f.buildMenus()

	splitter := wx.NewSplitterWindow(frame.Ptr(), wx.ID_ANY, -1, -1, -1, -1, wx.SP_3D)
	if splitter == nil {
		return nil
	}

	f.fontWindow = newFontWindow(splitter.Ptr())
	if f.fontWindow == nil {
		return nil
	}
	f.fontWindow.win.Bind(wx.EVT_BUTTON, f.onFontPanelApply)

	f.textCtrl = wx.NewTextCtrl(splitter.Ptr(), wx.ID_ANY,
		"Paste text here to see how it looks\nlike in the given font",
		-1, -1, -1, 6*frame.GetCharHeight(), wx.TE_MULTILINE)
	if f.textCtrl == nil {
		return nil
	}

	splitter.SplitHorizontally(f.fontWindow.win.Ptr(), f.textCtrl.Ptr(), 0)

	frame.CreateStatusBar(1, wx.STB_DEFAULT_STYLE, wx.ID_ANY)
	frame.SetStatusText("Welcome to wxWidgets font demo!", 0)

	best := splitter.GetBestSize()
	frame.SetClientSize(best.Width, best.Height)
	splitter.SetSashPosition(f.fontWindow.win.GetBestSize().Height, true)

	return f
}

// buildMenus creates the File, Font and Select menus and installs the handlers.
func (f *fontFrame) buildMenus() {
	menuFile := wx.NewMenu("", 0)
	menuFile.Append(idFontTestTextValue, "&Test text value",
		"Verify that getting and setting text value doesn't change it", false)
	menuFile.Append(idFontViewMsg, "&View...\tCtrl-V", "View an email message file", false)
	menuFile.AppendSeparator()
	menuFile.Append(idFontAbout, "&About\tCtrl-A", "Show about dialog", false)
	menuFile.AppendSeparator()
	menuFile.Append(idFontQuit, "E&xit\tAlt-X", "Quit this program", false)

	menuFont := wx.NewMenu("", 0)
	menuFont.Append(idFontIncSize, "&Increase font size by 2 points\tCtrl-I", "", false)
	menuFont.Append(idFontDecSize, "&Decrease font size by 2 points\tCtrl-D", "", false)
	menuFont.Append(idFontGetBaseFont, "Use &base version of the font\tCtrl-0", "", false)
	menuFont.AppendSeparator()
	menuFont.AppendCheckItem(idFontBold, "&Bold\tCtrl-B", "Toggle bold state")
	menuFont.AppendCheckItem(idFontLight, "&Light\tCtrl-L", "Toggle light state")
	menuFont.AppendSeparator()
	menuFont.AppendCheckItem(idFontItalic, "&Oblique\tCtrl-O", "Toggle italic state")
	// Under wxMSW slant == italic, so the original omits the extra menu item.
	menuFont.AppendSeparator()
	menuFont.AppendCheckItem(idFontUnderlined, "&Underlined\tCtrl-U", "Toggle underlined state")
	menuFont.AppendCheckItem(idFontStrikethrough, "&Strikethrough", "Toggle strikethrough state")
	menuFont.AppendSeparator()
	menuFont.Append(idFontSetNativeDesc, "Set native font &description\tShift-Ctrl-D", "", false)
	menuFont.Append(idFontSetNativeUserDesc, "Set &user font description\tShift-Ctrl-U", "", false)
	menuFont.AppendSeparator()
	menuFont.Append(idFontSetFamily, "Set font family", "", false)
	menuFont.Append(idFontSetFaceName, "Set font face name", "", false)
	menuFont.Append(idFontSetEncoding, "Set font &encoding\tShift-Ctrl-E", "", false)

	menuSelect := wx.NewMenu("", 0)
	menuSelect.Append(idFontChoose, "&Select font...\tCtrl-S", "Select a standard font", false)

	menuStdFonts := wx.NewMenu("", 0)
	menuStdFonts.Append(idFontNormal, "wxNORMAL_FONT", "Normal font used by wxWidgets", false)
	menuStdFonts.Append(idFontSmall, "wxSMALL_FONT", "Small font used by wxWidgets", false)
	menuStdFonts.Append(idFontItalicStock, "wxITALIC_FONT", "Italic font used by wxWidgets", false)
	menuStdFonts.Append(idFontSwiss, "wxSWISS_FONT", "Swiss font used by wxWidgets", false)
	menuStdFonts.Append(idFontDefault, "wxFont()", "wxFont constructed from default wxFontInfo", false)
	menuSelect.AppendSub(idFontStandard, "Standar&d fonts", menuStdFonts.Ptr(), "")

	menuSettingFonts := wx.NewMenu("", 0)
	menuSettingFonts.Append(idFontSysOEMFixed, "wxSYS_OEM_FIXED_FONT",
		"Original equipment manufacturer dependent fixed-pitch font.", false)
	menuSettingFonts.Append(idFontSysANSIFixed, "wxSYS_ANSI_FIXED_FONT",
		"Windows fixed-pitch (monospaced) font. ", false)
	menuSettingFonts.Append(idFontSysANSIVar, "wxSYS_ANSI_VAR_FONT",
		"Windows variable-pitch (proportional) font.", false)
	menuSettingFonts.Append(idFontSysSystem, "wxSYS_SYSTEM_FONT", "System font.", false)
	menuSettingFonts.Append(idFontSysDevice, "wxSYS_DEVICE_DEFAULT_FONT", "Device-dependent font.", false)
	menuSettingFonts.Append(idFontSysDefaultGUI, "wxSYS_DEFAULT_GUI_FONT",
		"Default font for user interface objects such as menus and dialog boxes. ", false)
	menuSelect.AppendSub(idFontSystemSettings, "System fonts", menuSettingFonts.Ptr(), "")

	menuSelect.AppendSeparator()
	menuSelect.Append(idFontEnumFamilies, "Enumerate font &families\tCtrl-F", "", false)
	menuSelect.Append(idFontEnumFixedFamilies, "Enumerate fi&xed font families\tCtrl-X", "", false)
	menuSelect.Append(idFontEnumEncodings, "Enumerate &encodings\tCtrl-E", "", false)
	menuSelect.Append(idFontEnumFamiliesForEncoding, "Find font for en&coding...\tCtrl-C",
		"Find font families for given encoding", false)

	// Try to use a private font from the current directory.
	privfont := "wxprivate.ttf"
	if wx.Font_AddPrivateFont(privfont) {
		menuSelect.AppendSeparator()
		menuSelect.Append(idFontPrivate, "Select private font",
			"Select a font available only in this application", false)
	} else {
		wx.LogWarning(fmt.Sprintf("Failed to add private font from \"%s\"", privfont))
	}

	menuBar := wx.NewMenuBar(0)
	menuBar.Append(menuFile.Ptr(), "&File")
	menuBar.Append(menuFont.Ptr(), "F&ont")
	menuBar.Append(menuSelect.Ptr(), "&Select")
	f.frame.SetMenuBar(menuBar.Ptr())

	bind := func(id int, fn wx.EventHandler) {
		f.frame.BindId(wx.EVT_MENU, id, fn)
	}

	bind(idFontQuit, f.onQuit)
	bind(idFontTestTextValue, f.onTestTextValue)
	bind(idFontViewMsg, f.onViewMsg)
	bind(idFontAbout, f.onAbout)

	bind(idFontGetBaseFont, f.onGetBaseFont)
	bind(idFontIncSize, func(evt *wx.Event) { f.doResizeFont(2) })
	bind(idFontDecSize, func(evt *wx.Event) { f.doResizeFont(-2) })

	bind(idFontBold, f.onBold)
	bind(idFontLight, f.onLight)
	bind(idFontItalic, f.onItalic)
	bind(idFontUnderlined, f.onUnderline)
	bind(idFontStrikethrough, f.onStrikethrough)

	bind(idFontNormal, f.onStockFont)
	bind(idFontSmall, f.onStockFont)
	bind(idFontItalicStock, f.onStockFont)
	bind(idFontSwiss, f.onStockFont)
	bind(idFontDefault, f.onFontDefault)

	bind(idFontSysOEMFixed, f.onSystemSettingsFont)
	bind(idFontSysANSIFixed, f.onSystemSettingsFont)
	bind(idFontSysANSIVar, f.onSystemSettingsFont)
	bind(idFontSysSystem, f.onSystemSettingsFont)
	bind(idFontSysDevice, f.onSystemSettingsFont)
	bind(idFontSysDefaultGUI, f.onSystemSettingsFont)

	bind(idFontSetNativeDesc, f.onSetNativeDesc)
	bind(idFontSetNativeUserDesc, f.onSetNativeUserDesc)
	bind(idFontSetFamily, f.onSetFamily)
	bind(idFontSetFaceName, f.onSetFaceName)
	bind(idFontSetEncoding, f.onSetEncoding)

	bind(idFontChoose, f.onSelectFont)
	bind(idFontEnumFamiliesForEncoding, f.onEnumerateFamiliesForEncoding)
	bind(idFontEnumFamilies, func(evt *wx.Event) { f.doEnumerateFamilies(false, wx.FONTENCODING_SYSTEM, false) })
	bind(idFontEnumFixedFamilies, func(evt *wx.Event) { f.doEnumerateFamilies(true, wx.FONTENCODING_SYSTEM, false) })
	bind(idFontEnumEncodings, f.onEnumerateEncodings)
	bind(idFontPrivate, f.onPrivateFont)
}

// newFontWindow is FontWindow::FontWindow().
func newFontWindow(parent unsafe.Pointer) *fontWindow {
	win := wx.NewWindow(parent, wx.ID_ANY, -1, -1, -1, -1, 0)
	if win == nil {
		return nil
	}

	fw := &fontWindow{win: win}

	fw.panel = newFontPanel(win.Ptr())
	if fw.panel == nil {
		return nil
	}

	line := wx.NewStaticLine(win.Ptr(), wx.ID_ANY, -1, -1, -1, -1, 0)
	if line == nil {
		return nil
	}

	fw.canvas = newFontCanvas(win.Ptr())
	if fw.canvas == nil {
		return nil
	}

	sizer := wx.NewBoxSizer(wx.VERTICAL)
	addWindow := func(w unsafe.Pointer, prop int, cfg func(*wx.SizerFlags)) {
		flags := wx.NewSizerFlags(prop)
		if cfg != nil {
			cfg(flags)
		}
		sizer.AddWindowFlags(w, flags.Ptr())
		flags.Delete()
	}

	addWindow(fw.panel.win.Ptr(), 0, func(s *wx.SizerFlags) {
		s.Expand()
		s.Border(wx.ALL)
	})
	addWindow(line.Ptr(), 0, func(s *wx.SizerFlags) { s.Expand() })
	addWindow(fw.canvas.win.Ptr(), 1, func(s *wx.SizerFlags) { s.Expand() })

	win.SetSizer(sizer.Ptr(), false)
	return fw
}

// updateFont is FontWindow::UpdateFont().
func (fw *fontWindow) updateFont(font *wx.Font, colour *wx.Colour) {
	fw.panel.showFont(font)

	fw.canvas.setTextFont(font)
	if colour != nil {
		fw.canvas.setColour(colour)
	}
	fw.canvas.win.Refresh(false)
}

// getTextFont returns the canvas font. The wrapper is borrowed; callers must
// clone it before modifying.
func (fw *fontWindow) getTextFont() *wx.Font {
	return fw.canvas.font
}

// getColour returns the canvas colour. The wrapper is owned by the canvas.
func (fw *fontWindow) getColour() *wx.Colour {
	return fw.canvas.colour
}

// makeNewFont is FontWindow::MakeNewFont().
func (fw *fontWindow) makeNewFont() *wx.Font {
	return fw.panel.getFontInfo()
}

// newFontPanel is FontPanel::FontPanel().
func newFontPanel(parent unsafe.Pointer) *fontPanel {
	win := wx.NewPanel(parent, wx.ID_ANY, -1, -1, -1, -1, wx.TAB_TRAVERSAL)
	if win == nil {
		return nil
	}

	p := &fontPanel{win: win}

	p.textFaceName = wx.NewTextCtrl(win.Ptr(), wx.ID_ANY, "", -1, -1, -1, -1, wx.TE_PROCESS_ENTER)
	if p.textFaceName == nil {
		return nil
	}
	p.textFaceName.Bind(wx.EVT_TEXT, func(evt *wx.Event) {
		if evt == nil {
			return
		}
		p.useFamily = false
	})

	p.choiceFamily = wx.NewChoice(win.Ptr(), wx.ID_ANY, -1, -1, -1, -1, nil, 0)
	if p.choiceFamily == nil {
		return nil
	}
	for _, name := range familyNames {
		p.choiceFamily.Append(name)
	}
	p.choiceFamily.Bind(wx.EVT_CHOICE, func(evt *wx.Event) {
		if evt == nil {
			return
		}
		p.useFamily = true
	})

	p.spinPointSize = wx.NewSpinCtrlDouble(win.Ptr(), wx.ID_ANY, "", -1, -1, -1, -1,
		wx.SP_ARROW_KEYS, 1.0, 100.0, 10.0, 0.1)
	if p.spinPointSize == nil {
		return nil
	}
	setSpinWidth(p.spinPointSize, win.Ptr(), "999.9")

	p.choiceStyle = wx.NewChoice(win.Ptr(), wx.ID_ANY, -1, -1, -1, -1, nil, 0)
	if p.choiceStyle == nil {
		return nil
	}
	for _, name := range styleNames {
		p.choiceStyle.Append(name)
	}

	p.choiceWeight = wx.NewChoice(win.Ptr(), wx.ID_ANY, -1, -1, -1, -1, nil, 0)
	if p.choiceWeight == nil {
		return nil
	}
	for _, name := range weightNames {
		p.choiceWeight.Append(name)
	}
	p.choiceWeight.Bind(wx.EVT_CHOICE, func(evt *wx.Event) {
		if evt == nil {
			return
		}
		p.useNumericWeight = false
	})

	p.spinWeight = wx.NewSpinCtrl(win.Ptr(), wx.ID_ANY, "", -1, -1, -1, -1,
		wx.SP_ARROW_KEYS, 1, wx.FONTWEIGHT_MAX, 1)
	if p.spinWeight == nil {
		return nil
	}
	setSpinWidth(p.spinWeight, win.Ptr(), "9999")
	p.spinWeight.Bind(wx.EVT_SPINCTRL, func(evt *wx.Event) {
		if evt == nil {
			return
		}
		p.useNumericWeight = true
	})

	p.checkUnderlined = wx.NewCheckBox(win.Ptr(), wx.ID_ANY, "", -1, -1, -1, -1, 0)
	p.checkStrikethrough = wx.NewCheckBox(win.Ptr(), wx.ID_ANY, "", -1, -1, -1, -1, 0)
	p.checkFixedWidth = wx.NewCheckBox(win.Ptr(), wx.ID_ANY, "", -1, -1, -1, -1, 0)
	if p.checkUnderlined == nil || p.checkStrikethrough == nil || p.checkFixedWidth == nil {
		return nil
	}
	p.checkFixedWidth.Disable()

	p.labelInfo = wx.NewStaticText(win.Ptr(), wx.ID_ANY, "\n\n\n", -1, -1, -1, -1, 0)
	if p.labelInfo == nil {
		return nil
	}

	border := wx.Sizer_GetDefaultBorder()

	// Columns are: label, control, gap, label, control, label, control.
	sizer := wx.NewFlexGridSizer(0, 7, border, border)

	addLabel := func(text string, cfg func(*wx.SizerFlags)) *wx.StaticText {
		label := wx.NewStaticText(win.Ptr(), wx.ID_ANY, text, -1, -1, -1, -1, 0)
		flags := wx.NewSizerFlags(0)
		if cfg == nil {
			flags.CenterVertical()
		} else {
			cfg(flags)
		}
		sizer.AddWindowFlags(label.Ptr(), flags.Ptr())
		flags.Delete()
		return label
	}
	addValue := func(w unsafe.Pointer) {
		flags := wx.NewSizerFlags(0)
		flags.Expand()
		flags.CenterVertical()
		sizer.AddWindowFlags(w, flags.Ptr())
		flags.Delete()
	}
	addGap := func(size int) {
		sizer.AddSpacer(size)
	}

	addLabel("Face &name:", nil)
	addValue(p.textFaceName.Ptr())
	addGap(2 * border)

	addLabel("&Family:", nil)
	addValue(p.choiceFamily.Ptr())

	addLabel("&Point size:", func(s *wx.SizerFlags) {
		s.DoubleBorder(wx.LEFT)
		s.CenterVertical()
	})
	addValue(p.spinPointSize.Ptr())

	addLabel("&Style:", nil)
	addValue(p.choiceStyle.Ptr())
	addGap(0)

	addLabel("&Weight:", nil)
	addValue(p.choiceWeight.Ptr())
	addLabel("or &raw value:", nil)
	addValue(p.spinWeight.Ptr())

	addLabel("&Underlined:", nil)
	addValue(p.checkUnderlined.Ptr())
	addGap(0)

	addLabel("&Strike through:", nil)
	addValue(p.checkStrikethrough.Ptr())
	addLabel("Fixed width:", nil)
	addValue(p.checkFixedWidth.Ptr())

	apply := wx.NewButton(win.Ptr(), wx.ID_APPLY, "&Apply changes", -1, -1, -1, -1, 0)
	if apply == nil {
		return nil
	}

	sizerTop := wx.NewBoxSizer(wx.VERTICAL)
	addSizer := func(s unsafe.Pointer, cfg func(*wx.SizerFlags)) {
		flags := wx.NewSizerFlags(0)
		if cfg != nil {
			cfg(flags)
		}
		sizerTop.AddSizerFlags(s, flags.Ptr())
		flags.Delete()
	}
	addControl := func(w unsafe.Pointer, cfg func(*wx.SizerFlags)) {
		flags := wx.NewSizerFlags(0)
		if cfg != nil {
			cfg(flags)
		}
		sizerTop.AddWindowFlags(w, flags.Ptr())
		flags.Delete()
	}

	addSizer(sizer.Ptr(), func(fl *wx.SizerFlags) {
		fl.Expand()
		fl.Border(wx.BOTTOM)
	})
	addControl(apply.Ptr(), func(fl *wx.SizerFlags) {
		fl.Border(wx.BOTTOM)
		fl.Center()
	})
	addControl(p.labelInfo.Ptr(), func(fl *wx.SizerFlags) {
		fl.Expand()
		fl.Border(wx.TOP)
	})

	win.SetSizer(sizerTop.Ptr(), false)

	p.showFont(wx.FontFromPtr(wx.NORMAL_FONT))
	return p
}

// setSpinWidth applies the C++ SetInitialSize(GetSizeFromTextSize(...)) idiom.
func setSpinWidth(ctrl interface {
	GetSizeFromTextSize(xlen int, ylen int, width unsafe.Pointer, height unsafe.Pointer)
	SetInitialSize(width int, height int) bool
}, winPtr unsafe.Pointer, sample string) {
	ext := windowTextExtent(winPtr, sample)
	var width, height int32
	ctrl.GetSizeFromTextSize(ext.Width, -1, unsafe.Pointer(&width), unsafe.Pointer(&height))
	ctrl.SetInitialSize(int(width), int(height))
}

// showFont is FontPanel::ShowFont().
func (p *fontPanel) showFont(font *wx.Font) {
	if font == nil {
		return
	}
	if p.font != nil {
		p.font.Delete()
	}
	p.font = cloneFont(font)
	p.doUpdate()
}

// doUpdate is FontPanel::DoUpdate().
func (p *fontPanel) doUpdate() {
	if p.font == nil {
		return
	}

	p.textFaceName.ChangeValue(p.font.GetFaceName())

	family := 0
	if fam := p.font.GetFamily(); fam == wx.FONTFAMILY_DECORATIVE {
		family = 1
	} else if fam == wx.FONTFAMILY_ROMAN {
		family = 2
	} else if fam == wx.FONTFAMILY_SCRIPT {
		family = 3
	} else if fam == wx.FONTFAMILY_SWISS {
		family = 4
	} else if fam == wx.FONTFAMILY_MODERN {
		family = 5
	} else if fam == wx.FONTFAMILY_TELETYPE {
		family = 6
	}
	p.choiceFamily.SetSelection(family)

	p.spinPointSize.SetValue(p.font.GetFractionalPointSize())

	style := 0
	if p.font.GetStyle() == wx.FONTSTYLE_ITALIC {
		style = 1
	} else if p.font.GetStyle() == wx.FONTSTYLE_SLANT {
		style = 2
	}
	p.choiceStyle.SetSelection(style)

	weight := 3
	if w := p.font.GetWeight(); w == wx.FONTWEIGHT_THIN {
		weight = 0
	} else if w == wx.FONTWEIGHT_EXTRALIGHT {
		weight = 1
	} else if w == wx.FONTWEIGHT_LIGHT {
		weight = 2
	} else if w == wx.FONTWEIGHT_MEDIUM {
		weight = 4
	} else if w == wx.FONTWEIGHT_SEMIBOLD {
		weight = 5
	} else if w == wx.FONTWEIGHT_BOLD {
		weight = 6
	} else if w == wx.FONTWEIGHT_EXTRABOLD {
		weight = 7
	} else if w == wx.FONTWEIGHT_HEAVY {
		weight = 8
	} else if w == wx.FONTWEIGHT_EXTRAHEAVY {
		weight = 9
	}
	p.choiceWeight.SetSelection(weight)
	p.spinWeight.SetValue(p.font.GetNumericWeight())

	if p.font.GetUnderlined() {
		p.checkUnderlined.SetValue(1)
	} else {
		p.checkUnderlined.SetValue(0)
	}
	if p.font.GetStrikethrough() {
		p.checkStrikethrough.SetValue(1)
	} else {
		p.checkStrikethrough.SetValue(0)
	}
	if p.font.IsFixedWidth() {
		p.checkFixedWidth.SetValue(1)
	} else {
		p.checkFixedWidth.SetValue(0)
	}

	pixelSize := p.font.GetPixelSize()

	dc := wx.NewClientDC(p.win.Ptr())
	dc.SetFont(p.font.Ptr())
	charWidth := dc.GetCharWidth()
	charHeight := dc.GetCharHeight()
	dc.Delete()

	p.labelInfo.SetLabel(fmt.Sprintf(
		"Font info string: %s\nSize in pixels: %d*%d, average char size: %d*%d",
		p.font.GetNativeFontInfoDesc(), pixelSize.Width, pixelSize.Height,
		charWidth, charHeight))
}

// getFontInfo is FontPanel::GetFontInfo().
func (p *fontPanel) getFontInfo() *wx.Font {
	info := wx.NewFontInfo(p.spinPointSize.GetValue())
	if info == nil {
		return nil
	}
	defer info.Delete()

	if p.useFamily {
		index := p.choiceFamily.GetSelection()
		if index < 0 || index >= len(familyNames) {
			index = 0
		}
		info.SetFamily(families()[index])
	} else {
		info.SetFaceName(p.textFaceName.GetValue())
	}

	switch p.choiceStyle.GetSelection() {
	case 1:
		info.SetItalic(true)
	case 2:
		info.SetSlant(true)
	}

	if p.useNumericWeight {
		info.SetWeight(p.spinWeight.GetValue())
	} else {
		info.SetWeight((p.choiceWeight.GetSelection() + 1) * 100)
	}

	if p.checkUnderlined.GetValue() {
		info.SetUnderlined(true)
	}
	if p.checkStrikethrough.GetValue() {
		info.SetStrikethrough(true)
	}

	return wx.NewFontFromInfo(info.Ptr())
}

// newFontCanvas is FontCanvas::FontCanvas().
func newFontCanvas(parent unsafe.Pointer) *fontCanvas {
	win := wx.NewWindow(parent, wx.ID_ANY, -1, -1, -1, -1, 0)
	if win == nil {
		return nil
	}

	c := &fontCanvas{win: win}
	c.font = stockFont(wx.NORMAL_FONT)
	c.colour = copyColour(wx.RED)
	if c.font == nil || c.colour == nil {
		return nil
	}

	// DoGetBestClientSize(): 80 average character widths by 15 character heights.
	win.SetMinSize(80*win.GetCharWidth(), 15*win.GetCharHeight())

	win.Bind(wx.EVT_PAINT, c.onPaint)
	return c
}

func (c *fontCanvas) setTextFont(font *wx.Font) {
	if font == nil {
		return
	}
	if c.font != nil {
		c.font.Delete()
	}
	c.font = cloneFont(font)
}

func (c *fontCanvas) setColour(colour *wx.Colour) {
	if colour == nil {
		return
	}
	if c.colour != nil {
		c.colour.Delete()
	}
	c.colour = copyColour(colour.Ptr())
}

// onPaint is FontCanvas::OnPaint().
func (c *fontCanvas) onPaint(evt *wx.Event) {
	if evt == nil {
		return
	}
	if c.font == nil || c.colour == nil {
		return
	}

	dc := wx.NewPaintDC(c.win.Ptr())
	if dc == nil {
		return
	}
	defer dc.Delete()

	c.win.PrepareDC(dc.Ptr())

	// Set the background.
	dc.SetBackground(wx.WHITE_BRUSH)
	dc.Clear()
	dc.SetFont(c.font.Ptr())

	// The current text origin.
	x, y := 5, 5

	dc.SetTextForeground(c.colour.Ptr())

	// The size of one cell (normally biggest char + small margin).
	maxExtent := dc.GetTextExtentSize("W")
	maxCharWidth, maxCharHeight := maxExtent.Width, maxExtent.Height
	w := maxCharWidth + 5
	h := maxCharHeight + 4

	// Print all font symbols from 32 to 255 in 7 rows of 32 chars each.
	for i := 0; i < 7; i++ {
		for j := 0; j < 32; j++ {
			ch := string(rune(32*(i+1) + j))
			extent := dc.GetTextExtentSize(ch)
			dc.DrawText(ch,
				x+w*j+(maxCharWidth-extent.Width)/2+1,
				y+h*i+(maxCharHeight-extent.Height)/2)
		}
	}

	// Draw the lines between them.
	dc.SetPen(wx.BLUE_PEN)

	// Horizontal.
	for l := 0; l < 8; l++ {
		yl := y + h*l - 2
		dc.DrawLine(x-2, yl, x+32*w-1, yl)
	}

	// And vertical.
	for l := 0; l < 33; l++ {
		xl := x + w*l - 2
		dc.DrawLine(xl, y-2, xl, y+7*h-1)
	}
}

// isChecked reads wxCommandEvent::IsChecked() from a menu event.
func isChecked(evt *wx.Event) bool {
	ce := wx.CommandEventFromPtr(evt.Ptr())
	if ce == nil {
		return false
	}
	return ce.IsChecked()
}

// doChangeFont is MyFrame::DoChangeFont().
func (f *fontFrame) doChangeFont(font *wx.Font, colour *wx.Colour) {
	if font == nil {
		return
	}

	f.fontWindow.updateFont(font, colour)

	f.textCtrl.SetFont(font.Ptr())
	if colour != nil {
		f.textCtrl.SetForegroundColour(colour.Ptr())
	}
	f.textCtrl.Refresh(false)

	// Update the state of the bold/italic/underlined menu items.
	if menuBar := wx.MenuBarFromPtr(f.frame.GetMenuBar()); menuBar != nil {
		menuBar.Check(idFontLight, font.GetWeight() == wx.FONTWEIGHT_LIGHT)
		menuBar.Check(idFontBold, font.GetWeight() == wx.FONTWEIGHT_BOLD)
		menuBar.Check(idFontItalic, font.GetStyle() == wx.FONTSTYLE_ITALIC)
		menuBar.Check(idFontUnderlined, font.GetUnderlined())
		menuBar.Check(idFontStrikethrough, font.GetStrikethrough())
	}
}

// doResizeFont is MyFrame::DoResizeFont().
func (f *fontFrame) doResizeFont(diff int) {
	font := cloneFont(f.fontWindow.getTextFont())
	if font == nil {
		return
	}
	font.SetPointSize(font.GetPointSize() + diff)
	f.doChangeFont(font, nil)
	font.Delete()
}

// onGetBaseFont is MyFrame::OnGetBaseFont().
func (f *fontFrame) onGetBaseFont(evt *wx.Event) {
	base := f.fontWindow.getTextFont().GetBaseFont()
	if base == nil {
		return
	}
	// GetBaseFont returns a new native font; copy it into an owned wrapper.
	font := wx.NewFontFromFont(base)
	if font == nil {
		return
	}
	f.doChangeFont(font, nil)
	font.Delete()
}

func (f *fontFrame) onBold(evt *wx.Event) {
	if evt == nil {
		return
	}
	font := cloneFont(f.fontWindow.getTextFont())
	if font == nil {
		return
	}
	if isChecked(evt) {
		font.SetWeight(wx.FONTWEIGHT_BOLD)
	} else {
		font.SetWeight(wx.FONTWEIGHT_NORMAL)
	}
	f.doChangeFont(font, nil)
	font.Delete()
}

func (f *fontFrame) onLight(evt *wx.Event) {
	if evt == nil {
		return
	}
	font := cloneFont(f.fontWindow.getTextFont())
	if font == nil {
		return
	}
	if isChecked(evt) {
		font.SetWeight(wx.FONTWEIGHT_LIGHT)
	} else {
		font.SetWeight(wx.FONTWEIGHT_NORMAL)
	}
	f.doChangeFont(font, nil)
	font.Delete()
}

func (f *fontFrame) onItalic(evt *wx.Event) {
	if evt == nil {
		return
	}
	font := cloneFont(f.fontWindow.getTextFont())
	if font == nil {
		return
	}
	if isChecked(evt) {
		font.SetStyle(wx.FONTSTYLE_ITALIC)
	} else {
		font.SetStyle(wx.FONTSTYLE_NORMAL)
	}
	f.doChangeFont(font, nil)
	font.Delete()
}

func (f *fontFrame) onUnderline(evt *wx.Event) {
	if evt == nil {
		return
	}
	font := cloneFont(f.fontWindow.getTextFont())
	if font == nil {
		return
	}
	if isChecked(evt) {
		font.SetUnderlined(1)
	} else {
		font.SetUnderlined(0)
	}
	f.doChangeFont(font, nil)
	font.Delete()
}

func (f *fontFrame) onStrikethrough(evt *wx.Event) {
	if evt == nil {
		return
	}
	font := cloneFont(f.fontWindow.getTextFont())
	if font == nil {
		return
	}
	font.SetStrikethrough(isChecked(evt))
	f.doChangeFont(font, nil)
	font.Delete()
}

// onStockFont is MyFrame::OnwxPointerFont().
func (f *fontFrame) onStockFont(evt *wx.Event) {
	if evt == nil {
		return
	}
	var stock unsafe.Pointer
	if id := evt.GetId(); id == idFontNormal {
		stock = wx.NORMAL_FONT
	} else if id == idFontSmall {
		stock = wx.SMALL_FONT
	} else if id == idFontItalicStock {
		stock = wx.ITALIC_FONT
	} else if id == idFontSwiss {
		stock = wx.SWISS_FONT
	} else {
		return
	}

	font := stockFont(stock)
	if font == nil {
		return
	}
	f.doChangeFont(font, nil)
	font.Delete()
}

// onFontDefault is MyFrame::OnFontDefault().
func (f *fontFrame) onFontDefault(evt *wx.Event) {
	info := wx.NewFontInfoDefault()
	if info == nil {
		return
	}
	font := wx.NewFontFromInfo(info.Ptr())
	info.Delete()
	if font == nil {
		return
	}
	f.doChangeFont(font, nil)
	font.Delete()
}

// onSystemSettingsFont is MyFrame::OnwxSystemSettingsFont().
func (f *fontFrame) onSystemSettingsFont(evt *wx.Event) {
	if evt == nil {
		return
	}

	var index int
	if id := evt.GetId(); id == idFontSysOEMFixed {
		index = wx.SYS_OEM_FIXED_FONT
	} else if id == idFontSysANSIFixed {
		index = wx.SYS_ANSI_FIXED_FONT
	} else if id == idFontSysANSIVar {
		index = wx.SYS_ANSI_VAR_FONT
	} else if id == idFontSysSystem {
		index = wx.SYS_SYSTEM_FONT
	} else if id == idFontSysDevice {
		index = wx.SYS_DEVICE_DEFAULT_FONT
	} else if id == idFontSysDefaultGUI {
		index = wx.SYS_DEFAULT_GUI_FONT
	} else {
		return
	}

	font := wx.NewFontDefault()
	if font == nil {
		return
	}
	wx.SystemSettings_GetFont(index, font.Ptr())

	f.doChangeFont(font, nil)
	font.Delete()
}

// onSetNativeDesc is MyFrame::OnSetNativeDesc().
func (f *fontFrame) onSetNativeDesc(evt *wx.Event) {
	fontInfo := wx.GetTextFromUserEx("Enter native font string", "Input font description",
		f.fontWindow.getTextFont().GetNativeFontInfoDesc(), f.frame.Ptr(), -1, -1, 1)
	if fontInfo == "" {
		return
	}

	font := wx.NewFontDefault()
	if font == nil {
		return
	}
	font.SetNativeFontInfo(fontInfo)
	if !font.IsOk() {
		wx.LogError(fmt.Sprintf("Font info string \"%s\" is invalid.", fontInfo))
		font.Delete()
		return
	}

	f.doChangeFont(font, nil)
	font.Delete()
}

// onSetNativeUserDesc is MyFrame::OnSetNativeUserDesc().
func (f *fontFrame) onSetNativeUserDesc(evt *wx.Event) {
	fontDesc := f.fontWindow.getTextFont().GetNativeFontInfoUserDesc()
	fontUserInfo := wx.GetTextFromUserEx("Here you can edit current font description",
		"Input font description", fontDesc, f.frame.Ptr(), -1, -1, 1)
	if fontUserInfo == "" {
		return
	}

	font := wx.NewFontDefault()
	if font == nil {
		return
	}
	if font.SetNativeFontInfoUserDesc(fontUserInfo) {
		f.doChangeFont(font, nil)
	} else {
		wx.MessageBox("Error trying to create a font with such description...",
			"", wx.OK, f.frame.Ptr(), -1, -1)
	}
	font.Delete()
}

// onSetFamily is MyFrame::OnSetFamily().
func (f *fontFrame) onSetFamily(evt *wx.Event) {
	family := f.getFamilyFromUser()

	font := cloneFont(f.fontWindow.getTextFont())
	if font == nil {
		return
	}
	font.SetFamily(family)
	f.doChangeFont(font, nil)
	font.Delete()
}

// onSetFaceName is MyFrame::OnSetFaceName().
func (f *fontFrame) onSetFaceName(evt *wx.Event) {
	facename := f.fontWindow.getTextFont().GetFaceName()
	newFaceName := wx.GetTextFromUserEx("Here you can edit current font face name.",
		"Input font facename", facename, f.frame.Ptr(), -1, -1, 1)
	if newFaceName == "" {
		return
	}

	font := cloneFont(f.fontWindow.getTextFont())
	if font == nil {
		return
	}
	font.SetFaceName(newFaceName)
	if font.IsOk() {
		f.doChangeFont(font, nil)
	} else {
		wx.MessageBox("There is no font with such face name...", "Invalid face name",
			wx.OK|wx.ICON_ERROR, f.frame.Ptr(), -1, -1)
	}
	font.Delete()
}

// onSetEncoding is MyFrame::OnSetEncoding().
func (f *fontFrame) onSetEncoding(evt *wx.Event) {
	enc := f.getEncodingFromUser()
	if enc == wx.FONTENCODING_SYSTEM {
		return
	}

	font := cloneFont(f.fontWindow.getTextFont())
	if font == nil {
		return
	}
	font.SetEncoding(enc)
	f.doChangeFont(font, nil)
	font.Delete()
}

// onPrivateFont is MyFrame::OnPrivateFont().
func (f *fontFrame) onPrivateFont(evt *wx.Event) {
	font := cloneFont(f.fontWindow.getTextFont())
	if font == nil {
		return
	}
	font.SetFaceName("wxprivate")
	if font.IsOk() {
		f.doChangeFont(font, nil)
	} else {
		wx.LogError("Failed to use private font.")
	}
	font.Delete()
}

// onQuit is MyFrame::OnQuit().
func (f *fontFrame) onQuit(evt *wx.Event) {
	f.frame.Close(true)
}

// onTestTextValue is MyFrame::OnTestTextValue().
func (f *fontFrame) onTestTextValue(evt *wx.Event) {
	value := f.textCtrl.GetValue()
	f.textCtrl.SetValue(value)
	if f.textCtrl.GetValue() != value {
		wx.LogError("Text value changed after getting and setting it")
	}
}

// getSingleChoiceIndex shows a wxSingleChoiceDialog and returns the chosen
// index, or -1 when the user cancelled.
func (f *fontFrame) getSingleChoiceIndex(message string, names []string) int {
	choices := wx.NewArrayString()
	if choices == nil {
		return -1
	}
	defer choices.Delete()
	for _, name := range names {
		choices.Add(name)
	}

	dlg := wx.NewSingleChoiceDialog(f.frame.Ptr(), message, sampleTitle,
		choices.Ptr(), wx.CHOICEDLG_STYLE)
	if dlg == nil {
		return -1
	}
	defer dlg.Delete()

	if dlg.ShowModal() != wx.ID_OK {
		return -1
	}
	return dlg.GetSelection()
}

// getEncodingFromUser is MyFrame::GetEncodingFromUser().
func (f *fontFrame) getEncodingFromUser() int {
	count := wx.FontMapper_GetSupportedEncodingsCount()
	names := make([]string, 0, count)
	encodings := make([]int, 0, count)

	for n := 0; n < count; n++ {
		enc := wx.FontMapper_GetEncoding(n)
		encodings = append(encodings, enc)
		names = append(names, wx.FontMapper_GetEncodingName(enc))
	}

	i := f.getSingleChoiceIndex("Choose the encoding", names)
	if i == -1 {
		return wx.FONTENCODING_SYSTEM
	}
	return encodings[i]
}

// getFamilyFromUser is MyFrame::GetFamilyFromUser().
func (f *fontFrame) getFamilyFromUser() int {
	familyValues := []int{
		wx.FONTFAMILY_DECORATIVE, wx.FONTFAMILY_ROMAN, wx.FONTFAMILY_SCRIPT,
		wx.FONTFAMILY_SWISS, wx.FONTFAMILY_MODERN, wx.FONTFAMILY_TELETYPE,
	}
	names := []string{"DECORATIVE", "ROMAN", "SCRIPT", "SWISS", "MODERN", "TELETYPE"}

	i := f.getSingleChoiceIndex("Choose the family", names)
	if i == -1 {
		return wx.FONTFAMILY_DEFAULT
	}
	return familyValues[i]
}

// doEnumerateFamilies is MyFrame::DoEnumerateFamilies().
func (f *fontFrame) doEnumerateFamilies(fixedWidthOnly bool, encoding int, silent bool) bool {
	enumerator := wx.NewFontEnumeratorDefault()
	if enumerator == nil {
		return false
	}
	defer enumerator.Delete()

	names := wx.NewArrayString()
	if names == nil {
		return false
	}
	defer names.Delete()

	count := enumerator.EnumerateFacenamesToArray(encoding, fixedWidthOnly, names.Ptr())

	if count > 0 {
		if !silent {
			prefix := ""
			if fixedWidthOnly {
				prefix = "fixed width "
			}
			wx.LogStatusWindow(f.frame.Ptr(), fmt.Sprintf("Found %d %sfonts", count, prefix))
		}

		facename := ""
		if silent {
			// Choose the first.
			facename = names.Item(0)
		} else {
			// Let the user choose, in alphabetical order.
			sorted := make([]string, count)
			for n := 0; n < count; n++ {
				sorted[n] = names.Item(n)
			}
			sort.Strings(sorted)

			if n := f.getSingleChoiceIndex("Choose a facename", sorted); n != -1 {
				facename = sorted[n]
			}
		}

		if facename != "" {
			info := wx.NewFontInfoDefault()
			if info != nil {
				info.SetFaceName(facename)
				info.SetEncoding(encoding)
				font := wx.NewFontFromInfo(info.Ptr())
				info.Delete()
				if font != nil {
					f.doChangeFont(font, nil)
					font.Delete()
				}
			}
		}

		return true
	} else if !silent {
		wx.LogWarning("No such fonts found.")
	}

	return false
}

// onEnumerateFamiliesForEncoding is MyFrame::OnEnumerateFamiliesForEncoding().
func (f *fontFrame) onEnumerateFamiliesForEncoding(evt *wx.Event) {
	enc := f.getEncodingFromUser()
	if enc != wx.FONTENCODING_SYSTEM {
		f.doEnumerateFamilies(false, enc, false)
	}
}

// onEnumerateEncodings is MyFrame::OnEnumerateEncodings(). The C++ version
// derives its text from wxFontEnumerator::OnFontEncoding(); the same pairs are
// reconstructed here from the callback-free ToArray helpers.
func (f *fontFrame) onEnumerateEncodings(evt *wx.Event) {
	enumerator := wx.NewFontEnumeratorDefault()
	if enumerator == nil {
		return
	}
	defer enumerator.Delete()

	facenames := wx.NewArrayString()
	if facenames == nil {
		return
	}
	defer facenames.Delete()

	faceCount := enumerator.EnumerateFacenamesToArray(wx.FONTENCODING_SYSTEM, false, facenames.Ptr())

	encodings := wx.NewArrayString()
	if encodings == nil {
		return
	}
	defer encodings.Delete()

	var text strings.Builder
	n := 0
	for i := 0; i < faceCount; i++ {
		facename := facenames.Item(i)
		count := enumerator.EnumerateEncodingsToArray(facename, encodings.Ptr())
		for j := 0; j < count; j++ {
			n++
			text.WriteString(fmt.Sprintf("Encoding %d: %s (available in facename '%s')\n",
				n, encodings.Item(j), facename))
		}
	}

	wx.LogMessage("Enumerating all available encodings:\n" + text.String())
}

// onSelectFont is MyFrame::OnSelectFont().
func (f *fontFrame) onSelectFont(evt *wx.Event) {
	data := wx.NewFontData()
	if data == nil {
		return
	}
	defer data.Delete()

	data.SetInitialFont(f.fontWindow.getTextFont().Ptr())
	data.SetColour(f.fontWindow.getColour().Ptr())

	dlg := wx.NewFontDialog(f.frame.Ptr(), data.Ptr())
	if dlg == nil {
		return
	}
	defer dlg.Delete()

	if dlg.ShowModal() == wx.ID_OK {
		retData := wx.NewFontData()
		if retData == nil {
			return
		}
		defer retData.Delete()

		dlg.GetFontData(retData.Ptr())

		font := wx.NewFontDefault()
		colour := wx.NewColourEmpty()
		if font != nil && colour != nil {
			retData.GetChosenFont(font.Ptr())
			retData.GetColour(colour.Ptr())
			f.doChangeFont(font, colour)
		}
		if font != nil {
			font.Delete()
		}
		if colour != nil {
			colour.Delete()
		}
	}
}

// onViewMsg is MyFrame::OnViewMsg(): it reads the charset header of an email
// message file, picks a matching font encoding and loads the file.
func (f *fontFrame) onViewMsg(evt *wx.Event) {
	dialog := wx.NewFileDialog(f.frame.Ptr(), "Open an email message file",
		f.lastDir, f.lastFile, "", -1, -1, wx.FD_DEFAULT_STYLE, -1, -1)
	if dialog == nil {
		return
	}
	if dialog.ShowModal() != wx.ID_OK {
		dialog.Delete()
		return
	}

	// Save for the next time.
	f.lastDir = dialog.GetDirectory()
	f.lastFile = dialog.GetFilename()
	filename := dialog.GetPath()
	dialog.Delete()

	file := wx.NewTextFileFromFile(filename)
	if file == nil {
		return
	}
	if !file.Open() {
		file.Delete()
		return
	}
	defer file.Delete()

	charset := ""
	const prefix = "Content-Type: text/plain; charset="

	count := file.GetLineCount()
	for n := 0; n < count; n++ {
		line := file.GetLine(n)

		if line == "" {
			// If it is an email message, headers are over, no need to parse
			// all the file.
			break
		}

		if strings.HasPrefix(line, prefix) {
			rest := line[len(prefix):]
			rest = strings.TrimPrefix(rest, "\"")
			if quote := strings.Index(rest, "\""); quote != -1 {
				rest = rest[:quote]
			}
			charset = rest
			break
		}
	}

	if charset == "" {
		wx.LogError(fmt.Sprintf("The file '%s' doesn't contain charset information.", filename))
		return
	}

	mapper := wx.FontMapperFromPtr(wx.FontMapper_Get())
	if mapper == nil {
		return
	}

	fontenc := mapper.CharsetToEncoding(charset, false)
	if fontenc == wx.FONTENCODING_SYSTEM {
		wx.LogError(fmt.Sprintf("Charset '%s' is unsupported.", charset))
		return
	}

	f.textCtrl.LoadFile(filename, 0)

	if fontenc == wx.FONTENCODING_UTF8 || !mapper.IsEncodingAvailable(fontenc, "") {
		// Try to find some similar encoding. wxEncodingConverter is not part of
		// the kwxGO surface, so the character conversion itself is skipped and
		// only the alternate encoding is selected.
		var alt int32
		if mapper.GetAltForEncoding(fontenc, unsafe.Pointer(&alt), "") {
			fontenc = int(alt)
		} else {
			wx.LogWarning(fmt.Sprintf("No fonts for encoding '%s' on this system.",
				wx.FontMapper_GetEncodingDescription(fontenc)))
		}
	}

	// And now create the correct font.
	if !f.doEnumerateFamilies(false, fontenc, true) {
		normalPointSize := wx.FontFromPtr(wx.NORMAL_FONT).GetPointSize()
		info := wx.NewFontInfo(float64(normalPointSize))
		if info == nil {
			return
		}
		info.SetEncoding(fontenc)
		font := wx.NewFontFromInfo(info.Ptr())
		info.Delete()

		if font != nil {
			if font.IsOk() {
				f.doChangeFont(font, nil)
			} else {
				wx.LogWarning(fmt.Sprintf("No fonts for encoding '%s' on this system.",
					wx.FontMapper_GetEncodingDescription(fontenc)))
			}
			font.Delete()
		}
	}
}

// onAbout is MyFrame::OnAbout().
func (f *fontFrame) onAbout(evt *wx.Event) {
	wx.MessageBox("wxWidgets font sample\n(c) 1999-2006 Vadim Zeitlin",
		"About "+sampleTitle, wx.OK|wx.ICON_INFORMATION, f.frame.Ptr(), -1, -1)
}

// onFontPanelApply is MyFrame::OnFontPanelApply().
func (f *fontFrame) onFontPanelApply(evt *wx.Event) {
	if evt == nil {
		return
	}
	font := f.fontWindow.makeNewFont()
	if font == nil {
		return
	}
	f.doChangeFont(font, nil)
	font.Delete()
}
