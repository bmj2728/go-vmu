package metadata

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

//builder functions maybe
