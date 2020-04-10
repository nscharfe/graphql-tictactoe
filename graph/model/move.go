package model

import "time"

type Move struct {
	ID          string `json:"id"`
	GameID      string `json:"game"`
	RowIndex    int    `json:"rowIndex"`
	ColumnIndex int    `json:"columnIndex"`
	Player      Player `json:"player"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
