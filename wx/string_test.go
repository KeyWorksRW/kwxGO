package wx

import (
	"testing"
	"unsafe"
)

// TestStringRoundTrip verifies that Go string -> wxString -> Go string works
func TestStringRoundTrip(t *testing.T) {
	tests := []string{
		"",
		"Hello",
		"Hello, World!",
		"Unicode: 你好世界",
		"Emoji: 😀🎉",
		"Multi\nline\ntext",
		"Tab\tseparated\tvalues",
	}

	for _, original := range tests {
		t.Run(original, func(t *testing.T) {
			// Create wxString
			wxStr := NewWxString(original)
			if wxStr == nil {
				t.Fatal("NewWxString returned nil")
			}
			defer wxStr.Free()

			// Convert back to Go
			result := wxStr.String()

			if result != original {
				t.Errorf("String round-trip failed\nOriginal: %q\nGot:      %q", original, result)
			}
		})
	}
}

// TestWxStringFree verifies that Free() can be called multiple times safely
func TestWxStringFree(t *testing.T) {
	wxStr := NewWxString("test")
	if wxStr == nil {
		t.Fatal("NewWxString returned nil")
	}

	// First free
	wxStr.Free()
	if wxStr.ptr != nil {
		t.Error("Free() should set ptr to nil")
	}

	// Second free should not panic
	wxStr.Free()
}

// TestWxStringToGoAndFree verifies the helper function
func TestWxStringToGoAndFree(t *testing.T) {
	// This test can't directly create a wxString* like wxFFI functions would,
	// so we'll just verify it handles nil properly
	result := WxStringToGoAndFree(nil)
	if result != "" {
		t.Errorf("WxStringToGoAndFree(nil) should return empty string, got %q", result)
	}
}

// TestWithWxString verifies the helper function
func TestWithWxString(t *testing.T) {
	original := "Test String"
	var capturedPtr unsafe.Pointer

	withWxString(original, func(ptr unsafe.Pointer) {
		// Just verify we got a non-nil pointer
		if ptr == nil {
			t.Error("withWxString passed nil pointer")
		}
		capturedPtr = ptr
	})

	// After the function returns, the wxString should be freed
	// We can't directly verify this, but at least the test runs without crash
	if capturedPtr == nil {
		t.Error("Callback was not executed")
	}
}
