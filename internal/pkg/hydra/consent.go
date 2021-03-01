package hydra

import (
	"context"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"log"
	"net/http"
)

var skipConsent = true

func GetConsent(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	consentRequest, err := adminClient.Admin.GetConsentRequest(&admin.GetConsentRequestParams{
		Context:          context.Background(),
		ConsentChallenge: query.Get("consent_challenge"),
	})
	if err != nil {
		log.Fatalf("failed to retreive Hydra consent request: %v", err)
	}

	consent := consentRequest.GetPayload()

	// Always skip consent for now
	if consent.Skip || skipConsent {

		acceptConsentResponse, err := adminClient.Admin.AcceptConsentRequest(&admin.AcceptConsentRequestParams{
			Body: &models.AcceptConsentRequest{
				GrantAccessTokenAudience: consent.RequestedAccessTokenAudience,
				GrantScope:               consent.RequestedScope,
				Remember:                 true,
			},
			ConsentChallenge: *consent.Challenge,
			Context:          context.Background(),
		})
		if err != nil {
			log.Fatalf("failed to accept Consent request: %v", err)
		}

		http.Redirect(w, r, *acceptConsentResponse.GetPayload().RedirectTo, http.StatusSeeOther)
	}
}
