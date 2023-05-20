package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/rashad-j/tiktokuploader/configs"
	"github.com/rashad-j/tiktokuploader/ffmpeg"
	"github.com/rashad-j/tiktokuploader/instagram"
	"github.com/rashad-j/tiktokuploader/tiktok"
	transloadit "github.com/rashad-j/tiktokuploader/transloadit"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// load the config
	cfg, err := configs.LoadConfigs()
	if err != nil {
		panic(err)
	}
	// url, err := transloadit.Subtitle(cfg, "media/quote/audio6.flac")
	// fmt.Println(url, err)
	// res, err := chatgpt.Ask(cfg.ChatGptApiKey, "generate a youtube video title for the following quote: How much more grievous are the consequences of anger than the causes of it.")
	// if err != nil {
	// 	panic(err)
	// }
	// text := strings.ReplaceAll(res.Choices[0].Text, "\n", "")
	// text = strings.ReplaceAll(text, "\"", "")
	// fmt.Println(text)
	// log.Info().Interface("res", res).Msg("")
	// os.Exit(0)
	// testing
	// ffmpeg.Run(cfg)
	// os.Exit(0)

	// res, _ := pexels.Search("landscape")
	// pexels.Save(res)
	// os.Exit(0)

	// resp, _ := pexels.SearchPhoto("galaxy")
	// if err := pexels.SavePhotos(resp); err != nil {
	// 	fmt.Println(err)
	// }
	// os.Exit(0)

	// if err := transloadit.Subtitle(cfg); err != nil {
	// 	fmt.Println(err)
	// }
	// os.Exit(0)

	// youtubeAPI.ChannelsListByUsername(service, []string{"snippet", "contentDetails", "statistics"}, "GoogleDevelopers")

	// tiktok
	// channel := tiktok.Channel{
	// 	Title:            "mindsetgame",
	// 	TransloaditCreds: "mindset-game",
	// 	TiktokUniqueID:   "mindset.game",
	// 	TiktokUserId:     7019876592412001285,
	// 	FollowingsNumber: 50,
	// 	Keywords:         "Motivation,Inspiration,Self-improvement,Success,Goal-setting,Mindset,Positive thinking,Productivity,Leadership,Confidence,Gratitude,Personal growth,Mindfulness,Habits,Time management", // TODO
	// 	Category:         "education",
	// 	InstaID:          "mindset.game2",
	// }
	// tiktok.FollowingList(cfg, channel)
	// os.Exit(0)

	// prepare logger
	loglevel, err := zerolog.ParseLevel("debug")
	if err != nil {
		loglevel = zerolog.InfoLevel
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		NoColor:    false,
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Logger().Level(loglevel)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// A) tiktok
	// tkUser := tiktok.TiktokUser{
	// 	ID:       "6816381254813828102",
	// 	UniqueID: "team.nevergiveup7",
	// }
	// // 1. Get the recent posts
	// tkPosts, err := tiktok.Posts(cfg, tkUser)
	// if err != nil {
	// 	log.Err(err).Interface("user", tkUser).Interface("posts", tkPosts).Msg("failed to get user posts")
	// }
	// // 2. Download the sounds
	// filepath := "./media/tiktok"
	// if err := tiktok.Audios(tkPosts, filepath); err != nil {
	// 	log.Err(err).Interface("user", tkUser).Interface("videos", tkPosts).Msg("failed to download audios")
	// }
	// os.Exit(0)

	// B) Pexels videos
	// resp, err := pexels.Search("city night rain") // winter snow nights, mountains snow storm
	// if err != nil {
	// 	log.Error().Err(err).Msg("failed to get response from pexels")
	// 	os.Exit(1)
	// }
	// // sleep a bit to avoid: Sorry This video does not exist.
	// time.Sleep(time.Minute * 5)
	// // save the videos
	// pexels.Save(resp)

	// C) quotes
	// this process must be supervised. I will expose an endpoint with a youtube video url, then it will be downloaded and divided
	// this is next TODO

	// D)
	// prepare video, with audios and subtitles
	// ffmpeg.RunLong(cfg)
	ffmpeg.Run(cfg)
	os.Exit(0)

	// redis database connection
	redisClient, err := tiktok.NewRedisClient(context.Background(), cfg.RedisURL)
	if err != nil {
		log.Err(err).Msg("failed to create redis client")
		panic(err)
	}

	// shut down context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	channels := []tiktok.Channel{
		{
			Title:            "foodlove",
			TransloaditCreds: "foodlove",
			TiktokUniqueID:   "food.love23",
			TiktokUserId:     7211873707819910149,
			FollowingsNumber: 50,                                                                                                                                          // max number of followings returned
			Keywords:         "Food,Cooking,Recipes,Kitchen,Baking,Healthy eating,Meal prep,Nutrition,Ingredients,Chef,Culinary,Gourmet,Appetizers,Desserts,Comfort food", // TODO
			Category:         "entertainment",
			InstaID:          "foodlove.u",
		},
		{
			Title:            "mindsetgame",
			TransloaditCreds: "mindset-game",
			TiktokUniqueID:   "mindset.game",
			TiktokUserId:     7019876592412001285,
			FollowingsNumber: 50,
			Keywords:         "Motivation,Inspiration,Self-improvement,Success,Goal-setting,Mindset,Positive thinking,Productivity,Leadership,Confidence,Gratitude,Personal growth,Mindfulness,Habits,Time management", // TODO
			Category:         "education",
			InstaID:          "mindset.game2",
		},
		{
			Title:            "howto",
			TransloaditCreds: "howto",
			TiktokUniqueID:   "howto.tiktok3",
			TiktokUserId:     7217356375314629638,
			FollowingsNumber: 50,
			Keywords:         "DIY,Fashion,Beauty,Style,Hair,Makeup,Crafting,Sewing,Design,Home decor,Organization,Upcycling,Refashion,Accessories,Jewelry", // TODO
			Category:         "howto & style",
			InstaID:          "everything.diy20", // TODO
		},
	}

	for _, channel := range channels {
		go runInsta(ctx, cfg, redisClient, channel)
		go runTiktok(ctx, cfg, redisClient, channel)
	}

	<-ctx.Done()
	log.Info().Msg("shutting down")
}

