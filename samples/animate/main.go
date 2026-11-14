// Command animate is the kwxGO port of the wxWidgets "animate" sample
// (samples/animate/anitest.cpp).
//
// It shows a wxAnimationCtrl playing throbber.gif (throbber_2x.gif on high DPI
// displays), and lets the user inspect the individual frames of the animation:
// the frame number spin control selects a frame, and the static texts and
// bitmap below it show that frame's delay and image. The Animation menu can
// play/stop the animation, replace it with an empty one, install the standard
// "missing image" bitmap as the inactive bitmap, switch auto-resizing off and
// change the control's background colour. File/Open loads an arbitrary
// animation file, and the frame's update-UI handler enables or disables the
// Play, Stop and "Set no autoresize" items depending on the playing state.
package main

import (
	"fmt"
	"unsafe"

	"github.com/KeyWorksRW/kwxGO/wx"
)

// Command IDs mirroring the anonymous enum in anitest.cpp. ID_USE_GENERIC (6)
// only exists when wxHAS_NATIVE_ANIMATIONCTRL is defined, which this build
// does not do: wxAnimationCtrl is the generic control here, so -- exactly as
// when the C++ sample is built for MSW -- the "Use generic animation" item and
// its handler are compiled out.
const (
	idPlay              = 1
	idSetNullAnimation  = 2
	idSetInactiveBitmap = 3
	idSetNoAutoResize   = 4
	idSetBgColor        = 5
)

const (
	// Dimensions of the original sample.xpm icon.
	sampleIconWidth  = 32
	sampleIconHeight = 32
)

// animateFrame holds the state MyFrame keeps in member variables.
type animateFrame struct {
	frame         *wx.Frame
	animationCtrl *wx.AnimationCtrl

	// Controls used to show the individual animation frames.
	frameBox    *wx.StaticBox
	frameSpin   *wx.SpinCtrl
	frameCount  *wx.StaticText
	frameDelay  *wx.StaticText
	frameBitmap *wx.StaticBitmap
}

func main() {
	wx.Run(func() bool {
		f := newAnimateFrame()
		if f == nil {
			return false
		}
		f.frame.Show()
		wx.SetTopWindow(f.frame.Ptr())
		return true
	})
}

