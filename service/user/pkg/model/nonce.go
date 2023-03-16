package model

import "time"

type Nonce struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"create_time"`
	UserID    string    `json:"user_id"`
	Value     string    `json:"value"`
}
