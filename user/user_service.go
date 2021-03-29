package user

import (
	"context"
	"errors"

	"github.com/go-ieat-api/prisma"
	"github.com/go-ieat-api/prisma/db"
)

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

var ctx = context.Background()
var ErrEmailExists = errors.New("Email already exists")

func createUser(user *User, hashed string) error {

	if exists := isEmailExist(user.Email); exists {
		return ErrEmailExists
	}

	created, err := prisma.Db.User.CreateOne(
		db.User.Email.Set(user.Email),
		db.User.Name.Set(user.Name),
		db.User.Hashedpassword.Set(hashed),
	).Exec(ctx)

	if err != nil {
		return err
	}

	user.Id = created.ID

	return nil
}

func isEmailExist(email string) bool {
	_, err := prisma.Db.User.FindUnique(
		db.User.Email.Equals(email),
	).Exec(ctx)

	if errors.Is(err, db.ErrNotFound) {
		return false
	}

	return true
}