// newAnimateFrame builds the frame, its menu bar, status bar and the sizer with
// the animation control and the frame-inspection controls -- MyFrame's
// constructor.
func newAnimateFrame() *animateFrame {
	frame := wx.NewFrame(nil, wx.ID_ANY, "Animation Demo",
		-1, -1, -1, -1, wx.DEFAULT_FRAME_STYLE)
	if frame == nil {
		return nil
	}

	f := &animateFrame{frame: frame}

	// SetIcon(wxICON(sample)): the XPM icon has been converted to sample.svg.
	if icon := loadSampleIcon(); icon != nil {
		frame.SetIcon(icon.Ptr())
		icon.Delete()
	}

	f.buildMenus()
	// CreateStatusBar(number = 1, style = wxSTB_DEFAULT_STYLE, winid = 0)
	frame.CreateStatusBar(1, wx.STB_DEFAULT_STYLE, 0)

	// Use a wxBoxSizer, otherwise wxFrame would automatically resize the
	// animation control to fill its client area on user resizes.
	sz := wx.NewBoxSizer(wx.VERTICAL)
	if sz == nil {
		return nil
	}

	label := wx.NewStaticText(frame.Ptr(), wx.ID_ANY, "wxAnimationCtrl:",
		-1, -1, -1, -1, 0)
	if label == nil {
		return nil
	}
	addWindow(sz, label.Ptr(), centreBorder)

	f.animationCtrl = wx.NewAnimationCtrl(frame.Ptr(), wx.ID_ANY,
		-1, -1, -1, -1, wx.AC_DEFAULT_STYLE)
	if f.animationCtrl == nil {
		return nil
	}

	// wxAnimationBundle animations;
	// wxAnimation throbber("throbber.gif");
	// wxAnimation throbber2x("throbber_2x.gif");
	// if (animations.IsOk()) { SetAnimation(animations); Play(); }
	if animations := newThrobberAnimationBundle(); animations != nil {
		if animations.IsOk() {
			f.animationCtrl.SetAnimationBundle(animations.Ptr())
			f.animationCtrl.Play()
		}
		animations.Delete()
	}

	addWindow(sz, f.animationCtrl.Ptr(), centreBorder)

	// Controls allowing the user to examine the individual animation frames.
	f.frameBox = wx.NewStaticBox(frame.Ptr(), wx.ID_ANY, "Animation frames",
		-1, -1, -1, -1, 0)
	if f.frameBox == nil {
		return nil
	}

	frameSizer := wx.NewStaticBoxSizer(f.frameBox.Ptr(), wx.VERTICAL)
	if frameSizer == nil {
		return nil
	}

	spinSizer := wx.NewBoxSizer(wx.HORIZONTAL)
	if spinSizer == nil {
		return nil
	}

	spinLabel := wx.NewStaticText(f.frameBox.Ptr(), wx.ID_ANY, "&Frame:",
		-1, -1, -1, -1, 0)
	if spinLabel == nil {
		return nil
	}
	addWindow(spinSizer, spinLabel.Ptr(), func(sf *wx.SizerFlags) {
		sf.CenterVertical()
		sf.Border(wx.RIGHT)
	})

	f.frameSpin = wx.NewSpinCtrl(f.frameBox.Ptr(), wx.ID_ANY, "",
		-1, -1, -1, -1, wx.SP_ARROW_KEYS, 0, 100, 0)
	if f.frameSpin == nil {
		return nil
	}
	// m_frameSpin->Bind(wxEVT_SPINCTRL, &MyFrame::OnFrameSpin, this)
	f.frameSpin.Bind(wx.EVT_SPINCTRL, f.onFrameSpin)
	addWindow(spinSizer, f.frameSpin.Ptr(), func(sf *wx.SizerFlags) {
		sf.CenterVertical()
	})

	f.frameCount = wx.NewStaticText(f.frameBox.Ptr(), wx.ID_ANY, "",
		-1, -1, -1, -1, 0)
	if f.frameCount == nil {
		return nil
	}
	addWindow(spinSizer, f.frameCount.Ptr(), func(sf *wx.SizerFlags) {
		sf.CenterVertical()
		sf.Border(wx.LEFT)
	})

	addSizer(frameSizer, spinSizer.Ptr(), centreBorder)

	f.frameDelay = wx.NewStaticText(f.frameBox.Ptr(), wx.ID_ANY, "",
		-1, -1, -1, -1, 0)
	if f.frameDelay == nil {
		return nil
	}
	addWindow(frameSizer, f.frameDelay.Ptr(), centreBorder)

	// new wxStaticBitmap(m_frameBox, wxID_ANY, wxBitmapBundle())
	f.frameBitmap = wx.NewStaticBitmap(f.frameBox.Ptr(), wx.ID_ANY,
		wx.Null_Bitmap(), -1, -1, -1, -1, 0)
	if f.frameBitmap == nil {
		return nil
	}
	addWindow(frameSizer, f.frameBitmap.Ptr(), centreBorder)

	addSizer(sz, frameSizer.Ptr(), func(sf *wx.SizerFlags) {
		sf.Expand()
		sf.Border(wx.ALL)
	})

	frame.SetSizer(sz.Ptr(), true)

	f.updateFrameControls()

	frame.SetSize(-1, -1, frame.FromDIP(500), frame.FromDIP(400), wx.SIZE_USE_EXISTING)

	f.bindEvents()
	return f
}

