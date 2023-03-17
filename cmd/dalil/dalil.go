package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aeon-fruit/dalil.git/internal/pkg/common/constants"
	"github.com/aeon-fruit/dalil.git/internal/pkg/config"
	"github.com/aeon-fruit/dalil.git/internal/pkg/middleware"
	controller "github.com/aeon-fruit/dalil.git/internal/pkg/tasks/controller"
	dao "github.com/aeon-fruit/dalil.git/internal/pkg/tasks/dao"
	service "github.com/aeon-fruit/dalil.git/internal/pkg/tasks/service"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func main() {
	appConfig := config.New(config.WithEnvVars())

	addr := fmt.Sprintf(":%v", appConfig.AppPort)
	handler := getHandler(appConfig)
	if err := http.ListenAndServe(addr, handler); err != nil {
		fmt.Printf("%v", err)
	}
}

func getHandler(appConfig config.AppConfig) http.Handler {
	tasksDAO := dao.New()
	tasksService := service.New(service.WithRepository(tasksDAO))
	tasksCtrl := controller.New(controller.WithService(tasksService))

	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(60 * time.Second))

	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/api/", func(r chi.Router) {
		r.Route("/tasks", func(r chi.Router) {
			r.Get("/", tasksCtrl.GetAll)
			r.Post("/", tasksCtrl.Add)

			r.Route("/{id}", func(r chi.Router) {
				r.Use(middleware.PathParamContextInt(constants.Id))
				r.Get("/", tasksCtrl.GetById)
				r.Put("/", tasksCtrl.Update)
				r.Delete("/", tasksCtrl.RemoveById)
			})
		})
	})

	return r
}
