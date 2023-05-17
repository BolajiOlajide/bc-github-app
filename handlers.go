package main

import (
	"bc-github-app/templates"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/google/go-github/v52/github"
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

type CreateCommitBody struct {
	Branch string
	Repo   string
	Owner  string
}

func handleCreateCommit(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("gh-token")

	dec := json.NewDecoder(r.Body)
	var body CreateCommitBody
	dec.Decode(&body)

	if body.Repo == "" {
		http.Error(w, "repo cannot be empty", http.StatusBadRequest)
	}

	if body.Branch == "" {
		http.Error(w, "branch cannot be empty", http.StatusBadRequest)
	}

	if body.Owner == "" {
		http.Error(w, "owner cannot be empty", http.StatusBadRequest)
	}

	ctx := context.Background()
	client := github.NewTokenClient(ctx, token)

	path := "file.md"
	content := []byte("# Hola\n* This is a file committed via Github's API\n")
	message := "hopefully signed commit"

	_, err := createBranch(ctx, client, body.Branch, body.Repo, body.Owner)
	if err != nil {
		log.Fatalf("error creating branch %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _, err = client.Repositories.CreateFile(ctx, body.Owner, body.Repo, path, &github.RepositoryContentFileOptions{
		Message: github.String(message),
		Content: content,
		Branch:  github.String(body.Branch),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("success! yaaaay!"))
}

func createBranch(ctx context.Context, client *github.Client, branch, repo, owner string) (*github.Reference, error) {
	// Fetching the latest commit on the master branch
	ref, _, err := client.Git.GetRef(ctx, owner, repo, "refs/heads/main")
	if err != nil {
		return nil, err
	}

	// Creating a new branch
	newRef, _, err := client.Git.CreateRef(ctx, owner, repo, &github.Reference{
		Ref:    github.String(fmt.Sprintf("refs/heads/%s", branch)),
		Object: &github.GitObject{SHA: ref.Object.SHA},
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return newRef, nil
}
