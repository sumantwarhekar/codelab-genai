package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/compute/metadata"
	"github.com/firebase/genkit/go/plugins/dotprompt"
	"github.com/firebase/genkit/go/plugins/vertexai"
)

func main() {
	ctx := context.Background()
	var projectId string
	var err error
	projectId = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectId == "" {
		projectId, err = metadata.ProjectIDWithContext(ctx)
		os.Setenv("GOOGLE_CLOUD_PROJECT", projectId)
		if err != nil {
			return
		}
	}

	if err := vertexai.Init(ctx, nil); err != nil {
		log.Fatal(err)
		return
	}

	dotprompt.SetDirectory("./")
	prompt, err := dotprompt.Open("animal-facts")
	if err != nil {
		log.Fatal(err)
		return
	}
	type AnimalPromptInput struct {
		Animal string `json:"animal"`
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		animal := r.URL.Query().Get("animal")
		if animal == "" {
			animal = "dog"
		}

		response, err := prompt.Generate(
			ctx,
			&dotprompt.PromptRequest{
				Variables: AnimalPromptInput{
					Animal: animal,
				},
			},
			nil,
		)

		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, err)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, response.Text())

	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
