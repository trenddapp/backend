package model

import "time"

type User struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"create_time"`
	UpdatedAt time.Time `json:"update_time"`
	Address   string    `json:"address"`
	Balance   int64     `json:"balance"`
}
