package model

import "time"

type Status int8

const (
	StatusUnspecified Status = iota
	StatusOpen
	StatusComplete
	StatusCanceled
)

type CharStatus int8

const (
	CharStatusUnspecified CharStatus = iota
	CharStatusNotPresent
	CharStatusPresent
	CharStatusCorrect
)

type Wordle struct {
	ID         string         `json:"id"`
	CreatedAt  time.Time      `json:"create_time"`
	UpdatedAt  time.Time      `json:"updated_time"`
	UserID     string         `json:"user_id"`
	Status     Status         `json:"status"`
	Solution   string         `json:"solution"`
	Guesses    []string       `json:"guesses"`
	CharStatus [][]CharStatus `json:"char_status"`
}
