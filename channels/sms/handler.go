package sms

import (
	"encoding/json"
	"fmt"
	"github.com/joecomscience/prom-webhook/notify_channel"
	prom "github.com/prometheus/alertmanager/template"
	"net/http"
	"os"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	promData := prom.Data{}
	if err := json.NewDecoder(r.Body).Decode(&promData); err != nil {
		fmt.Printf("Decode info to struct error: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	sms := notify_channel.SMS{
		PhoneNumbers: getPhoneNumbers(),
		Message: prepareMessage(promData.Alerts),
	}
	sms.Send()

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("done!"))
}

func prepareMessage(promAlertInformation prom.Alerts) []string {
	var messages []string
	for _, alertInfo := range promAlertInformation {
		alertMessage := fmt.Sprintf("Description: %s\n", alertInfo.Annotations["description"])
		messages = append(messages, alertMessage)
	}
	return messages
}

func getPhoneNumbers() []string {
	return strings.Split(os.Getenv("PHONES"), ",")
}
