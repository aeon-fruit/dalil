package middleware

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/aeon-fruit/dalil.git/internal/pkg/common/errors"
	"github.com/aeon-fruit/dalil.git/internal/pkg/model/marshaller"
	"github.com/aeon-fruit/dalil.git/internal/pkg/urlparams"
)

type ctxKey struct {
	key string
}

type proposition func() bool

func PathParamContextInt(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			value := urlparams.ParsePathParamInt(r, key)

			validateAndForward(w, r, next, key, value, intValidator(value))
		})
	}
}

func PathParamContextString(key string, pattern string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			value := urlparams.ParsePathParam(r, key)

			validateAndForward(w, r, next, key, value, stringValidator(value, pattern))
		})
	}
}

func GetPathParamInt(ctx context.Context, key string) (int, error) {
	if ctx == nil {
		return 0, errors.ErrNotFound
	}

	value, ok := ctx.Value(ctxKey{key: key}).(int)
	if !ok {
		return 0, errors.ErrNotFound
	}

	return value, nil
}

func GetPathParamString(ctx context.Context, key string) (string, error) {
	if ctx == nil {
		return "", errors.ErrNotFound
	}

	value, ok := ctx.Value(ctxKey{key: key}).(string)
	if !ok {
		return "", errors.ErrNotFound
	}

	return value, nil
}

func validateAndForward(w http.ResponseWriter, r *http.Request, next http.Handler, key string, value any, isValid proposition) {
	if !isValid() {
		_ = marshaller.SerializeFlatError(w, http.StatusBadRequest,
			fmt.Sprintf("%v could not be retrieved", key))
		return
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, ctxKey{key: key}, value)

	next.ServeHTTP(w, r.WithContext(ctx))
}

func intValidator(value int) proposition {
	return func() bool {
		return value != -1
	}
}

func stringValidator(value string, pattern string) proposition {
	return func() bool {
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
}
