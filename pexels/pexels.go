package pexels

// api key for pexel: grveYmbzwskVD1u2Pnm9fymi7rkLUX4wEz7q9UVD0nKTQwarVRgBQyqh
// https://www.pexels.com/api/documentation/

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/rashad-j/tiktokuploader/utils"
	"github.com/rs/zerolog/log"
)

type VideoFile struct {
	ID       int     `json:"id"`
	Quality  string  `json:"quality"`
	FileType string  `json:"file_type"`
	Width    int     `json:"width"`
	Height   int     `json:"height"`
	FPS      float32 `json:"fps"`
	Link     string  `json:"link"`
}

type VideoPicture struct {
	ID      int    `json:"id"`
	Nr      int    `json:"nr"`
	Picture string `json:"picture"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Video struct {
	ID            int            `json:"id"`
	Width         int            `json:"width"`
	Height        int            `json:"height"`
	Duration      int            `json:"duration"`
	FullRes       interface{}    `json:"full_res"`
	Tags          []interface{}  `json:"tags"`
	URL           string         `json:"url"`
	Image         string         `json:"image"`
	AvgColor      interface{}    `json:"avg_color"`
	User          User           `json:"user"`
	VideoFiles    []VideoFile    `json:"video_files"`
	VideoPictures []VideoPicture `json:"video_pictures"`
}

type VideoResponse struct {
	Page    int     `json:"page"`
	PerPage int     `json:"per_page"`
	Videos  []Video `json:"videos"`
}

type Photo struct {
	ID              int    `json:"id"`
	Width           int    `json:"width"`
	Height          int    `json:"height"`
	URL             string `json:"url"`
	Photographer    string `json:"photographer"`
	PhotographerURL string `json:"photographer_url"`
	PhotographerID  int    `json:"photographer_id"`
	AvgColor        string `json:"avg_color"`
	Src             struct {
		Original  string `json:"original"`
		Large2x   string `json:"large2x"`
		Large     string `json:"large"`
		Medium    string `json:"medium"`
		Small     string `json:"small"`
		Portrait  string `json:"portrait"`
		Landscape string `json:"landscape"`
		Tiny      string `json:"tiny"`
	} `json:"src"`
	Liked bool   `json:"liked"`
	Alt   string `json:"alt"`
}

type PhotoResponse struct {
	TotalResults int     `json:"total_results"`
	Page         int     `json:"page"`
	PerPage      int     `json:"per_page"`
	Photos       []Photo `json:"photos"`
	NextPage     string  `json:"next_page"`
}

func SearchVideo(query string) (VideoResponse, error) {
	orientation := "portrait" // optional parameter
	size := "large"           // optional parameter
	locale := "en-US"         // optional parameter
	page := 1                 // optional parameter
	perPage := 10             // optional parameter

	// Create a new HTTP request with the necessary query parameters
	req, err := http.NewRequest("GET", "https://api.pexels.com/videos/search", nil)
	if err != nil {
		return VideoResponse{}, err
	}

	q := req.URL.Query()
	q.Add("query", query)

	if orientation != "" {
		q.Add("orientation", orientation)
	}
	if size != "" {
		q.Add("size", size)
	}
	if locale != "" {
		q.Add("locale", locale)
	}
	if page > 1 {
		q.Add("page", fmt.Sprintf("%d", page))
	}
	if perPage > 0 {
		q.Add("per_page", fmt.Sprintf("%d", perPage))
	}

	req.URL.RawQuery = q.Encode()

	// Add the necessary headers to the request
	req.Header.Set("Authorization", "grveYmbzwskVD1u2Pnm9fymi7rkLUX4wEz7q9UVD0nKTQwarVRgBQyqh")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Send the request and get the response
	client := http.Client{}
	r, err := client.Do(req)
	if err != nil {
		return VideoResponse{}, err
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return VideoResponse{}, errors.New("non 200 returned, code" + strconv.Itoa(r.StatusCode))
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return VideoResponse{}, err
	}

	var resp VideoResponse
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		return VideoResponse{}, err
	}

	log.Info().Str("query", query).Msg("pexels videos retrieved")

	return resp, nil
}

func SaveVideo(resp VideoResponse) error {
	filepath := "./media/videos"

	lastVideoNumber, err := getLastMediaNumber(filepath, ".mp4")
	if err != nil {
		return err
	}
	lastVideoNumber++

	for i, v := range resp.Videos {
		fmt.Println(v.VideoFiles)
		for _, vl := range v.VideoFiles {
			fmt.Println(vl.Quality, vl.Link)
			if vl.Quality == "hd" {
				name := fmt.Sprintf("%d.mp4", lastVideoNumber+i)
				if err := utils.Download(filepath, name, vl.Link); err != nil {
					log.Error().
						Err(err).
						Str("name", name).
						Str("link", vl.Link).
						Msg("failed to download video from pexel")
				}
				log.Info().Str("link", vl.Link).Msg("pexels video saved")
				break
			}
		}
	}
	log.Info().Msg("pexels videos saved")
	return nil
}

func getLastMediaNumber(searchPath, suffix string) (int, error) {
	files, err := ioutil.ReadDir(searchPath)
	if err != nil {
		return 0, err
	}

	max := 0
	for _, f := range files {
		name := f.Name()
		if !strings.HasSuffix(name, suffix) { // e.g. ".mp4"
			continue
		}

		name = strings.TrimSuffix(name, suffix)
		num, err := strconv.Atoi(name)
		if err != nil {
			log.Printf("Ignoring invalid file name %s", f.Name())
			continue
		}

		if num > max {
			max = num
		}
	}

	if max == 0 {
		return 0, nil // errors.New("No video files found")
	}

	fmt.Printf("The highest file name is %d.%s", max, suffix)
	return max, nil
}

func SearchPhoto(query string) (PhotoResponse, error) {
	orientation := "portrait" // optional parameter
	size := "large"           // optional parameter
	locale := "en-US"         // optional parameter
	page := 1                 // optional parameter
	perPage := 20             // optional parameter

	// Create a new HTTP request with the necessary query parameters
	req, err := http.NewRequest("GET", "https://api.pexels.com/v1/search", nil)
	if err != nil {
		return PhotoResponse{}, err
	}

	q := req.URL.Query()
	q.Add("query", query)

	if orientation != "" {
		q.Add("orientation", orientation)
	}
	if size != "" {
		q.Add("size", size)
	}
	if locale != "" {
		q.Add("locale", locale)
	}
	if page > 1 {
		q.Add("page", fmt.Sprintf("%d", page))
	}
	if perPage > 0 {
		q.Add("per_page", fmt.Sprintf("%d", perPage))
	}

	req.URL.RawQuery = q.Encode()

	// Add the necessary headers to the request
	req.Header.Set("Authorization", "grveYmbzwskVD1u2Pnm9fymi7rkLUX4wEz7q9UVD0nKTQwarVRgBQyqh")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Send the request and get the response
	client := http.Client{}
	r, err := client.Do(req)
	if err != nil {
		return PhotoResponse{}, err
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return PhotoResponse{}, errors.New("non 200 returned, code" + strconv.Itoa(r.StatusCode))
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return PhotoResponse{}, err
	}

	var resp PhotoResponse
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		return PhotoResponse{}, err
	}

	log.Info().Str("query", query).Msg("pexels photos retrieved")

	// log.Info().Interface("photos", resp).Msg("photos")

	return resp, nil
}

func SavePhotos(resp PhotoResponse) error {
	filepath := "./media/photos"
	lastNumber, err := getLastMediaNumber(filepath, ".jpeg")
	if err != nil {
		return err
	}
	lastNumber++

	for i, photo := range resp.Photos {
		log.Info().Interface("src", photo.Src).Msg("src")
		name := fmt.Sprintf("%d.jpeg", lastNumber+i)
		url := photo.Src.Large2x
		if err := utils.Download(filepath, name, url); err != nil {
			log.Error().
				Err(err).
				Str("name", name).
				Str("link", url).
				Msg("failed to download file from pexel")
		}
		log.Info().Str("link", url).Msg("pexels video saved")
	}
	log.Info().Msg("pexels photo saved")
	return nil
}
