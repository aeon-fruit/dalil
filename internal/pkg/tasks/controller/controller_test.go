package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aeon-fruit/dalil.git/internal/pkg/common/constants"
	"github.com/aeon-fruit/dalil.git/internal/pkg/common/errors"
	reqctx "github.com/aeon-fruit/dalil.git/internal/pkg/context/request"
	errorModel "github.com/aeon-fruit/dalil.git/internal/pkg/model/error"
	"github.com/aeon-fruit/dalil.git/internal/pkg/tasks/controller"
	"github.com/aeon-fruit/dalil.git/internal/pkg/tasks/model"
	serviceMock "github.com/aeon-fruit/dalil.git/test/mocks/tasks/service"
)

var _ = Describe("Controller", func() {

	const url = "http://url"

	var (
		recorder    *httptest.ResponseRecorder
		mockCtrl    *gomock.Controller
		mockService *serviceMock.MockService
		tasksCtrl   controller.Controller
	)

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
		mockCtrl = gomock.NewController(GinkgoT())
		mockService = serviceMock.NewMockService(mockCtrl)
		tasksCtrl = controller.New(controller.WithService(mockService))
	})

	Describe("New", func() {
		It("returns a non-nil instance", func() {
			Expect(controller.New()).NotTo(BeNil())
		})
	})

	Describe("WithService", func() {
		It("changes a non-nil instance", func() {
			customErr := fmt.Errorf("some random error")
			mockService.EXPECT().GetAll().Return(nil, customErr)

			tasksCtrl.GetAll(recorder, &http.Request{})

			Expect(recorder.Body.String()).To(ContainSubstring(customErr.Error()))
		})
	})

	Describe("GetById", func() {

		var request *http.Request

		BeforeEach(func() {
			request = httptest.NewRequest("", url, strings.NewReader("{}"))
		})

		When("the id is not found", func() {
			It("responds with status InternalServerError and an error response payload", func() {
				tasksCtrl.GetById(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				var payload errorModel.Response
				err := json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&payload)

				Expect(err).ToNot(HaveOccurred())
				Expect(payload).NotTo(BeZero())
				Expect(payload.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		Context("the id is found", func() {

			BeforeEach(func() {
				ctx := reqctx.SetPathParam(request.Context(), constants.Id, "1")
				request = request.WithContext(ctx)
			})

			When("the entity is not found", func() {
				It("responds with status NotFound and no payload", func() {
					mockService.EXPECT().GetById(gomock.Any()).Return(model.GetTaskResponse{}, errors.ErrNotFound)

					tasksCtrl.GetById(recorder, request)

					Expect(recorder.Code).To(Equal(http.StatusNotFound))
					Expect(recorder.Body.String()).To(BeEmpty())
				})
			})

			When("an error happens while retrieving the entity", func() {
				It("responds with status InternalServerError and an error response payload", func() {
					customErr := fmt.Errorf("custom error")
					mockService.EXPECT().GetById(gomock.Any()).Return(model.GetTaskResponse{}, customErr)

					tasksCtrl.GetById(recorder, request)

					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

					var payload errorModel.Response
					err := json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&payload)

					Expect(err).ToNot(HaveOccurred())
					Expect(payload).NotTo(BeZero())
					Expect(payload.Code).To(Equal(http.StatusInternalServerError))
					Expect(payload.Message).To(Equal(customErr.Error()))
				})
			})

			When("the entity is found", func() {
				It("responds with status OK and the entity in the payload", func() {
					timestamp := time.UnixMilli(1679143523911)
					entity := model.GetTaskResponse{
						Id:          1,
						Name:        "A task",
						StatusId:    0,
						Description: "A short description of the task",
						CreatedAt:   timestamp,
						UpdatedAt:   timestamp.Add(2 * time.Hour),
					}
					mockService.EXPECT().GetById(gomock.Any()).Return(entity, nil)

					tasksCtrl.GetById(recorder, request)

					Expect(recorder.Code).To(Equal(http.StatusOK))

					var payload model.GetTaskResponse
					err := json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&payload)

					Expect(err).ToNot(HaveOccurred())
					Expect(payload).To(Equal(entity))
				})
			})

		})

	})

	Describe("GetAll", func() {

		var request *http.Request

		BeforeEach(func() {
			request = httptest.NewRequest("", url, strings.NewReader("{}"))
		})

		When("an error happens while retrieving the list of entities", func() {
			It("responds with status InternalServerError and an error response payload", func() {
				customErr := fmt.Errorf("custom error")
				mockService.EXPECT().GetAll().Return(nil, customErr)

				tasksCtrl.GetAll(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				var payload errorModel.Response
				err := json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&payload)

				Expect(err).ToNot(HaveOccurred())
				Expect(payload).NotTo(BeZero())
				Expect(payload.Code).To(Equal(http.StatusInternalServerError))
				Expect(payload.Message).To(Equal(customErr.Error()))
			})
		})

		When("the list of entities is empty", func() {
			It("responds with status NoContent and no payload", func() {
				mockService.EXPECT().GetAll().Return(nil, nil)

				tasksCtrl.GetAll(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusNoContent))
				Expect(recorder.Body.String()).To(BeEmpty())
			})
		})

		When("the list of entities is not empty", func() {
			It("responds with status OK and the full list in the payload", func() {
				timestamp := time.UnixMilli(1679143523911)
				list := []model.GetTaskResponse{
					{
						Id:          1,
						Name:        "A task",
						StatusId:    0,
						Description: "A short description of the task",
						CreatedAt:   timestamp,
						UpdatedAt:   timestamp.Add(2 * time.Hour),
					},
					{
						Id:          2,
						Name:        "Another task",
						StatusId:    0,
						Description: "A short description of the second task",
						CreatedAt:   timestamp.Add(27 * time.Hour),
						UpdatedAt:   timestamp.Add(27 * time.Hour),
					},
				}
				mockService.EXPECT().GetAll().Return(list, nil)

				tasksCtrl.GetAll(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusOK))

				var payload []model.GetTaskResponse
				err := json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&payload)

				Expect(err).ToNot(HaveOccurred())
				Expect(payload).To(Equal(list))
			})
		})

	})

	Describe("Add", func() {
		var request *http.Request

		BeforeEach(func() {
			entity := model.UpsertTaskRequest{
				Id:          nil,
				Name:        "A new task",
				StatusId:    0,
				Description: "A new short description of the task",
			}

			body, err := json.Marshal(entity)

			Expect(err).ToNot(HaveOccurred())

			request = httptest.NewRequest("", url, bytes.NewReader(body))
		})

		When("the request payload format is wrong", func() {
			It("responds with status BadRequest and an error response payload", func() {
				request = httptest.NewRequest("", url, strings.NewReader(`{"id":"error"}`))

				tasksCtrl.Add(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusBadRequest))

				var payload errorModel.Response
				err := json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&payload)

				Expect(err).ToNot(HaveOccurred())
				Expect(payload).NotTo(BeZero())
				Expect(payload.Code).To(Equal(http.StatusBadRequest))
				Expect(payload.Message).To(ContainSubstring("json"))
			})
		})

		When("the fields of the request payload are empty", func() {
			It("responds with status BadRequest and an error response payload", func() {
				body, err := json.Marshal(model.UpsertTaskRequest{})

				Expect(err).ToNot(HaveOccurred())

				request = httptest.NewRequest("", url, bytes.NewReader(body))

				tasksCtrl.Add(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusBadRequest))

				var payload errorModel.Response
				err = json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&payload)

				Expect(err).ToNot(HaveOccurred())
				Expect(payload).NotTo(BeZero())
				Expect(payload.Code).To(Equal(http.StatusBadRequest))
				Expect(payload.Message).NotTo(BeEmpty())
			})
		})

		When("id is nil in the request payload", func() {
			It("responds with status BadRequest and an error response payload", func() {
				id := 1
				entity := model.UpsertTaskRequest{
					Id:          &id,
					Name:        "A task",
					StatusId:    0,
					Description: "A short description of the task",
				}

				body, err := json.Marshal(entity)

				Expect(err).ToNot(HaveOccurred())

				request = httptest.NewRequest("", url, bytes.NewReader(body))

				tasksCtrl.Add(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusBadRequest))

				var payload errorModel.Response
				err = json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&payload)

				Expect(err).ToNot(HaveOccurred())
				Expect(payload).NotTo(BeZero())
				Expect(payload.Code).To(Equal(http.StatusBadRequest))
				Expect(payload.Message).NotTo(BeEmpty())
			})
		})

		When("an error happens while adding the entity", func() {
			It("responds with status InternalServerError and an error response payload", func() {
				customErr := fmt.Errorf("custom error")
				mockService.EXPECT().Upsert(gomock.Any()).Return(model.GetTaskResponse{}, customErr)

				tasksCtrl.Add(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				var payload errorModel.Response
				err := json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&payload)

				Expect(err).ToNot(HaveOccurred())
				Expect(payload).NotTo(BeZero())
				Expect(payload.Code).To(Equal(http.StatusInternalServerError))
				Expect(payload.Message).To(Equal(customErr.Error()))
			})
		})

		When("the entity is found and updated", func() {
			It("responds with status OK and the updated entity in the payload", func() {
				timestamp := time.UnixMilli(1679143523911)
				entity := model.GetTaskResponse{
					Id:          1,
					Name:        "New task",
					StatusId:    0,
					Description: "New short description of the task",
					CreatedAt:   timestamp,
					UpdatedAt:   timestamp,
				}
				mockService.EXPECT().Upsert(gomock.Any()).Return(entity, nil)

				tasksCtrl.Add(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusCreated))
				Expect(recorder.Header().Get("Location")).To(Equal("url/1"))

				var payload model.GetTaskResponse
				err := json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&payload)

				Expect(err).ToNot(HaveOccurred())
				Expect(payload).To(Equal(entity))
			})
		})

	})

	Describe("Update", func() {

		var request *http.Request

		BeforeEach(func() {
			id := 1
			entity := model.UpsertTaskRequest{
				Id:          &id,
				Name:        "An updated task",
				StatusId:    0,
				Description: "An updated short description of the task",
			}

			body, err := json.Marshal(entity)

			Expect(err).ToNot(HaveOccurred())

			request = httptest.NewRequest("", url, bytes.NewReader(body))
		})

		When("the id is not found", func() {
			It("responds with status InternalServerError and an error response payload", func() {
				tasksCtrl.Update(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				var payload errorModel.Response
				err := json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&payload)

				Expect(err).ToNot(HaveOccurred())
				Expect(payload).NotTo(BeZero())
				Expect(payload.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		Context("the id is found", func() {

			BeforeEach(func() {
				ctx := reqctx.SetPathParam(request.Context(), constants.Id, "1")
				request = request.WithContext(ctx)
			})

			When("the request payload format is wrong", func() {
				It("responds with status BadRequest and an error response payload", func() {
					request = httptest.NewRequest("", url, strings.NewReader(`{"id":"error"}`))
					ctx := reqctx.SetPathParam(request.Context(), constants.Id, "1")
					request = request.WithContext(ctx)

					tasksCtrl.Update(recorder, request)

					Expect(recorder.Code).To(Equal(http.StatusBadRequest))

					var payload errorModel.Response
					err := json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&payload)

					Expect(err).ToNot(HaveOccurred())
					Expect(payload).NotTo(BeZero())
					Expect(payload.Code).To(Equal(http.StatusBadRequest))
					Expect(payload.Message).To(ContainSubstring("json"))
				})
			})

			When("the fields of the request payload are empty", func() {
				It("responds with status BadRequest and an error response payload", func() {
					body, err := json.Marshal(model.UpsertTaskRequest{})

					Expect(err).ToNot(HaveOccurred())

					request = httptest.NewRequest("", url, bytes.NewReader(body))
					ctx := reqctx.SetPathParam(request.Context(), constants.Id, "1")
					request = request.WithContext(ctx)

					tasksCtrl.Update(recorder, request)

					Expect(recorder.Code).To(Equal(http.StatusBadRequest))

					var payload errorModel.Response
					err = json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&payload)

					Expect(err).ToNot(HaveOccurred())
					Expect(payload).NotTo(BeZero())
					Expect(payload.Code).To(Equal(http.StatusBadRequest))
					Expect(payload.Message).NotTo(BeEmpty())
				})
			})

			When("id is nil in the request payload", func() {
				It("responds with status BadRequest and an error response payload", func() {
					entity := model.UpsertTaskRequest{
						Id:          nil,
						Name:        "A task update",
						StatusId:    0,
						Description: "Short description of the task",
					}

					body, err := json.Marshal(entity)

					Expect(err).ToNot(HaveOccurred())

					request = httptest.NewRequest("", url, bytes.NewReader(body))
					ctx := reqctx.SetPathParam(request.Context(), constants.Id, "1")
					request = request.WithContext(ctx)

					tasksCtrl.Update(recorder, request)

					Expect(recorder.Code).To(Equal(http.StatusBadRequest))

					var payload errorModel.Response
					err = json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&payload)

					Expect(err).ToNot(HaveOccurred())
					Expect(payload).NotTo(BeZero())
					Expect(payload.Code).To(Equal(http.StatusBadRequest))
					Expect(payload.Message).NotTo(BeEmpty())
				})
			})

			When("the entity is not found", func() {
				It("responds with status NotFound and no payload", func() {
					mockService.EXPECT().Upsert(gomock.Any()).Return(model.GetTaskResponse{}, errors.ErrNotFound)

					tasksCtrl.Update(recorder, request)

					Expect(recorder.Code).To(Equal(http.StatusNotFound))
					Expect(recorder.Body.String()).To(BeEmpty())
				})
			})

			When("the entity is found but not modified", func() {
				It("responds with status NotModified and no payload", func() {
					mockService.EXPECT().Upsert(gomock.Any()).Return(model.GetTaskResponse{}, errors.ErrNotModified)

					tasksCtrl.Update(recorder, request)

					Expect(recorder.Code).To(Equal(http.StatusNotModified))
					Expect(recorder.Body.String()).To(BeEmpty())
				})
			})

			When("an error happens while retrieving the entity", func() {
				It("responds with status InternalServerError and an error response payload", func() {
					customErr := fmt.Errorf("custom error")
					mockService.EXPECT().Upsert(gomock.Any()).Return(model.GetTaskResponse{}, customErr)

					tasksCtrl.Update(recorder, request)

					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

					var payload errorModel.Response
					err := json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&payload)

					Expect(err).ToNot(HaveOccurred())
					Expect(payload).NotTo(BeZero())
					Expect(payload.Code).To(Equal(http.StatusInternalServerError))
					Expect(payload.Message).To(Equal(customErr.Error()))
				})
			})

			When("the entity is found and updated", func() {
				It("responds with status OK and the updated entity in the payload", func() {
					timestamp := time.UnixMilli(1679143523911)
					entity := model.GetTaskResponse{
						Id:          1,
						Name:        "A task update",
						StatusId:    0,
						Description: "Short description of the task",
						CreatedAt:   timestamp,
						UpdatedAt:   timestamp.Add(2 * time.Hour),
					}
					mockService.EXPECT().Upsert(gomock.Any()).Return(entity, nil)

					tasksCtrl.Update(recorder, request)

					Expect(recorder.Code).To(Equal(http.StatusOK))

					var payload model.GetTaskResponse
					err := json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&payload)

					Expect(err).ToNot(HaveOccurred())
					Expect(payload).To(Equal(entity))
				})
			})

		})

	})

	Describe("RemoveById", func() {

		var request *http.Request

		BeforeEach(func() {
			request = httptest.NewRequest("", url, strings.NewReader("{}"))
		})

		When("the id is not found", func() {
			It("responds with status InternalServerError and an error response payload", func() {
				tasksCtrl.RemoveById(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				var payload errorModel.Response
				err := json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&payload)

				Expect(err).ToNot(HaveOccurred())
				Expect(payload).NotTo(BeZero())
				Expect(payload.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		Context("the id is found", func() {

			BeforeEach(func() {
				ctx := reqctx.SetPathParam(request.Context(), constants.Id, "1")
				request = request.WithContext(ctx)
			})

			When("the entity is not found", func() {
				It("responds with status NotFound and no payload", func() {
					mockService.EXPECT().RemoveById(gomock.Any()).Return(errors.ErrNotFound)

					tasksCtrl.RemoveById(recorder, request)

					Expect(recorder.Code).To(Equal(http.StatusNotFound))
					Expect(recorder.Body.String()).To(BeEmpty())
				})
			})

			When("an error happens while removing", func() {
				It("responds with status InternalServerError and an error response payload", func() {
					customErr := fmt.Errorf("custom error")
					mockService.EXPECT().RemoveById(gomock.Any()).Return(customErr)

					tasksCtrl.RemoveById(recorder, request)

					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

					var payload errorModel.Response
					err := json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&payload)

					Expect(err).ToNot(HaveOccurred())
					Expect(payload).NotTo(BeZero())
					Expect(payload.Code).To(Equal(http.StatusInternalServerError))
					Expect(payload.Message).To(Equal(customErr.Error()))
				})
			})

			When("the entity is found", func() {
				It("responds with status NoContent and no payload", func() {
					mockService.EXPECT().RemoveById(gomock.Any()).Return(nil)

					tasksCtrl.RemoveById(recorder, request)

					Expect(recorder.Code).To(Equal(http.StatusNoContent))
					Expect(recorder.Body.String()).To(BeEmpty())
				})
			})

		})
	})

	Describe("RemoveByIds", func() {
		// Unimplemented
	})

})