// buildMenus recreates the File / Animation / Help menus of the C++ sample.
func (f *animateFrame) buildMenus() {
	fileMenu := wx.NewMenu("", 0)
	if fileMenu == nil {
		return
	}
	fileMenu.Append(wx.ID_OPEN, "&Open Animation...\tCtrl+O", "Loads an animation", false)
	fileMenu.Append(wx.ID_EXIT, "E&xit", "", false)

	playMenu := wx.NewMenu("", 0)
	if playMenu == nil {
		return
	}
	playMenu.Append(idPlay, "Play\tCtrl+P", "Play the animation", false)
	playMenu.Append(wx.ID_STOP, "Stop\tCtrl+S", "Stop the animation", false)
	playMenu.AppendSeparator()
	playMenu.Append(idSetNullAnimation, "Set null animation",
		"Sets the empty animation in the control", false)
	playMenu.AppendCheckItem(idSetInactiveBitmap, "Set inactive bitmap",
		"Sets an inactive bitmap for the control")
	playMenu.AppendCheckItem(idSetNoAutoResize, "Set no autoresize",
		"Tells the control not to resize automatically")
	playMenu.Append(idSetBgColor, "Set background colour...",
		"Sets the background colour of the control", false)

	helpMenu := wx.NewMenu("", 0)
	if helpMenu == nil {
		return
	}
	helpMenu.Append(wx.ID_ABOUT, "&About", "", false)

	menuBar := wx.NewMenuBar(0)
	if menuBar == nil {
		return
	}
	menuBar.Append(fileMenu.Ptr(), "&File")
	menuBar.Append(playMenu.Ptr(), "&Animation")
	menuBar.Append(helpMenu.Ptr(), "&Help")

	f.frame.SetMenuBar(menuBar.Ptr())
}

// bindEvents wires the C++ event table onto the frame.
func (f *animateFrame) bindEvents() {
	f.frame.BindId(wx.EVT_MENU, idPlay, f.onPlay)
	f.frame.BindId(wx.EVT_MENU, idSetNullAnimation, f.onSetNullAnimation)
	f.frame.BindId(wx.EVT_MENU, idSetInactiveBitmap, f.onSetInactiveBitmap)
	f.frame.BindId(wx.EVT_MENU, idSetNoAutoResize, f.onSetNoAutoResize)
	f.frame.BindId(wx.EVT_MENU, idSetBgColor, f.onSetBgColor)
	f.frame.BindId(wx.EVT_MENU, wx.ID_STOP, f.onStop)
	f.frame.BindId(wx.EVT_MENU, wx.ID_ABOUT, f.onAbout)
	f.frame.BindId(wx.EVT_MENU, wx.ID_EXIT, f.onQuit)
	f.frame.BindId(wx.EVT_MENU, wx.ID_OPEN, f.onOpen)

	// EVT_SIZE(MyFrame::OnSize) resolves to the inherited
	// wxTopLevelWindow::OnSize, which just calls Layout().
	f.frame.Bind(wx.EVT_SIZE, f.onSize)
	f.frame.Bind(wx.EVT_UPDATE_UI, f.onUpdateUI)
}

// ---------------------------------------------------------------------------
// Event handlers
// ---------------------------------------------------------------------------

func (f *animateFrame) onPlay(evt *wx.Event) {
	if evt == nil {
		return
	}
	if !f.animationCtrl.Play() {
		wx.LogError("Invalid animation")
	}
}

func (f *animateFrame) onStop(evt *wx.Event) {
	if evt == nil {
		return
	}
	f.animationCtrl.Stop()
}

func (f *animateFrame) onSetNullAnimation(evt *wx.Event) {
	if evt == nil {
		return
	}
	empty := wx.NewAnimation()
	if empty == nil {
		return
	}
	f.animationCtrl.SetAnimation(empty.Ptr())
	empty.Delete()

	f.updateFrameControls()
}

