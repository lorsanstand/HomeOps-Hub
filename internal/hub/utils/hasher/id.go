package hasher

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/lorsanstand/HomeOps-Hub/internal/domain"
)

func newSalt(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

func MakeID(info domain.HostInfo, AgentName string) (string, error) {
	salt, err := newSalt(10)
	if err != nil {
		return "", err
	}

	s := fmt.Sprintf("v1|host=%s|distro=%s|name=%s|", info.Hostname, info.Arch, AgentName)
	h := sha256.Sum256(append([]byte(s), salt...))
	return hex.EncodeToString(h[:16]), nil
}
