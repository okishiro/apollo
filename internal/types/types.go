package types

import "time"

type Movie struct {
	Accountname string
	Name        string `validate:"required"`
	Date        time.Time
	Comment     string
}

type Ids struct {
	Id   int64
	Name string `validate:"required"`
}
