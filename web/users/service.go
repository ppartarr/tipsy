package users

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/ppartarr/tipsy/checkers/typtop"
	"github.com/ppartarr/tipsy/config"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
)

// UserService represents the user service
type UserService struct {
	db     *bolt.DB
	config *config.Server
}

// User represents a user
type User struct {
	ID            int    `json:"id"`
	Email         string `json:"email"`
	PasswordHash  string `json:"password"`
	LoginAttempts int    `json:"loginAttempts"`
}

// Token represents a password reset token
type Token struct {
	ID        int           `json:"id"`
	Email     string        `json:"email"`
	Token     string        `json:"token"`
	CreatedAt time.Time     `json:"createdAt"`
	TTL       time.Duration `json:"ttl"`
}

// NewUserService returns a new UserService instance
func NewUserService(db *bolt.DB, tipsyConfig *config.Server) (userService *UserService) {
	db.Update(func(tx *bolt.Tx) error {
		// create users bucket
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return errors.New("create bucket: " + err.Error())
		}
		// create tokens bucket for password reset
		_, err = tx.CreateBucketIfNotExists([]byte("tokens"))
		if err != nil {
			return errors.New("create bucket: " + err.Error())
		}
		if tipsyConfig.Checker.TypTop != nil {
			// create typtop users bucket
			_, err = tx.CreateBucketIfNotExists([]byte("typtopUsers"))
			if err != nil {
				return errors.New("create bucket: " + err.Error())
			}
		}
		return nil
	})

	return &UserService{
		db:     db,
		config: tipsyConfig,
	}
}

// HashPassword hashes a password with bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash verifies that the password arguments matches the given hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (userService *UserService) incrementLoginAttempts(user *User) error {
	user.LoginAttempts++
	log.Println("incrementing user attempts: ", user.LoginAttempts)

	return userService.db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket
		// This should be created when the DB is first opened.
		bucket := tx.Bucket([]byte("users"))

		if bucket == nil {
			return errors.New("bucket users not found")
		}

		// Marshal user data into bytes
		buf, err := json.Marshal(user)
		if err != nil {
			return err
		}

		// Persist bytes to users bucket
		return bucket.Put([]byte(user.Email), buf)
	})
}

func (userService *UserService) incrementTypTopLoginAttempts(user *typtop.User) error {
	user.LoginAttempts++
	log.Println("incrementing user attempts: ", user.LoginAttempts)

	return userService.db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket
		// This should be created when the DB is first opened.
		bucket := tx.Bucket([]byte("users"))

		if bucket == nil {
			return errors.New("bucket users not found")
		}

		// Marshal user data into bytes
		buf, err := json.Marshal(user)
		if err != nil {
			return err
		}

		// Persist bytes to users bucket
		return bucket.Put([]byte(user.Email), buf)
	})
}
