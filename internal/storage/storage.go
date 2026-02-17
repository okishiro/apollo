package Storage

import (
	"time"

	"github.com/okishiro/pidgey/internal/types"
)

type Store interface {
	CreateMovieEntry(accountname string, name string, time time.Time, comment string) (int64, error)
	CreateAccount(name string) (int64, error)
	CreateTable(id int64, path string) error
	Getmovies(id int64) ([]types.GetDataResponse, error)
}
