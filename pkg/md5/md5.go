package md5

import (
	"crypto/md5"
	"encoding/hex"
	"log"
)

func HashWithMd5(text string) string {
	hash := md5.New()

	hash.Write([]byte(text))

	hashedData := hash.Sum(nil)

	res := hex.EncodeToString(hashedData[:])
	
	log.Printf("res : %s", res)

	return res
}