package users

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/ppartarr/tipsy/checkers/typtop"
)

// RegistrationForm represents a login form
type RegistrationForm struct {
	Email        string
	Password     string
	PasswordCopy string
	Errors       map[string]string
}

var rxEmail = regexp.MustCompile(".+@.+\\..+")

// Validate checks that the fields in the login form are set
func (form *RegistrationForm) Validate(submittedForm url.Values) bool {
	form.Errors = make(map[string]string)

	// check that username is an email address
	match := rxEmail.Match([]byte(form.Email))
	if match == false {
		form.Errors["Email"] = "Email must be valid"
	}

	if strings.TrimSpace(form.Password) == "" {
		form.Errors["Password"] = "Password cannot be empty"
	}

	// check that password & password copy match
	if submittedForm["password"][0] != submittedForm["password-copy"][0] {
		form.Errors["Password"] = "Passwords should match"
	}

	// TODO only register password if it is 3/4 in zxcbn-go

	return len(form.Errors) == 0
}

// Register allows a user to register a new account
func (userService *UserService) Register(w http.ResponseWriter, r *http.Request) (form *RegistrationForm, err error) {
	// TODO handle HEAD, PUT, and PATCH separately
	if r.Method != "POST" {
		http.Redirect(w, r, "/login.html", 301)
		return nil, errors.New("must use POST http method")
	}

	// check that a valid form was submitted
	r.ParseForm()
	log.Println(r.Form)
	form = &RegistrationForm{
		Email:        r.PostFormValue("email"),
		Password:     r.PostFormValue("password"),
		PasswordCopy: r.PostFormValue("password-copy"),
	}

	// validate form
	if form.Validate(r.Form) == false {
		log.Println(form.Errors)
		return form, errors.New("you must submit a valid form")
	}

	if userService.config.Checker.TypTop != nil {
		log.Println("running typtop mode")

		// check that email isn't already registered
		typtopUser, err := userService.getTypTopUser(r.Form["email"][0])
		if typtopUser != nil {
			return nil, errors.New("typtop user already with email " + r.Form["email"][0] + " is already registered")
		}

		// init the checker service
		log.Println("init checker service")
		checkerService := typtop.NewTypTopCheckerService(userService.config.Checker.TypTop)

		// register the password for typtop
		typtopState := checkerService.Register(r.Form["password"][0])

		// create new user from request then save in db
		typtopUser = &typtop.TypTopUser{
			Email:         r.Form["email"][0],
			LoginAttempts: 0,
			State:         *typtopState,
		}

		err = userService.createTypTopUser(typtopUser)
		if err != nil {
			return nil, errors.New("couldn't create typtop user")
		}
	} else {

		// check that email isn't already registered
		user, err := userService.getUser(r.Form["email"][0])
		if user != nil {
			return nil, errors.New("user already with email " + r.Form["email"][0] + " is already registered")
		}

		passwordHash, err := HashPassword(r.Form["password"][0])
		if err != nil {
			return nil, errors.New("couldn't hash password: " + err.Error())
		}

		// create new user from request then save in db
		user = &User{
			Email:         r.Form["email"][0],
			PasswordHash:  passwordHash,
			LoginAttempts: 0,
		}

		err = userService.createUser(user)
		if err != nil {
			return nil, errors.New("couldn't create user")
		}
	}

	// TODO return something other than form?
	return form, nil
}
