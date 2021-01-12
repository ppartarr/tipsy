package web

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/mcnijman/go-emailaddress"
	"github.com/ppartarr/tipsy/mail"
	"github.com/ppartarr/tipsy/web/session"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
)

// Service
type UserService struct {
	db *bolt.DB
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

// NewService returns a new Service instance
func NewUserService(db *bolt.DB) (userService *UserService) {
	db.Update(func(tx *bolt.Tx) error {
		// create users bucket
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		// create tokens bucket for password reset
		_, err := tx.CreateBucketIfNotExists([]byte("token"))
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
	log.Println("email: ", r.Form["email"])
	log.Println("password: ", r.Form["password"])

	// get user from email
	user, err = userService.getUser(r.Form["email"][0])

	if err != nil {
		return nil, fmt.Errorf("could not get user %q: %q", r.Form["email"][0], err)
	}

	// validate email
	_, err = emailaddress.Parse(r.Form["email"][0])
	if err != nil {
		return nil, fmt.Errorf("invalid email: %q", r.Form["email"][0])
	}

	// check password
	if !CheckPasswordHash(r.Form["password"][0], user.PasswordHash) {

		// increment login attempts
		if user.LoginAttempts < 10 {
			userService.incrementLoginAttempts(user)
		} else {
			userService.PasswordReset(w, r)
			return nil, fmt.Errorf("the limit of login attempts has been reached, please reset your password via the mail provided in the link %q", r.Form["email"][0])
		}

		return nil, fmt.Errorf("wrong password for user %q", r.Form["email"][0])
	}

	// init session
	_, err = session.SetUserID(w, r, strconv.Itoa(user.ID))
	if err != nil {
		return nil, fmt.Errorf("could not create a session for user %q: %q", r.Form["email"][0], err)
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

	// TODO only register password if it is 3/4 in zxcbn-go

	// TODO check that the email isn't already registered

	r.ParseForm()
	log.Println(r.Form)

	passwordHash, err := HashPassword(r.Form["password"][0])
	if err != nil {
		log.Println(err.Error())
	}

	// create new user from request then save in db
	user = &User{
		Email:         r.Form["email"][0],
		PasswordHash:  passwordHash,
		LoginAttempts: 0,
	}

	userService.createUser(user)

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

func (userService *UserService) PasswordReset(w http.ResponseWriter, r *http.Request) (token *Token, err error) {
	// get user from email
	user, err := userService.getUser(r.Form["email"][0])

	if err != nil {
		return nil, fmt.Errorf("could not get user %q: %q", r.Form["email"][0], err)
	}

	token = &Token{
		Email: r.Form["email"][0],
		TTL:   5 * time.Minute,
	}

	// generate token
	token.Token = generateToken(user.Email)

	// store the token
	err = userService.storeToken(token)
	if err != nil {
		return nil, fmt.Errorf("could not store token for user %q: %q", r.Form["email"][0], err)
	}

	// send password reset mail
	go func() {
		mail.Send(user.Email, "Tipsy password reset", mail.GeneratePasswordResetMail(token, url))
	}()
}

// hash email using bcrypt then hash again with md5 so that hash can be used as token in a url
func generateToken(email string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Hash to store:", string(hash))

	// TODO store hash
	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil))
}

func (userService *UserService) storeToken(token *Token) error {
	return userService.db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket
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
		return bucket.Put([]byte(token.Email), buf)
	})
}

func (userService *UserService) getToken(email string) (token *Token, err error) {
	token = &Token{}
	log.Println("getting token for", email)

	err = userService.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("tokens"))

		if bucket == nil {
			return fmt.Errorf("bucket %q not found", "users")
		}

		tokenBytes := bucket.Get([]byte(email))

		if len(tokenBytes) == 0 {
			return fmt.Errorf("no token for user with email %q in bucket %q", email, "users")
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

func (userService *UserService) createUser(user *User) error {
	return userService.db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket
		// This should be created when the DB is first opened.
		bucket := tx.Bucket([]byte("users"))

		// Generate ID for the user
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		id, _ := bucket.NextSequence()
		user.ID = int(id)

		// Marshal user data into bytes
		buf, err := json.Marshal(user)
		if err != nil {
			return err
		}

		// Persist bytes to users bucket
		return bucket.Put([]byte(user.Email), buf)
	})
}

func (userService *UserService) getUser(email string) (user *User, err error) {
	user = &User{}
	log.Println("getting user for", email)

	err = userService.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))

		if bucket == nil {
			return fmt.Errorf("bucket %q not found", "users")
		}

		userBytes := bucket.Get([]byte(email))

		if len(userBytes) == 0 {
			return fmt.Errorf("no user with email %q in bucket %q", email, "users")
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

	return userService.db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket
		// This should be created when the DB is first opened.
		bucket := tx.Bucket([]byte("users"))

		if bucket == nil {
			return fmt.Errorf("bucket %q not found", "users")
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
