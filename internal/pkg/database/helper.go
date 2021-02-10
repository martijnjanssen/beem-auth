package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"log"
)

// Starts a new postgres container in docker
// First return is a teardown function to be called when the user is done with the container
// Second returned value is the database
//
//
// A testing main to be used with this is for example:
//
//  func TestMain(m *testing.M) {
//    var td func()
//    td, db = StartTestPostgreSQL()
//    code := m.Run()
//    td()
//    os.Exit(code)
//  }
//
//
// In this example the teardown function td() is saved, the tests are run with m.Run(),
// teardown is done followed by the termination of the TestMain.
// The teardown CANNOT be deferred, since the os.Exit() call does not allow
// for deferred functions to be run.
func StartTestPostgreSQL() (func(), *sqlx.DB) {
	var db *sqlx.DB

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("postgres", "alpine", []string{"POSTGRES_PASSWORD=postgres", "POSTGRES_DB=beem_auth"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err = Connect("localhost", resource.GetPort("5432/tcp"), "postgres", "postgres", "beem_auth")
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	_, err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	if err != nil {
		log.Fatalf("Could not create extension")
	}

	// Return a teardown function to call when we are done with the db in the test
	return func() {
		// You can't defer this because os.Exit doesn't care for defer
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}, db
}
