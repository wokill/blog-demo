package common

import (
	"crypto/sha256"
	"fmt"
	"git.dustess.com/mk-base/util/crypto"
	"strings"
)

// Encrypt 加密密码
func Encrypt(pwd string) string {
	salt := crypto.RandID()
	return salt + "." + Sha256([]byte(salt+pwd))
}

// Verify 校验密码 `pwd` 是不是明文 `plaintext` 的密码
func Verify(plaintext, pwd string) bool {
	s := strings.Split(pwd, ".")
	if len(s) != 2 {
		return false
	}
	return s[1] == Sha256([]byte(s[0]+plaintext))
}

// Sha256 计算sha256
func Sha256(data []byte) string {
	h := sha256.New()
	h.Write(data)
	s := h.Sum(nil)
	return fmt.Sprintf("%x", s)
}
