package db

import (
	"database/sql"

	"github.com/dapper-labs-talent/cc_ashishchandra_BackendAPI/gtypes"
	_ "github.com/lib/pq"
)

type Db struct {
	Connection *sql.DB
}

// creates a user record in the database and encrypts the password for storage
func (db *Db) CreateUser(user *gtypes.Signup) (int64, error) {
	result, err := db.Connection.Query(`INSERT INTO users(email,password,first_name,last_name) VALUES($1,crypt($2,gen_salt('bf')),$3,$4) RETURNING ID`, user.Email, user.Password, user.FirstName, user.LastName)
	if err != nil {
		return -1, err
	}
	defer result.Close()

	var lastId int64

	for result.Next() {
		err := result.Scan(&lastId)
		if err != nil {
			return -1, err
		}
	}
	return lastId, nil
}

// Authenticates the password provided with the login against the password provided at signup
func (db *Db) AuthenticateUser(login gtypes.Login) (int64, error) {
	result, err := db.Connection.Query(
		`SELECT id FROM users WHERE email = $1 AND password = crypt($2, password)`, login.Email, login.Password)
	if err != nil {
		return -1, err
	}
	defer result.Close()

	var userId int64

	for result.Next() {
		err := result.Scan(&userId)
		if err != nil {
			return -1, err
		}
	}
	return userId, nil
}

// Updates the user's record (first name and last name) in the database. The record updated is the one associated with the email address associated with the provided JWT
func (db *Db) UpdateUser(user gtypes.UserRecord) error {
	result, err := db.Connection.Query(
		`UPDATE users SET first_name = $1, last_name = $2 WHERE email = $3 RETURNING first_name,last_name`, user.FirstName, user.LastName, user.Email)
	if err != nil {
		return err
	}
	defer result.Close()

	var fName, lName string

	for result.Next() {
		err := result.Scan(&fName, &lName)
		if err != nil || fName != user.FirstName || lName != user.LastName {
			return err
		}
	}
	return nil
}

// Gets a list of all the users in the database. This function is also called after the user is authorized
func (db *Db) GetUsers() (*[]gtypes.UserRecord, error) {

	retValArray := make([]gtypes.UserRecord, 0)

	result, err := db.Connection.Query(
		`SELECT email,first_name,last_name from users`)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	for result.Next() {
		retVal := gtypes.UserRecord{}
		err := result.Scan(&retVal.Email, &retVal.FirstName, &retVal.LastName)
		if err != nil {
			return nil, err
		}
		retValArray = append(retValArray, retVal)
	}
	return &retValArray, nil
}
