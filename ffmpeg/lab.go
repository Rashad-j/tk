package ffmpeg

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	chatgpt "github.com/rashad-j/tiktokuploader/chatGPT"
	"github.com/rashad-j/tiktokuploader/configs"
	"github.com/rashad-j/tiktokuploader/instagram"
	"github.com/rashad-j/tiktokuploader/tiktok"
	transloadit "github.com/rashad-j/tiktokuploader/transloadit"
	"github.com/rashad-j/tiktokuploader/utils"
	"github.com/rs/zerolog/log"
)

// subtitle represents a single subtitle
type subtitle struct {
	Index int    // subtitle index
	Start string // start time
	End   string // end time
	Text  string // subtitle text
}

var colors map[string]string = map[string]string{
	"Yellow":         "H0000FFFF",
	"Coral":          "H00507FFF",
	"DeepSkyBlue":    "H00FFBF00",
	"Gold":           "H0000D7FF",
	"Lavender":       "H00FAE6E6",
	"LimeGreen":      "H0032CD32",
	"MediumSeaGreen": "H0071B33C",
	"Orange":         "H0000A5FF",
}

var fonts []string = []string{"Bebas Neue", "Yuji Boku"} //, "Bruno Ace SC", "Dancing Script", "Roboto Medium", "Roboto Mono", "Roboto Slab"}

func RunLong(cfg configs.Config) error {
	ffmpegPath := "./ffmpeg/media"
	var err error
	// get an audio from quotes
	audio, _ := selectRandomFile("media/long/narrative")
	// get an audio from tk
	music, _ := selectRandomFile("media/long/audio")
	// get a video from pexels, say the first
	video, _ := selectRandomFile("media/long/video")

	withQuoteAudio := "ffmpeg/media/withQuoteAudio.mp4"
	if err := addQuoteAudio(audio, video, withQuoteAudio); err != nil {
		fmt.Println(err)
	}
	// add the tiktok audio with 50% volume
	withTkAudio := "ffmpeg/media/withTkAudio.mp4"
	if err := addTkAudio(music, withQuoteAudio, withTkAudio); err != nil {
		fmt.Println(err)
	}

	withDate := "ffmpeg/media/withDate.mp4"
	if err := addDate(withTkAudio, withDate); err != nil {
		fmt.Println(err)
	}

	// get the audio subtitle from transloadit
	_, err = getSubtitle(cfg, audio, ffmpegPath, "srtFile.srt")
	if err != nil {
		fmt.Println(err)
	}

	withSubtitle := fmt.Sprintf("ffmpeg/media/withSubtitle%d.mp4", 0)
	srtFile := "ffmpeg/media/srtFile.srt"
	if err := addSubtitle(withDate, srtFile, withSubtitle); err != nil {
		fmt.Println(err)
	}

	return nil
}

