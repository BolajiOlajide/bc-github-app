package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	appId := os.Getenv("APP_ID")
	privateKeyPath := os.Getenv("PRIVATE_KEY_PATH")
	privateKey, err := os.ReadFile(privateKeyPath)

	if err != nil {
		log.Fatal("Error reading private key: ", err)
	}

	fmt.Println(appId, string(privateKey))

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Get("/callback", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handling callback")
		w.Write([]byte("handling callback"))
	})
	r.Get("/webhook", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handling webhook")     // prints to the console / terminal
		w.Write([]byte("handling webhook")) // returns the response via HTTP
	})
	r.Get("/eventhandler", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handling event handler")
		w.Write([]byte("handling event"))
	})
	r.Get("/signin", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("signin-page").Parse(indexHTMLTemplate)
		if err != nil {
			http.Error(w, "An error occurred "+err.Error(), http.StatusBadRequest)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "An error occurred "+err.Error(), http.StatusBadRequest)
			return
		}
	})
	http.ListenAndServe(":3000", r)
}

const indexHTMLTemplate = `
<div>
	<button>Sign in via Github App</button>
</div>
`
