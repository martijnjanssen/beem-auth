package auth

import (
	"beem-auth/internal/pkg/database"
	"github.com/jmoiron/sqlx"
	"golang.org/x/oauth2"
	"log"
	"os"
	"testing"
)

var db *sqlx.DB
var closedDb *sqlx.DB

// A valid oauth2 client (check the store) that additionally requests an OpenID Connect id token
var clientConf = oauth2.Config{
	ClientID:     "my-client",
	ClientSecret: "foobar",
	RedirectURL:  "http://localhost:3846/callback",
	Scopes:       []string{"photos", "openid", "offline"},
	Endpoint: oauth2.Endpoint{
		TokenURL: "http://localhost:3846/oauth2/token",
		AuthURL:  "http://localhost:3846/oauth2/auth",
	},
}

//// The same thing (valid oauth2 client) but for using the client credentials grant
//var appClientConf = clientcredentials.Config{
//	ClientID:     "my-client",
//	ClientSecret: "foobar",
//	Scopes:       []string{"fosite"},
//	TokenURL:     "http://localhost:3846/oauth2/token",
//}

// Initializes the database for the tests run in this package
func TestMain(m *testing.M) {
	// Functions for teardown of started docker containers
	var td, closedTd func()

	td, db = database.StartTestPostgreSQL()
	closedTd, closedDb = database.StartTestPostgreSQL()
	if err := closedDb.Close(); err != nil {
		log.Fatalf("unable to close database: %s", err)
	}
	closedTd()

	if err := database.ApplyMigrations(db); err != nil {
		log.Fatalf("migration failed: %s", err)
	}

	code := m.Run()
	td()
	os.Exit(code)
}
