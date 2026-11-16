package resolvers

import "github.com/cristianrisueo/graph-bookstore-api/store"

// Este archivo sirve como inyección de dependencias para los resolvers.
// Recuerdo que en GraphQL, un resolver es el equivalente a un handler en REST.

type Resolver struct {
	Store *store.Store
}
