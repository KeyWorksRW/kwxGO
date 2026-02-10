package wx

import (
	"testing"
)

// TestEventCallbackRegistry verifies the callback registry works correctly.
func TestEventCallbackRegistry(t *testing.T) {
	called := false
	handler := EventHandler(func(evt *Event) {
		called = true
	})

	// Register
	id := registerCallback(handler)
	if id == 0 {
		t.Fatal("registerCallback returned 0")
	}

	// Verify registered
	callbackMutex.RLock()
	_, ok := callbackRegistry[id]
	callbackMutex.RUnlock()
	if !ok {
		t.Fatal("handler not found in registry after registration")
	}

	// Invoke the handler with a nil-ptr event (won't dereference)
	callbackRegistry[id](&Event{})
	if !called {
		t.Fatal("handler was not called")
	}

	// Unregister
	unregisterCallback(id)
	callbackMutex.RLock()
	_, ok = callbackRegistry[id]
	callbackMutex.RUnlock()
	if ok {
		t.Fatal("handler still in registry after unregisterCallback")
	}
}

// TestEventCallbackUniqueIDs verifies that each registration gets a unique ID.
func TestEventCallbackUniqueIDs(t *testing.T) {
	ids := make(map[uintptr]bool)
	handlers := make([]uintptr, 100)

	for i := 0; i < 100; i++ {
		id := registerCallback(func(evt *Event) {})
		if ids[id] {
			t.Fatalf("duplicate callback ID: %d", id)
		}
		ids[id] = true
		handlers[i] = id
	}

	// Cleanup
	for _, id := range handlers {
		unregisterCallback(id)
	}
}

// TestEventHandlerType verifies EventHandler type signature.
func TestEventHandlerType(t *testing.T) {
	// EventHandler is func(*Event)
	var handler EventHandler = func(evt *Event) {
		_ = evt.Ptr()
		_ = evt.GetEventType()
		_ = evt.GetId()
		_ = evt.GetTimestamp()
		_ = evt.GetEventObject()
		_ = evt.IsCommandEvent()
		_ = evt.GetSkipped()
		evt.Skip()
	}
	_ = handler
}

// TestEventType verifies Event methods exist and compile.
func TestEventType(t *testing.T) {
	e := &Event{}
	_ = e.Ptr
	_ = e.GetEventType
	_ = e.GetId
	_ = e.GetTimestamp
	_ = e.GetEventObject
	_ = e.IsCommandEvent
	_ = e.GetSkipped
	_ = e.Skip
}

// TestCommandEventType verifies CommandEvent methods exist and compile.
func TestCommandEventType(t *testing.T) {
	e := &Event{}
	cmd := AsCommandEvent(e)
	_ = cmd.GetSelection
	_ = cmd.GetString
	_ = cmd.GetInt
	_ = cmd.GetExtraLong
	_ = cmd.IsChecked
	_ = cmd.IsSelection
	// Base event methods are inherited
	_ = cmd.GetEventType
	_ = cmd.GetId
	_ = cmd.Skip
}

// TestCloseEventType verifies CloseEvent methods exist and compile.
func TestCloseEventType(t *testing.T) {
	e := &Event{}
	closeEvt := AsCloseEvent(e)
	_ = closeEvt.CanVeto
	_ = closeEvt.GetVeto
	_ = closeEvt.Veto
	_ = closeEvt.GetLoggingOff
	// Base event methods are inherited
	_ = closeEvt.GetEventType
	_ = closeEvt.Skip
}

// TestKeyEventType verifies KeyEvent methods exist and compile.
func TestKeyEventType(t *testing.T) {
	e := &Event{}
	key := AsKeyEvent(e)
	_ = key.GetKeyCode
	_ = key.GetX
	_ = key.GetY
	_ = key.GetPosition
	_ = key.AltDown
	_ = key.ControlDown
	_ = key.ShiftDown
	_ = key.MetaDown
	_ = key.GetModifiers
	_ = key.HasModifiers
	// Base event methods are inherited
	_ = key.GetEventType
	_ = key.Skip
}

