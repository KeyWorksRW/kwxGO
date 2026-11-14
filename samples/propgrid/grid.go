package main

import (
	"fmt"
	"unsafe"

	"github.com/KeyWorksRW/kwxGO/wx"
)

// newChoices builds a wxPGChoices handle from parallel label/value slices.
func newChoices(labels []string, values []int32) *wx.PGChoices {
	choices := wx.NewPGChoices()
	if choices == nil {
		return nil
	}
	for i, label := range labels {
		value := 0
		if i < len(values) {
			value = int(values[i])
		}
		choices.Add(label, value)
	}
	return choices
}

func newEnumProperty(label, name string, labels []string, values []int32, value int) *wx.EnumProperty {
	choices := newChoices(labels, values)
	if choices == nil {
		return nil
	}
	defer choices.Delete()
	return wx.NewEnumPropertyWithChoices(label, name, choices.Ptr(), value)
}

func newFlagsProperty(label, name string, labels []string, values []int32, value int) *wx.FlagsProperty {
	choices := newChoices(labels, values)
	if choices == nil {
		return nil
	}
	defer choices.Delete()
	return wx.NewFlagsPropertyWithChoices(label, name, choices.Ptr(), value)
}

func newEditEnumProperty(label, name string, labels []string, values []int32, value string) *wx.EditEnumProperty {
	choices := newChoices(labels, values)
	if choices == nil {
		return nil
	}
	defer choices.Delete()
	return wx.NewEditEnumPropertyWithChoices(label, name, choices.Ptr(), value)
}

func newArrayStringProp(label, name string, values []string) *wx.ArrayStringProperty {
	arr := wx.NewArrayString()
	if arr == nil {
		return nil
	}
	defer arr.Delete()
	for _, s := range values {
		arr.Add(s)
	}
	return wx.NewArrayStringPropertyFromArray(label, name, arr.Ptr())
}

func newMultiChoiceProp(label, name string, choicesLabels, initial []string) *wx.MultiChoiceProperty {
	choices := newChoices(choicesLabels, nil)
	if choices == nil {
		return nil
	}
	defer choices.Delete()
	arr := wx.NewArrayString()
	if arr == nil {
		return nil
	}
	defer arr.Delete()
	for _, s := range initial {
		arr.Add(s)
	}
	return wx.NewMultiChoicePropertyWithChoices(label, name, choices.Ptr(), arr.Ptr())
}

// newSystemColourProperty mirrors the C++ constructor overload that takes a
// wxColour (implicitly converted to wxColourPropertyValue by wxWidgets).
func newSystemColourProperty(label string, colour unsafe.Pointer) *wx.SystemColourProperty {
	cpv := wx.NewColourPropertyValueFromColour(colour)
	if cpv == nil {
		return nil
	}
	defer cpv.Delete()
	return wx.NewSystemColourProperty(label, wx.PG_LABEL, cpv.Ptr())
}

// newComposite builds a KwxPGProperty used in place of the C++ custom
// wxPGProperty subclasses, which cannot be derived from Go. The value is a
// display string; children are added by the caller.
func newComposite(label, name, display string) *wx.KwxPGProperty {
	v := wx.NewVariantFromString(display, "")
	defer v.Delete()
	return wx.NewKwxPGProperty(nil, label, name, v.Ptr())
}

func newVectorProperty(label, name string) *wx.KwxPGProperty {
	prop := newComposite(label, name, "0; 0; 0")
	if prop == nil {
		return nil
	}
	prop.AddPrivateChild(wx.NewFloatProperty("X", wx.PG_LABEL, 0).Ptr())
	prop.AddPrivateChild(wx.NewFloatProperty("Y", wx.PG_LABEL, 0).Ptr())
	prop.AddPrivateChild(wx.NewFloatProperty("Z", wx.PG_LABEL, 0).Ptr())
	return prop
}

func newTriangleProperty(label, name string) *wx.KwxPGProperty {
	prop := newComposite(label, name, "")
	if prop == nil {
		return nil
	}
	prop.AddPrivateChild(newVectorProperty("A", wx.PG_LABEL).Ptr())
	prop.AddPrivateChild(newVectorProperty("B", wx.PG_LABEL).Ptr())
	prop.AddPrivateChild(newVectorProperty("C", wx.PG_LABEL).Ptr())
	return prop
}

func pgIface(p unsafe.Pointer) *wx.PropertyGridInterface {
	return wx.PropertyGridInterfaceFromPtr(p)
}

func propByName(i *wx.PropertyGridInterface, name string) unsafe.Pointer {
	if i == nil {
		return nil
	}
	return i.GetPropertyByName(name)
}

func pageByName(m *wx.PropertyGridManager, name string) unsafe.Pointer {
	if m == nil {
		return nil
	}
	idx := m.GetPageByName(name)
	if idx < 0 {
		return nil
	}
	return m.GetPage(uint(idx))
}

func setAttrLong(i *wx.PropertyGridInterface, id unsafe.Pointer, attr string, value, flags int) {
	if i == nil || id == nil {
		return
	}
	v := wx.NewVariantFromLong(value, "")
	defer v.Delete()
	i.SetPropertyAttribute(id, attr, v.Ptr(), flags)
}

func setAttrBool(i *wx.PropertyGridInterface, id unsafe.Pointer, attr string, value bool, flags int) {
	if i == nil || id == nil {
		return
	}
	v := wx.NewVariantFromBool(value, "")
	defer v.Delete()
	i.SetPropertyAttribute(id, attr, v.Ptr(), flags)
}

func setAttrDouble(i *wx.PropertyGridInterface, id unsafe.Pointer, attr string, value float64, flags int) {
	if i == nil || id == nil {
		return
	}
	v := wx.NewVariantFromDouble(value, "")
	defer v.Delete()
	i.SetPropertyAttribute(id, attr, v.Ptr(), flags)
}

func setAttrString(i *wx.PropertyGridInterface, id unsafe.Pointer, attr, value string, flags int) {
	if i == nil || id == nil {
		return
	}
	v := wx.NewVariantFromString(value, "")
	defer v.Delete()
	i.SetPropertyAttribute(id, attr, v.Ptr(), flags)
}

