package youtube

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rashad-j/tiktokuploader/configs"
	"github.com/rashad-j/tiktokuploader/tiktok"
	logger "github.com/rs/zerolog/log"
	"google.golang.org/api/youtube/v3"
)

func UploadShort(service youtube.Service, video tiktok.TiktokVideo, filepath string) (err error) {
	// check if the video is a valid short video
	if video.Duration > 59 {
		logger.Warn().Float64("duration", video.Duration).
			Msg("skipping short video, longer than 60s")
		return nil
	}

	filename := filepath // "videos/food/0.mp4"
	title := video.Title
	description := video.Description
	category := "22"         // see https://gist.github.com/dgp/1b24bf2961521bd75d6c
	keywords := "motivation" // Comma separated list of video keywords
	privacy := "public"

	if filename == "" {
		err = errors.New("no video path assigned")
		return
	}

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       title,
			Description: description,
			CategoryId:  category,
		},
		Status: &youtube.VideoStatus{PrivacyStatus: privacy},
	}

	// The API returns a 400 Bad Request response if tags is an empty string.
	if strings.Trim(keywords, "") != "" {
		upload.Snippet.Tags = strings.Split(keywords, ",")
	}

	call := service.Videos.Insert([]string{"snippet", "status"}, upload)

	file, err := os.Open(filename)
	if err != nil {
		err = fmt.Errorf("error opening %v: %v", filename, err)
		return
	}
	defer file.Close()

	response, err := call.Media(file).Do()
	if err != nil {
		return
	}

	logger.Info().
		Str("response ID", response.Id).
		Str("file", filename).
		Msg("successfully uploaded video")

	return
}

func ChannelsListByUsername(service *youtube.Service, part []string, forUsername string) {
	call := service.Channels.List(part)
	call = call.ForUsername(forUsername)
	response, err := call.Do()
	if err == nil {
		fmt.Println(err)
	}
	fmt.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s', "+
		"and it has %d views.",
		response.Items[0].Id,
		response.Items[0].Snippet.Title,
		response.Items[0].Statistics.ViewCount))
}

func UploadVideo(tv []tiktok.TiktokVideo, cfg configs.Config) error {
	return nil
}
