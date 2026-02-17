package one

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	Storage "github.com/okishiro/pidgey/internal/storage"
	types "github.com/okishiro/pidgey/internal/types"
	"github.com/okishiro/pidgey/internal/utils/response"
)

func CreateMovie(datab Storage.Store) http.HandlerFunc { //pass the interface
	return func(w http.ResponseWriter, r *http.Request) {

		var recieved types.Movie
		accountname := r.PathValue("name")
		slog.Info("creating ")
		if accountname == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&recieved)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(recieved); err != nil {
			validateError := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationErrors(validateError))
			return
		}

		today := time.Now().UTC().Truncate(24 * time.Hour)
		fmt.Printf("%s helllo", recieved.Name)
		lastid, err := datab.CreateMovieEntry(
			accountname,
			recieved.Name,
			today,
			recieved.Comment,
		)
		if err != nil {
			slog.Error("Database Error", "details", err.Error()) // This will print the actual "no such table" message
			response.WriteJson(w, http.StatusInternalServerError, err.Error())
			return
		}

		// slog.Info(message, key, value)
		fmt.Printf("%s cmmmon", recieved.Comment)
		slog.Info("creating an entry", "id", lastid)
		response.WriteJson(w, http.StatusCreated, "")
	}
}

func GetData(datab Storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idstr := r.PathValue("id")
		if idstr == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID format: must be a number", http.StatusBadRequest)
			return
		}
		data, err := datab.Getmovies(id)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Send JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)

	}
}

func CreateAccount(datab Storage.Store, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}
		/*
			err := json.NewDecoder(r.Body).Decode(&recieved)
			if errors.Is(err, io.EOF) {
				response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
				return
			}

			if err := validator.New().Struct(recieved); err != nil {
				validateError := err.(validator.ValidationErrors)
				response.WriteJson(w, http.StatusBadRequest, response.ValidationErrors(validateError))
				return
			}
		*/

		lastid, err := datab.CreateAccount(
			name,
		)
		if err != nil {
			log.Printf("sqlite error: %v", err)
			response.WriteJson(w, http.StatusInternalServerError, err)
			return

		}
		err = datab.CreateTable(
			lastid,
			path,
		)
		if err != nil {
			log.Printf("sqlite error: %v", err)
			response.WriteJson(w, http.StatusInternalServerError, err)
			return

		}
		slog.Info("creating an account")
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastid})
	}
}
