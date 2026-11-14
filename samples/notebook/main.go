// Command notebook is the kwxGO port of the wxWidgets "notebook" sample
// (samples/notebook/notebook.cpp and notebook.h).
//
// It exercises every wxBookCtrlBase-derived control -- wxNotebook, wxListbook,
// wxChoicebook, wxTreebook, wxToolbook, wxAuiNotebook and wxSimplebook -- and
// lets the user switch between them at run time, changing the orientation,
// the notebook style flags, the page images, and the pages themselves.
package main

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/KeyWorksRW/kwxGO/wx"
)

// Page names. These are also used to decide which kind of page to create when
// the current book is cloned into another book control (CreatePage() in C++).
const (
	iWasInsertedPageName = "Inserted"
	radioButtonsPageName = "Radio buttons"
	vetoPageName         = "Veto"
	textPageName         = "Text"

	insertedPageName    = "Inserted "
	addedPageName       = "Added "
	addedPageNameBefore = " Inserted before "
	addedSubPageName    = " Inserted sub-page "
)

// Command IDs, mirroring the C++ ID_COMMANDS enum in notebook.h. They start at
// wxID_HIGHEST and must keep the same relative order, because the "Type" radio
// items derive the book type from the ID (OnType), and because
// ID_BOOK_NOTEBOOK + bookType is used to check the initial radio item.
var (
	idBookNotebook    = wx.ID_HIGHEST
	idBookListbook    = wx.ID_HIGHEST + 1
	idBookChoicebook  = wx.ID_HIGHEST + 2
	idBookTreebook    = wx.ID_HIGHEST + 3
	idBookToolbook    = wx.ID_HIGHEST + 4
	idBookAuiNotebook = wx.ID_HIGHEST + 5
	idBookSimplebook  = wx.ID_HIGHEST + 6
	idBookMax         = wx.ID_HIGHEST + 7 // ID_BOOK_MAX
	idOrientDefault   = wx.ID_HIGHEST + 8
	idOrientTop       = wx.ID_HIGHEST + 9
	idOrientBottom    = wx.ID_HIGHEST + 10
	idOrientLeft      = wx.ID_HIGHEST + 11
	idOrientRight     = wx.ID_HIGHEST + 12
	idOrientMax       = wx.ID_HIGHEST + 13 // ID_ORIENT_MAX
	idShowImages      = wx.ID_HIGHEST + 14
	idFixedWidth      = wx.ID_HIGHEST + 15
	idMulti           = wx.ID_HIGHEST + 16
	idNoPageTheme     = wx.ID_HIGHEST + 17
	idButtonBar       = wx.ID_HIGHEST + 18
	idHorzLayout      = wx.ID_HIGHEST + 19

	idAddPage         = wx.ID_HIGHEST + 20
	idAddPageNoSelect = wx.ID_HIGHEST + 21
	idInsertPage      = wx.ID_HIGHEST + 22
	idDeleteCurPage   = wx.ID_HIGHEST + 23
	idDeleteLastPage  = wx.ID_HIGHEST + 24
	idNextPage        = wx.ID_HIGHEST + 25
	idAddPageBefore   = wx.ID_HIGHEST + 26
	idAddSubPage      = wx.ID_HIGHEST + 27
	idChangeSelection = wx.ID_HIGHEST + 28
	idSetSelection    = wx.ID_HIGHEST + 29
	idGetPageSize     = wx.ID_HIGHEST + 30
	idSetPageSize     = wx.ID_HIGHEST + 31
	idContextHelp     = wx.ID_HIGHEST + 32
	idHitTest         = wx.ID_HIGHEST + 33
)

