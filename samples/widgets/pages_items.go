// Pages demonstrating the item-container controls: choice, listbox, combobox
// and spin controls (choice.cpp, listbox.cpp, combobox.cpp, spinbtn.cpp).
package main

import (
	"strconv"
	"unsafe"

	"github.com/KeyWorksRW/kwxGO/wx"
)

// --- Choice page (choice.cpp) ----------------------------------------------

var (
	idChoiceReset     = wx.ID_HIGHEST
	idChoiceAdd       = wx.ID_HIGHEST + 1
	idChoiceAddText   = wx.ID_HIGHEST + 2
	idChoiceAddMany   = wx.ID_HIGHEST + 3
	idChoiceClear     = wx.ID_HIGHEST + 4
	idChoiceChange    = wx.ID_HIGHEST + 5
	idChoiceChangeTxt = wx.ID_HIGHEST + 6
	idChoiceDelete    = wx.ID_HIGHEST + 7
	idChoiceDeleteTxt = wx.ID_HIGHEST + 8
	idChoiceDeleteSel = wx.ID_HIGHEST + 9
	idChoiceCtrl      = wx.ID_HIGHEST + 10
)

type choicePage struct {
	pageBase

	chkSort   *wx.CheckBox
	textAdd   *wx.TextCtrl
	textChg   *wx.TextCtrl
	textDel   *wx.TextCtrl
	choice    *wx.Choice
	sizerCh   *wx.BoxSizer
	itemCount int
}

func init() {
	registerPage("Choice", nativeCtrls|itemsCtrls, func(book unsafe.Pointer) *pageBase {
		pg := &choicePage{}
		pb := newPageWithIcon(book, "icons/choice.svg", 16, 16)
		if pb == nil {
			return nil
		}
		pg.pageBase = *pb
		pg.content = pg.create
		pg.getWidgets = func() []unsafe.Pointer {
			if pg.choice == nil {
				return nil
			}
			return []unsafe.Pointer{pg.choice.Ptr()}
		}
		pg.recreate = func(f *widgetsFrame) { pg.createChoice(f) }
		return &pg.pageBase
	})
}

func (pg *choicePage) create(f *widgetsFrame) {
	top := wx.NewBoxSizer(wx.HORIZONTAL)
	if top == nil {
		return
	}

	boxL, left := addStaticBox(pg.win.Ptr(), "&Set choice parameters", wx.VERTICAL)
	if left != nil {
		pg.chkSort = pg.createCheckBoxAndAddToSizer(left, "&Sort items", wx.ID_ANY, boxL.Ptr())
		reset := wx.NewButton(boxL.Ptr(), idChoiceReset, "&Reset", -1, -1, -1, -1, 0)
		addWin(left, reset.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.CenterHorizontal()
			sf.TripleBorder(wx.ALL)
		})
		reset.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			pg.reset()
			pg.createChoice(f)
		})
	}

	boxM, middle := addStaticBox(pg.win.Ptr(), "&Change choice contents", wx.VERTICAL)
	if middle != nil {
		addRow(middle, boxM.Ptr(), idChoiceAdd, "&Add this string", idChoiceAddText, &pg.textAdd, "test item 0")
		btnSeveral := wx.NewButton(boxM.Ptr(), idChoiceAddMany, "&Insert a few strings", -1, -1, -1, -1, 0)
		addExpandBorder(middle, btnSeveral.Ptr())
		btnMany := wx.NewButton(boxM.Ptr(), idChoiceAddMany+100, "Add &many strings", -1, -1, -1, -1, 0)
		addExpandBorder(middle, btnMany.Ptr())
		addRow(middle, boxM.Ptr(), idChoiceChange, "C&hange current", idChoiceChangeTxt, &pg.textChg, "")
		addRow(middle, boxM.Ptr(), idChoiceDelete, "&Delete this item", idChoiceDeleteTxt, &pg.textDel, "")
		btnDelSel := wx.NewButton(boxM.Ptr(), idChoiceDeleteSel, "Delete &selection", -1, -1, -1, -1, 0)
		addExpandBorder(middle, btnDelSel.Ptr())
		btnClear := wx.NewButton(boxM.Ptr(), idChoiceClear, "&Clear", -1, -1, -1, -1, 0)
		addExpandBorder(middle, btnClear.Ptr())

		wx.WindowFromPtr(btnSeveral.Ptr()).Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.choice == nil {
				return
			}
			ic := wx.ItemContainerFromPtr(pg.choice.Ptr())
			if ic == nil {
				return
			}
			ic.Insert("First", 0)
			ic.Insert("another one", 1)
			ic.Insert("and the last (very very very very very very very very very very long) one", 2)
		})
		wx.WindowFromPtr(btnMany.Ptr()).Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.choice == nil {
				return
			}
			for n := 0; n < 1000; n++ {
				pg.choice.Append("item #" + strconv.Itoa(n))
			}
		})
		wx.WindowFromPtr(btnDelSel.Ptr()).Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.choice == nil {
				return
			}
			if sel := pg.choice.GetSelection(); sel != wx.NOT_FOUND {
				pg.choice.Delete(sel)
			}
		})
		wx.WindowFromPtr(btnClear.Ptr()).Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.choice != nil {
				pg.choice.Clear()
			}
		})

		pg.win.BindId(wx.EVT_BUTTON, idChoiceChange, func(evt *wx.Event) {
			if pg.choice == nil || pg.textChg == nil {
				return
			}
			if sel := pg.choice.GetSelection(); sel != wx.NOT_FOUND {
				pg.choice.SetString(sel, pg.textChg.GetValue())
			}
		})
		pg.win.BindId(wx.EVT_BUTTON, idChoiceDelete, func(evt *wx.Event) {
			if pg.choice == nil || pg.textDel == nil {
				return
			}
			if n, err := strconv.Atoi(pg.textDel.GetValue()); err == nil && n < pg.choice.GetCount() {
				pg.choice.Delete(n)
			}
		})
		pg.win.BindId(wx.EVT_BUTTON, idChoiceAdd, func(evt *wx.Event) { pg.onAdd() })
		pg.win.BindId(wx.EVT_TEXT_ENTER, idChoiceAddText, func(evt *wx.Event) { pg.onAdd() })
		pg.win.BindId(wx.EVT_TEXT_ENTER, idChoiceDeleteTxt, func(evt *wx.Event) {
			if pg.choice == nil || pg.textDel == nil {
				return
			}
			if n, err := strconv.Atoi(pg.textDel.GetValue()); err == nil && n < pg.choice.GetCount() {
				pg.choice.Delete(n)
			}
		})
	}

	pg.sizerCh = wx.NewBoxSizer(wx.VERTICAL)
	if pg.sizerCh == nil {
		return
	}
	pg.sizerCh.SetMinSize(150, 0)
	pg.choice = wx.NewChoice(pg.win.Ptr(), idChoiceCtrl, -1, -1, -1, -1, nil, 0)
	addWin(pg.sizerCh, pg.choice.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})
	if pg.choice != nil {
		pg.choice.Bind(wx.EVT_CHOICE, func(evt *wx.Event) {
			cmd := wx.AsCommandEvent(evt)
			if cmd == nil || pg.choice == nil {
				return
			}
			if pg.textDel != nil {
				pg.textDel.SetValue(strconv.Itoa(cmd.GetSelection()))
			}
			f.logf("Choice item %d selected", cmd.GetSelection())
		})
	}

	addSub(top, left.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.LEFT)
	})
	addSub(top, middle.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL)
	})
	addSub(top, pg.sizerCh.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.RIGHT)
	})

	pg.reset()
	pg.win.SetSizer(top.Ptr(), false)
	bindStyleCheckboxes(pg.win, func() { pg.createChoice(f) })
}

