package reqctx

import (
	"context"
	"reflect"
	"strconv"

	"github.com/aeon-fruit/dalil.git/internal/pkg/common/errors"
)

type ctxKey struct {
	key string
}

type PathParam struct {
	value  string
	parsed map[reflect.Kind]any
}

func (param PathParam) Int() (int, error) {
	value, found := param.parsed[reflect.Int]
	if found {
		return value.(int), nil
	}

	intValue, err := strconv.Atoi(param.value)
	if err != nil {
		return 0, errors.ErrNotFound
	}

	if param.parsed == nil {
		param.parsed = map[reflect.Kind]any{}
	}

	param.parsed[reflect.Int] = intValue
	return intValue, nil

}

func (param PathParam) String() string {
	return param.value
}

func SetPathParam(ctx context.Context, key string, value string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, ctxKey{key: key}, PathParam{value: value})
}

func GetPathParam(ctx context.Context, key string) (PathParam, error) {
	if ctx == nil {
		return PathParam{}, errors.ErrNotFound
	}

	value, ok := ctx.Value(ctxKey{key: key}).(PathParam)
	if !ok {
		return PathParam{}, errors.ErrNotFound
	}

	return value, nil
}
