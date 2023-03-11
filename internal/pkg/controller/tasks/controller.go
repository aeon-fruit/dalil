package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aeon-fruit/dalil.git/internal/pkg/common"
	"github.com/aeon-fruit/dalil.git/internal/pkg/middleware"
	model "github.com/aeon-fruit/dalil.git/internal/pkg/model/tasks"
	service "github.com/aeon-fruit/dalil.git/internal/pkg/service/tasks"
)

type Controller interface {
	GetById(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	Add(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	RemoveById(w http.ResponseWriter, r *http.Request)
	RemoveByIds(w http.ResponseWriter, r *http.Request)
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
	id, stop := getIdOrStop(w, r)
	if stop {
		return
	}

	entity, err := ctrl.service.GetById(id)
	if err != nil {
		if err == common.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			_ = common.SerializeFlatError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	_ = common.SerializeEntity(w, entity)
}

func (ctrl *controllerImpl) GetAll(w http.ResponseWriter, r *http.Request) {
	entity, err := ctrl.service.GetAll()
	if err != nil {
		_ = common.SerializeFlatError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(entity) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_ = common.SerializeEntity(w, entity)
}

func (ctrl *controllerImpl) Add(w http.ResponseWriter, r *http.Request) {
	request := model.UpsertTaskRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		_ = common.SerializeFlatError(w, http.StatusBadRequest, err.Error())
		return
	}

	if !request.IsValid(nil) {
		_ = common.SerializeFlatError(w, http.StatusBadRequest, "Invalid request content")
		return
	}

	request.Id = nil

	entity, err := ctrl.service.Upsert(request)
	if err != nil {
		_ = common.SerializeFlatError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, entity.Id))
	w.WriteHeader(http.StatusCreated)
	_ = common.SerializeEntity(w, entity)
}

func (ctrl *controllerImpl) Update(w http.ResponseWriter, r *http.Request) {
	id, stop := getIdOrStop(w, r)
	if stop {
		return
	}

	request := model.UpsertTaskRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		_ = common.SerializeFlatError(w, http.StatusBadRequest, err.Error())
		return
	}

	if !request.IsValid(&id) {
		_ = common.SerializeFlatError(w, http.StatusBadRequest, "Invalid request content")
		return
	}

	request.Id = &id

	entity, err := ctrl.service.Upsert(request)
	if err != nil {
		if err == common.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			_ = common.SerializeFlatError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	_ = common.SerializeEntity(w, entity)
}

func (ctrl *controllerImpl) RemoveById(w http.ResponseWriter, r *http.Request) {
	id, stop := getIdOrStop(w, r)
	if stop {
		return
	}

	err := ctrl.service.RemoveById(id)
	if err != nil {
		if err == common.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			_ = common.SerializeFlatError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ctrl *controllerImpl) RemoveByIds(w http.ResponseWriter, r *http.Request) {
	panic("Unimplemented")
}

func getIdOrStop(w http.ResponseWriter, r *http.Request) (id int, stop bool) {
	id, err := middleware.GetId(r.Context())
	if err == nil {
		return
	}

	stop = true
	if err == common.ErrNotFound {
		_ = common.SerializeFlatError(w, http.StatusInternalServerError, "Unable to retrieve the Task Id")
	} else {
		_ = common.SerializeFlatError(w, http.StatusInternalServerError, err.Error())
	}
	return
}
