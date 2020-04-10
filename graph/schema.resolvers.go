package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/nscharfe/graphql-tictactoe/graph/generated"
	"github.com/nscharfe/graphql-tictactoe/graph/model"
)

func (r *gameResolver) Moves(ctx context.Context, obj *model.Game) ([]*model.Move, error) {
	return r.Backend.GetMovesForGame(ctx, obj.ID)
}

func (r *moveResolver) Game(ctx context.Context, obj *model.Move) (*model.Game, error) {
	return r.Backend.GetGame(ctx, obj.GameID)
}

func (r *mutationResolver) NewGame(ctx context.Context) (*model.Game, error) {
	return r.Backend.CreateGame(ctx)
}

func (r *mutationResolver) NewMove(ctx context.Context, gameID string, rowIndex int, columnIndex int, player model.Player) (*model.Move, error) {
	return r.Backend.CreateMove(ctx, gameID, rowIndex, columnIndex, player)
}

func (r *queryResolver) Game(ctx context.Context, id string) (*model.Game, error) {
	return r.Backend.GetGame(ctx, id)
}

// Game returns generated.GameResolver implementation.
func (r *Resolver) Game() generated.GameResolver { return &gameResolver{r} }

// Move returns generated.MoveResolver implementation.
func (r *Resolver) Move() generated.MoveResolver { return &moveResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type gameResolver struct{ *Resolver }
type moveResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