func Run(cfg configs.Config) error {
	ffmpegPath := "./ffmpeg/media"
	allPublishedTitles := ""
	var err error
	for {
		for i := 0; i < 3; i++ {
			// get an audio from quotes
			audio, _ := selectRandomFile("media/short/quote")
			// get an audio from tk
			tkAudio, _ := selectRandomFile("media/short/tiktok")
			// get a video from pexels, say the first
			video, _ := selectRandomFile("media/short/videos")
			// get a photo
			// photo, _ := selectRandomFile("media/photos")
			// // video from photo
			// fromPhotoVideo := "ffmpeg/media/fromPhoto.mp4"
			// // generate the video from this photo
			// if err := fromPhoto(photo, fromPhotoVideo); err != nil {
			// 	fmt.Println(err)
			// 	i--
			// 	continue
			// }
			// add the audio
			withQuoteAudio := "ffmpeg/media/withQuoteAudio.mp4"
			if err := addQuoteAudio(audio, video, withQuoteAudio); err != nil {
				fmt.Println(err)
				i--
				continue
			}
			// add the tiktok audio with 50% volume
			withTkAudio := "ffmpeg/media/withTkAudio.mp4"
			if err := addTkAudio(tkAudio, withQuoteAudio, withTkAudio); err != nil {
				fmt.Println("error: ", err)
				i--
				continue
			}

			withDate := "ffmpeg/media/withDate.mp4"
			if err := addDate(withTkAudio, withDate); err != nil {
				fmt.Println("error: ", err)
				i--
				continue
			}

			// get the audio subtitle from transloadit
			_, err = getSubtitle(cfg, audio, ffmpegPath, "srtFile.srt")
			if err != nil {
				fmt.Println(err)
				i--
				continue
			}

			withSubtitle := fmt.Sprintf("ffmpeg/media/withSubtitle%d.mp4", i)
			srtFile := "ffmpeg/media/srtFile.srt"
			if err := addSubtitle(withDate, srtFile, withSubtitle); err != nil {
				fmt.Println("error: ", err)
				i--
				continue
			}

			publishedTitle, err := publish(cfg, withSubtitle, srtFile)
			if err != nil {
				fmt.Println("error: ", err)
				i--
				continue
			}

			// get a thump nail
			thumpnailPath := fmt.Sprintf("ffmpeg/media/thumpnail%d.jpg", i)
			if err := generateThumbnail(withSubtitle, thumpnailPath); err != nil {
				fmt.Println("error: ", err)
				i--
				continue
			}
			// publish to instagram
			ur := instagram.UploadRequest{
				Cover:    thumpnailPath,
				Video:    withSubtitle,
				Caption:  publishedTitle,
				Username: "mindset.game2",
			}

			if err := instagram.UploadVideo(ur, "http://localhost:3000/video"); err != nil {
				fmt.Println("error: ", err)
				i--
				continue
			}

			fmt.Println("insta uploaded")

			//
			allPublishedTitles = fmt.Sprintf("%s\n%s", allPublishedTitles, publishedTitle)

			// transfer the srt file to .ass, make the words appear one by one
			// subtitlesList, err := parseSRTFile(srtFile)
			// if err != nil {
			// 	return err
			// }
			// log.Info().Interface("list", subtitlesList).Msg("subtitles list")

			// drawText, _ := generateTopDown(subtitlesList, 0)
			// fmt.Println(drawText)

			// assFile := "ffmpeg/media/ass.ass"
			// if err := generateAssFile(subtitlesList, assFile); err != nil {
			// 	return err
			// }

			// // add the subtitle to the video
			// if err := addAssFile(withTkAudio, assFile, withSubtitle); err != nil {
			// 	return err
			// }

			// publish it, otherwise everything will be overwritten

			// log.Info().Msg("starting insta upload")
			// ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
			// ur := instagram.UploadRequest{
			// 	Cover:    video.Cover,
			// 	PlayURL:  video.PlayURL,
			// 	Caption:  video.Title,
			// 	Duration: int(video.Duration),
			// 	Username: channel.InstaID,
			// }
			// if err := instagram.Upload(ctx, ur, cfg.InstaUploadURL); err != nil {
			// 	log.Err(err).Msg("failed to upload instagram video")
			// 	cancel()
			// }
			// log.Info().Msg("insta video uploaded")
			// cancel()

			fmt.Println("sleeping...")
			time.Sleep(time.Minute * 10)
		}
		// finalVideo, err := mergeVideos(ffmpegPath)
		// if err != nil {
		// 	log.Error().Err(err).Msg("failed to merge all videos")
		// }
		// finalTitle, err := publishMerged(cfg, allPublishedTitles, finalVideo)
		// if err != nil {
		// 	log.Error().Err(err).Msg("failed to publish long video")
		// }

		// log.Info().Str("all merged titles", allPublishedTitles).Str("final title", finalTitle).Msg("successfully published long video")

		fmt.Println("starting next round")
		time.Sleep(time.Hour * 12)
	}
}

