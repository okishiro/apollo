package types

import "time"

type Movie struct {
	// We don't need a JSON tag for Accountname if it comes from the URL path,
	// but we need them for the fields coming in the Request Body.
	Accountname string    `json:"accountname"`
	Name        string    `json:"name" validate:"required"`
	Date        time.Time `json:"date"`
	Comment     string    `json:"comment" validate:"required"`
}

type Ids struct {
	Id   int64
	Name string `validate:"required"`
}

type GetDataResponse struct {
	ID        int64     `json:"id"`
	MovieName string    `json:"name"`
	Timestamp time.Time `json:"time"`
	Comment   string    `json:"comment"`
}
