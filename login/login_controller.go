package login

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-ieat-api/prisma/db"
	"github.com/go-ieat-api/render"
	"golang.org/x/crypto/bcrypt"
)

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Router() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", authenticateUser)
	return r
}

func authenticateUser(w http.ResponseWriter, r *http.Request) {
	loginReq := UserLoginRequest{}

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		render.ErrorJSON(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	user, err := getUserHashedPassword(loginReq.Email)

	if errors.Is(err, db.ErrNotFound) {
		render.ErrorJSON(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	} else if err != nil {
		panic(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Hashedpassword), []byte(loginReq.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			render.ErrorJSON(w, http.StatusUnauthorized, "Incorrect email or password")
			return
		}
		panic(err)
	}

	token, err := createToken(user.ID)
	if err != nil {
		panic(err)
	}

	res := map[string]interface{}{
		"token": token,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
}
