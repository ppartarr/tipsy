package typtop

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"log"
)

func RSAEncrypt(publicKey *rsa.PublicKey, message []byte) []byte {
	encryptedEpsilon, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, message, nil)
	if err != nil {
		log.Fatal("failed encrypt the message")
	}
	return encryptedEpsilon
}

func RSADecrypt(privateKey *rsa.PrivateKey, ciphertext []byte) []byte {
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
	if err != nil {
		log.Fatal("failed to decrypt the ciphertext")
	}
	return plaintext
}
