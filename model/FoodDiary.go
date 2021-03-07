package model

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// FoodDiary is diary?
type FoodDiary struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	FoodName     string    `json:"foodName"`
	FoodCalories float32   `json:"calories"`
	Date         time.Time `json:"date"`
	Note         *string   `json:"note,omitempty"`
}

var db *gorm.DB

func init() {
	d, err := gorm.Open(sqlite.Open("Ieat.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db = d

	db.AutoMigrate(&FoodDiary{})
}

// GetFoodDiaryByID is ...
func GetFoodDiaryByID(id int, fd *FoodDiary) {
	db.First(fd, id)
}

// CreateFoodDiary is ...
func CreateFoodDiary(fd *FoodDiary) {
	db.Create(fd)
}

// GetAllFoodDiaries is ...
func GetAllFoodDiaries() []FoodDiary {
	fd := []FoodDiary{}
	db.Find(&fd)
	return fd
}
