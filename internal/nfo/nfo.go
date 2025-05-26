package nfo

import "encoding/xml"

type EpisodeDetails struct {
	XMLName   xml.Name `xml:"episodedetails"`
	Plot      string   `xml:"plot"`
	LockData  bool     `xml:"lockdata"`
	DateAdded string   `xml:"dateadded"`
	Title     string   `xml:"title"`
	Director  []string `xml:"director"`
	Writer    string   `xml:"writer"`
	Credits   string   `xml:"credits"`
	Rating    float64  `xml:"rating"`
	Year      int      `xml:"year"`
	MPAA      string   `xml:"mpaa,omitempty"`
	IMDBID    string   `xml:"imdbid"`
	TVDBID    string   `xml:"tvdbid"`
	Runtime   int      `xml:"runtime"`
	Genre     []string `xml:"genre"`
	Art       Art      `xml:"art"`
	Actor     []Actor  `xml:"actor"`
	ShowTitle string   `xml:"showtitle"`
	Episode   int      `xml:"episode"`
	Season    int      `xml:"season"`
	Aired     string   `xml:"aired"`
	FileInfo  FileInfo `xml:"fileinfo"`
}

type Art struct {
	Poster string `xml:"poster"`
}

type Actor struct {
	Name      string `xml:"name"`
	Role      string `xml:"role"`
	Type      string `xml:"type"`
	SortOrder int    `xml:"sortorder"`
	Thumb     string `xml:"thumb,omitempty"`
}

type FileInfo struct {
	StreamDetails StreamDetails `xml:"streamdetails"`
}

type StreamDetails struct {
	Video    []VideoStream    `xml:"video"`
	Audio    []AudioStream    `xml:"audio"`
	Subtitle []SubtitleStream `xml:"subtitle"`
}

type VideoStream struct {
	Codec             string  `xml:"codec"`
	MiCodec           string  `xml:"micodec"`
	Bitrate           int     `xml:"bitrate"`
	Width             int     `xml:"width"`
	Height            int     `xml:"height"`
	Aspect            string  `xml:"aspect"`
	AspectRatio       string  `xml:"aspectratio"`
	Framerate         float64 `xml:"framerate"`
	Language          string  `xml:"language,omitempty"`
	ScanType          string  `xml:"scantype"`
	Default           bool    `xml:"default"`
	Forced            bool    `xml:"forced"`
	Duration          int     `xml:"duration"`
	DurationInSeconds int     `xml:"durationinseconds"`
}

type AudioStream struct {
	Codec        string `xml:"codec"`
	MiCodec      string `xml:"micodec"`
	Bitrate      int    `xml:"bitrate"`
	Language     string `xml:"language"`
	ScanType     string `xml:"scantype"`
	Channels     int    `xml:"channels"`
	SamplingRate int    `xml:"samplingrate"`
	Default      bool   `xml:"default"`
	Forced       bool   `xml:"forced"`
}

type SubtitleStream struct {
	Codec    string `xml:"codec"`
	MiCodec  string `xml:"micodec"`
	Width    int    `xml:"width"`
	Height   int    `xml:"height"`
	Language string `xml:"language"`
	ScanType string `xml:"scantype"`
	Default  bool   `xml:"default"`
	Forced   bool   `xml:"forced"`
}
