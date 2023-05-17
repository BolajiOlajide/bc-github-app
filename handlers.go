package main

import (
	"bc-github-app/templates"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/google/go-github/v52/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
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
	// restClient := github.NewTokenClient(ctx, token)

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	fmt.Printf("I created the httpclient!")

	graphQLClient := githubv4.NewClient(httpClient)

	fmt.Printf("I created the graphqlclient!")

	// path := "file.md"
	// content := []byte("# Hola\n* This is a file committed via Github's API\n")
	// message := "hopefully signed commit"

	// _, err := createBranch(ctx, restClient, body.Branch, body.Repo, body.Owner)
	// if err != nil {
	// 	log.Fatalf("error creating branch %s", err.Error())
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// _, _, err = restClient.Repositories.CreateFile(ctx, body.Owner, body.Repo, path, &github.RepositoryContentFileOptions{
	// 	Message: github.String(message),
	// 	Content: content,
	// 	Branch:  github.String(body.Branch),
	// })

	err := createCommitOnBranch(ctx, graphQLClient, "becca-test", "This is a test!")
	if err != nil {
		fmt.Printf("There was an error!: %v", err)

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("success! yaaaay!"))
}

func createCommitOnBranch(ctx context.Context, client *githubv4.Client, branch, message string) (err error) {
	var mutation struct {
		CreateCommitOnBranch struct {
			Commit struct {
				ID githubv4.ID
				// Signature struct {
				// 	Email   string
				// 	IsValid bool
				// }
			}
		} `graphql:"createCommitOnBranch(input: $input)"`
	}
	input := githubv4.CreateCommitOnBranchInput{
		Branch: githubv4.CommittableBranch{
			RepositoryNameWithOwner: githubv4.NewString(githubv4.String("st0nebraker/food-diary")),
			BranchName:              githubv4.NewString(githubv4.String(branch)),
		},
		Message: githubv4.CommitMessage{
			Headline: *githubv4.NewString(githubv4.String(message)),
		},
		FileChanges: &githubv4.FileChanges{
			Additions: &[]githubv4.FileAddition{
				{
					Path:     "docs/README.txt",
					Contents: githubv4.Base64String(base64.StdEncoding.EncodeToString([]byte("Hello world!\n"))),
				},
			},
		},
		ExpectedHeadOid: githubv4.GitObjectID("d849462e17cda8ad9cc07ef38f3160f2c768154e"),
	}

	err = client.Mutate(context.Background(), &mutation, input, nil)
	if err != nil {
		return err
		// Handle error.
	}
	// fmt.Printf("Commit was signed by %v, right? %v!\n", mutation.CreateCommitOnBranch.Commit.Signature.Email, mutation.CreateCommitOnBranch.Commit.Signature.IsValid)

	return nil
}

func createBranch(ctx context.Context, restClient *github.Client, branch, repo, owner string) (*github.Reference, error) {
	// Fetching the latest commit on the master branch
	ref, _, err := restClient.Git.GetRef(ctx, owner, repo, "refs/heads/master")
	if err != nil {
		return nil, err
	}

	// Creating a new branch
	newRef, _, err := restClient.Git.CreateRef(ctx, owner, repo, &github.Reference{
		Ref:    github.String(fmt.Sprintf("refs/heads/%s", branch)),
		Object: &github.GitObject{SHA: ref.Object.SHA},
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return newRef, nil
}
