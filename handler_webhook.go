package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"github.com/Rajiv-test/chirpy/internal/auth"
	"github.com/google/uuid"

)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request){
	type request struct{
		Event string `json:"event"`
		Data struct{
			UserId uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	apikey,err := auth.GetAPIKey(r.Header)
	if err != nil || cfg.polkaKey != apikey {
		respondWithError(w,http.StatusUnauthorized,"couldn't get apikey",err)
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := request{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.db.UpgradeUser(r.Context(), params.Data.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}