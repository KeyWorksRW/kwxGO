// Book control pages: Notebook, Listbook and Choicebook (notebook.cpp).
package main

import (
	"strconv"
	"unsafe"

	"github.com/KeyWorksRW/kwxGO/wx"
)

var (
	idBookReset         = wx.ID_HIGHEST
	idBookSelectPage    = wx.ID_HIGHEST + 1
	idBookAddPage       = wx.ID_HIGHEST + 2
	idBookInsertPage    = wx.ID_HIGHEST + 3
	idBookRemovePage    = wx.ID_HIGHEST + 4
	idBookDeleteAll     = wx.ID_HIGHEST + 5
	idBookInsertText    = wx.ID_HIGHEST + 6
	idBookRemoveText    = wx.ID_HIGHEST + 7
	idBookSelectText    = wx.ID_HIGHEST + 8
	idBookNumPagesText  = wx.ID_HIGHEST + 9
	idBookCurSelectText = wx.ID_HIGHEST + 10
	idBookCtl           = wx.ID_HIGHEST + 11
)

// bookKind selects which concrete book control a bookPage creates.
type bookKind int

const (
	bookNotebook bookKind = iota
	bookListbook
	bookChoicebook
)

type bookPage struct {
	pageBase

	kind bookKind

	chkImages   *wx.CheckBox
	radioOrient *wx.RadioBox

	textInsert    *wx.TextCtrl
	textRemove    *wx.TextCtrl
	textSelect    *wx.TextCtrl
	textNumPages  *wx.TextCtrl
	textCurSelect *wx.TextCtrl

	book      *wx.BookCtrlBase
	sizerBook *wx.BoxSizer
	imageList *wx.ImageList
}

func newBookPage(kind bookKind, icon string) func(book unsafe.Pointer) *pageBase {
	return func(book unsafe.Pointer) *pageBase {
		pg := &bookPage{kind: kind}
		pb := newPageWithIcon(book, icon, 32, 32)
		if pb == nil {
			return nil
		}
		pg.pageBase = *pb
		pg.content = pg.create
		pg.getWidgets = func() []unsafe.Pointer {
			if pg.book == nil {
				return nil
			}
			return []unsafe.Pointer{pg.book.Ptr()}
		}
		pg.recreate = func(f *widgetsFrame) { pg.recreateBook(f) }
		return &pg.pageBase
	}
}

func init() {
	registerPage("Notebook", nativeCtrls|bookCtrls, newBookPage(bookNotebook, "icons/notebook.svg"))
	registerPage("Listbook", genericCtrls|bookCtrls, newBookPage(bookListbook, "icons/listbook.svg"))
	registerPage("Choicebook", genericCtrls|bookCtrls, newBookPage(bookChoicebook, "icons/choicebk.svg"))
}

