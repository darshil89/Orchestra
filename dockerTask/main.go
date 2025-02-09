package main

import (
	"fmt"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Docker Server is running"))
	})
	fmt.Printf("Server started at http://localhost:8082\n")
	http.ListenAndServe(":8082", r)
}
