package crawler

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/rashad-j/tiktokuploader/configs"
)

type TiktokExtender struct {
	gocrawl.DefaultExtender // Will use the default implementation of all but Visit and Filter
	ProfileURL              string
	Videos                  []string
	MaxVideos               int
}

func (t *TiktokExtender) TiktokCrawl(cfg configs.Config) error {
	// Set custom options
	opts := gocrawl.NewOptions(t)
	opts.RobotUserAgent = cfg.CrawlerRobotUserAgent
	opts.UserAgent = cfg.CrawlerUserAgent
	opts.CrawlDelay = cfg.CrawlerDelay
	opts.LogFlags = gocrawl.LogAll
	opts.MaxVisits = cfg.CrawlerMaxVisits

	// Create crawler and start at root of the profile
	c := gocrawl.NewCrawlerWithOptions(opts)
	if err := c.Run(t.ProfileURL); err != nil {
		return err
	}

	return nil
}

func (t *TiktokExtender) Log(logFlags gocrawl.LogFlags, msgLevel gocrawl.LogFlags, msg string) {
	if logFlags&msgLevel == msgLevel {
		// log.Println(msg)
	}
}

// Filter overrides the crawler Filter to our need
func (t *TiktokExtender) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	if isVisited {
		return false
	}
	url := fmt.Sprintf("%s", ctx.NormalizedURL().String())
	if strings.Contains(url, "video") {
		if t.MaxVideos != 0 {
			t.MaxVideos--
			t.Videos = append(t.Videos, url)
		}
	}
	return true
}
