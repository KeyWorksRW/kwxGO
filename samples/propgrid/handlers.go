package main

import (
	"fmt"
	"math/rand"
	"unsafe"

	"github.com/KeyWorksRW/kwxGO/wx"
)

// bindEvents wires the C++ event table onto the frame.
func (f *mainFrame) bindEvents() {
	fr := f.frame

	fr.Bind(wx.EVT_IDLE, f.onIdle)
	fr.Bind(wx.EVT_MOVE, f.onMove)
	fr.Bind(wx.EVT_SIZE, f.onResize)

	// Property grid events (all use the manager id).
	fr.BindId(wx.EVT_PG_SELECTED, idPGID, f.onPGSelect)
	fr.BindId(wx.EVT_PG_CHANGED, idPGID, f.onPGChange)
	fr.BindId(wx.EVT_PG_CHANGING, idPGID, f.onPGChanging)
	fr.BindId(wx.EVT_PG_HIGHLIGHTED, idPGID, f.onPGHighlight)
	fr.BindId(wx.EVT_PG_RIGHT_CLICK, idPGID, f.onPGRightClick)
	fr.BindId(wx.EVT_PG_DOUBLE_CLICK, idPGID, f.onPGDoubleClick)
	fr.BindId(wx.EVT_PG_PAGE_CHANGED, idPGID, f.onPGPageChange)
	fr.BindId(wx.EVT_PG_LABEL_EDIT_BEGIN, idPGID, f.onPGLabelEditBegin)
	fr.BindId(wx.EVT_PG_LABEL_EDIT_ENDING, idPGID, f.onPGLabelEditEnding)
	fr.BindId(wx.EVT_PG_ITEM_COLLAPSED, idPGID, f.onPGItemCollapse)
	fr.BindId(wx.EVT_PG_ITEM_EXPANDED, idPGID, f.onPGItemExpand)
	fr.BindId(wx.EVT_PG_COL_BEGIN_DRAG, idPGID, f.onPGColBeginDrag)
	fr.BindId(wx.EVT_PG_COL_DRAGGING, idPGID, f.onPGColDragging)
	fr.BindId(wx.EVT_PG_COL_END_DRAG, idPGID, f.onPGColEndDrag)

	fr.BindId(wx.EVT_BUTTON, idPGID, f.onPGButtonClick)
	fr.BindId(wx.EVT_TEXT, idPGID, f.onPGTextUpdate)
	fr.Bind(wx.EVT_KEY_DOWN, f.onPGKeyEvent)
	fr.Bind(wx.EVT_KEY_UP, f.onPGKeyEvent)

	// Basic menu.
	fr.BindId(wx.EVT_MENU, idAppendProp, f.onAppendPropClick)
	fr.BindId(wx.EVT_MENU, idAppendCat, f.onAppendCatClick)
	fr.BindId(wx.EVT_MENU, idInsertProp, f.onInsertPropClick)
	fr.BindId(wx.EVT_MENU, idInsertCat, f.onInsertCatClick)
	fr.BindId(wx.EVT_MENU, idDelete, f.onDelPropClick)
	fr.BindId(wx.EVT_MENU, idDeleteR, f.onDelPropRClick)
	fr.BindId(wx.EVT_MENU, idUnspecify, f.onMisc)
	fr.BindId(wx.EVT_MENU, idDeleteAll, f.onClearClick)
	fr.BindId(wx.EVT_MENU, idEnable, f.onEnableDisable)
	fr.BindId(wx.EVT_MENU, idSetReadOnly, f.onSetReadOnly)
	fr.BindId(wx.EVT_MENU, idHide, f.onHide)
	fr.BindId(wx.EVT_MENU, idBoolCheckbox, f.onBoolCheckbox)

	// Advanced menu.
	fr.BindId(wx.EVT_MENU, idIterate1, f.onIterate1Click)
	fr.BindId(wx.EVT_MENU, idIterate2, f.onIterate2Click)
	fr.BindId(wx.EVT_MENU, idIterate3, f.onIterate3Click)
	fr.BindId(wx.EVT_MENU, idIterate4, f.onIterate4Click)
	fr.BindId(wx.EVT_MENU, idOnExtendedKeyNav, f.onExtendedKeyNav)
	fr.BindId(wx.EVT_MENU, idSetBgColour, f.onSetBackgroundColour)
	fr.BindId(wx.EVT_MENU, idSetBgColourRecur, f.onSetBackgroundColour)
	fr.BindId(wx.EVT_MENU, idClearModif, f.onClearModifyStatusClick)
	fr.BindId(wx.EVT_MENU, idFreeze, f.onFreezeClick)
	fr.BindId(wx.EVT_MENU, idEnableLabelEditing, f.onEnableLabelEditing)
	fr.BindId(wx.EVT_MENU, idShowHeader, f.onShowHeader)
	fr.BindId(wx.EVT_MENU, idDumpList, f.onDumpList)
	fr.BindId(wx.EVT_MENU, idStaticLayout, f.onMisc)
	fr.BindId(wx.EVT_MENU, idCollapse, f.onMisc)
	fr.BindId(wx.EVT_MENU, idCollapseAll, f.onMisc)
	fr.BindId(wx.EVT_MENU, idGetValues, f.onMisc)
	fr.BindId(wx.EVT_MENU, idSetValues, f.onMisc)
	fr.BindId(wx.EVT_MENU, idSetValues2, f.onMisc)
	fr.BindId(wx.EVT_MENU, idFitColumns, f.onFitColumnsClick)
	fr.BindId(wx.EVT_MENU, idChangeFlagsItems, f.onChangeFlagsPropItemsClick)
	fr.BindId(wx.EVT_MENU, idTestInsertChoice, f.onInsertChoice)
	fr.BindId(wx.EVT_MENU, idTestDeleteChoice, f.onDeleteChoice)
	fr.BindId(wx.EVT_MENU, idInsertPage, f.onInsertPage)
	fr.BindId(wx.EVT_MENU, idRemovePage, f.onRemovePage)
	fr.BindId(wx.EVT_MENU, idSaveState, f.onSaveState)
	fr.BindId(wx.EVT_MENU, idRestoreState, f.onRestoreState)
	fr.BindId(wx.EVT_MENU, idSetSpinCtrlEditor, f.onSetSpinCtrlEditorClick)
	fr.BindId(wx.EVT_MENU, idTestReplace, f.onTestReplaceClick)
	fr.BindId(wx.EVT_MENU, idSetPropertyValue, f.onSetPropertyValue)
	fr.BindId(wx.EVT_MENU, idEnableCommonValues, f.onEnableCommonValues)
	fr.BindId(wx.EVT_MENU, idPopulate1, f.onPopulateClick)
	fr.BindId(wx.EVT_MENU, idPopulate2, f.onPopulateClick)
	fr.BindId(wx.EVT_MENU, idSetColumns, f.onSetColumns)
	fr.BindId(wx.EVT_MENU, idSetVirtWidth, f.onSetVirtualWidth)
	fr.BindId(wx.EVT_MENU, idSetPGDisabled, f.onSetGridDisabled)
	fr.BindId(wx.EVT_MENU, idTestXRC, f.onTestXRC)
	fr.BindId(wx.EVT_MENU, idSelectStyle, f.onSelectStyle)

	// Try These! / File / Help menus.
	fr.BindId(wx.EVT_MENU, idColourScheme1, f.onColourScheme)
	fr.BindId(wx.EVT_MENU, idColourScheme2, f.onColourScheme)
	fr.BindId(wx.EVT_MENU, idColourScheme3, f.onColourScheme)
	fr.BindId(wx.EVT_MENU, idColourScheme4, f.onColourScheme)
	fr.BindId(wx.EVT_MENU, idCatColours, f.onCatColours)
	fr.BindId(wx.EVT_MENU, idAbout, f.onAbout)
	fr.BindId(wx.EVT_MENU, idQuit, f.onCloseClick)
	fr.BindId(wx.EVT_MENU, idRunMinimal, f.onRunMinimalClick)

	fr.BindId(wx.EVT_UPDATE_UI, idCatColours, f.onCatColoursUpdateUI)
	fr.Bind(wx.EVT_CONTEXT_MENU, f.onContextMenu)
	fr.BindId(wx.EVT_BUTTON, idShowPopup, f.onShowPopup)
}

