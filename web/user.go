package web

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/mcnijman/go-emailaddress"
	"github.com/ppartarr/tipsy/config"
	"github.com/ppartarr/tipsy/mail"
	"github.com/ppartarr/tipsy/web/session"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
)

// Service
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

// NewService returns a new Service instance
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
		return nil
	})

	return &UserService{
		db:     db,
		config: tipsyConfig,
	}
}

// Login allows a user to login to their account
func (userService *UserService) Login(w http.ResponseWriter, r *http.Request) error {
	// TODO handle HEAD, PUT, and PATCH separately
	if r.Method != "POST" {
		http.Redirect(w, r, "/login.html", 301)
		return errors.New("must use POST http method")
	}

	// check that a valid form was submitted
	r.ParseForm()
	log.Println(r.Form)
	expectedValues := []string{"email", "password"}
	if !formIsValid(r.Form, expectedValues) {
		return errors.New("you must submit a valid form")
	}

	// get user from email
	user, err := userService.getUser(r.Form["email"][0])
	if err != nil {
		return errors.New("could not get user " + r.Form["email"][0] + ": " + err.Error())
	}

	// validate email
	log.Println("validating email address")
	_, err = emailaddress.Parse(r.Form["email"][0])
	if err != nil {
		return errors.New("invalid email: %q" + r.Form["email"][0] + ": " + err.Error())
	}

	// TODO plug-in checkers here
	// check password
	log.Println("checking password")
	if !CheckPasswordHash(r.Form["password"][0], user.PasswordHash) {

		// increment login attempts
		if user.LoginAttempts < userService.config.RateLimit {
			userService.incrementLoginAttempts(user)
		} else {
			log.Println("user has attempted to login too many times")
			go func() {
				userService.PasswordRecovery(w, r)
			}()
			log.Println("this should return...")
			return errors.New("the limit of login attempts has been reached, please reset your password via the mail provided in the link " + r.Form["email"][0])
		}

		return errors.New("wrong password for user " + r.Form["email"][0])
	}

	// init session
	_, err = session.SetUserID(w, r, strconv.Itoa(user.ID))
	if err != nil {
		return errors.New("could not create a session for user " + r.Form["email"][0] + ": " + err.Error())
	}

	return nil
}

// Register allows a user to register a new account
func (userService *UserService) Register(w http.ResponseWriter, r *http.Request) error {
	// TODO handle HEAD, PUT, and PATCH separately
	if r.Method != "POST" {
		http.Redirect(w, r, "/login.html", 301)
		return errors.New("must use POST http method")
	}

	// check that a valid form was submitted
	r.ParseForm()
	log.Println(r.Form)
	expectedValues := []string{"email", "password", "password-copy"}
	if !formIsValid(r.Form, expectedValues) {
		return errors.New("you must submit a valid form")
	}

	// check that password & password copy match
	if r.Form["password"][0] != r.Form["password-copy"][0] {
		return errors.New("password and password copy don't match")
	}

	// check that email isn't already registered
	user, err := userService.getUser(r.Form["email"][0])
	if user != nil {
		return errors.New("user already with email " + r.Form["email"][0] + " is already registered")
	}

	// TODO only register password if it is 3/4 in zxcbn-go

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

	return nil
}

// Logout logs a user out
func (userService *UserService) Logout(w http.ResponseWriter, r *http.Request) error {
	err := session.Destroy(w, r)
	if err != nil {
		return errors.New("could not destroy session")
	}

	return nil
}

