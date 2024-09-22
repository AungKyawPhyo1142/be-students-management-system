package helpers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func GetIDFromParams(r *http.Request) (int, error) {
	idParams := chi.URLParam(r, "id")
	id, error := strconv.Atoi(idParams)

	return id, error
}
