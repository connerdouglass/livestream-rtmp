package main

import (
	"fmt"

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

	// Create the CDN configuration
	cdnConfig := &rtmp.CdnHandlerConfig{}

	// Create the RTMP server
	rtmpServer := rtmp.Server{
		Address:          EnvOrDefault("RTMP_ADDR", ":1935"),
		Api:              apiClient,
		NewStreamHandler: rtmp.CdnStreamHandlerFactory(cdnConfig),
	}

	// Run the RTMP server. This blocks the main goroutine
	rtmpServer.Run()

}