func (pg *choicePage) reset() {
	if pg.chkSort != nil {
		pg.chkSort.SetValue(0)
	}
}

func (pg *choicePage) onAdd() {
	if pg.choice == nil || pg.textAdd == nil {
		return
	}
	pg.choice.Append(pg.textAdd.GetValue())
	pg.itemCount++
	pg.textAdd.SetValue("test item " + strconv.Itoa(pg.itemCount))
}

// createChoice is ChoiceWidgetsPage::CreateChoice().
func (pg *choicePage) createChoice(f *widgetsFrame) {
	if pg.sizerCh == nil {
		return
	}

	flags := f.attrs.defaultFlags
	if pg.chkSort != nil && pg.chkSort.GetValue() {
		flags |= wx.CB_SORT
	}

	var items []string
	sel := -1
	if pg.choice != nil {
		for n := 0; n < pg.choice.GetCount(); n++ {
			items = append(items, pg.choice.GetString(n))
		}
		sel = pg.choice.GetSelection()
		pg.choice.Destroy()
		pg.choice = nil
	}
	clearSizer(pg.sizerCh)

	pg.choice = wx.NewChoice(pg.win.Ptr(), idChoiceCtrl, -1, -1, -1, -1, nil, flags)
	if pg.choice == nil {
		return
	}
	for _, s := range items {
		pg.choice.Append(s)
	}
	if sel >= 0 && sel < len(items) {
		pg.choice.SetSelection(sel)
	}
	pg.choice.Bind(wx.EVT_CHOICE, func(evt *wx.Event) {
		cmd := wx.AsCommandEvent(evt)
		if cmd == nil || pg.choice == nil {
			return
		}
		if pg.textDel != nil {
			pg.textDel.SetValue(strconv.Itoa(cmd.GetSelection()))
		}
		f.logf("Choice item %d selected", cmd.GetSelection())
	})

	addWin(pg.sizerCh, pg.choice.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})
	pg.sizerCh.Layout()
	pg.notifyWidgetRecreation(f, pg.choice.Ptr())
}

// --- Listbox page (listbox.cpp) --------------------------------------------

var (
	idLbReset     = wx.ID_HIGHEST
	idLbAdd       = wx.ID_HIGHEST + 1
	idLbAddText   = wx.ID_HIGHEST + 2
	idLbAddMany   = wx.ID_HIGHEST + 3
	idLbClear     = wx.ID_HIGHEST + 4
	idLbChange    = wx.ID_HIGHEST + 5
	idLbChangeTxt = wx.ID_HIGHEST + 6
	idLbDelete    = wx.ID_HIGHEST + 7
	idLbDeleteTxt = wx.ID_HIGHEST + 8
	idLbDeleteSel = wx.ID_HIGHEST + 9
	idLbCtrl      = wx.ID_HIGHEST + 10
)

type listboxPage struct {
	pageBase

	chkVScroll *wx.CheckBox
	chkHScroll *wx.CheckBox
	chkSort    *wx.CheckBox
	chkOwner   *wx.CheckBox
	radioMode  *wx.RadioBox
	radioType  *wx.RadioBox
	textAdd    *wx.TextCtrl
	textChg    *wx.TextCtrl
	textDel    *wx.TextCtrl
	lbox       *wx.ListBox
	sizerLb    *wx.BoxSizer
	itemCount  int
}

func init() {
	registerPage("Listbox", nativeCtrls|itemsCtrls, func(book unsafe.Pointer) *pageBase {
		pg := &listboxPage{}
		pb := newPageWithIcon(book, "icons/listbox.svg", 32, 32)
		if pb == nil {
			return nil
		}
		pg.pageBase = *pb
		pg.content = pg.create
		pg.getWidgets = func() []unsafe.Pointer {
			if pg.lbox == nil {
				return nil
			}
			return []unsafe.Pointer{pg.lbox.Ptr()}
		}
		pg.recreate = func(f *widgetsFrame) { pg.createLbox(f) }
		return &pg.pageBase
	})
}

