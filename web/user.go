package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ppartarr/tipsy/web/session"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
)

// Service
type UserService struct {
	db *bolt.DB
}

// NewService returns a new Service instance
func NewUserService(db *bolt.DB) (userService *UserService) {
	// create users bucket
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return &UserService{
		db: db,
	}
}

// Login allows a user to login to their account
func (userService *UserService) Login(w http.ResponseWriter, r *http.Request) (user *User, err error) {
	log.Println("method: ", r.Method)

	// TODO handle HEAD, PUT, and PATCH separately
	if r.Method != "POST" {
		http.Redirect(w, r, "/login.html", 301)
		return
	}

	r.ParseForm()
	log.Println("username: ", r.Form["username"])
	log.Println("password: ", r.Form["password"])

	// get username
	user, err = userService.getUser(r.Form["username"][0])

	if err != nil {
		return nil, fmt.Errorf("could not get user %q: %q", r.Form["username"][0], err)
	}

	if !CheckPasswordHash(r.Form["password"][0], user.PasswordHash) {
		return nil, fmt.Errorf("wrong password for user %q", r.Form["username"][0])
	}

	_, err = session.SetUserID(w, r, strconv.Itoa(user.ID))
	if err != nil {
		return nil, fmt.Errorf("could not create a session for user %q: %q", r.Form["username"][0], err)
	}

	return
}

// Register allows a user to register a new account
func (userService *UserService) Register(w http.ResponseWriter, r *http.Request) (user *User, err error) {
	// TODO handle HEAD, PUT, and PATCH separately
	if r.Method != "POST" {
		http.Redirect(w, r, "/login.html", 301)
		return
	}

	r.ParseForm()
	log.Println(r.Form)

	passwordHash, err := HashPassword(r.Form["password"][0])
	if err != nil {
		log.Println(err.Error())
	}

	// create new user from request then save in db
	user = &User{
		Username:     r.Form["username"][0],
		PasswordHash: passwordHash,
	}
	CreateUser(user, userService.db)

	return user, nil
}

// Logout logs a user out
func (userService *UserService) Logout(w http.ResponseWriter, r *http.Request) (user *User, err error) {
	err = session.Destroy(w, r)
	if err != nil {
		return nil, fmt.Errorf("could not destroy session")
	}

	return
}

// CreateUser saves a user in the db
func CreateUser(user *User, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		bucket := tx.Bucket([]byte("users"))

		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		id, _ := bucket.NextSequence()
		user.ID = int(id)

		// Marshal user data into bytes.
		buf, err := json.Marshal(user)
		log.Println(buf)
		if err != nil {
			return err
		}

		// Persist bytes to users bucket.
		return bucket.Put([]byte(user.Username), buf)
	})
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

func (userService *UserService) getUser(username string) (user *User, err error) {
	user = &User{}
	log.Println("getting user for", username)

	err = userService.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))

		if bucket == nil {
			return fmt.Errorf("bucket %q not found", "users")
		}

		userBytes := bucket.Get([]byte(username))

		if len(userBytes) == 0 {
			return fmt.Errorf("no user with username %q in bucket %q", username, "users")
		}

		err := json.Unmarshal(userBytes, &user)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}
