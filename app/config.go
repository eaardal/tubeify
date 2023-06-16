package app

import (
	"flag"
	"log"
	"os"
)

var youTubeVideoUrlFlag = flag.String("youtube-video-url", "", "YouTube video URL(s) to scrape for songs")
var youTubeVideoIdFlag = flag.String("youtube-video-id", "", "The ID of the YouTube video to scrape for songs")
var spotifyPlaylistIdFlag = flag.String("spotify-playlist-id", "", "The ID of the Spotify playlist to save songs to")

const (
	youTubeApiKeyEnvName       = "YOUTUBE_API_KEY"
	spotifyClientIdEnvName     = "SPOTIFY_CLIENT_ID"
	spotifyClientSecretEnvName = "SPOTIFY_CLIENT_SECRET"
	spotifyPlaylistIdEnvName   = "SPOTIFY_PLAYLIST_ID"
)

type Config struct {
	YouTubeApiKey       string
	YouTubeVideoId      string
	YouTubeVideoUrl     string
	SpotifyClientId     string
	SpotifyClientSecret string
	SpotifyPlaylistId   string
}

func SetupEnvironment() *Config {
	flag.Parse()

	config := Config{}

	if val, exists := os.LookupEnv(youTubeApiKeyEnvName); exists {
		config.YouTubeApiKey = val
	} else {
		log.Fatalf("env %s is missing", youTubeApiKeyEnvName)
	}

	if val, exists := os.LookupEnv(spotifyClientIdEnvName); exists {
		config.SpotifyClientId = val
	} else {
		log.Fatalf("env %s is missing", spotifyClientIdEnvName)
	}

	if val, exists := os.LookupEnv(spotifyClientSecretEnvName); exists {
		config.SpotifyClientSecret = val
	} else {
		log.Fatalf("env %s is missing", spotifyClientSecretEnvName)
	}

	if spotifyPlaylistIdFlag != nil && *spotifyPlaylistIdFlag != "" {
		config.SpotifyPlaylistId = *spotifyPlaylistIdFlag
	} else if val, exists := os.LookupEnv(spotifyPlaylistIdEnvName); exists {
		config.SpotifyPlaylistId = val
	} else {
		log.Fatalf("env %s or flag --spotify-playlist-id must be set", spotifyClientIdEnvName)
	}

	if youTubeVideoIdFlag != nil && *youTubeVideoUrlFlag != "" {
		config.YouTubeVideoId = *youTubeVideoIdFlag
	}

	if youTubeVideoUrlFlag != nil && *youTubeVideoUrlFlag != "" {
		config.YouTubeVideoUrl = *youTubeVideoUrlFlag
	}

	return &config
}
