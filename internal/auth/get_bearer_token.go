package auth

import(
	"net/http"
	"strings"
	"errors"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == ""{
		return "", errors.New("Authorization header is missing")
	}

	token := authHeader[len("bearer "):]
	return strings.TrimSpace(token), nil
}