// sampleFrame holds the state MyFrame keeps in member and function-static
// variables in the C++ original.
type sampleFrame struct {
	frame      *wx.Frame
	panel      *wx.Panel
	sizerFrame *wx.BoxSizer

	logText *wx.TextCtrl
	logCtrl *wx.LogTextCtrl
	// icon is the frame icon built from sampleXPM; the frame copies it, so the
	// wrapper owns the native wxIcon and is released on close.
	icon *wx.Icon
	// logOld is m_logTargetOld: the log target that was active before the log
	// pane was installed. The close handler restores it and destroys the pane's
	// logger, which is what MyFrame::~MyFrame() does.
	logOld unsafe.Pointer

	book    book
	bookTyp bookType
	images  *wx.BitmapBundleArray
	// imageBundles holds the wxBitmapBundle temporaries pushed into images so
	// they are not freed before the frame is destroyed.
	imageBundles []*wx.BitmapBundle

	orient        int
	chkShowImages bool
	fixedWidth    bool
	multi         bool
	noPageTheme   bool
	buttonBar     bool
	horzLayout    bool

	menuType    *wx.Menu
	menuOrient  *wx.Menu
	menuStyle   *wx.Menu
	menuFile    *wx.Menu
	menuPageOps *wx.Menu

	// Function-level "static" counters from the C++ handlers.
	pageAdded          uint
	pageAddedNoSelect  uint
	pageInserted       uint
	subPageAdded       uint
	subPageAddedBefore uint
	eventNum           int

	// Idle-time title bookkeeping (static variables in OnIdle).
	titlePages int
	titleSel   int
	titleBook  book
}

func main() {
	wx.Run(func() bool {
		f := newSampleFrame()
		if f == nil {
			return false
		}
		f.frame.Show()
		wx.SetTopWindow(f.frame.Ptr())
		return true
	})
}

// newSampleFrame builds the frame, its menus, the log pane and the first book
// control -- MyFrame::MyFrame() plus CreateInitialPages().
func newSampleFrame() *sampleFrame {
	f := &sampleFrame{
		bookTyp:       typeNotebook,
		orient:        idOrientDefault,
		chkShowImages: true,
		titlePages:    wx.NOT_FOUND,
		titleSel:      wx.NOT_FOUND,
	}

	frame := wx.NewFrame(nil, wx.ID_ANY, "wxWidgets book controls sample",
		-1, -1, -1, -1, wx.DEFAULT_FRAME_STYLE)
	if frame == nil {
		return nil
	}
	f.frame = frame

	frame.SetExtraStyle(wx.FRAME_EX_CONTEXTHELP)

	// wxHelpProvider::Set(new wxSimpleHelpProvider)
	if hp := wx.NewSimpleHelpProvider(); hp != nil {
		hp.Set()
	}

	if icon := wx.NewIconFromXPM(sampleXPM); icon != nil {
		frame.SetIcon(icon.Ptr())
		f.icon = icon
	}

	f.buildMenus()

	// use some random images for the book control pages
	// (m_images.push_back(wxArtProvider::GetBitmapBundle(...)) x4, 32x32)
	f.images = wx.NewBitmapBundleArray()
	if f.images != nil {
		artIDs := []string{wx.ART_INFORMATION, wx.ART_QUESTION, wx.ART_WARNING, wx.ART_ERROR}
		for _, artID := range artIDs {
			bundle := wx.ArtProvider_GetBitmapBundle(artID, wx.ART_OTHER, 32, 32)
			if bundle == nil {
				continue
			}
			f.images.Add(bundle.Ptr())
			// Keep the bundle alive and undeleted for the life of the frame:
			// m_images is a member of MyFrame in the C++ original.
			f.imageBundles = append(f.imageBundles, bundle)
		}
	}

	panel := wx.NewPanel(frame.Ptr(), wx.ID_ANY, -1, -1, -1, -1, 0)
	if panel == nil {
		return nil
	}
	f.panel = panel

	// Log window that the sample routes wxLogMessage() into.
	f.logText = wx.NewTextCtrl(panel.Ptr(), wx.ID_ANY, "",
		-1, -1, -1, -1, wx.TE_MULTILINE|wx.TE_READONLY)
	if f.logText != nil {
		if lt := wx.NewLogTextCtrl(f.logText.Ptr()); lt != nil {
			f.logCtrl = lt
			f.logOld = wx.SetActiveLogTarget(lt.Ptr())
		}
	}

	f.sizerFrame = wx.NewBoxSizer(wx.VERTICAL)
	if f.logText != nil {
		wxf := wx.NewSizerFlags(1)
		wxf.Expand()
		f.sizerFrame.AddWindowFlags(f.logText.Ptr(), wxf.Ptr())
		wxf.Delete()
	}

	f.recreateBook()

	panel.SetSizer(f.sizerFrame.Ptr(), true)
	panel.Layout()

	sizer := wx.NewBoxSizer(wx.VERTICAL)
	if sizer != nil {
		wxf := wx.NewSizerFlags(1)
		wxf.Expand()
		sizer.AddWindowFlags(panel.Ptr(), wxf.Ptr())
		wxf.Delete()
		frame.SetSizerAndFit(sizer.Ptr(), true)
	}

	f.bindEvents()
	return f
}

