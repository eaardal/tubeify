package app

import (
	"context"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"log"
	"net/http"
	httpurl "net/url"
	"regexp"
	"strings"
)

const matchLineBeginningWithTimestamp = `^(\d{0,2}:\d{0,2})\s`

func ScrapeYouTubeVideoDescriptionForTracks(apiKey string, videoIds []string) (tracks []Track, err error) {
	yt, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	for _, videoId := range videoIds {
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

		t, err := findTracksInVideoDescription(youTubeVideo)
		if err != nil {
			log.Fatal(err)
		}

		tracks = append(tracks, t...)
	}

	return tracks, nil
}

func findTracksInVideoDescription(video *youtube.Video) (tracks []Track, err error) {
	description := video.Snippet.Description
	descriptionLines := strings.Split(description, "\n")

	for _, line := range descriptionLines {
		isMatch, artists, track := InterpretVideoDescriptionLine(line)
		if isMatch {
			track := Track{
				Artists:    artists,
				Track:      track,
				VideoId:    video.Id,
				VideoTitle: video.Snippet.Title,
			}
			logf("youtube", "track '%s' parsed to track name '%s' and artists '%s'", line, track.Track, track.ArtistsString())
			tracks = append(tracks, track)
		}
	}

	return tracks, nil
}

func extractTrack(text string) string {
	// Remove parens from text such as "(feat. xx)"
	cleaned := strings.ReplaceAll(text, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")

	// Normalize text by lowercasing it
	cleaned = strings.ToLower(cleaned)

	if strings.Contains(cleaned, "feat.") {
		parts := strings.Split(cleaned, "feat.")
		if len(parts) == 2 {
			return strings.TrimSpace(parts[0])
		} else {
			logf("youtube", "unexpectedly got %d parts when looking for track", len(parts))
		}
	} else {
		return cleaned
	}
	return ""
}

func extractFeaturedArtistsFromText(text string) (artists []*TrackArtist) {
	if !strings.Contains(text, "feat.") {
		return artists
	}

	parts := strings.Split(text, "feat.")

	if len(parts) == 2 {
		featuredArtist := strings.TrimSpace(parts[1])
		startIndex := 0

		if strings.Contains(featuredArtist, "[") {
			// Remove things in square brackets like
			for i, sym := range featuredArtist {
				s := string(sym)
				if s == "[" {
					startIndex = i
				}
			}
			featuredArtist = strings.TrimSpace(featuredArtist[0:startIndex])
		}

		artists = addArtistIfNotExists(artists, NewTrackArtist(featuredArtist, false))
	} else {
		logf("youtube", "unexpectedly got %d parts when looking for featured artists", len(parts))
	}

	return artists
}

func extractMainArtistsFromText(text string) []*TrackArtist {
	artists := make([]*TrackArtist, 0)

	if strings.Contains(text, "feat.") {
		parts := strings.Split(text, "feat.")
		if len(parts) == 2 {
			mainArtist := strings.TrimSpace(parts[0])
			artists = addArtistIfNotExists(artists, NewTrackArtist(mainArtist, true))
		} else {
			logf("youtube", "unexpectedly got %d parts when looking for main artist", len(parts))
		}
	}

	if strings.Contains(text, ",") {
		parts := strings.Split(text, ",")
		for i, part := range parts {
			artist := strings.TrimSpace(part)
			isMainArtist := i == 0
			artists = addArtistIfNotExists(artists, NewTrackArtist(artist, isMainArtist))
		}
	}

	if !strings.Contains(text, "feat.") && !strings.Contains(text, ",") {
		artistName := strings.TrimSpace(text)
		artists = addArtistIfNotExists(artists, NewTrackArtist(artistName, true))
	}

	return artists
}

func addArtistIfNotExists(artists []*TrackArtist, artist *TrackArtist) []*TrackArtist {
	if artists == nil {
		artists = make([]*TrackArtist, 0)
	}
	foundMatch := false
	for _, a := range artists {
		if a.Name == artist.Name {
			foundMatch = true
			break
		}
	}
	if !foundMatch {
		artists = append(artists, artist)
	}
	return artists
}

func extractArtists(text string, lookForMainArtist bool) (artists []*TrackArtist) {
	// Remove parens from text such as "(feat. xx)"
	cleaned := strings.ReplaceAll(text, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")

	// Normalize text by lowercasing it
	cleaned = strings.ToLower(cleaned)

	// Replace ampersands with comma to make it easier to split text only on comma to find all artists
	cleaned = strings.ReplaceAll(cleaned, "&", ",")

	// TODO: Find a better solution...
	if lookForMainArtist {
		artists = append(artists, extractMainArtistsFromText(cleaned)...)
	}

	artists = append(artists, extractFeaturedArtistsFromText(cleaned)...)

	return artists
}

func InterpretVideoDescriptionLine(line string) (isMatch bool, artists []*TrackArtist, track string) {
	if isMatch, _ := regexp.MatchString(matchLineBeginningWithTimestamp, line); isMatch {
		lineWithoutTimestamp := strings.TrimSpace(line[5:])
		lineParts := strings.Split(lineWithoutTimestamp, "-")
		artistPart := strings.TrimSpace(lineParts[0])
		trackPart := strings.TrimSpace(lineParts[1])

		artists = append(artists, extractArtists(artistPart, true)...)
		artists = append(artists, extractArtists(trackPart, false)...)
		track = extractTrack(trackPart)

		return true, artists, track
	}
	return false, nil, ""
}

type Track struct {
	Artists    []*TrackArtist
	Track      string
	VideoId    string
	VideoTitle string
}

func NewTrack(track string, artists ...*TrackArtist) *Track {
	return &Track{
		Artists:    artists,
		Track:      track,
		VideoId:    "",
		VideoTitle: "",
	}
}

func (t Track) MainArtist() *TrackArtist {
	for _, artist := range t.Artists {
		if artist.IsMainTrackArtist {
			return artist
		}
	}
	return nil
}

func (t Track) ArtistsString() string {
	str := ""
	for i, artist := range t.Artists {
		if i+1 == len(t.Artists) {
			str += artist.Name
		} else {
			str += fmt.Sprintf("%s, ", artist.Name)
		}
	}
	return str
}

type TrackArtist struct {
	Name              string
	IsMainTrackArtist bool
}

func NewTrackArtist(name string, isMainTrackArtist bool) *TrackArtist {
	return &TrackArtist{
		Name:              name,
		IsMainTrackArtist: isMainTrackArtist,
	}
}

func FindYouTubeVideoIds(videoIdCLIArg string, videoUrlCLIArg string) []string {
	ids := make([]string, 0)
	ids = append(ids, parseVideoIdCLIArg(videoIdCLIArg)...)
	ids = append(ids, parseVideoUrlCLIArg(videoUrlCLIArg)...)
	return ids
}

func parseVideoIdCLIArg(videoIdCLIArg string) (ids []string) {
	if videoIdCLIArg != "" {
		if strings.Contains(videoIdCLIArg, ",") {
			parts := strings.Split(videoIdCLIArg, ",")
			for _, part := range parts {
				ids = append(ids, strings.TrimSpace(part))
			}
		} else {
			ids = append(ids, videoIdCLIArg)
		}
	}
	return ids
}

func parseVideoUrlCLIArg(videoUrlCLIArg string) (ids []string) {
	if videoUrlCLIArg != "" {
		if strings.Contains(videoUrlCLIArg, ",") {
			parts := strings.Split(videoUrlCLIArg, ",")
			for _, part := range parts {
				ids = append(ids, findVideoIdInUrl(part))
			}
		} else {
			ids = append(ids, findVideoIdInUrl(videoUrlCLIArg))
		}
	}
	return ids
}

func findVideoIdInUrl(url string) string {
	trimmedUrl := strings.TrimSpace(url)
	uri, err := httpurl.Parse(trimmedUrl)
	if err != nil {
		log.Fatalf("failed to parse url '%s': %s", trimmedUrl, err)
	}
	return uri.Query().Get("v")
}
