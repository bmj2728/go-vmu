package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadata_ToMap_Empty(t *testing.T) {
	// Test with an empty metadata struct
	meta := NewMetadata()

	result, err := meta.ToMap()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result, "Map should be empty for an empty metadata struct")
}

func TestMetadata_ToMap_Nil(t *testing.T) {
	// Test with a nil metadata struct
	var meta *Metadata = nil

	result, err := meta.ToMap()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "metadata is nil")
}

func TestMetadata_ToMap_FullData(t *testing.T) {
	// Test with a fully populated metadata struct
	meta := &Metadata{
		Title:     "Test Title",
		Plot:      "Test Plot",
		Runtime:   120,
		ShowTitle: "Test Show",
		Season:    1,
		Episode:   2,
		Genres:    "Action, Comedy",
		IMDBID:    "tt1234567",
		TVDBID:    "12345",
		Year:      2023,
		Writer:    "Test Writer",
		Credits:   "Test Credits",
		Directors: "Director 1, Director 2",
		Actors:    "Actor 1, Actor 2",
	}

	result, err := meta.ToMap()

	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify all fields are correctly mapped
	assert.Equal(t, "Test Title", result["title"])
	assert.Equal(t, "Test Plot", result["plot"])
	assert.Equal(t, 120, result["runtime"])
	assert.Equal(t, "Test Show", result["showtitle"])
	assert.Equal(t, 1, result["season"])
	assert.Equal(t, 2, result["episode"])
	assert.Equal(t, "Action, Comedy", result["genre"])
	assert.Equal(t, "tt1234567", result["imdb_id"])
	assert.Equal(t, "12345", result["tvdb_id"])
	assert.Equal(t, 2023, result["year"])
	assert.Equal(t, "Test Writer", result["writer"])
	assert.Equal(t, "Test Credits", result["credits"])
	assert.Equal(t, "Director 1, Director 2", result["director"])
	assert.Equal(t, "Actor 1, Actor 2", result["actor"])
}

func TestMetadata_ToMap_PartialData(t *testing.T) {
	// Test with a partially populated metadata struct
	meta := &Metadata{
		Title:   "Test Title",
		Runtime: 120,
		Season:  1,
		IMDBID:  "tt1234567",
		Year:    2023,
	}

	result, err := meta.ToMap()

	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify only non-empty fields are in the map
	assert.Equal(t, 5, len(result))
	assert.Equal(t, "Test Title", result["title"])
	assert.Equal(t, 120, result["runtime"])
	assert.Equal(t, 1, result["season"])
	assert.Equal(t, "tt1234567", result["imdb_id"])
	assert.Equal(t, 2023, result["year"])

	// Verify empty fields are not in the map
	_, hasPlot := result["plot"]
	assert.False(t, hasPlot)
	_, hasShowTitle := result["showtitle"]
	assert.False(t, hasShowTitle)
}
