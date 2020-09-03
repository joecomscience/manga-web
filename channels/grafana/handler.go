package grafana

import (
	"encoding/json"
	"fmt"
	"github.com/joecomscience/prom-webhook/notify_channel"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(b))

	payload := Payload{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		fmt.Printf("Decode info to struct error: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()
	imageFileName := payload.GetImageFileName()
	if err := payload.DownloadImage(imageFileName); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	lineManager := notify_channel.LineManager{
		Message:   payload.Message,
		ImageFile: imageFileName,
	}

	if err := lineManager.Send(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := payload.RemoveImageFile(imageFileName); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type Payload struct {
	DashboardID int `json:"dashboardId"`
	EvalMatches []struct {
		Value  int         `json:"value"`
		Metric string      `json:"metric"`
		Tags   interface{} `json:"tags"`
	} `json:"evalMatches"`
	ImageURL string `json:"imageUrl"`
	Message  string `json:"message"`
	OrgID    int    `json:"orgId"`
	PanelID  int    `json:"panelId"`
	RuleID   int    `json:"ruleId"`
	RuleName string `json:"ruleName"`
	RuleURL  string `json:"ruleUrl"`
	State    string `json:"state"`
	Tags     struct {
	} `json:"tags"`
	Title string `json:"title"`
}

func (payload Payload) DownloadImage(filename string) error {
	resp, err := http.Get(payload.ImageURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	output, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, resp.Body)
	return err
}

func (payload Payload) RemoveImageFile(filename string) error {
	return os.Remove(filename)
}

func (payload Payload) GetImageFileName() string {
	filenames := strings.Split(payload.ImageURL, "/")
	lastIndex := len(filenames) - 1
	return filenames[lastIndex]
}
