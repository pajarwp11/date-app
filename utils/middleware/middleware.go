package middleware

import (
	"date-app/models"
	"date-app/utils/jwt"
	"encoding/json"
	"net/http"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var res models.GeneralResponse
		w.Header().Set("Content-Type", "application/json")
		err := jwt.VerifyToken(r)
		if err != nil {
			res.Code = http.StatusUnauthorized
			res.Message = "invalid token"
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(res)
			return
		}
		next.ServeHTTP(w, r)
	})
}
