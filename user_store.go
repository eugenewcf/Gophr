package main

import (
	"encoding/json"
	// "io/ioutil"
	"os"
	"fmt"
	"strings"
)

type UserStore interface {
	Find(string) (*User, error)
	FindByEmail(string) (*User, error)
	FindByUsername(string) (*User, error)
	Save(User) error
}

type FileUserStore struct {
	FileStore
	Users    map[string]User
}

var globalUserStore UserStore
var globalUserEmailMapping map[string] *User
var globalUserUsernameMapping map[string] *User

func init() {
	store, err := NewFileUserStore("./data/users.json")
	if err != nil {
		panic(fmt.Errorf("Error creating user store: %s", err))
	}
	globalUserStore = store
	globalUserEmailMapping = map[string] *User{}
	globalUserUsernameMapping = map[string] *User{}
	for _, user := range store.Users {
		globalUserEmailMapping[strings.ToLower(user.Email)] = &user
		globalUserUsernameMapping[strings.ToLower(user.Username)] = &user
	}
}

func NewFileUserStore(filename string) (*FileUserStore, error) {
	store := &FileUserStore{
		Users: map[string]User{},
		FileStore: FileStore{
			filename: filename,
		},
	}

	// contents, err := ioutil.ReadFile(filename)
	contents, err := store.FileStore.Read()

	if err != nil {
		// If it's a matter of the file not existing, that's ok
		if os.IsNotExist(err) {
			return store, nil
		}
		return nil, err
	}
	err = json.Unmarshal(contents, store)
	if err != nil {
		return nil, err
	}
	return store, nil
}

func (store FileUserStore) Save(user User) error {
	store.Users[user.ID] = user
	globalUserEmailMapping[strings.ToLower(user.Email)] = &user
	globalUserUsernameMapping[strings.ToLower(user.Username)] = &user

	contents, err := json.MarshalIndent(store, "", "	")
	if err != nil {
		return err
	}

	// err = ioutil.WriteFile(store.filename, contents, 0660)
	err = store.FileStore.Write(contents)
	if err != nil {
		return err
	}

	return nil
}

func (store FileUserStore) Find(id string) (*User, error) {
	user, ok := store.Users[id]
	if ok {
		return &user, nil
	}
	return nil, nil
}

func (store FileUserStore) FindByUsername(username string) (*User, error) {
	if username == "" {
		return nil, nil
	}

	user, ok := globalUserUsernameMapping[strings.ToLower(username)]
	if ok {
		return user, nil
	}

	return nil, nil
}

func (store FileUserStore) FindByEmail(email string) (*User, error) {
	if email == "" {
		return nil, nil
	}

	user, ok := globalUserEmailMapping[strings.ToLower(email)]
	if ok {
		return user, nil
	}

	return nil, nil
}
