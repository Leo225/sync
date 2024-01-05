package sync

import (
	"crypto/rand"
	"fmt"
)

func sessionID() string {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", buf)
}
