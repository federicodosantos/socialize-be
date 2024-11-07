package md5

import (
	"crypto/md5"
	"encoding/hex"
)

func HashWithMd5(text string) string {
	hash := md5.New()

	data := []byte(text)

	hashedData := hash.Sum(data)

	return hex.EncodeToString(hashedData[:])
}

