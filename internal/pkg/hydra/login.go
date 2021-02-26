package hydra

import (
	"context"
	"github.com/gorilla/csrf"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"html/template"
	"log"
	"net/http"
)

var form = `
<html>
<head>
<title>Log in!</title>
</head>
<body>
<form method="POST" action="/login" accept-charset="UTF-8">
<input type="text" name="email">
<input type="password" name="password">
<input type="hidden" name="challenge" value="{{.hydraChallenge}}">
<!--
The default template tag used by the CSRF middleware .
This will be replaced with a hidden <input> field containing the
masked CSRF token.
-->
{{ .csrfField }}
<input type="submit" value="Login">
</form>
</body>
</html>
`

var t = template.Must(template.New("signup_form.tmpl").Parse(form))

func GetLogin(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	loginResponse, err := adminClient.Admin.GetLoginRequest(&admin.GetLoginRequestParams{
		Context:        context.Background(),
		LoginChallenge: query.Get("login_challenge"),
	})
	if err != nil {
		log.Fatalf("failed to parse Hydra login Challenge: %v", err)
	}

	loginRequest := loginResponse.GetPayload()

	if *loginRequest.Skip {
		// Hydra already authenticated, no need to check credentials
		acceptLoginResponse, err := adminClient.Admin.AcceptLoginRequest(&admin.AcceptLoginRequestParams{
			Context: context.Background(),
			Body: &models.AcceptLoginRequest{
				Subject: loginRequest.Subject,
			},
			LoginChallenge: *loginRequest.Challenge,
		})
		if err != nil {
			log.Fatalf("failed to accept Login request: %v", err)
		}
		http.Redirect(w, r, *acceptLoginResponse.GetPayload().RedirectTo, http.StatusSeeOther)
	} else {
		// signup_form.tmpl just needs a {{ .csrfField }} template tag for
		// csrf.TemplateField to inject the CSRF token into. Easy!
		err := t.ExecuteTemplate(w, "signup_form.tmpl", map[string]interface{}{
			csrf.TemplateTag: csrf.TemplateField(r),
			"hydraChallenge": *loginRequest.Challenge,
		})
		if err != nil {
			log.Fatalf("failed to parse login template: %v", err)
		}
	}
}

func PostLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("email")
	if username == "test@mail.com" && r.FormValue("password") == "pass123" {
		loginResponse, err := adminClient.Admin.GetLoginRequest(&admin.GetLoginRequestParams{
			Context:        context.Background(),
			LoginChallenge: r.FormValue("challenge"),
		})
		if err != nil {
			log.Fatalf("failed to parse Hydra login Challenge: %v", err)
		}

		loginRequest := loginResponse.GetPayload()

		acceptLoginResponse, err := adminClient.Admin.AcceptLoginRequest(&admin.AcceptLoginRequestParams{
			Context: context.Background(),
			Body: &models.AcceptLoginRequest{
				Subject: &username,
			},
			LoginChallenge: *loginRequest.Challenge,
		})
		if err != nil {
			log.Fatalf("failed to accept Login request: %v", err)
		}
		http.Redirect(w, r, *acceptLoginResponse.GetPayload().RedirectTo, http.StatusSeeOther)
	}
}