// lengthUpdate make the video matches the audio length, loop if necessary
func addQuoteAudio(audioPath, videoPath, outputPath string) error {
	cmd := exec.Command("ffmpeg",
		"-stream_loop", "-1",
		"-i", videoPath,
		"-i", audioPath,
		"-map", "0:v:0",
		"-map", "1:a:0",
		"-c:v", "copy",
		"-c:a", "aac",
		"-shortest",
		"-tune",
		"film",
		outputPath, "-y",
	)
	// Set stdout and stderr to os.Stdout and os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// execute FFmpeg command
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func addTkAudio(audioPath, videoPath, outputPath string) error {
	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-stream_loop", "-1",
		"-i", audioPath,
		"-filter_complex", "[0:a]volume=1[a];[1:a]volume=0.25[b];[a][b]amerge=inputs=2[aout]",
		"-map", "0:v:0",
		"-map", "[aout]",
		"-c:v", "copy",
		"-c:a", "aac",
		outputPath, "-y",
	)

	// Set stdout and stderr to os.Stdout and os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Run the command and check for errors
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running ffmpeg command:", err)
		return err
	}

	return nil
}

func getSubtitle(cfg configs.Config, audioPath, subtitlePath, subtitleName string) (string, error) {
	subtitle := fmt.Sprintf("%s/%s", subtitlePath, subtitleName)
	url, err := transloadit.Subtitle(cfg, audioPath)
	if err != nil {
		return "", err
	}
	// download the subtitle
	if err := utils.Download(subtitlePath, subtitleName, url); err != nil {
		return "", err
	}
	return subtitle, nil
}

// parseSRTFile parses an .srt file and returns a slice of subtitles
func parseSRTFile(filePath string) ([]subtitle, error) {
	var subtitles []subtitle

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var idx int
	var start, end string
	var text string

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case idx == 0 && line == "":
			// Skip the first empty line
			continue
		case idx == 0 && isNumeric(line):
			// The first line is the index
			idx, _ = strconv.Atoi(line)
		case start == "" && strings.Contains(line, "-->"):
			// Parse the start and end time
			times := strings.Split(line, "-->")
			start, _ = parseSRTTime(times[0])
			end, _ = parseSRTTime(times[1])
		case line == "":
			// Add the subtitle to the slice
			subtitles = append(subtitles, subtitle{idx, start, end, text})
			idx = 0
			start = ""
			end = ""
			text = ""
		default:
			// Add the subtitle text
			if len(text) > 0 {
				text += "\n"
			}
			text += line
		}
	}

	// Add the last subtitle to the slice
	if idx != 0 && start != "0" && end != "0" && len(text) > 0 {
		subtitles = append(subtitles, subtitle{idx, start, end, text})
	}

	return subtitles, nil
}

// isNumeric checks whether a string is numeric
func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// parseSRTTime parses an .srt time string and returns a duration
func parseSRTTime(timeStr string) (string, error) {
	timeStr = strings.TrimSpace(timeStr)
	re := regexp.MustCompile(`\d{2}:\d{2}:\d{2},\d{3}`)
	if !re.MatchString(timeStr) {
		return "", fmt.Errorf("invalid SRT time format: %s", timeStr)
	}

	parts := strings.Split(timeStr, ":")
	hours, _ := strconv.Atoi(parts[0])
	minutes, _ := strconv.Atoi(parts[1])
	secondsAndMs := strings.Split(parts[2], ",")
	seconds, _ := strconv.Atoi(secondsAndMs[0])
	millis, _ := strconv.Atoi(secondsAndMs[1])

	return fmt.Sprintf("%d:%02d:%02d.%02d", hours, minutes, seconds, millis/10), nil
}

