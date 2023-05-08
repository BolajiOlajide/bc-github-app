package main

import (
	"os"

	"github.com/joho/godotenv"
)

var env struct {
	AppID        string
	PrivateKey   string
	ClientID     string
	ClientSecret string
}

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	env.AppID = os.Getenv("APP_ID")
	privateKeyPath := os.Getenv("PRIVATE_KEY_PATH")
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return err
	}
	env.PrivateKey = string(privateKey)
	env.ClientID = os.Getenv("CLIENT_ID")
	env.ClientSecret = os.Getenv("CLIENT_SECRET")
	return nil
}
