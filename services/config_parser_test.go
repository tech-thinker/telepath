package services

import (
    "os"
    "testing"
)

func TestParseConfigComments(t *testing.T) {
    raw := `[
        // single line comment
        {
            "name": "test",
            "type": "L",
            "localPort": 8080,
            "localHost": "127.0.0.1",
            "remotePort": 80,
            "remoteHost": "example.com",
            "server": {
                "host": "ssh.example.com",
                "port": 22,
                "username": "user",
                "authType": "PASS",
                "password": "cGFzc3dvcmQ="
            }
        } /* multi line comment */
    ]`
    f, err := os.CreateTemp("", "config-*.json")
    if err != nil {
        t.Fatalf("temp file: %v", err)
    }
    defer os.Remove(f.Name())
    if _, err := f.Write([]byte(raw)); err != nil {
        t.Fatalf("write temp: %v", err)
    }
    f.Close()

    cfgs, err := ParseConfig(f.Name())
    if err != nil {
        t.Fatalf("ParseConfig error: %v", err)
    }
    if len(cfgs) != 1 {
        t.Fatalf("expected 1 config, got %d", len(cfgs))
    }
    if cfgs[0].Name != "test" {
        t.Fatalf("unexpected config name: %s", cfgs[0].Name)
    }
    if cfgs[0].Server.AuthType != "PASS" {
        t.Fatalf("unexpected authType: %s", cfgs[0].Server.AuthType)
    }
}

func TestParseConfigMissingFile(t *testing.T) {
    _, err := ParseConfig("/tmp/nonexistent-config-xyz.json")
    if err == nil {
        t.Fatal("expected error for missing file")
    }
}

func TestParseConfigInvalidJSON(t *testing.T) {
    f, _ := os.CreateTemp("", "bad-json")
    f.WriteString("{ invalid json }")
    f.Close()
    defer os.Remove(f.Name())

    _, err := ParseConfig(f.Name())
    if err == nil {
        t.Fatal("expected JSON parse error")
    }
}

func TestParseConfigValidWithComments(t *testing.T) {
    f, _ := os.CreateTemp("", "config")
    f.WriteString(`[
        // comment
        {
            "name": "test",
            "type": "L",
            "localPort": 8080,
            "localHost": "127.0.0.1",
            "remotePort": 80,
            "remoteHost": "127.0.0.1",
            "server": {
                "host": "h",
                "port": 22,
                "username": "u",
                "authType": "PASS",
                "password": "p"
            }
        } /* comment */
    ]`)
    f.Close()
    defer os.Remove(f.Name())

    cfgs, err := ParseConfig(f.Name())
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(cfgs) != 1 {
        t.Fatalf("expected 1 config, got %d", len(cfgs))
    }
}