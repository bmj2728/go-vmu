package metadata

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
)

type Checker interface {
	Compare() bool
}

type MetaChecker struct {
	ExistingMetadata map[string]interface{}
	CompareMetadata  map[string]interface{}
}

func NewMetaChecker(existing map[string]interface{}, compare map[string]interface{}) *MetaChecker {
	return &MetaChecker{
		ExistingMetadata: existing,
		CompareMetadata:  compare,
	}
}

func (m *MetaChecker) Compare() bool {
	//compare is new data since we want this to be the data we check each value
	//normalize the data
	normalizedExisting := make(map[string]interface{})
	for k, v := range m.ExistingMetadata {
		normalizedExisting[strings.ToUpper(k)] = fmt.Sprintf("%v", v)
	}
	normalizedCompare := make(map[string]interface{})
	for k, v := range m.CompareMetadata {
		normalizedCompare[strings.ToUpper(k)] = fmt.Sprintf("%v", v)
	}

	for k, v := range normalizedCompare {
		if normalizedExisting[k] != v {
			log.Info().Msgf("Inconsistency found: %s - Old:%v New:%v", k, m.ExistingMetadata[k], v)
			return false
		}
	}
	//if map processes without a false return true
	log.Info().Msg("No inconsistencies found - skipping file")
	return true
}
