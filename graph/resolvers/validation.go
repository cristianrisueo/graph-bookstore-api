package resolvers

import (
	"fmt"
	"strings"
	"time"
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

// validateBookTitle recorta el título y comprueba que sea válido. Devuelve
// el título recortado y, si no es válido, un inputError listo para devolver.
func validateBookTitle(title string) (string, error) {
	title = strings.TrimSpace(title)

	if title == "" {
		return "", inputError("title", "title must not be empty")
	}
	if utf8.RuneCountInString(title) > 300 {
		return "", inputError("title", "title must be at most 300 characters")
	}

	return title, nil
}

// validateBookYear comprueba que el año esté entre 1 y el año actual.
func validateBookYear(year int) error {
	currentYear := time.Now().Year()

	if year < 1 || year > currentYear {
		return inputError("year", fmt.Sprintf("year must be between 1 and %d", currentYear))
	}

	return nil
}