// buildMenus recreates the File / Pages / Operations menu tree verbatim.
func (f *sampleFrame) buildMenus() {
	f.menuType = wx.NewMenu("", 0)
	if f.menuType == nil {
		return
	}
	f.menuType.AppendRadioItem(idBookNotebook, "&Notebook\tCtrl-1", "")
	f.menuType.AppendRadioItem(idBookListbook, "&Listbook\tCtrl-2", "")
	f.menuType.AppendRadioItem(idBookChoicebook, "&Choicebook\tCtrl-3", "")
	f.menuType.AppendRadioItem(idBookTreebook, "&Treebook\tCtrl-4", "")
	f.menuType.AppendRadioItem(idBookToolbook, "T&oolbook\tCtrl-5", "")
	f.menuType.AppendRadioItem(idBookAuiNotebook, "&AuiNotebook\tCtrl-6", "")
	f.menuType.AppendRadioItem(idBookSimplebook, "&Simple book\tCtrl-7", "")
	f.menuType.Check(idBookNotebook+int(f.bookTyp), true)

	f.menuOrient = wx.NewMenu("", 0)
	if f.menuOrient == nil {
		return
	}
	f.menuOrient.AppendRadioItem(idOrientDefault, "&Default\tAlt-0", "")
	f.menuOrient.AppendRadioItem(idOrientTop, "&Top\tAlt-1", "")
	f.menuOrient.AppendRadioItem(idOrientBottom, "&Bottom\tAlt-2", "")
	f.menuOrient.AppendRadioItem(idOrientLeft, "&Left\tAlt-3", "")
	f.menuOrient.AppendRadioItem(idOrientRight, "&Right\tAlt-4", "")

	f.menuStyle = wx.NewMenu("", 0)
	if f.menuStyle == nil {
		return
	}
	f.menuStyle.AppendCheckItem(idFixedWidth, "&Fixed Width (wxNotebook)", "")
	f.menuStyle.AppendCheckItem(idMulti, "&Multiple lines (wxNotebook)", "")
	f.menuStyle.AppendCheckItem(idNoPageTheme, "&No Page Theme (wxNotebook)", "")
	f.menuStyle.AppendCheckItem(idButtonBar, "&Button Bar (wxToolbook)", "")
	f.menuStyle.AppendCheckItem(idHorzLayout, "&Horizontal layout (wxToolbook)", "")

	f.menuPageOps = wx.NewMenu("", 0)
	if f.menuPageOps == nil {
		return
	}
	f.menuPageOps.Append(idAddPage, "&Add page\tAlt-A", "", false)
	f.menuPageOps.Append(idAddPageNoSelect, "&Add page (don't select)\tAlt-B", "", false)
	f.menuPageOps.Append(idInsertPage, "&Insert page\tAlt-I", "", false)
	f.menuPageOps.Append(idDeleteCurPage, "&Delete current page\tAlt-D", "", false)
	f.menuPageOps.Append(idDeleteLastPage, "D&elete last page\tAlt-L", "", false)
	f.menuPageOps.Append(idNextPage, "&Next page\tAlt-N", "", false)
	f.menuPageOps.AppendSeparator()
	f.menuPageOps.Append(idAddPageBefore, "Insert page &before\tAlt-B", "", false)
	f.menuPageOps.Append(idAddSubPage, "Add s&ub page\tAlt-U", "", false)
	f.menuPageOps.AppendSeparator()
	f.menuPageOps.Append(idChangeSelection, "&Change selection to 0\tCtrl-0", "", false)
	f.menuPageOps.Append(idSetSelection, "&Set selection to 0\tShift-Ctrl-0", "", false)
	f.menuPageOps.AppendSeparator()
	f.menuPageOps.Append(idGetPageSize, "Sho&w page size", "", false)
	f.menuPageOps.Append(idSetPageSize, "Set &page size", "", false)

	menuOperations := wx.NewMenu("", 0)
	if menuOperations == nil {
		return
	}
	menuOperations.Append(idContextHelp, "&Context help\tCtrl-F1", "", false)
	menuOperations.Append(idHitTest, "&Hit test\tCtrl-H", "", false)

	f.menuFile = wx.NewMenu("", 0)
	if f.menuFile == nil {
		return
	}
	f.menuFile.AppendSub(wx.ID_ANY, "&Type", f.menuType.Ptr(), "Type of control")
	f.menuFile.AppendSub(wx.ID_ANY, "&Orientation", f.menuOrient.Ptr(), "Orientation of control")
	f.menuFile.AppendCheckItem(idShowImages, "&Show images\tAlt-S", "")
	f.menuFile.AppendSub(wx.ID_ANY, "&Style", f.menuStyle.Ptr(), "Style of control")
	f.menuFile.AppendSeparator()
	f.menuFile.Append(wx.ID_EXIT, "E&xit", "Quits the application", false)
	f.menuFile.Check(idShowImages, f.chkShowImages)

	menuBar := wx.NewMenuBar(0)
	if menuBar == nil {
		return
	}
	menuBar.Append(f.menuFile.Ptr(), "&File")
	menuBar.Append(f.menuPageOps.Ptr(), "&Pages")
	menuBar.Append(menuOperations.Ptr(), "&Operations")
	f.frame.SetMenuBar(menuBar.Ptr())
}

