package database

import (
	"fmt"

	"github.com/google/uuid"
)

var db *Database

type Database struct {
	Users map[string]User
}

type User struct {
	ID       string
	Email    string
	Password string
}

func NewDatabase() *Database {
	db = &Database{
		Users: map[string]User{},
	}
	return db
}

func (db *Database) NewUser(u User) error {
	userid := uuid.New().String()
	db.Users[userid] = User{
		ID:       userid,
		Email:    u.Email,
		Password: u.Password,
	}

	return nil
}

func (db *Database) GetUserByEmail(email string) (User, error) {
	for _, u := range db.Users {
		if u.Email == email {
			return u, nil
		}
	}
	return User{}, fmt.Errorf("User not found")
}

func (db *Database) GetUserById(userID string) (User, error) {
	u, ok := db.Users[userID]
	if !ok {
		return User{}, fmt.Errorf("User not found")
	}
	return u, nil
}
