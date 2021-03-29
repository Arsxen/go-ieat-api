package user

import (
	"context"

	"github.com/go-ieat-api/prisma"
	"github.com/go-ieat-api/prisma/db"
)

type User struct {
	Id    int    `json:id`
	Email string `json:email`
	Name  string `json:name`
}

var ctx = context.Background()

func createUser(user *User, hashed string) error {
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
