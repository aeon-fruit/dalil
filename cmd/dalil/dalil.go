package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/aeon-fruit/dalil.git/internal/pkg/common/constants"
	"github.com/aeon-fruit/dalil.git/internal/pkg/config"
	"github.com/aeon-fruit/dalil.git/internal/pkg/log"
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

	logger := log.New(appConfig, os.Stderr)

	addr := fmt.Sprintf(":%v", appConfig.AppPort)
	handler := getHandler(logger)
	if err := http.ListenAndServe(addr, handler); err != nil {
		logger.Error(err, "Failed to start the server", "addr", addr)
	}
}

func getHandler(logger log.Logger) http.Handler {
	chiMiddleware.DefaultLogger = chiMiddleware.RequestLogger(&chiMiddleware.DefaultLogFormatter{
		Logger:  logger,
		NoColor: runtime.GOOS != "windows",
	})

	r := chi.NewRouter()

	r.Use(middleware.LoggingContext(logger))
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Timeout(60 * time.Second))

	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/api/", func(r chi.Router) {
		r.Route("/v1/", v1())
	})

	return r
}

func v1() func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/tasks", tasksRouter())
	}
}

func tasksRouter() func(r chi.Router) {
	tasksDAO := dao.New()
	tasksService := service.New(service.WithRepository(tasksDAO))
	tasksCtrl := controller.New(controller.WithService(tasksService))

	return func(r chi.Router) {
		r.Get("/", tasksCtrl.GetAll)
		r.Post("/", tasksCtrl.Add)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(middleware.PathParamContextInt(constants.Id))
			r.Get("/", tasksCtrl.GetById)
			r.Put("/", tasksCtrl.Update)
			r.Delete("/", tasksCtrl.RemoveById)
		})
	}
}
