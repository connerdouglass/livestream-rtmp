package dashmpd

import (
	"encoding/xml"
	"time"

	"github.com/godocompany/livestream-rtmp/hls/internal/ratedetect"
)

type MPD struct {
	XMLName  xml.Name `xml:"urn:mpeg:dash:schema:mpd:2011 MPD"`
	ID       string   `xml:"id,attr"`
	Profiles string   `xml:"profiles,attr"`
	Type     string   `xml:"type,attr"`

	AvailabilityStartTime time.Time `xml:"availabilityStartTime,attr"`
	PublishTime           time.Time `xml:"publishTime,attr"`
	MinimumUpdatePeriod   Duration  `xml:"minimumUpdatePeriod,attr"`
	MaxSegmentDuration    Duration  `xml:"maxSegmentDuration,attr"`
	MinBufferTime         Duration  `xml:"minBufferTime,attr"`
	TimeShiftBufferDepth  Duration  `xml:"timeShiftBufferDepth,attr"`

	Period    []Period
	UTCTiming *UTCTiming
}

type Period struct {
	ID    string   `xml:"id,attr"`
	Start Duration `xml:"start,attr"`

	AdaptationSet []AdaptationSet
}

type AdaptationSet struct {
	ContentType      string          `xml:"contentType,attr"`
	Lang             string          `xml:"lang,attr,omitempty"`
	SegmentAlignment bool            `xml:"segmentAlignment,attr"`
	MaxFrameRate     ratedetect.Rate `xml:"maxFrameRate,attr,omitempty"`
	MaxWidth         int             `xml:"maxWidth,attr,omitempty"`
	MaxHeight        int             `xml:"maxHeight,attr,omitempty"`
	PAR              string          `xml:"par,attr,omitempty"`

	SegmentTemplate SegmentTemplate
	Representation  []Representation
}

type SegmentTemplate struct {
	Duration       int    `xml:"duration,attr,omitempty"`
	Initialization string `xml:"initialization,attr"`
	Media          string `xml:"media,attr"`
	StartNumber    int    `xml:"startNumber,attr"`
	Timescale      int    `xml:"timescale,attr"`

	AvailabilityTimeComplete string  `xml:"availabilityTimeComplete,attr,omitempty"`
	AvailabilityTimeOffset   float64 `xml:"availabilityTimeOffset,attr,omitempty"`

	SegmentTimeline *SegmentTimeline
}

type SegmentTimeline struct {
	Segments []Segment `xml:"S"`
}

type Segment struct {
	Time     uint64 `xml:"t,attr,omitempty"`
	Duration int    `xml:"d,attr,omitempty"`
	Repeat   int    `xml:"r,attr,omitempty"`
}

type Representation struct {
	ID                string          `xml:"id,attr"`
	AudioSamplingRate int             `xml:"audioSamplingRate,attr,omitempty"`
	Bandwidth         int             `xml:"bandwidth,attr"`
	Codecs            string          `xml:"codecs,attr"`
	MimeType          string          `xml:"mimeType,attr"`
	FrameRate         ratedetect.Rate `xml:"frameRate,attr,omitempty"`
	Width             int             `xml:"width,attr,omitempty"`
	Height            int             `xml:"height,attr,omitempty"`
	SAR               string          `xml:"sar,attr,omitempty"`

	AudioChannelConfiguration *AudioChannelConfiguration
}

type AudioChannelConfiguration struct {
	SchemeID string `xml:"schemeIdUri,attr"`
	Value    int    `xml:"value,attr"`
}

type UTCTiming struct {
	Scheme string `xml:"schemeIdUri,attr"`
	Value  string `xml:"value,attr"`
}
