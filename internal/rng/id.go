package rng

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func SessionID() (string, error) {
	b := make([]byte, 4)
	_, err := io.LimitReader(rand.Reader, 4).Read(b)
	if err != nil {
		return "", fmt.Errorf("read rand id: %v", err)
	}
	return "s_" + hex.EncodeToString(b), nil
}

func ViewID() (string, error) {
	b := make([]byte, 4)
	_, err := io.LimitReader(rand.Reader, 4).Read(b)
	if err != nil {
		return "", fmt.Errorf("read rand id: %v", err)
	}
	return "v_" + hex.EncodeToString(b), nil
}
