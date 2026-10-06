package services

import (
    "os"
    "testing"
    "github.com/tech-thinker/telepath/models"
    "github.com/tech-thinker/telepath/constants"
)

func TestCreateSSHClientPasswordBranch(t *testing.T) {
    srv := &srv{}
    server := &models.Server{
        Host:     "invalid",
        Port:     22,
        Username: "user",
        AuthType: constants.CREDIENTIAL_PASS,
        Password: "cGFzc3dvcmQ=",
    }
    _, err := srv.createSSHClient(server, nil)
    if err == nil {
        t.Fatal("expected error")
    }
}

func TestCreateSSHClientKeyMissing(t *testing.T) {
    srv := &srv{}
    server := &models.Server{
        Host:     "invalid",
        Port:     22,
        Username: "user",
        AuthType: constants.CREDIENTIAL_KEY,
        Key:      "/tmp/nonexistent-key-xyz",
    }
    _, err := srv.createSSHClient(server, nil)
    if err == nil {
        t.Fatal("expected error for missing key")
    }
}

func TestCreateSSHClientKeyInvalid(t *testing.T) {
    f, _ := os.CreateTemp("", "bad-key")
    f.WriteString("INVALID KEY")
    f.Close()
    defer os.Remove(f.Name())

    srv := &srv{}
    server := &models.Server{
        Host:     "invalid",
        Port:     22,
        Username: "user",
        AuthType: constants.CREDIENTIAL_KEY,
        Key:      f.Name(),
    }
    _, err := srv.createSSHClient(server, nil)
    if err == nil {
        t.Fatal("expected error for invalid key")
    }
}

func TestCreateSSHClientKeyAuthNoPassphrase(t *testing.T) {
    f, _ := os.CreateTemp("", "key-no-pass")
    keyData := `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA0Z3VS5JJcds3s5uEsw9Z3vNaHxnhW9VV7iHZJoB6f8bC5v3u
rKz0dK7vLJnCKMhfJzV5GvKKHHMjxKaQUKWZB4YZvKQZq5J8pZ5sJ6L5K5V5K5V5
K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5
K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5
K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5
QIDAQABAoIBAC5N2sVRVXZCaAf0wKH2rrVVl5XZE8rPl5YzqzqZJmUQR5pP5sV5
K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5
K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5
K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5
K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5
-----END RSA PRIVATE KEY-----`
    f.WriteString(keyData)
    f.Close()
    defer os.Remove(f.Name())

    srv := &srv{}
    server := &models.Server{
        Host:     "invalid",
        Port:     22,
        Username: "user",
        AuthType: constants.CREDIENTIAL_KEY,
        Key:      f.Name(),
    }
    _, err := srv.createSSHClient(server, nil)
    if err == nil {
        t.Fatal("expected error for invalid host")
    }
}

func TestCreateSSHClientWithPassphraseKey(t *testing.T) {
    f, _ := os.CreateTemp("", "test-key-*.pem")
    keyData := `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA0Z3VS5JJcds3s5uEsw9Z3vNaHxnhW9VV7iHZJoB6f8bC5v3u
rKz0dK7vLJnCKMhfJzV5GvKKHHMjxKaQUKWZB4YZvKQZq5J8pZ5sJ6L5K5V5K5V5
K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5
K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5
K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5
QIDAQABAoIBAC5N2sVRVXZCaAf0wKH2rrVVl5XZE8rPl5YzqzqZJmUQR5pP5sV5
K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5
K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5
K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5
K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5K5V5
-----END RSA PRIVATE KEY-----`
    f.WriteString(keyData)
    f.Close()
    defer os.Remove(f.Name())

    srv := &srv{}
    server := &models.Server{
        Host:       "invalid",
        Port:       22,
        Username:   "user",
        AuthType:   constants.CREDIENTIAL_KEY,
        Key:        f.Name(),
        Passphrase: "cGFzc3BocmFzZQ==",
    }
    _, err := srv.createSSHClient(server, nil)
    if err == nil {
        t.Fatal("expected error for invalid host")
    }
}

func TestMakeHostChainNil(t *testing.T) {
    s := &srv{}
    chain := s.makeHostChain(nil)
    if len(chain) != 0 {
        t.Fatalf("expected 0, got %d", len(chain))
    }
}

func TestMakeHostChainOrdering(t *testing.T) {
    f := &models.Server{Host: "final"}
    h := &models.Server{Host: "hop", Jump: f}
    s := &srv{}
    chain := s.makeHostChain(h)
    if len(chain) != 2 || chain[0].Host != "final" || chain[1].Host != "hop" {
        t.Fatal("chain order wrong")
    }
}

func TestMakeHostChainThreeLevels(t *testing.T) {
    f := &models.Server{Host: "final"}
    h2 := &models.Server{Host: "hop2", Jump: f}
    h1 := &models.Server{Host: "hop1", Jump: h2}
    s := &srv{}
    chain := s.makeHostChain(h1)
    if len(chain) != 3 {
        t.Fatalf("expected 3, got %d", len(chain))
    }
    if chain[0].Host != "final" || chain[1].Host != "hop2" || chain[2].Host != "hop1" {
        t.Fatal("wrong order")
    }
}

func TestBuildSSHClientWithMultipleHops(t *testing.T) {
    cfg := models.Config{
        Server: &models.Server{
            Host: "final",
            Port: 22,
            Username: "u",
            AuthType: constants.CREDIENTIAL_PASS,
            Password: "p",
            Jump: &models.Server{
                Host: "jump",
                Port: 22,
                Username: "u",
                AuthType: constants.CREDIENTIAL_PASS,
                Password: "p",
            },
        },
    }
    s := &srv{}
    _, err := s.buildSSHClient(cfg)
    if err == nil {
        t.Fatal("expected error for invalid host")
    }
}