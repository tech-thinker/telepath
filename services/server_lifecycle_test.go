package services

import (
    "context"
    "sync"
    "testing"
    "time"
    "github.com/tech-thinker/telepath/models"
)

func TestNewServerInitialization(t *testing.T) {
    cfgs := []models.Config{{Name: "c1"}, {Name: "c2"}}
    s := NewServer(cfgs)
    if s == nil {
        t.Fatal("NewServer returned nil")
    }
}

func TestStartWithNilServer(t *testing.T) {
    s := &srv{cfg: []models.Config{{Name: "test"}}, liveConnections: make(map[string]*LiveConnection)}
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    s.start(ctx, models.Config{Name: "test", Server: nil})
}

func TestStartTypeR(t *testing.T) {
    s := &srv{cfg: []models.Config{}, liveConnections: make(map[string]*LiveConnection)}
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    cfg := models.Config{Name: "test", Type: "R", Server: &models.Server{Host: "x", Port: 22, Username: "u", AuthType: "PASS", Password: "p"}, LocalPort: 1, RemotePort: 1}
    s.start(ctx, cfg)
}

func TestStartTypeL(t *testing.T) {
    s := &srv{cfg: []models.Config{}, liveConnections: make(map[string]*LiveConnection)}
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    cfg := models.Config{Name: "test", Type: "L", Server: &models.Server{Host: "x", Port: 22, Username: "u", AuthType: "PASS", Password: "p"}, LocalPort: 1, RemotePort: 1}
    s.start(ctx, cfg)
}

func TestBuildSSHClientNilServer(t *testing.T) {
    s := &srv{}
    c, err := s.buildSSHClient(models.Config{Server: nil})
    if c != nil || err != nil {
        t.Fatalf("expected nil,nil got %v,%v", c, err)
    }
}

func TestStartAllMultiple(t *testing.T) {
    cfgs := []models.Config{{Name: "c1", Server: nil}, {Name: "c2", Server: nil}}
    s := &srv{cfg: cfgs, liveConnections: make(map[string]*LiveConnection)}
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    var wg sync.WaitGroup
    s.StartAll(ctx, &wg)
    wg.Wait()
}

func TestStopAllMultiple(t *testing.T) {
    s := &srv{liveConnections: map[string]*LiveConnection{
        "c1": NewLiveConnection(func() {}),
        "c2": NewLiveConnection(func() {}),
    }}
    s.StopAll()
}

func TestStartWithServerErrorAfterSSH(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    cfg := models.Config{
        Name:       "ssh-error",
        Type:       "L",
        Server:     &models.Server{Host: "invalid.host", Port: 22, Username: "u", AuthType: "PASS", Password: "p"},
        LocalPort:  8080,
        RemotePort: 80,
    }
    s := &srv{liveConnections: make(map[string]*LiveConnection)}
    s.start(ctx, cfg)
    time.Sleep(10 * time.Millisecond)
}

func TestStartWithAlreadyActiveConnection(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    cfg := models.Config{Name: "active", Server: nil}
    s := &srv{liveConnections: make(map[string]*LiveConnection)}
    active := NewLiveConnection(func() {})
    s.liveConnections[cfg.Name] = active
    s.start(ctx, cfg)
    if _, ok := s.liveConnections[cfg.Name]; !ok {
        t.Fatal("expected active connection to remain")
    }
}

func TestStartWithTypeRLocalBoth(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    s := &srv{liveConnections: make(map[string]*LiveConnection)}
    cfgR := models.Config{Name: "r", Type: "R", Server: &models.Server{Host: "x", Port: 22, Username: "u", AuthType: "PASS", Password: "p"}, LocalPort: 1, RemotePort: 1}
    s.start(ctx, cfgR)
    s2 := &srv{liveConnections: make(map[string]*LiveConnection)}
    cfgL := models.Config{Name: "l", Type: "L", Server: &models.Server{Host: "x", Port: 22, Username: "u", AuthType: "PASS", Password: "p"}, LocalPort: 1, RemotePort: 1}
    s2.start(ctx, cfgL)
    s3 := &srv{liveConnections: make(map[string]*LiveConnection)}
    cfgEmpty := models.Config{Name: "empty", Type: "", Server: &models.Server{Host: "x", Port: 22, Username: "u", AuthType: "PASS", Password: "p"}, LocalPort: 1, RemotePort: 1}
    s3.start(ctx, cfgEmpty)
    time.Sleep(10 * time.Millisecond)
}

func TestStartAllWithMixedConfigs(t *testing.T) {
    cfgs := []models.Config{
        {Name: "c1", Server: nil},
        {Name: "c2", Server: &models.Server{Host: "invalid", Port: 22, Username: "u", AuthType: "PASS", Password: "p"}},
        {Name: "c3", Server: nil},
    }
    s := &srv{cfg: cfgs, liveConnections: make(map[string]*LiveConnection)}
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    var wg sync.WaitGroup
    s.StartAll(ctx, &wg)
    wg.Wait()
}