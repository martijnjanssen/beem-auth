package hydra

import (
	"github.com/ory/hydra-client-go/client"
	"log"
	"net/url"
)

var adminClient *client.OryHydra

func init() {
	adminURL, err := url.Parse("http://hydra:4445")
	if err != nil {
		log.Fatalf("failed to parse Hydra adminURL: %v", err)
	}
	adminClient = client.NewHTTPClientWithConfig(nil, &client.TransportConfig{Schemes: []string{adminURL.Scheme}, Host: adminURL.Host, BasePath: adminURL.Path})
}
