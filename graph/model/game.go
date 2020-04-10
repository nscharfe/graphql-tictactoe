package model

import "time"

type Game struct {
	ID     string  `json:"id"`
	Winner *Player `json:"winner"`
	Turn   Player  `json:"turn"`
	Status Status  `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
