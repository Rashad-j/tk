package transload

import (
	"context"
	"errors"
	"fmt"

	"github.com/rashad-j/tiktokuploader/configs"
	"github.com/rashad-j/tiktokuploader/tiktok"
	logger "github.com/rs/zerolog/log"
	transloadit "github.com/transloadit/go-sdk"
)

var ErrAssembly error = errors.New("failed to complete the assembly")

// UploadYtURL uploads short or long youtube video. The youtube will decide on the type from the length and video dimensions
func UploadYtURL(video tiktok.TiktokVideo, cfg configs.Config, channel tiktok.Channel) error {
	options := transloadit.DefaultConfig
	options.AuthKey = cfg.TransloaditAuthKey
	options.AuthSecret = cfg.TransloaditAuthSecret
	client := transloadit.NewClient(options)
	// Initialize new assembly
	assembly := transloadit.NewAssembly()
	// get the video
	assembly.AddStep("imported", map[string]interface{}{
		"robot":  "/http/import",
		"result": true,
		"url":    video.PlayURL,
	})
	// checks for the Title & Description
	if video.Title == "" {
		return errors.New("no title found, skipping this video")
	}
	// set description
	video.Description = video.Title
	// check if the video description valid
	if len(video.Description) > int(cfg.YoutubeDescLength) {
		video.Description = video.Description[:cfg.YoutubeDescLength]
	}
	// check if title is valid
	if len(video.Title) > int(cfg.YoutubeTitleLength) {
		video.Title = video.Title[:cfg.YoutubeTitleLength]
	}
	// check if its long video, if so use Chat GPT to generate description
	if channel.TransloaditCreds == "mindset-game" && video.Duration > 60 {
		query := fmt.Sprintf("By using this youtube video title:'%s', generate an advice in 2-3 sentences, make it motivating, exciting and engaging. Don't use phrases like 'Hey there!'. Make the description talking to people in 2nd person. Also add a few quotes that related to the title at the end. Request them to subscribe to this channel, which is named Mindset Game", video.Title)
		description, err := tiktok.AskGPT(cfg, query)
		if err == nil {
			video.Description = description
		} else {
			logger.Err(fmt.Errorf("generating description failed: %w", err))
		}
	}
	// upload the video
	assembly.AddStep("uploaded", map[string]interface{}{
		"robot":       "/youtube/store",
		"credentials": channel.TransloaditCreds,
		"use":         "imported",
		"title":       video.Title,
		"description": video.Description,
		"category":    channel.Category,
		"visibility":  "public",
		"keywords":    channel.Keywords, // TODO ask ai to generate a list of keywords
	})

	info, err := client.StartAssembly(context.Background(), assembly)
	if err != nil {
		return err
	}

	info, err = client.WaitForAssembly(context.Background(), info)
	if err != nil {
		return err
	}

	if info.Ok != "ASSEMBLY_COMPLETED" {
		logger.Err(errors.New(info.Error)).
			Str("transloadit", info.Message).
			Str("title", video.Title).
			Msg("failed to upload")
		return ErrAssembly
	}

	logger.Info().
		Str("transloadit", info.Message).
		Str("title", video.Title).
		Str("desc", video.Description).
		Float64("duration", video.Duration).
		Msg("youtube short video uploaded")

	return nil
}

func UploadYtVideo(video tiktok.TiktokVideo, cfg configs.Config, channel tiktok.Channel, file string) error {
	options := transloadit.DefaultConfig
	options.AuthKey = cfg.TransloaditAuthKey
	options.AuthSecret = cfg.TransloaditAuthSecret
	client := transloadit.NewClient(options)
	// Initialize new assembly
	assembly := transloadit.NewAssembly()
	// get the video
	if err := assembly.AddFile("imported", file); err != nil {
		return err
	}
	// checks for the Title & Description
	if video.Title == "" {
		return errors.New("no title found, skipping this video")
	}
	// set description
	video.Description = video.Title
	// check if the video description valid
	if len(video.Description) > int(cfg.YoutubeDescLength) {
		video.Description = video.Description[:cfg.YoutubeDescLength]
	}
	// check if title is valid
	if len(video.Title) > int(cfg.YoutubeTitleLength) {
		video.Title = video.Title[:cfg.YoutubeTitleLength]
	}
	// upload the video
	assembly.AddStep("uploaded", map[string]interface{}{
		"robot":       "/youtube/store",
		"credentials": channel.TransloaditCreds,
		"use":         ":original",
		"title":       video.Title,
		"description": video.Description,
		"category":    channel.Category,
		"visibility":  "public",
		"keywords":    channel.Keywords, // TODO ask ai to generate a list of keywords
	})

	info, err := client.StartAssembly(context.Background(), assembly)
	if err != nil {
		return err
	}

	info, err = client.WaitForAssembly(context.Background(), info)
	if err != nil {
		return err
	}

	if info.Ok != "ASSEMBLY_COMPLETED" {
		logger.Err(errors.New(info.Error)).
			Str("transloadit", info.Message).
			Str("title", video.Title).
			Msg("failed to upload")
		return ErrAssembly
	}

	logger.Info().
		Str("transloadit", info.Message).
		Str("title", video.Title).
		Str("desc", video.Description).
		Float64("duration", video.Duration).
		Msg("youtube short video uploaded")

	return nil
}

