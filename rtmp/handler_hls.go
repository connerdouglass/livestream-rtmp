package rtmp

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/godocompany/livestream-rtmp/api"
	"github.com/godocompany/livestream-rtmp/hls"
	"github.com/nareix/joy4/av"
)

type HlsHandlerConfig struct {
	WorkDir string
}

type HlsStreamHandlerFactory struct {
	BasePath    string
	WorkDir     string
	handlers    []*HlsStreamHandler
	handlersMut sync.RWMutex
}

func (hf *HlsStreamHandlerFactory) NewHandler(streamConfig *api.StreamPublishConfig) (StreamHandler, error) {

	// Create the HLS publisher
	hlsPub := &hls.Publisher{
		WorkDir: path.Join(hf.WorkDir, fmt.Sprintf("stream_%s", streamConfig.StreamID)),
	}

	// Make sure the directory exists
	if err := os.MkdirAll(hlsPub.WorkDir, 0700); err != nil {
		return nil, err
	}

	// Create the handler instance
	handler := &HlsStreamHandler{
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

func (hf *HlsStreamHandlerFactory) parseStreamIDFromRequest(req *http.Request) string {

	// If there is no request or no URL, return empty string
	if req == nil || req.URL == nil {
		return ""
	}

	// Get the portion of the URL after the base path
	relativePath := strings.TrimPrefix(req.URL.Path, hf.BasePath)

	// Split the URL up by slashes
	pathParts := strings.Split(relativePath, "/")

	// If there are no parts
	if len(pathParts) == 0 {
		return ""
	}

	// Return the first part otherwise
	return pathParts[0]

}

func (hf *HlsStreamHandlerFactory) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	// Get the stream ID from the URL
	streamID := hf.parseStreamIDFromRequest(req)

	// If there is an identifier
	if len(streamID) > 0 {

		// Acquire access to the handlers slice as a reader
		hf.handlersMut.RLock()
		defer hf.handlersMut.RUnlock()

		// Loop through the handlers
		for _, handler := range hf.handlers {
			if handler.streamConfig.StreamID == streamID {
				handler.hlsPub.ServeHTTP(rw, req)
				return
			}
		}

	}

	// If we get here, serve a 404 error
	rw.WriteHeader(http.StatusNotFound)
	rw.Write([]byte("Not found"))

}

// HlsStreamHandler is the stream handler
type HlsStreamHandler struct {
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
	h.hlsPub.Close()
	return nil
}
