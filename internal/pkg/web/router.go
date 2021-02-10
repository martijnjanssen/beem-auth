package web

import (
	"log"
	"net/http"

	"beem-auth/internal/pkg/auth"
	"github.com/gorilla/mux"
)

func ListenHTTP() func() {
	r := mux.NewRouter()

	r.HandleFunc("/challenge/complete", wrapTransaction(challengeCompleteHandler))
	r.HandleFunc("/oauth2/auth", auth.AuthorizeHandlerFunc)
	r.HandleFunc("/oauth2/token", auth.TokenHandlerFunc)

	srv := &http.Server{Addr: ":8081", Handler: r}

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8081", nil))

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
