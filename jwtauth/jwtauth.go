package jwtauth

import (
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt"
)

const TOKEN_VALIDITY_IN_SECONDS uint = 600 // Valid for 10 minutes

// This signing authority struct will hold private and public keys for signing and verification of the token
type SigningAuthority struct {
	verifyPubKey   *ecdsa.PublicKey
	signingPrivKey *ecdsa.PrivateKey
}

// Custom claim instantiated from jwt.StandardClaims type
type DapperLabsCustomClaims struct {
	jwt.StandardClaims        // inherit StandardClaims from jwt package
	Data               string `json:"data"` // Custom storage field for this derived type. We will use this to store email
}

// LoadKeys loads the private and public keys from keys' files on the server
func (s *SigningAuthority) LoadKeys(privKeyPath string, pubKeyPath string) error {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		return err
	}

	s.signingPrivKey, err = jwt.ParseECPrivateKeyFromPEM(signBytes)
	if err != nil {
		return err
	}

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		return err
	}

	s.verifyPubKey, err = jwt.ParseECPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return err
	}

	return nil
}

// generateJWT is called on a SigningAuthority object and generates the JWT using our server's ECDSA private key for signing the token
func (sa *SigningAuthority) GenerateJWT(email string) (string, error) {
	claims := DapperLabsCustomClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Duration(TOKEN_VALIDITY_IN_SECONDS) * time.Second).Unix(),
		},
		Data: email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES512, claims)

	// Sign with our ECDSA private key and return the token
	tokenStr, err := token.SignedString(sa.signingPrivKey)

	return tokenStr, err
}

// parseJWT is called on a SigningAuthority object and generates the JWT using our server's ECDSA private key for signing the token
func (sa *SigningAuthority) ParseJWT(tokenString string) (*DapperLabsCustomClaims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("invalid or empty token")
	}

	token, err := jwt.ParseWithClaims(tokenString, &DapperLabsCustomClaims{}, func(token *jwt.Token) (interface{}, error) {

		if token.Method.Alg() != jwt.SigningMethodES512.Name {
			return nil, fmt.Errorf("token signing method does not match expected signing method: %v", token.Header["alg"])
		}

		return sa.verifyPubKey, nil
	})

	if claims, ok := token.Claims.(*DapperLabsCustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

// This is a middleware function that is invoked on an object of type SigningAuthority and validate the JWT provided in the x-authentication-token header of the request
func (sa *SigningAuthority) AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := sa.ParseJWT(ctx.GetHeader("x-authentication-token"))
		if err != nil || token.Valid() != nil {
			fmt.Println(err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		} else {
			ctx.Set("email", token.Data)
		}
	}
}