func (f *animateFrame) onSetInactiveBitmap(evt *wx.Event) {
	if evt == nil {
		return
	}
	cmd := wx.AsCommandEvent(evt)
	if cmd == nil {
		return
	}
	if cmd.IsChecked() {
		// Set a dummy bitmap as the inactive bitmap.
		bmp := wx.ArtProvider_GetBitmap(wx.ART_MISSING_IMAGE, wx.ART_OTHER, -1, -1)
		if bmp != nil {
			f.animationCtrl.SetInactiveBitmap(bmp.Ptr())
			bmp.Delete()
		}
	} else {
		f.animationCtrl.SetInactiveBitmap(wx.Null_Bitmap())
	}
}

func (f *animateFrame) onSetNoAutoResize(evt *wx.Event) {
	if evt == nil {
		return
	}
	cmd := wx.AsCommandEvent(evt)
	if cmd == nil {
		return
	}

	// Recreate the control with the new flag if necessary.
	style := wx.AC_DEFAULT_STYLE
	if cmd.IsChecked() {
		style |= wx.AC_NO_AUTORESIZE
	}

	ctrlWindow := wx.WindowFromPtr(f.animationCtrl.Ptr())
	if ctrlWindow != nil && style != ctrlWindow.GetWindowStyle() {
		f.recreateAnimation(style)
	}
}

func (f *animateFrame) onSetBgColor(evt *wx.Event) {
	if evt == nil {
		return
	}
	ctrlWindow := wx.WindowFromPtr(f.animationCtrl.Ptr())
	if ctrlWindow == nil {
		return
	}

	current := wx.NewColourRGB(0, 0, 0, 255)
	if current == nil {
		return
	}
	defer current.Delete()
	ctrlWindow.GetBackgroundColour(current.Ptr())

	clr := wx.NewColourRGB(0, 0, 0, 255)
	if clr == nil {
		return
	}
	defer clr.Delete()

	wx.GetColourFromUserEx(f.frame.Ptr(), current.Ptr(),
		"Choose the background colour", clr.Ptr())
	if clr.IsOk() {
		ctrlWindow.SetBackgroundColour(clr.Ptr())
	}
}

func (f *animateFrame) onFrameSpin(evt *wx.Event) {
	if evt == nil {
		return
	}
	f.updateFrameInfo()
}

func (f *animateFrame) onQuit(evt *wx.Event) {
	if evt == nil {
		return
	}
	f.frame.Close(false)
}

func (f *animateFrame) onAbout(evt *wx.Event) {
	if evt == nil {
		return
	}
	info := wx.NewAboutDialogInfo()
	if info == nil {
		return
	}
	defer info.Delete()

	info.SetName("wxAnimationCtrl and wxAnimation sample")
	info.SetDescription("This sample program demonstrates the usage of wxAnimationCtrl")
	info.SetCopyright("(C) 2006 Julian Smart")
	info.AddDeveloper("Julian Smart")
	info.AddDeveloper("Guillermo Rodriguez Garcia")
	info.AddDeveloper("Francesco Montorsi")

	wx.AboutBox(info.Ptr(), f.frame.Ptr())
}

func (f *animateFrame) onOpen(evt *wx.Event) {
	if evt == nil {
		return
	}

	dlg := wx.NewFileDialog(f.frame.Ptr(), "Please choose an animation",
		"", "", "*.gif;*.ani;*.webp", -1, -1, wx.FD_OPEN, -1, -1)
	if dlg == nil {
		return
	}
	if dlg.ShowModal() == wx.ID_OK {
		filename := dlg.GetPath()

		// wxAnimation temp(m_animationCtrl->CreateAnimation());
		temp := wx.NewAnimation()
		if temp != nil {
			f.animationCtrl.CreateAnimation(temp.Ptr())
			if !temp.LoadFile(filename, wx.ANIMATION_TYPE_ANY) {
				wx.LogError("Sorry, this animation is not a valid format for wxAnimation.")
			} else {
				f.animationCtrl.SetAnimation(temp.Ptr())
				f.animationCtrl.Play()

				if sizer := wx.SizerFromPtr(f.frame.GetSizer()); sizer != nil {
					sizer.Layout()
				}

				f.updateFrameControls()
			}
			temp.Delete()
		}
	}

	dlg.Destroy()
}

