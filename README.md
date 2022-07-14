## Assignment Summary
The focus of this assignment is to see how you would set up an API specifically for a web application.

The majority of our apps require an user login flow. This assignment will cover the signup, login endpoints as well as an additional resource `users` that can be accessed if you are logged in.

To determine whether a request is from a logged in user or not, we'll be using Json Web Tokens (https://jwt.io/). The frontend will be sending requests with the JWT in the `x-authentication-token` header.

For the database, we like to start projects off using PostgreSQL. Feel free to use something like https://github.com/go-pg/pg to get yourself started.

Also feel free to use whatever open source packages you're comfortable with.

## API Specs

### `POST /signup`
Endpoint to create an user row in postgres db. The payload should have the following fields:

```json
{
  "email": "test@axiomzen.co",
  "password": "axiomzen",
  "firstName": "Alex",
  "lastName": "Zimmerman"
}
```

where `email` is an unique key in the database.

The response body should return a JWT on success that can be used for other endpoints:

```json
{
  "token": "some_jwt_token" 
}
```

### `POST /login`
Endpoint to log an user in. The payload should have the following fields:

```json
{
  "email": "test@axiomzen.co",
  "password": "axiomzen"
}
```

The response body should return a JWT on success that can be used for other endpoints:

```json
{
  "token": "some_jwt_token"
}
```

### `GET /users`
Endpoint to retrieve a json of all users. This endpoint requires a valid `x-authentication-token` header to be passed in with the request.

The response body should look like:
```json
{
  "users": [
    {
      "email": "test@axiomzen.co",
      "firstName": "Alex",
      "lastName": "Zimmerman"
    }
  ]
}
```

### `PUT /users`
Endpoint to update the current user `firstName` or `lastName` only. This endpoint requires a valid `x-authentication-token` header to be passed in and it should only update the user of the JWT being passed in. The payload can have the following fields:

```json
{
  "firstName": "NewFirstName",
  "lastName": "NewLastName"
}
```

The response can body can be empty.

## Last notes
Create an INSTRUCTIONS.md file if there are any special steps we need to take to run your assignment.

Please reach out if you have any questions regarding the assignment! We're happy to help.

Try to spend approximately 4 hours on this assignment. Use your best judgement when estimating time spent.

Good luck!
