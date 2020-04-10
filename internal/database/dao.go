package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/nscharfe/graphql-tictactoe/graph/model"
)

type Dao interface {
	CreateGame(*sql.Tx, *model.Game) (*model.Game, error)
	GetGame(*sql.Tx, string) (*model.Game, error)
	UpdateGame(*sql.Tx, *model.Game) (*model.Game, error)

	CreateMove(*sql.Tx, *model.Move) (*model.Move, error)
	GetMove(*sql.Tx, string) (*model.Move, error)
	GetMovesForGame(*sql.Tx, string) ([]*model.Move, error)
}

type dao struct{}

func NewDao() Dao {
	return dao{}
}

func (c dao) CreateGame(tx *sql.Tx, game *model.Game) (*model.Game, error) {
	insert := fmt.Sprintf(`
       INSERT INTO %s (id, winner, turn, status, created_at, updated_at)
       VALUES ($1, $2, $3, $4, $5, $6)
    `, tables.Games)

	if _, err := tx.Exec(insert, game.ID, game.Winner, game.Turn, game.Status, time.Now(), time.Now()); err != nil {
		return nil, err
	}

	return c.GetGame(tx, game.ID)
}

func (c dao) GetGame(tx *sql.Tx, id string) (*model.Game, error) {
	var result model.Game
	query := fmt.Sprintf(`
        SELECT id, winner, turn, status, created_at, updated_at 
        FROM %s 
        WHERE id = $1`,
		tables.Games)
	if err := tx.QueryRow(query, id).Scan(&result.ID, &result.Winner, &result.Turn, &result.Status, &result.CreatedAt, &result.UpdatedAt); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c dao) UpdateGame(tx *sql.Tx, game *model.Game) (*model.Game, error) {
	update := fmt.Sprintf(`
        UPDATE %s
        SET winner = $1,
            turn = $2,
            status = $3,
            updated_at = $4
        WHERE id = $5
    `, tables.Games)

	if _, err := tx.Exec(update, game.Winner, game.Turn, game.Status, time.Now(), game.ID); err != nil {
		return nil, err
	}
	return c.GetGame(tx, game.ID)
}

func (c dao) CreateMove(tx *sql.Tx, move *model.Move) (*model.Move, error) {
	insert := fmt.Sprintf(`
       INSERT INTO %s (id, game_id, row_index, column_index, player, created_at, updated_at)
       VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, tables.Moves)

	if _, err := tx.Exec(insert, move.ID, move.GameID, move.RowIndex, move.ColumnIndex, move.Player, time.Now(), time.Now()); err != nil {
		return nil, err
	}

	return c.GetMove(tx, move.ID)
}

func (c dao) GetMove(tx *sql.Tx, id string) (*model.Move, error) {
	var result model.Move
	query := fmt.Sprintf(`
        SELECT id, game_id, row_index, column_index, player, created_at, updated_at 
        FROM %s 
        WHERE id = $1`,
		tables.Moves)
	if err := tx.QueryRow(query, id).Scan(&result.ID, &result.GameID, &result.RowIndex, &result.ColumnIndex, &result.Player, &result.CreatedAt, &result.UpdatedAt); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c dao) GetMovesForGame(tx *sql.Tx, gameId string) ([]*model.Move, error) {
	var result []*model.Move
	query := fmt.Sprintf(`
        SELECT id, game_id, row_index, column_index, player, created_at, updated_at 
        FROM %s 
        WHERE game_id = $1`,
		tables.Moves)
	rows, err := tx.Query(query, gameId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		move := model.Move{}
		if err = rows.Scan(&move.ID, &move.GameID, &move.RowIndex, &move.ColumnIndex, &move.Player, &move.CreatedAt, &move.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, &move)
	}
	return result, nil
}
