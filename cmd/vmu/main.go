package main

import (
	"github.com/rs/zerolog/log"
	"go-vmu/internal/config"
	"go-vmu/internal/logger"
	"go-vmu/internal/metadata"
	"go-vmu/internal/nfo"
)

func main() {
	cfg, err := config.Load("././config.toml")
	if err != nil {
		panic(err)
	}
	logger.Setup(&cfg.Logger)

	log.Info().Str("startup", "logger").Msg("Logger started")
	log.Info().Msgf("Config: %v", cfg)

	testPath := "/mnt/eagle5/videos/Tv/Westworld/Season 1/Westworld - S01E01 - The Original Bluray-1080p.mkv"

	nfoPath, err := nfo.MatchEpisodeFile(testPath)
	if err != nil {
		log.Error().Err(err).Msg("Error matching NFO file")
	}
	data, err := nfo.ParseEpisodeNFO(nfoPath)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing NFO file")
	}

	adapter := metadata.NewNFOAdapter(data)
	meta, err := adapter.TranslateNFO()
	if err != nil {
		log.Error().Err(err).Msg("Error translating NFO file")
	}
	log.Info().Msgf("\nNFO:\n%v\n", data)
	log.Info().Msgf("\nMetadata:\n%v\n", meta)

}
