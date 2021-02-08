package email

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
)

var apiURL = "https://api.sendinblue.com/v3"

const sendEmailTemplate = `{
  "sender":{
    "name": "Beem-auth",
    "email":"no-reply@beem.com"
  },
  "to": [
    {"email":"{{.Recipient}}"}
  ],
  "subject":"{{.Subject}}",
  "htmlContent":"{{.Content}}"
}`

type sendinblue struct {
	apiKey   string
	template *template.Template
}

func NewSendinblue(apiKey string) *sendinblue {
	t := template.Must(template.New("Email").Parse(sendEmailTemplate))
	return &sendinblue{apiKey: apiKey, template: t}
}

func (s *sendinblue) SendEmail(email Email) error {
	data := &strings.Builder{}
	s.template.Execute(data, email)

	req, err := http.NewRequest(http.MethodPost, getURL("/smtp/email"), strings.NewReader(data.String()))
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-key", s.apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed reading body: %w", err)
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("error while sending email: %s", body)
	}

	return nil
}

func getURL(path string) string {
	return fmt.Sprintf("%s%s", apiURL, path)
}
