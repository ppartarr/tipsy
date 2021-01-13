package users

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"

	bolt "go.etcd.io/bbolt"
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

func (userService *UserService) storeToken(token *Token) error {
	log.Println("storing token")
	return userService.db.Update(func(tx *bolt.Tx) error {
		// Retrieve the tokens bucket
		// This should be created when the DB is first opened.
		bucket := tx.Bucket([]byte("tokens"))

		// Generate ID for the user
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		id, _ := bucket.NextSequence()
		token.ID = int(id)

		// Marshal user data into bytes
		buf, err := json.Marshal(token)
		if err != nil {
			return err
		}

		// Persist bytes to users bucket
		return bucket.Put([]byte(token.Token), buf)
	})
}

func (userService *UserService) deleteToken(token *Token) error {
	log.Println("deleting token")
	return userService.db.Update(func(tx *bolt.Tx) error {
		// Retrieve the tokens bucket
		// This should be created when the DB is first opened.
		bucket := tx.Bucket([]byte("tokens"))

		// delete token from bucket
		return bucket.Delete([]byte(token.Token))
	})
}

func (userService *UserService) getToken(tokenHash string) (token *Token, err error) {
	token = &Token{}
	log.Println("getting token for", tokenHash)

	err = userService.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("tokens"))

		if bucket == nil {
			return errors.New("bucket tokens not found")
		}

		tokenBytes := bucket.Get([]byte(tokenHash))

		if tokenBytes == nil || len(tokenBytes) == 0 {
			return errors.New("no token for token hash " + tokenHash + " in bucket tokens")
		}

		err := json.Unmarshal(tokenBytes, &token)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
