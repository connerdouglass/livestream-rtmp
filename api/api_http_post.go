package api

import (
	"bytes"
	"encoding/json"
	"errors"
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

	// Full response body
	var fullResponse map[string]interface{}

	// Unmarshal the JSON into the result
	if err := json.Unmarshal(buf.Bytes(), &fullResponse); err != nil {
		return err
	}

	// If there is an error key in the root
	if errVal, ok := fullResponse["error"]; ok {
		if errStr, ok := errVal.(string); ok {
			return errors.New(errStr)
		}
		return fmt.Errorf("something went wrong: %+v", errVal)
	}

	// Get the data value
	dataVal, ok := fullResponse["data"]
	if !ok {
		return errors.New("response contains no data value")
	}

	// Marshal it back to JSON
	jsonDataBytes, err := json.Marshal(dataVal)
	if err != nil {
		return err
	}

	// Unmarshal the values into the final output response
	if err := json.Unmarshal(jsonDataBytes, &result); err != nil {
		return err
	}

	// Return without error
	return nil

}
