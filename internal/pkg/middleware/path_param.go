package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	reqctx "github.com/aeon-fruit/dalil.git/internal/pkg/context/request"
	errorModel "github.com/aeon-fruit/dalil.git/internal/pkg/model/error"
	"github.com/aeon-fruit/dalil.git/internal/pkg/model/marshaller"
	"github.com/aeon-fruit/dalil.git/internal/pkg/urlparams"
)

const (
	patternInt = "[0-9]+"
)

func PathParamContextInt(key string) func(http.Handler) http.Handler {
	return PathParamContextString(key, patternInt)
}

func PathParamContextString(key string, pattern string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			value := urlparams.ParsePathParam(r, key)

			if !isValid(value, pattern) {
				_ = marshaller.SerializeError(w, errorModel.New(http.StatusBadRequest,
					fmt.Sprintf("%v could not be retrieved", key)))
				return
			}

			ctx := reqctx.SetPathParam(r.Context(), key, value)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func isValid(value string, pattern string) bool {
	if value == "" {
		return false
	}

	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return true
	}

	match, err := regexp.MatchString(pattern, value)
	return match && err == nil
}
