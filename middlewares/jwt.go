package middlewares

import (
	"fmt"
	"log"
	"strings"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/helper"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		message := map[string]interface{}{
			"message": nil,
		}

		c, err := r.Cookie("token")
		if err != nil {
			log.Println("Missing token cookie:", err)
			message["message"] = "Token is missing"
			helper.Response(w, message, http.StatusUnauthorized)
			return 
		}

		// mengambil value token
		tokenString := c.Value
		claims := &config.JWTClaim{}
		//parsing token jwt
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error){
			return config.JWT_KEY, nil
		})
		if err != nil {
			switch err {
			case jwt.ErrTokenSignatureInvalid:
				log.Println("Invalid token signature:", err)
				message["message"] = "Unauthorized"
				helper.Response(w, message, http.StatusUnauthorized)
				return
			case jwt.ErrTokenExpired:
				log.Println("Token have expired", err)
				message["message"] = "Unauthorized"
				helper.Response(w, message, http.StatusUnauthorized)
			default:
				log.Println("Error Parsing token:", err)
				message["message"] = err.Error()
				helper.Response(w, message, http.StatusUnauthorized)
				return
			}
		}

		// memeriksa apakah token available and valid
		if claims, ok := token.Claims.(*config.JWTClaim); ok && token.Valid{
			role := claims.Role
			endpoint := r.URL.Path

			fmt.Println(claims)

			if err := endPoinCanAccess(role, endpoint); err != nil {
				log.Println("Access denied to endpoint:", err)
				message["message"] = err.Error()
				helper.Response(w, message, http.StatusUnauthorized)
				return
			}
		} else {
			message["message"] = "Unauthorized"
			helper.Response(w, message, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
} 

func endPoinCanAccess(role, endpoint string) error {
	var endpoin = []string{
		"/berkahjaya/gifts/have/change/user",
		"/berkahjaya/users/data",
		"/berkahjaya/proses/poin/verify",
		"/berkahjaya/scan/poin",
		"/berkahjaya/tukar/poin/hadiah",
		"/berkahjaya/user/proses/hadiah",
		"/berkahjaya/user/remove/nota/not/valid/",
		"/berkahjaya/change/password",
	}

	if role == "Customers" {
		for _, en := range endpoin {
			if strings.HasPrefix(endpoint, en) {
				return nil
			}
		}
	}

	return fmt.Errorf("access denied to endpoint: %s", endpoint)
}