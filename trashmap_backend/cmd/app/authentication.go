package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
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

func verifyAPIKey(endpointHandler func(context *gin.Context)) gin.HandlerFunc {

	return gin.HandlerFunc(func(context *gin.Context) {
		apiKeyHeader := "APIKey" // make dynamic later // context.GetHeader("APIKey")
		apiKeys := map[string]string{
			"first": os.Getenv("API_KEYS"),
		} // creating a new map[string]string because I currently have no config/Database to store and load the API keys from

		decodedAPIKeys := make(map[string][]byte)
		for name, value := range apiKeys {
			decodedKey, err := hex.DecodeString(value)
			if err != nil {
				context.String(http.StatusInternalServerError, "error decoding available API keys")
				return
			}

			decodedAPIKeys[name] = decodedKey
		}

		apiKey, err := bearerToken(context, apiKeyHeader)
		fmt.Print("API Key: ", apiKey)
		fmt.Print("error: ", err)
		if err != nil {
			context.String(http.StatusUnauthorized, "invalid API key")
			return

		}

		if _, ok := apiKeyIsValid(apiKey, decodedAPIKeys); !ok {
			fmt.Print("API Key is not valid")
			context.String(http.StatusUnauthorized, "invalid API key")
			return
		}
		fmt.Print("API Key is valid")
		endpointHandler(context)
	})

}

func bearerToken(context *gin.Context, apiKeyHeader string) (string, error) {
	rawToken := context.GetHeader(apiKeyHeader)
	pieces := strings.SplitN(rawToken, ".", 2)

	if len(pieces) < 2 {
		return "", errors.New("token with incorrect bearer format")
	}

	token := strings.TrimSpace(pieces[1])

	return token, nil
}

// api key will be sent as a plain, not encoded string. It will then be hashed within this function and compared to the hashed key in my database
func apiKeyIsValid(rawKey string, availableKeys map[string][]byte) (string, bool) {
	hash := sha256.Sum256([]byte(rawKey))
	key := hash[:]
	fmt.Print("Key: ", key)
	fmt.Println("availableKeys: ", availableKeys)
	for name, value := range availableKeys {
		contentEqual := subtle.ConstantTimeCompare(value, key) == 1

		if contentEqual {
			return name, true
		}
	}

	return "", false
}
