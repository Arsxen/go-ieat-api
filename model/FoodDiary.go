package model

import (
	"context"
	"fmt"
	"time"

	"github.com/go-ieat-api/prisma/db"
	"github.com/joho/godotenv"
)

// FoodDiary is ...
type FoodDiary struct {
	ID       uint      `json:"id"`
	FoodName string    `json:"foodName"`
	Calories float64   `json:"calories"`
	Date     time.Time `json:"date"`
	Note     *string   `json:"note,omitempty"`
}

var client *db.PrismaClient
var ctx context.Context

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	client = db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}

	ctx = context.Background()
}

// DisconnectPrisma - disconnect the prisma client
func DisconnectPrisma() {
	if err := client.Prisma.Disconnect(); err != nil {
		panic(err)
	}
	fmt.Println("Prisma client disconnected.")
}

// GetFoodDiaryByID is ...
func GetFoodDiaryByID(id int) (*FoodDiary, error) {
	fd, err := client.FoodDiary.FindUnique(
		db.FoodDiary.ID.Equals(id),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	food := &FoodDiary{}

	food.fromFoodDiaryModel(fd)

	return food, nil
}

// CreateFoodDiary is ...
func CreateFoodDiary(fd *FoodDiary) error {
	created, err := client.FoodDiary.CreateOne(
		db.FoodDiary.FoodName.Set(fd.FoodName),
		db.FoodDiary.Calories.Set(fd.Calories),
		db.FoodDiary.Date.Set(fd.Date),
		db.FoodDiary.Note.SetIfPresent(fd.Note),
	).Exec(ctx)

	if err != nil {
		return err
	}

	fd.ID = uint(created.ID)

	return nil
}

// GetAllFoodDiaries is ...
func GetAllFoodDiaries() ([]FoodDiary, error) {
	foodDiaries, err := client.FoodDiary.FindMany().Exec(ctx)

	if err != nil {
		return nil, err
	}

	return modelsToFoodDiaries(foodDiaries), nil
}

func modelsToFoodDiaries(models []db.FoodDiaryModel) []FoodDiary {
	diaries := make([]FoodDiary, len(models))

	for i := range diaries {
		f := FoodDiary{}
		f.fromFoodDiaryModel(&models[i])
		diaries[i] = f
	}

	return diaries
}

func (fd *FoodDiary) fromFoodDiaryModel(model *db.FoodDiaryModel) {
	fd.ID = uint(model.ID)
	fd.FoodName = model.FoodName
	fd.Calories = model.Calories
	fd.Date = model.Date

	if n, ok := model.Note(); ok {
		fd.Note = &n
	} else {
		fd.Note = nil
	}
}
