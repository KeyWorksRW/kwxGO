package main

import (
	"fmt"
	"math"
	"strings"
	"unsafe"

	"github.com/KeyWorksRW/kwxGO/wx"
)

// gridPopup is the PropertyGridPopup struct in propgrid.cpp: a wxPopupWindow
// holding a scrolled window with its own three-column wxPropertyGrid.
type gridPopup struct {
	popup *wx.PopupWindow
	panel *wx.ScrolledWindow
	grid  *wx.PropertyGrid
	sizer *wx.BoxSizer
}

// onShowPopup is FormMain::OnShowPopup().
func (f *mainFrame) onShowPopup(evt *wx.Event) {
	if evt == nil {
		return
	}
	if f.popup != nil {
		if f.popup.popup != nil {
			f.popup.popup.Destroy()
		}
		f.popup = nil
		return
	}
	f.popup = newGridPopup(f.frame)
	if f.popup == nil {
		return
	}
	pt := wx.GetMousePositionPoint()
	f.popup.popup.Position(pt.X, pt.Y, 0, 0)
	f.popup.popup.Show()
}

func newGridPopup(parent *wx.Frame) *gridPopup {
	popup := wx.NewPopupWindow(parent.Ptr(), wx.BORDER_NONE|wx.WANTS_CHARS)
	if popup == nil {
		return nil
	}
	p := &gridPopup{popup: popup}

	p.panel = wx.NewScrolledWindow(popup.Ptr(), wx.ID_ANY, -1, -1, 200, 200, 0)
	if p.panel == nil {
		return nil
	}
	p.grid = wx.NewPropertyGrid(p.panel.Ptr(), idPopupGrid, -1, -1, 400, 400,
		wx.PG_SPLITTER_AUTO_CENTER)
	if p.grid == nil {
		return nil
	}
	p.grid.SetColumnCount(3)

	mi := pgIface(p.grid.Ptr())

	prop := p.grid.Append(wx.NewStringProperty("test_name", wx.PG_LABEL, "test_value").Ptr())
	if prop != nil {
		v := wx.NewVariantFromString("type", "")
		mi.SetPropertyAttribute(prop, wx.PG_ATTR_UNITS, v.Ptr(), 0)
		v.Delete()
	}
	prop1 := mi.AppendIn(prop, wx.NewStringProperty("sub_name1", wx.PG_LABEL, "sub_value1").Ptr())
	if prop1 != nil {
		mi.AppendIn(prop1, wx.NewColourProperty("Cell Colour", wx.PG_LABEL,
			p.grid.GetCellBackgroundColour()).Ptr())
	}
	prop2 := mi.AppendIn(prop, wx.NewStringProperty("sub_name2", wx.PG_LABEL, "sub_value2").Ptr())
	if prop2 != nil {
		mi.AppendIn(prop2, wx.NewStringProperty("sub_name21", wx.PG_LABEL, "sub_value21").Ptr())
	}

	// wxArrayDoubleProperty is a custom class; an array string property shows
	// the same values.
	mi.AppendIn(prop, newArrayStringProp("ArrayDoubleProperty", wx.PG_LABEL,
		[]string{"-1", "-0.5", "0", "0.5", "1"}).Ptr())

	font := wx.NewFontDefault()
	mi.AppendIn(prop, wx.NewFontProperty("Font", wx.PG_LABEL, font.Ptr()).Ptr())
	font.Delete()

	if prop2 != nil {
		mi.AppendIn(prop2, wx.NewStringProperty("sub_name22", wx.PG_LABEL, "sub_value22").Ptr())
		mi.AppendIn(prop2, wx.NewStringProperty("sub_name23", wx.PG_LABEL, "sub_value23").Ptr())
		wx.PGPropertyFromPtr(prop2).SetExpanded(false)
	}

	p.sizer = wx.NewBoxSizer(wx.VERTICAL)
	if p.sizer == nil {
		return nil
	}
	flags := wx.NewSizerFlags(0)
	flags.Expand()
	p.sizer.AddWindowFlags(p.grid.Ptr(), flags.Ptr())
	flags.Delete()

	p.panel.SetAutoLayout(true)
	p.panel.SetSizer(p.sizer.Ptr(), true)
	p.fit()

	p.popup.BindId(wx.EVT_PG_ITEM_COLLAPSED, idPopupGrid, func(evt *wx.Event) {
		wx.LogMessage("OnCollapse")
		p.fit()
	})
	p.popup.BindId(wx.EVT_PG_ITEM_EXPANDED, idPopupGrid, func(evt *wx.Event) {
		wx.LogMessage("OnExpand")
		p.fit()
	})
	return p
}

