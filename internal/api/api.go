package api

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

	"github.com/duchoang206h/send-cli/config"
)

type (
	Config     struct{}
	ResultHttp struct {
		Result string `json:"result"`
	}
)

func UploadFile(ctx context.Context, filePath string) (*ResultHttp, error) {
	uploadUrl := fmt.Sprintf("%s/api/file", config.Config("API_URL"))
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	if err != nil {
		return nil, err
	}
	err = writer.Close()

	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, "POST", uploadUrl, &requestBody)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var result ResultHttp
	err = json.Unmarshal(body, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}
