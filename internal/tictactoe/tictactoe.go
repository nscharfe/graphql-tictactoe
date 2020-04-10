package tictactoe

import (
	"errors"

	"github.com/nscharfe/graphql-tictactoe/graph/model"
)

var (
	errorGameNotActive = errors.New("game not active")
	errorCollision     = errors.New("move already occupied")
	errorWrongTurn     = errors.New("player cannot move out of turn")
)

type tictactoe struct {
	board [][]model.Player
	moves []*model.Move
}

type GameClient interface {
	MakeMove(*model.Move) error
	CurrentTurn() model.Player
	Winner() *model.Player
	Status() model.Status
}

func New(moves []*model.Move) GameClient {
	board := make([][]model.Player, 3)
	for i := range board {
		board[i] = make([]model.Player, 3)
	}

	for _, move := range moves {
		board[move.RowIndex][move.ColumnIndex] = move.Player
	}

	return &tictactoe{
		board: board,
		moves: moves,
	}
}

func (t *tictactoe) MakeMove(move *model.Move) error {
	if err := t.isValidMove(move); err != nil {
		return err
	}

	t.board[move.RowIndex][move.ColumnIndex] = move.Player
	t.moves = append(t.moves, move)
	return nil
}

func (t *tictactoe) CurrentTurn() model.Player {
	if len(t.moves) == 0 {
		return model.PlayerX
	}

	lastTurn := t.moves[len(t.moves)-1]
	switch lastTurn.Player {
	case model.PlayerX:
		return model.PlayerO
	case model.PlayerO:
		return model.PlayerX
	default:
		panic("should never happen")
	}
}

func (t *tictactoe) Winner() *model.Player {
	// hacky way of determining winner

	// first check for row / column win
	for i := 0; i < 3; i++ {
		if t.board[i][0].String() != "" && t.board[i][0] == t.board[i][1] && t.board[i][1] == t.board[i][2] {
			return &t.board[i][0]
		}
		if t.board[i][0].String() != "" && t.board[0][i] == t.board[1][i] && t.board[1][i] == t.board[2][i] {
			return &t.board[i][0]
		}
	}

	// top left to bottom right
	if t.board[0][0].String() != "" && t.board[0][0] == t.board[1][1] && t.board[1][1] == t.board[2][2] {
		return &t.board[0][0]
	}

	// bottom left to top right
	if t.board[0][2].String() != "" && t.board[0][2] == t.board[1][1] && t.board[1][1] == t.board[2][0] {
		return &t.board[0][0]
	}

	return nil
}

func (t *tictactoe) Status() model.Status {
	if t.isWinner() {
		return model.StatusWinner
	}
	if t.allMovesTaken() {
		return model.StatusDraw
	}
	return model.StatusActive
}

func (t *tictactoe) isValidMove(move *model.Move) error {
	if t.isGameOver() {
		return errorGameNotActive
	}
	if t.board[move.RowIndex][move.ColumnIndex].String() != "" {
		return errorCollision
	}
	if move.Player != t.CurrentTurn() {
		return errorWrongTurn
	}
	return nil
}

func (t *tictactoe) isGameOver() bool {
	return t.isWinner() || t.allMovesTaken()
}

func (t *tictactoe) isWinner() bool {
	return t.Winner() != nil
}

func (t *tictactoe) allMovesTaken() bool {
	for _, row := range t.board {
		for _, cell := range row {
			if cell.String() == "" {
				return false
			}
		}
	}

	return true
}