func (pg *listboxPage) create(f *widgetsFrame) {
	top := wx.NewBoxSizer(wx.HORIZONTAL)
	if top == nil {
		return
	}

	boxL, left := addStaticBox(pg.win.Ptr(), "&Set listbox parameters", wx.VERTICAL)
	if left != nil {
		pg.chkVScroll = pg.createCheckBoxAndAddToSizer(left, "Always show &vertical scrollbar", wx.ID_ANY, boxL.Ptr())
		pg.chkHScroll = pg.createCheckBoxAndAddToSizer(left, "Show &horizontal scrollbar", wx.ID_ANY, boxL.Ptr())
		pg.chkSort = pg.createCheckBoxAndAddToSizer(left, "&Sort items", wx.ID_ANY, boxL.Ptr())
		pg.chkOwner = pg.createCheckBoxAndAddToSizer(left, "&Owner drawn", wx.ID_ANY, boxL.Ptr())
		pg.chkHScroll.SetValue(1)

		addSpacerSizer(left, 5)
		pg.radioMode = newRadioBox(boxL.Ptr(), wx.ID_ANY, "Selection &mode:",
			[]string{"single", "extended", "multiple"}, 1, wx.RA_SPECIFY_COLS)
		addExpandBorder(left, pg.radioMode.Ptr())
		addSpacerSizer(left, 5)
		pg.radioType = newRadioBox(boxL.Ptr(), wx.ID_ANY, "&List type:",
			[]string{"list box", "check list box", "rearrange list"}, 1, wx.RA_SPECIFY_COLS)
		addExpandBorder(left, pg.radioType.Ptr())

		reset := wx.NewButton(boxL.Ptr(), idLbReset, "&Reset", -1, -1, -1, -1, 0)
		addWin(left, reset.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.CenterHorizontal()
			sf.TripleBorder(wx.ALL)
		})
		reset.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			pg.reset()
			pg.createLbox(f)
		})
	}

	boxM, middle := addStaticBox(pg.win.Ptr(), "&Change listbox contents", wx.VERTICAL)
	if middle != nil {
		addRow(middle, boxM.Ptr(), idLbAdd, "&Add this string", idLbAddText, &pg.textAdd, "test item \t0")
		btnSeveral := wx.NewButton(boxM.Ptr(), idLbAddMany, "&Insert a few strings", -1, -1, -1, -1, 0)
		addExpandBorder(middle, btnSeveral.Ptr())
		btnMany := wx.NewButton(boxM.Ptr(), idLbAddMany+100, "Add &many strings", -1, -1, -1, -1, 0)
		addExpandBorder(middle, btnMany.Ptr())
		addRow(middle, boxM.Ptr(), idLbChange, "C&hange current", idLbChangeTxt, &pg.textChg, "")
		addRow(middle, boxM.Ptr(), idLbDelete, "&Delete this item", idLbDeleteTxt, &pg.textDel, "")
		btnDelSel := wx.NewButton(boxM.Ptr(), idLbDeleteSel, "Delete &selection", -1, -1, -1, -1, 0)
		addExpandBorder(middle, btnDelSel.Ptr())
		btnClear := wx.NewButton(boxM.Ptr(), idLbClear, "&Clear", -1, -1, -1, -1, 0)
		addExpandBorder(middle, btnClear.Ptr())

		wx.WindowFromPtr(btnSeveral.Ptr()).Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.lbox == nil {
				return
			}
			ic := wx.ItemContainerFromPtr(pg.lbox.Ptr())
			if ic == nil {
				return
			}
			ic.Insert("First", 0)
			ic.Insert("another one", 1)
			ic.Insert("and the last (very very very very very very very very very very long) one", 2)
		})
		wx.WindowFromPtr(btnMany.Ptr()).Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.lbox == nil {
				return
			}
			for n := 0; n < 1000; n++ {
				pg.lbox.Append("item #" + strconv.Itoa(n))
			}
		})
		wx.WindowFromPtr(btnDelSel.Ptr()).Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			pg.deleteSelection()
		})
		wx.WindowFromPtr(btnClear.Ptr()).Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.lbox != nil {
				pg.lbox.Clear()
			}
		})

		pg.win.BindId(wx.EVT_BUTTON, idLbAdd, func(evt *wx.Event) { pg.onAdd() })
		pg.win.BindId(wx.EVT_TEXT_ENTER, idLbAddText, func(evt *wx.Event) { pg.onAdd() })
		pg.win.BindId(wx.EVT_BUTTON, idLbDelete, func(evt *wx.Event) { pg.onDelete() })
		pg.win.BindId(wx.EVT_TEXT_ENTER, idLbDeleteTxt, func(evt *wx.Event) { pg.onDelete() })
		pg.win.BindId(wx.EVT_BUTTON, idLbChange, func(evt *wx.Event) {
			if pg.lbox == nil || pg.textChg == nil {
				return
			}
			if sel := pg.lbox.GetSelection(); sel != wx.NOT_FOUND {
				wx.ItemContainerFromPtr(pg.lbox.Ptr()).SetString(uint(sel), pg.textChg.GetValue())
			}
		})
	}

	pg.sizerLb = wx.NewBoxSizer(wx.VERTICAL)
	if pg.sizerLb == nil {
		return
	}
	pg.sizerLb.SetMinSize(150, 0)
	pg.lbox = wx.NewListBox(pg.win.Ptr(), idLbCtrl, -1, -1, -1, -1, nil, 0)
	addWin(pg.sizerLb, pg.lbox.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})
	pg.bindListboxEvents(f)

	addSub(top, left.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.LEFT)
	})
	addSub(top, middle.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL)
	})
	addSub(top, pg.sizerLb.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.RIGHT)
	})

	pg.reset()
	pg.win.SetSizer(top.Ptr(), false)
	bindStyleCheckboxes(pg.win, func() { pg.createLbox(f) })
}

func (pg *listboxPage) reset() {
	if pg.radioMode != nil {
		pg.radioMode.SetSelection(0)
	}
	if pg.radioType != nil {
		pg.radioType.SetSelection(0)
	}
	if pg.chkVScroll != nil {
		pg.chkVScroll.SetValue(0)
	}
	if pg.chkHScroll != nil {
		pg.chkHScroll.SetValue(1)
	}
	if pg.chkSort != nil {
		pg.chkSort.SetValue(0)
	}
	if pg.chkOwner != nil {
		pg.chkOwner.SetValue(0)
	}
}

func (pg *listboxPage) onAdd() {
	if pg.lbox == nil || pg.textAdd == nil {
		return
	}
	pg.lbox.Append(pg.textAdd.GetValue())
	pg.itemCount++
	pg.textAdd.SetValue("test item \t" + strconv.Itoa(pg.itemCount))
}

func (pg *listboxPage) onDelete() {
	if pg.lbox == nil || pg.textDel == nil {
		return
	}
	if n, err := strconv.Atoi(pg.textDel.GetValue()); err == nil && n < pg.lbox.GetCount() {
		wx.ItemContainerFromPtr(pg.lbox.Ptr()).Delete(uint(n))
	}
}

func (pg *listboxPage) deleteSelection() {
	if pg.lbox == nil {
		return
	}
	ic := wx.ItemContainerFromPtr(pg.lbox.Ptr())
	for sel := pg.lbox.GetCount() - 1; sel >= 0; sel-- {
		if pg.lbox.GetSelection() == sel {
			ic.Delete(uint(sel))
		}
	}
}

func (pg *listboxPage) bindListboxEvents(f *widgetsFrame) {
	if pg.lbox == nil {
		return
	}
	pg.lbox.Bind(wx.EVT_LISTBOX, func(evt *wx.Event) {
		cmd := wx.AsCommandEvent(evt)
		if cmd == nil {
			return
		}
		if pg.textDel != nil {
			pg.textDel.SetValue(strconv.Itoa(cmd.GetSelection()))
		}
		f.logf("Listbox item %d selected", cmd.GetSelection())
	})
	pg.lbox.Bind(wx.EVT_LISTBOX_DCLICK, func(evt *wx.Event) {
		cmd := wx.AsCommandEvent(evt)
		if cmd == nil {
			return
		}
		f.logf("Listbox item %d double clicked", cmd.GetInt())
	})
}

