package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/tasks/v1"
)

type GoogleAuth struct {
	Config *oauth2.Config
	Token  *oauth2.Token
}

func NewGoogleAuth() *GoogleAuth {
	secretBytes, err := os.ReadFile("config/client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(secretBytes, tasks.TasksScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	return &GoogleAuth{Config: config}
}

func (g *GoogleAuth) GetTokenFromWeb() *oauth2.Token {
	authUrl := g.Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Visit the following URL for the authentication: \n%v\n", authUrl)

	var code string
	log.Print("Enter the authorization code: ")
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	token, err := g.Config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token: %v", err)
	}

	return token
}

func (g *GoogleAuth) SaveToken(file string, token *oauth2.Token) {
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Failed to save token to file %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
