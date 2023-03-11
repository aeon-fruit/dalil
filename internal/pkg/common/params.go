package common

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func GetPathParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func GetPathParamInt(r *http.Request, key string) int {
	return toInt(chi.URLParam(r, key))
}

func GetQueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func GetQueryParamInt(r *http.Request, key string) int {
	return toInt(r.URL.Query().Get(key))
}

func GetQueryParamBool(r *http.Request, key string) bool {
	query := r.URL.Query()
	if !query.Has(key) {
		return false
	}
	return toBool(query.Get(key))
}

func toInt(strValue string) int {
	value, err := strconv.Atoi(strings.ToLower(strValue))
	if err != nil {
		return -1
	}
	return value
}

func toBool(strValue string) bool {
	value, err := strconv.ParseBool(strings.ToLower(strValue))
	if err != nil {
		return len(strValue) == 0
	}
	return value
}
