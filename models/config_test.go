package models

import (
    "testing"
    "github.com/tech-thinker/telepath/constants"
)

func TestServerValidateSuccessKey(t *testing.T) {
    s := &Server{
        Host:     "example.com",
        Port:     22,
        Username: "user",
        AuthType: constants.CREDIENTIAL_KEY,
        Key:      "/tmp/key",
    }
    if err := s.Validate("test"); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}

func TestServerValidateSuccessPass(t *testing.T) {
    s := &Server{
        Host:     "example.com",
        Port:     22,
        Username: "user",
        AuthType: constants.CREDIENTIAL_PASS,
        Password: "cGFzc3dvcmQ=", // "password" base64
    }
    if err := s.Validate("test"); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}

func TestServerValidateMissingFields(t *testing.T) {
    tests := []struct {
        name   string
        server *Server
        want   string
    }{
        {"missing host", &Server{Port: 22, Username: "u", AuthType: constants.CREDIENTIAL_PASS, Password: "cGFz"}, "`test.host` must be present."},
        {"invalid port", &Server{Host: "h", Username: "u", AuthType: constants.CREDIENTIAL_PASS, Password: "cGFz"}, "`test.port` is invalid."},
        {"missing username", &Server{Host: "h", Port: 22, AuthType: constants.CREDIENTIAL_PASS, Password: "cGFz"}, "`test.username` must be present."},
        {"invalid auth", &Server{Host: "h", Port: 22, Username: "u", AuthType: "XYZ", Password: "cGFz"}, "`test.authType` must be present."},
        {"key auth missing key", &Server{Host: "h", Port: 22, Username: "u", AuthType: constants.CREDIENTIAL_KEY}, "`test.key` must be present for authType KEY."},
        {"pass auth missing password", &Server{Host: "h", Port: 22, Username: "u", AuthType: constants.CREDIENTIAL_PASS}, "`test.password` must be present for authType PASS."},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.server.Validate("test")
            if err == nil || err.Error() != tt.want {
                t.Fatalf("expected error %q, got %v", tt.want, err)
            }
        })
    }
}

func TestServerValidateJump(t *testing.T) {
    // nested jump chain with valid inner server
    inner := &Server{Host: "inner", Port: 22, Username: "u", AuthType: constants.CREDIENTIAL_PASS, Password: "cGFz"}
    top := &Server{Host: "top", Port: 22, Username: "u", AuthType: constants.CREDIENTIAL_PASS, Password: "cGFz", Jump: inner}
    if err := top.Validate("test"); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}

func TestConfigValidateSuccess(t *testing.T) {
    cfg := Config{
        Name:      "cfg1",
        Type:      "L",
        LocalPort: 8080,
        RemotePort: 80,
        Server:    &Server{Host: "h", Port: 22, Username: "u", AuthType: constants.CREDIENTIAL_PASS, Password: "cGFz"},
    }
    if err := cfg.Validate(); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}

func TestConfigValidateErrors(t *testing.T) {
    baseServer := &Server{Host: "h", Port: 22, Username: "u", AuthType: constants.CREDIENTIAL_PASS, Password: "cGFz"}
    tests := []struct {
        name string
        cfg  Config
        want string
    }{
        {"missing name", Config{Type: "L", LocalPort: 1, RemotePort: 1, Server: baseServer}, "`name` must be present."},
        {"invalid type", Config{Name: "n", Type: "X", LocalPort: 1, RemotePort: 1, Server: baseServer}, "`type` must be L or R."},
        {"invalid localPort", Config{Name: "n", Type: "L", LocalPort: 0, RemotePort: 1, Server: baseServer}, "`localPort` is invalid."},
        {"invalid remotePort", Config{Name: "n", Type: "L", LocalPort: 1, RemotePort: 0, Server: baseServer}, "`remotePort` is invalid."},
        {"missing server", Config{Name: "n", Type: "L", LocalPort: 1, RemotePort: 1}, "`server` must be present."},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.cfg.Validate()
            if err == nil || err.Error() != tt.want {
                t.Fatalf("expected %q, got %v", tt.want, err)
            }
        })
    }
}