func commandEvent(evt *wx.Event) *wx.CommandEvent {
	if evt == nil {
		return nil
	}
	return wx.AsCommandEvent(evt)
}

func propertyGridEvent(evt *wx.Event) *wx.PropertyGridEvent {
	if evt == nil {
		return nil
	}
	return wx.AsPropertyGridEvent(evt)
}

func (f *mainFrame) managerIface() *wx.PropertyGridInterface {
	if f.manager == nil {
		return nil
	}
	return wx.PropertyGridInterfaceFromPtr(f.manager.Ptr())
}

func (f *mainFrame) gridIface() *wx.PropertyGridInterface {
	if f.grid == nil {
		return nil
	}
	return wx.PropertyGridInterfaceFromPtr(f.grid.Ptr())
}

// onMove is FormMain::OnMove().
func (f *mainFrame) onMove(evt *wx.Event) {
	if evt == nil {
		return
	}
	if f.manager == nil {
		evt.Skip(true)
		return
	}
	pos := f.frame.GetPosition()
	mi := f.managerIface()
	if p := propByName(mi, "X"); p != nil {
		mi.SetPropertyValueInt(p, pos.X)
	}
	if p := propByName(mi, "Y"); p != nil {
		mi.SetPropertyValueInt(p, pos.Y)
	}
	if p := propByName(mi, "Position"); p != nil {
		mi.SetPropertyValueString(p, fmt.Sprintf("(%d, %d)", pos.X, pos.Y))
		prop := wx.PGPropertyFromPtr(p)
		if c := prop.Item(0); c != nil {
			mi.SetPropertyValueInt(c, pos.X)
		}
		if c := prop.Item(1); c != nil {
			mi.SetPropertyValueInt(c, pos.Y)
		}
	}
	evt.Skip(true)
}

// onResize is FormMain::OnResize().
func (f *mainFrame) onResize(evt *wx.Event) {
	if evt == nil {
		return
	}
	if f.manager == nil {
		evt.Skip(true)
		return
	}
	sz := f.frame.GetSize()
	mi := f.managerIface()
	if p := propByName(mi, "Width"); p != nil && !wx.PGPropertyFromPtr(p).IsValueUnspecified() {
		mi.SetPropertyValueInt(p, sz.Width)
	}
	if p := propByName(mi, "Height"); p != nil && !wx.PGPropertyFromPtr(p).IsValueUnspecified() {
		mi.SetPropertyValueInt(p, sz.Height)
	}
	if p := propByName(mi, "Size"); p != nil {
		mi.SetPropertyValueString(p, fmt.Sprintf("%d x %d", sz.Width, sz.Height))
		prop := wx.PGPropertyFromPtr(p)
		if c := prop.Item(0); c != nil {
			mi.SetPropertyValueInt(c, sz.Width)
		}
		if c := prop.Item(1); c != nil {
			mi.SetPropertyValueInt(c, sz.Height)
		}
	}
	evt.Skip(true)
}

// onPGChanging is FormMain::OnPropertyGridChanging().
func (f *mainFrame) onPGChanging(evt *wx.Event) {
	pe := propertyGridEvent(evt)
	if pe == nil {
		return
	}
	prop := wx.PGPropertyFromPtr(pe.GetProperty())
	if prop == nil || prop.GetName() != "Font" {
		return
	}
	typ := ""
	if v := wx.VariantFromPtr(pe.GetValue()); v != nil {
		typ = v.GetType()
	}
	msg := fmt.Sprintf("'%s' is about to change (to variant of type '%s')\n\nAllow or deny?",
		prop.GetName(), typ)
	res := wx.MessageBox(msg, "Testing wxEVT_PG_CHANGING", wx.YES_NO, f.manager.Ptr(), -1, -1)
	if res == wx.NO {
		// Since we ask a question, it is better if we omit any validation
		// failure behaviour.
		pe.Veto(true)
		pe.SetValidationFailureBehavior(wx.PG_VFB_NULL)
	}
}

