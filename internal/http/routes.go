package http

import (
	"fmt"
	"log"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

// TestAuth tests that authentication is working
func TestAuth(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	log.Printf("user: %+v", user)
	fmt.Fprintf(w, "TEST")
}
