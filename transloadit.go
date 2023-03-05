package main

import (
	"context"
	"fmt"

	transloadit "github.com/transloadit/go-sdk"
)

func transloaditCheck() error {
	// Create client
	options := transloadit.DefaultConfig
	options.AuthKey = "fcffdd9e24214d62aac5bda1f465d1f6"
	options.AuthSecret = "f849729daeeea784feb60cd6efd35142872f3edb"
	client := transloadit.NewClient(options)

	// Initialize new assembly
	assembly := transloadit.NewAssembly()

	assembly.AddStep("preroll_imported", map[string]interface{}{
		"robot":  "/http/import",
		"result": true,
		"url":    "https://v39-eu.tiktokcdn.com/654824710daf1d76f17945df3a64fc7d/6403fb59/video/tos/useast2a/tos-useast2a-ve-0068c001/oYjW6FC2o7WMBwMEQOHysBAegR8tFJbemQWn1D/?a=1340&ch=0&cr=3&dr=0&lr=all&cd=0%7C0%7C0%7C3&cv=1&br=3312&bt=1656&cs=0&ds=6&ft=kLey-ymRZNr0PD13d_RXg9wTWXNcBEeC~&mime_type=video_mp4&qs=0&rc=OjQ8Nzo0OjQ8aDk1OWQ7Z0BpMzRnbTQ6Zm1wZzMzNzczM0AuXjE1Yl9fXzYxXzZhMC8wYSMwZG02cjRnbmBgLS1kMTZzcw%3D%3D&l=20230304201532684D657DA7CB747A2247&btag=80000&cc=1e",
	})

	// Add Instructions, e.g. resize image to 75x75px
	assembly.AddStep("preroll_resized", map[string]interface{}{
		"robot":           "/video/encode",
		"use":             "preroll_imported",
		"result":          true,
		"ffmpeg_stack":    "v4.3.1",
		"preset":          "ipad-high",
		"width":           480,
		"height":          270,
		"resize_strategy": "pad",
		"background":      "#000000",
		"turbo":           true,
	})

	// assembly.AddStep("exported", map[string]interface{}{
	// 	"use": [
	// 		"preroll_imported",
	// 		"preroll_resized",
	// 	  ],
	// 	  "robot": "/s3/store",
	// 	  "credentials": "YOUR_AWS_CREDENTIALS",
	// 	  "url_prefix": "https://demos.transloadit.com/"
	// })

	// Start the upload
	info, err := client.StartAssembly(context.Background(), assembly)
	if err != nil {
		panic(err)
	}

	// All files have now been uploaded and the assembly has started but no
	// results are available yet since the conversion has not finished.
	// WaitForAssembly provides functionality for polling until the assembly
	// has ended.
	info, err = client.WaitForAssembly(context.Background(), info)
	if err != nil {
		panic(err)
	}

	fmt.Println(info.Results)
	//fmt.Printf("You can view the result at: %s\n", info.Results["resize"][0].SSLURL)

	return nil
}
