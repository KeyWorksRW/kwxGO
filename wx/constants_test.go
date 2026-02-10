package wx

import (
	"testing"
)

// TestStandardIDs verifies standard ID constants are non-zero
func TestStandardIDs(t *testing.T) {
	ids := map[string]int{
		"ID_ANY":       ID_ANY,
		"ID_OK":        ID_OK,
		"ID_CANCEL":    ID_CANCEL,
		"ID_YES":       ID_YES,
		"ID_NO":        ID_NO,
		"ID_CLOSE":     ID_CLOSE,
		"ID_HELP":      ID_HELP,
		"ID_NEW":       ID_NEW,
		"ID_OPEN":      ID_OPEN,
		"ID_SAVE":      ID_SAVE,
		"ID_EXIT":      ID_EXIT,
		"ID_ABOUT":     ID_ABOUT,
	}

	for name, value := range ids {
		if value == 0 {
			t.Errorf("%s should not be 0", name)
		}
	}
}

// TestWindowStyles verifies window style constants are non-zero
func TestWindowStyles(t *testing.T) {
	styles := map[string]int{
		"DEFAULT_FRAME_STYLE": DEFAULT_FRAME_STYLE,
		"CAPTION":             CAPTION,
		"RESIZE_BORDER":       RESIZE_BORDER,
		"SYSTEM_MENU":         SYSTEM_MENU,
		"MINIMIZE_BOX":        MINIMIZE_BOX,
		"MAXIMIZE_BOX":        MAXIMIZE_BOX,
		"CLOSE_BOX":           CLOSE_BOX,
	}

	for name, value := range styles {
		if value == 0 {
			t.Errorf("%s should not be 0", name)
		}
	}
}

// TestSizerFlags verifies sizer flag constants are non-zero
func TestSizerFlags(t *testing.T) {
	flags := map[string]int{
		"EXPAND": EXPAND,
		"ALL":    ALL,
		"LEFT":   LEFT,
		"RIGHT":  RIGHT,
		"TOP":    TOP,
		"BOTTOM": BOTTOM,
	}

	for name, value := range flags {
		if value == 0 {
			t.Errorf("%s should not be 0", name)
		}
	}
}

// TestAlignment verifies alignment constants are non-zero
func TestAlignment(t *testing.T) {
	alignments := map[string]int{
		"ALIGN_LEFT":              ALIGN_LEFT,
		"ALIGN_RIGHT":             ALIGN_RIGHT,
		"ALIGN_TOP":               ALIGN_TOP,
		"ALIGN_BOTTOM":            ALIGN_BOTTOM,
		"ALIGN_CENTER":            ALIGN_CENTER,
		"ALIGN_CENTER_HORIZONTAL": ALIGN_CENTER_HORIZONTAL,
		"ALIGN_CENTER_VERTICAL":   ALIGN_CENTER_VERTICAL,
	}

	for name, value := range alignments {
		if value == 0 {
			t.Errorf("%s should not be 0", name)
		}
	}

	// Verify CENTER and CENTRE are equivalent
	if ALIGN_CENTER != ALIGN_CENTRE {
		t.Error("ALIGN_CENTER and ALIGN_CENTRE should be equivalent")
	}
	if ALIGN_CENTER_HORIZONTAL != ALIGN_CENTRE_HORIZONTAL {
		t.Error("ALIGN_CENTER_HORIZONTAL and ALIGN_CENTRE_HORIZONTAL should be equivalent")
	}
	if ALIGN_CENTER_VERTICAL != ALIGN_CENTRE_VERTICAL {
		t.Error("ALIGN_CENTER_VERTICAL and ALIGN_CENTRE_VERTICAL should be equivalent")
	}
}

// TestOrientation verifies orientation constants are non-zero
func TestOrientation(t *testing.T) {
	orientations := map[string]int{
		"HORIZONTAL": HORIZONTAL,
		"VERTICAL":   VERTICAL,
		"BOTH":       BOTH,
	}

	for name, value := range orientations {
		if value == 0 {
			t.Errorf("%s should not be 0", name)
		}
	}
}

// TestTextControlStyles verifies text control style constants are non-zero
func TestTextControlStyles(t *testing.T) {
	styles := map[string]int{
		"TE_MULTILINE":     TE_MULTILINE,
		"TE_READONLY":      TE_READONLY,
		"TE_PROCESS_ENTER": TE_PROCESS_ENTER,
		"TE_PASSWORD":      TE_PASSWORD,
	}

	for name, value := range styles {
		if value == 0 {
			t.Errorf("%s should not be 0", name)
		}
	}
}

// TestEventTypes verifies event type constants are non-zero
func TestEventTypes(t *testing.T) {
	events := map[string]int{
		"EVT_COMMAND_BUTTON_CLICKED": EVT_COMMAND_BUTTON_CLICKED,
		"EVT_CLOSE_WINDOW":           EVT_CLOSE_WINDOW,
		"EVT_SIZE":                   EVT_SIZE,
		"EVT_BUTTON":                 EVT_BUTTON,
		"EVT_TEXT":                   EVT_TEXT,
	}

	for name, value := range events {
		if value == 0 {
			t.Errorf("%s should not be 0", name)
		}
	}

	// Verify aliases
	if EVT_BUTTON != EVT_COMMAND_BUTTON_CLICKED {
		t.Error("EVT_BUTTON should be an alias for EVT_COMMAND_BUTTON_CLICKED")
	}
	if EVT_TEXT != EVT_COMMAND_TEXT_UPDATED {
		t.Error("EVT_TEXT should be an alias for EVT_COMMAND_TEXT_UPDATED")
	}
}

// TestMessageBoxStyles verifies message box style constants are non-zero
func TestMessageBoxStyles(t *testing.T) {
	styles := map[string]int{
		"OK":               OK,
		"CANCEL":           CANCEL,
		"YES_NO":           YES_NO,
		"ICON_INFORMATION": ICON_INFORMATION,
		"ICON_WARNING":     ICON_WARNING,
		"ICON_ERROR":       ICON_ERROR,
		"ICON_QUESTION":    ICON_QUESTION,
	}

	for name, value := range styles {
		if value == 0 {
			t.Errorf("%s should not be 0", name)
		}
	}
}

// TestBorderStyles verifies border style constants exist
func TestBorderStyles(t *testing.T) {
	// BORDER_NONE can legitimately be 0, so we just check it's defined
	_ = BORDER_NONE

	styles := map[string]int{
		"BORDER_DEFAULT": BORDER_DEFAULT,
		"BORDER_SIMPLE":  BORDER_SIMPLE,
		"BORDER_SUNKEN":  BORDER_SUNKEN,
		"BORDER_RAISED":  BORDER_RAISED,
	}

	for name, value := range styles {
		if value == 0 {
			t.Errorf("%s should not be 0", name)
		}
	}
}

// TestButtonStyles verifies button style constants exist
func TestButtonStyles(t *testing.T) {
	// Just verify they're accessible
	_ = BU_LEFT
	_ = BU_RIGHT
	_ = BU_TOP
	_ = BU_BOTTOM
	_ = BU_EXACTFIT
}
