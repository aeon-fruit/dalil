package marshaller

import (
	"encoding/json"
	"net/http"

	errorModel "github.com/aeon-fruit/dalil.git/internal/pkg/model/error"
)

func SerializeEntity(w http.ResponseWriter, entity any) error {
	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(entity)
}

func SerializeError(w http.ResponseWriter, errorResponse errorModel.Response) error {
	w.WriteHeader(errorResponse.Code)
	err := SerializeEntity(w, errorResponse)
	if err != nil {
		httpStatus := http.StatusInternalServerError
		http.Error(w, http.StatusText(httpStatus), httpStatus)
	}
	return err
}

func SerializeFlatError(w http.ResponseWriter, httpStatusCode int, message string) error {
	return SerializeError(w, errorModel.New(httpStatusCode).WithMessage(message))
}
