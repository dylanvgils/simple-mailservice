package main

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"time"

	sendgrid "github.com/sendgrid/sendgrid-go"
)

const requestTemplate = `
	{
		"personalizations": [{
			"to": [{
				"email": "{{ .To.Address }}"
			}],
			"subject": "{{ .Subject }}"
		}],
		"from": {
			"name": "{{ .From.Name }}",
			"email": "{{ .From.Address }}"
		},
		"content": [{
			"type": "text/plain",
			"value": "{{ .Content }}\n------------------------------------\nSenderIP: {{ .SenderIP }}; Date: {{ .DateTime }}"
		}, {
			"type": "text/html",
			"value": "<html><body>{{ .Content }}<hr />SenderIP: {{ .SenderIP }}; Date: {{ .DateTime }}</body></html>"
		}],
		"tracking_settings": {
			"click_tracking": { "enable": false },
			"open_tracking": { "enable": false },
			"subscription_tracking": { "enable": false },
			"ganalytics": { "enable": false }
		}
	}`

type mail struct {
	SenderIP string    `json:"-"`
	DateTime time.Time `json:"-"`
	From     email     `json:"from,omitempty"`
	To       email     `json:"-"`
	Subject  string    `json:"subject,omitempty"`
	Content  string    `json:"content,omitempty"`
}

type email struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
}

func newMail(senderIP string) *mail {
	return &mail{
		SenderIP: senderIP,
		DateTime: time.Now(),
		To:       email{Address: os.Getenv("TO_EMAIL")},
	}
}

func (m *mail) parseTemplate() []byte {
	var tpl bytes.Buffer

	t := template.Must(template.New("Mail Request").Parse(requestTemplate))
	err := t.Execute(&tpl, m)
	if err != nil {
		log.Fatalln("Error: Encounterd error while parsing template, errror:", err)
	}

	return tpl.Bytes()
}

func (m *mail) sendMail() {
	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = m.parseTemplate()
	response, err := sendgrid.API(request)
	if err != nil {
		log.Println("Error: Encounterd error while sending request to Sendgrid API, error:", err)
	} else {
		log.Println("Success: Request on Sendgrid API fulfilled with code:", response.StatusCode)
	}
}
