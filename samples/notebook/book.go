package main

import (
	"strings"
	"unsafe"

	"github.com/KeyWorksRW/kwxGO/wx"
)

// bookType mirrors the C++ MyFrame::BookType enum. The numeric order matters:
// OnType() derives the type from the menu ID, and the initial radio item is
// checked with ID_BOOK_NOTEBOOK + bookType.
type bookType int

const (
	typeNotebook bookType = iota
	typeListbook
	typeChoicebook
	typeTreebook
	typeToolbook
	typeAuiNotebook
	typeSimplebook
)

// book is the common surface of every wxBookCtrlBase subclass used by the
// sample, so the page operations and the shared geometry queries can be
// expressed once.
type book interface {
	Ptr() unsafe.Pointer
	AddPage(page unsafe.Pointer, text string, sel bool, imageId int) bool
	InsertPage(index int, page unsafe.Pointer, text string, sel bool, imageId int) bool
	DeletePage(nPage int) bool
	GetPage(nPage int) unsafe.Pointer
	GetPageCount() int
	GetPageText(nPage int) string
	GetSelection() int
	SetSelection(nPage int) int
	ChangeSelection(nPage int) int
	AdvanceSelection(forward bool)
	SetPageSize(width int, height int)
	SetImages(images unsafe.Pointer)
	GetImageCount() int
	GetSize() wx.Size
	ScreenToClient(x int, y int) wx.Point
	HitTestFlags(x int, y int, flags unsafe.Pointer) int
	Hide() bool
	Destroy() bool
}

// createBook performs the DISPATCH_ON_TYPE() construction from the C++ sample.
// Each constructor is checked for nil before being stored in the interface so a
// failed creation yields a nil interface rather than a typed nil pointer.
func (f *sampleFrame) createBook(flags int) book {
	parent := f.panel.Ptr()
	switch f.bookTyp {
	case typeNotebook:
		if w := wx.NewNotebook(parent, wx.ID_ANY, -1, -1, -1, -1, flags); w != nil {
			return w
		}
	case typeListbook:
		if w := wx.NewListbook(parent, wx.ID_ANY, -1, -1, -1, -1, flags); w != nil {
			return w
		}
	case typeChoicebook:
		if w := wx.NewChoicebook(parent, wx.ID_ANY, -1, -1, -1, -1, flags); w != nil {
			return w
		}
	case typeTreebook:
		if w := wx.NewTreebook(parent, wx.ID_ANY, -1, -1, -1, -1, flags); w != nil {
			return w
		}
	case typeToolbook:
		if w := wx.NewToolbook(parent, wx.ID_ANY, -1, -1, -1, -1, flags); w != nil {
			return w
		}
	case typeAuiNotebook:
		if w := wx.NewAuiNotebook(parent, wx.ID_ANY, -1, -1, -1, -1, flags); w != nil {
			return w
		}
	case typeSimplebook:
		if w := wx.NewSimplebook(parent, wx.ID_ANY, -1, -1, -1, -1, flags); w != nil {
			return w
		}
	}
	return nil
}

