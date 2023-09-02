package graph

import (
	"github.com/vikstrous/dataloadgen-example/graph/storage"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct {
	UserStorage *storage.UserStorage
	TodoStorage *storage.TodoStorage
}
