package api

// Client is the livestream-api client
type Client struct {
	Hostname     string
	RtmpPasscode string
}

type StreamPublishConfig struct {
	StreamID string `json:"stream_id"`
}

// GetStreamConfig gets the data for a stream from its stream key
func (c *Client) GetStreamConfig(streamKey string) (*StreamPublishConfig, error) {

	// Create the request data
	req := map[string]interface{}{
		"stream_key": streamKey,
	}

	// Send the request and handle error
	var res StreamPublishConfig
	if err := c.request("/v1/rtmp/stream/get-config", req, &res); err != nil {
		return nil, err
	}

	// Return the response data
	return &res, nil

}

func (c *Client) SetStreaming(streamID string, streaming bool) error {
	return c.request(
		"/v1/rtmp/stream/set-streaming",
		map[string]interface{}{
			"stream_id": streamID,
			"streaming": streaming,
		},
		nil,
	)
}