// onPGChange is FormMain::OnPropertyGridChange().
func (f *mainFrame) onPGChange(evt *wx.Event) {
	pe := propertyGridEvent(evt)
	if pe == nil {
		return
	}
	prop := wx.PGPropertyFromPtr(pe.GetProperty())
	if prop == nil {
		return
	}
	name := prop.GetName()
	v := wx.VariantFromPtr(pe.GetValue())
	if v == nil || v.IsNull() {
		return
	}

	switch name {
	case "X":
		f.frame.SetSize(v.GetLong(), -1, -1, -1, wx.SIZE_USE_EXISTING)
	case "Y":
		f.frame.SetSize(-1, v.GetLong(), -1, -1, wx.SIZE_USE_EXISTING)
	case "Width":
		f.frame.SetSize(-1, -1, v.GetLong(), -1, wx.SIZE_USE_EXISTING)
	case "Height":
		f.frame.SetSize(-1, -1, -1, v.GetLong(), wx.SIZE_USE_EXISTING)
	case "Label":
		f.frame.SetTitle(v.GetString())
	case "Password":
		// Do something very secret here...
	case "Font":
		// wxVariant::GetFont() returns a new wxFont; take ownership of the
		// wrapper so Delete() releases it after the manager copies it.
		if ptr := v.GetFont(); ptr != nil {
			if fnt := wx.FontFromPtr(ptr); fnt != nil {
				f.manager.SetFont(fnt.Ptr())
				fnt.SetBorrowed(false)
				fnt.Delete()
			}
		}
	case "Margin Colour":
		if cpv := wx.ColourPropertyValueFromPtr(v.GetColourPropertyValue()); cpv != nil {
			f.grid.SetMarginColour(cpv.GetColour())
		}
	case "Cell Colour":
		if cpv := wx.ColourPropertyValueFromPtr(v.GetColourPropertyValue()); cpv != nil {
			f.grid.SetCellBackgroundColour(cpv.GetColour())
		}
	case "Line Colour":
		if cpv := wx.ColourPropertyValueFromPtr(v.GetColourPropertyValue()); cpv != nil {
			f.grid.SetLineColour(cpv.GetColour())
		}
	case "Cell Text Colour":
		if cpv := wx.ColourPropertyValueFromPtr(v.GetColourPropertyValue()); cpv != nil {
			f.grid.SetCellTextColour(cpv.GetColour())
		}
	}
}

// onPGSelect is FormMain::OnPropertyGridSelect().
func (f *mainFrame) onPGSelect(evt *wx.Event) {
	pe := propertyGridEvent(evt)
	if pe == nil {
		return
	}
	propPtr := pe.GetProperty()
	if propPtr != nil {
		if f.itemEnable != nil {
			f.itemEnable.Enable(true)
			if wx.PGPropertyFromPtr(propPtr).IsEnabled() {
				f.itemEnable.SetItemLabel("Disable")
			} else {
				f.itemEnable.SetItemLabel("Enable")
			}
		}
		if mi := f.managerIface(); mi != nil {
			f.setStatusText("Selected: " + mi.GetPropertyLabel(propPtr))
		}
	} else if f.itemEnable != nil {
		f.itemEnable.Enable(false)
	}
}

// onPGPageChange is FormMain::OnPropertyGridPageChange().
func (f *mainFrame) onPGPageChange(evt *wx.Event) {
	if evt == nil || f.manager == nil {
		return
	}
	f.setStatusText("Page Changed: " + f.manager.GetPageName(f.manager.GetSelectedPage()))
}

func (f *mainFrame) onPGLabelEditBegin(evt *wx.Event) {
	pe := propertyGridEvent(evt)
	if pe == nil {
		return
	}
	if prop := wx.PGPropertyFromPtr(pe.GetProperty()); prop != nil {
		wx.LogMessage("wxPG_EVT_LABEL_EDIT_BEGIN(" + prop.GetLabel() + ")")
	}
}

func (f *mainFrame) onPGLabelEditEnding(evt *wx.Event) {
	pe := propertyGridEvent(evt)
	if pe == nil {
		return
	}
	if prop := wx.PGPropertyFromPtr(pe.GetProperty()); prop != nil {
		wx.LogMessage("wxPG_EVT_LABEL_EDIT_ENDING(" + prop.GetLabel() + ")")
	}
}

func (f *mainFrame) onPGHighlight(evt *wx.Event) {
}

func (f *mainFrame) onPGRightClick(evt *wx.Event) {
	pe := propertyGridEvent(evt)
	if pe == nil {
		return
	}
	propPtr := pe.GetProperty()
	if propPtr == nil {
		f.setStatusText("")
		return
	}
	text := "Right-clicked: " + wx.PGPropertyFromPtr(propPtr).GetLabel() + ", name="
	if mi := f.managerIface(); mi != nil {
		text += mi.GetPropertyName(propPtr)
	}
	f.setStatusText(text)
}

func (f *mainFrame) onPGDoubleClick(evt *wx.Event) {
	pe := propertyGridEvent(evt)
	if pe == nil {
		return
	}
	propPtr := pe.GetProperty()
	if propPtr == nil {
		f.setStatusText("")
		return
	}
	text := "Double-clicked: " + wx.PGPropertyFromPtr(propPtr).GetLabel() + ", name="
	if mi := f.managerIface(); mi != nil {
		text += mi.GetPropertyName(propPtr)
	}
	f.setStatusText(text)
}

func (f *mainFrame) onPGButtonClick(evt *wx.Event) {
	mi := f.managerIface()
	if mi == nil {
		return
	}
	prop := mi.GetSelection()
	if prop == nil {
		wx.MessageBox("SHOULD NOT HAPPEN!!!", "", wx.OK, f.frame.Ptr(), -1, -1)
		return
	}
	f.setStatusText("Button clicked: " + mi.GetPropertyLabel(prop) + ", name=" +
		mi.GetPropertyName(prop))
}

func (f *mainFrame) onPGItemCollapse(evt *wx.Event) {
	wx.LogMessage("Item was Collapsed")
}

func (f *mainFrame) onPGItemExpand(evt *wx.Event) {
	wx.LogMessage("Item was Expanded")
}

func (f *mainFrame) onPGColBeginDrag(evt *wx.Event) {
	pe := propertyGridEvent(evt)
	if pe == nil {
		return
	}
	if f.itemVetoDrag != nil && f.itemVetoDrag.IsChecked() {
		wx.LogMessage(fmt.Sprintf("Splitter %d resize was vetoed", pe.GetColumn()))
		pe.Veto(true)
	} else {
		wx.LogDebug(fmt.Sprintf("Splitter %d resize began", pe.GetColumn()))
	}
}

func (f *mainFrame) onPGColDragging(evt *wx.Event) {
}

func (f *mainFrame) onPGColEndDrag(evt *wx.Event) {
	pe := propertyGridEvent(evt)
	if pe == nil {
		return
	}
	wx.LogDebug(fmt.Sprintf("Splitter %d resize ended", pe.GetColumn()))
}

func (f *mainFrame) onPGTextUpdate(evt *wx.Event) {
	if evt != nil {
		evt.Skip(true)
	}
}

func (f *mainFrame) onPGKeyEvent(evt *wx.Event) {
}