func generateAssFile(subtitles []subtitle, assFilePath string) error {
	assLines := make([]string, 0)

	for _, sub := range subtitles {
		// convert start and end time to seconds
		startTime, _ := convertTimeToSeconds(sub.Start)
		endTime, _ := convertTimeToSeconds(sub.End)

		// split the text into words
		words := strings.Split(sub.Text, " ")

		// calculate the duration for each word
		originalWordsDur := float64(endTime-startTime) / float64(len(words))
		wordsDuration := originalWordsDur / float64(2.5)

		text := ""
		// loop through each word and create an ASS line
		wordStartTime := startTime
		wordEndTime := endTime

		everyQuickWords := 0
		quickWords := 3

		for j, word := range words {
			text = fmt.Sprintf("%s %s", text, word)

			if everyQuickWords == 0 {
				wordStartTime = startTime + float64(j)*originalWordsDur
			} else {
				wordStartTime = wordStartTime + wordsDuration
			}

			wordEndTime = wordStartTime + wordsDuration

			if everyQuickWords == quickWords {
				wordEndTime = startTime + float64(j+1)*originalWordsDur
				everyQuickWords = 0
			} else {
				everyQuickWords++
			}

			if j+1 == len(words) {
				wordEndTime = endTime
			}

			assLine := fmt.Sprintf("Dialogue: 0,%s,%s,Default,,0,0,0,,%s", convertSecondsToTime(wordStartTime), convertSecondsToTime(wordEndTime), text)
			assLines = append(assLines, assLine)
		}
	}

	// join all lines and print to console
	assContent := strings.Join(assLines, "\n")

	file, err := os.Create(assFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// write the text to the file
	_, err = file.WriteString("[Script Info]\n; Script generated by FFmpeg/LavcLIBAVCODEC_VERSION\nScriptType: v4.00+\nPlayResX: 384\nPlayResY: 288\nScaledBorderAndShadow: yes\n\n[V4+ Styles]\nFormat: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding\nStyle: Default,Bebas Neue,12,&H0000FFFF,&H0000FFFF,&H0,&H0,0,0,0,0,100,100,0,0,1,0.5,0,4,10,10,10,0\n\n[Events]\nFormat: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text\n")
	if err != nil {
		return err
	}
	_, err = file.WriteString(assContent)
	if err != nil {
		return err
	}

	return nil
}

// helper function to convert time to seconds
func convertTimeToSeconds(timeStr string) (float64, error) {
	timeFormat := "15:04:05.00"
	time, err := time.Parse(timeFormat, timeStr)

	if err != nil {
		return 0, err
	}

	return float64(time.Hour()*3600+time.Minute()*60+time.Second()) + float64(time.Nanosecond())/1e9, nil
}

// helper function to convert seconds to time
func convertSecondsToTime(seconds float64) string {
	hours := int(seconds) / 3600
	minutes := (int(seconds) % 3600) / 60
	secondsFloat := seconds - float64(hours*3600+minutes*60)
	secondsStr := fmt.Sprintf("%.2f", secondsFloat)

	return fmt.Sprintf("%d:%02d:%s", hours, minutes, secondsStr)
}

func addAssFile(videoPath, assFile, outputPath string) error {
	// prepare the ffmpeg command
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-vf", fmt.Sprintf("ass=%s", assFile), outputPath, "-y")

	// set the output of the command to the console
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// execute the command
	err := cmd.Run()
	if err != nil {
		return err
	}

	return err
}

func selectRandomFile(dirPath string) (string, error) {
	// Read the list of files in the directory
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return "", err
	}

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Select a random file from the list
	fileIndex := rand.Intn(len(files))
	selectedFile := files[fileIndex]

	// Return the full absolute path for the selected file
	return filepath.Join(dirPath, selectedFile.Name()), nil
}

func addSubtitle(videoPath, subtitle, outputPath string) error {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
	// Get a random key
	keys := make([]string, 0, len(colors))
	for k := range colors {
		keys = append(keys, k)
	}
	randomKey := keys[rand.Intn(len(keys))]
	color := fmt.Sprintf("&%s", colors[randomKey])
	fmt.Println("color picked", color)
	// &H0000FFFF original
	// &H00FFBF00 picked example

	font := fonts[rand.Intn(len(fonts))]

	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-vf", fmt.Sprintf("subtitles=%s:force_style='Alignment=10,Fontname=%s,Fontsize=22,Outline=0.5,PrimaryColour=%s'", subtitle, font, color),
		outputPath, "-y",
	)

	// set the output of the command to the console
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// execute the command
	err := cmd.Run()
	if err != nil {
		return err
	}

	return err
}