func setLongByName(i *wx.PropertyGridInterface, name, attr string, value, flags int) {
	setAttrLong(i, propByName(i, name), attr, value, flags)
}

func setStringByName(i *wx.PropertyGridInterface, name, attr, value string, flags int) {
	setAttrString(i, propByName(i, name), attr, value, flags)
}

func setBoolByName(i *wx.PropertyGridInterface, name, attr string, value bool, flags int) {
	setAttrBool(i, propByName(i, name), attr, value, flags)
}

func helpString(i *wx.PropertyGridInterface, name, text string) {
	if p := propByName(i, name); p != nil {
		i.SetPropertyHelpString(p, text)
	}
}

// eachProperty visits every descendant of root (excluding root itself) in
// display order.
func eachProperty(root unsafe.Pointer, fn func(p unsafe.Pointer)) {
	if root == nil {
		return
	}
	var walk func(p unsafe.Pointer)
	walk = func(p unsafe.Pointer) {
		prop := wx.PGPropertyFromPtr(p)
		if prop == nil {
			return
		}
		fn(p)
		n := prop.GetChildCount()
		for j := 0; j < n; j++ {
			walk(prop.Item(uint(j)))
		}
	}
	rp := wx.PGPropertyFromPtr(root)
	n := rp.GetChildCount()
	for i := 0; i < n; i++ {
		walk(rp.Item(uint(i)))
	}
}

// createGrid is FormMain::CreateGrid().
func (f *mainFrame) createGrid(style, extraStyle int) {
	if style == -1 {
		style = wx.PG_BOLD_MODIFIED | wx.PG_SPLITTER_AUTO_CENTER | wx.PG_AUTO_SORT |
			wx.PG_TOOLBAR | wx.PG_DESCRIPTION
	}
	if extraStyle == -1 {
		extraStyle = wx.PG_EX_MODE_BUTTONS | wx.PG_EX_NATIVE_DOUBLE_BUFFERING |
			wx.PG_EX_MULTIPLE_SELECTION
	}

	manager := wx.NewPropertyGridManager(f.panel.Ptr(), idPGID, -1, -1, -1, -1, style)
	if manager == nil {
		f.manager = nil
		return
	}
	f.manager = manager
	f.grid = wx.PropertyGridFromPtr(manager.GetGrid())

	manager.SetExtraStyle(extraStyle)

	mi := pgIface(manager.Ptr())
	// This is the default validation failure behaviour.
	mi.SetValidationFailureBehavior(wx.PG_VFB_MARK_CELL | wx.PG_VFB_SHOW_MESSAGEBOX)

	f.grid.SetVerticalSpacing(2)

	// Set somewhat different unspecified value appearance.
	cell := wx.NewPGCell()
	if cell != nil {
		cell.SetText("Unspecified")
		cell.SetFgCol(wx.LIGHT_GREY)
		f.grid.SetUnspecifiedValueAppearance(cell.Ptr())
		cell.Delete()
	}

	f.populateGrid()

	f.grid.MakeColumnEditable(0, f.labelEditingEnabled)
	manager.ShowHeader(f.hasHeader)
	if f.hasHeader {
		manager.SetColumnTitle(2, "Units")
	}
}

// replaceGrid is FormMain::ReplaceGrid().
func (f *mainFrame) replaceGrid(style, extraStyle int) {
	old := f.manager
	f.createGrid(style, extraStyle)
	if old != nil && f.manager != nil && f.topSizer != nil {
		f.topSizer.ReplaceWindow(old.Ptr(), f.manager.Ptr(), false)
		old.Destroy()
	}
	if f.manager != nil {
		f.manager.SetFocus()
	}
	f.panel.Layout()
}

// populateGrid is FormMain::PopulateGrid().
func (f *mainFrame) populateGrid() {
	f.manager.AddPage("Standard Items", nil)
	f.populateWithStandardItems()

	f.manager.AddPage("wxWidgets Library Config", nil)
	f.populateWithLibraryConfig()

	// wxMyPropertyGridPage is a wxPropertyGridPage with event logging and
	// IsHandlingAllEvents() == false. Deriving from wxPropertyGridPage from Go
	// needs C callback pointers, so a plain page is used instead.
	myPage := wx.NewPropertyGridPage()
	if myPage != nil {
		pi := pgIface(myPage.Ptr())
		pi.Append(wx.NewIntProperty("IntProperty", wx.PG_LABEL, 12345678).Ptr())
		f.manager.AddPageObject("Examples", nil, myPage.Ptr())
	}

	f.populateWithExamples()
}

