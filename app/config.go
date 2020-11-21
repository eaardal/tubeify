package app

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	youTubeApiKeyEnvName       = "YOUTUBE_API_KEY"
	spotifyClientIdEnvName     = "SPOTIFY_CLIENT_ID"
	spotifyClientSecretEnvName = "SPOTIFY_CLIENT_SECRET"
)

var youTubeApiKey string
var spotifyClientId string
var spotifyClientSecret string

func SetupEnvironment() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	if val, exists := os.LookupEnv(youTubeApiKeyEnvName); exists {
		youTubeApiKey = val
	} else {
		log.Fatalf("env %s is missing", youTubeApiKeyEnvName)
	}

	if val, exists := os.LookupEnv(spotifyClientIdEnvName); exists {
		spotifyClientId = val
	} else {
		log.Fatalf("env %s is missing", spotifyClientIdEnvName)
	}

	if val, exists := os.LookupEnv(spotifyClientSecretEnvName); exists {
		spotifyClientSecret = val
	} else {
		log.Fatalf("env %s is missing", spotifyClientSecretEnvName)
	}
}
