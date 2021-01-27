package typtop

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"log"
	"math"
	"math/big"

	"github.com/ppartarr/tipsy/correctors"
)

// CacheState represents the state of the cache
type CacheState struct {
	// w
	RegisteredPassword string
	// F
	TypoCacheFrequency []int
	// M
	WaitListFrequency map[string]int
	// U
	TypoIndexPairs []TypoIndexPair
}

// TypoIndexPair represents a typo and the index of that typo in the typo cache
type TypoIndexPair struct {
	typo  string
	index int
}

func encryptCacheState(publicKey *rsa.PublicKey, cacheState *CacheState) []byte {
	encodedCacheState, err := json.Marshal(cacheState)
	if err != nil {
		log.Println("failed to encode the cache state into json")
	}
	return RSAEncrypt(publicKey, encodedCacheState)
}

func decryptCacheState(privateKey *rsa.PrivateKey, encryptedCacheState []byte) *CacheState {
	cacheState := &CacheState{}

	decryptedCacheState := RSADecrypt(privateKey, encryptedCacheState)

	err := json.Unmarshal(decryptedCacheState, &cacheState)
	if err != nil {
		log.Println("failed to decode the cache state from json")
	}

	return cacheState
}

func addPasswordsToTypoCache(pairs []TypoIndexPair, privateKey *rsa.PrivateKey, typoCache [][]byte) [][]byte {
	encodedPrivateKey := encodeKey(privateKey)
	for _, pair := range pairs {
		fmt.Println("adding password to typo cache: ", pair)
		typoCache[pair.index] = aesEncrypt(pair.typo, encodedPrivateKey)
	}
	return typoCache
}

func initTypoCache(typoCache [][]byte, privateKey *rsa.PrivateKey) [][]byte {
	encodedPrivateKey := encodeKey(privateKey)

	for i := 1; i < len(typoCache); i++ {
		// length between 10 and 26
		length, err := rand.Int(rand.Reader, big.NewInt(16))
		if err != nil {
			log.Println("failed to generate random int")
		}

		// generate random string
		ciphertext := generateRandomStringFromRunes(int(length.Int64()) + 10)

		log.Println(i)
		log.Println(ciphertext)

		// store in typo cache
		encryptedCiphertext := aesEncrypt(ciphertext, encodedPrivateKey)
		typoCache[i] = encryptedCiphertext
	}

	return typoCache
}

func generateRandomStringFromRunes(length int) string {
	b := make([]rune, length)
	for i := range b {
		randomLength, err := rand.Int(rand.Reader, big.NewInt(int64(len(correctors.LetterRunes))))
		if err != nil {
			log.Println("failed to generate random int")
		}
		b[i] = correctors.LetterRunes[randomLength.Int64()]
	}
	return string(b)
}

// CacheInit in Chatterjee et al.
func (checker *Checker) initCache(password string) *CacheState {
	typoCacheFrequency := make([]int, checker.config.TypoCache.Length)

	// TODO implement cache warm up
	for i := 0; i < checker.config.TypoCache.Length; i++ {
		typoCacheFrequency[i] = 0
	}

	// create typo based on correctors
	typoIndexPairs := make([]TypoIndexPair, 0)

	if checker.config.TypoCache.WarmUp {
		corrections := correctors.GetNBestCorrectors(checker.config.TypoCache.Length, checker.typoFrequency)
		// we start from 1 to avoid overwritting the submitted password in TypoCache[0]
		for i := 1; i < checker.config.TypoCache.Length; i++ {
			// TODO if > 10 fill 10 with best correctors and rest with random strings
			correctedPassword := correctors.ApplyCorrectionFunction(corrections[i], password)
			log.Println(correctedPassword)
			typoIndexPair := TypoIndexPair{
				typo:  correctedPassword,
				index: i,
			}
			typoIndexPairs = append(typoIndexPairs, typoIndexPair)
		}
	}

	state := &CacheState{
		RegisteredPassword: password,
		TypoCacheFrequency: typoCacheFrequency,
		TypoIndexPairs:     typoIndexPairs,
	}

	return state
}

// CacheUpdt in Chatterjee et al.
func (checker *Checker) updateCache(pi []int, cacheState *CacheState, typoIndexPair *TypoIndexPair, waitList [][]byte) *CacheState {
	// update the frequency of the submitted typo in typo cache
	if typoIndexPair.index > 0 {
		cacheState.TypoCacheFrequency[typoIndexPair.index]++
	}

	// iterate over decrypted wait list
	for j := 0; j < len(waitList); j++ {
		if string(waitList[j]) != "" && checker.valid(typoIndexPair.typo, string(waitList[j])) {
			log.Println("incrementing frequency of typo in wait list")
			// increment wait list frequency
			cacheState.WaitListFrequency[string(waitList[j])]++
		}
	}

	// sort frequencies in frequency wait list
	log.Println("sorting wait list frequency")
	// ss := correctors.ConvertMapToSortedSlice(cacheState.WaitListFrequency)

	// TODO iterate in decreasing order
	// iterate over wait listtypoCacheFrequency
	for typo, frequency := range cacheState.WaitListFrequency {
		if frequency > 0 {
			// k is the frequency of the least frequently used typo in the typo cache
			k := min(cacheState.TypoCacheFrequency)

			nu := calculateFrequency(cacheState.WaitListFrequency[typo], k)

			// update the typo cache
			if nu >= 0.5 {
				cacheState.TypoCacheFrequency[k] = cacheState.TypoCacheFrequency[k] + cacheState.WaitListFrequency[typo]

				// append new typo index pair
				newTypoIndexPair := &TypoIndexPair{
					typo:  typoIndexPair.typo,
					index: k,
				}
				cacheState.TypoIndexPairs = append(cacheState.TypoIndexPairs, *newTypoIndexPair)
			}
		}
	}

	// create new typo cache frequency list
	newTypoCacheFrequency := make([]int, len(cacheState.TypoCacheFrequency))
	for index, frequency := range cacheState.TypoCacheFrequency {
		newTypoCacheFrequency[pi[index]] = frequency
	}

	log.Println(newTypoCacheFrequency)

	cacheState.TypoCacheFrequency = newTypoCacheFrequency

	return cacheState
}

func calculateFrequency(waitListFrequency int, leastUsedTypoFrequency int) float64 {
	return float64(waitListFrequency) / float64(leastUsedTypoFrequency+waitListFrequency)
}

func min(values []int) int {
	if len(values) == 0 {
		return -1
	}

	min := math.MaxInt32
	for _, v := range values {
		if v < min {
			min = v
		}
	}

	return min
}