// bindEvents wires the C++ event table onto the frame.
func (f *sampleFrame) bindEvents() {
	// File menu
	f.frame.BindIdRange(wx.EVT_MENU, idBookNotebook, idBookMax, f.onType)
	f.frame.BindIdRange(wx.EVT_MENU, idOrientDefault, idOrientMax, f.onOrient)
	f.frame.BindId(wx.EVT_MENU, idShowImages, f.onShowImages)
	f.frame.BindIdRange(wx.EVT_MENU, idFixedWidth, idHorzLayout, f.onStyle)
	f.frame.BindId(wx.EVT_MENU, wx.ID_EXIT, f.onExit)

	// Pages menu
	f.frame.BindId(wx.EVT_MENU, idAddPage, f.onAddPage)
	f.frame.BindId(wx.EVT_MENU, idAddPageNoSelect, f.onAddPageNoSelect)
	f.frame.BindId(wx.EVT_MENU, idInsertPage, f.onInsertPage)
	f.frame.BindId(wx.EVT_MENU, idDeleteCurPage, f.onDeleteCurPage)
	f.frame.BindId(wx.EVT_MENU, idDeleteLastPage, f.onDeleteLastPage)
	f.frame.BindId(wx.EVT_MENU, idNextPage, f.onNextPage)
	f.frame.BindId(wx.EVT_MENU, idAddPageBefore, f.onAddPageBefore)
	f.frame.BindId(wx.EVT_MENU, idAddSubPage, f.onAddSubPage)
	f.frame.BindId(wx.EVT_MENU, idChangeSelection, f.onChangeSelection)
	f.frame.BindId(wx.EVT_MENU, idSetSelection, f.onSetSelection)
	f.frame.BindId(wx.EVT_MENU, idGetPageSize, f.onGetPageSize)
	f.frame.BindId(wx.EVT_MENU, idSetPageSize, f.onSetPageSize)

	// Operations menu
	f.frame.BindId(wx.EVT_MENU, idContextHelp, f.onContextHelp)
	f.frame.BindId(wx.EVT_MENU, idHitTest, f.onHitTest)

	// Book controls: the C++ event table connects every page event type on the
	// frame with wxID_ANY.
	f.frame.Bind(wx.EVT_NOTEBOOK_PAGE_CHANGED, f.bookHandler("wxNotebook", false))
	f.frame.Bind(wx.EVT_NOTEBOOK_PAGE_CHANGING, f.bookHandler("wxNotebook", true))
	f.frame.Bind(wx.EVT_LISTBOOK_PAGE_CHANGED, f.bookHandler("wxListbook", false))
	f.frame.Bind(wx.EVT_LISTBOOK_PAGE_CHANGING, f.bookHandler("wxListbook", true))
	f.frame.Bind(wx.EVT_CHOICEBOOK_PAGE_CHANGED, f.bookHandler("wxChoicebook", false))
	f.frame.Bind(wx.EVT_CHOICEBOOK_PAGE_CHANGING, f.bookHandler("wxChoicebook", true))
	f.frame.Bind(wx.EVT_TREEBOOK_PAGE_CHANGED, f.bookHandler("wxTreebook", false))
	f.frame.Bind(wx.EVT_TREEBOOK_PAGE_CHANGING, f.bookHandler("wxTreebook", true))
	f.frame.Bind(wx.EVT_TOOLBOOK_PAGE_CHANGED, f.bookHandler("wxToolbook", false))
	f.frame.Bind(wx.EVT_TOOLBOOK_PAGE_CHANGING, f.bookHandler("wxToolbook", true))
	f.frame.Bind(wx.EVT_AUINOTEBOOK_PAGE_CHANGED, f.bookHandler("wxAuiNotebook", false))
	f.frame.Bind(wx.EVT_AUINOTEBOOK_PAGE_CHANGING, f.bookHandler("wxAuiNotebook", true))

	// EVT_UPDATE_UI_RANGE(ID_ADD_PAGE_BEFORE, ID_ADD_SUB_PAGE, OnUpdateTreeMenu)
	f.frame.BindIdRange(wx.EVT_UPDATE_UI, idAddPageBefore, idAddSubPage, f.onUpdateTreeMenu)

	// Update title in idle time
	f.frame.Bind(wx.EVT_IDLE, f.onIdle)

	f.frame.Bind(wx.EVT_CLOSE_WINDOW, func(evt *wx.Event) {
		if evt == nil {
			return
		}
		// MyFrame::~MyFrame(): delete wxLog::SetActiveTarget(m_logTargetOld).
		wx.SetActiveLogTarget(f.logOld)
		if f.logCtrl != nil {
			f.logCtrl.Delete()
			f.logCtrl = nil
		}
		// MyFrame::~MyFrame() also destroys the page image array and the art
		// bundles held in it. wxBitmapBundle is reference counted, so the book
		// control's own copies stay valid until the frame is destroyed.
		if f.images != nil {
			f.images.Delete()
			f.images = nil
		}
		for _, bundle := range f.imageBundles {
			bundle.Delete()
		}
		f.imageBundles = nil
		if f.icon != nil {
			f.icon.Delete()
			f.icon = nil
		}
		wx.ExitMainLoop()
	})
}