func (f *mainFrame) onContextMenu(evt *wx.Event) {
	// FormMain::OnContextMenu() only logs the position and does not skip.
}

func (f *mainFrame) onIdle(evt *wx.Event) {
	if evt != nil {
		evt.Skip(true)
	}
}

// onCloseClick is FormMain::OnCloseClick().
func (f *mainFrame) onCloseClick(evt *wx.Event) {
	if evt == nil {
		return
	}
	f.frame.Close(false)
}

// --- Basic menu ------------------------------------------------------------

func (f *mainFrame) onInsertPropClick(evt *wx.Event) {
	mi := f.managerIface()
	if mi == nil {
		return
	}
	if !wx.PGPropertyFromPtr(f.grid.GetRoot()).HasAnyChild() {
		wx.MessageBox("No items to relate - first add some with Append.", "", wx.OK, nil, -1, -1)
		return
	}
	id := mi.GetSelection()
	if id == nil {
		wx.MessageBox("First select a property - new one will be inserted right before that.",
			"", wx.OK, nil, -1, -1)
		return
	}
	label := uniqueLabel(mi, "Property")
	mi.Insert2(mi.GetPropertyParent(id), int(wx.PGPropertyFromPtr(id).GetIndexInParent()),
		wx.NewStringProperty(label, wx.PG_LABEL, "").Ptr())
}

func (f *mainFrame) onAppendPropClick(evt *wx.Event) {
	mi := f.managerIface()
	if mi == nil {
		return
	}
	label := uniqueLabel(mi, "Property")
	mi.Append(wx.NewStringProperty(label, wx.PG_LABEL, "").Ptr())
	f.manager.Refresh(true)
}

func (f *mainFrame) onClearClick(evt *wx.Event) {
	if gi := f.gridIface(); gi != nil {
		gi.Clear()
	}
}

func (f *mainFrame) onAppendCatClick(evt *wx.Event) {
	mi := f.managerIface()
	if mi == nil {
		return
	}
	label := uniqueLabel(mi, "Category")
	mi.Append(wx.NewPropertyCategory(label).Ptr())
	f.manager.Refresh(true)
}

func (f *mainFrame) onInsertCatClick(evt *wx.Event) {
	mi := f.managerIface()
	if mi == nil {
		return
	}
	if !wx.PGPropertyFromPtr(f.grid.GetRoot()).HasAnyChild() {
		wx.MessageBox("No items to relate - first add some with Append.", "", wx.OK, nil, -1, -1)
		return
	}
	id := mi.GetSelection()
	if id == nil {
		wx.MessageBox("First select a property - new one will be inserted right before that.",
			"", wx.OK, nil, -1, -1)
		return
	}
	label := uniqueLabel(mi, "Category")
	mi.Insert2(mi.GetPropertyParent(id), int(wx.PGPropertyFromPtr(id).GetIndexInParent()),
		wx.NewPropertyCategory(label).Ptr())
}

func (f *mainFrame) onDelPropClick(evt *wx.Event) {
	mi := f.managerIface()
	if mi == nil {
		return
	}
	id := mi.GetSelection()
	if id == nil {
		wx.MessageBox("First select a property.", "", wx.OK, nil, -1, -1)
		return
	}
	mi.DeleteProperty(id)
}

func (f *mainFrame) onDelPropRClick(evt *wx.Event) {
	mi := f.managerIface()
	if mi == nil {
		return
	}
	p := f.grid.GetRoot()
	for p != nil {
		pp := wx.PGPropertyFromPtr(p)
		if !pp.HasAnyChild() {
			break
		}
		n := pp.GetChildCount()
		if n == 0 {
			break
		}
		p = pp.Item(uint(rand.Intn(n)))
		cp := wx.PGPropertyFromPtr(p)
		if !cp.IsCategory() {
			label := cp.GetLabel()
			mi.DeleteProperty(p)
			wx.LogMessage("Property deleted: " + label)
			break
		}
	}
}

func (f *mainFrame) onEnableDisable(evt *wx.Event) {
	mi := f.managerIface()
	if mi == nil {
		return
	}
	id := mi.GetSelection()
	if id == nil {
		wx.MessageBox("First select a property.", "", wx.OK, nil, -1, -1)
		return
	}
	if mi.IsPropertyEnabled(id) {
		mi.DisableProperty(id)
		if f.itemEnable != nil {
			f.itemEnable.SetItemLabel("Enable")
		}
	} else {
		mi.EnableProperty(id, true)
		if f.itemEnable != nil {
			f.itemEnable.SetItemLabel("Disable")
		}
	}
}

func (f *mainFrame) onSetReadOnly(evt *wx.Event) {
	mi := f.managerIface()
	if mi == nil {
		return
	}
	p := mi.GetSelection()
	if p == nil {
		wx.MessageBox("First select a property.", "", wx.OK, nil, -1, -1)
		return
	}
	mi.SetPropertyReadOnly(p, true, pgDontRecurse)
}

func (f *mainFrame) onHide(evt *wx.Event) {
	mi := f.managerIface()
	if mi == nil {
		return
	}
	id := mi.GetSelection()
	if id == nil {
		wx.MessageBox("First select a property.", "", wx.OK, nil, -1, -1)
		return
	}
	mi.HideProperty(id, true, pgDontRecurse)
}

func (f *mainFrame) onBoolCheckbox(evt *wx.Event) {
	cmd := commandEvent(evt)
	mi := f.managerIface()
	if cmd == nil || mi == nil {
		return
	}
	v := wx.NewVariantFromBool(cmd.IsChecked(), "")
	mi.SetPropertyAttributeAll(wx.PG_BOOL_USE_CHECKBOX, v.Ptr())
	v.Delete()
}

func (f *mainFrame) onSetBackgroundColour(evt *wx.Event) {
	cmd := commandEvent(evt)
	gi := f.gridIface()
	if cmd == nil || gi == nil {
		return
	}
	prop := gi.GetSelection()
	if prop == nil {
		wx.MessageBox("First select a property.", "", wx.OK, nil, -1, -1)
		return
	}
	white := wx.NewColourRGB(255, 255, 255, 255)
	defer white.Delete()
	col := wx.NewColourEmpty()
	defer col.Delete()
	wx.GetColourFromUserEx(f.frame.Ptr(), white.Ptr(), "Choose colour", col.Ptr())
	if !col.IsOk() {
		return
	}
	flags := pgDontRecurse
	if cmd.GetId() == idSetBgColourRecur {
		flags = pgRecurse
	}
	gi.SetPropertyBackgroundColour(prop, col.Ptr(), flags)
}

