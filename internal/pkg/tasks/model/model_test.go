package model_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aeon-fruit/dalil.git/internal/pkg/tasks/dao/entity"
	"github.com/aeon-fruit/dalil.git/internal/pkg/tasks/model"
)

var _ = Describe("Model", func() {

	const (
		id          = 10
		name        = "Task"
		description = "Task description"
		statusId    = 1
	)

	Describe("EntityToGetTaskResponse", func() {

		When("called", func() {
			It("returns a TaskResponse that contains the same values as the entity's matching fields", func() {
				now := time.Now()
				e := entity.Task{
					Id:   id,
					Name: name,
					Status: entity.Status{
						Id:          statusId,
						Name:        "Status",
						Description: "Status description",
						CreatedAt:   now.Add(-time.Hour),
						UpdatedAt:   now.Add(-time.Hour),
					},
					StatusId:    statusId,
					Description: description,
					CreatedAt:   now.Add(-25 * time.Minute),
					UpdatedAt:   now,
				}

				expected := model.GetTaskResponse{
					Id:          id,
					Name:        name,
					StatusId:    statusId,
					Description: description,
					CreatedAt:   now.Add(-25 * time.Minute),
					UpdatedAt:   now,
				}

				m := model.EntityToGetTaskResponse(e)

				Expect(m).To(Equal(expected))
			})
		})

	})

	Describe("UpsertTaskRequest", func() {

		var m model.UpsertTaskRequest

		BeforeEach(func() {
			id := id
			m = model.UpsertTaskRequest{
				Id:          &id,
				Name:        name,
				StatusId:    statusId,
				Description: description,
			}
		})

		Describe("IsValid", func() {

			var inputId *int

			BeforeEach(func() {
				id := *m.Id + 1
				inputId = &id
			})

			When("the model id is nil", func() {

				BeforeEach(func() {
					m.Id = nil
				})

				When("the id argument is nil", func() {
					It("returns true", func() {
						Expect(m.IsValid(nil)).To(BeTrue())
					})
				})

				When("the id argument is non-nil", func() {
					It("returns false", func() {
						Expect(m.IsValid(inputId)).To(BeFalse())
					})
				})
			})

			When("the model id is non-nil", func() {
				When("the id argument is nil", func() {
					It("returns false", func() {
						Expect(m.IsValid(nil)).To(BeFalse())
					})
				})

				When("the id argument is non-nil and equal to the model's id", func() {
					It("returns true", func() {
						Expect(m.IsValid(m.Id)).To(BeTrue())
					})
				})

				When("the id argument is non-nil and not equal to the model's id", func() {
					It("returns false", func() {
						Expect(m.IsValid(inputId)).To(BeFalse())
					})
				})
			})

		})

		Describe("ToEntity", func() {

			var expected entity.Task

			BeforeEach(func() {
				expected = entity.Task{
					Id:          id,
					Name:        name,
					StatusId:    statusId,
					Description: description,
				}
			})

			When("the id field is nil", func() {
				It("returns an entity having a zero id", func() {
					m.Id = nil
					expected.Id = 0

					e := m.ToEntity()

					Expect(e).To(Equal(expected))
				})
			})

			When("the id field is non-nil", func() {
				It("returns an entity having a zero id", func() {
					e := m.ToEntity()

					Expect(e).To(Equal(expected))
				})
			})

		})

	})

})
