package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func generateJWT() (string, error) {
	secretKey := os.Getenv("SECRET_KEY")
	secretKeyBytes, err := hex.DecodeString(secretKey)
	if err != nil {
		return "", err
	}
	secretKeyEd := ed25519.PrivateKey(secretKeyBytes)

	token := jwt.New(jwt.SigningMethodEdDSA)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(1 * time.Hour)

	tokenString, err := token.SignedString(secretKeyEd)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyJWT(endpointHandler func(context *gin.Context)) gin.HandlerFunc {
	publicKey := os.Getenv("PUBLIC_KEY")
	publicKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		println("Error decoding secret key")
	}
	publicKeyEd := ed25519.PublicKey(publicKeyBytes)
	return gin.HandlerFunc(func(context *gin.Context) {
		if context.GetHeader("Token") != "" {
			token, err := jwt.Parse(context.GetHeader("Token"), func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodEd25519)
				if !ok {
					return nil, errors.New("") // no error string to avoid leaking the algorithm
				}

				return publicKeyEd, nil

			})
			if err != nil {
				context.String(http.StatusUnauthorized, "You're Unauthorized due to error parsing the JWT "+err.Error())
			} else {
				if token.Valid {
					endpointHandler(context) //if token is valid, call the actual handler function for the request
				} else {
					context.String(http.StatusUnauthorized, "You're Unauthorized due to invalid token")
				}
			}
		} else {
			context.String(http.StatusUnauthorized, "You're Unauthorized due to No token in the header")
		}
	})
}
