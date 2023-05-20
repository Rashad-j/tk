package ffmpeg

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"

	ff "github.com/u2takey/ffmpeg-go"
)

func convertToMp4(filename string) error {
	// Open the list.txt file
	file, err := os.Open("list.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	// Build the ffmpeg command
	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", "list.txt", "-c", "copy", "output.mp4")

	// Set the current working directory to the directory of the input file
	cmd.Dir = filepath.Dir(filename)

	// Set the command output to the standard output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func test() {
	err := convertToMp4("/path/to/input/file.wav")
	if err != nil {
		fmt.Println("Error converting file:", err)
	} else {
		fmt.Println("File converted successfully")
	}
}

func removeDir(path string) error {
	// Get the absolute path of the directory
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	// Check if the directory exists
	_, err = os.Stat(absPath)
	if os.IsNotExist(err) {
		return nil
	}

	// Remove the directory and its content
	err = os.RemoveAll(absPath)
	if err != nil {
		return err
	}

	return nil
}

func testRemove() {
	err := removeDir("/path/to/directory")
	if err != nil {
		fmt.Println("Error removing directory:", err)
	} else {
		fmt.Println("Directory removed successfully")
	}
}

func Segments() {
	file, err := os.Open("ffmpeg/silence.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	inputAudio := "ffmpeg/output/output.flac"

	scanner := bufio.NewScanner(file)

	var start float64
	var end float64
	var count int

	re := regexp.MustCompile(`silence_end: (\d+\.\d+)`)

	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindStringSubmatch(line)
		if len(match) > 1 {
			end, err = strconv.ParseFloat(match[1], 64)
			if err != nil {
				panic(err)
			}
			outputAudio := fmt.Sprintf("ffmpeg/output/audio%d.flac", count)

			duration := end - start
			durationStr := fmt.Sprintf("%.4f", duration)
			err = ff.Input(inputAudio, ff.KwArgs{"ss": start}).Output(outputAudio, ff.KwArgs{"t": durationStr}).OverWriteOutput().Run()
			if err != nil {
				fmt.Println(err)
			}

			// duration := end - start
			// use -t duration
			// cmd := exec.Command("ffmpeg", "-y", "-ss", fmt.Sprintf("%f", start), "-i", inputAudio, "-t", fmt.Sprintf("%f", duration), "-c", "copy", outputAudio)
			// use -to
			//cmd := exec.Command("ffmpeg", "-y", "-ss", fmt.Sprintf("%f", start), "-i", inputAudio, "-to", fmt.Sprintf("%f", end), "-c", "copy", outputAudio)
			// err := cmd.Run()
			// if err != nil {
			// 	panic(err)
			// }

			start = end
			count++
		}
	}
}
