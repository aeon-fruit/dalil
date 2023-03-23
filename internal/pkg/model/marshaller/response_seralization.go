package marshaller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aeon-fruit/dalil.git/internal/pkg/common/errors"
	errorModel "github.com/aeon-fruit/dalil.git/internal/pkg/model/error"
)

func SerializeEntity(w http.ResponseWriter, entity any) error {
	if w == nil {
		return fmt.Errorf("%w: 1st argument should be non-nil", errors.ErrInvalidArgument)
	}

	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(entity)
}

func SerializeError(w http.ResponseWriter, errorResponse errorModel.Response) error {
	if w == nil {
		return fmt.Errorf("%w: 1st argument should be non-nil", errors.ErrInvalidArgument)
	}

	statusCode := errorResponse.Code
	if http.StatusText(statusCode) == "" {
		statusCode = http.StatusInternalServerError
	}
	w.WriteHeader(statusCode)

	err := SerializeEntity(w, errorResponse)
	if err != nil {
		httpStatus := http.StatusInternalServerError
		http.Error(w, http.StatusText(httpStatus), httpStatus)
	}
	return err
}