// createLbox is ListboxWidgetsPage::CreateLbox().
func (pg *listboxPage) createLbox(f *widgetsFrame) {
	if pg.sizerLb == nil {
		return
	}

	flags := f.attrs.defaultFlags
	if pg.radioMode != nil {
		switch pg.radioMode.GetSelection() {
		case 1:
			flags |= wx.LB_EXTENDED
		case 2:
			flags |= wx.LB_MULTIPLE
		default:
			flags |= wx.LB_SINGLE
		}
	}
	if pg.chkVScroll != nil && pg.chkVScroll.GetValue() {
		flags |= wx.LB_ALWAYS_SB
	}
	if pg.chkSort != nil && pg.chkSort.GetValue() {
		flags |= wx.LB_SORT
	}
	if pg.chkOwner != nil && pg.chkOwner.GetValue() {
		flags |= wx.LB_OWNERDRAW
	}

	var items []string
	if pg.lbox != nil {
		for n := 0; n < pg.lbox.GetCount(); n++ {
			if s := listboxString(pg.lbox, n); s != "" {
				items = append(items, s)
			}
		}
		pg.lbox.Destroy()
		pg.lbox = nil
	}
	clearSizer(pg.sizerLb)

	pg.lbox = wx.NewListBox(pg.win.Ptr(), idLbCtrl, -1, -1, -1, -1, items, flags)
	if pg.lbox == nil {
		return
	}
	pg.bindListboxEvents(f)

	addWin(pg.sizerLb, pg.lbox.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})
	pg.sizerLb.Layout()
	pg.notifyWidgetRecreation(f, pg.lbox.Ptr())
}

func listboxString(lb *wx.ListBox, n int) string {
	ic := wx.ItemContainerFromPtr(lb.Ptr())
	if ic == nil {
		return ""
	}
	return ic.GetString(uint(n))
}

// --- Combobox page (combobox.cpp) ------------------------------------------

var (
	idComboReset     = wx.ID_HIGHEST
	idComboPopup     = wx.ID_HIGHEST + 1
	idComboDismiss   = wx.ID_HIGHEST + 2
	idComboSetCur    = wx.ID_HIGHEST + 3
	idComboCurText   = wx.ID_HIGHEST + 4
	idComboInsPtText = wx.ID_HIGHEST + 5
	idComboInsert    = wx.ID_HIGHEST + 6
	idComboInsertTxt = wx.ID_HIGHEST + 7
	idComboAdd       = wx.ID_HIGHEST + 8
	idComboAddText   = wx.ID_HIGHEST + 9
	idComboSetFirst  = wx.ID_HIGHEST + 10
	idComboSetFirstT = wx.ID_HIGHEST + 11
	idComboMany      = wx.ID_HIGHEST + 12
	idComboChange    = wx.ID_HIGHEST + 13
	idComboChangeTxt = wx.ID_HIGHEST + 14
	idComboDelete    = wx.ID_HIGHEST + 15
	idComboDeleteTxt = wx.ID_HIGHEST + 16
	idComboDeleteSel = wx.ID_HIGHEST + 17
	idComboClear     = wx.ID_HIGHEST + 18
	idComboSetValue  = wx.ID_HIGHEST + 19
	idComboValueText = wx.ID_HIGHEST + 20
	idComboCtrl      = wx.ID_HIGHEST + 21
	idComboDynamic   = wx.ID_HIGHEST + 22
)

type comboboxPage struct {
	pageBase

	radioKind   *wx.RadioBox
	chkSort     *wx.CheckBox
	chkReadonly *wx.CheckBox
	chkEnter    *wx.CheckBox
	textCur     *wx.TextCtrl
	textInsPt   *wx.TextCtrl
	textInsert  *wx.TextCtrl
	textAdd     *wx.TextCtrl
	textFirst   *wx.TextCtrl
	textChange  *wx.TextCtrl
	textDelete  *wx.TextCtrl
	textValue   *wx.TextCtrl
	combo       *wx.ComboBox
	combo1      *wx.ComboBox
	sizerCb     *wx.BoxSizer
	itemCount   int
}

func init() {
	registerPage("Combobox", nativeCtrls|itemsCtrls|comboCtrls, func(book unsafe.Pointer) *pageBase {
		pg := &comboboxPage{}
		pb := newPageWithIcon(book, "icons/combobox.svg", 32, 32)
		if pb == nil {
			return nil
		}
		pg.pageBase = *pb
		pg.content = pg.create
		pg.getWidgets = func() []unsafe.Pointer {
			if pg.combo == nil {
				return nil
			}
			return []unsafe.Pointer{pg.combo.Ptr()}
		}
		pg.getEntry = func() unsafe.Pointer {
			if pg.combo == nil {
				return nil
			}
			return pg.combo.Ptr()
		}
		pg.recreate = func(f *widgetsFrame) { pg.createCombo(f) }
		return &pg.pageBase
	})
}