// populateWithStandardItems is FormMain::PopulateWithStandardItems().
func (f *mainFrame) populateWithStandardItems() {
	page := pageByName(f.manager, "Standard Items")
	if page == nil {
		return
	}
	pg := pgIface(page)

	pg.Append(wx.NewPropertyCategory("Appearance").Ptr())

	pg.Append(wx.NewStringProperty("Label", wx.PG_LABEL, f.frame.GetTitle()).Ptr())
	font := wx.NewFontDefault()
	pg.Append(wx.NewFontProperty("Font", wx.PG_LABEL, font.Ptr()).Ptr())
	font.Delete()
	helpString(pg, "Font", "Editing this will change font used in the property grid.")

	pg.Append(newSystemColourProperty("Margin Colour", f.grid.GetMarginColour()).Ptr())
	pg.Append(newSystemColourProperty("Cell Colour", f.grid.GetCellBackgroundColour()).Ptr())
	pg.Append(newSystemColourProperty("Cell Text Colour", f.grid.GetCellTextColour()).Ptr())
	pg.Append(newSystemColourProperty("Line Colour", f.grid.GetLineColour()).Ptr())

	winStyleLabels := []string{
		"wxSIMPLE_BORDER", "wxDOUBLE_BORDER", "wxSUNKEN_BORDER", "wxRAISED_BORDER",
		"wxNO_BORDER", "wxTAB_TRAVERSAL", "wxWANTS_CHARS", "wxVSCROLL",
		"wxALWAYS_SHOW_SB", "wxCLIP_CHILDREN", "wxFULL_REPAINT_ON_RESIZE",
	}
	winStyleValues := []int32{
		int32(wx.BORDER_SIMPLE), int32(wx.DOUBLE_BORDER), int32(wx.BORDER_SUNKEN),
		int32(wx.BORDER_RAISED), int32(wx.BORDER_NONE), int32(wx.TAB_TRAVERSAL),
		int32(wx.WANTS_CHARS), int32(wx.VSCROLL), int32(wx.ALWAYS_SHOW_SB),
		int32(wx.CLIP_CHILDREN), int32(wx.FULL_REPAINT_ON_RESIZE),
	}
	frameStyleLabels := []string{
		"wxCAPTION", "wxMINIMIZE", "wxMAXIMIZE", "wxCLOSE_BOX", "wxSTAY_ON_TOP",
		"wxSYSTEM_MENU", "wxRESIZE_BORDER", "wxFRAME_TOOL_WINDOW",
		"wxFRAME_NO_TASKBAR", "wxFRAME_FLOAT_ON_PARENT", "wxFRAME_SHAPED",
	}
	frameStyleValues := []int32{
		int32(wx.CAPTION), int32(wx.ICONIZE), int32(wx.MAXIMIZE), int32(wx.CLOSE_BOX),
		int32(wx.STAY_ON_TOP), int32(wx.SYSTEM_MENU), int32(wx.RESIZE_BORDER),
		int32(wx.FRAME_TOOL_WINDOW), int32(wx.FRAME_NO_TASKBAR),
		int32(wx.FRAME_FLOAT_ON_PARENT), int32(wx.FRAME_SHAPED),
	}
	labels := append(append([]string{}, winStyleLabels...), frameStyleLabels...)
	values := append(append([]int32{}, winStyleValues...), frameStyleValues...)
	pg.Append(newFlagsProperty("Window Styles", wx.PG_LABEL, labels, values,
		f.frame.GetWindowStyle()).Ptr())

	pg.Append(wx.NewCursorProperty("Cursor", wx.PG_LABEL, wx.CURSOR_ARROW).Ptr())

	positionCat := wx.NewPropertyCategory("Position")
	positionCat.SetName("PositionCategory")
	pg.Append(positionCat.Ptr())
	helpString(pg, "PositionCategory",
		"Change in items in this category will cause respective changes in frame.")

	height := wx.NewIntProperty("Height", wx.PG_LABEL, 480)
	pg.Append(height.Ptr())
	setAttrLong(pg, height.Ptr(), wx.PG_ATTR_MIN, 10, 0)
	setAttrLong(pg, height.Ptr(), wx.PG_ATTR_MAX, 2048, 0)
	setAttrString(pg, height.Ptr(), wx.PG_ATTR_UNITS, "Pixels", 0)
	pg.SetPropertyValueUnspecified(height.Ptr())
	setAttrString(pg, height.Ptr(), wx.PG_ATTR_HINT, "Enter new height for window", 0)
	helpString(pg, "Height", "This property uses attributes \"Units\" and \"Hint\".")

	width := wx.NewIntProperty("Width", wx.PG_LABEL, 640)
	pg.Append(width.Ptr())
	setAttrLong(pg, width.Ptr(), wx.PG_ATTR_MIN, 10, 0)
	setAttrLong(pg, width.Ptr(), wx.PG_ATTR_MAX, 2048, 0)
	setAttrString(pg, width.Ptr(), wx.PG_ATTR_UNITS, "Pixels", 0)
	pg.SetPropertyValueUnspecified(width.Ptr())
	setAttrString(pg, width.Ptr(), wx.PG_ATTR_HINT, "Enter new width for window", 0)
	helpString(pg, "Width", "This property uses attributes \"Units\" and \"Hint\".")

	x := wx.NewIntProperty("X", wx.PG_LABEL, 10)
	pg.Append(x.Ptr())
	setAttrString(pg, x.Ptr(), wx.PG_ATTR_UNITS, "Pixels", 0)
	helpString(pg, "X", "This property uses \"Units\" attribute.")

	y := wx.NewIntProperty("Y", wx.PG_LABEL, 10)
	pg.Append(y.Ptr())
	setAttrString(pg, y.Ptr(), wx.PG_ATTR_UNITS, "Pixels", 0)
	helpString(pg, "Y", "This property uses \"Units\" attribute.")

	disabledHelp := "This property is simply disabled. In order to have label disabled as well, " +
		"you need to set wxPG_EX_GREY_LABEL_WHEN_DISABLED using SetExtraStyle."

	pg.Append(wx.NewPropertyCategory("Environment").Ptr())
	pg.Append(wx.NewStringProperty("Operating System", wx.PG_LABEL, wx.GetOsDescription()).Ptr())
	pg.Append(wx.NewStringProperty("User Id", wx.PG_LABEL, wx.GetUserId()).Ptr())
	pg.Append(wx.NewDirProperty("User Home", wx.PG_LABEL, wx.GetUserHome("")).Ptr())
	pg.Append(wx.NewStringProperty("User Name", wx.PG_LABEL, wx.GetUserName()).Ptr())

	for _, name := range []string{"Operating System", "User Id", "User Name"} {
		if p := propByName(pg, name); p != nil {
			pg.DisableProperty(p)
		}
		helpString(pg, name, disabledHelp)
	}

	pg.Append(wx.NewPropertyCategory("More Examples").Ptr())

	// wxFontDataProperty (a font with a Colour child) is a custom class; a
	// plain wxFontProperty is used instead.
	fdFont := wx.NewFontDefault()
	pg.Append(wx.NewFontProperty("FontDataProperty", wx.PG_LABEL, fdFont.Ptr()).Ptr())
	fdFont.Delete()
	helpString(pg, "FontDataProperty",
		"This demonstrates wxFontDataProperty class defined in this sample app. "+
			"It is exactly like wxFontProperty from the library, but also has colour sub-property.")

	// wxDirsProperty is an array-of-strings property with a custom browser
	// action; a stock wxArrayStringProperty is used instead.
	pg.Append(newArrayStringProp("DirsProperty", wx.PG_LABEL, nil).Ptr())
	helpString(pg, "DirsProperty",
		"This demonstrates wxDirsProperty class defined in this sample app. "+
			"It is built with WX_PG_IMPLEMENT_ARRAYSTRING_PROPERTY_WITH_VALIDATOR macro, "+
			"with custom action (dir dialog popup) defined.")

	// wxArrayDoubleProperty is a custom list-editor property; an array string
	// property with the same initial values stands in.
	pg.Append(newArrayStringProp("ArrayDoubleProperty", wx.PG_LABEL,
		[]string{"-1", "-0.5", "0", "0.5", "1"}).Ptr())
	helpString(pg, "ArrayDoubleProperty",
		"This demonstrates wxArrayDoubleProperty class defined in this sample app. "+
			"It is an example of a custom list editor property.")

	pg.Append(wx.NewLongStringProperty("Information", wx.PG_LABEL,
		"Editing properties will have immediate effect on this window, "+
			"and vice versa (at least in most cases, that is).").Ptr())
	helpString(pg, "Information", "This property is read-only.")
	if p := propByName(pg, "Information"); p != nil {
		pg.SetPropertyReadOnly(p, true, pgDontRecurse)
	}

	// Set test information for cells in columns 3 and 4 (column 2 is reserved
	// for displaying units).
	var folder *wx.BitmapBundle
	if b := wx.ArtProvider_GetBitmapBundle(wx.ART_FOLDER, wx.ART_OTHER, 0, 0); b != nil {
		folder = b
		defer folder.Delete()
	}
	pageRoot := wx.PropertyGridPageFromPtr(page).GetRoot()
	eachProperty(pageRoot, func(p unsafe.Pointer) {
		prop := wx.PGPropertyFromPtr(p)
		if prop == nil || prop.IsCategory() {
			return
		}
		var bmp unsafe.Pointer
		if folder != nil {
			bmp = folder.Ptr()
		}
		pg.SetPropertyCell(p, 3, "Cell 3", bmp, nil, nil)
		pg.SetPropertyCell(p, 4, "Cell 4", nil, wx.WHITE, wx.BLACK)
	})
}

