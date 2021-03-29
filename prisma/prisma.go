package prisma

import (
	"github.com/go-ieat-api/prisma/db"
	"github.com/joho/godotenv"
)

var Db *db.PrismaClient

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	Db = db.NewClient()
}

func DisconnectDb() {
	if err := Db.Prisma.Disconnect(); err != nil {
		panic(err)
	}
}
