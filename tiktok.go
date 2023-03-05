package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/gocrawl"
)

// Create the Extender implementation, based on the gocrawl-provided DefaultExtender,
// because we don't want/need to override all methods.
type TiktokExtender struct {
	gocrawl.DefaultExtender // Will use the default implementation of all but Visit and Filter
	baseURL                 string
	originalURLS            []string
	publicURLS              []string
	maxURLS                 int8
}

// Override Filter for our need.
func (t *TiktokExtender) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	if isVisited {
		return false
	}
	url := fmt.Sprintf("%s", ctx.NormalizedURL().String())
	if strings.Contains(url, "video") {
		if t.maxURLS != 0 {
			t.maxURLS--
			t.originalURLS = append(t.originalURLS, url)
		}
	}
	return true
}

func TiktokCrawl(tiktokExtender *TiktokExtender) error {
	log.Println("crawl started")
	// Set custom options
	opts := gocrawl.NewOptions(tiktokExtender)
	opts.RobotUserAgent = "Crawler"
	opts.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/110.0"
	opts.CrawlDelay = 2 * time.Second
	opts.LogFlags = gocrawl.LogAll
	opts.MaxVisits = 2

	// Create crawler and start at root of duckduckgo
	c := gocrawl.NewCrawlerWithOptions(opts)
	if err := c.Run(tiktokExtender.baseURL); err != nil {
		return err
	}

	log.Println("crawl finished")
	return nil
}

func publicURL(orignialURL string) (string, error) {
	log.Println("getting public url for", orignialURL)

	url := fmt.Sprintf("https://tiktok-video-no-watermark2.p.rapidapi.com/?url=%s", orignialURL)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "596fd604f9msh9205a55cc027e9ap1bdb1ejsn79730d1f6b79")
	req.Header.Add("X-RapidAPI-Host", "tiktok-video-no-watermark2.p.rapidapi.com")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	publicURL := data["data"].(map[string]interface{})["play"].(string)

	return publicURL, nil
}

func download(filepath string, name string, url string) error {
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

	log.Println("downloaded file", filepath)

	return nil
}

func uploadToYoutube() {

}

func uploadToInstagram() {

}
