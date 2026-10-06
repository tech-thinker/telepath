package services

import "testing"

func TestLiveConnectionIsActive(t *testing.T) {
    lc := &LiveConnection{}
    if lc.IsActive() {
        t.Fatalf("expected inactive when cancel is nil")
    }
    lc2 := NewLiveConnection(func() {})
    if !lc2.IsActive() {
        t.Fatalf("expected active when cancel is set")
    }
}

func TestNewLiveConnection(t *testing.T) {
    cancelCalled := false
    lc := NewLiveConnection(func() { cancelCalled = true })
    if lc.cancel == nil {
        t.Fatal("cancel should be set")
    }
    lc.cancel()
    if !cancelCalled {
        t.Fatal("cancel function not called")
    }
}