func (pg *comboboxPage) create(f *widgetsFrame) {
	top := wx.NewBoxSizer(wx.HORIZONTAL)
	if top == nil {
		return
	}

	boxLT, leftTop := addStaticBox(pg.win.Ptr(), "&Popup", wx.VERTICAL)
	if leftTop != nil {
		show := wx.NewButton(boxLT.Ptr(), idComboPopup, "&Show", -1, -1, -1, -1, 0)
		addWin(leftTop, show.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Border(wx.ALL)
			sf.Center()
		})
		hide := wx.NewButton(boxLT.Ptr(), idComboDismiss, "&Hide", -1, -1, -1, -1, 0)
		addWin(leftTop, hide.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Border(wx.ALL)
			sf.Center()
		})
		show.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.combo != nil {
				pg.combo.Popup()
			}
		})
		hide.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.combo != nil {
				pg.combo.Dismiss()
			}
		})
	}

	boxLB, leftBottom := addStaticBox(pg.win.Ptr(), "&Set style", wx.VERTICAL)
	if leftBottom != nil {
		pg.chkSort = pg.createCheckBoxAndAddToSizer(leftBottom, "&Sort items", wx.ID_ANY, boxLB.Ptr())
		pg.chkReadonly = pg.createCheckBoxAndAddToSizer(leftBottom, "&Read only", wx.ID_ANY, boxLB.Ptr())
		pg.chkEnter = pg.createCheckBoxAndAddToSizer(leftBottom, "Process &Enter", wx.ID_ANY, boxLB.Ptr())
		addSpacerSizer(leftBottom, 15)

		pg.radioKind = newRadioBox(boxLB.Ptr(), wx.ID_ANY, "Combobox &kind:",
			[]string{"default", "simple", "drop down"}, 1, wx.RA_SPECIFY_COLS)
		addExpandBorder(leftBottom, pg.radioKind.Ptr())

		reset := wx.NewButton(boxLB.Ptr(), idComboReset, "&Reset", -1, -1, -1, -1, 0)
		addWin(leftBottom, reset.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.CenterHorizontal()
			sf.TripleBorder(wx.ALL)
		})
		reset.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			pg.reset()
			pg.createCombo(f)
		})
	}

	left := wx.NewBoxSizer(wx.VERTICAL)
	if left != nil {
		addSizer(left, leftTop.Ptr(), 0, nil)
		addSpacerSizer(left, 10)
		addSizer(left, leftBottom.Ptr(), 0, func(sf *wx.SizerFlags) { sf.Expand() })
	}

	boxM, middle := addStaticBox(pg.win.Ptr(), "&Change combobox contents", wx.VERTICAL)
	if middle != nil {
		addRow(middle, boxM.Ptr(), idComboSetCur, "Current &selection", idComboCurText, &pg.textCur, "")
		addRow(middle, boxM.Ptr(), idComboInsPtText+1000, "Insertion Point", idComboInsPtText, &pg.textInsPt, "")
		if pg.textInsPt != nil {
			pg.textInsPt.SetEditable(false)
		}
		addRow(middle, boxM.Ptr(), idComboInsert, "&Insert this string", idComboInsertTxt, &pg.textInsert, "")
		addRow(middle, boxM.Ptr(), idComboAdd, "&Add this string", idComboAddText, &pg.textAdd, "")
		addRow(middle, boxM.Ptr(), idComboSetFirst, "Change &1st string", idComboSetFirstT, &pg.textFirst, "")
		btnSeveral := wx.NewButton(boxM.Ptr(), idComboMany, "&Append a few strings", -1, -1, -1, -1, 0)
		addExpandBorder(middle, btnSeveral.Ptr())
		btnMany := wx.NewButton(boxM.Ptr(), idComboMany+100, "Append &many strings", -1, -1, -1, -1, 0)
		addExpandBorder(middle, btnMany.Ptr())
		addRow(middle, boxM.Ptr(), idComboChange, "C&hange current", idComboChangeTxt, &pg.textChange, "")
		addRow(middle, boxM.Ptr(), idComboDelete, "&Delete this item", idComboDeleteTxt, &pg.textDelete, "")
		btnDelSel := wx.NewButton(boxM.Ptr(), idComboDeleteSel, "Delete &selection", -1, -1, -1, -1, 0)
		addExpandBorder(middle, btnDelSel.Ptr())
		btnClear := wx.NewButton(boxM.Ptr(), idComboClear, "&Clear", -1, -1, -1, -1, 0)
		addExpandBorder(middle, btnClear.Ptr())
		addRow(middle, boxM.Ptr(), idComboSetValue, "SetValue", idComboValueText, &pg.textValue, "")

		wx.WindowFromPtr(btnSeveral.Ptr()).Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.combo == nil {
				return
			}
			pg.combo.Append("First")
			pg.combo.Append("another one")
			pg.combo.Append("and the last (very very very very very very very very very very long) one")
		})
		wx.WindowFromPtr(btnMany.Ptr()).Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.combo == nil {
				return
			}
			for n := 0; n < 1000; n++ {
				pg.combo.Append("item #" + strconv.Itoa(n))
			}
		})
		wx.WindowFromPtr(btnDelSel.Ptr()).Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.combo == nil {
				return
			}
			if sel := pg.combo.GetSelection(); sel != wx.NOT_FOUND {
				pg.combo.Delete(sel)
			}
		})
		wx.WindowFromPtr(btnClear.Ptr()).Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.combo != nil {
				pg.combo.Clear()
			}
		})

		pg.win.BindId(wx.EVT_BUTTON, idComboSetCur, func(evt *wx.Event) {
			if pg.combo == nil || pg.textCur == nil {
				return
			}
			if n, err := strconv.Atoi(pg.textCur.GetValue()); err == nil {
				pg.combo.SetSelection(n)
			}
		})
		pg.win.BindId(wx.EVT_BUTTON, idComboInsert, func(evt *wx.Event) { pg.onInsert() })
		pg.win.BindId(wx.EVT_TEXT_ENTER, idComboInsertTxt, func(evt *wx.Event) { pg.onInsert() })
		pg.win.BindId(wx.EVT_BUTTON, idComboAdd, func(evt *wx.Event) { pg.onAdd() })
		pg.win.BindId(wx.EVT_TEXT_ENTER, idComboAddText, func(evt *wx.Event) { pg.onAdd() })
		pg.win.BindId(wx.EVT_BUTTON, idComboChange, func(evt *wx.Event) {
			if pg.combo == nil || pg.textChange == nil {
				return
			}
			if sel := pg.combo.GetSelection(); sel != wx.NOT_FOUND {
				pg.combo.SetString(sel, pg.textChange.GetValue())
			}
		})
		pg.win.BindId(wx.EVT_BUTTON, idComboDelete, func(evt *wx.Event) {
			if pg.combo == nil || pg.textDelete == nil {
				return
			}
			if n, err := strconv.Atoi(pg.textDelete.GetValue()); err == nil && n < pg.combo.GetCount() {
				pg.combo.Delete(n)
			}
		})
		pg.win.BindId(wx.EVT_BUTTON, idComboSetFirst, func(evt *wx.Event) {
			if pg.combo == nil || pg.textFirst == nil {
				return
			}
			if pg.combo.GetCount() == 0 {
				f.logf("The combobox is empty, can't change the first string.")
				return
			}
			pg.combo.SetString(0, pg.textFirst.GetValue())
		})
		pg.win.BindId(wx.EVT_BUTTON, idComboSetValue, func(evt *wx.Event) {
			if pg.combo != nil && pg.textValue != nil {
				pg.combo.SetValue(pg.textValue.GetValue())
			}
		})
	}

	pg.sizerCb = wx.NewBoxSizer(wx.VERTICAL)
	if pg.sizerCb != nil {
		pg.sizerCb.SetMinSize(150, 0)
		pg.combo = wx.NewComboBox(pg.win.Ptr(), idComboCtrl, "", -1, -1, -1, -1, nil, 0)
		addWin(pg.sizerCb, pg.combo.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})
		pg.combo1 = wx.NewComboBox(pg.win.Ptr(), idComboDynamic, "", -1, -1, -1, -1, nil, 0)
		if pg.combo1 != nil {
			pg.combo1.Append("Dynamic ComboBox Test - Click me!")
			pg.combo1.SetSelection(0)
		}
		addSpacerSizer(pg.sizerCb, 20)
		addWin(pg.sizerCb, pg.combo1.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})
		pg.bindComboEvents(f)
	}

	addSub(top, left.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.LEFT)
	})
	addSub(top, middle.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL)
	})
	addSub(top, pg.sizerCb.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.RIGHT)
	})

	pg.reset()
	pg.win.SetSizer(top.Ptr(), false)
	bindStyleCheckboxes(pg.win, func() { pg.createCombo(f) })
}

