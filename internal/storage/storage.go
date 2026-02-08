package Storage

import "time"

type Store interface {
	CreateMovie(accountname string, name string, time time.Time, comment string) (int64, error)
	CreateAccount(name string) (int64, error)
	CreateTable(id int64, path string) error
}