// PasswordRecovery sends the user an email containing a password reset link
func (userService *UserService) PasswordRecovery(w http.ResponseWriter, r *http.Request) error {
	// check that a valid form was submitted
	r.ParseForm()
	log.Println(r.Form)
	expectedValues := []string{"email"}
	if !formIsValid(r.Form, expectedValues) {
		return errors.New("you must submit a valid form")
	}

	// get user from db using submitted email
	user, err := userService.getUser(r.Form["email"][0])
	if err != nil {
		return errors.New("could not get user" + r.Form["email"][0] + ": " + err.Error())
	}

	log.Println("there")

	log.Println(userService.config.TokenValidity)
	log.Println(userService.config.TokenValidity * time.Minute)

	token := &Token{
		Email:     r.Form["email"][0],
		TTL:       userService.config.TokenValidity,
		CreatedAt: time.Now().Local(),
	}

	// generate token
	token.Token = generateToken(user.Email)

	// store the token
	log.Println(token)
	err = userService.storeToken(token)
	if err != nil {
		return errors.New("could not store token for user " + r.Form["email"][0] + ": " + err.Error())
	}

	// send password reset mail
	go func() {
		log.Println("sending password reset email")
		// i hope this is safe
		url := "http://localhost:8000/reset.html?token=" + token.Token
		mail.Send(user.Email, "Tipsy password reset", mail.GeneratePasswordResetMail(url))
	}()

	return nil
}

func formIsValid(form url.Values, expectedValues []string) bool {
	for _, value := range expectedValues {
		if len(form) == 0 || len(form[value]) == 0 {
			return false
		}
	}

	return true
}

// PasswordReset validates the token then updates the user's password
func (userService *UserService) PasswordReset(w http.ResponseWriter, r *http.Request) error {

	// check that a valid form was submitted
	r.ParseForm()
	log.Println(r.Form)
	expectedValues := []string{"email", "password", "password-copy"}
	if !formIsValid(r.Form, expectedValues) {
		return errors.New("you must submit a valid form")
	}

	// check that password & password copy match
	if r.Form["password"][0] != r.Form["password-copy"][0] {
		return errors.New("password and password copy don't match")
	}

	// get token hash from form & get associated token from bucket
	tokenHash := r.Form["token"][0]
	token, err := userService.getToken(tokenHash)
	if err != nil {
		return errors.New("could not get token")
	}

	log.Println(token)

	// check that the token matches the submitted user
	if r.Form["email"][0] != token.Email {
		return errors.New("submitted email & token email don't match")
	}

	// check if the TTL is still valid
	log.Println(token.TTL)
	log.Println(token.TTL * time.Minute)
	log.Println(token.CreatedAt.Add(token.TTL))
	log.Println(time.Now().Local())
	log.Println(token.CreatedAt.Add(token.TTL).Before(time.Now().Local()))
	if token.CreatedAt.Add(token.TTL).Before(time.Now().Local()) {
		return errors.New("the token expired at " + token.CreatedAt.Add(token.TTL).String())
	}

	// TODO invalidate token & replace old password hash with new password hash
	err = userService.deleteToken(token)
	if err != nil {
		return errors.New("error deleting the token" + err.Error())
	}

	// TODO improve this updatePassword -> updateUser
	err = userService.updatePassword(token.Email, r.Form["password"][0])

	log.Println("password has been reset")

	return nil
}

func (userService *UserService) updatePassword(email string, password string) error {
	// get user from db
	user, err := userService.getUser(email)
	if err != nil {
		return errors.New("error getting the user: " + err.Error())
	}

	// hash the input password
	passwordHash, err := HashPassword(password)
	if err != nil {
		return errors.New("failed to hash the input password: " + err.Error())
	}

	// set new password & set login attempts to 0
	user.PasswordHash = passwordHash
	user.LoginAttempts = 0

	err = userService.updateUser(user)
	if err != nil {
		return errors.New("failed to update the user's password: " + err.Error())
	}

	return nil
}

func (userService *UserService) updateUser(user *User) error {
	log.Println("updating user for", user.Email)

	return userService.db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket
		// This should be created when the DB is first opened.
		bucket := tx.Bucket([]byte("users"))

		// Marshal user data into bytes
		buf, err := json.Marshal(user)
		if err != nil {
			return err
		}

		// Persist bytes to users bucket
		return bucket.Put([]byte(user.Email), buf)
	})
}

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
			return errors.New("bucket users not found")
		}

		userBytes := bucket.Get([]byte(email))

		if userBytes == nil || len(userBytes) == 0 {
			return errors.New("no user with email " + email + " in bucket users")
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
	log.Println(user.LoginAttempts)

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