func topDowntext(videoPath, subtitlePath, outputPath string) error {
	// Define the video, subtitle, and output paths
	// Parse the subtitle file and extract the subtitle segments
	segments, err := parseSRTFile(subtitlePath)
	if err != nil {
		return err
	}

	// Define the subtitle style
	style := "force_style='Alignment=10,Fontname=Bebas Neue,Fontsize=20,Outline=0.5,PrimaryColour=&H0000FFFF'"

	// Define the subtitle filter
	var filter []string
	for _, segment := range segments {
		// Add the subtitle segment to the filter
		filter = append(filter, fmt.Sprintf("subtitles=%s:enable='between(t,%f,%f)':%s", subtitlePath, segment.Start, segment.End, style))

		// Add some vertical spacing between the subtitle segments
		filter = append(filter, "split[v"+strconv.Itoa(segment.Index)+"]")
		if segment.Index > 1 {
			filter = append(filter, "[v"+strconv.Itoa(segment.Index-1)+"][v"+strconv.Itoa(segment.Index)+"]vstack=inputs=2")
		}
	}

	// Build the ffmpeg command
	cmdArgs := []string{"-i", videoPath, "-vf", strings.Join(filter, ","), outputPath, "-y"}
	cmd := exec.Command("ffmpeg", cmdArgs...)
	// Execute the ffmpeg command and print the output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()
	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func textHeightwidth() {
	cmd := exec.Command("ffprobe", "-f", "lavfi", "-i", "color=c=black:s=640x480,drawtext=fontsize=30:fontfile=FreeSerif.ttf:text='hello world':x=(w-text_w)/2:y=(h-text_h)/2", "-show_entries", "frame=pkt_pts_time,w,h", "-of", "default=noprint_wrappers=1", "-v", "quiet")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Parse the output to get the values of text_w and text_h
	// Output format: n:1 text_w=x text_h=y w=z h=w
	// Extract the values of text_w and text_h from the output
	// and convert them to integers
	var text_w, text_h int
	// n, err := fmt.Sscanf(out.String(), "%d:1 text_w=%d text_h=%d w=%d h=%d", &n, &text_w, &text_h)
	// if err != nil || n != 5 {
	// 	fmt.Println("Error parsing output:", err)
	// 	return
	// }

	fmt.Println("text_w:", text_w)
	fmt.Println("text_h:", text_h)
}

func executeTopDown(videoPath, outputPath, cmd string) {
	args := []string{
		"-i", videoPath,
		"-vf", fmt.Sprintf("'%s'", cmd),
		"-c:a", "copy",
		"-y", outputPath,
	}
	command := exec.Command("ffmpeg", args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		return
	}
}

func generateTopDown(subtitles []subtitle, videoEnd float64) (string, error) {
	var result string
	baseMargin := 200
	end := videoEnd
	if end == 0 {
		duration, err := time.ParseDuration(subtitles[len(subtitles)-1].End)
		if err != nil {
			log.Error().Err(err)
			return "", err
		}
		end = duration.Seconds()
		end += 3 // 3 seconds usually needed
	}

	for i, sub := range subtitles {
		duration, err := time.ParseDuration(sub.Start)
		if err != nil {
			log.Error().Err(err)
			return "", err
		}
		start := duration.Seconds()
		fmt.Println("start", start, "end", end)

		result += fmt.Sprintf("drawtext=fontfile=ffmpeg/fonts/BebasNeue-Regular.ttf:text='%s':x=(w-text_w)/2:y=(h-%d):fontsize=180:fontcolor=0xFFFF00:enable='between(t,%f,%f)'", sub.Text, sub.Index*baseMargin, start, end)
		if i < len(subtitles)-1 {
			result += ","
		}
	}

	return result, nil
}

func SrtAsText(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var sb strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			// Empty line means end of subtitle block
			// sb.WriteString("")
		} else if _, err := strconv.Atoi(line); err == nil {
			// Line contains the subtitle index, skip it
		} else if strings.Contains(line, "-->") {
			// Line contains the subtitle timecode, skip it
		} else {
			// Line contains the subtitle text, append it
			sb.WriteString(line)
			sb.WriteString(" ")
		}
	}

	return sb.String(), nil
}