func UploadYoutubeVideo(tv []tiktok.TiktokVideo, cfg configs.Config) error {
	// Create client
	options := transloadit.DefaultConfig
	options.AuthKey = cfg.TransloaditAuthKey
	options.AuthSecret = cfg.TransloaditAuthSecret
	client := transloadit.NewClient(options)
	// Initialize new assembly
	assembly := transloadit.NewAssembly()
	// the encodedVideos that we will concatenate
	var encodedVideos []map[string]interface{}
	// import all videos, give them names
	for index, tv := range tv {
		name := fmt.Sprintf("import_%d", index+1)
		encodeName := fmt.Sprintf("%s_encoded", name)
		as := fmt.Sprintf("video_%d", index+1)
		// import
		assembly.AddStep(name, map[string]interface{}{
			"robot":  "/http/import",
			"result": true,
			"url":    tv.PlayURL,
		})
		// change hight and width to match
		assembly.AddStep(encodeName, map[string]interface{}{
			"robot":        "/video/encode",
			"use":          name,
			"result":       true,
			"width":        720,
			"height":       1280,
			"ffmpeg_stack": "v4.3.1",
			"preset":       "ipad-high",
			// "resize_strategy": "fillcrop",
		})
		// add the video for later use
		v := map[string]interface{}{
			"name": encodeName,
			"as":   as,
		}
		encodedVideos = append(encodedVideos, v)
	}

	// Add a step to concatenate all videos
	assembly.AddStep("concatenated", map[string]interface{}{
		"robot":  "/video/concat",
		"preset": "ipad-high",
		"use": map[string]interface{}{
			"steps":        encodedVideos,
			"bundle_steps": true,
		},
	})

	// get title a desc //TODO do a better approach
	// rand.Seed(time.Now().UnixNano()) // seed the random number generator with the current time
	// index := rand.Intn(len(tv))      // generate a random number between 0 and len(list)
	title := tv[0].Title // TODO read the todo about ideas for titles

	desc := ""
	for _, video := range tv {
		desc = desc + " " + video.Description
	}
	// check if the title is valid
	if len(title) > int(cfg.YoutubeTitleLength) {
		title = title[:cfg.YoutubeTitleLength]
	}
	// check if the video description valid
	if len(desc) > int(cfg.YoutubeDescLength) {
		desc = desc[:cfg.YoutubeDescLength]
	}

	// upload to youtube
	assembly.AddStep("uploaded", map[string]interface{}{
		"robot":       "/youtube/store",
		"credentials": cfg.TransloaditYoutubeCreds,
		"use":         "concatenated",
		"title":       title,
		"description": desc,
		"category":    "entertainment",
		"visibility":  "public",
		"keywords":    "food",
	})

	info, err := client.StartAssembly(context.Background(), assembly)
	if err != nil {
		return err
	}

	info, err = client.WaitForAssembly(context.Background(), info)
	if err != nil {
		return err
	}

	logger.Info().
		Str("transloadit", info.Message).
		Str("title", title).
		Str("desc", desc).
		Msg("youtube video uploaded")

	return nil
}

func Subtitle(cfg configs.Config, file string) (string, error) {
	options := transloadit.DefaultConfig
	options.AuthKey = cfg.TransloaditAuthKey
	options.AuthSecret = cfg.TransloaditAuthSecret
	client := transloadit.NewClient(options)
	// Initialize new assembly
	assembly := transloadit.NewAssembly()

	if err := assembly.AddFile("audio", file); err != nil {
		return "", err
	}

	assembly.AddStep("transcribed", map[string]interface{}{
		"use":      ":original",
		"robot":    "/speech/transcribe",
		"provider": "aws",
		"format":   "srt",
		"result":   true,
	})
	info, err := client.StartAssembly(context.Background(), assembly)
	if err != nil {
		return "", err
	}

	info, err = client.WaitForAssembly(context.Background(), info)
	if err != nil {
		return "", err
	}

	if info.Ok != "ASSEMBLY_COMPLETED" {
		logger.Err(errors.New(info.Error)).
			Str("transloadit", info.Message).
			Msg("failed to upload")
		return "", ErrAssembly
	}

	// fmt.Printf("%+v", info)

	result, ok := info.Results["transcribed"]
	if !ok {
		return "", errors.New("failed to get transcription url")
	}
	url := result[0].SSLURL
	return url, err
	// Get the URL from the response (assume `url` is already obtained)
	// resp, err := http.Get(url)
	// if err != nil {
	// 	return err
	// }
	// defer resp.Body.Close()

	// // Create a file to save the transcription
	// file, err := os.Create(fileName)
	// if err != nil {
	// 	return err
	// }
	// defer file.Close()

	// // Copy the response body to the file
	// _, err = io.Copy(file, resp.Body)
	// if err != nil {
	// 	return err
	// }

}
