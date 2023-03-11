package middleware

import (
	"context"
	"net/http"

	"github.com/aeon-fruit/dalil.git/internal/pkg/common"
)

type _idKey struct{}

func IdContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := common.GetPathParamInt(r, "id")
		if id == -1 {
			_ = common.SerializeFlatError(w, http.StatusBadRequest, "Task Id should be an integer")
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, _idKey{}, id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetId(ctx context.Context) (int, error) {
	if ctx == nil {
		return 0, common.ErrNotFound
	}

	value, ok := ctx.Value(_idKey{}).(int)
	if !ok {
		return 0, common.ErrNotFound
	}

	return value, nil

}
