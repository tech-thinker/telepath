package utils

import "encoding/base64"

func Base64Encode(plain string) string {
	return base64.StdEncoding.EncodeToString([]byte(plain))
}

func Base64Decode(cipher string) string {
	data, err := base64.StdEncoding.DecodeString(cipher)
	if err != nil {
		return ""
	}
	return string(data)
}