// fit is PropertyGridPopup::Fit().
func (p *gridPopup) fit() {
	if p == nil || p.grid == nil || p.sizer == nil || p.panel == nil || p.popup == nil {
		return
	}
	p.setMinSize()
	p.sizer.Fit(p.panel.Ptr())

	pos := p.popup.GetPosition()
	size := p.panel.GetSize()
	p.popup.SetSize(pos.X, pos.Y, size.Width, size.Height, wx.SIZE_USE_EXISTING)
}

// getColumnWidths is the C++ GetColumnWidths() helper: it measures the widest
// content in each of the three columns and descends into expanded children.
func getColumnWidths(state *wx.PropertyGridPageState, root unsafe.Pointer, width *[3]int) {
	if state == nil || root == nil {
		return
	}
	minWidths := [3]int{
		state.GetColumnMinWidth(0),
		state.GetColumnMinWidth(1),
		state.GetColumnMinWidth(2),
	}
	prop := wx.PGPropertyFromPtr(root)
	if prop == nil {
		return
	}
	count := prop.GetChildCount()
	for i := 0; i < count; i++ {
		child := prop.Item(uint(i))
		if child == nil {
			continue
		}
		for c := uint(0); c < 3; c++ {
			if w := state.GetColumnFullWidth(child, c); w > width[c] {
				width[c] = w
			}
		}
	}
	for i := 0; i < count; i++ {
		child := prop.Item(uint(i))
		if child == nil {
			continue
		}
		if cp := wx.PGPropertyFromPtr(child); cp != nil && cp.IsExpanded() {
			var w [3]int
			getColumnWidths(state, child, &w)
			for c := 0; c < 3; c++ {
				if w[c] > width[c] {
					width[c] = w[c]
				}
			}
		}
	}
	for c := 0; c < 3; c++ {
		if minWidths[c] > width[c] {
			width[c] = minWidths[c]
		}
	}
}

// setMinSize is the C++ SetMinSize() helper.
func (p *gridPopup) setMinSize() {
	root := p.grid.GetRoot()
	if root == nil {
		return
	}
	gi := pgIface(p.grid.Ptr())
	realRoot := gi.GetFirstChild(root)

	first := gi.GetFirst(wx.PG_ITERATE_ALL)
	last := p.grid.GetLastItem(wx.PG_ITERATE_DEFAULT)
	rect := p.grid.GetPropertyRect(first, last)
	height := rect.Height + 2*p.grid.GetVerticalSpacing()

	// Add some height when the root item is collapsed, so the vertical scroll
	// bar does not show.
	if realRoot != nil && !gi.IsPropertyExpanded(realRoot) {
		height += 2 * p.grid.GetVerticalSpacing()
	}

	var width [3]int
	state := wx.PropertyGridPageStateFromPtr(gi.GetState())
	getColumnWidths(state, root, &width)

	minWidth := wx.SystemSettings_GetMetric(wx.SYS_SCREEN_X) * 3 / 2
	minHeight := wx.SystemSettings_GetMetric(wx.SYS_SCREEN_Y) * 3 / 2

	sizeWidth := width[0] + width[1] + width[2] + p.grid.GetMarginWidth()
	if sizeWidth > minWidth {
		sizeWidth = minWidth
	}
	sizeHeight := height
	if sizeHeight > minHeight {
		sizeHeight = minHeight
	}
	p.grid.SetMinSize(sizeWidth, sizeHeight)

	if sizeWidth > 0 {
		proportions := [3]int{}
		proportions[0] = int(math.Round(float64(width[0]) / float64(sizeWidth) * 100.0))
		proportions[1] = int(math.Round(float64(width[1]) / float64(sizeWidth) * 100.0))
		proportions[2] = 100 - proportions[0] - proportions[1]
		if proportions[2] < 0 {
			proportions[2] = 0
		}
		gi.SetColumnProportion(0, proportions[0])
		gi.SetColumnProportion(1, proportions[1])
		gi.SetColumnProportion(2, proportions[2])
	}
}

