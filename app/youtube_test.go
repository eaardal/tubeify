package app_test

import (
	"fmt"
	"log"
	"testing"
	"tubeify/app"
)

/*
```
0:00 Seven Lions, Last Heroes & HALIENE - Don't Wanna Fall
3:22 Nurko feat. RØRY - Better of Lonely
7:26 Soar & Nytrix - Illuminate
11:53 Last Heroes - Love Like Us (feat. RUNN)
15:23 ADVENT - Last Mistake (feat. Akacia)
18:28 Yetep & RUNN - Alright
21:25 NIKAI - Rain Inside (feat. Nomeli)
24:45 Sabai - Memories (feat. Claire Ridgely)
27:29 Seven Lions - Only Now (Feat. Tyler Graves) [MitiS Remix]
```
*/
func TestInterpretVideoDescriptionLine(t *testing.T) {
	trackList := map[string]*app.Track{
		"0:00 Seven Lions, Last Heroes & HALIENE - Don't Wanna Fall": app.NewTrack("don't wanna fall",
			app.NewTrackArtist("seven lions", true),
			app.NewTrackArtist("last heroes", false),
			app.NewTrackArtist("haliene", false)),

		"3:22 Nurko feat. RØRY - Better of Lonely": app.NewTrack("better of lonely",
			app.NewTrackArtist("nurko", true),
			app.NewTrackArtist("røry", false)),

		"7:26 Soar & Nytrix - Illuminate": app.NewTrack("illuminate",
			app.NewTrackArtist("soar", true),
			app.NewTrackArtist("nytrix", false)),

		"11:53 Last Heroes - Love Like Us (feat. RUNN)": app.NewTrack("love like us",
			app.NewTrackArtist("last heroes", true),
			app.NewTrackArtist("runn", false)),

		"15:23 ADVENT - Last Mistake (feat. Akacia)": app.NewTrack("last mistake",
			app.NewTrackArtist("advent", true),
			app.NewTrackArtist("akacia", false)),

		"18:28 Yetep & RUNN - Alright": app.NewTrack("alright",
			app.NewTrackArtist("yetep", true),
			app.NewTrackArtist("runn", false)),

		"21:25 NIKAI - Rain Inside (feat. Nomeli)": app.NewTrack("rain inside",
			app.NewTrackArtist("nikai", true),
			app.NewTrackArtist("nomeli", false)),

		"24:45 Sabai - Memories (feat. Claire Ridgely)": app.NewTrack("memories",
			app.NewTrackArtist("sabai", true),
			app.NewTrackArtist("claire ridgely", false)),

		"27:29 Seven Lions - Only Now (Feat. Tyler Graves) [MitiS Remix]": app.NewTrack("only now",
			app.NewTrackArtist("seven lions", true),
			app.NewTrackArtist("tyler graves", false)),
	}

	for trackText, expectedTrack := range trackList {
		t.Run(fmt.Sprintf("should parse track '%s'", trackText), func(t *testing.T) {
			isMatch, artists, track := app.InterpretVideoDescriptionLine(trackText)
			if !isMatch {
				t.Fatalf("expected isMatch to be true")
			}

			if track != expectedTrack.Track {
				t.Fatalf("expected track to be '%s' but got '%s'", expectedTrack.Track, track)
			}

			assertArtists(t, expectedTrack.Artists, artists)
		})

	}
}

func assertArtists(t *testing.T, expectedArtists []*app.TrackArtist, actualArtists []*app.TrackArtist) {
	if len(actualArtists) != len(expectedArtists) {
		for i, a := range actualArtists {
			log.Printf("artist %d: %s (%t)", i+1, a.Name, a.IsMainTrackArtist)
		}
		t.Fatalf("expected %d artists but got %d", len(expectedArtists), len(actualArtists))
	}

	for _, expectedArtist := range expectedArtists {
		foundMatch := false
		var match *app.TrackArtist

		for _, actualArtist := range actualArtists {
			if actualArtist.Name == expectedArtist.Name {
				foundMatch = true
				match = actualArtist
				break
			}
		}

		if !foundMatch || match == nil {
			for i, a := range actualArtists {
				log.Printf("artist %d: %s (%t)", i+1, a.Name, a.IsMainTrackArtist)
			}
			t.Fatalf("expected to find artist '%s' but it wasn't in the list", expectedArtist.Name)
		}

		if match.IsMainTrackArtist != expectedArtist.IsMainTrackArtist {
			t.Fatalf("expected artist '%s' to be the main track artist", expectedArtist.Name)
		}
	}
}
