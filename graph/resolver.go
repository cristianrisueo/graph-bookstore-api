package graph

import "github.com/cristianrisueo/graph-bookstore-api/store"

// Este archivo no se regenera automáticamente.
//
// Sirve como inyección de dependencias de la aplicación: añade aquí las
// dependencias que necesiten los resolvers.

type Resolver struct {
	Store *store.Store
}
