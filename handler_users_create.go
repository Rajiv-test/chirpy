package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Rajiv-test/chirpy/internal/auth"
	"github.com/Rajiv-test/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	IsChirpyRed bool    `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	if params.Password == "" {
		// Handle missing password case
		// Perhaps return a 400 Bad Request with a helpful message
		respondWithError(w,http.StatusBadRequest,"Please provide password",err)
		return
	}
	hashed_password,err  := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w,http.StatusBadRequest,"Error hashing password",err)
		return 
	}

	user, err := cfg.db.CreateUser(r.Context(),database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashed_password,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			IsChirpyRed: user.IsChirpyRed.Bool,
		},
	})
}
