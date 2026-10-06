package services

import (
    "context"
    "errors"
    "net"
    "sync"
    "testing"
    "time"
    "github.com/tech-thinker/telepath/models"
)

// mockListener implements net.Listener for testing acceptLoop.
type mockListener struct {
    acceptCh chan net.Conn
    closeCh  chan struct{}
    addr     net.Addr
    errOnce  error
}

func (m *mockListener) Accept() (net.Conn, error) {
    if m.errOnce != nil {
        err := m.errOnce
        m.errOnce = nil
        return nil, err
    }
    select {
    case <-m.closeCh:
        return nil, errors.New("closed")
    case c := <-m.acceptCh:
        return c, nil
    }
}
func (m *mockListener) Close() error {
    close(m.closeCh)
    return nil
}
func (m *mockListener) Addr() net.Addr { return m.addr }

func newMockListener() *mockListener {
    return &mockListener{acceptCh: make(chan net.Conn, 1), closeCh: make(chan struct{}), addr: &net.TCPAddr{}}
}

func TestStartLocalWithListenerError(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    s := &srv{}
    cfg := models.Config{Name: "bad", LocalPort: -1, LocalHost: "127.0.0.1", RemotePort: 80, RemoteHost: "127.0.0.1"}
    s.startLocal(ctx, cfg, nil)
}

func TestStartLocalWithValidListener(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    cfg := models.Config{
        Name:       "valid-local",
        LocalPort:  0,
        LocalHost:  "127.0.0.1",
        RemotePort: 80,
        RemoteHost: "127.0.0.1",
    }
    s := &srv{liveConnections: make(map[string]*LiveConnection)}
    go func() {
        time.Sleep(10 * time.Millisecond)
        cancel()
    }()
    s.startLocal(ctx, cfg, nil)
}

func TestStartRemoteWithNilServerSkipped(t *testing.T) {
    t.Skip("startRemote requires real SSH client; cannot test with nil client without production code changes")
}

func TestAcceptLoopContextCancel(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    ml := newMockListener()
    done := make(chan struct{})
    go func() {
        s := &srv{}
        s.acceptLoop(ctx, models.Config{Name: "test"}, ml, func(c net.Conn) { c.Close() })
        close(done)
    }()
    time.Sleep(10 * time.Millisecond)
    cancel()
    ml.Close()
    select {
    case <-done:
    case <-time.After(500 * time.Millisecond):
        t.Fatal("acceptLoop did not exit after context cancel")
    }
}

func TestAcceptLoopErrorExits(t *testing.T) {
    ctx := context.Background()
    ml := newMockListener()
    ml.errOnce = errors.New("forced error")
    s := &srv{}
    s.acceptLoop(ctx, models.Config{Name: "test"}, ml, func(c net.Conn) { c.Close() })
    ml.Close()
}

func TestAcceptLoopMultipleAccepts(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    ml := newMockListener()
    cfg := models.Config{Name: "multi-accept"}
    s := &srv{}
    var mu sync.Mutex
    handlerCalls := 0
    handler := func(c net.Conn) {
        mu.Lock()
        handlerCalls++
        mu.Unlock()
        c.Close()
    }
    done := make(chan struct{})
    go func() {
        s.acceptLoop(ctx, cfg, ml, handler)
        close(done)
    }()
    for i := 0; i < 3; i++ {
        c1, c2 := net.Pipe()
        ml.acceptCh <- c1
        c2.Close()
    }
    time.Sleep(50 * time.Millisecond)
    cancel()
    ml.Close()
    select {
    case <-done:
    case <-time.After(500 * time.Millisecond):
        t.Fatal("acceptLoop didn't exit")
    }
    mu.Lock()
    if handlerCalls != 3 {
        t.Fatalf("expected 3 handler calls, got %d", handlerCalls)
    }
    mu.Unlock()
}

func TestAcceptLoopAcceptErrorNonContext(t *testing.T) {
    ctx := context.Background()
    ml := newMockListener()
    ml.errOnce = errors.New("non-context accept error")
    cfg := models.Config{Name: "accept-error"}
    s := &srv{}
    s.acceptLoop(ctx, cfg, ml, func(c net.Conn) { c.Close() })
    ml.Close()
}

func TestPipeCopiesBidirectional(t *testing.T) {
    a, b := net.Pipe()
    defer a.Close()
    defer b.Close()
    go pipe(a, b)
    msg1 := []byte("hello from a")
    msg2 := []byte("reply from b")
    if _, err := a.Write(msg1); err != nil {
        t.Fatalf("write a->b failed: %v", err)
    }
    buf := make([]byte, len(msg1))
    if _, err := b.Read(buf); err != nil {
        t.Fatalf("read b failed: %v", err)
    }
    if string(buf) != string(msg1) {
        t.Fatalf("expected %s, got %s", msg1, buf)
    }
    if _, err := b.Write(msg2); err != nil {
        t.Fatalf("write b->a failed: %v", err)
    }
    buf2 := make([]byte, len(msg2))
    if _, err := a.Read(buf2); err != nil {
        t.Fatalf("read a failed: %v", err)
    }
    if string(buf2) != string(msg2) {
        t.Fatalf("expected %s, got %s", msg2, buf2)
    }
    time.Sleep(10 * time.Millisecond)
}

func TestPipeWithImmediateClose(t *testing.T) {
    a, b := net.Pipe()
    b.Close()
    pipe(a, b)
}