func (pg *bookPage) create(f *widgetsFrame) {
	top := wx.NewBoxSizer(wx.HORIZONTAL)
	if top == nil {
		return
	}

	boxL, left := addStaticBox(pg.win.Ptr(), "&Set style", wx.VERTICAL)
	if left != nil {
		pg.chkImages = wx.NewCheckBox(boxL.Ptr(), wx.ID_ANY, "Show &images", -1, -1, -1, -1, 0)
		addWin(left, pg.chkImages.Ptr(), 0, func(sf *wx.SizerFlags) { sf.Border(wx.ALL) })
		addSpacerSizer(left, 15)
		pg.radioOrient = newRadioBox(boxL.Ptr(), wx.ID_ANY, "&Tab orientation",
			[]string{"&top", "&bottom", "&left", "&right"}, 1, wx.RA_SPECIFY_COLS)
		addWin(left, pg.radioOrient.Ptr(), 0, func(sf *wx.SizerFlags) { sf.Border(wx.ALL) })

		reset := wx.NewButton(boxL.Ptr(), idBookReset, "&Reset", -1, -1, -1, -1, 0)
		addWin(left, reset.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.CenterHorizontal()
			sf.TripleBorder(wx.ALL)
		})
		reset.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			pg.reset()
			pg.recreateBook(f)
		})
	}

	boxM, middle := addStaticBox(pg.win.Ptr(), "&Contents", wx.VERTICAL)
	if middle != nil {
		var numText *wx.TextCtrl
		row := pg.createSizerWithTextAndLabel("Number of pages: ", idBookNumPagesText, &numText, boxM.Ptr())
		pg.textNumPages = numText
		if numText != nil {
			numText.SetEditable(false)
			addSub(middle, row.Ptr(), 0, func(sf *wx.SizerFlags) {
				sf.Expand()
				sf.Border(wx.ALL)
			})
		}

		var curText *wx.TextCtrl
		row = pg.createSizerWithTextAndLabel("Current selection: ", idBookCurSelectText, &curText, boxM.Ptr())
		pg.textCurSelect = curText
		if curText != nil {
			curText.SetEditable(false)
			addSub(middle, row.Ptr(), 0, func(sf *wx.SizerFlags) {
				sf.Expand()
				sf.Border(wx.ALL)
			})
		}

		row = pg.createSizerWithTextAndButton(idBookSelectPage, "&Select page",
			idBookSelectText, &pg.textSelect, boxM.Ptr())
		if row != nil {
			addSub(middle, row.Ptr(), 0, func(sf *wx.SizerFlags) {
				sf.Expand()
				sf.Border(wx.ALL)
			})
		}

		addPage := wx.NewButton(boxM.Ptr(), idBookAddPage, "&Add page", -1, -1, -1, -1, 0)
		addWin(middle, addPage.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})
		addPage.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.book != nil {
				pg.book.AddPage(pg.createNewPage(), "Added page", false, pg.iconIndex())
			}
		})

		pg.textInsert = nil
		row = pg.createSizerWithTextAndButton(idBookInsertPage, "&Insert page at",
			idBookInsertText, &pg.textInsert, boxM.Ptr())
		if row != nil {
			addSub(middle, row.Ptr(), 0, func(sf *wx.SizerFlags) {
				sf.Expand()
				sf.Border(wx.ALL)
			})
		}

		pg.textRemove = nil
		row = pg.createSizerWithTextAndButton(idBookRemovePage, "&Remove page",
			idBookRemoveText, &pg.textRemove, boxM.Ptr())
		if row != nil {
			addSub(middle, row.Ptr(), 0, func(sf *wx.SizerFlags) {
				sf.Expand()
				sf.Border(wx.ALL)
			})
		}

		deleteAll := wx.NewButton(boxM.Ptr(), idBookDeleteAll, "&Delete All", -1, -1, -1, -1, 0)
		addWin(middle, deleteAll.Ptr(), 0, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})
		deleteAll.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			if pg.book != nil {
				pg.book.DeleteAllPages()
			}
		})
	}

	pg.sizerBook = wx.NewBoxSizer(wx.HORIZONTAL)

	addSub(top, left.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL & ^wx.LEFT)
	})
	addSub(top, middle.Ptr(), 0, func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.DoubleBorder(wx.ALL)
	})
	if pg.sizerBook != nil {
		addSub(top, pg.sizerBook.Ptr(), 1, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.DoubleBorder(wx.ALL & ^wx.RIGHT)
		})
	}

	pg.recreateBook(f)
	pg.reset()
	pg.createImageList()

	pg.win.SetSizer(top.Ptr(), false)

	pg.win.BindId(wx.EVT_BUTTON, idBookSelectPage, func(evt *wx.Event) {
		if pos := textValue(pg.textSelect); pg.book != nil && pg.validValue(pos) {
			pg.book.SetSelection(pos)
		}
	})
	pg.win.BindId(wx.EVT_BUTTON, idBookInsertPage, func(evt *wx.Event) {
		if pos := textValue(pg.textInsert); pg.book != nil && pg.validValue(pos) {
			pg.book.InsertPage(pos, pg.createNewPage(), "Inserted page", false, pg.iconIndex())
		}
	})
	pg.win.BindId(wx.EVT_BUTTON, idBookRemovePage, func(evt *wx.Event) {
		if pos := textValue(pg.textRemove); pg.book != nil && pg.validValue(pos) {
			pg.book.DeletePage(pos)
		}
	})
	pg.win.BindId(wx.EVT_UPDATE_UI, idBookNumPagesText, func(evt *wx.Event) {
		if ui := wx.AsUpdateUIEvent(evt); ui != nil && pg.book != nil {
			ui.SetText(strconv.Itoa(pg.book.GetPageCount()))
		}
	})
	pg.win.BindId(wx.EVT_UPDATE_UI, idBookCurSelectText, func(evt *wx.Event) {
		if ui := wx.AsUpdateUIEvent(evt); ui != nil && pg.book != nil {
			ui.SetText(strconv.Itoa(pg.book.GetSelection()))
		}
	})
	pg.win.BindId(wx.EVT_UPDATE_UI, idBookSelectPage, func(evt *wx.Event) {
		if ui := wx.AsUpdateUIEvent(evt); ui != nil {
			ui.Enable(pg.validValue(textValue(pg.textSelect)))
		}
	})
	pg.win.BindId(wx.EVT_UPDATE_UI, idBookInsertPage, func(evt *wx.Event) {
		if ui := wx.AsUpdateUIEvent(evt); ui != nil {
			ui.Enable(pg.validValue(textValue(pg.textInsert)))
		}
	})
	pg.win.BindId(wx.EVT_UPDATE_UI, idBookRemovePage, func(evt *wx.Event) {
		if ui := wx.AsUpdateUIEvent(evt); ui != nil {
			ui.Enable(pg.validValue(textValue(pg.textRemove)))
		}
	})
	pg.win.Bind(wx.EVT_CHECKBOX, func(evt *wx.Event) { pg.recreateBook(f) })
	pg.win.Bind(wx.EVT_RADIOBOX, func(evt *wx.Event) { pg.recreateBook(f) })
}

