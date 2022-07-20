package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/dapper-labs-talent/cc_ashishchandra_BackendAPI/gtypes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Tests signup functionality to create a new user in the database. This test passes if the user is successfully created
func TestSignupUser(t *testing.T) {
	log.Println()
	log.Println("Testing signup of new user function")
	log.Println()

	router := gin.Default()
	router.POST("/signup", SignupUser)

	session := time.Now().Nanosecond() // unique session Id
	jsonString := fmt.Sprintf(`{
		"email":"testemail%s@testing.com",
		"password":"mypassword",
		"firstName":"Test",
		"lastName":"User"
	}`, strconv.Itoa(session))
	postData := []byte(jsonString)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(postData))

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// Tests whether an unauthenticated user can call protected methods such as GET /users. This test passes if the GET /users fails with HTTP 401 (Unauthorized)
func TestGetAllUsersUnathenticatedUser(t *testing.T) {
	log.Println()
	log.Println("Testing get all users (unauthenticated user) function")
	log.Println()

	router := gin.Default()
	router.GET("/users", jwtSigningAuthority.AuthRequired(), GetAllUsers)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Tests whether an authenticated user can call protected methods such as GET /users. This test passes if the GET /users succeeds
func TestGetAllUsersAthenticatedUser(t *testing.T) {
	log.Println()
	log.Println("Testing get all users (authenticated user) function")
	log.Println()

	router := gin.Default()
	router.GET("/users", jwtSigningAuthority.AuthRequired(), GetAllUsers)

	w := httptest.NewRecorder()
	var tokenStr string

	id, _ := databaseInstance.AuthenticateUser(gtypes.Login{
		Email:    "dapperlabs@testing.com",
		Password: "mypassword",
	})
	if id <= 0 {
		router.POST("/signup", SignupUser)

		jsonString := `{
		"email":"dapperlabs@testing.com",
		"password":"mypassword",
		"firstName":"Test",
		"lastName":"User"
		}`
		postData := []byte(jsonString)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(postData))
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Skip("Test setup failed. Cannot successfully create new user dapperlabs@testing.com")
		}

		var jwtToken map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &jwtToken)

		if err != nil {
			t.Skip("Cannot get JWT for this test")
		}
		tokenStr = jwtToken["token"]
	} else {
		router.POST("/login", LoginUser)

		jsonString := `{
		"email":"dapperlabs@testing.com",
		"password":"mypassword"
		}`
		postData := []byte(jsonString)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(postData))
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Skip("Test setup failed. Cannot successfully authenticate user dapperlabs@testing.com")
		}

		var jwtToken map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &jwtToken)

		if err != nil {
			t.Skip("Cannot get JWT for this test")
		}
		tokenStr = jwtToken["token"]
	}

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("x-authentication-token", tokenStr)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// This tests the login functionality and succeeds if the user is able to log in and a JWT is returned
func TestLoginUser(t *testing.T) {
	log.Println()
	log.Println("Testing login function")
	log.Println()

	router := gin.Default()
	router.POST("/login", LoginUser)

	jsonString := `{
		"email":"dapperlabs@testing.com",
		"password":"mypassword"
		}`
	postData := []byte(jsonString)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(postData))

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// Tests the update user functionality. This test succeeds if we are able to log in and then call the PUT /users API with the JWT returned from the login method
func TestUpdateUser(t *testing.T) {
	log.Println()
	log.Println("Testing update user function")
	log.Println()

	// Create router for authenticating a user who will submit the update
	router := gin.Default()
	router.POST("/login", LoginUser)

	// set up test data
	jsonString := `{
		"email":"dapperlabs@testing.com",
		"password":"mypassword"
		}`
	postData := []byte(jsonString)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(postData))

	router.ServeHTTP(w, req)

	// Skip test if login fails
	if w.Code != http.StatusOK {
		t.Skip("Test setup failed. Cannot successfully authenticate user dapperlabs@testing.com")
	}

	// Get the token from the response body
	var jwtToken map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &jwtToken)

	// Skip test if we do not have a JWT
	if err != nil {
		t.Skip("Cannot get JWT for this test")
	}
	tokenStr := jwtToken["token"]

	router.PUT("/users", jwtSigningAuthority.AuthRequired(), UpdateUser)

	jsonString = `{
		"firstName":"Dapper",
		"lastName":"Labs"
		}`
	postData = []byte(jsonString)

	req, _ = http.NewRequest(http.MethodPut, "/users", bytes.NewBuffer(postData))
	req.Header.Set("x-authentication-token", tokenStr)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	log.Println(w.Body.String())
}
