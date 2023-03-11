package entity

import "time"

type Task struct {
	Id          int       `json:"id" gorm:"column:id;type:int;primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"column:name;type:varchar;size:255"`
	Status      Status    `json:"status"`
	StatusId    int       `json:"statusId" gorm:"column:status_id;type:text"`
	Description string    `json:"description,omitempty" gorm:"column:description;type:varchar;size:255"`
	CreatedAt   time.Time `json:"createdAt,omitempty" gorm:"column:created_at;type:timestamp"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" gorm:"column:updated_at;type:timestamp"`
}
