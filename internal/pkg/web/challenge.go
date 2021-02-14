package web

import (
	"context"
	"log"
	"net/http"

	"beem-auth/internal/pkg/database"
)

func challengeCompleteHandler(ctx context.Context, db database.Queryer, w http.ResponseWriter, r *http.Request) (func(), error) {
	key := r.FormValue("key")

	if key == "" {
		writeResponse(w, http.StatusBadRequest, "key must be set")
		return nilFn, handleError
	}

	userId, err := database.ChallengeComplete(ctx, db, key)
	if err != nil {
		log.Printf("unable to complete challenge: %s", err)
		writeResponse(w, http.StatusInternalServerError, "unable to handle request")
		return nilFn, handleError
	}

	err = database.UserSetValid(ctx, db, userId)
	if err != nil {
		log.Printf("unable to set user valid: %s", err)
		writeResponse(w, http.StatusInternalServerError, "unable to handle request")
		return nilFn, handleError
	}

	return func() {
		writeResponse(w, http.StatusOK, "email validated")
	}, nil
}
