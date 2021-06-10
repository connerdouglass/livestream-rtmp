package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) post(url string, body interface{}) (*http.Response, error) {

	// Format the request URL
	fullUrl := fmt.Sprintf("%s%s", c.Hostname, url)

	// The body string to send
	bodyBytes := []byte("{}")
	if body != nil {
		bytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyBytes = bytes
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.RtmpPasscode))

	// Send the request
	return http.DefaultClient.Do(req)

}

func (c *Client) request(url string, body interface{}, result interface{}) error {

	// Send the post request
	res, err := c.post(url, body)
	if err != nil {
		return err
	}

	// If the result passed in is nil, ignore it
	if result == nil {
		return nil
	}

	// Create the buffer
	buf := bytes.Buffer{}
	if _, err := buf.ReadFrom(res.Body); err != nil {
		return err
	}

	// Unmarshal the JSON into the result
	return json.Unmarshal(
		buf.Bytes(),
		result,
	)

}