func runTiktok(ctx context.Context, cfg configs.Config, redisClient *tiktok.RedisClient, channel tiktok.Channel) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// get all users that I am following
			followings, err := tiktok.FollowingList(cfg, channel)
			if err != nil {
				log.Err(nil).Interface("followings error", followings).Msg("no followings returned from the API")
			}

			if followings.Data.Followings == nil {
				log.Err(nil).Interface("followings", followings).Msg("no followings returned from the API")
			}

			for _, followedUser := range followings.Data.Followings {
				videosResponse, err := tiktok.Posts(cfg, followedUser)
				if err != nil {
					log.Err(err).Interface("following", followedUser).Interface("videos", videosResponse).Msg("failed to get user videos")
				}

				// get the user latest videos
				videos := videosResponse.Data.Videos
				if len(videos) > cfg.TiktokVideosPerProfile {
					videos = videos[:cfg.TiktokVideosPerProfile]
				}
				for _, video := range videos {
					// fmt.Println(video.VideoID, followedUser.UniqueID, video.Duration, video.Title)
					key := fmt.Sprintf("tiktok:%s", followedUser.UniqueID)
					updated, err := redisClient.CheckVideo(cfg, key, video)
					if err != nil {
						log.Err(err).Interface("video", video).Msg("failed to check on redis")
					}
					if updated {
						log.Info().Msg("starting youtube upload...")
						if err := transloadit.UploadYtURL(video, cfg, channel); err != nil {
							log.Err(err).Interface("video", video).
								Msg("failed to upload youtube short")
							if err == transloadit.ErrAssembly {
								// TODO try upload the video again??
								videos = append(videos, video)
								log.Warn().Str("channel", channel.TiktokUniqueID).Msg("upload limit reached, going into a stale period")
								time.Sleep(time.Duration(6) * time.Hour)
							}
						}
						// download
						// filepath := fmt.Sprintf("./tiktok/videos/%s", followedUser.UniqueID)
						// name := fmt.Sprintf("%d.mp4", index)
						// url := video.PlayURL

						// err := utils.Download(filepath, name, url)
						// if err != nil {
						// 	log.Err(err).Interface("video", video).Msg("failed to download video")
						// }
						// start a delay
						randomDelay := rand.Intn(20) + 40
						log.Info().Int("minutes", randomDelay).Msg("tiktok next upload delay running...")
						// sleep for the random delay
						time.Sleep(time.Duration(randomDelay) * time.Minute)
					} else {
						log.Info().Str("user", followedUser.UniqueID).Str("video ID", video.VideoID).Msg("video was uploaded before, skipping")
					}
				}
				// TODO merge all videos for that user and upload them
			}
			// delay next check
			log.Info().Msg("tiktok starting next round delay")
			time.Sleep(time.Duration(3) * time.Hour)
		}
	}

}

