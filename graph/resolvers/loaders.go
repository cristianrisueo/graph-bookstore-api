package resolvers

// Este archivo es propio del paquete: gqlgen no lo genera ni lo toca.

import (
	"context"
	"net/http"
	"time"

	"github.com/vikstrous/dataloadgen"

	"github.com/cristianrisueo/graph-bookstore-api/store"
)

// loaders agrupa los dataloaders de una petición.
type loaders struct {
	authorByID      *dataloadgen.Loader[int, *store.Author]
	booksByAuthorID *dataloadgen.Loader[int, []*store.Book]
}

// newLoaders crea un conjunto de loaders que agrupan las peticiones recibidas
// en el mismo milisegundo en una sola consulta por lote.
func newLoaders(st *store.Store) *loaders {
	return &loaders{
		authorByID: dataloadgen.NewLoader(
			func(ctx context.Context, ids []int) ([]*store.Author, []error) {
				authors, err := st.GetAuthorsByIDs(ctx, ids)
				if err != nil {
					errs := make([]error, len(ids))
					for i := range errs {
						errs[i] = err
					}
					return make([]*store.Author, len(ids)), errs
				}

				byID := make(map[int]*store.Author, len(authors))
				for _, a := range authors {
					byID[a.ID] = a
				}

				results := make([]*store.Author, len(ids))
				errs := make([]error, len(ids))
				for i, id := range ids {
					a, ok := byID[id]
					if !ok {
						errs[i] = store.ErrAuthorNotFound
						continue
					}
					results[i] = a
				}
				return results, errs
			},
			dataloadgen.WithWait(time.Millisecond),
		),
		booksByAuthorID: dataloadgen.NewLoader(
			func(ctx context.Context, authorIDs []int) ([][]*store.Book, []error) {
				books, err := st.ListBooksByAuthorIDs(ctx, authorIDs)
				if err != nil {
					errs := make([]error, len(authorIDs))
					for i := range errs {
						errs[i] = err
					}
					return make([][]*store.Book, len(authorIDs)), errs
				}

				byAuthorID := make(map[int][]*store.Book, len(authorIDs))
				for _, b := range books {
					byAuthorID[b.AuthorID] = append(byAuthorID[b.AuthorID], b)
				}

				results := make([][]*store.Book, len(authorIDs))
				for i, id := range authorIDs {
					results[i] = byAuthorID[id]
				}
				return results, make([]error, len(authorIDs))
			},
			dataloadgen.WithWait(time.Millisecond),
		),
	}
}

// loadersCtxKey es una clave de tipo privado para guardar los loaders en el
// contexto, evitando colisiones con otras claves.
type loadersCtxKey struct{}

// LoadersMiddleware crea unos loaders nuevos por petición y los guarda en su
// contexto. La caché de un dataloader no debe sobrevivir a la petición: si se
// compartiera entre peticiones, una mutación en una petición dejaría datos
// obsoletos en la caché de otra.
func LoadersMiddleware(st *store.Store, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), loadersCtxKey{}, newLoaders(st))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// loadersFor recupera los loaders de la petición actual del contexto.
func loadersFor(ctx context.Context) *loaders {
	return ctx.Value(loadersCtxKey{}).(*loaders)
}
