package nfo

import (
	"encoding/xml"
	"fmt"
	"github.com/bmj2728/go-vmu/internal/utils"
	"github.com/rs/zerolog/log"
	"os"
)

// NFONotFoundError represents an error message when an NFO file cannot be located.
// NFOReadError represents an error message for issues reading an NFO file.
// NFOUnMarshallError represents an error message when unmarshalling an NFO file fails.
const (
	NFONotFoundError   = "error locating nfo file"
	NFOReadError       = "error reading nfo file"
	NFOUnMarshallError = "error unmarshalling nfo file"
)

// ParseEpisodeNFO parses the given NFO file path into an EpisodeDetails struct.
// Returns an error if the file cannot be opened, read, or unmarshalled into the struct.
func ParseEpisodeNFO(path string) (*EpisodeDetails, error) {
	//We're testing the validity of the nfo file path in the matching step
	log.Debug().Strs("nfo", []string{path, "loading"}).Msg("Parsing NFO file")
	//Open the file
	file, err := os.Open(path)
	if err != nil {
		log.Error().Err(err).Msg(NFOReadError)
		return nil, fmt.Errorf(NFOReadError+": %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error().Err(err).Msg(NFOReadError)
		}
	}(file)
	//Parse the xml using std lib
	details := &EpisodeDetails{}
	decoder := xml.NewDecoder(file)
	if err := decoder.Decode(details); err != nil {
		log.Error().Err(err).Msg(NFOUnMarshallError)
		return nil, fmt.Errorf(NFOUnMarshallError+": %v", err)
	}
	return details, nil
}

// MatchEpisodeFile attempts to find an NFO file corresponding to the given episode file path.
// It checks both the existence of the given file and the deduced NFO file in the same directory.
// Returns the path to the located NFO file or an error if it does not exist.
func MatchEpisodeFile(path string) (string, error) {
	log.Debug().Strs("nfo", []string{path, "matching"}).Msg("Matching NFO file")
	//validate existence or just give up
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Error().Err(err).Msg(NFONotFoundError)
		return "", fmt.Errorf(NFONotFoundError+": %v", err)
	}
	//get working directory
	/*old
	//wd := filepath.Dir(path)
	//base := strings.Split(filepath.Base(path), ".")[0]
	//nfoPath := wd + "/" + base + ".nfo"
	//log.Debug().Strs("nfo", []string{path, "working directory " + wd, "base " + base, "expected " + nfoPath}).Msg("Working directory")
	////Test our guess
	//if _, err := os.Stat(nfoPath); os.IsNotExist(err) {
	//	log.Error().Err(err).Msg(NFONotFoundError)
	//	return "", fmt.Errorf(NFONotFoundError+": %v", err)
	//}
	*/
	return utils.NFOPath(path)
}
