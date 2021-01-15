package typtop

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"log"

	"golang.org/x/crypto/pbkdf2"
)

type State struct {
	Publickey *rsa.PublicKey
	// cipheredCacheState
	// TypoCache
	// WaitList
	gamma int
}

func Register(password string, keyLength int, typoCacheLength int) (state *State) {
	// generate key
	privateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		log.Printf("Failed to generate RSA key\n")
	}
	publicKey := &privateKey.PublicKey

	// TODO generate different salt for every user & store
	salt := []byte("salt")

	typoCache := make([][]byte, typoCacheLength)

	// encrypt the private key using pbkdf, using the password
	typoCache[0] = pbkdf2.Key([]byte(password), salt, 4096, 32, sha1.New)

	// warm up typo cache
	// warmTypoCache(typoCache)

	state = &State{
		Publickey: publicKey,
	}

	return state
}

// func warmTypoCache(typoCache [][]byte) [][]byte {
// 	for i := range typoCache {

// 	}
// }