// --- Advanced menu ---------------------------------------------------------

func (f *mainFrame) onPopulateClick(evt *wx.Event) {
	cmd := commandEvent(evt)
	gi := f.gridIface()
	if cmd == nil || gi == nil {
		return
	}
	gi.Clear()
	f.grid.Freeze()
	switch cmd.GetId() {
	case idPopulate1:
		f.populateWithStandardItems()
	case idPopulate2:
		f.populateWithLibraryConfig()
	}
	f.grid.Thaw()
}

func (f *mainFrame) onClearModifyStatusClick(evt *wx.Event) {
	mi := f.managerIface()
	if mi == nil {
		return
	}
	mi.ClearModifiedStatus()
	f.manager.Refresh(true)
}

func (f *mainFrame) onFreezeClick(evt *wx.Event) {
	cmd := commandEvent(evt)
	if cmd == nil || f.manager == nil {
		return
	}
	if cmd.IsChecked() {
		if !f.manager.IsFrozen() {
			f.manager.Freeze()
		}
	} else {
		if f.manager.IsFrozen() {
			f.manager.Thaw()
			f.manager.Refresh(true)
		}
	}
}

func (f *mainFrame) onEnableLabelEditing(evt *wx.Event) {
	cmd := commandEvent(evt)
	if cmd == nil {
		return
	}
	f.labelEditingEnabled = cmd.IsChecked()
	f.grid.MakeColumnEditable(0, f.labelEditingEnabled)
}

func (f *mainFrame) onShowHeader(evt *wx.Event) {
	cmd := commandEvent(evt)
	if cmd == nil || f.manager == nil {
		return
	}
	f.hasHeader = cmd.IsChecked()
	f.manager.ShowHeader(f.hasHeader)
	if f.hasHeader {
		f.manager.SetColumnTitle(2, "Units")
	}
}

func (f *mainFrame) onInsertPage(evt *wx.Event) {
	if f.manager != nil {
		f.manager.AddPage("New Page", nil)
	}
}

func (f *mainFrame) onRemovePage(evt *wx.Event) {
	if f.manager != nil {
		f.manager.RemovePage(f.manager.GetSelectedPage())
	}
}

func (f *mainFrame) onSaveState(evt *wx.Event) {
	mi := f.managerIface()
	if mi == nil {
		return
	}
	f.savedState = mi.SaveEditableState(pgAllStates)
	wx.LogDebug("Saved editable state string: \"" + f.savedState + "\"")
}

func (f *mainFrame) onRestoreState(evt *wx.Event) {
	mi := f.managerIface()
	if mi == nil {
		return
	}
	mi.RestoreEditableState(f.savedState, pgAllStates)
}

func (f *mainFrame) onSetSpinCtrlEditorClick(evt *wx.Event) {
	mi := f.managerIface()
	if mi == nil {
		return
	}
	if id := mi.GetSelection(); id != nil {
		mi.SetPropertyEditorByName(id, "SpinCtrl")
	} else {
		wx.MessageBox("First select a property", "", wx.OK, nil, -1, -1)
	}
}

func (f *mainFrame) onTestReplaceClick(evt *wx.Event) {
	mi := f.managerIface()
	if mi == nil {
		return
	}
	id := mi.GetSelection()
	if id == nil {
		wx.MessageBox("First select a property", "", wx.OK, nil, -1, -1)
		return
	}
	labels := []string{"Flag 0", "Flag 1", "Flag 2", "Flag 3"}
	values := []int32{0x0001, 0x0002, 0x0004, 0x0008}
	propName := "ReplaceFlagsProperty"
	for idx := 1; propByName(mi, propName) != nil; idx++ {
		propName = fmt.Sprintf("ReplaceFlagsProperty %d", idx)
	}
	propVal := wx.GetLocalTime()%0x000F + 1
	newID := mi.ReplaceProperty(id,
		newFlagsProperty(propName, wx.PG_LABEL, labels, values, propVal).Ptr())
	if newID != nil {
		setAttrBool(mi, newID, wx.PG_BOOL_USE_CHECKBOX, true, pgRecurse)
		f.grid.SelectProperty(newID, true)
	}
}

func (f *mainFrame) onEnableCommonValues(evt *wx.Event) {
	mi := f.managerIface()
	if mi == nil {
		return
	}
	if prop := mi.GetSelection(); prop != nil {
		wx.PGPropertyFromPtr(prop).EnableCommonValue(true)
	} else {
		wx.MessageBox("First select a property", "", wx.OK, nil, -1, -1)
	}
}

func (f *mainFrame) onMisc(evt *wx.Event) {
	cmd := commandEvent(evt)
	if cmd == nil {
		return
	}
	switch cmd.GetId() {
	case idStaticLayout:
		if f.manager == nil {
			return
		}
		wsf := f.manager.GetWindowStyleFlag()
		if cmd.IsChecked() {
			f.manager.SetWindowStyleFlag(wsf | wx.PG_STATIC_LAYOUT)
		} else {
			f.manager.SetWindowStyleFlag(wsf &^ wx.PG_STATIC_LAYOUT)
		}
	case idCollapseAll:
		if f.manager == nil {
			return
		}
		page := f.manager.GetCurrentPage()
		if page == nil {
			return
		}
		eachProperty(wx.PropertyGridPageFromPtr(page).GetRoot(), func(p unsafe.Pointer) {
			wx.PGPropertyFromPtr(p).SetExpanded(false)
		})
		if mi := f.managerIface(); mi != nil {
			mi.RefreshGrid(nil)
		}
	case idGetValues:
		gi := f.gridIface()
		if gi == nil {
			return
		}
		if f.storedValues != nil {
			f.storedValues.Delete()
		}
		f.storedValues = wx.VariantFromPtr(gi.GetPropertyValues("Test",
			f.grid.GetRoot(), pgKeepStructure|pgIncAttributes))
	case idSetValues:
		gi := f.gridIface()
		if gi == nil {
			return
		}
		if f.storedValues != nil && f.storedValues.IsType("list") {
			gi.SetPropertyValues(f.storedValues.Ptr(), nil)
		} else {
			wx.MessageBox("First use Get Property Values.", "", wx.OK, nil, -1, -1)
		}
	case idSetValues2:
		gi := f.gridIface()
		if gi == nil {
			return
		}
		list := wx.NewVariant()
		if list == nil {
			return
		}
		list.NullList()
		v1 := wx.NewVariantFromLong(1234, "VariantLong")
		v2 := wx.NewVariantFromBool(true, "VariantBool")
		v3 := wx.NewVariantFromString("Test Text", "VariantString")
		list.Append(v1.Ptr())
		list.Append(v2.Ptr())
		list.Append(v3.Ptr())
		gi.SetPropertyValues(list.Ptr(), nil)
		v1.Delete()
		v2.Delete()
		v3.Delete()
		list.Delete()
	case idCollapse:
		if mi := f.managerIface(); mi != nil {
			if prop := mi.GetSelection(); prop != nil {
				mi.Collapse(prop)
			}
		}
	case idUnspecify:
		if mi := f.managerIface(); mi != nil {
			if prop := mi.GetSelection(); prop != nil {
				mi.SetPropertyValueUnspecified(prop)
				wx.PGPropertyFromPtr(prop).RefreshEditor()
			}
		}
	}
}

