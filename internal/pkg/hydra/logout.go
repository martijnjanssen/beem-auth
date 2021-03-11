package hydra

import (
	"context"
	"github.com/ory/hydra-client-go/client/admin"
	"log"
	"net/http"
)

func GetLogout(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	//logoutRequest, err := adminClient.Admin.GetLogoutRequest(&admin.GetLogoutRequestParams{
	//	Context:         context.Background(),
	//	LogoutChallenge: query.Get("logout_challenge"),
	//})
	//if err != nil {
	//	log.Fatalf("failed to retreive Hydra logout request: %v", err)
	//}

	acceptLogoutResponse, err := adminClient.Admin.AcceptLogoutRequest(&admin.AcceptLogoutRequestParams{
		LogoutChallenge: query.Get("logout_challenge"),
		Context:         context.Background(),
	})
	if err != nil {
		log.Fatalf("failed to accept Consent request: %v", err)
	}

	http.Redirect(w, r, *acceptLogoutResponse.GetPayload().RedirectTo, http.StatusSeeOther)

}