// populateWithLibraryConfig is FormMain::PopulateWithLibraryConfig().
func (f *mainFrame) populateWithLibraryConfig() {
	page := pageByName(f.manager, "wxWidgets Library Config")
	if page == nil {
		return
	}
	pg := pgIface(page)

	if f.manager.HasFlag(wx.PG_SPLITTER_AUTO_CENTER) {
		pg.SetColumnProportion(0, 3)
		pg.SetColumnProportion(1, 1)
	}

	var bmp *wx.BitmapBundle
	if b := wx.ArtProvider_GetBitmapBundle(wx.ART_REPORT_VIEW, wx.ART_OTHER, 0, 0); b != nil {
		bmp = b
		defer bmp.Delete()
	}

	var bmpPtr unsafe.Pointer
	if bmp != nil {
		bmpPtr = bmp.Ptr()
	}

	var pid unsafe.Pointer

	var captionFont *wx.Font
	if ptr := f.grid.GetCaptionFont(); ptr != nil {
		captionFont = wx.FontFromPtr(ptr)
		captionFont.SetStyle(wx.FONTSTYLE_ITALIC)
	}

	italicFontHelp := "Font of this property's wxPGCell has " +
		"been modified. Obtain property's cell " +
		"with wxPGProperty::" +
		"GetOrCreateCell(column)."

	addGroup := func(label string) unsafe.Pointer {
		cat := wx.NewPropertyCategory(label)
		// AppendIn() returns the category that actually ended up in the page;
		// when a category with this name already exists wxWidgets deletes the
		// one passed in and reuses the existing one.
		catPtr := pg.AppendIn(pid, cat.Ptr())
		if catPtr == nil {
			return nil
		}
		pg.SetPropertyCell(catPtr, 0, wx.PG_LABEL, bmpPtr, nil, nil)
		if captionFont != nil {
			if cell := wx.PGCellFromPtr(wx.PGPropertyFromPtr(catPtr).GetCell(0)); cell != nil {
				cell.SetFont(captionFont.Ptr())
			}
		}
		wx.PGPropertyFromPtr(catPtr).SetHelpString(italicFontHelp)
		return catPtr
	}

	addBool := func(name string, value bool) unsafe.Pointer {
		p := wx.NewBoolProperty(name, wx.PG_LABEL, value)
		pg.Append(p.Ptr())
		return p.Ptr()
	}

	addNoDef := func(name string) {
		p := wx.NewBoolProperty(name, wx.PG_LABEL, false)
		pg.Append(p.Ptr())
		if ptr := propByName(pg, name); ptr != nil {
			pg.DisableProperty(ptr)
		}
	}

	pidProp := wx.NewPropertyCategory("wxWidgets Library Configuration")
	pid = pg.Append(pidProp.Ptr())
	if pid == nil {
		return
	}
	pg.SetPropertyCell(pid, 0, wx.PG_LABEL, bmpPtr, nil, nil)
	pg.SetPropertyCell(pid, 1, "Is Enabled", nil, nil, nil)
	pg.SetPropertyValueString(pid, "Is Enabled")

	addGroup("Global Settings")
	addBool("wxUSE_GUI", true)

	addGroup("Compatibility Settings")
	addBool("WXWIN_COMPATIBILITY_3_0", false)
	addBool("WXWIN_COMPATIBILITY_3_2", true)
	addNoDef("wxFONT_SIZE_COMPATIBILITY")
	addBool("wxDIALOG_UNIT_COMPATIBILITY", false)

	addGroup("Debugging Settings")
	addBool("wxUSE_DEBUGREPORT", true)
	addBool("wxUSE_STACKWALKER", true)
	addBool("wxUSE_ON_FATAL_EXCEPTION", true)

	addGroup("Unicode Support")
	addBool("wxUSE_UNICODE", true)

	addGroup("Global Features")
	addBool("wxUSE_EXCEPTIONS", true)
	addBool("wxUSE_EXTENDED_RTTI", false)
	addBool("wxUSE_STD_CONTAINERS", true)
	addBool("wxUSE_LOG", true)
	addBool("wxUSE_LOGWINDOW", true)
	addBool("wxUSE_LOGGUI", true)
	addBool("wxUSE_LOG_DIALOG", true)
	addBool("wxUSE_CMDLINE_PARSER", true)
	addBool("wxUSE_THREADS", true)
	addBool("wxUSE_STREAMS", true)
	addBool("wxUSE_STD_IOSTREAM", true)

	addGroup("Non-GUI Features")
	for _, name := range []string{
		"wxUSE_FILE", "wxUSE_FFILE", "wxUSE_FSVOLUME", "wxUSE_TEXTBUFFER",
		"wxUSE_TEXTFILE", "wxUSE_INTL", "wxUSE_DATETIME", "wxUSE_TIMER",
		"wxUSE_STOPWATCH", "wxUSE_CONFIG", "wxUSE_CONFIG_NATIVE",
		"wxUSE_DIALUP_MANAGER", "wxUSE_DYNLIB_CLASS", "wxUSE_DYNAMIC_LOADER",
		"wxUSE_SOCKETS", "wxUSE_FILESYSTEM", "wxUSE_FS_ZIP", "wxUSE_FS_INET",
		"wxUSE_ZIPSTREAM", "wxUSE_ZLIB", "wxUSE_APPLE_IEEE", "wxUSE_JOYSTICK",
		"wxUSE_FONTMAP", "wxUSE_MIMETYPE", "wxUSE_PROTOCOL",
		"wxUSE_PROTOCOL_FILE", "wxUSE_PROTOCOL_FTP", "wxUSE_PROTOCOL_HTTP",
		"wxUSE_URL",
	} {
		addBool(name, true)
	}
	addBool("wxUSE_URL_NATIVE", false)
	for _, name := range []string{"wxUSE_REGEX", "wxUSE_SYSTEM_OPTIONS", "wxUSE_SOUND", "wxUSE_XRC", "wxUSE_XML"} {
		addBool(name, true)
	}

	// Set them to use check box.
	setAttrBool(pg, pid, wx.PG_BOOL_USE_CHECKBOX, true, pgRecurse)
}