// showDumpDialog is FormMain::OnDumpList().
func (f *mainFrame) onDumpList(evt *wx.Event) {
	if evt == nil {
		return
	}
	gi := f.gridIface()
	if gi == nil {
		return
	}
	values := wx.VariantFromPtr(gi.GetPropertyValues("list", nil, pgIncAttributes))
	text := "This only tests that wxVariant related routines do not crash.\n"
	if values != nil {
		count := values.GetCount()
		for i := 0; i < count; i++ {
			v := wx.VariantFromPtr(values.Item(i))
			if v == nil {
				continue
			}
			strValue := v.GetString()
			if strings.HasSuffix(v.GetName(), "@attr") {
				text += "Attributes:\n"
				attrCount := v.GetCount()
				for n := 0; n < attrCount; n++ {
					if a := wx.VariantFromPtr(v.Item(n)); a != nil {
						text += fmt.Sprintf("  attribute %d: name=%q  (type=%q  value=%q)\n",
							n, a.GetName(), a.GetType(), a.GetString())
						a.SetBorrowed(false)
						a.Delete()
					}
				}
			} else {
				text += fmt.Sprintf("%d: name=%q  type=%q  value=%q\n",
					i, v.GetName(), v.GetType(), strValue)
			}
			// Variant.Item returns a new variant; release the wrapper.
			v.SetBorrowed(false)
			v.Delete()
		}
		values.SetBorrowed(false)
		values.Delete()
	}

	dlg := wx.NewDialog(f.frame.Ptr(), wx.ID_ANY, "wxVariant Test", -1, -1, -1, -1,
		wx.DEFAULT_DIALOG_STYLE|wx.RESIZE_BORDER)
	if dlg == nil {
		return
	}

	spacing := 8
	topsizer := wx.NewBoxSizer(wx.VERTICAL)
	rowsizer := wx.NewBoxSizer(wx.HORIZONTAL)
	ed := wx.NewTextCtrl(dlg.Ptr(), wx.ID_ANY, text, -1, -1, -1, -1,
		wx.TE_MULTILINE|wx.TE_READONLY)
	if ed != nil {
		ef := wx.NewSizerFlags(1)
		ef.Expand()
		ef.BorderSize(wx.ALL, spacing)
		rowsizer.AddWindowFlags(ed.Ptr(), ef.Ptr())
		ef.Delete()
	}
	rf := wx.NewSizerFlags(1)
	rf.Expand()
	topsizer.AddSizerFlags(rowsizer.Ptr(), rf.Ptr())
	rf.Delete()

	rowsizer2 := wx.NewBoxSizer(wx.HORIZONTAL)
	if ok := wx.NewButton(dlg.Ptr(), wx.ID_OK, "Ok", -1, -1, -1, -1, 0); ok != nil {
		bf := wx.NewSizerFlags(0)
		bf.CenterHorizontal()
		bf.CenterVertical()
		bf.BorderSize(wx.BOTTOM|wx.LEFT|wx.RIGHT, spacing)
		rowsizer2.AddWindowFlags(ok.Ptr(), bf.Ptr())
		bf.Delete()
	}
	right := wx.NewSizerFlags(0)
	right.Right()
	topsizer.AddSizerFlags(rowsizer2.Ptr(), right.Ptr())
	right.Delete()

	dlg.SetSizer(topsizer.Ptr(), true)
	topsizer.SetSizeHints(dlg.Ptr())
	dlg.SetSize(-1, -1, 400, 300, wx.SIZE_USE_EXISTING)
	dlg.Centre(wx.BOTH)
	dlg.ShowModal()
	dlg.Delete()
}

