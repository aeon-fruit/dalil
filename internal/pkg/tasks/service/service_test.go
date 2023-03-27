package service_test

import (
	"fmt"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aeon-fruit/dalil.git/internal/pkg/tasks/dao/entity"
	"github.com/aeon-fruit/dalil.git/internal/pkg/tasks/model"
	"github.com/aeon-fruit/dalil.git/internal/pkg/tasks/service"
	daoMock "github.com/aeon-fruit/dalil.git/test/mocks/tasks/dao"
)

var _ = Describe("Service", func() {

	const (
		id              = 1
		taskName        = "A task"
		taskDescription = "A short description of the task"
	)

	var (
		customErr      error
		mockCtrl       *gomock.Controller
		mockRepository *daoMock.MockRepository
		tasksSvc       service.Service
	)

	BeforeEach(func() {
		customErr = fmt.Errorf("custom error")

		mockCtrl = gomock.NewController(GinkgoT())
		mockRepository = daoMock.NewMockRepository(mockCtrl)
		tasksSvc = service.New(service.WithRepository(mockRepository))
	})

	Describe("New", func() {
		It("returns a non-nil instance", func() {
			Expect(service.New()).NotTo(BeNil())
		})
	})

	Describe("WithRepository", func() {
		It("changes a non-nil instance", func() {
			mockRepository.EXPECT().GetAll().Return(nil, customErr)

			Expect(tasksSvc.GetAll()).Error().To(Equal(customErr))
		})
	})

	Describe("GetById", func() {

		When("an error happens while retrieving the dao", func() {
			It("returns an empty dto plus the error", func() {
				mockRepository.EXPECT().GetById(gomock.Any()).Return(entity.Task{}, customErr)

				Expect(tasksSvc.GetById(id)).Error().To(Equal(customErr))
			})
		})

		When("retrieving the doa is successful", func() {
			It("returns a dto and no error", func() {
				timestamp := time.UnixMilli(1679143523911)
				createdAt := timestamp
				updatedAt := timestamp.Add(2 * time.Hour)
				dao := entity.Task{
					Id:   id,
					Name: taskName,
					Status: entity.Status{
						Id:          0,
						Name:        "Status",
						Description: "Status description",
						CreatedAt:   timestamp.Add(-time.Hour),
						UpdatedAt:   timestamp.Add(-time.Hour),
					},
					StatusId:    0,
					Description: taskDescription,
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
				}
				dto := model.GetTaskResponse{
					Id:          id,
					Name:        taskName,
					StatusId:    0,
					Description: taskDescription,
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
				}
				mockRepository.EXPECT().GetById(gomock.Any()).Return(dao, nil)

				Expect(tasksSvc.GetById(id)).To(Equal(dto))
			})
		})

	})

	Describe("GetAll", func() {

		When("an error happens while retrieving the list of dao", func() {
			It("returns nil and the error", func() {
				mockRepository.EXPECT().GetAll().Return(nil, customErr)

				Expect(tasksSvc.GetAll()).Error().To(Equal(customErr))
			})
		})

		When("the list of dao is empty", func() {
			It("returns nil and no error", func() {
				mockRepository.EXPECT().GetAll().Return(nil, nil)

				Expect(tasksSvc.GetAll()).To(BeNil())
			})
		})

		When("the list of dao is not empty", func() {
			It("returns a list of dto and no error", func() {
				const (
					taskOneId          = 1
					taskOneName        = "A task"
					taskOneDescription = "A short description of the task"
					taskTwoId          = 2
					taskTwoName        = "Example"
					taskTwoDescription = "An example task"
				)

				timestamp := time.UnixMilli(1679143523911)
				taskOneCreatedAt := timestamp
				taskOneUpdatedAt := timestamp.Add(2 * time.Hour)
				taskTwoCreatedAt := timestamp.Add(24 * time.Hour)
				taskTwoUpdatedAt := timestamp.Add(25 * time.Hour)
				daoList := []entity.Task{
					{
						Id:          taskOneId,
						Name:        taskOneName,
						StatusId:    0,
						Description: taskOneDescription,
						CreatedAt:   taskOneCreatedAt,
						UpdatedAt:   taskOneUpdatedAt,
					},
					{
						Id:          taskTwoId,
						Name:        taskTwoName,
						StatusId:    0,
						Description: taskTwoDescription,
						CreatedAt:   taskTwoCreatedAt,
						UpdatedAt:   taskTwoUpdatedAt,
					},
				}
				dtoList := []model.GetTaskResponse{
					{
						Id:          taskOneId,
						Name:        taskOneName,
						StatusId:    0,
						Description: taskOneDescription,
						CreatedAt:   taskOneCreatedAt,
						UpdatedAt:   taskOneUpdatedAt,
					},
					{
						Id:          taskTwoId,
						Name:        taskTwoName,
						StatusId:    0,
						Description: taskTwoDescription,
						CreatedAt:   taskTwoCreatedAt,
						UpdatedAt:   taskTwoUpdatedAt,
					},
				}

				mockRepository.EXPECT().GetAll().Return(daoList, nil)

				Expect(tasksSvc.GetAll()).To(Equal(dtoList))
			})
		})

	})

	Describe("Upsert", func() {
		var dao entity.Task
		var inDto model.UpsertTaskRequest
		var outDto model.GetTaskResponse

		BeforeEach(func() {
			timestamp := time.UnixMilli(1679143523911)
			createdAt := timestamp
			updatedAt := timestamp.Add(2 * time.Hour)
			dao = entity.Task{
				Id:   id,
				Name: taskName,
				Status: entity.Status{
					Id:          0,
					Name:        "Status",
					Description: "Status description",
					CreatedAt:   timestamp.Add(-time.Hour),
					UpdatedAt:   timestamp.Add(-time.Hour),
				},
				StatusId:    0,
				Description: taskDescription,
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
			}
			inDto = model.UpsertTaskRequest{
				Id:          nil,
				Name:        taskName,
				StatusId:    0,
				Description: taskDescription,
			}
			outDto = model.GetTaskResponse{
				Id:          id,
				Name:        taskName,
				StatusId:    0,
				Description: taskDescription,
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
			}
		})

		When("the id in the dto is nil", func() {

			It("tries to insert and no update is called", func() {
				mockRepository.EXPECT().Insert(gomock.Any()).Times(1)
				mockRepository.EXPECT().Update(gomock.Any()).Times(0)

				tasksSvc.Upsert(inDto)
			})

			When("an error happens while inserting the new dao", func() {
				It("returns an empty dto and the error", func() {
					mockRepository.EXPECT().Insert(gomock.Any()).Return(entity.Task{}, customErr)

					Expect(tasksSvc.Upsert(inDto)).Error().To(Equal(customErr))
				})
			})

			When("inserting the new dao is successful", func() {
				It("returns the new dto with an id and no error", func() {
					mockRepository.EXPECT().Insert(gomock.Any()).Return(dao, nil)

					Expect(tasksSvc.Upsert(inDto)).To(Equal(outDto))
				})
			})
		})

		When("the id in the dto is not nil", func() {

			BeforeEach(func() {
				dtoId := id
				inDto.Id = &dtoId
			})

			It("tries to insert and no update is called", func() {
				mockRepository.EXPECT().Update(gomock.Any()).Times(1)
				mockRepository.EXPECT().Insert(gomock.Any()).Times(0)

				tasksSvc.Upsert(inDto)
			})

			When("an error happens while updating the new dao", func() {
				It("returns an empty dto and the error", func() {
					mockRepository.EXPECT().Update(gomock.Any()).Return(entity.Task{}, customErr)

					Expect(tasksSvc.Upsert(inDto)).Error().To(Equal(customErr))
				})
			})

			When("updating the dao is successful", func() {
				It("returns the updated dto and no error", func() {
					mockRepository.EXPECT().Update(gomock.Any()).Return(dao, nil)

					Expect(tasksSvc.Upsert(inDto)).To(Equal(outDto))
				})
			})
		})

	})

	Describe("RemoveById", func() {

		When("an error happens while removing the dao", func() {
			It("returns the error", func() {
				mockRepository.EXPECT().RemoveById(id).Return(entity.Task{}, customErr)

				Expect(tasksSvc.RemoveById(id)).To(Equal(customErr))
			})
		})

		When("removing the dao is successful", func() {
			It("returns no error", func() {
				mockRepository.EXPECT().RemoveById(id).Return(entity.Task{}, nil)

				Expect(tasksSvc.RemoveById(id)).ToNot(HaveOccurred())
			})
		})

	})

})