// ---------------------------------------------------------------------------
// Event handlers
// ---------------------------------------------------------------------------

func (f *sampleFrame) onType(evt *wx.Event) {
	if evt == nil {
		return
	}
	f.bookTyp = bookType(evt.GetId() - idBookNotebook)

	if f.book != nil {
		f.sizerFrame.HideWindow(f.book.Ptr())
	}
	f.recreateBook()
}

// onUpdateTreeMenu is OnUpdateTreeMenu(): the treebook-only commands are only
// enabled while a treebook is the current book control.
func (f *sampleFrame) onUpdateTreeMenu(evt *wx.Event) {
	if evt == nil {
		return
	}
	updateEvt := wx.AsUpdateUIEvent(evt)
	if updateEvt == nil {
		return
	}
	updateEvt.Enable(f.bookTyp == typeTreebook)
}

func (f *sampleFrame) onOrient(evt *wx.Event) {
	if evt == nil {
		return
	}
	f.orient = evt.GetId()
	f.recreateBook()
	f.sizerFrame.Layout()
}

func (f *sampleFrame) onShowImages(evt *wx.Event) {
	cmdEvt := wx.AsCommandEvent(evt)
	if cmdEvt == nil {
		return
	}
	f.chkShowImages = cmdEvt.IsChecked()
	f.recreateBook()
	f.sizerFrame.Layout()
}

func (f *sampleFrame) onStyle(evt *wx.Event) {
	cmdEvt := wx.AsCommandEvent(evt)
	if cmdEvt == nil {
		return
	}
	checked := cmdEvt.IsChecked()

	id := evt.GetId()
	if id == idFixedWidth {
		f.fixedWidth = checked
	} else if id == idMulti {
		f.multi = checked
	} else if id == idNoPageTheme {
		f.noPageTheme = checked
	} else if id == idButtonBar {
		f.buttonBar = checked
	} else if id == idHorzLayout {
		f.horzLayout = checked
	}

	f.recreateBook()
	f.sizerFrame.Layout()
}

func (f *sampleFrame) onExit(evt *wx.Event) {
	if evt == nil {
		return
	}
	// MyFrame::OnExit() calls Close(); force stays false so the close event is
	// generated and the log target is restored before the loop ends.
	f.frame.Close(false)
}

