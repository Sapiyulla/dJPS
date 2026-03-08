package domain

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// @Entity
type User struct {
	id       uuid.UUID
	name     string
	email    string
	password []byte
}

func NewUser(name, email, password string) *User {
	passwordHash, err := hashPassword(password)
	if err != nil {
		passwordHash = []byte(password)
	}
	return &User{
		id:       uuid.New(),
		name:     name,
		email:    email,
		password: passwordHash,
	}
}

func RehydrateUser(id uuid.UUID, name, email, password string) *User {
	return &User{
		id:       id,
		name:     name,
		email:    email,
		password: []byte(password),
	}
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return u.email
}

// return password Hash
func (u *User) Password() []byte {
	return u.password
}

func hashPassword(password string) ([]byte, error) {
	pwhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return pwhash, nil
}

func CheckPasswordHash(password string, hashedPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	return err == nil
}
