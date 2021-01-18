package typtop

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"log"
	"math/big"

	"golang.org/x/crypto/pbkdf2"
)

// State represents the server state in Chatterjee et al. typo correction paper
type State struct {
	Publickey          *rsa.PublicKey
	CipheredCacheState []byte
	TypoCache          [][]byte
	WaitList           [][]byte
	gamma              int
}

// CacheState represents the state of the cache
type CacheState struct {
	RegisteredPassword string
	TypoCacheFrequency []int
	TypoIndexPairs     map[string]int
}

// Register initialises the user's state
func (checker *TypTopCheckerService) Register(password string) (state *State) {
	// generate key
	privateKey, err := rsa.GenerateKey(rand.Reader, checker.config.PublicKeyEncryption.KeyLength)
	if err != nil {
		log.Printf("Failed to generate RSA key\n")
	}
	publicKey := &privateKey.PublicKey

	// TODO generate different salt for every user & store
	salt := []byte("salt")

	typoCache := make([][]byte, checker.config.TypoCache.Length)

	// encrypt the private key using pbkdf and add resulting ciphertext to the typocache
	typoCache[0] = pbkdf2.Key([]byte(password), salt, 4096, 32, sha1.New)

	// fill typo cache with random ciphertext
	initTypoCache(typoCache)

	waitList := make([][]byte, checker.config.WaitList.Length)

	// fill wait list with empty string
	initWaitList(waitList, state.Publickey)

	// init cache
	cacheState, typoIndexPairs := checker.initCache(password)

	// encrypt cache state
	encryptedCacheState := encryptCacheState(publicKey, cacheState)

	// iterate over the typoIndexPairs
	addPasswordsToTypoCache(typoIndexPairs, privateKey, typoCache)

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

	return state
}

func encryptCacheState(publicKey *rsa.PublicKey, cacheState *CacheState) []byte {
	encodedCacheState, err := json.Marshal(cacheState)
	if err != nil {
		log.Println("failed to encode the cache state into json")
	}
	encryptedCacheState, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, encodedCacheState, nil)
	if err != nil {
		log.Fatal("failed to encrypt cache state")
	}
	return encryptedCacheState
}

func addPasswordsToTypoCache(pairs map[string]int, privateKey *rsa.PrivateKey, typoCache [][]byte) [][]byte {
	for _, index := range pairs {
		// encode private key
		encodedPrivateKey, err := json.Marshal(privateKey)
		if err != nil {
			log.Fatal("failed to encode private key")
		}
		encryptedSecretKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &privateKey.PublicKey, encodedPrivateKey, nil)
		typoCache[index] = encryptedSecretKey
	}
	return typoCache
}

func initTypoCache(typoCache [][]byte) [][]byte {
	for i := 1; i < len(typoCache); i++ {
		// length between 10 and 26
		length, err := rand.Int(rand.Reader, big.NewInt(16))
		if err != nil {
			log.Println("failed to generate random int")
		}

		// generate random string
		log.Println("generating random string")
		ciphertext := generateRandomStringFromRunes(int(length.Int64()) + 10)

		// TODO get good random salt of at least 8 bytes
		salt := []byte("salt")

		// store in typo cache
		typoCache[i] = pbkdf2.Key([]byte(ciphertext), salt, 4096, 32, sha1.New)
	}

	return typoCache
}

// TODO extend to use any rune instead of US alphanumerics
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`1234567890-=[]\\;',./~!@#$%^&*()_+{}|:\"<>?")

func generateRandomStringFromRunes(length int) string {
	log.Println(length)
	b := make([]rune, length)
	for i := range b {
		randomLength, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterRunes))))
		if err != nil {
			log.Println("failed to generate random int")
		}
		b[i] = letterRunes[randomLength.Int64()]
	}
	return string(b)
}

func initWaitList(waitList [][]byte, publicKey *rsa.PublicKey) [][]byte {
	for i := 0; i < len(waitList); i++ {
		encryptedEpsilon, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, nil, nil)
		if err != nil {
			log.Fatal("failed encrypt the empty message in the wait list")
		}
		waitList[i] = encryptedEpsilon
	}
	return waitList
}

func (checker *TypTopCheckerService) initCache(password string) (*CacheState, map[string]int) {
	typoCacheFrequency := make([]int, checker.config.TypoCache.Length)

	// TODO implement cache warm up
	for i := 0; i < checker.config.TypoCache.Length; i++ {
		typoCacheFrequency[i] = 0
	}

	// init the typo index pairs & the cache state
	typoIndexPairs := make(map[string]int)

	state := &CacheState{
		RegisteredPassword: password,
		TypoCacheFrequency: typoCacheFrequency,
	}

	return state, typoIndexPairs
}

// func warmTypoCache([][]byte) [][]byte {

// 	// get n best correctors
// 	nBestCorrectors := correctors.GetNBestCorrectors(len(checker.Correctors), checker.TypoFrequency)

// }