func (f *sampleFrame) onAddPage(evt *wx.Event) {
	if evt == nil || f.book == nil {
		return
	}
	f.pageAdded++
	f.book.AddPage(createNewPage(f.book.Ptr()),
		fmt.Sprintf("%s%d", addedPageName, f.pageAdded), true, f.iconIndex())
}

func (f *sampleFrame) onAddPageNoSelect(evt *wx.Event) {
	if evt == nil || f.book == nil {
		return
	}
	f.pageAddedNoSelect++
	f.book.AddPage(createNewPage(f.book.Ptr()),
		fmt.Sprintf("%s%d", addedPageName, f.pageAddedNoSelect), false, f.iconIndex())
}

func (f *sampleFrame) onInsertPage(evt *wx.Event) {
	if evt == nil || f.book == nil {
		return
	}
	f.pageInserted++
	f.book.InsertPage(0, createUserCreatedPage(f.book.Ptr()),
		fmt.Sprintf("%s%d", insertedPageName, f.pageInserted), false, f.iconIndex())
	f.book.SetSelection(0)
}

func (f *sampleFrame) onDeleteCurPage(evt *wx.Event) {
	if evt == nil || f.book == nil {
		return
	}
	if sel := f.book.GetSelection(); sel != wx.NOT_FOUND {
		f.book.DeletePage(sel)
	}
}

func (f *sampleFrame) onDeleteLastPage(evt *wx.Event) {
	if evt == nil || f.book == nil {
		return
	}
	if n := f.book.GetPageCount(); n != 0 {
		f.book.DeletePage(n - 1)
	}
}

func (f *sampleFrame) onNextPage(evt *wx.Event) {
	if evt == nil || f.book == nil {
		return
	}
	f.book.AdvanceSelection(true)
}

func (f *sampleFrame) onChangeSelection(evt *wx.Event) {
	if evt == nil || f.book == nil {
		return
	}
	f.book.ChangeSelection(0)
}

func (f *sampleFrame) onSetSelection(evt *wx.Event) {
	if evt == nil || f.book == nil {
		return
	}
	f.book.SetSelection(0)
}

func (f *sampleFrame) onGetPageSize(evt *wx.Event) {
	if evt == nil || f.book == nil {
		return
	}
	page := wx.WindowFromPtr(f.book.GetPage(0))
	if page == nil {
		return
	}
	sizePage := page.GetSize()
	sizeBook := f.book.GetSize()

	wx.LogMessage(fmt.Sprintf("Page size is (%d, %d), book size (%d, %d)",
		sizePage.Width, sizePage.Height, sizeBook.Width, sizeBook.Height))
}

func (f *sampleFrame) onSetPageSize(evt *wx.Event) {
	if evt == nil || f.book == nil {
		return
	}
	f.book.SetPageSize(300, 300)
	wx.LogMessage("Page size set to (300, 300)")
}

func (f *sampleFrame) onAddSubPage(evt *wx.Event) {
	if evt == nil {
		return
	}
	tbk, ok := f.book.(*wx.Treebook)
	if !ok {
		return
	}
	sel := tbk.GetSelection()
	if sel == wx.NOT_FOUND {
		wx.LogError("Select the parent page first!")
		return
	}
	f.subPageAdded++
	tbk.InsertSubPage(sel, createNewPage(tbk.Ptr()),
		fmt.Sprintf("%s%d", addedSubPageName, f.subPageAdded), true, f.iconIndex())
}

func (f *sampleFrame) onAddPageBefore(evt *wx.Event) {
	if evt == nil || f.book == nil {
		return
	}
	sel := f.book.GetSelection()
	if sel == wx.NOT_FOUND {
		wx.LogError("Select the parent page first!")
		return
	}
	f.subPageAddedBefore++
	f.book.InsertPage(sel, createNewPage(f.book.Ptr()),
		fmt.Sprintf("%s%d", addedPageNameBefore, f.subPageAddedBefore), true, f.iconIndex())
}

func (f *sampleFrame) onContextHelp(evt *wx.Event) {
	if evt == nil {
		return
	}
	// wxContextHelp ch(this) runs a local event loop; the object is deleted
	// afterwards, mirroring the C++ stack object's destructor.
	if ch := wx.NewContextHelp(f.frame.Ptr(), true); ch != nil {
		ch.Delete()
	}
}