func (pg *comboboxPage) reset() {
	if pg.chkSort != nil {
		pg.chkSort.SetValue(0)
	}
	if pg.chkReadonly != nil {
		pg.chkReadonly.SetValue(0)
	}
	if pg.chkEnter != nil {
		pg.chkEnter.SetValue(0)
	}
}

func (pg *comboboxPage) onAdd() {
	if pg.combo == nil || pg.textAdd == nil {
		return
	}
	pg.combo.Append(pg.textAdd.GetValue())
	pg.itemCount++
	pg.textAdd.SetValue("test item " + strconv.Itoa(pg.itemCount))
}

func (pg *comboboxPage) onInsert() {
	if pg.combo == nil || pg.textInsert == nil {
		return
	}
	pos := pg.combo.GetSelection()
	if pos < 0 {
		pos = 0
	}
	wx.ItemContainerFromPtr(pg.combo.Ptr()).Insert(pg.textInsert.GetValue(), uint(pos))
	pg.itemCount++
	pg.textInsert.SetValue("test item " + strconv.Itoa(pg.itemCount))
}

func (pg *comboboxPage) bindComboEvents(f *widgetsFrame) {
	if pg.combo == nil {
		return
	}
	pg.combo.Bind(wx.EVT_COMBOBOX, func(evt *wx.Event) {
		cmd := wx.AsCommandEvent(evt)
		if cmd == nil || pg.combo == nil {
			return
		}
		sel := cmd.GetSelection()
		if pg.textDelete != nil {
			pg.textDelete.SetValue(strconv.Itoa(sel))
		}
		if pg.textCur != nil {
			pg.textCur.SetValue(strconv.Itoa(sel))
		}
		f.logf("Combobox item %d selected, client data = %s", sel, pg.combo.GetValue())
	})
	pg.combo.Bind(wx.EVT_COMBOBOX_DROPDOWN, func(evt *wx.Event) {
		f.logf("Combobox dropped down")
	})
	pg.combo.Bind(wx.EVT_COMBOBOX_CLOSEUP, func(evt *wx.Event) {
		f.logf("Combobox closed up")
	})
	pg.combo.Bind(wx.EVT_TEXT, func(evt *wx.Event) {
		cmd := wx.AsCommandEvent(evt)
		if cmd != nil {
			f.logf("Combobox text changed to \"%s\"", cmd.GetString())
		}
	})
}

// createCombo is ComboboxWidgetsPage::CreateCombo().
func (pg *comboboxPage) createCombo(f *widgetsFrame) {
	if pg.sizerCb == nil || pg.combo == nil {
		return
	}

	flags := f.attrs.defaultFlags
	if pg.chkSort != nil && pg.chkSort.GetValue() {
		flags |= wx.CB_SORT
	}
	if pg.chkReadonly != nil && pg.chkReadonly.GetValue() {
		flags |= wx.CB_READONLY
	}
	if pg.chkEnter != nil && pg.chkEnter.GetValue() {
		flags |= wx.TE_PROCESS_ENTER
	}
	if pg.radioKind != nil {
		switch pg.radioKind.GetSelection() {
		case 1:
			flags |= wx.CB_SIMPLE
		case 2:
			flags = wx.CB_DROPDOWN
		}
	}

	var items []string
	for n := 0; n < pg.combo.GetCount(); n++ {
		ic := wx.ItemContainerFromPtr(pg.combo.Ptr())
		items = append(items, ic.GetString(uint(n)))
	}
	sel := pg.combo.GetSelection()

	pg.combo.Destroy()
	pg.combo = nil
	// keep the dynamic combobox in place: only the first child is replaced
	clearSizer(pg.sizerCb)

	pg.combo = wx.NewComboBox(pg.win.Ptr(), idComboCtrl, "", -1, -1, -1, -1, nil, flags)
	if pg.combo == nil {
		return
	}
	for _, s := range items {
		pg.combo.Append(s)
	}
	if sel >= 0 && sel < len(items) {
		pg.combo.SetSelection(sel)
	}
	pg.bindComboEvents(f)

	addWin(pg.sizerCb, pg.combo.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})
	addSpacerSizer(pg.sizerCb, 20)
	addWin(pg.sizerCb, pg.combo1.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})
	pg.sizerCb.Layout()
	pg.notifyWidgetRecreation(f, pg.combo.Ptr())
}

// --- Spin page (spinbtn.cpp) -----------------------------------------------

var (
	idSpinReset     = wx.ID_HIGHEST
	idSpinClear     = wx.ID_HIGHEST + 1
	idSpinSetValue  = wx.ID_HIGHEST + 2
	idSpinSetMinMax = wx.ID_HIGHEST + 3
	idSpinSetBase   = wx.ID_HIGHEST + 4
	idSpinSetIncr   = wx.ID_HIGHEST + 5
	idSpinCurText   = wx.ID_HIGHEST + 6
	idSpinValueText = wx.ID_HIGHEST + 7
	idSpinMinText   = wx.ID_HIGHEST + 8
	idSpinMaxText   = wx.ID_HIGHEST + 9
	idSpinBaseText  = wx.ID_HIGHEST + 10
	idSpinIncrText  = wx.ID_HIGHEST + 11
	idSpinBtn       = wx.ID_HIGHEST + 12
	idSpinCtrl      = wx.ID_HIGHEST + 13
	idSpinCtrlDbl   = wx.ID_HIGHEST + 14
)

