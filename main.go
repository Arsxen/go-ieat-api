package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-ieat-api/login"
	"github.com/go-ieat-api/prisma"
	"github.com/go-ieat-api/render"
	"github.com/go-ieat-api/user"
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

	r.Mount("/user", user.Router())
	r.Mount("/login", login.Router())

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

	prisma.DisconnectDb()
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

				render.ErrorJSON(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
