package loader

import (
	"context"
	"net/http"
	"time"

	"github.com/vikstrous/dataloadgen"
	"github.com/vikstrous/dataloadgen-example/graph/model"
	"github.com/vikstrous/dataloadgen-example/graph/storage"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

// userReader is used to give the getUsers function access to the underlying storage system and can be used to group multiple related fetch functions with similar storage system access patterns.
type userReader struct {
	userStorage *storage.UserStorage
}

// getUsers retrieves multiple users at the same time from the underlying storage system.
func (u userReader) getUsers(ctx context.Context, userIDs []string) ([]*model.User, []error) {
	users, errs := u.userStorage.GetMulti(userIDs)
	return users, errs
}

// Loaders wrap your data loaders to inject via middleware
type Loaders struct {
	UserLoader *dataloadgen.Loader[string, *model.User]
}

// NewLoaders instantiates data loaders for the middleware
func NewLoaders(s *storage.UserStorage) *Loaders {
	// define the data loader
	ur := &userReader{userStorage: s}
	return &Loaders{
		UserLoader: dataloadgen.NewLoader(ur.getUsers, dataloadgen.WithWait(time.Millisecond)),
	}
}

// Middleware injects data loaders into the context
func Middleware(userStorage *storage.UserStorage, next http.Handler) http.Handler {
	// return a middleware that injects the loader to the request context
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note that the loaders are being created per-request. This is important because they contain caching and batching logic that must be request-scoped.
		loaders := NewLoaders(userStorage)
		r = r.WithContext(context.WithValue(r.Context(), loadersKey, loaders))
		next.ServeHTTP(w, r)
	})
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

// GetUser returns single user by id efficiently
func GetUser(ctx context.Context, userID string) (*model.User, error) {
	loaders := For(ctx)
	return loaders.UserLoader.Load(ctx, userID)
}

// PrimeUser primes the user loader cache
func PrimeUser(ctx context.Context, user *model.User) bool {
	loaders := For(ctx)
	return loaders.UserLoader.Prime(user.ID, user)
}

// GetUsers returns many users by ids efficiently
func GetUsers(ctx context.Context, userIDs []string) ([]*model.User, error) {
	loaders := For(ctx)
	return loaders.UserLoader.LoadAll(ctx, userIDs)
}
