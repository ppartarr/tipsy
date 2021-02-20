package typtop

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"

	"golang.org/x/crypto/scrypt"
)

func encodeKey(privateKey *rsa.PrivateKey) []byte {
	encodedKey := x509.MarshalPKCS1PrivateKey(privateKey)
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: encodedKey,
		},
	)
}

func decodeKey(pemBytes []byte) (key *rsa.PrivateKey) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		log.Println("failed to decode from pem")
	}
	// log.Println(len(block.Bytes))
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Println("failed to decode from PKCS1")
	}
	return key
}

// Will derive an AES key from a password using scrypt
func deriveAESKey(password string, salt []byte) ([]byte, []byte) {
	if salt == nil {
		salt = make([]byte, 8)
		_, err := rand.Read(salt)
		if err != nil {
			log.Println(err)
		}
	}

	// TODO make PBE configurable
	derivedKey, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		log.Println("failed to create PBE key with scrypt")
	}
	return derivedKey, salt
	// return pbkdf2.Key([]byte(password), salt, 4096, 32, sha256.New), salt
}

func aesEncrypt(password string, message []byte) []byte {
	salt := make([]byte, 8)
	key, salt := deriveAESKey(password, salt)

	// create the aes block
	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err)
	}

	// specify gcm block mode
	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		log.Println(err)
	}

	// init the initialisation vector
	initialisationVector := make([]byte, gcm.NonceSize())
	_, err = rand.Read(initialisationVector)
	if err != nil {
		log.Println(err)
	}

	// save the nonce as a prefix to the encrypted data
	return append(salt, gcm.Seal(initialisationVector, initialisationVector, []byte(message), nil)...)
}

func aesDecrypt(password string, saltAndCiphertext []byte) ([]byte, error) {
	salt := saltAndCiphertext[:8]
	ciphertext := saltAndCiphertext[8:]
	key, salt := deriveAESKey(password, salt)

	// create the aes block
	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err)
	}

	// specify gcm block mode
	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		log.Println(err)
	}

	// init the initialisation vector
	initialisationVector, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]

	return gcm.Open(nil, initialisationVector, ciphertext, nil)
}
