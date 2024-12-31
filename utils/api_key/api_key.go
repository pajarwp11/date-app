package api_key

import (
	"errors"
	"net/http"
	"os"
)

func VerifyApiKey(r *http.Request) error {
	authHeader := r.Header.Get("api-key")
	if authHeader == "" {
		return errors.New("api key header is missing")
	}

	if authHeader != os.Getenv("API_KEY") {
		return errors.New("invalid api key")
	}

	return nil
}
