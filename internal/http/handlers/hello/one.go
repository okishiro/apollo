package one

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"
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

		lastid, err := datab.CreateMovie(
			accountname,
			recieved.Name,
			today,
			recieved.Comment,
		)

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		// slog.Info(message, key, value)
		slog.Info("creating an account", "id", lastid)
		response.WriteJson(w, http.StatusCreated, "")
	}
}

/*
func GGetmovie(datab Storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		//slog.Info("getting info of movie with id %v", id)

		intid, err := strconv.ParseInt(id, 10, 64)
		moviee, err := datab.GetMovieByID(intid)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		}

		response.WriteJson(w, http.StatusOK, moviee)
	}
}
*/

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
