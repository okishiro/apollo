package one

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	Storage "github.com/okishiro/pidgey/internal/storage"
	types "github.com/okishiro/pidgey/internal/types"
	"github.com/okishiro/pidgey/internal/ui"
	"github.com/okishiro/pidgey/internal/utils/response"
)

func CreateMovie(datab Storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var movieTitle, comment string

		// 1. Get accountname from the URL {name}
		accountname := r.PathValue("name")

		// 2. Decide if we are reading JSON or FORM data
		if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
			// Data from your GOTH frontend
			movieTitle = r.FormValue("name")
			comment = r.FormValue("comment")
		} else {
			// Data from your JSON stress tests
			var recieved types.Movie
			if err := json.NewDecoder(r.Body).Decode(&recieved); err == nil {
				movieTitle = recieved.Name
				comment = recieved.Comment
			}
		}

		// 3. Generate the 3rd argument (time)
		today := time.Now().UTC().Truncate(24 * time.Hour)

		// 4. Pass the 4 arguments to your existing DB logic
		lastid, err := datab.CreateMovieEntry(
			accountname, // Arg 1
			movieTitle,  // Arg 2
			today,       // Arg 3
			comment,     // Arg 4
		)

		if err != nil {
			slog.Error("Database Error", "details", err.Error())
			http.Error(w, "DB Fail", 500)
			return
		}

		// 5. Response logic
		if r.Header.Get("HX-Request") == "true" {
			// For the frontend: Send back a single row
			ui.MovieRow(movieTitle, accountname).Render(r.Context(), w)
		} else {
			// For the API: Send back JSON
			response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastid})
		}
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
