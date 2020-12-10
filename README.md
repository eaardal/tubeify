# Tubeify

Scrapes track lists in YouTube video descriptions and adds the songs to a Spotify playlist if matches are found.
Made with compilation videos in mind.

### Example

The video https://www.youtube.com/watch?v=_EopQtSVZQY has a video description that contains the track list:
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
This app will try to find a track list such as this and parse it into a list of artist(s) and song names.
It will then try to search Spotify for the scraped tracks and if it finds a match, add the tracks to the provided Spotify playlist. 

## Setup

### YouTube Authentication

- Create or open an existing [Google Cloud project](https://console.cloud.google.com/) you own.
- Search for `YouTube Data API v3` (it's under _API & Services_) and click the `Enable` button.
- Click the `Manage` button that replaced the `Enable` button if the page did not redirect you automatically.
- You should be on an overview dashboard page for YouTube Data API v3. Click the `Credentials` item on the left side menu.
- Click `+ Create Credentials` in the top menu and select `API Key`. Copy the API key.
- Paste in the API key in a .env file for this project, or an environment variable named `YOUTUBE_API_KEY`. See below.

### Spotify Authentication

- Log in to the Spotify Developer Dashboard: https://developer.spotify.com/dashboard/
- Create a new app
- Copy the `Client ID` and paste it into your .env file (see below).
- Click the `Show Client Secret` button and copy & paste it into the .env file.
- Click the `Edit Settings` button and under `Redirect URIs`, add `http://localhost:8080`

### Local Environment

- Create a `.env` file

```
YOUTUBE_API_KEY=your_key
SPOTIFY_CLIENT_ID=your_id
SPOTIFY_CLIENT_SECRET=your_secret
```

- Optional: Create a Makefile for convenience:

```
scrape:
	go run ./main.go \
		-youtube-video-id your_video_id \
		-spotify-playlist-id your_playlist_id
``` 

### Running

- With the environment variables set up, run the app with `go run main.go` or `make scrape`

## Resources

### Spotify
- [Spotify Golang SDK](https://github.com/zmb3/spotify)
- [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/login)

### YouTube
- [YouTube API Reference: list call](https://developers.google.com/youtube/v3/docs/channels/list?apix_params=%7B%22part%22%3A%5B%22snippet%2CcontentDetails%2Cstatistics%22%5D%2C%22id%22%3A%5B%22UC_x5XG1OV2P6uZZ5FSM9Ttw%22%5D%7D)
- [YouTube API Sample Requests](https://developers.google.com/youtube/v3/sample_requests)
- [YouTube API Code Samples](https://developers.google.com/youtube/v3/code_samples/go)