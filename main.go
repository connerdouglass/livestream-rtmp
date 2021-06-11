package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/godocompany/livestream-rtmp/api"
	"github.com/godocompany/livestream-rtmp/rtmp"
	"github.com/joho/godotenv"
)

func main() {

	// Load the .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file: ", err)
	}

	// Create the livestream-api client
	apiClient := &api.Client{
		Hostname:     EnvOrDefault("API_HOSTNAME", "http://localhost:8080"),
		RtmpPasscode: RequireEnv("API_PASSWORD"),
	}

	// Create the handler factory
	handlerFactory := &rtmp.HlsStreamHandlerFactory{
		BasePath: "/hls/",
		WorkDir:  "./hlsdata",
	}

	// Create the RTMP server
	rtmpServer := rtmp.Server{
		Address:          EnvOrDefault("RTMP_ADDR", ":1935"),
		Api:              apiClient,
		NewStreamHandler: handlerFactory.NewHandler,
	}

	go func() {
		http.Handle(handlerFactory.BasePath, handlerFactory)
		http.Handle("/", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			r := strings.NewReader(home)
			http.ServeContent(rw, req, "index.html", time.Time{}, r)
		}))
		http.ListenAndServe(":8080", nil)
	}()

	// Run the RTMP server. This blocks the main goroutine
	rtmpServer.Run()

}

const home = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>HLS demo</title>
<script src="https://cdn.jsdelivr.net/npm/hls.js@latest"></script>
</head>
<body>
<video id="video" muted autoplay controls></video>
<script>
let hls = new Hls();
hls.loadSource('/hls/helloworld/index.m3u8');
hls.attachMedia(document.getElementById('video'));
// hls.on(Hls.Events.MANIFEST_PARSED, () => video.play());
</script>
</body>
</html>
`
