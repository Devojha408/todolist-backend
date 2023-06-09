package helpers

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payload struct {
	FirstName string
	LastName  string
	Email     string
}

type Claims struct {
	FirstName string             `json:"firstname"`
	LastName  string             `json:"lastname"`
	Email     string             `json:"email"`
	Id        primitive.ObjectID `json:"_id"`
	jwt.StandardClaims
}

var JWT_SECRET string

func GenerateJwtToken(payload Payload) (string, error) {
	if JWT_SECRET = os.Getenv("JWT_SECRET"); JWT_SECRET == "" {
		log.Fatal("[ ERROR ] JWT_SECRET environment variable not provided!\n")
	}

	key := []byte(JWT_SECRET)

	expirationTime := time.Now().Add(7 * 24 * 60 * time.Minute)

	claims := &Claims{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	UnsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	SignedToken, err := UnsignedToken.SignedString(key)
	if err != nil {
		return "", err
	}

	return SignedToken, nil
}

func VerifyJwtToken(strToken string) (*Claims, error) {
	if JWT_SECRET = os.Getenv("JWT_SECRET"); JWT_SECRET == "" {
		log.Fatal("[ ERROR ] JWT_SECRET environment variable not provided!\n")
	}

	key := []byte(JWT_SECRET)

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(strToken, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return claims, errors.New("invalid token signature")
			//fmt.Errorf("invalid token signature")
		}
	}

	if !token.Valid {
		return claims, errors.New("invalid token")
	}

	return claims, nil
}
