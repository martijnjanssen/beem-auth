package web

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func ListenHTTP() func() {
	r := mux.NewRouter()

	r.HandleFunc("/challenge/complete", wrapTransaction(challengeCompleteHandler))

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
