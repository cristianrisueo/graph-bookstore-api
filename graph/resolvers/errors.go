package resolvers

import (
	"context"
	"errors"
	"log/slog"
	"runtime/debug"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/cristianrisueo/graph-bookstore-api/store"
)

// Códigos de error expuestos en extensions.code.
const (
	codeNotFound            = "NOT_FOUND"
	codeAuthorHasBooks      = "AUTHOR_HAS_BOOKS"
	codeBadUserInput        = "BAD_USER_INPUT"
	codeInternalServerError = "INTERNAL_SERVER_ERROR"
)

// errInternal es el error que devuelve NewRecoverFunc tras registrar un panic.
var errInternal = errors.New("internal server error")

// inputError crea un error de validación ligado a un campo de entrada.
func inputError(field, msg string) error {
	return &gqlerror.Error{
		Message:    msg,
		Extensions: map[string]any{"code": codeBadUserInput, "field": field},
	}
}

// NewErrorPresenter asigna extensions.code a cada error y oculta los detalles internos.
func NewErrorPresenter(logger *slog.Logger) graphql.ErrorPresenterFunc {
	return func(ctx context.Context, err error) *gqlerror.Error {
		gqlErr := graphql.DefaultErrorPresenter(ctx, err)
		var numErr *strconv.NumError

		switch {
		// Errores de dominio conocidos.
		case errors.Is(err, store.ErrAuthorNotFound), errors.Is(err, store.ErrBookNotFound):
			setExtension(gqlErr, "code", codeNotFound)
		case errors.Is(err, store.ErrAuthorHasBooks):
			setExtension(gqlErr, "code", codeAuthorHasBooks)

		// Ya traen código: los de inputError y los de parseo/validación de gqlgen.
		case gqlErr.Extensions["code"] != nil:

		// ID que no se puede convertir a número (graphql.IntID). Es el único error de
		// conversión que ocurre al ejecutar; los demás tipos los rechaza la validación.
		case errors.As(err, &numErr):
			gqlErr.Message = "invalid ID"
			setExtension(gqlErr, "code", codeBadUserInput)

		// Panic ya registrado por NewRecoverFunc.
		case errors.Is(err, errInternal):
			setExtension(gqlErr, "code", codeInternalServerError)

		// Cualquier otro error: se registra y se oculta al cliente.
		default:
			logger.Error("unhandled graphql error", "path", gqlErr.Path.String(), "err", err)
			gqlErr.Message = errInternal.Error()
			setExtension(gqlErr, "code", codeInternalServerError)
		}

		return gqlErr
	}
}

// NewRecoverFunc registra el panic con su traza y devuelve errInternal.
func NewRecoverFunc(logger *slog.Logger) graphql.RecoverFunc {
	return func(ctx context.Context, p any) error {
		logger.Error("panic recovered",
			"path", graphql.GetPath(ctx).String(),
			"panic", p,
			"stack", string(debug.Stack()),
		)
		return errInternal
	}
}

// setExtension añade o sobrescribe una clave de extensions.
func setExtension(err *gqlerror.Error, key string, value any) {
	if err.Extensions == nil {
		err.Extensions = map[string]any{}
	}
	err.Extensions[key] = value
}