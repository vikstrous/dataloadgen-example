package loader

import (
	"context"
	"net/http"

	"github.com/vikstrous/dataloadgen"
	"github.com/vikstrous/dataloadgen-example/graph/model"
	"github.com/vikstrous/dataloadgen-example/graph/storage"
)

// Get returns the Loaders bundle from the context. It must be used only in graphql resolvers where Middleware has put the Loaders struct into the context already.
func Get(ctx context.Context) *Loaders {
	return ctx.Value(ctxKey{}).(*Loaders)
}

// Loaders provide access for loading various objects from the underlying object's storage system while batching concurrent requests and caching responses.
type Loaders struct {
	User *dataloadgen.Loader[string, *model.User]
}

// Middleware injects data loaders into the context
func Middleware(userStorage *storage.UserStorage, next http.Handler) http.Handler {
	// return a middleware that injects the loader to the request context
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note that the loaders are being created per-request. This is important because they contain caching and batching logic that must be request-scoped.
		loaders := newLoaders(userStorage)
		r = r.WithContext(context.WithValue(r.Context(), ctxKey{}, loaders))
		next.ServeHTTP(w, r)
	})
}

type ctxKey struct{}

// newLoaders creates the Loaders struct
func newLoaders(userStorage *storage.UserStorage) *Loaders {
	userFetcher := userFetcher{userStorage: userStorage}
	loaders := &Loaders{
		User: dataloadgen.NewLoader(userFetcher.fetchUsers),
	}
	return loaders
}

// userFetcher is used to give the fetchUsers function access to the underlying storage system and can be used to group multiple related fetch functions with similar storage system access patterns.
type userFetcher struct {
	userStorage *storage.UserStorage
}

// fetchUsers retrieves multiple users at the same time from the underlying storage system.
func (u userFetcher) fetchUsers(ctx context.Context, userIDs []string) ([]*model.User, []error) {
	users, errs := u.userStorage.GetMulti(userIDs)
	return users, errs
}