func (f *animateFrame) onSize(evt *wx.Event) {
	if evt == nil {
		return
	}
	f.frame.Layout()
}

func (f *animateFrame) onUpdateUI(evt *wx.Event) {
	if evt == nil {
		return
	}
	menuBar := wx.MenuBarFromPtr(f.frame.GetMenuBar())
	if menuBar == nil {
		return
	}

	playing := f.animationCtrl.IsPlaying()

	// GetMenuBar()->FindItem(id)->Enable(...) in anitest.cpp. The item is
	// addressed by id: wxMenuBar::Enable(id, enable) is exactly
	// FindItem(id)->Enable(enable) inside wxWidgets
	// (wxMenuBarBase::Enable, src/common/menucmn.cpp), while the binding's
	// MenuBar.FindItem allocates a throwaway wxMenu on every call, which leaks
	// an HMENU per update-UI event until CreatePopupMenu() fails.
	menuBar.EnableItem(wx.ID_STOP, playing)
	menuBar.EnableItem(idPlay, !playing)
	menuBar.EnableItem(idSetNoAutoResize, !playing)
}

// ---------------------------------------------------------------------------
// Frame inspection
// ---------------------------------------------------------------------------

// recreateAnimation destroys the animation control and creates a new one with
// the given style, preserving the inactive bitmap and the background colour.
func (f *animateFrame) recreateAnimation(style int) {
	old := f.animationCtrl
	oldWindow := wx.WindowFromPtr(old.Ptr())

	// Save the status of the control before destroying it. The animation
	// itself is intentionally not preserved: in the C++ sample the current
	// animation is only carried over when switching between native and generic
	// controls, and this build always uses the generic one.
	inactive := wx.NewBitmapDefault()
	if inactive == nil {
		return
	}
	defer inactive.Delete()
	old.GetInactiveBitmap(inactive.Ptr())

	bg := wx.NewColourRGB(0, 0, 0, 255)
	if bg == nil {
		return
	}
	defer bg.Delete()
	if oldWindow != nil {
		oldWindow.GetBackgroundColour(bg.Ptr())
	}

	// Destroy and rebuild.
	sizer := wx.SizerFromPtr(f.frame.GetSizer())

	newCtrl := wx.NewAnimationCtrl(f.frame.Ptr(), wx.ID_ANY,
		-1, -1, -1, -1, style)
	if newCtrl == nil {
		return
	}

	if sizer != nil {
		sizer.ReplaceWindow(old.Ptr(), newCtrl.Ptr(), true)
	}

	// Load the old status into the new control.
	newCtrl.SetInactiveBitmap(inactive.Ptr())
	if newWindow := wx.WindowFromPtr(newCtrl.Ptr()); newWindow != nil {
		newWindow.SetBackgroundColour(bg.Ptr())
	}

	if oldWindow != nil {
		oldWindow.Destroy()
	}

	f.animationCtrl = newCtrl

	if sizer != nil {
		sizer.Layout()
	}

	f.updateFrameControls()
}

// updateFrameControls updates the frame controls after the animation changes.
func (f *animateFrame) updateFrameControls() {
	anim := f.currentAnimation()
	count := 0
	if anim != nil && anim.IsOk() {
		count = anim.GetFrameCount()
	}
	if anim != nil {
		defer anim.Delete()
	}

	f.frameBox.Enable(count != 0)

	if count != 0 {
		f.frameSpin.SetRange(0, count-1)
		f.frameCount.SetLabel(fmt.Sprintf("(of %d)", count))
	} else {
		f.frameSpin.SetRange(0, 0)
		f.frameCount.SetLabel("(no frames)")
	}

	f.frameSpin.SetValue(0)

	f.updateFrameInfo()
}

