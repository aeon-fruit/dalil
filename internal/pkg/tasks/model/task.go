package model

import (
	"time"

	"github.com/aeon-fruit/dalil.git/internal/pkg/tasks/dao/entity"
)

type GetTaskResponse struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	StatusId    int       `json:"statusId"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func EntityToGetTaskResponse(entity entity.Task) GetTaskResponse {
	return GetTaskResponse{
		Id:          entity.Id,
		Name:        entity.Name,
		StatusId:    entity.StatusId,
		Description: entity.Description,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

type UpsertTaskRequest struct {
	Id          *int   `json:"id,omitempty"`
	Name        string `json:"name"`
	StatusId    int    `json:"statusId"`
	Description string `json:"description,omitempty"`
}

func (dto *UpsertTaskRequest) IsValid(id *int) bool {
	return dto != nil && ((id == nil && dto.Id == nil) || (*id == *dto.Id))
}

func (dto *UpsertTaskRequest) ToEntity() entity.Task {
	var id int
	if dto.Id != nil {
		id = *dto.Id
	}

	return entity.Task{
		Id:          id,
		Name:        dto.Name,
		StatusId:    dto.StatusId,
		Description: dto.Description,
	}
}
