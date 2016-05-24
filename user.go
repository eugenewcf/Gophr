package main

import "golang.org/x/crypto/bcrypt"

const (
	passwordLength = 8
	hashCost       = 10
	userIDLength
)

type User struct {
	Username     string
	Email        string
	HashPassword string
	ID           string
}

func NewUser(username, email, password string) (User, []error) {
	user := User{
		Username: username,
		Email:    email,
	}

	var errs []error

	if username == "" {
		// return user, errNoUsername
		errs = append(errs, errNoUsername)
	}

	if email == "" {
		// return user, errNoEmail
		errs = append(errs, errNoEmail)
	}

	if password == "" {
		// return user, errNoPassword
		errs = append(errs, errNoPassword)
	}else if len(password) < passwordLength {
		// return user, errPasswordTooShort
		errs = append(errs, errPasswordTooShort)
	}

	// Check if the username exits
	existingUser, err := globalUserStore.FindByUsername(username)
	if err != nil {
		// return user, err
		errs = append(errs, err)
	}
	if existingUser != nil {
		// return user, errUsernameExist
		errs = append(errs, errUsernameExist)
	}

	// Check if the email exists
	existingUser, err = globalUserStore.FindByEmail(email)
	if err != nil {
		// return user, err
		errs = append(errs, err)
	}
	if existingUser != nil {
		// return user, errEmailExists
		errs = append(errs, errEmailExists)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	user.HashPassword = string(hashPassword)
	user.ID = GenerateID("usr", userIDLength)

	return user, errs
}
