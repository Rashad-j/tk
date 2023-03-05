package main

import (
	"io/ioutil"
	"strings"
)

func profiles() (map[string][]string, error) {
	// specify the directory path
	dirPath := "profiles"
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
			// split the content into lines
			lines := strings.Split(string(content), "\n")
			// add the lines to the map
			filename := file.Name()
			filename = filename[:len(filename)-4]
			profiles[filename] = lines
		}
	}
	// print the result
	return profiles, nil
}
