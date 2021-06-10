package api

import "fmt"

// Client is the livestream-api client
type Client struct {
	Hostname     string
	RtmpPasscode string
}

type StreamPublishConfig struct {
	StreamID string `json:"stream_id"`
}

// GetStreamPublishData gets the data for a stream from its stream key
func (c *Client) GetStreamPublishData(streamKey string) (*StreamPublishConfig, error) {
	return &StreamPublishConfig{
		StreamID: "helloworld",
	}, nil
}

func (c *Client) MarkStreamStarted(streamID string) error {
	fmt.Println("Marking stream as STARTED: ", streamID)
	return nil
}

func (c *Client) MarkStreamEnded(streamID string) error {
	fmt.Println("Marking stream as ENDED: ", streamID)
	return nil
}