func (f *mainFrame) onFitColumnsClick(evt *wx.Event) {
	if f.manager == nil {
		return
	}
	page := f.manager.GetCurrentPage()
	if page == nil {
		return
	}
	f.manager.SetWindowStyle(f.manager.GetWindowStyle() &^ wx.PG_SPLITTER_AUTO_CENTER)

	oldGridSize := f.grid.GetClientSize()
	oldFullSize := f.frame.GetSize()
	f.frame.SetSize(-1, -1, 1000, oldFullSize.Height, wx.SIZE_USE_EXISTING)

	newSz := wx.PropertyGridPageFromPtr(page).FitColumns()

	dx := oldFullSize.Width - oldGridSize.Width
	dy := oldFullSize.Height - oldGridSize.Height

	newSz.Width += dx
	newSz.Height += dy
	f.frame.SetSize(-1, -1, newSz.Width, newSz.Height, wx.SIZE_USE_EXISTING)
}

func (f *mainFrame) onChangeFlagsPropItemsClick(evt *wx.Event) {
	mi := f.managerIface()
	if mi == nil {
		return
	}
	p := propByName(mi, "Window Styles")
	if p == nil {
		return
	}
	choices := wx.NewPGChoices()
	if choices == nil {
		return
	}
	defer choices.Delete()
	choices.Add("Fast", 0x1)
	choices.Add("Powerful", 0x2)
	choices.Add("Safe", 0x4)
	choices.Add("Sleek", 0x8)

	if !wx.PGPropertyFromPtr(p).SetChoices(choices.Ptr()) {
		wx.LogMessage("SetChoices failed")
	}
}

func (f *mainFrame) onSetColumns(evt *wx.Event) {
	if f.manager == nil {
		return
	}
	colCount := wx.GetNumberFromUser("Enter number of columns (2-20).", "Columns:",
		"Change Columns", f.manager.GetColumnCount(-1), 2, 20, f.frame.Ptr(), -1, -1)
	if colCount >= 2 {
		f.manager.SetColumnCount(colCount, -1)
	}
}

func (f *mainFrame) onSetVirtualWidth(evt *wx.Event) {
	mi := f.managerIface()
	if mi == nil || f.grid == nil {
		return
	}
	state := wx.PropertyGridPageStateFromPtr(mi.GetState())
	if state == nil {
		return
	}
	oldWidth := state.GetVirtualWidth()
	newWidth := oldWidth
	if dlg := wx.NewNumberEntryDialog(f.frame.Ptr(), "Enter virtual width (-1-2000).",
		"Width:", "Change Virtual Width", oldWidth, -1, 2000, -1, -1); dlg != nil {
		if dlg.ShowModal() == wx.ID_OK {
			newWidth = dlg.GetValue()
		}
		dlg.Delete()
	}
	if newWidth != oldWidth {
		f.grid.SetVirtualWidth(newWidth)
	}
}

func (f *mainFrame) onSetGridDisabled(evt *wx.Event) {
	cmd := commandEvent(evt)
	if cmd == nil || f.manager == nil {
		return
	}
	f.manager.Enable(!cmd.IsChecked())
}

func (f *mainFrame) onSetPropertyValue(evt *wx.Event) {
	gi := f.gridIface()
	if gi == nil {
		return
	}
	selected := gi.GetSelection()
	if selected == nil {
		return
	}
	value := wx.GetTextFromUserEx("Enter new value:", "", "", f.frame.Ptr(), -1, -1, 1)
	gi.SetPropertyValueString(selected, value)
}

func (f *mainFrame) onInsertChoice(evt *wx.Event) {
	gi := f.gridIface()
	if gi == nil {
		return
	}
	selected := gi.GetSelection()
	if selected != nil {
		prop := wx.PGPropertyFromPtr(selected)
		choices := wx.PGChoicesFromPtr(prop.GetChoices())
		if choices != nil {
			ok := choices.IsOk()
			// Insert new choice to the center of list.
			pos := int(choices.GetCount()) / 2
			choices.Delete()
			if ok {
				prop.InsertChoice("New Choice", pos, -1)
				return
			}
		}
	}
	wx.MessageBox("First select a property with some choices.", "", wx.OK, nil, -1, -1)
}

func (f *mainFrame) onDeleteChoice(evt *wx.Event) {
	gi := f.gridIface()
	if gi == nil {
		return
	}
	selected := gi.GetSelection()
	if selected != nil {
		prop := wx.PGPropertyFromPtr(selected)
		choices := wx.PGChoicesFromPtr(prop.GetChoices())
		if choices != nil {
			ok := choices.IsOk()
			// Deletes choice from the center of list.
			pos := int(choices.GetCount()) / 2
			choices.Delete()
			if ok {
				prop.DeleteChoice(pos)
				return
			}
		}
	}
	wx.MessageBox("First select a property with some choices.", "", wx.OK, nil, -1, -1)
}

func (f *mainFrame) onExtendedKeyNav(evt *wx.Event) {
	if f.grid == nil {
		return
	}
	f.grid.AddActionTrigger(pgNextProperty, wx.K_RETURN, 0)
	f.grid.DedicateKey(wx.K_RETURN)
	f.grid.DedicateKey(wx.K_UP)
	f.grid.DedicateKey(wx.K_DOWN)
}

// --- Iteration -------------------------------------------------------------

const (
	iterAll = iota
	iterVisible
	iterCategories
)

