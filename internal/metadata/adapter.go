package metadata

import (
	"fmt"
	"github.com/bmj2728/go-vmu/internal/nfo"
	"github.com/rs/zerolog/log"
	"strings"
)

type NFOAdapter struct {
	Details  *nfo.EpisodeDetails
	Metadata *Metadata
}

func NewNFOAdapter(details *nfo.EpisodeDetails) *NFOAdapter {
	return &NFOAdapter{
		Details:  details,
		Metadata: NewMetadata(),
	}
}

func (a *NFOAdapter) TranslateNFO() (*Metadata, error) {
	if a.Details == nil {
		log.Error().Msg("NFO details not set")
		return nil, fmt.Errorf("NFO details not set")
	}

	//convert to strings

	var actorsNames []string

	for _, actor := range a.Details.Actor {
		log.Debug().Msgf("Actor: %v", actor)
		actorsNames = append(actorsNames, actor.Name)
	}

	actors := strings.Join(actorsNames, ", ")
	log.Debug().Msgf("Actors: %s", actors)
	genres := strings.Join(a.Details.Genre, ", ")
	log.Debug().Msgf("Genres: %s", genres)
	directors := strings.Join(a.Details.Director, ", ")
	log.Debug().Msgf("Directors: %s", directors)

	a.Metadata.Title = a.Details.Title
	log.Debug().Msgf("Title: %s", a.Metadata.Title)
	a.Metadata.Plot = a.Details.Plot
	log.Debug().Msgf("Plot: %s", a.Metadata.Plot)
	a.Metadata.Runtime = a.Details.Runtime
	log.Debug().Msgf("Runtime: %d", a.Metadata.Runtime)
	a.Metadata.ShowTitle = a.Details.ShowTitle
	log.Debug().Msgf("ShowTitle: %s", a.Metadata.ShowTitle)
	a.Metadata.Season = a.Details.Season
	log.Debug().Msgf("Season: %d", a.Metadata.Season)
	a.Metadata.Episode = a.Details.Episode
	log.Debug().Msgf("Episode: %d", a.Metadata.Episode)
	a.Metadata.Genres = genres
	log.Debug().Msgf("Genres: %v", a.Metadata.Genres)
	a.Metadata.IMDBID = a.Details.IMDBID
	log.Debug().Msgf("IMDBID: %s", a.Metadata.IMDBID)
	a.Metadata.TVDBID = a.Details.TVDBID
	log.Debug().Msgf("TVDBID: %s", a.Metadata.TVDBID)
	a.Metadata.Year = a.Details.Year
	log.Debug().Msgf("Year: %d", a.Metadata.Year)
	a.Metadata.Writer = a.Details.Writer
	log.Debug().Msgf("Writer: %s", a.Metadata.Writer)
	a.Metadata.Credits = a.Details.Credits
	log.Debug().Msgf("Credits: %s", a.Metadata.Credits)
	a.Metadata.Directors = directors
	log.Debug().Msgf("Directors: %v", a.Metadata.Directors)
	a.Metadata.Actors = actors
	log.Debug().Msgf("Actors: %v", a.Metadata.Actors)

	return a.Metadata, nil
}
