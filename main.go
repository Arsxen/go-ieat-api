package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-ieat-api/model"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/food", func(r chi.Router) {
		r.Post("/", createFoodDiary)
		r.Get("/", getAllFoodDiary)
	})

	log.Fatal(http.ListenAndServe(":8000", r))
}

func createFoodDiary(w http.ResponseWriter, r *http.Request) {
	data := &model.FoodDiary{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	model.CreateFoodDiary(data)

	json.NewEncoder(w).Encode(data)
}

func getAllFoodDiary(w http.ResponseWriter, r *http.Request) {
	fd := model.GetAllFoodDiaries()

	json.NewEncoder(w).Encode(fd)
}
