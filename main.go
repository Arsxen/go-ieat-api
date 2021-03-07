package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-ieat-api/model"
)

// FoodDiaryRequest is ...
type FoodDiaryRequest struct {
	FoodName     string    `json:"foodName"`
	FoodCalories float32   `json:"calories"`
	Date         time.Time `json:"date"`
	Note         *string   `json:"note"`
}

type ctxKey int

const (
	foodDiaryKey ctxKey = iota
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json; charset=utf-8"))

	r.Route("/food", func(r chi.Router) {
		r.Post("/", createFoodDiary)
		r.Get("/", getAllFoodDiary)
		r.With(foodDiaryCtx).Get("/{foodDiaryID}", getFoodDiary)
	})

	log.Fatal(http.ListenAndServe(":8000", r))
}

func createFoodDiary(w http.ResponseWriter, r *http.Request) {
	data := &FoodDiaryRequest{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	fd := model.FoodDiary{FoodName: data.FoodName, FoodCalories: data.FoodCalories, Date: data.Date, Note: data.Note}
	model.CreateFoodDiary(&fd)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fd)
}

func getAllFoodDiary(w http.ResponseWriter, r *http.Request) {
	fd := model.GetAllFoodDiaries()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fd)
}

func foodDiaryCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fd := &model.FoodDiary{}

		diaryID := chi.URLParam(r, "foodDiaryID")

		if diaryID == "" {
			errorJSON(rw, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		}

		id, err := strconv.Atoi(diaryID)

		if err != nil {
			panic(err)
		}

		model.GetFoodDiaryByID(id, fd)

		ctx := context.WithValue(r.Context(), foodDiaryKey, fd)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

func getFoodDiary(w http.ResponseWriter, r *http.Request) {
	fd := r.Context().Value(foodDiaryKey).(*model.FoodDiary)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fd)
}

type jsonError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func errorJSON(w http.ResponseWriter, code int, msg string) {
	err := jsonError{Code: code, Message: msg}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}
