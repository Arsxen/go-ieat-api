package login

import (
	"context"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-ieat-api/prisma"
	"github.com/go-ieat-api/prisma/db"
)

var ctx = context.Background()

func getUserHashedPassword(email string) (*db.UserModel, error) {
	user, err := prisma.Db.User.FindUnique(
		db.User.Email.Equals(email),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func createToken(userId int) (string, error) {
	cliams := jwt.MapClaims{}
	cliams["user_id"] = userId
	cliams["exp"] = time.Now().Add(time.Minute * 20).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cliams)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
