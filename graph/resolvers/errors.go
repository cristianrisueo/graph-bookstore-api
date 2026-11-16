package resolvers

import (
	"context"
	"errors"
	"log/slog"
	"runtime/debug"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
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

// validationError es un error de validación de entrada ligado a un campo
// concreto. Los resolvers lo crean con inputError; el presenter se encarga
// de traducirlo a BAD_USER_INPUT con extensions.field.
type validationError struct {
	field   string
	message string
}

func (e *validationError) Error() string { return e.message }

// inputError crea un error de validación para que los resolvers lo
// devuelvan cuando el valor de un campo de entrada no es válido.
func inputError(field, msg string) error {
	return &validationError{field: field, message: msg}
}

// NewErrorPresenter clasifica los errores que llegan al cliente GraphQL y
// les asigna un código en extensions.code, sin filtrar detalles internos.
func NewErrorPresenter(logger *slog.Logger) graphql.ErrorPresenterFunc {
	return func(ctx context.Context, err error) *gqlerror.Error {
		gqlErr := graphql.DefaultErrorPresenter(ctx, err)
		if gqlErr == nil {
			return nil
		}

		cause := gqlErr.Err
		if cause == nil {
			cause = err
		}

		var valErr *validationError
		var numErr *strconv.NumError

		switch {
		case errors.Is(cause, store.ErrAuthorNotFound), errors.Is(cause, store.ErrBookNotFound):
			setCode(gqlErr, codeNotFound)

		case errors.Is(cause, store.ErrAuthorHasBooks):
			setCode(gqlErr, codeAuthorHasBooks)

		case errors.As(cause, &valErr):
			gqlErr.Message = valErr.message
			setCode(gqlErr, codeBadUserInput)
			setExtension(gqlErr, "field", valErr.field)

		case errors.As(cause, &numErr):
			// Error de conversión de argumentos generado por gqlgen
			// (p. ej. id: "abc"), sin detalles de Go para el cliente.
			gqlErr.Message = "invalid value for argument " + argumentName(gqlErr.Path)
			setCode(gqlErr, codeBadUserInput)

		default:
			logger.Error("unhandled graphql error", "path", gqlErr.Path.String(), "err", cause)
			gqlErr.Message = "internal server error"
			setCode(gqlErr, codeInternalServerError)
		}

		return gqlErr
	}
}

// NewRecoverFunc registra el panic y devuelve un error genérico que el
// presenter termina mostrando como INTERNAL_SERVER_ERROR.
func NewRecoverFunc(logger *slog.Logger) graphql.RecoverFunc {
	return func(ctx context.Context, err any) error {
		logger.Error("panic recovered",
			"path", graphql.GetPath(ctx).String(),
			"err", err,
			"stack", string(debug.Stack()),
		)
		return errors.New("internal server error")
	}
}

// setCode añade o sobrescribe extensions.code en err.
func setCode(err *gqlerror.Error, code string) {
	if err.Extensions == nil {
		err.Extensions = map[string]any{}
	}
	err.Extensions["code"] = code
}

// setExtension añade o sobrescribe una clave de extensions en err.
func setExtension(err *gqlerror.Error, key string, value any) {
	if err.Extensions == nil {
		err.Extensions = map[string]any{}
	}
	err.Extensions[key] = value
}

// argumentName devuelve el nombre del último segmento de path, que es el
// argumento cuyo valor gqlgen no pudo convertir.
func argumentName(path ast.Path) string {
	if len(path) == 0 {
		return "input"
	}
	if name, ok := path[len(path)-1].(ast.PathName); ok {
		return string(name)
	}
	return "input"
}
