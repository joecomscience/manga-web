package channels

import (
	"bytes"
	"encoding/json"
	"fmt"
	promTemplate "github.com/prometheus/alertmanager/template"
	"html/template"
	"net/http"
	"os"
	"strings"
)

type smsInfo struct {
	message string
	to      string
}

func (p *smsInfo) getPlayload() ([]byte, error) {
	tmpl, err := template.New("sms_playload").Parse(`
		<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ser="http://services.ge.com" xmlns:bean="http://bean.ge.com">
			<soapenv:Header/>
			<soapenv:Body>
			<ser:saveSMS>
				<ser:bean>
					<bean:accountNo></bean:accountNo>
					<bean:appowner>DIGITAL_TECH</bean:appowner>
					<bean:busMsgCode>IT</bean:busMsgCode>
					<bean:cardNo></bean:cardNo>
					<bean:channel>INFO</bean:channel>
					<bean:lang>THA</bean:lang>
					<bean:msg>{{.message}}</bean:msg>
					<bean:profile>KRUNGSRIGRP_EN</bean:profile>
					<bean:refMsgCode></bean:refMsgCode>
					<bean:scheduling></bean:scheduling>
					<bean:subApp>OPENSHIFT</bean:subApp>
					<bean:telNo>{{.to}}</bean:telNo>
				</ser:bean>
			</ser:saveSMS>
			</soapenv:Body>
		</soapenv:Envelope>`)
	if err != nil {
		return nil, err
	}
	var strPlayLoad bytes.Buffer
	if err = tmpl.Execute(&strPlayLoad, p); err != nil {
		return nil, err
	}
	return strPlayLoad.Bytes(), nil
}

func (p *smsInfo) sendToSmsGateway(c chan string) {
	url := os.Getenv("SMS_URL")
	pl, err := p.getPlayload()
	if err != nil {
		fmt.Printf("Get data playload error: %s\n", err.Error())
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(pl))
	req.Header.Set("Content-type", "text/xml")
	req.Header.Set("SOAPAction", "saveSMS")

	if err != nil {
		fmt.Printf("Error on creating request object: %s\n", err.Error())
		return
	}

	_, err = client.Do(req)
	if err != nil {
		fmt.Printf("Error on dispatching request: %s\n ", err.Error())
		return
	}
	c <- "message: " + p.message + ", to: " + p.to
}

func SmsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	alertInfo := promTemplate.Data{}
	if err := json.NewDecoder(r.Body).Decode(&alertInfo); err != nil {
		fmt.Printf("Decode info to struct error: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//var sms []smsInfo
	result := make(chan string)
	for _, alert := range alertInfo.Alerts {
		phones := strings.Split(alert.Annotations["phones"], ",")
		des := alert.Annotations["description"]
		sum := alert.Annotations["summary"]
		msg := "message=Description: " + des + "; Summary: " + sum

		for _, p := range phones {
			sms := smsInfo{message:msg, to: strings.TrimSpace(p)}
			sms.sendToSmsGateway(result)
			//sms = append(sms, smsInfo{message: msg, to: strings.TrimSpace(string(p))})
		}
	}

	fmt.Printf("%s\n", <-result)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
