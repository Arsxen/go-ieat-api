package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/badoux/checkmail"
	"github.com/go-chi/chi/v5"
	"github.com/go-ieat-api/render"
	"golang.org/x/crypto/bcrypt"
)

type UserRegistrationRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func Router() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/regis", registerUser)
	return r
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	userReq := UserRegistrationRequest{}

	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		render.ErrorJSON(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}

	if err := checkmail.ValidateFormat(userReq.Email); err != nil {
		render.ErrorJSON(w, http.StatusBadRequest, "Invalid email address")
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	user := User{Email: userReq.Email, Name: userReq.Name}

	err = createUser(&user, string(hashed))

	if err != nil {
		if errors.Is(err, ErrEmailExists) {
			render.ErrorJSON(w, http.StatusBadRequest, "Email already exists")
			return
		}

		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}
