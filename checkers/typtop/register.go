package typtop

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"math/big"
)

// State represents the server state in Chatterjee et al. typo correction paper
type State struct {
	Publickey          *rsa.PublicKey `json:"publicKey"`
	CipheredCacheState []byte         `json:"cipheredCacheState"`
	TypoCache          [][]byte       `json:"typoCache"`
	WaitList           [][]byte       `json:"waitList"`
	Gamma              int            `json:"gamma"`
}

// Register initialises the user's state
func (checker *Checker) Register(password string) (state *State, privateKey *rsa.PrivateKey) {
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
	log.Println("init the typo cache")
	typoCache = initTypoCache(typoCache, privateKey)

	// fill wait list with empty string
	log.Println("init the wait list")
	waitList = initWaitList(waitList, publicKey)

	// init cache
	log.Println("init the cache state")
	cacheState := checker.initCache(password)

	// encrypt cache state
	log.Println("encrypt cache state")
	encryptedCacheState := encryptCacheState(publicKey, cacheState)

	// warm up the cache
	log.Println("warm up typo cache")
	if checker.config.TypoCache.WarmUp && len(cacheState.TypoIndexPairs) != 0 {
		typoCache = addPasswordsToTypoCache(cacheState.TypoIndexPairs, privateKey, typoCache)
	}

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
		Gamma:              int(gamma.Int64()),
	}

	return state, privateKey
}
