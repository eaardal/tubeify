package app

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
)

var spotifyOAuthRedirectUrl = "http://localhost:8080"
var spotifyOAuthState = uuid.Must(uuid.NewV4(), nil).String()

func CreateSpotifyAuthenticator() *spotify.Authenticator {
	auth := spotify.NewAuthenticator(spotifyOAuthRedirectUrl, spotify.ScopeUserReadPrivate, spotify.ScopePlaylistModifyPublic)
	auth.SetAuthInfo(spotifyClientId, spotifyClientSecret)
	return &auth
}

func CreateSpotifyOAuthUrl(auth *spotify.Authenticator) string {
	return auth.AuthURL(spotifyOAuthState)
}

func CreateSpotifyClient(auth *spotify.Authenticator, r *http.Request) (*spotify.Client, error) {
	token, err := auth.Token(spotifyOAuthState, r)
	if err != nil {
		return nil, err
	}

	client := auth.NewClient(token)
	return &client, nil
}

func PrintLoggedInSpotifyUserName(client *spotify.Client) {
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	logf("spotify", "logged in as %s", user.DisplayName)
}

func AddTracksToSpotifyPlaylist(client *spotify.Client, youTubeTracks []Track, playlistId string) error {
	spotifyTracks, err := searchSpotifyForTracks(client, youTubeTracks, TakeFirst)
	if err != nil {
		return err
	}

	if _, err := client.AddTracksToPlaylist(spotify.ID(playlistId), spotifyTracks...); err != nil {
		return err
	}

	logf("spotify", "added %d tracks to playlist", len(spotifyTracks))
	return nil
}

func searchSpotifyForTracks(client *spotify.Client, tracks []Track, searchFor SpotifySearch) ([]spotify.ID, error) {
	trackIds := make([]spotify.ID, 0)

	for _, track := range tracks {
		query := fmt.Sprintf("%s %s", track.MainArtist().Name, track.Track)
		search, err := client.Search(query, spotify.SearchTypeTrack)
		if err != nil {
			return nil, fmt.Errorf("spotify search failed: %v", err)
		}

		if search == nil || search.Tracks == nil || len(search.Tracks.Tracks) == 0 {
			logf("spotify", "search for '%s' returned empty", query)
			continue
		}

		switch searchFor {
		case TakeFirst:
			logf("spotify", "search for '%s': found %d tracks, taking first track only", query, search.Tracks.Total)
			first := search.Tracks.Tracks[0]
			trackIds = append(trackIds, first.ID)
			break
		case TakeAll:
			logf("spotify", "search for '%s': found %d tracks, taking all tracks found", query, search.Tracks.Total)
			for i, track := range search.Tracks.Tracks {
				logf("spotify", "%d: %s - %s", i+1, track.Artists[0].Name, track.Name)
				trackIds = append(trackIds, track.ID)
			}
			break
		default:
			return nil, fmt.Errorf("unknown searchFor: %s", searchFor)
		}
	}

	return trackIds, nil
}

type SpotifySearch string

const (
	TakeFirst SpotifySearch = "take_first"
	TakeAll   SpotifySearch = "take_all"
)
