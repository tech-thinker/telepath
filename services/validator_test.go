package services

import (
    "testing"
    "github.com/tech-thinker/telepath/models"
    "github.com/tech-thinker/telepath/constants"
)

func TestValidateConfigSuccess(t *testing.T) {
    cfg := models.Config{
        Name:      "valid",
        Type:      "L",
        LocalPort: 8080,
        RemotePort: 80,
        Server:    &models.Server{Host: "h", Port: 22, Username: "u", AuthType: constants.CREDIENTIAL_PASS, Password: "cGFz"},
    }
    valid, err := ValidateConfig([]models.Config{cfg})
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(valid) != 1 {
        t.Fatalf("expected 1 valid config, got %d", len(valid))
    }
}

func TestValidateConfigError(t *testing.T) {
    invalid := models.Config{Type: "L", LocalPort: 8080, RemotePort: 80, Server: &models.Server{Host: "h", Port: 22, Username: "u", AuthType: constants.CREDIENTIAL_PASS, Password: "cGFz"}}
    _, err := ValidateConfig([]models.Config{invalid})
    if err == nil {
        t.Fatalf("expected error for invalid config, got nil")
    }
}