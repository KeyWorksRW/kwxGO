package wx

import (
	"testing"
	"unsafe"
)

func TestPoint(t *testing.T) {
	tests := []struct {
		name string
		p    Point
		x    int
		y    int
	}{
		{"Zero", Point{0, 0}, 0, 0},
		{"Positive", Point{100, 200}, 100, 200},
		{"Negative", Point{-50, -100}, -50, -100},
		{"DefaultPosition", DefaultPosition, -1, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.p.X != tt.x {
				t.Errorf("X mismatch: got %d, want %d", tt.p.X, tt.x)
			}
			if tt.p.Y != tt.y {
				t.Errorf("Y mismatch: got %d, want %d", tt.p.Y, tt.y)
			}
		})
	}
}

func TestSize(t *testing.T) {
	tests := []struct {
		name string
		s    Size
		w    int
		h    int
	}{
		{"Zero", Size{0, 0}, 0, 0},
		{"Positive", Size{800, 600}, 800, 600},
		{"DefaultSize", DefaultSize, -1, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.s.W != tt.w {
				t.Errorf("W mismatch: got %d, want %d", tt.s.W, tt.w)
			}
			if tt.s.H != tt.h {
				t.Errorf("H mismatch: got %d, want %d", tt.s.H, tt.h)
			}
		})
	}
}

func TestRect(t *testing.T) {
	r := Rect{X: 10, Y: 20, W: 300, H: 400}

	if r.X != 10 || r.Y != 20 || r.W != 300 || r.H != 400 {
		t.Errorf("Rect values incorrect: got (%d,%d,%d,%d), want (10,20,300,400)",
			r.X, r.Y, r.W, r.H)
	}
}

func TestBaseWindow(t *testing.T) {
	// Test nil window
	var w BaseWindow
	if !w.IsNil() {
		t.Error("New BaseWindow should have nil pointer")
	}
	if w.Ptr() != nil {
		t.Error("Ptr() should return nil for uninitialized window")
	}

	// Test SetPtr with a dummy value
	// We use the address of a local variable as a safe test pointer
	dummy := 42
	dummyPtr := unsafe.Pointer(&dummy)
	w.SetPtr(dummyPtr)
	if w.IsNil() {
		t.Error("IsNil() should return false after SetPtr")
	}
	if w.Ptr() != dummyPtr {
		t.Error("Ptr() should return the pointer set by SetPtr")
	}

	// Test SetPtr with nil
	w.SetPtr(nil)
	if !w.IsNil() {
		t.Error("IsNil() should return true after SetPtr(nil)")
	}
}

func TestDefaultPositionAndSize(t *testing.T) {
	// Verify default values are -1 (wxDefaultCoord)
	if DefaultPosition.X != -1 || DefaultPosition.Y != -1 {
		t.Errorf("DefaultPosition should be (-1, -1), got (%d, %d)",
			DefaultPosition.X, DefaultPosition.Y)
	}
	if DefaultSize.W != -1 || DefaultSize.H != -1 {
		t.Errorf("DefaultSize should be (-1, -1), got (%d, %d)",
			DefaultSize.W, DefaultSize.H)
	}
}
