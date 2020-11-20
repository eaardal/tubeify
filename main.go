package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/zmb3/spotify"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

//const apiKey = "AIzaSyBsSEblCxpRaJTPTDjUibqcyLvLq2SOX3o"
//const spotifyClientId = "57fd5f017fa447c1a5299e3a9c6f1c1a"
//const spotifyClientSecret = "b69619806d3347aa99a6dd44cfd89f4f"

const (
	youTubeApiKeyEnvName = "YOUTUBE_API_KEY"
	spotifyClientIdEnvName = "SPOTIFY_CLIENT_ID"
	spotifyClientSecretEnvName = "SPOTIFY_CLIENT_SECRET"
)

var youTubeApiKey string
var spotifyClientId string
var spotifyClientSecret string

var matchLineBeginningWithTimestamp = `^(\d{0,2}:\d{0,2})\s`
var spotifyOAuthRedirectUrl = "http://localhost:8080"
var spotifyOAuthState = uuid.Must(uuid.NewV4(), nil).String()

// _EopQtSVZQY
var youTubeVideoId = flag.String("youtube-video-id", "", "The ID of the YouTube video to scrape for songs")
// 38q3HMXsxSpXMvR0cqCQR6
var spotifyPlaylistId = flag.String("spotify-playlist-id", "", "The ID of the Spotify playlist to save songs to")

func main() {
	verifyEnvironment()
	setupSpotifyClient()
}

func verifyEnvironment() {
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

	flag.Parse()

	if youTubeVideoId == nil || *youTubeVideoId == "" {
		log.Fatalf("youtube-video-id cli arg is missing")
	}

	if spotifyPlaylistId == nil || *spotifyPlaylistId == "" {
		log.Fatalf("spotify-playlist-id cli arg is missing")
	}

	log.Printf("%+v", os.Environ())
}

func setupSpotifyClient() {
	auth := spotify.NewAuthenticator(spotifyOAuthRedirectUrl, spotify.ScopeUserReadPrivate, spotify.ScopePlaylistModifyPublic)
	auth.SetAuthInfo(spotifyClientId, spotifyClientSecret)

	url := auth.AuthURL(spotifyOAuthState)
	println("open url in browser to authenticate:")
	println(url)

	http.HandleFunc("/", spotifyOAuthCallback(auth))

	log.Printf("starting http server at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func spotifyOAuthCallback(auth spotify.Authenticator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client, err := createSpotifyClient(auth, r)
		if err != nil {
			log.Fatal(err)
		}

		printLoggedInSpotifyUserName(client)

		youTubeTracks, err := scrapeYouTubeVideoDescriptionForTracks(*youTubeVideoId)
		if err != nil {
			log.Fatal(err)
		}

		if err := addTracksToSpotifyPlaylist(client, youTubeTracks, *spotifyPlaylistId); err != nil {
			log.Fatal(err)
		}
	}
}

func scrapeYouTubeVideoDescriptionForTracks(videoId string) ([]Track, error) {
	yt, err := youtube.NewService(context.Background(), option.WithAPIKey(youTubeApiKey))
	if err != nil {
		log.Fatal(err)
	}

	res, err := yt.Videos.
		List([]string{"snippet"}).
		Id(videoId).
		Do()

	if err != nil {
		log.Fatal(err)
	}

	if res.HTTPStatusCode != http.StatusOK {
		log.Fatalf("youtube response statuscode %d", res.HTTPStatusCode)
	}

	if len(res.Items) != 1 {
		log.Fatalf("expected to find 1 YouTube video but got %d videos", len(res.Items))
	}

	youTubeVideo := res.Items[0]
	return findTracksInVideoDescription(youTubeVideo)
}

func findTracksInVideoDescription(video *youtube.Video) (tracks []Track, err error) {
	description := video.Snippet.Description
	descriptionLines := strings.Split(description, "\n")

	for _, line := range descriptionLines {
		isMatch, artist, track := interpretVideoDescriptionLine(line)
		if isMatch {
			track := Track{
				Artist:     artist,
				Track:      track,
				VideoId:    video.Id,
				VideoTitle: video.Snippet.Title,
			}
			tracks = append(tracks, track)
		}
	}

	return tracks, nil
}

func interpretVideoDescriptionLine(line string) (isMatch bool, artist string, track string) {
	if isMatch, _ := regexp.MatchString(matchLineBeginningWithTimestamp, line); isMatch {
		lineWithoutTimestamp := strings.TrimSpace(line[5:])
		lineParts := strings.Split(lineWithoutTimestamp, "-")
		artist := strings.TrimSpace(lineParts[0])
		track := strings.TrimSpace(lineParts[1])
		return true, artist, track
	}
	return false, "", ""
}

func createSpotifyClient(auth spotify.Authenticator, r *http.Request) (*spotify.Client, error) {
	token, err := auth.Token(spotifyOAuthState, r)
	if err != nil {
		return nil, err
	}

	client := auth.NewClient(token)
	return &client, nil
}

func printLoggedInSpotifyUserName(client *spotify.Client) {
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("logged in as %s", user.DisplayName)
}

func addTracksToSpotifyPlaylist(client *spotify.Client, youTubeTracks []Track, playlistId string) error {
	spotifyTracks, err := searchSpotifyForTracks(client, youTubeTracks)
	if err != nil {
		return err
	}

	if _, err := client.AddTracksToPlaylist(spotify.ID(playlistId), spotifyTracks...); err != nil {
		return err
	}

	log.Printf("added %d tracks to playlist", len(spotifyTracks))
	return nil
}

func searchSpotifyForTracks(client *spotify.Client, tracks []Track) ([]spotify.ID, error) {
	trackIds := make([]spotify.ID, 0)

	for _, track := range tracks {
		query := fmt.Sprintf("%s %s", track.Artist, track.Track)
		search, err := client.Search(query, spotify.SearchTypeTrack)
		if err != nil {
			return nil, err
		}

		if search == nil || search.Tracks == nil {
			continue
		}

		log.Printf("search for '%s': found %d tracks", query, search.Tracks.Total)

		for i, track := range search.Tracks.Tracks {
			log.Printf("%d: %s - %s", i +1, track.Artists[0].Name, track.Name)
			trackIds = append(trackIds, track.ID)
		}
	}

	return trackIds, nil
}

type Track struct {
	Artist string
	Track string
	VideoId string
	VideoTitle string
}