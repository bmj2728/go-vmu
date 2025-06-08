package metadata

import (
	"testing"

	"github.com/bmj2728/go-vmu/internal/nfo"
	"github.com/stretchr/testify/assert"
)

func TestNFOAdapter_TranslateNFO_Nil(t *testing.T) {
	// Test with nil NFO details
	adapter := NewNFOAdapter(nil)

	result, err := adapter.TranslateNFO()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "NFO details not set")
}

func TestNFOAdapter_TranslateNFO_Empty(t *testing.T) {
	// Test with empty NFO details
	details := &nfo.EpisodeDetails{}
	adapter := NewNFOAdapter(details)

	result, err := adapter.TranslateNFO()

	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify all fields are empty
	assert.Empty(t, result.Title)
	assert.Empty(t, result.Plot)
	assert.Zero(t, result.Runtime)
	assert.Empty(t, result.ShowTitle)
	assert.Zero(t, result.Season)
	assert.Zero(t, result.Episode)
	assert.Empty(t, result.Genres)
	assert.Empty(t, result.IMDBID)
	assert.Empty(t, result.TVDBID)
	assert.Zero(t, result.Year)
	assert.Empty(t, result.Writer)
	assert.Empty(t, result.Credits)
	assert.Empty(t, result.Directors)
	assert.Empty(t, result.Actors)
}

func TestNFOAdapter_TranslateNFO_FullData(t *testing.T) {
	// Test with fully populated NFO details
	details := &nfo.EpisodeDetails{
		Title:     "Test Title",
		Plot:      "Test Plot",
		Runtime:   120,
		ShowTitle: "Test Show",
		Season:    1,
		Episode:   2,
		Genre:     []string{"Action", "Comedy"},
		IMDBID:    "tt1234567",
		TVDBID:    "12345",
		Year:      2023,
		Writer:    "Test Writer",
		Credits:   "Test Credits",
		Director:  []string{"Director 1", "Director 2"},
		Actor: []nfo.Actor{
			{Name: "Actor 1", Role: "Role 1"},
			{Name: "Actor 2", Role: "Role 2"},
		},
	}

	adapter := NewNFOAdapter(details)

	result, err := adapter.TranslateNFO()

	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify all fields are correctly mapped
	assert.Equal(t, "Test Title", result.Title)
	assert.Equal(t, "Test Plot", result.Plot)
	assert.Equal(t, 120, result.Runtime)
	assert.Equal(t, "Test Show", result.ShowTitle)
	assert.Equal(t, 1, result.Season)
	assert.Equal(t, 2, result.Episode)
	assert.Equal(t, "Action, Comedy", result.Genres)
	assert.Equal(t, "tt1234567", result.IMDBID)
	assert.Equal(t, "12345", result.TVDBID)
	assert.Equal(t, 2023, result.Year)
	assert.Equal(t, "Test Writer", result.Writer)
	assert.Equal(t, "Test Credits", result.Credits)
	assert.Equal(t, "Director 1, Director 2", result.Directors)
	assert.Equal(t, "Actor 1, Actor 2", result.Actors)
}

func TestNFOAdapter_TranslateNFO_PartialData(t *testing.T) {
	// Test with partially populated NFO details
	details := &nfo.EpisodeDetails{
		Title:    "Test Title",
		Runtime:  120,
		Season:   1,
		IMDBID:   "tt1234567",
		Year:     2023,
		Director: []string{"Director 1"},
		Actor: []nfo.Actor{
			{Name: "Actor 1", Role: "Role 1"},
		},
	}

	adapter := NewNFOAdapter(details)

	result, err := adapter.TranslateNFO()

	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify populated fields are correctly mapped
	assert.Equal(t, "Test Title", result.Title)
	assert.Equal(t, 120, result.Runtime)
	assert.Equal(t, 1, result.Season)
	assert.Equal(t, "tt1234567", result.IMDBID)
	assert.Equal(t, 2023, result.Year)
	assert.Equal(t, "Director 1", result.Directors)
	assert.Equal(t, "Actor 1", result.Actors)

	// Verify empty fields remain empty
	assert.Empty(t, result.Plot)
	assert.Empty(t, result.ShowTitle)
	assert.Zero(t, result.Episode)
	assert.Empty(t, result.Genres)
	assert.Empty(t, result.TVDBID)
	assert.Empty(t, result.Writer)
	assert.Empty(t, result.Credits)
}
