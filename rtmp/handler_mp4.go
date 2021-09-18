package rtmp

import (
	"os"

	"github.com/connerdouglass/livestream-rtmp/api"
	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/format/mp4"
)

type Mp4HandlerConfig struct {
	Filename string
}

// Mp4StreamHandlerFactory returns a stream handler that writes to an MP4 file
func Mp4StreamHandlerFactory(mp4Config *Mp4HandlerConfig) func(*api.StreamPublishConfig) (StreamHandler, error) {
	return func(streamConfig *api.StreamPublishConfig) (StreamHandler, error) {
		return NewMp4StreamHandler(mp4Config, streamConfig)
	}
}

// Mp4StreamHandler is the MP4 stream handler
type Mp4StreamHandler struct {
	file  *os.File
	muxer *mp4.Muxer
}

func NewMp4StreamHandler(
	mp4Config *Mp4HandlerConfig,
	streamConfig *api.StreamPublishConfig,
) (*Mp4StreamHandler, error) {

	// Create the file with the filename
	file, err := os.Create(mp4Config.Filename)
	if err != nil {
		return nil, err
	}

	// Create the muxer for the file
	muxer := mp4.NewMuxer(file)

	// Return the handler
	return &Mp4StreamHandler{
		file:  file,
		muxer: muxer,
	}, nil

}

// WriteHeader writes the header data for the streams
func (h *Mp4StreamHandler) WriteHeader(streams []av.CodecData) error {
	return h.muxer.WriteHeader(streams)
}

// WritePacket writes a stream packet to the MP$
func (h *Mp4StreamHandler) WritePacket(packet av.Packet) error {
	return h.muxer.WritePacket(packet)
}

func (h *Mp4StreamHandler) WriteTrailer() error {
	return h.muxer.WriteTrailer()
}

func (h *Mp4StreamHandler) Close() error {
	return h.file.Close()
}