type spinPage struct {
	pageBase

	chkVert    *wx.CheckBox
	chkArrows  *wx.CheckBox
	chkWrap    *wx.CheckBox
	chkEnter   *wx.CheckBox
	radioAlign *wx.RadioBox
	textCur    *wx.TextCtrl
	textValue  *wx.TextCtrl
	textMin    *wx.TextCtrl
	textMax    *wx.TextCtrl
	textBase   *wx.TextCtrl
	textIncr   *wx.TextCtrl
	spinBtn    *wx.SpinButton
	spinCtrl   *wx.SpinCtrl
	spinDbl    *wx.SpinCtrlDouble
	sizerSpin  *wx.BoxSizer

	min, max, base, increment int
}

func init() {
	registerPage("Spin", nativeCtrls|editableCtrls, func(book unsafe.Pointer) *pageBase {
		pg := &spinPage{max: 10, base: 10, increment: 1}
		pb := newPageWithIcon(book, "icons/spinbtn.svg", 32, 32)
		if pb == nil {
			return nil
		}
		pg.pageBase = *pb
		pg.content = pg.create
		pg.getWidgets = func() []unsafe.Pointer {
			var widgets []unsafe.Pointer
			if pg.spinBtn != nil {
				widgets = append(widgets, pg.spinBtn.Ptr())
			}
			if pg.spinCtrl != nil {
				widgets = append(widgets, pg.spinCtrl.Ptr())
			}
			if pg.spinDbl != nil {
				widgets = append(widgets, pg.spinDbl.Ptr())
			}
			return widgets
		}
		pg.recreate = func(f *widgetsFrame) { pg.createSpin(f) }
		return &pg.pageBase
	})
}

func (pg *spinPage) create(f *widgetsFrame) {
	top := wx.NewBoxSizer(wx.HORIZONTAL)
	if top == nil {
		return
	}

	boxL, left := addStaticBox(pg.win.Ptr(), "&Set style", wx.VERTICAL)
	if left != nil {
		pg.chkVert = pg.createCheckBoxAndAddToSizer(left, "&Vertical", wx.ID_ANY, boxL.Ptr())
		pg.chkVert.SetValue(1)
		pg.chkArrows = pg.createCheckBoxAndAddToSizer(left, "&Arrow Keys", wx.ID_ANY, boxL.Ptr())
		pg.chkArrows.SetValue(1)
		pg.chkWrap = pg.createCheckBoxAndAddToSizer(left, "&Wrap", wx.ID_ANY, boxL.Ptr())
		pg.chkEnter = pg.createCheckBoxAndAddToSizer(left, "Process &Enter", wx.ID_ANY, boxL.Ptr())
		addSpacerSizer(left, 5)

		pg.radioAlign = newRadioBox(boxL.Ptr(), wx.ID_ANY, "&Text alignment",
			[]string{"left", "centre", "right"}, 1, 0)
		if pg.radioAlign != nil {
			pg.radioAlign.SetSelection(2)
		}
		addExpandBorder(left, pg.radioAlign.Ptr())
		addSpacerSizer(left, 5)

		reset := wx.NewButton(boxL.Ptr(), idSpinReset, "&Reset", -1, -1, -1, -1, 0)
		addWin(left, reset.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.CenterHorizontal()
			sf.TripleBorder(wx.ALL)
		})
		reset.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			pg.reset()
			pg.createSpin(f)
		})
	}

	boxM, middle := addStaticBox(pg.win.Ptr(), "&Change spinbtn value", wx.VERTICAL)
	if middle != nil {
		addRow(middle, boxM.Ptr(), idSpinCurText+1000, "Current value", idSpinCurText, &pg.textCur, "")
		if pg.textCur != nil {
			pg.textCur.SetEditable(false)
		}
		addRow(middle, boxM.Ptr(), idSpinSetValue, "Set &value", idSpinValueText, &pg.textValue, "")
		row := pg.createSizerWithTextAndButton(idSpinSetMinMax, "&Min and max", idSpinMinText, &pg.textMin, boxM.Ptr())
		addSub(middle, row.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})
		if pg.textMin != nil {
			pg.textMin.SetValue("0")
		}
		pg.textMax = wx.NewTextCtrl(boxM.Ptr(), idSpinMaxText, "10", -1, -1, -1, -1, 0)
		addWin(row, pg.textMax.Ptr(), 1, func(sf *wx.SizerFlags) {
			sf.CenterVertical()
			sf.Border(wx.LEFT)
		})
		addRow(middle, boxM.Ptr(), idSpinSetBase, "Set &base", idSpinBaseText, &pg.textBase, "10")
		addRow(middle, boxM.Ptr(), idSpinSetIncr, "Set Increment", idSpinIncrText, &pg.textIncr, "1")

		pg.win.BindId(wx.EVT_BUTTON, idSpinSetValue, func(evt *wx.Event) {
			if pg.textValue == nil {
				return
			}
			if pg.textValue.GetValue() == "" {
				pg.setValue(0)
				return
			}
			if n, err := strconv.Atoi(pg.textValue.GetValue()); err == nil {
				pg.setValue(n)
			}
		})
		pg.win.BindId(wx.EVT_BUTTON, idSpinSetMinMax, func(evt *wx.Event) {
			if pg.textMin == nil || pg.textMax == nil {
				return
			}
			mn, err1 := strconv.Atoi(pg.textMin.GetValue())
			mx, err2 := strconv.Atoi(pg.textMax.GetValue())
			if err1 != nil || err2 != nil {
				return
			}
			if mn > mx {
				f.logf("Invalid min/max range.")
				return
			}
			pg.min, pg.max = mn, mx
			if pg.spinBtn != nil {
				pg.spinBtn.SetRange(pg.min, pg.max)
			}
			if pg.spinCtrl != nil {
				pg.spinCtrl.SetRange(pg.min, pg.max)
			}
			if pg.spinDbl != nil {
				pg.spinDbl.SetRange(float64(pg.min), float64(pg.max))
			}
			pg.win.Layout()
		})
		pg.win.BindId(wx.EVT_BUTTON, idSpinSetBase, func(evt *wx.Event) {
			if pg.textBase == nil || pg.spinCtrl == nil {
				return
			}
			if n, err := strconv.Atoi(pg.textBase.GetValue()); err == nil && n != 0 {
				pg.base = n
				if !pg.spinCtrl.SetBase(n) {
					f.logf("Invalid base.")
				}
			}
		})
		pg.win.BindId(wx.EVT_BUTTON, idSpinSetIncr, func(evt *wx.Event) {
			if pg.textIncr == nil || pg.spinCtrl == nil {
				return
			}
			if n, err := strconv.Atoi(pg.textIncr.GetValue()); err == nil && n != 0 {
				pg.increment = n
				pg.spinCtrl.SetIncrement(n)
			}
		})
	}

	pg.sizerSpin = wx.NewBoxSizer(wx.VERTICAL)
	if pg.sizerSpin != nil {
		pg.sizerSpin.SetMinSize(150, 0)
	}

	pg.reset()
	pg.createSpin(f)

	addSub(top, left.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.LEFT)
	})
	addSub(top, middle.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL)
	})
	addSub(top, pg.sizerSpin.Ptr(), 1, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.RIGHT)
	})

	pg.win.SetSizer(top.Ptr(), false)
	bindStyleCheckboxes(pg.win, func() { pg.createSpin(f) })
}

