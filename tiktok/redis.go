package tiktok

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/rashad-j/tiktokuploader/configs"
	logger "github.com/rs/zerolog/log"
)

type RedisClient struct {
	rdb *redis.Client
}

// NewRedisClient Returns a connected RedisClient with the properties of the given configuration.
func NewRedisClient(ctx context.Context, redisURL string) (*RedisClient, error) {
	u, err := url.Parse(redisURL)
	if err != nil {
		return &RedisClient{}, err
	}

	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return &RedisClient{}, err
	}
	if len(m) == 0 {
		return &RedisClient{}, fmt.Errorf("no redis parameters found")
	}

	db, ok := m["db"]
	if !ok {
		return &RedisClient{}, fmt.Errorf("no redis db param found")
	}
	if _, err := strconv.Atoi(db[0]); err != nil {
		return &RedisClient{}, err
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return &RedisClient{}, err
	}

	rdb := redis.NewClient(opt)

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	if pong != "PONG" {
		return nil, fmt.Errorf("redis ping failed: %s", pong)
	}

	logger.Debug().Str("Redis", redisURL).Msg("successfully connected")

	return &RedisClient{rdb: rdb}, nil
}

func (c RedisClient) CheckVideo(cfg configs.Config, key string, tkv TiktokVideo) (bool, error) {
	// Save the new TiktokVideo object
	updated, err := c.saveTiktokVideo(cfg, tkv, key)
	return updated, err
}

func (c RedisClient) saveTiktokVideo(cfg configs.Config, video TiktokVideo, key string) (bool, error) {
	// Check if the video is already saved
	profileVideos, err := c.getTiktokVideos(c.rdb, key)
	if err != nil {
		return false, err
	}

	if len(profileVideos) != 0 {
		for _, v := range profileVideos {
			if v.VideoID == video.VideoID {
				return false, nil
			}
		}
	}

	// Append the new video to the profileVideos list
	profileVideos = append(profileVideos, video)

	// Remove the oldest videos
	// limit the videos. Save only the last 10
	maxVideos := 10
	if len(profileVideos) > maxVideos {
		// Sort the videos by createdAt timestamp in ascending order
		sort.Slice(profileVideos, func(i, j int) bool {
			return profileVideos[i].CreateTime < profileVideos[j].CreateTime
		})
		profileVideos = profileVideos[:maxVideos]
	}

	// Save the updated profileVideos list to Redis
	data, err := json.Marshal(profileVideos)
	if err != nil {
		return false, err
	}
	err = c.rdb.Set(context.Background(), key, data, 0).Err()
	if err != nil {
		return false, err
	}

	logger.Info().Str("key", key).Msg("profile videos updated in redis")

	return true, nil
}

func (c RedisClient) getTiktokVideos(rdb *redis.Client, key string) ([]TiktokVideo, error) {
	// Get all TiktokVideo objects for the given topic and profile
	// Retrieve the profileVideos data from Redis
	data, err := c.rdb.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			// Handle case where data is not found in Redis
			return []TiktokVideo{}, nil
		}
		return nil, err
	}

	// Unmarshal the JSON-encoded data into a []TiktokVideo
	var profileVideos []TiktokVideo
	err = json.Unmarshal([]byte(data), &profileVideos)
	if err != nil {
		return nil, err
	}

	logger.Info().Str("key", key).Int("videos number", len(profileVideos)).Msg("retrieved profiles videos from redis")

	return profileVideos, nil

}

func flushDB(rdb *redis.Client) error {
	_, err := rdb.FlushDB(context.Background()).Result()
	if err != nil {
		return err
	}
	return nil
}

// TODO
// create a go routine to clean up the database every few days