func publish(cfg configs.Config, finalVideoPath, srtFilePath string) (string, error) {
	channel := tiktok.Channel{
		Title:            "mindsetgame",
		TransloaditCreds: "mindset-game",
		TiktokUniqueID:   "mindset.game",
		TiktokUserId:     7019876592412001285,
		FollowingsNumber: 50,
		Keywords:         "Motivation,Inspiration,Self-improvement,Success,Goal-setting,Mindset,Positive thinking,Productivity,Leadership,Confidence,Gratitude,Personal growth,Mindfulness,Habits,Time management", // TODO
		Category:         "education",
		InstaID:          "mindset.game2",
	}

	srtText, err := SrtAsText(srtFilePath)
	if err != nil {
		return "", err
	}
	titlePrompt := fmt.Sprintf("generate a youtube video title for this quote: %s. Make it in one sentence and less than 10 words. Write words in small letters", srtText)
	res, err := chatgpt.Ask(cfg.ChatGptApiKey, titlePrompt)
	if err != nil {
		return "", err
	}
	title := strings.ReplaceAll(res.Choices[0].Text, "\n", "")
	title = strings.ReplaceAll(title, "\"", "")

	authorPrompt := fmt.Sprintf("get the author of this quote: %s. If no author found return '#quote'", srtText)
	res, err = chatgpt.Ask(cfg.ChatGptApiKey, authorPrompt)
	if err != nil {
		return "", err
	}
	author := strings.ReplaceAll(res.Choices[0].Text, "\n", "")
	author = strings.ReplaceAll(author, "\"", "")

	hashtagPrompt := fmt.Sprintf("generate a random success and motivation hashtags given this quote: %s. Make sure every hashtag is one word only and all lower case letters, e.g. #success.", srtText)
	res, err = chatgpt.Ask(cfg.ChatGptApiKey, hashtagPrompt)
	if err != nil {
		return "", err
	}
	hashtags := strings.ReplaceAll(res.Choices[0].Text, "\n", "")
	hashtags = strings.ReplaceAll(hashtags, "\"", "")

	videoTitle := fmt.Sprintf("%s, %s %s", title, author, hashtags)

	video := tiktok.TiktokVideo{
		Title:       videoTitle,
		Description: videoTitle,
	}

	if err := transloadit.UploadYtVideo(video, cfg, channel, finalVideoPath); err != nil {
		return "", err
	}

	return videoTitle, nil
}

func publishMerged(cfg configs.Config, desc, videoPath string) (string, error) {
	channel := tiktok.Channel{
		Title:            "mindsetgame",
		TransloaditCreds: "mindset-game",
		TiktokUniqueID:   "mindset.game",
		TiktokUserId:     7019876592412001285,
		FollowingsNumber: 50,
		Keywords:         "Motivation,Inspiration,Self-improvement,Success,Goal-setting,Mindset,Positive thinking,Productivity,Leadership,Confidence,Gratitude,Personal growth,Mindfulness,Habits,Time management", // TODO
		Category:         "education",
		InstaID:          "mindset.game2",
	}

	titlePrompt := fmt.Sprintf("generate a motivational youtube video title, make it engaging, enthusiastic, and address 2nd person. Make it in one sentence. Use the following for keywords :%s", desc)
	res, err := chatgpt.Ask(cfg.ChatGptApiKey, titlePrompt)
	if err != nil {
		return "", err
	}
	title := strings.ReplaceAll(res.Choices[0].Text, "\n", "")
	title = strings.ReplaceAll(title, "\"", "")

	video := tiktok.TiktokVideo{
		Title:       title,
		Description: desc,
	}

	if err := transloadit.UploadYtVideo(video, cfg, channel, videoPath); err != nil {
		return "", err
	}

	return title, nil
}