// updateFrameInfo updates the controls showing the currently selected frame.
func (f *animateFrame) updateFrameInfo() {
	anim := f.currentAnimation()
	if anim == nil {
		f.frameDelay.SetLabel("Delay: N/A")
		f.clearFrameBitmap()
		f.frame.Layout()
		return
	}
	defer anim.Delete()

	if anim.IsOk() && anim.GetFrameCount() != 0 {
		frame := f.frameSpin.GetValue()

		delay := anim.GetDelay(frame)
		if delay == -1 {
			f.frameDelay.SetLabel("Delay: forever")
		} else {
			f.frameDelay.SetLabel(fmt.Sprintf("Delay: %d ms", delay))
		}

		image := wx.NewImageDefault()
		if image == nil {
			f.clearFrameBitmap()
			f.frame.Layout()
			return
		}
		anim.GetFrame(frame, image.Ptr())
		if image.IsOk() {
			bitmap := wx.NewBitmapFromImage(image.Ptr(), -1)
			if bitmap != nil {
				f.frameBitmap.SetBitmap(bitmap.Ptr())
				bitmap.Delete()
			}
		} else {
			wx.LogError(fmt.Sprintf("Failed to get animation frame %d.", frame))
			f.clearFrameBitmap()
		}
		image.Delete()
	} else {
		f.frameDelay.SetLabel("Delay: N/A")
		f.clearFrameBitmap()
	}

	f.frame.Layout()
}

// currentAnimation returns a copy of the animation currently set in the control.
func (f *animateFrame) currentAnimation() *wx.Animation {
	anim := wx.NewAnimation()
	if anim == nil {
		return nil
	}
	f.animationCtrl.GetAnimation(anim.Ptr())
	return anim
}

// clearFrameBitmap is SetBitmap(wxBitmapBundle()) on the frame preview.
func (f *animateFrame) clearFrameBitmap() {
	f.frameBitmap.SetBitmap(wx.Null_Bitmap())
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// centreBorder is wxSizerFlags().Centre().Border().
func centreBorder(sf *wx.SizerFlags) {
	sf.Center()
	sf.Border(wx.ALL)
}

// sizerTarget is implemented by every sizer wrapper through its embedded
// wx.Sizer, so box and static-box sizers can both be filled by the helpers
// below.
type sizerTarget interface {
	AddWindowFlags(window unsafe.Pointer, flags unsafe.Pointer)
	AddSizerFlags(sizer unsafe.Pointer, flags unsafe.Pointer)
}

// addWindow adds a window to the sizer with the given wxSizerFlags modifiers.
func addWindow(s sizerTarget, window unsafe.Pointer, apply func(*wx.SizerFlags)) {
	flags := wx.NewSizerFlags(0)
	if flags == nil {
		return
	}
	apply(flags)
	s.AddWindowFlags(window, flags.Ptr())
	flags.Delete()
}

// addSizer adds a nested sizer with the given wxSizerFlags modifiers.
func addSizer(s sizerTarget, child unsafe.Pointer, apply func(*wx.SizerFlags)) {
	flags := wx.NewSizerFlags(0)
	if flags == nil {
		return
	}
	apply(flags)
	s.AddSizerFlags(child, flags.Ptr())
	flags.Delete()
}

// newThrobberAnimationBundle builds the wxAnimationBundle of the C++
// constructor: throbber.gif and, when it loads, throbber_2x.gif. The bundle
// itself picks the animation matching the display's DPI scale factor.
func newThrobberAnimationBundle() *wx.AnimationBundle {
	animations := wx.NewAnimationBundle()
	if animations == nil {
		return nil
	}

	for _, name := range []string{"throbber.gif", "throbber_2x.gif"} {
		anim := wx.NewAnimationFromFile(name, wx.ANIMATION_TYPE_ANY)
		if anim == nil {
			continue
		}
		if anim.IsOk() {
			animations.Add(anim.Ptr())
		}
		anim.Delete()
	}

	return animations
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
