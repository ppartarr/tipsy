package typtop

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"log"
)

func initWaitList(waitList [][]byte, publicKey *rsa.PublicKey) [][]byte {
	for i := 0; i < len(waitList); i++ {
		waitList[i] = RSAEncrypt(publicKey, []byte(""))
	}
	return waitList
}

func decryptWaitList(privateKey *rsa.PrivateKey, waitList [][]byte) [][]byte {
	typos := make([][]byte, len(waitList))

	// update wait list
	for j := 1; j < len(waitList); j++ {
		// public key decryption of wait list j
		typo, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, waitList[j], nil)
		if err != nil {
			log.Fatal("failed to decrypt wait list typo")
		}
		typos = append(typos, typo)
	}
	return typos
}