// populateWithExamples is FormMain::PopulateWithExamples().
func (f *mainFrame) populateWithExamples() {
	page := pageByName(f.manager, "Examples")
	if page == nil {
		return
	}
	pg := pgIface(page)

	// SpinCtrl editor.
	spin := wx.NewIntProperty("SpinCtrl", wx.PG_LABEL, 0)
	pg.Append(spin.Ptr())
	if p := propByName(pg, "SpinCtrl"); p != nil {
		pg.SetPropertyEditorByName(p, "SpinCtrl")
		setAttrLong(pg, p, wx.PG_ATTR_MIN, -2, 0)
		setAttrLong(pg, p, wx.PG_ATTR_MAX, 16384, 0)
		setAttrLong(pg, p, wx.PG_ATTR_SPINCTRL_STEP, 2, 0)
		setAttrBool(pg, p, wx.PG_ATTR_SPINCTRL_MOTION, true, 0)
	}
	helpString(pg, "SpinCtrl",
		"This is regular wxIntProperty, which editor has been "+
			"changed to wxPGEditor_SpinCtrl. Note however that "+
			"static wxPropertyGrid::RegisterAdditionalEditors() "+
			"needs to be called prior to using it.")

	pg.Append(wx.NewBoolProperty("BoolProperty", wx.PG_LABEL, false).Ptr())

	pg.Append(wx.NewBoolProperty("BoolProperty with CheckBox", wx.PG_LABEL, false).Ptr())
	setBoolByName(pg, "BoolProperty with CheckBox", wx.PG_BOOL_USE_CHECKBOX, true, 0)
	helpString(pg, "BoolProperty with CheckBox",
		"Property attribute wxPG_BOOL_USE_CHECKBOX has been set to true.")

	floatProp := wx.NewFloatProperty("FloatProperty", wx.PG_LABEL, 1234500.23)
	pg.Append(floatProp.Ptr())
	setAttrDouble(pg, floatProp.Ptr(), wx.PG_ATTR_MIN, -100.12, 0)

	pg.Append(wx.NewLongStringProperty("LongStringProperty", "LongStringProp",
		"This is much longer string than the first one. Edit it by clicking the button.").Ptr())

	pg.Append(newArrayStringProp("ArrayStringProperty", wx.PG_LABEL,
		[]string{"String 1", "String 2", "String 3"}).Ptr())

	fileProp := wx.NewFileProperty("FileProperty", "TextFile", "")
	pg.Append(fileProp.Ptr())
	setAttrString(pg, fileProp.Ptr(), wx.PG_FILE_WILDCARD, "Text Files (*.txt)|*.txt", 0)
	setAttrString(pg, fileProp.Ptr(), wx.PG_DIALOG_TITLE, "Custom File Dialog Title", 0)
	setAttrBool(pg, fileProp.Ptr(), wx.PG_FILE_SHOW_FULL_PATH, false, 0)
	setAttrString(pg, fileProp.Ptr(), wx.PG_FILE_SHOW_RELATIVE_PATH, "C:\\Windows", 0)
	pg.SetPropertyValueString(fileProp.Ptr(), "C:\\Windows\\System32\\msvcrt71.dll")

	pg.Append(wx.NewImageFileProperty("ImageFile", wx.PG_LABEL, "").Ptr())

	pg.Append(wx.NewColourProperty("ColourProperty", wx.PG_LABEL, wx.RED).Ptr())
	if p := propByName(pg, "ColourProperty"); p != nil {
		pg.SetPropertyEditorByName(p, "ComboBox")
		wx.PGPropertyFromPtr(p).SetAutoUnspecified(true)
	}
	helpString(pg, "ColourProperty",
		"wxPropertyGrid::SetPropertyEditor method has been used to change "+
			"editor of this property to wxPGEditor_ComboBox)")

	alphaColour := wx.NewColourRGB(15, 200, 95, 128)
	pg.Append(wx.NewColourProperty("ColourPropertyWithAlpha", wx.PG_LABEL, alphaColour.Ptr()).Ptr())
	alphaColour.Delete()
	setBoolByName(pg, "ColourPropertyWithAlpha", wx.PG_COLOUR_HAS_ALPHA, true, 0)
	helpString(pg, "ColourPropertyWithAlpha",
		"Attribute \"HasAlpha\" is set to true for this property.")

	pg.Append(wx.NewColourProperty("ColourProperty2", wx.PG_LABEL, wx.GREEN).Ptr())

	// EnumProperty uses the label/value pairs from the sample.
	enumLabels := []string{"One Item", "Another Item", "One More", "This Is Last"}
	enumValues := []int32{40, 80, 120, 160}
	pg.Append(newEnumProperty("EnumProperty", wx.PG_LABEL, enumLabels, enumValues, 80).Ptr())

	// The C++ sample builds wxPGChoices with custom colours and a bitmap in
	// the entries.
	socLabels := []string{
		"One Item", "Another Item", "One More", "This Is Last",
		"Look, it continues", "Even More", "And More", "", "True End of the List",
	}
	socValues := []int32{40, 80, 120, 160, 200, 240, 280, 300, 320}

	soc := newChoices(socLabels, socValues)
	if soc == nil {
		return
	}
	defer soc.Delete()

	// Test custom colours ([] operator of wxPGChoices returns references to
	// wxPGChoiceEntry, so the entries are borrowed).
	if e := wx.PGChoiceEntryFromPtr(soc.Item(1)); e != nil {
		e.SetFgCol(wx.RED)
		e.SetBgCol(wx.LIGHT_GREY)
	}
	if e := wx.PGChoiceEntryFromPtr(soc.Item(2)); e != nil {
		e.SetFgCol(wx.GREEN)
		e.SetBgCol(wx.LIGHT_GREY)
	}
	if e := wx.PGChoiceEntryFromPtr(soc.Item(3)); e != nil {
		e.SetFgCol(wx.BLUE)
		e.SetBgCol(wx.LIGHT_GREY)
	}
	if e := wx.PGChoiceEntryFromPtr(soc.Item(4)); e != nil {
		// wxPGCell::SetBitmap takes a wxBitmapBundle in wx 3.3; the C++
		// sample relies on the implicit wxBitmap -> wxBitmapBundle
		// conversion, so ask the art provider for the bundle directly.
		if b := wx.ArtProvider_GetBitmapBundle(wx.ART_FOLDER, wx.ART_OTHER, 0, 0); b != nil {
			e.SetBitmap(b.Ptr())
			b.Delete()
		}
	}

	if enum2 := wx.NewEnumPropertyWithChoices("EnumProperty 2", wx.PG_LABEL, soc.Ptr(), 240); enum2 != nil {
		pg.Append(enum2.Ptr())
		if p := propByName(pg, "EnumProperty 2"); p != nil {
			wx.PGPropertyFromPtr(p).AddChoice("Testing Extra", 360)
		}
	}

	if p := wx.NewEnumPropertyWithChoices("EnumProperty 3", wx.PG_LABEL, soc.Ptr(), 240); p != nil {
		pg.Append(p.Ptr())
		if id := propByName(pg, "EnumProperty 3"); id != nil {
			v := wx.NewVariantFromString("Dummy Hint", "")
			wx.PGPropertyFromPtr(id).SetAttribute(wx.PG_ATTR_HINT, v.Ptr())
			v.Delete()
		}
	}
	helpString(pg, "EnumProperty 3", "This property uses \"Hint\" attribute.")

	if p := wx.NewEnumPropertyWithChoices("EnumProperty 4", wx.PG_LABEL, soc.Ptr(), 240); p != nil {
		pg.Append(p.Ptr())
		if id := propByName(pg, "EnumProperty 4"); id != nil {
			wx.PGPropertyFromPtr(id).AddChoice("4th only", 360)
		}
	}
	helpString(pg, "EnumProperty 4",
		"Should have one extra item when compared to EnumProperty 3")

	pg.Append(wx.NewEnumPropertyWithChoices("EnumProperty With Bitmap", "EnumProperty 5", soc.Ptr(), 280).Ptr())
	helpString(pg, "EnumProperty 5", "Should have bitmap in front of the displayed value")
	if b := wx.ArtProvider_GetBitmapBundle(wx.ART_REMOVABLE, wx.ART_OTHER, 0, 0); b != nil {
		if p := propByName(pg, "EnumProperty 5"); p != nil {
			pg.SetPropertyImage(p, b.Ptr())
		}
		b.Delete()
	}

	pg.Append(wx.NewStringProperty("Password", wx.PG_LABEL, "password").Ptr())
	setBoolByName(pg, "Password", wx.PG_STRING_PASSWORD, true, 0)
	helpString(pg, "Password", "Has attribute wxPG_STRING_PASSWORD set to true")

	pg.Append(wx.NewDirProperty("DirProperty", wx.PG_LABEL, wx.GetUserHome("")).Ptr())
	setStringByName(pg, "DirProperty", wx.PG_DIALOG_TITLE,
		"This is a custom dir dialog title", 0)

	pg.Append(wx.NewStringProperty("StringProperty", wx.PG_LABEL, "").Ptr())
	if p := propByName(pg, "StringProperty"); p != nil {
		pg.SetPropertyMaxLength(p, 6)
	}
	helpString(pg, "StringProperty",
		"Max length of this text has been limited to 6, using wxPropertyGrid::SetPropertyMaxLength.")
	if p := propByName(pg, "StringProperty"); p != nil {
		pg.SetPropertyValueString(p, "some text")
	}

	// AutoComplete attribute.
	pg.Append(wx.NewStringProperty("StringProperty AutoComplete", wx.PG_LABEL, "").Ptr())
	autoArr := wx.NewArrayString()
	if autoArr != nil {
		for _, s := range []string{
			"One choice", "Another choice", "Another choice, yeah",
			"Yet another choice", "Yet another choice, bear with me",
		} {
			autoArr.Add(s)
		}
		av := wx.NewVariantFromArrayString(autoArr.Ptr(), "")
		if p := propByName(pg, "StringProperty AutoComplete"); p != nil && av != nil {
			pg.SetPropertyAttribute(p, wx.PG_ATTR_AUTOCOMPLETE, av.Ptr(), 0)
		}
		if av != nil {
			av.Delete()
		}
		autoArr.Delete()
	}
	helpString(pg, "StringProperty AutoComplete",
		"AutoComplete attribute has been set for this property "+
			"(try writing something beginning with 'a', 'o' or 'y').")

	// The C++ version draws two bitmaps for this property; the closest
	// standard art is used instead.
	pg.Append(wx.NewStringProperty("StringPropertyWithBitmap", wx.PG_LABEL, "Test Text").Ptr())
	if b := wx.ArtProvider_GetBitmapBundle(wx.ART_INFORMATION, wx.ART_OTHER, 0, 0); b != nil {
		if p := propByName(pg, "StringPropertyWithBitmap"); p != nil {
			pg.SetPropertyImage(p, b.Ptr())
		}
		b.Delete()
	}

	flagsLabels := []string{"Bit 0", "Bit 1", "Bit 2", "Bit 3"}
	flagsValues := []int32{0x01, 0x02, 0x04, 0x08}
	pg.Append(newFlagsProperty("FlagsProperty", wx.PG_LABEL, flagsLabels, flagsValues, 0x03).Ptr())
	setBoolByName(pg, "FlagsProperty", wx.PG_BOOL_USE_CHECKBOX, true, 0)

	tchoices := []string{"Cabbage", "Carrot", "Onion", "Potato", "Strawberry"}
	pg.Append(newEnumProperty("EnumProperty X", wx.PG_LABEL, tchoices,
		[]int32{0, 1, 2, 3, 4}, 0).Ptr())

	pg.Append(newMultiChoiceProp("MultiChoiceProperty", wx.PG_LABEL, tchoices,
		[]string{"Carrot", "Potato"}).Ptr())
	setLongByName(pg, "MultiChoiceProperty", wx.PG_ATTR_MULTICHOICE_USERSTRINGMODE, 1, 0)

	// wxSizeProperty / wxPointProperty are custom composite classes. KwxPGProperty
	// with Int children reproduces their structure.
	sz := f.frame.GetSize()
	sizeProp := newComposite("SizeProperty", "Size", fmt.Sprintf("%d x %d", sz.Width, sz.Height))
	if sizeProp != nil {
		sizeProp.AddPrivateChild(wx.NewIntProperty("Width", wx.PG_LABEL, sz.Width).Ptr())
		sizeProp.AddPrivateChild(wx.NewIntProperty("Height", wx.PG_LABEL, sz.Height).Ptr())
		pg.Append(sizeProp.Ptr())
	}
	pos := f.frame.GetPosition()
	pointProp := newComposite("PointProperty", "Position", fmt.Sprintf("(%d, %d)", pos.X, pos.Y))
	if pointProp != nil {
		pointProp.AddPrivateChild(wx.NewIntProperty("X", wx.PG_LABEL, pos.X).Ptr())
		pointProp.AddPrivateChild(wx.NewIntProperty("Y", wx.PG_LABEL, pos.Y).Ptr())
		pg.Append(pointProp.Ptr())
	}

	pg.Append(wx.NewUIntProperty("UIntProperty", wx.PG_LABEL, 0xFEEEFEEEFEEE).Ptr())
	setLongByName(pg, "UIntProperty", wx.PG_UINT_PREFIX, pgPrefixNone, 0)
	setLongByName(pg, "UIntProperty", wx.PG_UINT_BASE, pgBaseHex, 0)

	pg.Append(newEditEnumProperty("EditEnumProperty", wx.PG_LABEL,
		[]string{"Choice 1", "Choice 2", "Choice 3"}, []int32{0, 1, 2},
		"Choice not in the list").Ptr())
	if p := propByName(pg, "EditEnumProperty"); p != nil {
		v := wx.NewVariantFromString("Dummy Hint", "")
		wx.PGPropertyFromPtr(p).SetAttribute(wx.PG_ATTR_HINT, v.Ptr())
		v.Delete()
	}

	now := wx.NewDateTime()
	if now != nil {
		now.Now()
		pg.Append(wx.NewDateProperty("DateProperty", wx.PG_LABEL, now.Ptr()).Ptr())
		now.Delete()
	}
	if p := propByName(pg, "DateProperty"); p != nil {
		setAttrLong(pg, p, wx.PG_DATE_PICKER_STYLE,
			wx.DP_DROPDOWN|wx.DP_SHOWCENTURY|wx.DP_ALLOWNONE, 0)
	}
	helpString(pg, "DateProperty",
		"Attribute wxPG_DATE_PICKER_STYLE has been set to (long)"+
			"(wxDP_DROPDOWN | wxDP_SHOWCENTURY | wxDP_ALLOWNONE).")

	// 3D Object: custom vector/triangle properties are replaced by KwxPGProperty
	// composites.
	topID := pg.Append(wx.NewStringProperty("3D Object", wx.PG_LABEL, "<composed>").Ptr())

	pid := pg.AppendIn(topID, wx.NewStringProperty("Triangle 1", "Triangle 1", "<composed>").Ptr())
	if pid != nil {
		pg.AppendIn(pid, newVectorProperty("A", wx.PG_LABEL).Ptr())
		pg.AppendIn(pid, newVectorProperty("B", wx.PG_LABEL).Ptr())
		pg.AppendIn(pid, newVectorProperty("C", wx.PG_LABEL).Ptr())
	}
	pg.AppendIn(topID, newTriangleProperty("Triangle 2", "Triangle 2").Ptr())

	helpString(pg, "3D Object",
		"3D Object is wxStringProperty with value \"<composed>\". Two of its children are similar wxStringProperties with "+
			"three wxVectorProperty children, and other two are custom wxTriangleProperties.")

	pid = pg.AppendIn(topID, wx.NewStringProperty("Triangle 3", "Triangle 3", "<composed>").Ptr())
	if pid != nil {
		pg.AppendIn(pid, newVectorProperty("A", wx.PG_LABEL).Ptr())
		pg.AppendIn(pid, newVectorProperty("B", wx.PG_LABEL).Ptr())
		pg.AppendIn(pid, newVectorProperty("C", wx.PG_LABEL).Ptr())
	}
	pg.AppendIn(topID, newTriangleProperty("Triangle 4", "Triangle 4").Ptr())

	// Car (doc sample test).
	carProp := pg.Append(wx.NewStringProperty("Car", wx.PG_LABEL, "<composed>").Ptr())
	pg.AppendIn(carProp, wx.NewStringProperty("Model", wx.PG_LABEL,
		"Lamborghini Diablo SV").Ptr())
	pg.AppendIn(carProp, wx.NewIntProperty("Engine Size (cc)", wx.PG_LABEL, 5707).Ptr())
	speedsProp := pg.AppendIn(carProp, wx.NewStringProperty("Speeds", wx.PG_LABEL, "<composed>").Ptr())
	pg.AppendIn(speedsProp, wx.NewIntProperty("Max. Speed (mph)", wx.PG_LABEL, 290).Ptr())
	pg.AppendIn(speedsProp, wx.NewFloatProperty("0-100 mph (sec)", wx.PG_LABEL, 3.9).Ptr())
	pg.AppendIn(speedsProp, wx.NewFloatProperty("1/4 mile (sec)", wx.PG_LABEL, 8.6).Ptr())

	// This is how a child property can be referred to by name.
	if p := propByName(pg, "Car.Speeds.Max. Speed (mph)"); p != nil {
		pg.SetPropertyValueInt(p, 300)
	}

	pg.AppendIn(carProp, wx.NewIntProperty("Price ($)", wx.PG_LABEL, 300000).Ptr())
	pg.AppendIn(carProp, wx.NewBoolProperty("Convertible", wx.PG_LABEL, false).Ptr())

	// Test the custom multi-button editor.
	mb := wx.NewLongStringProperty("MultipleButtons", wx.PG_LABEL, "")
	pg.Append(mb.Ptr())
	if f.sampleEditor != nil {
		pg.SetPropertyEditor(mb.Ptr(), f.sampleEditor)
	}

	// SingleChoiceProperty (a wxStringProperty with a choice-dialog adapter).
	// The adapter needs a wxPGEditorDialogAdapter subclass, so a plain string
	// property with the same editor button is used.
	sc := wx.NewStringProperty("SingleChoiceProperty", wx.PG_LABEL, "")
	pg.Append(sc.Ptr())
	pg.SetPropertyEditorByName(sc.Ptr(), "TextCtrlAndButton")

	// Test adding variable height bitmaps in wxPGChoices.
	bc := wx.NewPGChoices()
	if bc != nil {
		defer bc.Delete()
		bcLabels := []string{"Wee", "Not so wee", "Friggin' huge"}
		bcIDs := []string{wx.ART_CDROM, wx.ART_FLOPPY, wx.ART_HARDDISK}
		bcSizes := []int{16, 32, 64}
		for i, label := range bcLabels {
			bc.Add(label, i)
		}
		for i, id := range bcIDs {
			if e := wx.PGChoiceEntryFromPtr(bc.Item(uint(i))); e != nil {
				if b := wx.ArtProvider_GetBitmapBundle(id, wx.ART_OTHER, bcSizes[i], bcSizes[i]); b != nil {
					e.SetBitmap(b.Ptr())
					b.Delete()
				}
			}
		}
		if p := wx.NewEnumPropertyWithChoices("Variable Height Bitmaps", wx.PG_LABEL, bc.Ptr(), 0); p != nil {
			pg.Append(p.Ptr())
		}
	}

	// Test how non-editable composite strings appear.
	traits := wx.NewStringProperty("wxWidgets Traits", wx.PG_LABEL, "<composed>")
	pg.SetPropertyReadOnly(traits.Ptr(), true, pgDontRecurse)
	traits.AppendChild(wx.NewStringProperty("Latest Release", wx.PG_LABEL, "3.1.2").Ptr())
	traits.AppendChild(wx.NewBoolProperty("Win API", wx.PG_LABEL, true).Ptr())
	pg.Append(traits.Ptr())
	pg.AppendIn(traits.Ptr(), wx.NewBoolProperty("QT", wx.PG_LABEL, true).Ptr())
	pg.AppendIn(traits.Ptr(), wx.NewBoolProperty("Cocoa", wx.PG_LABEL, true).Ptr())
	pg.AppendIn(traits.Ptr(), wx.NewBoolProperty("Haiku", wx.PG_LABEL, false).Ptr())
	pg.AppendIn(traits.Ptr(), wx.NewStringProperty("Trunk Version", wx.PG_LABEL,
		wx.VersionString()).Ptr())
	pg.AppendIn(traits.Ptr(), wx.NewBoolProperty("GTK+", wx.PG_LABEL, true).Ptr())
	pg.AppendIn(traits.Ptr(), wx.NewBoolProperty("Android", wx.PG_LABEL, false).Ptr())

	f.addTestProperties(pg)
}

// addTestProperties is FormMain::AddTestProperties().
func (f *mainFrame) addTestProperties(pg *wx.PropertyGridInterface) {
	// MyColourProperty customises the choice list of wxColourProperty; a plain
	// wxColourProperty is used instead.
	pg.Append(wx.NewColourProperty("CustomColourProperty", wx.PG_LABEL, wx.GREEN).Ptr())
	if p := propByName(pg, "CustomColourProperty"); p != nil {
		wx.PGPropertyFromPtr(p).SetAutoUnspecified(true)
		pg.SetPropertyEditorByName(p, "ComboBox")
	}
	helpString(pg, "CustomColourProperty",
		"This is a MyColourProperty from the sample app. "+
			"It is built by subclassing wxColourProperty.")
}
