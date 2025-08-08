package main

import (
	"dirigeant/worker"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	api := &worker.Api{
		Worker: worker.NewWorker(),
	}
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", api.HandleListTasks)
		r.Get("/{id}", api.HandleGetTask)
		r.Post("/", api.HandleCreateTask)
		r.Delete("/{id}", api.HandleDeleteTask)
	})

	http.ListenAndServe(":8080", r)
}
