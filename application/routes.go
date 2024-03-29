package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yaninyzwitty/CRUD-go/handler"
)

func loadRoutes() *chi.Mux {

	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))

	})

	router.Route("/orders", loadOrderRoutes)

	return router

}

func loadOrderRoutes(router chi.Router) {

	orderHandler := &handler.Order{}

	router.Get("/", orderHandler.List)
	router.Post("/{id}", orderHandler.Create)
	router.Get("/{id}", orderHandler.GetById)
	router.Put("/{id}", orderHandler.Update)
	router.Delete("/{id}", orderHandler.Delete)

}
