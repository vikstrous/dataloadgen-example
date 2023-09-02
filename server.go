package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vikstrous/dataloadgen-example/graph"
	"github.com/vikstrous/dataloadgen-example/graph/loader"
	"github.com/vikstrous/dataloadgen-example/graph/model"
	"github.com/vikstrous/dataloadgen-example/graph/storage"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	userStorage := storage.NewUserStorage()
	todoStorage := storage.NewTodoStorage()
	err := userStorage.Put(model.User{
		ID:   "alice",
		Name: "Alice",
	})
	if err != nil {
		panic(err)
	}
	err = userStorage.Put(model.User{
		ID:   "bob",
		Name: "Bob",
	})
	if err != nil {
		panic(err)
	}

	srv := loader.Middleware(userStorage, handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		UserStorage: userStorage,
		TodoStorage: todoStorage,
	}})))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
