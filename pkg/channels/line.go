package channels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/prometheus/alertmanager/template"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	url = "https://notify-api.line.me/api/notify"
	token = "Bearer " + "20zl2k8gtimiX1js3vWxxm0XPDAjRPLDUKQQA87y4Kz"
)

func sendLineNotify(msg []byte) error {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(msg))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", token)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	fmt.Printf("line notify response: %v\n", string(body))
	return nil
}

func LineHandler(w http.ResponseWriter, r *http.Request)  {
	defer r.Body.Close()

	alertInfo := template.Data{}
	if err := json.NewDecoder(r.Body).Decode(&alertInfo); err != nil {
		fmt.Errorf("Decode info to struct error: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var messages []string

	for _, alert := range alertInfo.Alerts {
		description := alert.Annotations["description"]
		summary := alert.Annotations["summary"]
		message := "Description: " + description + "; Summary: " + summary
		messages = append(messages, message)
	}

	fmt.Printf("time_stamp: %v, data: %v\n", time.Now(), messages)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("ok"))
}