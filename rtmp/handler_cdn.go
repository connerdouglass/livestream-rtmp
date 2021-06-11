package rtmp

import (
	"github.com/godocompany/livestream-rtmp/api"
	"github.com/nareix/joy4/av"
)

type CdnHandlerConfig struct {
	// CDN API details
}

// CdnStreamHandlerFactory returns a stream handler that sends
func CdnStreamHandlerFactory(cdnConfig *CdnHandlerConfig) func(*api.StreamPublishConfig) (StreamHandler, error) {
	return func(streamConfig *api.StreamPublishConfig) (StreamHandler, error) {

		return &CdnStreamHandler{}, nil

	}
}

// CdnStreamHandler is the CDN stream handler
type CdnStreamHandler struct {
	streams          []av.CodecData
	videoStreamIndex int
}

// WriteHeader writes the header data for the streams
func (h *CdnStreamHandler) WriteHeader(streams []av.CodecData) error {

	// Copy the streams
	h.streams = streams

	// Identify the video stream index
	for i, stream := range streams {
		if stream.Type().IsVideo() {
			h.videoStreamIndex = i
		}
	}

	// No error
	return nil

}

// WritePacket writes a stream packet to the CDN
func (h *CdnStreamHandler) WritePacket(packet av.Packet) error {

	// Write the packet
	// ...

	// Return without error
	return nil

}

func (h *CdnStreamHandler) WriteTrailer() error {
	return nil
}
