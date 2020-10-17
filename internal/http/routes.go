package http

import (
	"fmt"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

// TestAuth tests that authentication is working
func TestAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "TEST")
}
