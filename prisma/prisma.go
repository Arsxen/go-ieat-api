package prisma

import (
	"fmt"

	"github.com/go-ieat-api/prisma/db"
	"github.com/joho/godotenv"
)

var Db *db.PrismaClient

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	Db = db.NewClient()
	if err := Db.Prisma.Connect(); err != nil {
		panic(err)
	}
}

func DisconnectDb() {
	if err := Db.Prisma.Disconnect(); err != nil {
		panic(err)
	}
	fmt.Println("Disconnect prisma client")
}