func mergeVideos(videosPath string) (string, error) {
	// merge all videos
	outFile := fmt.Sprintf("%s/long.mp4", videosPath)
	// Get a list of all the input video file names in the directory
	files, err := os.ReadDir(videosPath)
	if err != nil {
		return "", err
	}

	// Initialize a slice to hold the input video file names
	inputFiles := make([]string, 0)
	// Loop through the files in the directory and add any video files to the slice
	for _, file := range files {
		if !file.IsDir() && (file.Name()[len(file.Name())-4:] == ".mp4") && strings.Contains(file.Name(), "Subtitle") {
			inputFiles = append(inputFiles, file.Name())
		}
	}

	// Open the list for appending or create it if it doesn't exist
	listPath := fmt.Sprintf("%s/%s", videosPath, "list.txt")
	list, err := os.OpenFile(listPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer list.Close()

	// Clear the contents of the file
	list.Truncate(0)
	list.Seek(0, 0)

	for _, file := range inputFiles {
		// inputFullpath := fmt.Sprintf("%s/%s", videosPath, file)
		text := fmt.Sprintf("file '%s'\n", file)
		if _, err := list.WriteString(text); err != nil {
			return "", err
		}
	}

	// ffmpeg -f concat -safe 0 -i mylist.txt -c copy output.mp4
	cmd := exec.Command("ffmpeg",
		"-i", listPath,
		"-f", "concat",
		"-safe", "0",
		"-c:v", "libx264",
		"-b:v", "5000k",
		"-c:a", "copy",
		outFile, "-y",
	)

	// Set stdout and stderr to os.Stdout and os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Run the command and check for errors
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	fmt.Println("Merged videos successfully!")
	// get a video title from the desc
	return outFile, nil
}

func fromPhoto(photoPath, outputPath string) error {
	// prepare the ffmpeg command
	cmd := exec.Command("ffmpeg",
		"-i", photoPath,
		"-vf", "scale=w=trunc(iw/2)*2:h=trunc(ih/2)*2",
		"-framerate", "10",
		outputPath, "-y")

	// set the output of the command to the console
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// execute the command
	err := cmd.Run()
	if err != nil {
		return err
	}

	return err
}

func addDate(inputPath, outputPath string) error {
	now := time.Now()
	// year, day, hour, minute := now.Year(), now.YearDay(), now.Hour(), now.Minute()
	year, day := now.Year(), now.YearDay()
	fontPath := "ffmpeg/fonts/BebasNeue-Regular.ttf"
	text1 := fmt.Sprintf("Year %d, Day %d", year, day)
	// text2 := fmt.Sprintf("Time %d %d", hour, minute)
	text2 := fmt.Sprintf("")
	fontSize := "60"
	fontColor := "0xFFFF00"

	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-vf", fmt.Sprintf("drawtext=fontfile=%s:text='%s':x=(w-text_w)/2:y=(text_h+110):fontsize=%s:fontcolor=%s,drawtext=fontfile=%s:text='%s':x=(w-text_w)/2:y=(text_h+text_h+140):fontsize=%s:fontcolor=%s", fontPath, text1, fontSize, fontColor, fontPath, text2, fontSize, fontColor),
		"-c:a", "copy",
		outputPath, "-y")

	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Println("Data text added successfully!")

	return nil
}

func generateThumbnail(videoPath, thumpnailPath string) error {
	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-ss", "00:00:01",
		"-vframes", "1",
		thumpnailPath, "-y")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute FFmpeg command: %s", err)
		return err
	}
	fmt.Println("thumpnail generated")
	return nil
}
