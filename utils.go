package main

import "encoding/base64"

// Error struct
type Error struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// Base64Decode method
func Base64Decode(message []byte) (b []byte, err error) {
	var l int
	b = make([]byte, base64.StdEncoding.DecodedLen(len(message)))
	l, err = base64.StdEncoding.Decode(b, message)
	if err != nil {
		return
	}
	return b[:l], nil
}
