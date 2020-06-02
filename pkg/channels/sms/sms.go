package sms

import (
	"bytes"
	"fmt"
	prom "github.com/prometheus/alertmanager/template"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"sync"
	"text/template"
)

type SMS struct {
	Phones  []string
	Message []string
}

func (s *SMS) AddMsg(alerts prom.Alerts) {
	s.Phones = strings.Split(os.Getenv("PHONES"), ",")
	for _, a := range alerts {
		msg := fmt.Sprintf("Description: %s\n", a.Annotations["description"])
		s.Message = append(s.Message, msg)
	}
}

func (s *SMS) Send() {
	var wg sync.WaitGroup
	for _, p := range s.Phones {
		for _, m := range s.Message {
			wg.Add(1)
			go sendMessageToSMSGateway(m, p, &wg)
		}
	}
}

func getPayload(msg string, to string) ([]byte, error) {
	p := &bytes.Buffer{}
	tmpl, err := template.New("sms_payload").Parse(`
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
					<bean:msg>{{.Msg}}</bean:msg>
					<bean:profile>KRUNGSRIGRP_EN</bean:profile>
					<bean:refMsgCode></bean:refMsgCode>
					<bean:scheduling></bean:scheduling>
					<bean:subApp>OPENSHIFT</bean:subApp>
					<bean:telNo>{{.To}}</bean:telNo>
				</ser:bean>
			</ser:saveSMS>
			</soapenv:Body>
		</soapenv:Envelope>`)
	if err != nil {
		return nil, err
	}

	i := struct {
		Msg string
		To      string
	}{msg, to}

	if err = tmpl.Execute(p, i); err != nil {
		return nil, err
	}

	return p.Bytes(), nil
}

func sendMessageToSMSGateway(msg string, to string, wg *sync.WaitGroup) {
	defer wg.Done()
	url := os.Getenv("SMS_URL")
	payload, err := getPayload(msg, to)
	if err != nil {
		logrus.Errorf("Get data payload error: %s\n", err.Error())
		return
	}

	c := &http.Client{}
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	r.Header.Set("Content-type", "text/xml")
	r.Header.Set("SOAPAction", "saveSMS")
	if err != nil {
		logrus.Errorf("Error on creating request object: %s\n", err.Error())
		return
	}

	logrus.Info("done!")
	_, err = c.Do(r)
	if err != nil {
		logrus.Errorf("Error on dispatching request: %s\n ", err.Error())
		return
	}
}
