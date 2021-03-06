package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// FoodDiary is diary?
type FoodDiary struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	FoodName string `json:"foodName"`
	// FoodCalories float32   `json:"calories"`
	// Date         time.Time `json:"date"`
	// Note         string    `json:"note"`
}

var db *gorm.DB

func init() {
	d, err := gorm.Open(sqlite.Open("IEat.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db = d

	db.AutoMigrate(&FoodDiary{})
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