func runInsta(ctx context.Context, cfg configs.Config, redisClient *tiktok.RedisClient, channel tiktok.Channel) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// get all users that I am following
			followings, err := tiktok.FollowingList(cfg, channel)
			if err != nil {
				log.Err(nil).Interface("followings error", followings).Msg("no followings returned from the API")
			}

			if followings.Data.Followings == nil {
				log.Err(nil).Interface("followings", followings).Msg("no followings returned from the API")
			}

			for _, followedUser := range followings.Data.Followings {
				videosResponse, err := tiktok.Posts(cfg, followedUser)
				if err != nil {
					log.Err(err).Interface("following", followedUser).Interface("videos", videosResponse).Msg("failed to get user videos")
				}

				// get the user latest videos
				videos := videosResponse.Data.Videos
				if len(videos) > cfg.TiktokVideosPerProfile {
					videos = videos[:cfg.TiktokVideosPerProfile]
				}
				for _, video := range videos {
					// fmt.Println(video.VideoID, followedUser.UniqueID, video.Duration, video.Title)
					key := fmt.Sprintf("insta:%s", followedUser.UniqueID)
					updated, err := redisClient.CheckVideo(cfg, key, video)
					if err != nil {
						log.Err(err).Interface("video", video).Msg("failed to check video on instaRun")
					}
					if updated {
						if channel.InstaID != "" {
							if video.Duration > 59 {
								log.Warn().Float64("length", video.Duration).Msg("insta video length is more than 60 seconds")
								continue
							}
							log.Info().Msg("starting insta upload")
							ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
							ur := instagram.UploadRequest{
								Cover:    video.Cover,
								Video:    video.PlayURL,
								Caption:  video.Title,
								Duration: int(video.Duration),
								Username: channel.InstaID,
							}
							if err := instagram.UploadUrl(ctx, ur, cfg.InstaUploadURL); err != nil {
								log.Err(err).Msg("failed to upload instagram video")
								cancel()
							}
							log.Info().Msg("insta video uploaded")
							cancel()
							// start a delay
							randomDelay := rand.Intn(20) + 40
							log.Info().Int("minutes", randomDelay).Msg("insta next upload delay running...")
							// sleep for the random delay
							time.Sleep(time.Duration(randomDelay) * time.Minute)
						} else {
							log.Info().Str("channel", channel.TiktokUniqueID).Msg("no insta channel found")
						}
					} else {
						log.Info().Str("user", followedUser.UniqueID).Str("video ID", video.VideoID).Msg("video was uploaded before, skipping")
					}
				}
			}
			// delay next check
			log.Info().Msg("insta starting next round delay")
			time.Sleep(time.Duration(3) * time.Hour)
		}
	}

}
