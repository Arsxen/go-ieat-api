package user

import (
	"encoding/json"
	"net/http"

	"github.com/badoux/checkmail"
	"github.com/go-chi/chi/v5"
	"github.com/go-ieat-api/render"
	"golang.org/x/crypto/bcrypt"
)

type UserRegistrationRequest struct {
	Id       int
	Email    string
	Name     string
	Password string
}

func Router() {
	r := chi.NewRouter()
	r.Post("/regis", registerUser)
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	userReq := UserRegistrationRequest{}

	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		render.ErrorJSON(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}

	if err := checkmail.ValidateFormat(userReq.Email); err != nil {
		render.ErrorJSON(w, http.StatusUnprocessableEntity, "Invalid email address")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	user := User{Email: userReq.Email, Name: userReq.Name}

	if err := createUser(&user, string(hashed)); err != nil {
		panic(err)
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}
