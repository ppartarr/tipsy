package typtop

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"math/big"
)

// State represents the server state in Chatterjee et al. typo correction paper
type State struct {
	Publickey          *rsa.PublicKey
	CipheredCacheState []byte
	TypoCache          [][]byte
	WaitList           [][]byte
	gamma              int
}

// Register initialises the user's state
func (checker *CheckerService) Register(password string) (state *State, privateKey *rsa.PrivateKey) {
	// generate key
	privateKey, err := rsa.GenerateKey(rand.Reader, checker.config.PublicKeyEncryption.KeyLength)
	if err != nil {
		log.Printf("Failed to generate RSA key\n")
	}
	publicKey := &privateKey.PublicKey

	typoCache := make([][]byte, checker.config.TypoCache.Length)
	waitList := make([][]byte, checker.config.WaitList.Length)

	// encode the private key and save in typo cache
	typoCache[0] = aesEncrypt(password, encodeKey(privateKey))

	// fill typo cache with random ciphertext
	initTypoCache(typoCache, publicKey)

	// fill wait list with empty string
	initWaitList(waitList, publicKey)

	// TODO warm up cache
	// init cache
	cacheState := checker.initCache(password)

	// encrypt cache state
	encryptedCacheState := encryptCacheState(publicKey, cacheState)

	// TODO add back in when warming up cache
	// iterate over the typoIndexPairs
	// addPasswordsToTypoCache(emptyTypoIndexPairs, privateKey, typoCache)

	// get random number between 0 and t
	gamma, err := rand.Int(rand.Reader, big.NewInt(int64(checker.config.TypoCache.Length)))
	if err != nil {
		log.Println("failed to generate random int")
	}

	state = &State{
		Publickey:          publicKey,
		CipheredCacheState: encryptedCacheState,
		TypoCache:          typoCache,
		WaitList:           waitList,
		gamma:              int(gamma.Int64()),
	}

	return state, privateKey
}
