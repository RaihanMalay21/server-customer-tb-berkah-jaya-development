package helper

import (
	"net/http"
	"log"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/RaihanMalay21/config-tb-berkah-jaya"
)

func GetIDFromToken(r *http.Request) (uint, error) {
	// mengambil cookie dari http request
	c, err := r.Cookie("token")
	if err != nil {
		log.Println("Missing token cookie:", err)
		return 0, fmt.Errorf(err.Error()) 
	}

	// mengambil token value
	tokenString := c.Value
	claims := &config.JWTClaim{}
	//parsing token jwt
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error){
		return config.JWT_KEY, nil
	})
	if err != nil || !token.Valid {
		log.Println("Error parsing token:", err)
		return 0, fmt.Errorf(err.Error())
	}
	
	IDUser := claims.ID
	if IDUser == 0 {
		log.Println("Invalid user ID in token")
		return 0, fmt.Errorf("Error you not have in session")
	}

	return IDUser, nil
}