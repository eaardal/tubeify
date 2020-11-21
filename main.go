package main

import (
	"flag"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
	"tubeify/app"
)

var youTubeVideoId = flag.String("youtube-video-id", "", "The ID of the YouTube video to scrape for songs")
var spotifyPlaylistId = flag.String("spotify-playlist-id", "", "The ID of the Spotify playlist to save songs to")

func main() {
	app.SetupEnvironment()
	parseFlags()
	runApp()
}

func parseFlags() {
	flag.Parse()

	if youTubeVideoId == nil || *youTubeVideoId == "" {
		log.Fatalf("youtube-video-id cli arg is missing")
	}

	if spotifyPlaylistId == nil || *spotifyPlaylistId == "" {
		log.Fatalf("spotify-playlist-id cli arg is missing")
	}
}

func runApp() {
	auth := app.CreateSpotifyAuthenticator()
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

		youTubeTracks, err := app.ScrapeYouTubeVideoDescriptionForTracks(*youTubeVideoId)
		if err != nil {
			log.Fatal(err)
		}

		if err := app.AddTracksToSpotifyPlaylist(client, youTubeTracks, *spotifyPlaylistId); err != nil {
			log.Fatal(err)
		}

		println("all done, press CTRL+C to exit")
	}
}
