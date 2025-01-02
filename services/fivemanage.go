package services

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type FiveManageService struct {
	APIKey string
	Debug  bool
}

type FMResponse struct {
	URL string `json:"url"`
	ID  string `json:"id"`
}

func init() {
	// Force HTTP/1.1
	http.DefaultTransport.(*http.Transport).TLSNextProto = map[string]func(authority string, c *tls.Conn) http.RoundTripper{}
}

func (fm *FiveManageService) debugLog(format string, args ...interface{}) {
	if fm.Debug {
		fmt.Printf("FiveManage Debug - "+format+"\n", args...)
	}
}

func (fm *FiveManageService) UploadImage(file multipart.File, metadata map[string]interface{}) (*FMResponse, error) {
	// Create a custom client with HTTP/1.1
	client := &http.Client{
		Transport: &http.Transport{
			TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
		},
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file with original filename
	filename := "image.jpg"
	if fn, ok := metadata["filename"].(string); ok {
		filename = fn
	}
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("error creating form file: %w", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("error copying file: %w", err)
	}

	// Add metadata if provided
	if metadata != nil {
		metadataBytes, err := json.Marshal(metadata)
		if err != nil {
			return nil, fmt.Errorf("error marshaling metadata: %w", err)
		}
		err = writer.WriteField("metadata", string(metadataBytes))
		if err != nil {
			return nil, fmt.Errorf("error writing metadata field: %w", err)
		}
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing writer: %w", err)
	}

	// Debug the request
	fm.debugLog("Uploading to URL: %s", "https://api.fivemanage.com/api/image")
	fm.debugLog("Content-Type: %s", writer.FormDataContentType())
	fm.debugLog("Metadata: %+v", metadata)

	req, err := http.NewRequest("POST", "https://api.fivemanage.com/api/image", body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fm.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	fm.debugLog("Response Status: %d", resp.StatusCode)
	fm.debugLog("Response Body: %s", string(respBody))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse the response
	var fmResp FMResponse
	err = json.Unmarshal(respBody, &fmResp)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	fm.debugLog("Response Headers: %+v", resp.Header)
	fm.debugLog("Parsed Response - URL: %s, ID: %s", fmResp.URL, fmResp.ID)

	if fmResp.ID == "" || fmResp.URL == "" {
		return nil, fmt.Errorf("invalid response format: %s", string(respBody))
	}

	fm.debugLog("Successfully uploaded image. URL: %s, ID: %s", fmResp.URL, fmResp.ID)
	return &fmResp, nil
}

func (fm *FiveManageService) DeleteImage(imageID string) error {
	fm.debugLog("Attempting to delete image with ID: %s", imageID)

	client := &http.Client{
		Transport: &http.Transport{
			TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
		},
	}

	// Use the correct endpoint format from docs: /api/image/delete/:id
	req, err := http.NewRequest("DELETE",
		fmt.Sprintf("https://api.fivemanage.com/api/image/delete/%s", imageID),
		nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fm.APIKey)
	req.Header.Set("Content-Type", "application/json")

	fm.debugLog("Request URL: %s", req.URL.String())
	fm.debugLog("Auth Header: %s", req.Header.Get("Authorization"))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	fm.debugLog("Response Status: %d", resp.StatusCode)
	fm.debugLog("Response Body: %s", string(body))

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusNoContent &&
		resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(body))
	}

	fm.debugLog("Successfully deleted image with ID: %s", imageID)
	return nil
}
