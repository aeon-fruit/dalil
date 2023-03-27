package middleware

import (
	"net/http"

	"github.com/aeon-fruit/dalil.git/internal/pkg/common/constants"
	"github.com/aeon-fruit/dalil.git/internal/pkg/log"
	"github.com/go-logr/logr"
)

func LoggingContext(logger log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = logr.NewContext(ctx, logger.WithName(constants.AppName))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
