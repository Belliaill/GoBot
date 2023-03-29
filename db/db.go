package db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/exp/slices"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	IsBanned bool   `json:"banned"`
}

type DB struct {
	path  string
	users []User
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	
	return false, err
}

func NewDB(path string) *DB {
	db := DB{path: path, users: make([]User, 0)}
	ok, err := Exists(path)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	if ok {
		db.Pull()
	} else {
		db.Push()
	}
	return &db
}

func (db *DB) Pull() {
	data, err := os.ReadFile(db.path)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	json.Unmarshal(data, &db.users)
}

func (db *DB) Push() {
	data, err := json.Marshal(db.users)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	f, err := os.Create(db.path)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Fprint(f, data)
}

func (db *DB) AppendUser(user User) {
	db.Pull()
	db.users = append(db.users, user)
	db.Push()
}

func (db *DB) RemoveUser(index int) {
	db.Pull()
	slices.Delete(db.users, index, index+1)
	db.Push()
}

func (db *DB) GetUsers() []User {
	return db.users
}