// recreateBook is MyFrame::RecreateBook(): it builds the book control for the
// current type/orientation/style, moves the pages over, and swaps it into the
// sizer.
func (f *sampleFrame) recreateBook() {
	flags := wx.BK_DEFAULT
	switch f.orient {
	case idOrientTop:
		flags = wx.BK_TOP
	case idOrientBottom:
		flags = wx.BK_BOTTOM
	case idOrientLeft:
		flags = wx.BK_LEFT
	case idOrientRight:
		flags = wx.BK_RIGHT
	}

	if f.fixedWidth && f.bookTyp == typeNotebook {
		flags |= wx.NB_FIXEDWIDTH
	}
	if f.multi && f.bookTyp == typeNotebook {
		flags |= wx.NB_MULTILINE
	}
	if f.noPageTheme && f.bookTyp == typeNotebook {
		flags |= wx.NB_NOPAGETHEME
	}
	if f.buttonBar && f.bookTyp == typeToolbook {
		flags |= wx.TBK_BUTTONBAR
	}
	if f.horzLayout && f.bookTyp == typeToolbook {
		flags |= wx.TBK_HORZ_LAYOUT
	}

	oldBook := f.book
	f.book = nil

	newBook := f.createBook(flags)
	if newBook == nil {
		return
	}
	f.book = newBook

	newBook.Hide()

	// wxToolbook doesn't work without icons so always use them for it.
	if f.chkShowImages || f.bookTyp == typeToolbook {
		if f.images != nil {
			newBook.SetImages(f.images.Ptr())
		}
	}

	if oldBook != nil {
		// wxTreebook complication: pages may have a parent page.
		var oldTree *wx.Treebook
		if f.bookTyp == typeTreebook {
			oldTree, _ = oldBook.(*wx.Treebook)
		}

		count := oldBook.GetPageCount()
		for n := 0; n < count; n++ {
			image := f.iconIndex()
			str := oldBook.GetPageText(n)

			page := createPage(newBook.Ptr(), str)
			if page == nil {
				continue
			}

			if oldTree != nil {
				if parent := oldTree.GetPageParent(n); parent != wx.NOT_FOUND {
					if tbk, ok := newBook.(*wx.Treebook); ok {
						tbk.InsertSubPage(parent, page, str, false, image)
						// skip adding it again below
						continue
					}
				}
			}

			newBook.AddPage(page, str, false, image)
		}

		if sel := oldBook.GetSelection(); sel != wx.NOT_FOUND {
			newBook.SetSelection(sel)
		}

		f.sizerFrame.DetachWindow(oldBook.Ptr())
		// The C++ code deletes the detached book; Destroy() is the wxWidgets
		// sanctioned equivalent (deferred delete).
		oldBook.Destroy()
	} else {
		f.createInitialPages()
	}

	wxf := wx.NewSizerFlags(5)
	wxf.Expand()
	wxf.Border(wx.ALL)
	f.sizerFrame.InsertWindowFlags(0, newBook.Ptr(), wxf.Ptr())
	wxf.Delete()

	f.sizerFrame.ShowWindow(newBook.Ptr(), true, false)
	f.sizerFrame.Layout()
}

// createInitialPages is CreateInitialPages() from the C++ sample.
func (f *sampleFrame) createInitialPages() {
	if f.book == nil {
		return
	}

	f.book.AddPage(createRadioButtonsPage(f.book.Ptr()),
		radioButtonsPageName, false, f.iconIndex())
	f.book.AddPage(createVetoPage(f.book.Ptr()),
		vetoPageName, false, f.iconIndex())
	f.book.AddPage(createFullPageText(f.book.Ptr()),
		textPageName, false, f.iconIndex())
	f.book.InsertPage(0, createInsertPage(f.book.Ptr()),
		iWasInsertedPageName, false, f.iconIndex())

	f.book.SetSelection(1)
}

// iconIndex is GetIconIndex(): the image index for the next page is the page
// count modulo the number of images the current book control holds.
func (f *sampleFrame) iconIndex() int {
	if f.book != nil {
		nImages := f.book.GetImageCount()
		if nImages > 0 {
			return f.book.GetPageCount() % nImages
		}
	}
	return wx.NOT_FOUND
}

// bookHandler captures the control name and the changed/changing direction of
// the page event it is bound to (OnNotebook() and friends in the C++ sample).
func (f *sampleFrame) bookHandler(name string, changing bool) wx.EventHandler {
	return func(evt *wx.Event) {
		f.onBookCtrl(evt, name, changing)
	}
}

// createPage is CreatePage(): pick a page factory based on the page label of
// the page being cloned into the new book control.
func createPage(parent unsafe.Pointer, pageName string) unsafe.Pointer {
	if strings.Contains(pageName, insertedPageName) ||
		strings.Contains(pageName, addedPageName) ||
		strings.Contains(pageName, addedSubPageName) ||
		strings.Contains(pageName, addedPageNameBefore) {
		return createUserCreatedPage(parent)
	}

	if pageName == iWasInsertedPageName {
		return createInsertPage(parent)
	}
	if pageName == vetoPageName {
		return createVetoPage(parent)
	}
	if pageName == radioButtonsPageName {
		return createRadioButtonsPage(parent)
	}
	if pageName == textPageName {
		return createFullPageText(parent)
	}

	wx.LogError("unknown page name")
	return nil
}
