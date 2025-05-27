package utils

import (
	"fmt"
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
	name := strings.Replace(file, "."+ext, "", 1)
	newPath := dir + "/" + name + "." + tag + "." + ext
	return newPath
}

// NFOPath returns the path to the nfo file for a given video file if it exists
// assumes that nfo is a sidecar with the same name /foo/bar/my_video.avi & /foo/bar/my_video.nfo
func NFOPath(path string) (string, error) {
	//get working directory + name
	dir, file := filepath.Split(path)
	ext := filepath.Ext(path)
	name := strings.Replace(file, "."+ext, "", 1)
	nfoPath := dir + "/" + name + ".nfo"
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
