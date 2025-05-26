package metadata

import "fmt"

type Metadata struct {
	Title     string
	Plot      string
	Runtime   int
	ShowTitle string
	Season    int
	Episode   int
	Genres    []string
	IMDBID    string
	TVDBID    string
	Year      int
	Writer    string
	Credits   string
	Directors []string
	//Will need to process this from actor structs - only need names
	Actors []string
}

func NewMetadata() *Metadata {
	return &Metadata{}
}

func (m *Metadata) ToMap() (map[string]interface{}, error) {
	//We need this lil guy
	var metaFields map[string]interface{}

	// do nothing if metadata is nil
	if m == nil {
		return nil, fmt.Errorf("metadata is nil")
	}

	//dynamicaly build the map
	if m.Title != "" {
		metaFields["title"] = m.Title
	}
	if m.Plot != "" {
		metaFields["plot"] = m.Plot
	}
	if m.Runtime != 0 {
		metaFields["runtime"] = m.Runtime
	}
	if m.ShowTitle != "" {
		metaFields["showtitle"] = m.ShowTitle
	}
	if m.Season != 0 {
		metaFields["season"] = m.Season
	}
	if m.Episode != 0 {
		metaFields["episode"] = m.Episode
	}
	if m.Genres != nil {
		metaFields["genre"] = m.Genres
	}
	if m.IMDBID != "" {
		metaFields["imdb_id"] = m.IMDBID
	}
	if m.TVDBID != "" {
		metaFields["tvdb_id"] = m.TVDBID
	}
	if m.Year != 0 {
		metaFields["year"] = m.Year
	}
	if m.Writer != "" {
		metaFields["writer"] = m.Writer
	}
	if m.Credits != "" {
		metaFields["credits"] = m.Credits
	}
	if m.Directors != nil {
		metaFields["director"] = m.Directors
	}
	if m.Actors != nil {
		metaFields["actor"] = m.Actors
	}
	//return the map
	return metaFields, nil
}
