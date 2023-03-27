package middleware

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/aeon-fruit/dalil.git/internal/pkg/common/constants"
	reqctx "github.com/aeon-fruit/dalil.git/internal/pkg/context/request"
	errorModel "github.com/aeon-fruit/dalil.git/internal/pkg/model/error"
	"github.com/aeon-fruit/dalil.git/internal/pkg/model/marshaller"
	"github.com/aeon-fruit/dalil.git/internal/pkg/urlparams"
	"github.com/go-logr/logr"
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

			ctx := r.Context()
			if !isValid(ctx, value, pattern) {
				_ = marshaller.SerializeError(w, errorModel.New(http.StatusBadRequest,
					fmt.Sprintf("%v could not be retrieved", key)))
				return
			}

			ctx = reqctx.SetPathParam(ctx, key, value)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func isValid(ctx context.Context, value string, pattern string) bool {
	if value == "" {
		return false
	}

	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return true
	}

	match, err := regexp.MatchString(pattern, value)
	if err != nil {
		logr.FromContextOrDiscard(ctx).Error(err, "Regex matching failed",
			constants.Pattern, pattern,
			constants.Value, value)
	}
	return match && err == nil
}
