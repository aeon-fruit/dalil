package urlparams

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func ParsePathParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func ParsePathParamInt(r *http.Request, key string) int {
	return toInt(chi.URLParam(r, key))
}

func ParseQueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func ParseQueryParamInt(r *http.Request, key string) int {
	return toInt(r.URL.Query().Get(key))
}

func ParseQueryParamBool(r *http.Request, key string) bool {
	query := r.URL.Query()
	if !query.Has(key) {
		return false
	}

	strValue := query.Get(key)
	value, err := strconv.ParseBool(strings.ToLower(strValue))
	if err != nil {
		return len(strValue) == 0
	}
	return value
}

func toInt(strValue string) int {
	value, err := strconv.Atoi(strings.ToLower(strValue))
	if err != nil {
		return -1
	}
	return value
}