// collectProps walks the current page and returns the properties selected by
// filter, in display order.
func (f *mainFrame) collectProps(page unsafe.Pointer, filter int) []unsafe.Pointer {
	var out []unsafe.Pointer
	eachProperty(wx.PropertyGridPageFromPtr(page).GetRoot(), func(p unsafe.Pointer) {
		prop := wx.PGPropertyFromPtr(p)
		if prop == nil {
			return
		}
		if filter == iterCategories {
			if prop.IsCategory() {
				out = append(out, p)
			}
			return
		}
		if prop.IsCategory() {
			return
		}
		if filter == iterVisible {
			// Skip hidden items and the children of collapsed parents.
			if !isVisible(f.managerIface(), p) {
				return
			}
		}
		out = append(out, p)
	})
	return out
}

func isVisible(mi *wx.PropertyGridInterface, p unsafe.Pointer) bool {
	if mi == nil || !mi.IsPropertyShown(p) {
		return false
	}
	prop := wx.PGPropertyFromPtr(p)
	parent := prop.GetParent()
	for parent != nil {
		pp := wx.PGPropertyFromPtr(parent)
		if pp == nil {
			break
		}
		// The hidden root has no parent; never treat it as collapsed.
		if pp.GetParent() == nil {
			break
		}
		if !pp.IsExpanded() {
			return false
		}
		parent = pp.GetParent()
	}
	return true
}

func iterateMessage(prop unsafe.Pointer) int {
	p := wx.PGPropertyFromPtr(prop)
	if p == nil {
		return wx.CANCEL
	}
	s := fmt.Sprintf("\"%s\" valuetype = %s", p.GetLabel(), p.GetValueType())
	return wx.MessageBox(s, "Iterating... (press CANCEL to end)", wx.OK|wx.CANCEL, nil, -1, -1)
}

func (f *mainFrame) runIteration(filter int, reverse bool) {
	if f.manager == nil {
		return
	}
	page := f.manager.GetCurrentPage()
	if page == nil {
		return
	}
	props := f.collectProps(page, filter)
	if reverse {
		for i, j := 0, len(props)-1; i < j; i, j = i+1, j-1 {
			props[i], props[j] = props[j], props[i]
		}
	}
	for _, p := range props {
		if iterateMessage(p) == wx.CANCEL {
			break
		}
	}
}

func (f *mainFrame) onIterate1Click(evt *wx.Event) { f.runIteration(iterAll, false) }
func (f *mainFrame) onIterate2Click(evt *wx.Event) { f.runIteration(iterVisible, false) }
func (f *mainFrame) onIterate3Click(evt *wx.Event) { f.runIteration(iterAll, true) }
func (f *mainFrame) onIterate4Click(evt *wx.Event) { f.runIteration(iterCategories, false) }

// --- Colour schemes and categories -----------------------------------------

func (f *mainFrame) onColourScheme(evt *wx.Event) {
	cmd := commandEvent(evt)
	if cmd == nil || f.manager == nil {
		return
	}
	grid := f.grid
	switch cmd.GetId() {
	case idColourScheme1:
		grid.ResetColours()
	case idColourScheme2:
		// white
		myGrey1 := wx.NewColourRGB(212, 208, 200, 255)
		myGrey3 := wx.NewColourRGB(113, 111, 100, 255)
		f.manager.Freeze()
		grid.SetMarginColour(wx.WHITE)
		grid.SetCaptionBackgroundColour(wx.WHITE)
		grid.SetCellBackgroundColour(wx.WHITE)
		grid.SetCellTextColour(myGrey3.Ptr())
		grid.SetLineColour(myGrey1.Ptr())
		f.manager.Thaw()
		myGrey1.Delete()
		myGrey3.Delete()
	case idColourScheme3:
		// .NET
		myGrey1 := wx.NewColourRGB(212, 208, 200, 255)
		myGrey2 := wx.NewColourRGB(236, 233, 216, 255)
		f.manager.Freeze()
		grid.SetMarginColour(myGrey1.Ptr())
		grid.SetCaptionBackgroundColour(myGrey1.Ptr())
		grid.SetLineColour(myGrey2.Ptr())
		f.manager.Thaw()
		myGrey1.Delete()
		myGrey2.Delete()
	case idColourScheme4:
		// cream
		myGrey1 := wx.NewColourRGB(212, 208, 200, 255)
		myGrey2 := wx.NewColourRGB(241, 239, 226, 255)
		myGrey3 := wx.NewColourRGB(113, 111, 100, 255)
		f.manager.Freeze()
		grid.SetMarginColour(wx.WHITE)
		grid.SetCaptionBackgroundColour(wx.WHITE)
		grid.SetCellBackgroundColour(myGrey2.Ptr())
		grid.SetCellTextColour(myGrey3.Ptr())
		grid.SetLineColour(myGrey1.Ptr())
		f.manager.Thaw()
		myGrey1.Delete()
		myGrey2.Delete()
		myGrey3.Delete()
	}
}

func (f *mainFrame) onCatColoursUpdateUI(evt *wx.Event) {
	if f.itemCatColours == nil {
		return
	}
	gi := f.gridIface()
	if gi == nil {
		return
	}
	f.itemCatColours.SetCheckable(
		propByName(gi, "Appearance") != nil &&
			propByName(gi, "PositionCategory") != nil &&
			propByName(gi, "Environment") != nil &&
			propByName(gi, "More Examples") != nil)
}

func (f *mainFrame) onCatColours(evt *wx.Event) {
	cmd := commandEvent(evt)
	gi := f.gridIface()
	if cmd == nil || gi == nil || f.manager == nil {
		return
	}
	if propByName(gi, "Appearance") == nil ||
		propByName(gi, "PositionCategory") == nil ||
		propByName(gi, "Environment") == nil ||
		propByName(gi, "More Examples") == nil {
		wx.MessageBox("First switch to 'Standard Items' page!", "", wx.OK, nil, -1, -1)
		return
	}

	f.manager.Freeze()

	if cmd.IsChecked() {
		paint := func(name string, r, g, b uint8, flags int) {
			if p := propByName(gi, name); p != nil {
				col := wx.NewColourRGB(r, g, b, 255)
				gi.SetPropertyTextColour(p, col.Ptr(), flags)
				col.Delete()
			}
		}
		bgpaint := func(name string, r, g, b uint8) {
			if p := propByName(gi, name); p != nil {
				col := wx.NewColourRGB(r, g, b, 255)
				gi.SetPropertyBackgroundColour(p, col.Ptr(), pgDontRecurse)
				col.Delete()
			}
		}
		paint("Appearance", 255, 0, 0, pgDontRecurse)
		bgpaint("Appearance", 255, 255, 183)
		paint("Appearance", 255, 0, 183, pgDontRecurse)
		paint("PositionCategory", 0, 255, 0, pgDontRecurse)
		bgpaint("PositionCategory", 255, 226, 190)
		paint("PositionCategory", 255, 0, 190, pgDontRecurse)
		paint("Environment", 0, 0, 255, pgDontRecurse)
		bgpaint("Environment", 208, 240, 175)
		paint("Environment", 255, 255, 255, pgDontRecurse)
		bgpaint("More Examples", 172, 237, 255)
		paint("More Examples", 172, 0, 255, pgDontRecurse)
	} else {
		defaults := func(name string, flags int) {
			if p := propByName(gi, name); p != nil {
				gi.SetPropertyColoursToDefault(p, flags)
			}
		}
		defaults("Appearance", pgDontRecurse)
		defaults("Appearance", pgRecurse)
		defaults("PositionCategory", pgDontRecurse)
		defaults("PositionCategory", pgRecurse)
		defaults("Environment", pgDontRecurse)
		defaults("Environment", pgRecurse)
		defaults("More Examples", pgRecurse)
	}
	f.manager.Thaw()
	f.manager.Refresh(true)
}