func (f *sampleFrame) onHitTest(evt *wx.Event) {
	if evt == nil || f.book == nil {
		return
	}
	pt := wx.GetMousePositionPoint()
	clientPt := f.book.ScreenToClient(pt.X, pt.Y)

	var hitFlags int32
	pagePos := f.book.HitTestFlags(clientPt.X, clientPt.Y, unsafe.Pointer(&hitFlags))
	flags := int(hitFlags)

	flagsStr := ""
	flagsStr = addFlagStrIfFlagPresent(flagsStr, flags, wx.BK_HITTEST_NOWHERE, "wxBK_HITTEST_NOWHERE")
	flagsStr = addFlagStrIfFlagPresent(flagsStr, flags, wx.BK_HITTEST_ONICON, "wxBK_HITTEST_ONICON")
	flagsStr = addFlagStrIfFlagPresent(flagsStr, flags, wx.BK_HITTEST_ONLABEL, "wxBK_HITTEST_ONLABEL")
	flagsStr = addFlagStrIfFlagPresent(flagsStr, flags, wx.BK_HITTEST_ONITEM, "wxBK_HITTEST_ONITEM")
	flagsStr = addFlagStrIfFlagPresent(flagsStr, flags, wx.BK_HITTEST_ONPAGE, "wxBK_HITTEST_ONPAGE")

	wx.LogMessage(fmt.Sprintf("HitTest at (%d,%d): %d: %s", pt.X, pt.Y, pagePos, flagsStr))
}

// addFlagStrIfFlagPresent is MyFrame::AddFlagStrIfFlagPresent().
func addFlagStrIfFlagPresent(flagStr string, flags int, flag int, flagName string) string {
	if flags&flag != flag {
		return flagStr
	}
	if flagStr != "" {
		flagStr += " | "
	}
	return flagStr + flagName
}

// onBookCtrl is MyFrame::OnBookCtrl(): log the page change and, for a Changing
// event coming from the "Veto" page, ask whether the user really wants to
// leave it.
func (f *sampleFrame) onBookCtrl(evt *wx.Event, controlName string, changing bool) {
	if evt == nil || f.book == nil || f.logText == nil {
		return
	}
	bk := wx.AsBookCtrlEvent(evt)
	if bk == nil {
		return
	}

	// notebook.cpp uses static_cast<wxBookCtrlBase*>(event.GetEventObject())
	// rather than the current book, so a stale event still reads its own book.
	book := f.book
	if obj := evt.GetEventObject(); obj != nil {
		if base := wx.BookCtrlBaseFromPtr(obj); base != nil {
			book = base
		}
	}

	nameEvent := "Changed"
	veto := ""
	if changing {
		nameEvent = "Changing"
		idx := bk.GetOldSelection()
		if idx != wx.NOT_FOUND && book.GetPageText(idx) == vetoPageName {
			answer := wx.MessageBox("Are you sure you want to leave this page?\n(This demonstrates veto-ing)",
				"Notebook sample", wx.ICON_QUESTION|wx.YES_NO, f.frame.Ptr(), -1, -1)
			if answer != wx.YES {
				bk.Veto()
				veto = " (vetoed)"
			}
		}
	}

	f.eventNum++
	wx.LogMessage(fmt.Sprintf("Event #%d: %s: %s (%d) new sel %d, old %d, current %d%s",
		f.eventNum, controlName, nameEvent, evt.GetEventType(),
		bk.GetSelection(), bk.GetOldSelection(), book.GetSelection(), veto))
	f.logText.SetInsertionPointEnd()
}

func (f *sampleFrame) onIdle(evt *wx.Event) {
	if evt == nil {
		return
	}
	nPages := 0
	nSel := wx.NOT_FOUND
	if f.book != nil {
		nPages = f.book.GetPageCount()
		nSel = f.book.GetSelection()
	}

	if nPages != f.titlePages || nSel != f.titleSel || f.titleBook != f.book {
		f.titlePages = nPages
		f.titleSel = nSel
		f.titleBook = f.book

		selection := "not found"
		if nSel != wx.NOT_FOUND {
			selection = strconv.Itoa(nSel)
		}
		f.frame.SetTitle(fmt.Sprintf("Notebook and friends (%d pages, selection: %s)",
			nPages, selection))
	}
}
