package sms

import (
	"encoding/json"
	"fmt"
	prom "github.com/prometheus/alertmanager/template"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	info := prom.Data{}
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		fmt.Printf("Decode info to struct error: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	sms := SMS{}
	sms.AddMsg(info.Alerts)
	sms.Send()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("done!"))
}
