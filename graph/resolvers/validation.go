package resolvers

import (
	"strings"
	"unicode/utf8"
)

// validateAuthorName recorta el nombre y comprueba que sea válido. Devuelve
// el nombre recortado y, si no es válido, un inputError listo para devolver.
func validateAuthorName(name string) (string, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return "", inputError("name", "name must not be empty")
	}
	if utf8.RuneCountInString(name) > 200 {
		return "", inputError("name", "name must be at most 200 characters")
	}

	return name, nil
}
