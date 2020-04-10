package backend

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/nscharfe/graphql-tictactoe/internal/database"

	"github.com/nscharfe/graphql-tictactoe/internal/tictactoe"

	"github.com/nscharfe/graphql-tictactoe/graph/model"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type Client interface {
	CreateGame(context.Context) (*model.Game, error)
	GetGame(context.Context, string) (*model.Game, error)

	CreateMove(context.Context, string, int, int, model.Player) (*model.Move, error)
	GetMovesForGame(context.Context, string) ([]*model.Move, error)
}

type client struct {
	db  database.DB
	dao database.Dao
}

func New(db database.DB, dbClient database.Dao) Client {
	return &client{
		db:  db,
		dao: dbClient,
	}
}

func (c *client) CreateGame(ctx context.Context) (*model.Game, error) {
	gameId := fmt.Sprintf("T%d", rand.Int())
	game := &model.Game{
		ID:     gameId,
		Winner: nil,
		Turn:   model.PlayerX,
		Status: model.StatusActive,
	}
	var err error
	if err := c.db.RunTransaction(ctx, nil, func(tx *sql.Tx) error {
		game, err = c.dao.CreateGame(tx, game)
		return err
	}); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return game, nil
}

func (c *client) GetGame(ctx context.Context, id string) (*model.Game, error) {
	var game *model.Game
	var err error

	if err := c.db.RunTransaction(ctx, nil, func(tx *sql.Tx) error {
		game, err = c.dao.GetGame(tx, id)
		return err
	}); err != nil {
		return nil, err
	}

	return game, nil
}

func (c *client) CreateMove(ctx context.Context, gameID string, rowIndex int, columnIndex int, player model.Player) (*model.Move, error) {
	move := &model.Move{
		GameID:      gameID,
		RowIndex:    rowIndex,
		ColumnIndex: columnIndex,
		Player:      player,
		ID:          fmt.Sprintf("T%d", rand.Int()),
	}

	// save move and updated status
	if err := c.db.RunTransaction(ctx, nil, func(tx *sql.Tx) error {
		var movesForGame []*model.Move
		var err error

		// validate and make move and get new game status
		if movesForGame, err = c.dao.GetMovesForGame(tx, gameID); err != nil {
			return err
		}
		gameClient := tictactoe.New(movesForGame)
		if err := gameClient.MakeMove(move); err != nil {
			return err
		}
		// create new move
		if move, err = c.dao.CreateMove(tx, move); err != nil {
			return err
		}

		// fetch game and update to reflect new game status
		var game *model.Game
		if game, err = c.dao.GetGame(tx, gameID); err != nil {
			return err
		}
		game.Winner = gameClient.Winner()
		game.Status = gameClient.Status()
		game.Turn = gameClient.CurrentTurn()
		if _, err := c.dao.UpdateGame(tx, game); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return move, nil
}
func (c *client) GetMovesForGame(ctx context.Context, gameId string) ([]*model.Move, error) {
	var moves []*model.Move
	var err error

	if err := c.db.RunTransaction(ctx, nil, func(tx *sql.Tx) error {
		if moves, err = c.dao.GetMovesForGame(tx, gameId); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return moves, nil
}
