package utils

import (
    "encoding/base64"
    "testing"
)

func TestBase64EncodeDecode(t *testing.T) {
    plain := "hello world"
    encoded := Base64Encode(plain)
    expected := base64.StdEncoding.EncodeToString([]byte(plain))
    if encoded != expected {
        t.Fatalf("expected %s, got %s", expected, encoded)
    }

    decoded := Base64Decode(encoded)
    if decoded != plain {
        t.Fatalf("expected decoded %s, got %s", plain, decoded)
    }
}

func TestBase64DecodeInvalid(t *testing.T) {
    // non‑base64 string should return empty string
    if got := Base64Decode("@@@invalid@@@"); got != "" {
        t.Fatalf("expected empty string for invalid base64, got %s", got)
    }
}
