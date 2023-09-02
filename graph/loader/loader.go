package loader

import (
	"context"
	"net/http"

	"github.com/vikstrous/dataloadgen"
	"github.com/vikstrous/dataloadgen-example/graph/model"
	"github.com/vikstrous/dataloadgen-example/graph/storage"
)

type ctxKey struct{}

type Loaders struct {
	User *dataloadgen.Loader[string, *model.User]
}

// Middleware injects data loaders into the context
func Middleware(userStorage *storage.UserStorage, next http.Handler) http.Handler {
	userFetcher := &userFetcher{
		userStorage: userStorage,
	}
	// return a middleware that injects the loader to the request context
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note that the loaders are being created per-request. This is important because they contain caching and batching logic that must be request-scoped.
		loaders := &Loaders{
			User: dataloadgen.NewLoader(userFetcher.fetch),
		}
		r = r.WithContext(context.WithValue(r.Context(), ctxKey{}, loaders))
		next.ServeHTTP(w, r)
	})
}

// Get returns the dataloader for a given context
func Get(ctx context.Context) *Loaders {
	return ctx.Value(ctxKey{}).(*Loaders)
}

type userFetcher struct {
	userStorage *storage.UserStorage
}

func (u *userFetcher) fetch(ctx context.Context, userIDs []string) ([]*model.User, []error) {
	users, errs := u.userStorage.GetMulti(userIDs)
	return users, errs
}