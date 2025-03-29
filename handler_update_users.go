package main

import (
	"encoding/json"
	"net/http"

	"github.com/Rajiv-test/chirpy/internal/auth"
	"github.com/Rajiv-test/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUsers(w http.ResponseWriter,r *http.Request){

	type request struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}
	req_token ,err := auth.GetBearerToken(r.Header)
	
	if err != nil {
		respondWithError(w,http.StatusUnauthorized,"couldn't find token",err)
		return
	}
	var req request
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	
	if err != nil{
		respondWithError(w,http.StatusInternalServerError,"error decoding request body",err)
		return
	}
	
	if req.Email == "" || req.Password == "" {
		respondWithError(w,http.StatusBadRequest,"pr0vide valid email and password",err)
		return
	}
	userId,err := auth.ValidateJWT(req_token,cfg.jwtSecret)
	
	if err != nil {
		respondWithError(w,http.StatusUnauthorized,"couldn't validate access token",err)
		return
	}
	hashedPassword,err := auth.HashPassword(req.Password)
	
	if err != nil {
		respondWithError(w,http.StatusInternalServerError,"couldn't hash the password",err)
		return
	}
	user,err := cfg.db.UpdateUser(r.Context(),database.UpdateUserParams{
		HashedPassword: hashedPassword,
		Email: req.Email,
		ID: userId,
	})
	if err != nil {
		respondWithError(w,http.StatusInternalServerError,"couldn't update user password",err)
		return
	}
	respondWithJSON(w,http.StatusOK,User{
		ID:user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})

}