func textValue(t *wx.TextCtrl) int {
	if t == nil {
		return -1
	}
	if v, err := strconv.Atoi(t.GetValue()); err == nil {
		return v
	}
	return -1
}

func (pg *bookPage) validValue(v int) bool {
	return pg.book != nil && v >= 0 && v < pg.book.GetPageCount()
}

func (pg *bookPage) reset() {
	setChk(pg.chkImages, 1)
	if pg.radioOrient != nil {
		pg.radioOrient.SetSelection(0)
	}
}

// createImageList is BookWidgetsPage::CreateImageList().
func (pg *bookPage) createImageList() {
	if !checked(pg.chkImages) {
		pg.imageList = nil
		return
	}
	if pg.imageList == nil {
		pg.imageList = wx.NewImageList(32, 32, 1, 4)
		if pg.imageList == nil {
			return
		}
		for _, id := range []string{wx.ART_INFORMATION, wx.ART_QUESTION, wx.ART_WARNING, wx.ART_ERROR} {
			if bmp := wx.ArtProvider_GetBitmap(id, wx.ART_OTHER, 32, 32); bmp != nil {
				pg.imageList.Add(bmp.Ptr(), nil)
				bmp.Delete()
			}
		}
	}
	if pg.book != nil {
		pg.book.SetImageList(pg.imageList.Ptr())
	}
}

// iconIndex is BookWidgetsPage::GetIconIndex().
func (pg *bookPage) iconIndex() int {
	if pg.book != nil && checked(pg.chkImages) {
		if n := pg.book.GetImageCount(); n > 0 {
			return pg.book.GetPageCount() % n
		}
	}
	return -1
}

func (pg *bookPage) createNewPage() unsafe.Pointer {
	if pg.book == nil {
		return nil
	}
	tc := wx.NewTextCtrl(pg.book.Ptr(), wx.ID_ANY, "I'm a book page", -1, -1, -1, -1, 0)
	if tc == nil {
		return nil
	}
	return tc.Ptr()
}

// recreateBook is BookWidgetsPage::RecreateBook().
func (pg *bookPage) recreateBook(f *widgetsFrame) {
	if pg.radioOrient == nil {
		return
	}

	flags := f.attrs.defaultFlags
	switch selection(pg.radioOrient) {
	case 1:
		flags |= wx.BK_BOTTOM
	case 2:
		flags |= wx.BK_LEFT
	case 3:
		flags |= wx.BK_RIGHT
	default:
		flags |= wx.BK_TOP
	}

	old := pg.book
	pg.book = nil

	var created *wx.BookCtrlBase
	switch pg.kind {
	case bookNotebook:
		if nb := wx.NewNotebook(pg.win.Ptr(), idBookCtl, -1, -1, -1, -1, flags); nb != nil {
			created = wx.BookCtrlBaseFromPtr(nb.Ptr())
		}
	case bookListbook:
		if lb := wx.NewListbook(pg.win.Ptr(), idBookCtl, -1, -1, -1, -1, flags); lb != nil {
			created = wx.BookCtrlBaseFromPtr(lb.Ptr())
		}
	case bookChoicebook:
		if cb := wx.NewChoicebook(pg.win.Ptr(), idBookCtl, -1, -1, -1, -1, flags); cb != nil {
			created = wx.BookCtrlBaseFromPtr(cb.Ptr())
		}
	}
	if created == nil {
		return
	}
	pg.book = created
	pg.notifyWidgetRecreation(f, pg.book.Ptr())
	pg.createImageList()

	if old != nil {
		sel := old.GetSelection()
		count := old.GetPageCount()
		for n := 0; n < count; n++ {
			pg.book.AddPage(pg.createNewPage(), old.GetPageText(n), false, pg.iconIndex())
		}
		if pg.sizerBook != nil {
			pg.sizerBook.Remove(0)
		}
		// the book was wrapped with *FromPtr, so Delete() is a no-op: destroy the
		// window explicitly, which also destroys the pages it owns.
		wx.WindowFromPtr(old.Ptr()).Destroy()
		if sel != -1 {
			pg.book.SetSelection(sel)
		}
	}

	if pg.sizerBook != nil {
		addWin(pg.sizerBook, pg.book.Ptr(), 1, func(sf *wx.SizerFlags) {
			sf.Expand()
			sf.Border(wx.ALL)
		})
		pg.sizerBook.SetMinSize(150, 0)
		pg.sizerBook.Layout()
	}
}