func (pg *spinPage) reset() {
	if pg.chkVert != nil {
		pg.chkVert.SetValue(1)
	}
	if pg.chkArrows != nil {
		pg.chkArrows.SetValue(1)
	}
	if pg.chkWrap != nil {
		pg.chkWrap.SetValue(0)
	}
	if pg.chkEnter != nil {
		pg.chkEnter.SetValue(0)
	}
	if pg.radioAlign != nil {
		pg.radioAlign.SetSelection(2)
	}
}

func (pg *spinPage) setValue(v int) {
	if pg.spinBtn != nil {
		pg.spinBtn.SetValue(v)
	}
	if pg.spinCtrl != nil {
		pg.spinCtrl.SetValue(v)
	}
	if pg.spinDbl != nil {
		pg.spinDbl.SetValue(float64(v))
	}
}

// createSpin is SpinBtnWidgetsPage::CreateSpin().
func (pg *spinPage) createSpin(f *widgetsFrame) {
	if pg.sizerSpin == nil {
		return
	}

	flags := f.attrs.defaultFlags
	if pg.chkArrows != nil && pg.chkArrows.GetValue() {
		flags |= wx.SP_ARROW_KEYS
	}
	if pg.chkWrap != nil && pg.chkWrap.GetValue() {
		flags |= wx.SP_WRAP
	}
	if pg.chkEnter != nil && pg.chkEnter.GetValue() {
		flags |= wx.TE_PROCESS_ENTER
	}

	textFlags := 0
	if pg.radioAlign != nil {
		switch pg.radioAlign.GetSelection() {
		case 1:
			textFlags = wx.ALIGN_CENTRE_HORIZONTAL
		case 2:
			textFlags = wx.ALIGN_RIGHT
		default:
			textFlags = wx.ALIGN_LEFT
		}
	}

	val := pg.min
	if pg.spinBtn != nil {
		val = pg.spinBtn.GetValue()
		pg.spinBtn.Destroy()
		pg.spinBtn = nil
	}
	destroyWin(pg.spinCtrl)
	destroyWin(pg.spinDbl)
	pg.spinCtrl = nil
	pg.spinDbl = nil
	// the sizer also holds the three static labels, which are recreated below
	clearSizerDelete(pg.sizerSpin)

	pg.spinBtn = wx.NewSpinButton(pg.win.Ptr(), idSpinBtn, -1, -1, -1, -1, flags)
	if pg.spinBtn != nil {
		pg.spinBtn.SetValue(pg.min)
		pg.spinBtn.SetRange(pg.min, pg.max)
	}

	pg.spinCtrl = wx.NewSpinCtrl(pg.win.Ptr(), idSpinCtrl, strconv.Itoa(val), -1, -1, -1, -1,
		flags|textFlags, pg.min, pg.max, val)
	pg.spinDbl = wx.NewSpinCtrlDouble(pg.win.Ptr(), idSpinCtrlDbl, strconv.Itoa(val), -1, -1, -1, -1,
		flags|textFlags, float64(pg.min), float64(pg.max), float64(val), 0.1)

	// layout: stretch, label, control, repeated (as the C++ sizer)
	children := []struct {
		label string
		win   unsafe.Pointer
	}{
		{"wxSpinButton", ptrOf(pg.spinBtn)},
		{"wxSpinCtrl", ptrOf(pg.spinCtrl)},
		{"wxSpinCtrlDouble", ptrOf(pg.spinDbl)},
	}
	pg.sizerSpin.AddStretchSpacer(0)
	for _, ch := range children {
		st := wx.NewStaticText(pg.win.Ptr(), wx.ID_ANY, ch.label, -1, -1, -1, -1, 0)
		addWin(pg.sizerSpin, st.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Center()
			sf.Border(wx.ALL)
		})
		addWin(pg.sizerSpin, ch.win, 0, func(sf *wx.SizerFlags) {
			sf.Center()
			sf.Border(wx.ALL)
		})
		pg.sizerSpin.AddStretchSpacer(0)
	}
	pg.sizerSpin.Layout()

	for _, handle := range []unsafe.Pointer{ptrOf(pg.spinBtn), ptrOf(pg.spinCtrl), ptrOf(pg.spinDbl)} {
		pg.notifyWidgetRecreation(f, handle)
	}
}

func ptrOf[T any](p *T) unsafe.Pointer {
	if p == nil {
		return nil
	}
	type ptrable interface{ Ptr() unsafe.Pointer }
	if v, ok := any(p).(ptrable); ok {
		return v.Ptr()
	}
	return nil
}

// --- shared page helpers ---------------------------------------------------

// addExpandBorder adds the window with Expand().Border().
func addExpandBorder(s sizerTarget, window unsafe.Pointer) {
	addWin(s, window, 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})
}

// addSpacerSizer adds a fixed vertical spacer to the sizer.
func addSpacerSizer(s sizerTarget, size int) {
	switch v := s.(type) {
	case *wx.BoxSizer:
		v.AddSpacer(size)
	case *wx.StaticBoxSizer:
		v.AddSpacer(size)
	}
}

// addRow creates the "button + text ctrl" row used by the item pages.
func addRow(sizer sizerTarget, parent unsafe.Pointer, idBtn int, label string, idText int, ppText **wx.TextCtrl, value string) {
	btn := wx.NewButton(parent, idBtn, label, -1, -1, -1, -1, 0)
	if btn == nil {
		return
	}
	text := wx.NewTextCtrl(parent, idText, "", -1, -1, -1, -1, wx.TE_PROCESS_ENTER)
	if text == nil {
		return
	}
	row := wx.NewBoxSizer(wx.HORIZONTAL)
	if row == nil {
		return
	}
	addWin(row, btn.Ptr(), 0, func(sf *wx.SizerFlags) { sf.Border(wx.RIGHT) })
	addWin(row, text.Ptr(), 1, func(sf *wx.SizerFlags) { sf.Border(wx.LEFT) })
	addSub(sizer, row.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})
	if value != "" {
		text.SetValue(value)
	}
	if ppText != nil {
		*ppText = text
	}
}
