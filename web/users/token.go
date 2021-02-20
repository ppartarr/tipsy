package users

import (
	"crypto/md5"
	"encoding/hex"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// hash email using bcrypt then hash again with md5 so that hash can be used as token in a url
func generateToken(email string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	// TODO store hash
	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil))
}
