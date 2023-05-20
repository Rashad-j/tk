package configs

import (
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
)

type Config struct {
	// logger
	ZerologLevel zerolog.Level `env:"ZEROLOG_LEVEL" envDefault:"0"` // DEBUG 0

	// crawler
	CrawlerVideosPerProfile int           `env:"CRAWLER_VIDEOS_PER_PROFILE" envDefault:"1"`
	CrawlerRobotUserAgent   string        `env:"CRAWLER_ROBOT_USER_AGENT" envDefault:"Crawler"`
	CrawlerUserAgent        string        `env:"CRAWLER_USER_AGENT" envDefault:"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/110.0"`
	CrawlerDelay            time.Duration `env:"CRAWLER_DELAY" envDefault:"2s"`
	CrawlerMaxVisits        int           `env:"CRAWLER_MAX_VISITS" envDefault:"2"`

	// rapid API
	RapidAPIKey  string `env:"RAPID_API_KEY" envDefault:"key-string"`
	RapidAPIHost string `env:"RAPID_API_HOST" envDefault:"host-string"`

	// transloadit
	TransloaditAuthKey      string `env:"TRANSLOADIT_AUTH_KEY" envDefault:"KEY"`
	TransloaditAuthSecret   string `env:"TRANSLOADIT_AUTH_SECRET" envDefault:"SEC"`
	TransloaditYoutubeCreds string `env:"TRANSLOADIT_YOUTUBE_CREDS" envDefault:"CRED"`

	// youtube
	YoutubeTitleLength  int    `env:"YOUTUBE_TITLE_LENGTH" envDefault:"99"`
	YoutubeDescLength   int    `env:"YOUTUBE_DESC_LENGTH" envDefault:"5000"`
	YoutubeClientSecret string `env:"YOUTUBE_CLIENT_SECRET" envDefault:"youtube/client_secret.json"`

	// tiktok
	TiktokProfileDir       string `env:"TIKTOK_PROFILE_DIR" envDefault:"tiktok/profiles"`
	TiktokCacheDir         string `env:"TIKTOK_CACHE_DIR" envDefault:"tiktok/cache/cache.json"`
	TiktokVideosPerProfile int    `env:"TIKTOK_VIDEOS_PER_PROFILE" envDefault:"20"`

	// chatGPT
	ChatGptApiKey string `env:"CHAT_GPT_API_KEY" envDefault:"api-key"`

	// others
	VideoDownloadDir string        `env:"VIDEO_DOWNLOAD_DIR" envDefault:"tiktok/videos"`
	UploadDelay      time.Duration `env:"UPLOAD_DELAY" envDefault:"5s"`
	RedisURL         string        `env:"REDIS_URL" envDefault:"redis://localhost:6379?db=3"` // db is the database index
	InstaUploadURL   string        `env:"INSTA_UPLOAD_URL" envDefault:"http://localhost:3000/upload"`
}

func LoadConfigs() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	// TODO check if the configs were loaded successfully
	logger.Info().
		//Interface("config", cfg).
		Msg("loaded config")
	return cfg, nil
}
