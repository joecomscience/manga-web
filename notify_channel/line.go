package notify_channel

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type LineManager struct {
	Message   string
	ImageFile string
}

func (manager LineManager) Send() error {
	lineAPIUrl := os.Getenv("LINE_URL")
	lineToken := os.Getenv("LINE_TOKEN")
	imageFile := fmt.Sprintf("./%s", manager.ImageFile)

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("message", "joewalker")

	file, err := os.Open(imageFile)
	if err != nil {
		return err
	}
	defer file.Close()

	part2, err := writer.CreateFormFile("imageFile", filepath.Base(imageFile))
	_, err = io.Copy(part2, file)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, lineAPIUrl, payload)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", lineToken))

	resp, err := client.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	return err
}
