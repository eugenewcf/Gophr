package main

import (
	"crypto/md5"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	// "io"
)

const (
	passwordLength = 8
	hashCost       = 10
	userIDLength
)

type User struct {
	Username       string
	Email          string
	HashedPassword string
	ID             string
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
	} else if len(password) < passwordLength {
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
	user.HashedPassword = string(hashPassword)
	user.ID = GenerateID("usr", userIDLength)

	return user, errs
}

func FindUser(username, password string) (*User, error) {
	out := &User{
		Username: username,
	}

	existingUser, err := globalUserStore.FindByUsername(username)
	if err != nil {
		return out, err
	}
	if existingUser == nil {
		return out, errCredentialsIncorrect
	}

	if bcrypt.CompareHashAndPassword(
		[]byte(existingUser.HashedPassword),
		[]byte(password),
	) != nil {
		return out, errCredentialsIncorrect
	}
	return existingUser, nil
}

func UpdateUser(user *User, email, currentPassword, newPassword string) (User, error) {
	out := *user
	out.Email = email

	// Check if the email exists
	existingUser, err := globalUserStore.FindByEmail(email)
	if err != nil {
		return out, err
	}
	if existingUser != nil && existingUser.ID != user.ID {
		return out, errEmailExists
	}

	// At this point, we can update the email address
	user.Email = email

	// No current password? Don't try update the password
	if currentPassword == "" {
		return out, nil
	}

	if bcrypt.CompareHashAndPassword(
		[]byte(user.HashedPassword),
		[]byte(currentPassword),
	) != nil {
		return out, errPasswordIncorrect
	}

	if newPassword == "" {
		return out, errNoPassword
	}

	if len(newPassword) < passwordLength {
		return out, errPasswordTooShort
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), hashCost)
	user.HashedPassword = string(hashedPassword)
	return out, err
}

func (user *User) AvatarURL() string {
	return fmt.Sprintf(
		"//www.gravatar.com/avatar/%x",
		md5.Sum([]byte(user.Email)),
	)
}

func (user *User) ImageRoute() string {
	return "/user/" + user.ID
}
