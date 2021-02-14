package web

import (
	"context"
	"errors"
	"log"
	"net/http"

	"beem-auth/internal/pkg/database"

	"github.com/jmoiron/sqlx"
)

var d *sqlx.DB

// To return to the wrapTransaction when there is no action to perform after committing
var nilFn func() = func() {}

// handleError signals that there was an error while processing,
// signalling a rollback of the transsaction
var handleError error = errors.New("error while handling request")

func SetDB(db *sqlx.DB) {
	d = db
}

func wrapTransaction(handler func(context.Context, database.Queryer, http.ResponseWriter, *http.Request) (func(), error)) func(http.ResponseWriter, *http.Request) {
	if d == nil {
		log.Fatal("SetDB() has to be used before wrapping requests with transactions")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Start transaction and attach to context
		tx, err := d.Beginx()
		if err != nil {
			log.Printf("unable to start transaction: %s", err)
			writeResponse(w, http.StatusInternalServerError, "unable to handle request")
			return
		}

		// Call handler of the route to handle the request, if there was an error,
		// roll back the transaction. The successFunc is only called after the transaction
		// is committed, since committing can still fail, resulting in the request not being
		// processed correctly.
		successFunc, err := handler(r.Context(), tx, w, r)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Printf("error while rolling back transaction: %s", rollbackErr)
			}
			return
		}

		// Try to commit the transaction
		commitErr := tx.Commit()
		if commitErr != nil {
			log.Printf("unable to commit transaction: %s", commitErr)
			writeResponse(w, http.StatusInternalServerError, "unable to handle request")
			return
		}

		successFunc()
	}
}

func writeResponse(w http.ResponseWriter, s int, m string) {
	w.WriteHeader(s)
	w.Write([]byte(m))
}
