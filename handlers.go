package main

import (
	"bc-github-app/templates"
	"fmt"
	"html/template"
	"net/http"
)

type AuthDetails struct {
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	TokenType             string `json:"token_type"`
}

func handleCallbackRequest(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	result, err := exchangeCode(code)
	if err != nil {
		http.Error(w, "An error occurred wile exchanging code: "+err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFS(templates.FS, "callback.html")
	if err != nil {
		http.Error(w, "An error occurred "+err.Error(), http.StatusBadRequest)
		return
	}
	err = tmpl.Execute(w, result)
	if err != nil {
		http.Error(w, "An error occurred "+err.Error(), http.StatusBadRequest)
		return
	}
}

func handleSignIn(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(templates.FS, "index.html")
	if err != nil {
		http.Error(w, "An error occurred "+err.Error(), http.StatusBadRequest)
		return
	}
	err = tmpl.Execute(w, struct {
		ClientID string
	}{
		ClientID: env.ClientID,
	})
	if err != nil {
		http.Error(w, "An error occurred "+err.Error(), http.StatusBadRequest)
		return
	}
}

func handleWelcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome"))
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handling webhook")     // prints to the console / terminal
	w.Write([]byte("handling webhook")) // returns the response via HTTP
}
