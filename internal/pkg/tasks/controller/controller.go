package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aeon-fruit/dalil.git/internal/pkg/common/constants"
	"github.com/aeon-fruit/dalil.git/internal/pkg/common/errors"
	reqctx "github.com/aeon-fruit/dalil.git/internal/pkg/context/request"
	errorModel "github.com/aeon-fruit/dalil.git/internal/pkg/model/error"
	"github.com/aeon-fruit/dalil.git/internal/pkg/model/marshaller"
	model "github.com/aeon-fruit/dalil.git/internal/pkg/tasks/model"
	service "github.com/aeon-fruit/dalil.git/internal/pkg/tasks/service"
	"github.com/go-logr/logr"
)

const (
	getByIdFailed    = "GetById failed"
	getByIdResponse  = "GetById response"
	getAllFailed     = "GetAll failed"
	getAllResponse   = "GetAll response"
	addFailed        = "Add failed"
	addResponse      = "Add response"
	updateFailed     = "Update failed"
	updateResponse   = "Update response"
	removeByIdFailed = "RemoveById failed"
)

type Controller interface {
	GetById(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	Add(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	RemoveById(w http.ResponseWriter, r *http.Request)
}

type controllerImpl struct {
	service service.Service
}

type ControllerOption func(*controllerImpl)

func New(options ...ControllerOption) Controller {
	instance := controllerImpl{}

	for _, option := range options {
		if option != nil {
			option(&instance)
		}
	}

	return &instance
}

func WithService(service service.Service) ControllerOption {
	return func(controller *controllerImpl) {
		if controller != nil {
			controller.service = service
		}
	}
}

func (ctrl *controllerImpl) GetById(w http.ResponseWriter, r *http.Request) {
	logger := logr.FromContextOrDiscard(r.Context())

	id, stop := getIdOrStop(w, r)
	if stop {
		return
	}

	entity, err := ctrl.service.GetById(id)
	if err != nil {
		if err == errors.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			logger.Error(err, getAllFailed, constants.Id, id)
			_ = marshaller.SerializeError(w, errorModel.New(http.StatusInternalServerError, err.Error()))
		}
		return
	}

	logger.V(1).Info(getByIdResponse, constants.Payload, entity)

	_ = marshaller.SerializeEntity(w, entity)
}

func (ctrl *controllerImpl) GetAll(w http.ResponseWriter, r *http.Request) {
	logger := logr.FromContextOrDiscard(r.Context())

	entity, err := ctrl.service.GetAll()
	if err != nil {
		logger.Error(err, getAllFailed)
		_ = marshaller.SerializeError(w, errorModel.New(http.StatusInternalServerError, err.Error()))
		return
	}

	if len(entity) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	logger.V(1).Info(getAllResponse, constants.Payload, entity)

	_ = marshaller.SerializeEntity(w, entity)
}

func (ctrl *controllerImpl) Add(w http.ResponseWriter, r *http.Request) {
	logger := logr.FromContextOrDiscard(r.Context())

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error(err, addFailed)
		_ = marshaller.SerializeError(w, errorModel.New(http.StatusBadRequest, err.Error()))
		return
	}

	request := model.UpsertTaskRequest{}
	err = json.Unmarshal(body, &request)
	if err != nil {
		logger.Error(err, addFailed, constants.Body, string(body))
		_ = marshaller.SerializeError(w, errorModel.New(http.StatusBadRequest, err.Error()))
		return
	}

	if !request.IsValid(nil) {
		logger.Error(errors.ErrInvalidArgument, addFailed, constants.Field, constants.Id)
		_ = marshaller.SerializeError(w, errorModel.New(http.StatusBadRequest, "Invalid request content"))
		return
	}

	request.Id = nil

	entity, err := ctrl.service.Upsert(request)
	if err != nil {
		logger.Error(err, addFailed)
		_ = marshaller.SerializeError(w, errorModel.New(http.StatusInternalServerError, err.Error()))
		return
	}

	location := fmt.Sprintf("%s/%d", r.Host, entity.Id)
	logger.V(1).Info("Added entity location", constants.Location, location)
	logger.V(1).Info(addResponse, constants.Payload, entity)

	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)
	_ = marshaller.SerializeEntity(w, entity)
}

func (ctrl *controllerImpl) Update(w http.ResponseWriter, r *http.Request) {
	logger := logr.FromContextOrDiscard(r.Context())

	id, stop := getIdOrStop(w, r)
	if stop {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error(err, addFailed)
		_ = marshaller.SerializeError(w, errorModel.New(http.StatusBadRequest, err.Error()))
		return
	}

	request := model.UpsertTaskRequest{}
	err = json.Unmarshal(body, &request)
	if err != nil {
		logger.Error(err, updateFailed, constants.Body, string(body))
		_ = marshaller.SerializeError(w, errorModel.New(http.StatusBadRequest, err.Error()))
		return
	}

	if !request.IsValid(&id) {
		logger.Error(errors.ErrInvalidArgument, updateFailed, constants.Field, constants.Id)
		_ = marshaller.SerializeError(w, errorModel.New(http.StatusBadRequest, "Invalid request content"))
		return
	}

	request.Id = &id

	entity, err := ctrl.service.Upsert(request)
	if err != nil {
		if err == errors.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else if err == errors.ErrNotModified {
			w.WriteHeader(http.StatusNotModified)
		} else {
			logger.Error(err, updateFailed)
			_ = marshaller.SerializeError(w, errorModel.New(http.StatusInternalServerError, err.Error()))
		}
		return
	}

	logger.V(1).Info(updateResponse, constants.Payload, entity)

	_ = marshaller.SerializeEntity(w, entity)
}

func (ctrl *controllerImpl) RemoveById(w http.ResponseWriter, r *http.Request) {
	logger := logr.FromContextOrDiscard(r.Context())

	id, stop := getIdOrStop(w, r)
	if stop {
		return
	}

	err := ctrl.service.RemoveById(id)
	if err != nil {
		if err == errors.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			logger.Error(err, removeByIdFailed, constants.Id, id)
			_ = marshaller.SerializeError(w, errorModel.New(http.StatusInternalServerError, err.Error()))
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getIdOrStop(w http.ResponseWriter, r *http.Request) (id int, stop bool) {
	ctx := r.Context()
	value, err := reqctx.GetPathParam(ctx, constants.Id)
	if err == nil {
		id, err = value.Int()
		if err == nil {
			return id, false
		}
	}

	logger := logr.FromContextOrDiscard(ctx)
	logger.Error(err, "Cannot retrieve field from context", constants.Field, constants.Id)

	_ = marshaller.SerializeError(w, errorModel.New(http.StatusInternalServerError,
		"Unable to retrieve the Task Id"))
	return 0, true
}
