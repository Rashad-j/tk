package main

import (
	"fmt"
	"log"
)

func main() {

	const baseVideosPath = "videos"

	profiles, err := profiles()
	if err != nil {
		log.Println(err)
	}

	for topic, urls := range profiles {
		topicVideos := make([]string, 0)
		for _, url := range urls {
			te := TiktokExtender{
				baseURL: url,
				maxURLS: 5,
			}

			if err := TiktokCrawl(&te); err != nil {
				fmt.Println(err)
			}
			topicVideos = append(topicVideos, te.originalURLS...)
		}
		// publish this topic
		fmt.Println("topic initialized", topic, topicVideos)
		for index, video := range topicVideos {
			// get public url
			publicURL, err := publicURL(video)
			if err != nil {
				log.Println(err)
			}
			filepath := fmt.Sprintf("%s/%s", baseVideosPath, topic)
			name := fmt.Sprintf("%d.mp4", index)
			// download the video
			if err := download(filepath, name, publicURL); err != nil {
				log.Println(err)
			}
			break // TODO remove
		}
	}

	// url, err := publicURL(te.originalURLS[0])
	// if err != nil {
	// 	log.Println("error in main", err)
	// }

	// filepath := "videos/video.mp4" // Replace with the desired file path and name for the saved video
	// err = DownloadFile(filepath, url)
	// if err != nil {
	// 	fmt.Println("Error downloading file:", err)
	// 	return
	// }
	// fmt.Println("Video downloaded successfully")

	// te.printURLS()
	// transloaditCheck()
	// getDownloadURL()
	// getDownloadURLv1()
}
