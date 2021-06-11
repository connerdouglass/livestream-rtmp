package rtmp

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/godocompany/livestream-rtmp/api"
	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/format"
	joyRtmp "github.com/nareix/joy4/format/rtmp"
)

type activeStream struct {
	StreamKey     string
	StreamID      string
	PublishConfig *api.StreamPublishConfig
	StreamHandler StreamHandler
}

// StreamHandler processes incoming stream packets
type StreamHandler interface {
	av.Muxer
}

type Server struct {
	Address          string
	Api              *api.Client
	NewStreamHandler func(*api.StreamPublishConfig) (StreamHandler, error)
	activeStreams    map[string]*activeStream
	activeStreamsMut sync.RWMutex
}

func (s *Server) Run() {

	// Create the active streams map
	s.activeStreams = make(map[string]*activeStream)

	// Prepare all of the file formats and codec handlers
	format.RegisterAll()

	// Setup the RTMP server
	server := &joyRtmp.Server{
		Addr:          s.Address,
		HandlePublish: s.handlePublish,
	}

	// Listen and serve the RTMP server
	if err := server.ListenAndServe(); err != nil {
		fmt.Fprintln(os.Stderr, "Fatal: Couldn't run the stream server:", err)
	}

}

// handlePublish handles a connection by a streamer publishing to the RTMP server
func (s *Server) handlePublish(conn *joyRtmp.Conn) {

	// Close the connection when everything is done
	defer conn.Close()

	// Get the stream key from the connection
	streamKey := strings.TrimPrefix(conn.URL.Path, "/")

	// Get the stream info from the stream key
	config, err := s.Api.GetStreamPublishData(streamKey)
	if err != nil {
		fmt.Println("Error getting stream config: ", err.Error())
		return
	}
	if config == nil {
		fmt.Println("No stream config for stream key: ", streamKey)
		return
	}

	// Add the stream to the slice
	stream := &activeStream{
		StreamKey:     streamKey,
		StreamID:      config.StreamID,
		PublishConfig: config,
	}

	// Add the stream to the slice
	s.activeStreamsMut.Lock()
	s.activeStreams[config.StreamID] = stream
	s.activeStreamsMut.Unlock()

	// Remove the stream from the slice when things end
	defer func() {
		s.activeStreamsMut.Lock()
		defer s.activeStreamsMut.Unlock()
		delete(s.activeStreams, config.StreamID)
	}()

	// Get the streams from the connection
	streams, err := conn.Streams()
	if err != nil {
		fmt.Println("Error getting streams: ", err.Error())
		return
	}

	// Let the API know the stream has started
	if err := s.Api.MarkStreamStarted(config.StreamID); err != nil {
		fmt.Println("Error marking stream as started: ", err.Error())
		return
	}

	// Defer the cleanup function
	defer func() {

		// Let the API know the stream has ended
		if err := s.Api.MarkStreamEnded(config.StreamID); err != nil {
			fmt.Println("Error marking stream as ended: ", err.Error())
		}

	}()

	fmt.Println("Info: The stream has started")

	// Create the packet writer to write packets to the CDN
	writer, err := s.NewStreamHandler(config)
	if err != nil {
		fmt.Println("Error creating packet writer: ", writer)
		return
	}

	// Write the headers to the queue
	if err := writer.WriteHeader(streams); err != nil {
		fmt.Println("Error writing stream headers: ", err.Error())
		return
	}

	// Copy all of the packets from the stream, until it has concluded
	if err := avutil.CopyPackets(writer, conn); err == io.EOF {
		fmt.Println("Info: The server has stopped streaming.")
	} else if err != nil {
		fmt.Println("Stream ended with error: ", err.Error())
	}

}
