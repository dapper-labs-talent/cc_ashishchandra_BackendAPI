package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/dapper-labs-talent/cc_ashishchandra_BackendAPI/db"
	"github.com/dapper-labs-talent/cc_ashishchandra_BackendAPI/gtypes"
	"github.com/dapper-labs-talent/cc_ashishchandra_BackendAPI/jwtauth"
	"github.com/gin-gonic/gin"
)

// Paths to the private and public key files. These are the signing authorities we will use to sign and verify our JWT
const (
	privKeyPath = "keys/app-ecdsa"
	pubKeyPath  = "keys/app-ecdsa.pub"
)

var jwtSigningAuthority *jwtauth.SigningAuthority
var databaseInstance *db.Db
var PORT string

// This function checks the environment variables, sets up the database and loads the private/public keypair in memory
func init() {
	var err error

	databaseInstance = &db.Db{}

	dbUser, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		fmt.Printf("environment variable POSTGRES_USER is not set")
		os.Exit(1)
	}
	dbPassword, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		fmt.Printf("environment variable POSTGRES_PASSWORD is not set")
		os.Exit(1)
	}
	dbHost, ok := os.LookupEnv("POSTGRES_HOST")
	if !ok {
		fmt.Printf("environment variable POSTGRES_HOST is not set")
		os.Exit(1)
	}
	dbPort, ok := os.LookupEnv("POSTGRES_PORT")
	if !ok {
		fmt.Printf("environment variable POSTGRES_PORT is not set")
		os.Exit(1)
	}
	dbDatabase, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		fmt.Printf("environment variable POSTGRES_DB is not set")
		os.Exit(1)
	}
	PORT, ok = os.LookupEnv("PORT")
	if !ok {
		log.Println("environment variable PORT is not set. Using default of 3000")
		PORT = "3000"
	}

	// Connection string for connecting to the database
	connStr := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbDatabase + "?sslmode=disable"

	databaseInstance.Connection, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("Could not open a connection to the postgres database. Please check that the environment variables POSTGRES_USER and POSTGRES_PASSWORD are exported and correct. Error: %s\n", err)
		os.Exit(1)
	}
	_, err = os.Stat("db.sql")
	if err != nil {
		fmt.Printf("Could not find the file db.sql that holds the database setup queries. Irrecoverable error. Exiting. Error: %s\n", err)
		os.Exit(1)
	}
	sqlToExecute, err := os.ReadFile("db.sql")
	if err != nil {
		fmt.Printf("Could not read the file db.sql that holds the database setup queries. Irrecoverable error. Exiting. Error: %s\n", err)
		os.Exit(1)
	}
	_, err = databaseInstance.Connection.Exec(string(sqlToExecute))
	if err != nil {
		fmt.Printf("Could not create database objects at initialization. Irrecoverable error. Exiting. Error: %s\n", err)
		os.Exit(1)
	}

	jwtSigningAuthority = &jwtauth.SigningAuthority{}

	err = jwtSigningAuthority.LoadKeys(privKeyPath, pubKeyPath)
	gtypes.Fatal(err)
}

func main() {

	r := gin.Default() // creates a router with Logger and Recovery middleware added

	// Auth required for the following REST calls
	authorizedGroup := r.Group("/")

	authorizedGroup.Use(jwtSigningAuthority.AuthRequired())
	{
		authorizedGroup.GET("/users", GetAllUsers)
		authorizedGroup.PUT("/users", UpdateUser)
	}

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "UP")
	})
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "UP")
	})
	r.POST("/signup", SignupUser)
	r.POST("/login", LoginUser)

	r.Run(":" + PORT)

}

// Handles the GET /users request and returns all users in the database. This function is invoked after the user is authenticated via the AuthRequired middleware
func GetAllUsers(ctx *gin.Context) {
	userRecords, err := databaseInstance.GetUsers()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "error fetching user records:" + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"users": userRecords})

}

// Handles the PUT /users request and updates the first and last names of the current user. This function is invoked after the user is authenticated via the AuthRequired middleware
func UpdateUser(ctx *gin.Context) {

	var updateUser gtypes.UserRecord

	err := ctx.BindJSON(&updateUser)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "error while binding:" + err.Error()})
		return
	}

	// Now tack on the email associated with the token to the record being updated
	updateUser.Email = ctx.GetString("email") // we had saved this in the context when authorizing the user

	err = databaseInstance.UpdateUser(updateUser)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "error updating user record:" + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Update successful"})

}

// This function handles the POST /signup request and returns a JWT if the signup is successful
func SignupUser(c *gin.Context) {
	var signupBody gtypes.Signup

	err := c.BindJSON(&signupBody)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "error while binding:" + err.Error()})
		return
	}

	insertedUser, err := databaseInstance.CreateUser(&signupBody)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error while creating user:" + err.Error()})
		return
	}

	// we will get an ID (primary key) back from the database if the record was created successfully
	if insertedUser > 0 { // record was inserted successfully and we got a primary key id back
		jwToken, err := jwtSigningAuthority.GenerateJWT(signupBody.Email)
		gtypes.Fatal(err)
		c.JSON(http.StatusOK, gin.H{"token": jwToken})
		return
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not validate user was created successfully"})
		return
	}
}

// This function logs in a user using their email and password. It returns a JWT if the email and password match in the database
func LoginUser(c *gin.Context) {
	var loginBody gtypes.Login

	err := c.BindJSON(&loginBody)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "error while binding:" + err.Error()})
		return
	}

	authUserId, err := databaseInstance.AuthenticateUser(loginBody)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error while logging in user:" + err.Error()})
		return
	}

	// we will get an ID (primary key) back from the database if the user was authenticated against the password
	if authUserId > 0 { // we got a primary key id back
		jwToken, err := jwtSigningAuthority.GenerateJWT(loginBody.Email)
		gtypes.Fatal(err)
		c.JSON(http.StatusOK, gin.H{"token": jwToken})
		return
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not login user"})
		return
	}
}
