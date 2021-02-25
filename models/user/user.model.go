package user

import (
	"admin/database"
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User contains all nessecary information about a user
type User struct {
	ID       uint      `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Created  time.Time `json:"created"`
	Admin    bool      `json:"admin"`
}

// HashPassword generates a hash from a string
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// ComparePassword compares a hash and a string and checks the validity
func ComparePassword(savedPassword string, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(savedPassword), []byte(password)); err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

// GetUserByUsername select a user from database by its username
func (u *User) GetUserByUsername() (User, error) {
	var user User
	// "root:password@/admin_db?parseTime=true"
	db := database.Connect(&database.SQLConfig{
		User:     "root",
		Password: "password",
		Database: "admin_db",
	})
	defer db.Close()

	row := db.QueryRow("SELECT * FROM users WHERE username=?", u.Username)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Created,
		&user.Admin)

	switch {
	case err == sql.ErrNoRows:
		return User{}, err
	case err != nil:
		return User{}, err
	default:
		return user, nil
	}
}

// CreateNewUser inserts a new user row into the database
func (u *User) CreateNewUser() error {
	db := database.Connect(&database.SQLConfig{
		User:     "root",
		Password: "password",
		Database: "admin_db",
	})
	defer db.Close()

	hash, err := HashPassword(u.Password)

	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		u.Username, u.Email, hash)

	if err != nil {
		return err
	}

	return nil
}