// TestMouseEventType verifies MouseEvent methods exist and compile.
func TestMouseEventType(t *testing.T) {
	e := &Event{}
	mouse := AsMouseEvent(e)
	_ = mouse.GetX
	_ = mouse.GetY
	_ = mouse.GetPosition
	_ = mouse.AltDown
	_ = mouse.ControlDown
	_ = mouse.ShiftDown
	_ = mouse.MetaDown
	_ = mouse.LeftDown
	_ = mouse.LeftUp
	_ = mouse.LeftDClick
	_ = mouse.LeftIsDown
	_ = mouse.MiddleDown
	_ = mouse.MiddleUp
	_ = mouse.MiddleIsDown
	_ = mouse.RightDown
	_ = mouse.RightUp
	_ = mouse.Dragging
	_ = mouse.Moving
	_ = mouse.Entering
	_ = mouse.Leaving
	_ = mouse.IsButton
	_ = mouse.Button
	_ = mouse.ButtonDown
	_ = mouse.ButtonUp
	_ = mouse.ButtonDClick
	_ = mouse.ButtonIsDown
	// Base event methods are inherited
	_ = mouse.GetEventType
	_ = mouse.Skip
}

// TestSizeEventType verifies SizeEvent methods exist and compile.
func TestSizeEventType(t *testing.T) {
	e := &Event{}
	sizeEvt := AsSizeEvent(e)
	_ = sizeEvt.GetSize
	// Base event methods are inherited
	_ = sizeEvt.GetEventType
	_ = sizeEvt.Skip
}

// TestBindMethodsExist verifies Bind/Unbind methods are on BaseWindow.
func TestBindMethodsExist(t *testing.T) {
	w := &BaseWindow{}
	_ = w.Bind
	_ = w.BindId
	_ = w.BindIdRange
	_ = w.Unbind
	_ = w.UnbindId
}

// TestEventTypeConstants verifies event type constants are accessible.
func TestEventTypeConstants(t *testing.T) {
	// Command events
	_ = EVT_BUTTON
	_ = EVT_TEXT
	_ = EVT_MENU
	_ = EVT_COMMAND_BUTTON_CLICKED
	_ = EVT_COMMAND_CHECKBOX_CLICKED
	_ = EVT_COMMAND_CHOICE_SELECTED
	_ = EVT_COMMAND_LISTBOX_SELECTED
	_ = EVT_COMMAND_TEXT_UPDATED
	_ = EVT_COMMAND_TEXT_ENTER
	_ = EVT_COMMAND_MENU_SELECTED
	_ = EVT_COMMAND_SLIDER_UPDATED
	_ = EVT_COMMAND_RADIOBOX_SELECTED
	_ = EVT_COMMAND_RADIOBUTTON_SELECTED

	// Window events
	_ = EVT_CLOSE_WINDOW
	_ = EVT_SIZE
	_ = EVT_MOVE
	_ = EVT_PAINT
	_ = EVT_SHOW
	_ = EVT_ACTIVATE
}

// TestEventConversionRoundtrip verifies As*Event preserves the pointer.
func TestEventConversionRoundtrip(t *testing.T) {
	baseEvt := &Event{}
	ptr := baseEvt.Ptr()

	cmd := AsCommandEvent(baseEvt)
	if cmd.Ptr() != ptr {
		t.Fatal("AsCommandEvent lost pointer")
	}

	closeEvt := AsCloseEvent(baseEvt)
	if closeEvt.Ptr() != ptr {
		t.Fatal("AsCloseEvent lost pointer")
	}

	keyEvt := AsKeyEvent(baseEvt)
	if keyEvt.Ptr() != ptr {
		t.Fatal("AsKeyEvent lost pointer")
	}

	mouseEvt := AsMouseEvent(baseEvt)
	if mouseEvt.Ptr() != ptr {
		t.Fatal("AsMouseEvent lost pointer")
	}

	sizeEvt := AsSizeEvent(baseEvt)
	if sizeEvt.Ptr() != ptr {
		t.Fatal("AsSizeEvent lost pointer")
	}
}