// newMinimalFrame is DisplayMinimalFrame()/MyFrame in propgrid_minimal.cpp.
func newMinimalFrame(parent unsafe.Pointer) *wx.Frame {
	frame := wx.NewFrame(parent, wx.ID_ANY, "PropertyGrid Test", -1, -1, -1, -1,
		wx.DEFAULT_FRAME_STYLE)
	if frame == nil {
		return nil
	}

	menu := wx.NewMenu("", 0)
	if menu != nil {
		menu.Append(wx.ID_HIGHEST, "Action", "", false)
		mb := wx.NewMenuBar(0)
		if mb != nil {
			mb.Append(menu.Ptr(), "Action")
			frame.SetMenuBar(mb.Ptr())
		}
	}

	pg := wx.NewPropertyGrid(frame.Ptr(), wx.ID_ANY, -1, -1, 400, 400,
		wx.PG_SPLITTER_AUTO_CENTER|wx.PG_BOLD_MODIFIED)
	if pg == nil {
		return nil
	}
	gi := pgIface(pg.Ptr())
	gi.Append(wx.NewStringProperty("String Property", wx.PG_LABEL, "").Ptr())
	gi.Append(wx.NewIntProperty("Int Property", wx.PG_LABEL, 0).Ptr())
	gi.Append(wx.NewBoolProperty("Bool Property", wx.PG_LABEL, false).Ptr())

	frame.SetClientSize(frame.FromDIP(400), frame.FromDIP(300))

	frame.Bind(wx.EVT_PG_CHANGED, func(evt *wx.Event) {
		pe := wx.AsPropertyGridEvent(evt)
		if pe == nil {
			return
		}
		p := wx.PGPropertyFromPtr(pe.GetProperty())
		if p != nil {
			wx.LogVerbose("OnPropertyGridChange(" + p.GetName() + ", value=" +
				p.GetValueAsString() + ")")
		} else {
			wx.LogVerbose("OnPropertyGridChange(nullptr)")
		}
	})
	frame.Bind(wx.EVT_PG_CHANGING, func(evt *wx.Event) {
		pe := wx.AsPropertyGridEvent(evt)
		if pe == nil {
			return
		}
		if p := wx.PGPropertyFromPtr(pe.GetProperty()); p != nil {
			wx.LogVerbose("OnPropertyGridChanging(" + p.GetName() + ")")
		}
	})
	frame.BindId(wx.EVT_MENU, wx.ID_HIGHEST, func(evt *wx.Event) {})

	frame.Show()
	return frame
}

// onRunMinimalClick is FormMain::OnRunMinimalClick().
func (f *mainFrame) onRunMinimalClick(evt *wx.Event) {
	if evt == nil {
		return
	}
	newMinimalFrame(f.frame.Ptr())
}

// onTestXRC is FormMain::OnTestXRC().
func (f *mainFrame) onTestXRC(evt *wx.Event) {
	if evt == nil {
		return
	}
	xml := wx.XmlResourceFromPtr(wx.XmlResource_Get())
	if xml == nil {
		wx.MessageBox("Sorry, XRC support is not available in this build.", "", wx.OK,
			f.frame.Ptr(), -1, -1)
		return
	}
	if !f.xrcHandlersAdded {
		f.xrcHandlersAdded = true
		xml.InitAllHandlers()
		// This needs to be done in addition to initializing all the standard
		// handlers above.
		if h := wx.NewPropertyGridXmlHandler(); h != nil {
			xml.AddHandler(h.Ptr())
		}
	}
	if !xml.Load("propgrid.xrc") {
		wx.LogError("XRC file couldn't be loaded.")
		return
	}
	ptr := xml.LoadFrame(f.frame.Ptr(), "propgrid")
	if ptr == nil {
		wx.LogError("Property grid frame couldn't be loaded from XRC")
		return
	}
	if frame := wx.FrameFromPtr(ptr); frame != nil {
		frame.Show()
	}
}
