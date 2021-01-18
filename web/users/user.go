package users

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/ppartarr/tipsy/checkers/typtop"

	bolt "go.etcd.io/bbolt"
)

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
			// log.Fatal("user not found")
			return errors.New("no user with email " + email + " in bucket users")
		}

		err := json.Unmarshal(userBytes, &user)

		if err != nil {
			return err
		}

		return nil
	})

	log.Println("this should run")

	if err != nil {
		return nil, err
	}

	log.Println(user.Email)

	return user, nil
}

func (userService *UserService) createTypTopUser(user *typtop.TypTopUser) error {
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

func (userService *UserService) getTypTopUser(email string) (user *typtop.TypTopUser, err error) {
	user = &typtop.TypTopUser{}
	log.Println("getting typtop user for", email)

	err = userService.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))

		if bucket == nil {
			return errors.New("bucket users not found")
		}

		userBytes := bucket.Get([]byte(email))

		// TODO find out what's weird about this
		if userBytes == nil || len(userBytes) == 0 {
			// log.Fatal("user bytes null")
			return errors.New("no user with email " + email + " in bucket users")
		}

		err := json.Unmarshal(userBytes, &user)

		if err != nil {
			return err
		}

		return nil
	})

	log.Println("waiting for error")

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return user, nil
}
