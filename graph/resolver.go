package graph

import "github.com/nscharfe/graphql-tictactoe/internal/backend"

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Backend backend.Client
}
