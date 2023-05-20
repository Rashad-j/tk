package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/rashad-j/tiktokuploader/configs"
	logger "github.com/rs/zerolog/log"
)

func Profiles(cfg configs.Config) (map[string][]string, error) {
	// specify the directory path
	dirPath := cfg.TiktokProfileDir
	// create an empty map to hold the data
	profiles := make(map[string][]string)
	// read the directory contents
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	// loop through each file
	for _, file := range files {
		// check if the file is a regular file (not a directory or symlink)
		if file.Mode().IsRegular() {
			// read the file content
			content, err := ioutil.ReadFile(dirPath + "/" + file.Name())
			if err != nil {
				return nil, err
			}
			// split the content into allLines
			allLines := strings.Split(string(content), "\n")
			filteredLines := make([]string, 0)
			for _, line := range allLines {
				if strings.HasPrefix(line, "#") {
					continue // ignore line
				}
				// process non-comment line
				filteredLines = append(filteredLines, line)
			}
			// add the lines to the map
			filename := file.Name()
			filename = filename[:len(filename)-4]
			profiles[filename] = filteredLines
		}
	}
	// print the result
	return profiles, nil
}

func Download(filepath string, name string, url string) error {
	err := os.MkdirAll(filepath, os.ModePerm)
	if err != nil {
		return err
	}
	// Create the file
	out, err := os.Create(fmt.Sprintf("%s/%s", filepath, name))
	if err != nil {
		return err
	}
	defer out.Close()

	// Make the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the response body to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	logger.Info().
		Str("path", filepath).
		Str("name", name).
		Msg("file downloaded")

	return nil
}
