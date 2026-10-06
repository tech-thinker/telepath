package services

import (
    "testing"
    "github.com/tech-thinker/telepath/models"
)

func TestMakeHostChainOrder(t *testing.T) {
    hop2 := &models.Server{Host: "hop2", Port: 2222}
    hop1 := &models.Server{Host: "hop1", Port: 2221, Jump: hop2}
    final := &models.Server{Host: "final", Port: 22, Jump: hop1}

    s := &srv{cfg: nil, liveConnections: nil}
    chain := s.makeHostChain(final)
    if len(chain) != 3 {
        t.Fatalf("expected chain length 3, got %d", len(chain))
    }
    expected := []string{"hop2", "hop1", "final"}
    for i, srv := range chain {
        if srv.Host != expected[i] {
            t.Fatalf("at index %d expected host %s, got %s", i, expected[i], srv.Host)
        }
    }
}