package hlsserver

import (
	_ "embed"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/godocompany/livestream-rtmp/rtmp"
)

//go:embed player.html
var htmlTemplate string

type Server struct {
	HlsFactory   *rtmp.HlsStreamHandlerFactory
	EnablePlayer bool
}

func (s *Server) Run(addr string) {

	// Create the gin server
	r := gin.Default()

	// Configure CORS for the HLS server
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOriginFunc = s.allowCorsOrigin
	corsCfg.AllowCredentials = true
	corsCfg.AddAllowHeaders("Accept", "User-Agent", "Authorization")
	r.Use(cors.New(corsCfg))

	// Add the HLS serving endpoint
	r.GET("/hls/:streamID/*filename", s.serveHls)

	// Add the player endpoint
	if s.EnablePlayer {
		r.GET("/play/:streamID", s.streamDemoViewer)
	}

	// Run the server
	if err := r.Run(addr); err != nil {
		fmt.Println(err.Error())
	}

}

func (s *Server) allowCorsOrigin(origin string) bool {
	return true
}

func (s *Server) serveHls(c *gin.Context) {

	// Get the stream ID from the URL
	streamID := c.Param("streamID")

	// Get the HLS handler with the stream ID
	handler := s.HlsFactory.GetHandler(streamID)
	if handler != nil {
		handler.ServeHTTP(c.Writer, c.Request)
		return
	}

	// Send a 404 message instead
	c.Status(http.StatusNotFound)
	c.Writer.Write([]byte("Not found"))

}

func (s *Server) streamDemoViewer(c *gin.Context) {

	// Get the stream ID from the URL
	streamID := c.Param("streamID")

	// Inject data into the template
	html := strings.ReplaceAll(htmlTemplate, "{{ STREAM_ID }}", streamID)

	// Serve the HTML contents
	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/html")
	c.Writer.Write([]byte(html))

}
