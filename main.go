package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-ieat-api/model"
	"github.com/go-ieat-api/prisma/db"
)

// FoodDiaryRequest is ...
type FoodDiaryRequest struct {
	FoodName     string    `json:"foodName"`
	FoodCalories float64   `json:"calories"`
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
	r.Use(recoverWithJSON)
	r.Use(middleware.SetHeader("Content-Type", "application/json; charset=utf-8"))

	r.Route("/food", func(r chi.Router) {
		r.Post("/", createFoodDiary)
		r.Get("/", getAllFoodDiary)
		r.With(foodDiaryCtx).Get("/{foodDiaryID}", getFoodDiary)
	})

	server := &http.Server{
		Handler: r,
		Addr:    ":8000",
	}

	go func() {
		fmt.Println(server.ListenAndServe())
	}()

	//Gracefully Shutdown
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}

	model.DisconnectPrisma()
}

func recoverWithJSON(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {

				logEntry := middleware.GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else {
					middleware.PrintPrettyStack(rvr)
				}

				errorJSON(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func createFoodDiary(w http.ResponseWriter, r *http.Request) {
	data := &FoodDiaryRequest{}

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		errorJSON(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	fd := model.FoodDiary{FoodName: data.FoodName, Calories: data.FoodCalories, Date: data.Date, Note: data.Note}

	if err := model.CreateFoodDiary(&fd); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(fd)
}

func getAllFoodDiary(w http.ResponseWriter, _ *http.Request) {
	fd, err := model.GetAllFoodDiaries()
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(fd)
}

func foodDiaryCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		diaryID := chi.URLParam(r, "foodDiaryID")

		if diaryID == "" {
			errorJSON(rw, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		id, err := strconv.Atoi(diaryID)

		if err != nil {
			errorJSON(rw, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		var fd *model.FoodDiary

		fd, err = model.GetFoodDiaryByID(id)

		if err == db.ErrNotFound {
			errorJSON(rw, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		ctx := context.WithValue(r.Context(), foodDiaryKey, fd)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

func getFoodDiary(w http.ResponseWriter, r *http.Request) {
	fd := r.Context().Value(foodDiaryKey).(*model.FoodDiary)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(fd)
}

type jsonError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func errorJSON(w http.ResponseWriter, code int, msg string) {
	err := jsonError{Code: code, Message: msg}
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(err)
}