// onSelectStyle is FormMain::OnSelectStyle().
func (f *mainFrame) onSelectStyle(evt *wx.Event) {
	if f.manager == nil {
		return
	}

	styleLabels := []string{
		"wxPG_HIDE_CATEGORIES", "wxPG_AUTO_SORT", "wxPG_BOLD_MODIFIED",
		"wxPG_SPLITTER_AUTO_CENTER", "wxPG_TOOLTIPS", "wxPG_STATIC_SPLITTER",
		"wxPG_HIDE_MARGIN", "wxPG_LIMITED_EDITING", "wxPG_TOOLBAR",
		"wxPG_DESCRIPTION", "wxPG_NO_INTERNAL_BORDER",
	}
	styleValues := []int{
		wx.PG_HIDE_CATEGORIES, wx.PG_AUTO_SORT, wx.PG_BOLD_MODIFIED,
		wx.PG_SPLITTER_AUTO_CENTER, wx.PG_TOOLTIPS, wx.PG_STATIC_SPLITTER,
		wx.PG_HIDE_MARGIN, wx.PG_LIMITED_EDITING, wx.PG_TOOLBAR,
		wx.PG_DESCRIPTION, wx.PG_NO_INTERNAL_BORDER,
	}
	style, ok := chooseFlags(f.frame.Ptr(), "Select window styles to use",
		"wxPropertyGrid Window Style", styleLabels, styleValues,
		f.manager.GetWindowStyle())
	if !ok {
		return
	}

	exLabels := []string{
		"wxPG_EX_INIT_NOCAT", "wxPG_EX_NO_FLAT_TOOLBAR", "wxPG_EX_MODE_BUTTONS",
		"wxPG_EX_HELP_AS_TOOLTIPS", "wxPG_EX_NATIVE_DOUBLE_BUFFERING",
		"wxPG_EX_AUTO_UNSPECIFIED_VALUES", "wxPG_EX_WRITEONLY_BUILTIN_ATTRIBUTES",
		"wxPG_EX_HIDE_PAGE_BUTTONS", "wxPG_EX_MULTIPLE_SELECTION",
		"wxPG_EX_ENABLE_TLP_TRACKING", "wxPG_EX_NO_TOOLBAR_DIVIDER",
		"wxPG_EX_TOOLBAR_SEPARATOR", "wxPG_EX_ALWAYS_ALLOW_FOCUS",
	}
	exValues := []int{
		wx.PG_EX_INIT_NOCAT, wx.PG_EX_NO_FLAT_TOOLBAR, wx.PG_EX_MODE_BUTTONS,
		wx.PG_EX_HELP_AS_TOOLTIPS, wx.PG_EX_NATIVE_DOUBLE_BUFFERING,
		wx.PG_EX_AUTO_UNSPECIFIED_VALUES, wx.PG_EX_WRITEONLY_BUILTIN_ATTRIBUTES,
		wx.PG_EX_HIDE_PAGE_BUTTONS, wx.PG_EX_MULTIPLE_SELECTION,
		wx.PG_EX_ENABLE_TLP_TRACKING, wx.PG_EX_NO_TOOLBAR_DIVIDER,
		wx.PG_EX_TOOLBAR_SEPARATOR, wx.PG_EX_ALWAYS_ALLOW_FOCUS,
	}
	extraStyle, ok := chooseFlags(f.frame.Ptr(), "Select extra window styles to use",
		"wxPropertyGrid Extra Style", exLabels, exValues, f.manager.GetExtraStyle())
	if !ok {
		return
	}

	f.replaceGrid(style, extraStyle)
}

// chooseFlags shows a wxMultiChoiceDialog and returns the OR of the values whose
// entries were selected.
func chooseFlags(parent unsafe.Pointer, message, caption string, labels []string,
	values []int, current int) (int, bool) {
	sel := wx.NewArrayInt()
	if sel == nil {
		return 0, false
	}
	defer sel.Delete()
	for i, v := range values {
		if current&v == v {
			sel.Add(i)
		}
	}
	// The wxMultiChoiceDialog default style is
	// wxDEFAULT_DIALOG_STYLE|wxRESIZE_BORDER|wxOK|wxCANCEL|wxCENTRE; the OK and
	// CANCEL bits are what make wxAnyChoiceDialog create its button sizer.
	dlg := wx.NewMultiChoiceDialog(parent, message, caption, labels,
		wx.DEFAULT_DIALOG_STYLE|wx.RESIZE_BORDER|wx.OK|wx.CANCEL|wx.CENTRE, -1, -1)
	if dlg == nil {
		return 0, false
	}
	defer dlg.Delete()
	dlg.SetSelections(sel.Ptr())
	if dlg.ShowModal() == wx.ID_CANCEL {
		return 0, false
	}
	dlg.GetSelections(sel.Ptr())
	result := 0
	for ind := 0; ind < sel.GetCount(); ind++ {
		idx := sel.Item(ind)
		if idx >= 0 && idx < len(values) {
			result |= values[idx]
		}
	}
	return result, true
}

func uniqueLabel(mi *wx.PropertyGridInterface, base string) string {
	if mi.GetPropertyByLabel(base) == nil {
		return base
	}
	for i := 0; ; i++ {
		cand := fmt.Sprintf("%s%d", base, i)
		if mi.GetPropertyByLabel(cand) == nil {
			return cand
		}
	}
}
