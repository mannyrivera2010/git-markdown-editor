package auth

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"os"
	"strings"
)

const shadowFileName = "shadow_file"

type AuthService struct{}

func (a *AuthService) Init() error {
	if _, err := os.Stat(shadowFileName); os.IsNotExist(err) {
		salt := generateSalt()
		hash := hashPassword("password", salt)
		line := fmt.Sprintf("admin:%s:%s\n", base64.StdEncoding.EncodeToString(salt), base64.StdEncoding.EncodeToString(hash))
		return os.WriteFile(shadowFileName, []byte(line), 0600)
	}
	return nil
}
func (a *AuthService) Authenticate(username, password string) bool {
	file, err := os.Open(shadowFileName)
	if err != nil {
		return false
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")
		if len(parts) != 3 {
			continue
		}
		if parts[0] == username {
			salt, _ := base64.StdEncoding.DecodeString(parts[1])
			hash, _ := base64.StdEncoding.DecodeString(parts[2])
			if bytes.Equal(argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32), hash) {
				return true
			}
		}
	}
	return false
}
func generateSalt() []byte                   { b := make([]byte, 16); rand.Read(b); return b }
func hashPassword(p string, s []byte) []byte { return argon2.IDKey([]byte(p), s, 1, 64*1024, 4, 32) }
