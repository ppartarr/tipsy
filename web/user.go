package web

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"

	bolt "go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
)

// Service
type UserService struct {
	db *bolt.DB
}

// NewService returns a new Service instance
func NewUserService(db *bolt.DB) (userService *UserService) {
	return &UserService{
		db: db,
	}
}

// Login allows a user to login to their account
func (userService *UserService) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method: ", r.Method)

	// TODO handle HEAD, PUT, and PATCH separately
	if r.Method != "POST" {
		http.Redirect(w, r, "/login.html", 301)
		return
	}

	r.ParseForm()
	fmt.Println("username: ", r.Form["username"])
	fmt.Println("password: ", r.Form["password"])
}

// Register allows a user to register a new account
func (userService *UserService) Register(w http.ResponseWriter, r *http.Request) {
	// TODO handle HEAD, PUT, and PATCH separately
	if r.Method != "POST" {
		http.Redirect(w, r, "/login.html", 301)
		return
	}

	r.ParseForm()
	fmt.Println(r.Form)

	passwordHash, err := HashPassword(r.Form["password"][0])
	if err != nil {
		fmt.Println(err.Error())
	}

	// create new user from request then save in db
	user := &User{
		Username:     r.Form["username"][0],
		PasswordHash: passwordHash,
	}
	CreateUser(user, userService.db)
}

// CreateUser saves a user in the db
func CreateUser(user *User, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket([]byte("users"))

		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		id, _ := b.NextSequence()
		user.ID = int(id)

		// Marshal user data into bytes.
		buf, err := json.Marshal(user)
		if err != nil {
			return err
		}

		// Persist bytes to users bucket.
		return b.Put(itob(user.ID), buf)
	})
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// User represents a user
type User struct {
	ID           int    `json:"id" storm:"id,increment"`
	Username     string `json:"username" storm:"unique"`
	PasswordHash string `json:"password"`
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
