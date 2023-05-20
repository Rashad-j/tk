package instagram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type UploadRequest struct {
	Cover    string `json:"cover"`
	Video    string `json:"play_url"`
	Caption  string `json:"caption"`
	Duration int    `json:"duration"`
	Username string `json:"username"`
}

func UploadUrl(ctx context.Context, uploadReq UploadRequest, uploadURL string) error {
	uploadReqBody, err := json.Marshal(uploadReq)
	if err != nil {
		err = fmt.Errorf("error marshaling request body: %w", err)
		return err
	}

	req, err := http.NewRequest("POST", uploadURL, bytes.NewBuffer(uploadReqBody))
	if err != nil {
		err = fmt.Errorf("error creating request: %w", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		err = fmt.Errorf("error making request: %w", err)
		return err
	}

	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		err := fmt.Errorf("none 200 returned from express.js: %s", resp.Status)
		return err
	}

	return nil
}

func UploadVideo(uploadReq UploadRequest, uploadURL string) error {
	endpoint := uploadURL // "http://localhost:3000/video"
	requestData := map[string]string{
		"caption":  uploadReq.Caption,
		"username": uploadReq.Username,
	}

	filePaths := map[string]string{
		"image": uploadReq.Cover,
		"video": uploadReq.Video,
	}

	// Create a new HTTP request with a POST method
	request, err := newFileUploadRequest(endpoint, requestData, filePaths)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	// Send the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}

	// Print the response
	fmt.Println("Response:", string(body))
	return nil
}

func newFileUploadRequest(uri string, requestData map[string]string, filePaths map[string]string) (*http.Request, error) {
	body, writer := io.Pipe()
	multipartWriter := multipart.NewWriter(writer)

	// Create a goroutine to write the form data to the multipart writer
	go func() {
		defer writer.Close()
		defer multipartWriter.Close()

		for key, value := range requestData {
			_ = multipartWriter.WriteField(key, value)
		}

		for key, filePath := range filePaths {
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Println("Error opening file:", err)
				continue
			}
			defer file.Close()

			part, err := multipartWriter.CreateFormFile(key, filepath.Base(filePath))
			if err != nil {
				fmt.Println("Error creating form file:", err)
				continue
			}

			_, err = io.Copy(part, file)
			if err != nil {
				fmt.Println("Error copying file to form part:", err)
				continue
			}
		}
	}()

	request, err := http.NewRequest(http.MethodPost, uri, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	return request, nil
}
