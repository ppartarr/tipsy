package typtop

import (
	"crypto/rsa"
	"log"
	mrand "math/rand"

	"github.com/hbollon/go-edlib"
	"github.com/nbutton23/zxcvbn-go"
	"gonum.org/v1/gonum/stat/combin"
)

// Login checks if the typtop user's password
func (checker *Checker) Login(state *State, submittedPassword string, privateKey *rsa.PrivateKey) (bool, *State) {
	success := false

	// decrypt typo cache
	for i := 0; i < len(state.TypoCache); i++ {
		encodedPrivateKey, err := aesDecrypt(submittedPassword, state.TypoCache[i])

		// if the decoded & decrypted private key matches
		if err == nil && privateKey.Equal(decodeKey(encodedPrivateKey)) {
			log.Println("private keys match!")
			success = true

			// get random permutation
			permutations := generatePermutations(checker.config.TypoCache.Length)
			pi := permutations[mrand.Intn(len(permutations))]

			// decrypt ciphered cache state
			log.Println("decrypting cache state")
			cacheState := decryptCacheState(privateKey, state.CipheredCacheState)
			log.Println("cache state")

			// decrypt typos from the wait list
			log.Println("decrypting wait list")
			typos := decryptWaitList(privateKey, state.WaitList)
			log.Println("wait list: ", typos)

			typoIndexPair := &TypoIndexPair{
				typo:  submittedPassword,
				index: i,
			}

			// update cache
			cacheState = checker.updateCache(pi, cacheState, typoIndexPair, typos)

			// encrypt cache state
			log.Println("encrypting cache state")
			encryptedCacheState := encryptCacheState(&privateKey.PublicKey, cacheState)

			// randomize typo order in typo cache
			newTypoCache := make([][]byte, len(state.TypoCache))
			for index, typo := range state.TypoCache {
				newTypoCache[pi[index]] = typo
			}

			// clear the wait list (same as init)
			log.Println("clear the wait list")
			initWaitList(state.WaitList, &privateKey.PublicKey)

			// update state
			state.CipheredCacheState = encryptedCacheState
			state.TypoCache = newTypoCache

			return success, state
		}
	}

	if !success {
		// add typo to wait list
		state.WaitList[state.gamma] = RSAEncrypt(&privateKey.PublicKey, []byte(submittedPassword))

		// increment gamma
		state.gamma = state.gamma + 1%len(state.WaitList)
	}

	return success, state
}

// GeneratePermutations given an int, will generate every permutation of set of that length
// e.g. given 3 will return [[0, 1, 2], [0, 2, 1], [1, 0, 2], [1, 2, 0], [2, 0, 1], [2, 1, 0]]
func generatePermutations(length int) [][]int {
	return combin.Permutations(length, length)
}

// validates that a typo is an acceptable password submission
// 1. check damereau-levenshtein distance < 2
// 2. check that the password strength
func (checker *Checker) valid(password string, typo string) bool {

	// TODO make this configurable
	distance := edlib.DamerauLevenshteinDistance(password, typo)
	if distance <= checker.config.EditDistance {
		log.Println("damereau levenshtein distance is too large")
		return false
	}

	typoStrengthEstimation := zxcvbn.PasswordStrength(typo, nil)
	passwordStrengthEstimation := zxcvbn.PasswordStrength(password, nil)
	// TODO use guessability here instead
	if typoStrengthEstimation.Score < checker.config.Zxcvbn {
		log.Println(typoStrengthEstimation.Score)
		log.Println("typo strength estimation is too low to add the typo to the typo", typoStrengthEstimation.Score)
		return false
	}

	if typoStrengthEstimation.Score < passwordStrengthEstimation.Score {
		log.Println("typo strength estimation is lower than the registered password's")
		return false
	}

	return true
}
