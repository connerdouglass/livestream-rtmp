package rtmp

import (
	"os"

	"github.com/godocompany/livestream-rtmp/api"
	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/format/mp4"
)

type CdnHandlerConfig struct {
	Filename string
}

// CdnStreamHandlerFactory returns a stream handler that sends
func CdnStreamHandlerFactory(cdnConfig *CdnHandlerConfig) func(*api.StreamPublishConfig) (StreamHandler, error) {
	return func(streamConfig *api.StreamPublishConfig) (StreamHandler, error) {
		return NewCdnStreamHandler(cdnConfig, streamConfig)
	}
}

// CdnStreamHandler is the CDN stream handler
type CdnStreamHandler struct {
	// streams          []av.CodecData
	// videoStreamIndex int
	file  *os.File
	muxer *mp4.Muxer
}

func NewCdnStreamHandler(
	cdnConfig *CdnHandlerConfig,
	streamConfig *api.StreamPublishConfig,
) (*CdnStreamHandler, error) {

	// Create the file with the filename
	file, err := os.Create(cdnConfig.Filename)
	if err != nil {
		return nil, err
	}

	// Create the muxer for the file
	muxer := mp4.NewMuxer(file)

	// Return the handler
	return &CdnStreamHandler{
		file:  file,
		muxer: muxer,
	}, nil

}

// WriteHeader writes the header data for the streams
func (h *CdnStreamHandler) WriteHeader(streams []av.CodecData) error {
	return h.muxer.WriteHeader(streams)

	// // Copy the streams
	// h.streams = streams

	// // Identify the video stream index
	// for i, stream := range streams {
	// 	if stream.Type().IsVideo() {
	// 		h.videoStreamIndex = i
	// 	}
	// }

	// // No error
	// return nil

}

// WritePacket writes a stream packet to the CDN
func (h *CdnStreamHandler) WritePacket(packet av.Packet) error {
	return h.muxer.WritePacket(packet)
}

func (h *CdnStreamHandler) WriteTrailer() error {
	return h.muxer.WriteTrailer()
}

func (h *CdnStreamHandler) Close() error {
	return h.file.Close()
}
