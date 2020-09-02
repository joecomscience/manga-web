package notify_channel

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"sync"
	"text/template"
)

type SMS struct {
	PhoneNumbers []string
	Message      []string
}

func (sms SMS) Send() {
	var wg sync.WaitGroup
	for _, phoneNumber := range sms.PhoneNumbers {
		for _, message := range sms.Message {
			wg.Add(1)
			go sendMessageToSMSGateway(message, phoneNumber, &wg)
		}
	}
}

func getPayload(msg string, to string) ([]byte, error) {
	bufferPayload := &bytes.Buffer{}
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
					<bean:msg>{{.Message}}</bean:msg>
					<bean:profile>KRUNGSRIGRP_EN</bean:profile>
					<bean:refMsgCode></bean:refMsgCode>
					<bean:scheduling></bean:scheduling>
					<bean:subApp>OPENSHIFT</bean:subApp>
					<bean:telNo>{{.SendTo}}</bean:telNo>
				</ser:bean>
			</ser:saveSMS>
			</soapenv:Body>
		</soapenv:Envelope>`)
	if err != nil {
		return nil, err
	}

	templateData := struct {
		Message string
		SendTo  string
	}{msg, to}

	if err = tmpl.Execute(bufferPayload, templateData); err != nil {
		return nil, err
	}

	return bufferPayload.Bytes(), nil
}

func sendMessageToSMSGateway(message string, sendTo string, wg *sync.WaitGroup) {
	defer wg.Done()
	url := os.Getenv("SMS_URL")
	payload, err := getPayload(message, sendTo)
	if err != nil {
		logrus.Errorf("Get data payload error: %s\n", err.Error())
		return
	}

	c := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		logrus.Errorf("Error on creating request object: %s\n", err.Error())
		return
	}
	req.Header.Set("Content-type", "text/xml")
	req.Header.Set("SOAPAction", "saveSMS")

	logrus.Info("done!")
	_, err = c.Do(req)
	if err != nil {
		logrus.Errorf("Error on dispatching request: %s\n ", err.Error())
		return
	}
}
