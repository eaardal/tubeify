package main

import (
	"github.com/zmb3/spotify"
	"log"
	"net/http"
	"tubeify/app"
)

var config = app.SetupEnvironment()

func main() {
	runApp()
}

func runApp() {
	auth := app.CreateSpotifyAuthenticator(config.SpotifyClientId, config.SpotifyClientSecret)
	url := app.CreateSpotifyOAuthUrl(auth)
	println("open url in browser to authenticate (and return to this terminal immediately):")
	println(url)

	http.HandleFunc("/", spotifyOAuthCallback(auth))

	log.Printf("starting http server at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func spotifyOAuthCallback(auth *spotify.Authenticator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client, err := app.CreateSpotifyClient(auth, r)
		if err != nil {
			log.Fatal(err)
		}

		app.PrintLoggedInSpotifyUserName(client)

		youTubeVideoIds := app.FindYouTubeVideoIds(config.YouTubeVideoId, config.YouTubeVideoUrl)
		youTubeTracks, err := app.ScrapeYouTubeVideoDescriptionForTracks(config.YouTubeApiKey, youTubeVideoIds)
		if err != nil {
			log.Fatal(err)
		}

		if err := app.AddTracksToSpotifyPlaylist(client, youTubeTracks, config.SpotifyPlaylistId); err != nil {
			log.Fatal(err)
		}

		println("all done, press CTRL+C to exit")
	}
}
