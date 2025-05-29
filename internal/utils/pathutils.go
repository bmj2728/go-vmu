package utils

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"strings"
)

// InsertTagToFileName insert user defined text between the filename and extension
// example: newPath := utils.InsertTagToFileName("/home/user/Videos/some_video_file.mkv", "my-tag")
// output: /home/user/Videos/some_video_file.my-tag.mkv
func InsertTagToFileName(path string, tag string) string {
	dir, file := filepath.Split(path)
	ext := filepath.Ext(path)
	name := strings.Replace(file, ext, "", 1)
	newPath := dir + name + "." + tag + ext
	return newPath
}

// NFOPath returns the path to the nfo file for a given video file if it exists
// assumes that nfo is a sidecar with the same name /foo/bar/my_video.avi & /foo/bar/my_video.nfo
func NFOPath(path string) (string, error) {
	//get working directory + name
	dir, file := filepath.Split(path)
	ext := filepath.Ext(path)
	name := strings.Replace(file, ext, "", 1)
	nfoPath := dir + name + ".nfo"
	ok, err := os.Stat(nfoPath)
	if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("error checking if %s exists: %v", nfoPath, err)
	} else if err == nil && ok.IsDir() {
		return "", fmt.Errorf("%s is a directory", nfoPath)
	} else if err != nil && os.IsNotExist(err) {
		return "", fmt.Errorf("%s does not exist", nfoPath)
	}

	return nfoPath, nil
}

func GetFiles(path string) ([]string, int, error) {
	//to store paths
	var files []string
	//go get em
	err := filepath.Walk(path,
		// Helper function to collect paths
		func(path string, info os.FileInfo, err error) error {
			log.Debug().Msgf("Checking %s", path)
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			} else if strings.HasSuffix(info.Name(), ".avi") ||
				strings.HasSuffix(info.Name(), ".mp4") ||
				strings.HasSuffix(info.Name(), ".mkv") ||
				strings.HasSuffix(info.Name(), ".mpg") ||
				strings.HasSuffix(info.Name(), ".mov") ||
				strings.HasSuffix(info.Name(), ".wmv") ||
				strings.HasSuffix(info.Name(), ".flv") ||
				strings.HasSuffix(info.Name(), ".m4v") {
				files = append(files, path)
				log.Debug().Msgf("Found file: %s", path)
			}
			return nil
		})

	//check for errors
	if err != nil {
		return nil, 0, err
	}
	log.Debug().Msgf("Found %d files", len(files))
	return files, len(files), nil
}
