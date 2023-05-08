package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	err := loadEnv()
	if err != nil {
		log.Fatal("error loading env", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", handleSignIn)
	r.Get("/callback", handleCallbackRequest)
	r.Get("/webhook", handleWebhook)

	http.ListenAndServe(":4000", r)
}
