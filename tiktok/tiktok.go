package tiktok

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/rashad-j/tiktokuploader/configs"
	"github.com/rashad-j/tiktokuploader/utils"

	"github.com/rs/zerolog/log"
	logger "github.com/rs/zerolog/log"
)

type TiktokUser struct {
	ID                  string `json:"id"`
	Region              string `json:"region"`
	UniqueID            string `json:"unique_id"`
	Nickname            string `json:"nickname"`
	AwemeCount          int    `json:"aweme_count"`
	FollowerCount       int    `json:"follower_count"`
	FavoritingCount     int    `json:"favoriting_count"`
	TotalFavorited      int    `json:"total_favorited"`
	YoutubeChannelID    string `json:"youtube_channel_id"`
	YoutubeChannelTitle string `json:"youtube_channel_title"`
}

type FollowingListResponse struct {
	Code          int     `json:"code"`
	Message       string  `json:"msg"`
	ProcessedTime float64 `json:"processed_time"`
	Data          struct {
		Followings []TiktokUser `json:"followings"`
	} `json:"data"`
}
type TiktokVideo struct {
	VideoID       string    `json:"video_id"`
	Region        string    `json:"region"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Cover         string    `json:"cover"`
	Duration      float64   `json:"duration"`
	PlayURL       string    `json:"play"`
	PlayCount     int       `json:"play_count"`
	CommentCount  int       `json:"comment_count"`
	ShareCount    int       `json:"share_count"`
	DownloadCount int       `json:"download_count"`
	ForwardCount  int       `json:"forward_count"`
	CreateTime    int       `json:"create_time"`
	Music         string    `json:"music"`
	MusicInfo     MusicInfo `json:"music_info"`
}
type MusicInfo struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Play   string `json:"play"`
	Cover  string `json:"cover"`
	Author string `json:"author"`
}
type VideosListResponse struct {
	Code          int     `json:"code"`
	Message       string  `json:"msg"`
	ProcessedTime float64 `json:"processed_time"`
	Data          struct {
		Videos []TiktokVideo `json:"videos"`
	} `json:"data"`
}

// channels
type Channel struct {
	Title            string
	TransloaditCreds string
	TiktokUniqueID   string
	TiktokUserId     int
	FollowingsNumber int
	Keywords         string
	Category         string
	InstaID          string
}

func FollowingList(cfg configs.Config, channel Channel) (FollowingListResponse, error) {
	// get all accounts that I am following with Mindset.game account
	userID := channel.TiktokUserId
	count := channel.FollowingsNumber
	time := 0
	url := fmt.Sprintf("https://tiktok-video-no-watermark2.p.rapidapi.com/user/following?user_id=%d&count=%d&time=%d", userID, count, time)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-RapidAPI-Key", cfg.RapidAPIKey)
	req.Header.Add("X-RapidAPI-Host", cfg.RapidAPIHost)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return FollowingListResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return FollowingListResponse{}, fmt.Errorf("non 200 returned: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return FollowingListResponse{}, err
	}

	response := FollowingListResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return FollowingListResponse{}, err
	}

	// fmt.Println(string(body))

	logger.Info().Int("len", len(response.Data.Followings)).Msg("fetched the list of followings")

	return response, nil
}

func Posts(cfg configs.Config, tku TiktokUser) (VideosListResponse, error) {
	// from the following list I get two things for each user: unique_id and user_id
	uniqueID := tku.UniqueID
	userID := tku.ID
	if uniqueID == "" || userID == "" {
		return VideosListResponse{}, errors.New("empty unique/user ID")
	}

	url := fmt.Sprintf("https://tiktok-video-no-watermark2.p.rapidapi.com/user/posts?unique_id=%s&user_id=%s&count=%d&cursor=0", uniqueID, userID, cfg.TiktokVideosPerProfile)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-RapidAPI-Key", cfg.RapidAPIKey)
	req.Header.Add("X-RapidAPI-Host", cfg.RapidAPIHost)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return VideosListResponse{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return VideosListResponse{}, fmt.Errorf("non 200 returned: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return VideosListResponse{}, err
	}

	response := VideosListResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return VideosListResponse{}, err
	}

	logger.Info().Int("number of videos", len(response.Data.Videos)).Msg("user videos fetched")

	return response, nil
}

func Audios(videosResponse VideosListResponse, filepath string) error {
	lastAudioNumber, err := getLastAudioNumber(filepath)
	if err != nil {
		return err
	}
	lastAudioNumber++

	videos := videosResponse.Data.Videos
	for index, video := range videos {
		musicUrl := video.Music
		name := fmt.Sprintf("%d.mp3", index+lastAudioNumber)
		err := utils.Download(filepath, name, musicUrl)
		if err != nil {
			log.Err(err).Interface("video", video).Msg("failed to download music")
		}
		log.Info().Str("music", musicUrl).Str("filepath", filepath).Msg("download success")
	}
	return nil
}

func getLastAudioNumber(searchPath string) (int, error) {
	files, err := ioutil.ReadDir(searchPath)
	if err != nil {
		return 0, err
	}

	max := 0
	for _, f := range files {
		name := f.Name()
		if !strings.HasSuffix(name, ".mp3") {
			continue
		}

		name = strings.TrimSuffix(name, ".mp3")
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
		return 0, nil
	}

	fmt.Printf("The highest video name is %d.mp3", max)
	return max, nil
}
