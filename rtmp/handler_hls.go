package rtmp

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	"github.com/godocompany/livestream-rtmp/api"
	"github.com/godocompany/livestream-rtmp/hls"
	"github.com/nareix/joy4/av"
)

type HlsHandlerConfig struct {
	WorkDir string
}

type HlsStreamHandlerFactory struct {
	WorkDir     string
	handlers    []*HlsStreamHandler
	handlersMut sync.RWMutex
}

func (hf *HlsStreamHandlerFactory) NewHandler(streamConfig *api.StreamPublishConfig) (StreamHandler, error) {

	// Create the HLS publisher
	hlsPub := &hls.Publisher{
		FragmentLength: 5 * time.Second,
		WorkDir:        path.Join(hf.WorkDir, fmt.Sprintf("stream_%s", streamConfig.StreamID)),
	}

	// Make sure the directory exists
	if err := os.MkdirAll(hlsPub.WorkDir, 0700); err != nil {
		return nil, err
	}

	// Create the handler instance
	handler := &HlsStreamHandler{
		factory:      hf,
		hlsPub:       hlsPub,
		streamConfig: streamConfig,
	}

	// Add it to the slice
	hf.handlersMut.Lock()
	hf.handlers = append(hf.handlers, handler)
	hf.handlersMut.Unlock()

	// Return the handler
	return handler, nil

}

func (hf *HlsStreamHandlerFactory) GetHandler(streamID string) *HlsStreamHandler {

	// Acquire access to the handlers slice as a reader
	hf.handlersMut.RLock()
	defer hf.handlersMut.RUnlock()

	// Loop through the handlers
	for _, handler := range hf.handlers {
		if handler.streamConfig.StreamID == streamID {
			return handler
		}
	}

	// Return nil instead
	return nil

}

// streamClosed called when a stream has been closed
func (hf *HlsStreamHandlerFactory) streamClosed(stream *HlsStreamHandler) {

	// Acquire access to the handlers slice as a writer
	hf.handlersMut.Lock()
	defer hf.handlersMut.Unlock()

	// Create a new slice for the values, minus this stream
	newHandlers := []*HlsStreamHandler{}
	for _, handler := range hf.handlers {
		if handler.streamConfig.StreamID != stream.streamConfig.StreamID {
			newHandlers = append(newHandlers, handler)
		}
	}
	hf.handlers = newHandlers

}

// HlsStreamHandler is the stream handler
type HlsStreamHandler struct {
	factory      *HlsStreamHandlerFactory
	hlsPub       *hls.Publisher
	streamConfig *api.StreamPublishConfig
}

// WriteHeader writes the header data for the streams
func (h *HlsStreamHandler) WriteHeader(streams []av.CodecData) error {
	return h.hlsPub.WriteHeader(streams)
}

// WritePacket writes a stream packet
func (h *HlsStreamHandler) WritePacket(packet av.Packet) error {
	return h.hlsPub.WritePacket(packet)
}

func (h *HlsStreamHandler) WriteTrailer() error {
	return h.hlsPub.WriteTrailer()
}

func (h *HlsStreamHandler) Close() error {

	// Close the HLS publisher
	h.hlsPub.Close()

	// Remove this stream from the factory
	if h.factory != nil {
		h.factory.streamClosed(h)
	}

	return nil
}

func (h *HlsStreamHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	h.hlsPub.ServeHTTP(rw, req)
}
