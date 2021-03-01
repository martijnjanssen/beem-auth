package web

import (
	"beem-auth/internal/pkg/hydra"
	"github.com/gorilla/csrf"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func ListenHTTP() func() {
	r := mux.NewRouter()
	CSRF := csrf.Protect([]byte("RandomBeemAuthKeyTing"))

	r.HandleFunc("/challenge/complete", wrapTransaction(challengeCompleteHandler))

	loginEndpoint := r.PathPrefix("/login").Subrouter()
	loginEndpoint.Use(CSRF)
	loginEndpoint.HandleFunc("", hydra.GetLogin).Methods("GET")
	loginEndpoint.HandleFunc("", hydra.PostLogin).Methods("POST")

	r.HandleFunc("/consent", hydra.GetConsent)

	srv := &http.Server{Addr: ":8081", Handler: r}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("error listening: %s", err)
		}
	}()

	return func() {
		err := srv.Close()
		if err != nil {
			log.Fatalf("error closing server: %s", err)
		}
	}
}
