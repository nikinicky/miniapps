# miniapps
A simple API. User can create, edit, and delete profiles. For now, this code only run localy.
First, you need to create database and table user. I provide query to create table user below.
```sql
CREATE TABLE public.users (
	id serial NOT NULL DEFAULT nextval('users_id_seq1'::regclass),
	email varchar NULL,
	username varchar NOT NULL,
	"password" varchar NOT NULL,
	address varchar NULL,
	CONSTRAINT users_pkey PRIMARY KEY (id),
	CONSTRAINT users_username_key UNIQUE (username)
);
```

___

## Endpoints
* `POST /signup` create a new user and get token to access others endpoint
* `POST /login` get token 
* `GET /users` returns list of users as JSON (requeire token)
* `GET /users/{id}` get a user detail (require token)
* `PUT /users/{id}` update an users (require token)
* `DELETE /users/{id}` delete user record (require token)


## Data Types
A user object should look like this
```json
{
    "id": 1,
    "email": "example@gmail.com",
    "username": "example",
    "address": "Sudirman Street"
}
```

## Example
### Signup
request
```
curl --request POST \
  --url http://localhost:8080/signup \
  --header 'Content-Type: application/json' \
  --data '{
	  "username": "example2",
	  "password": "youknowmypassword",
	  "email": "example2@gmail.com"
  }'
```
response
```json
{
  "id": 79,
  "email": "example2@gmail.com",
  "username": "example2",
  "address": "",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6ImV4YW1wbGUyQGdtYWlsLmNvbSIsImV4cCI6MTYxMTU1Mjg1NywiaXNzIjoiQXV0aFNlcnZpY2UifQ.J7yKDDZvizddaXLI_LlZrVWqH2xRgbuZPTlRkAJvlbg"
}
```

### Login
request
```
curl --request POST \
  --url http://localhost:8080/login \
  --header 'Content-Type: application/json' \
  --data '{
	  "username": "example2",
	  "password": "youknowmypassword"
  }'
```
response
```json
{
  "id": 79,
  "email": "example2@gmail.com",
  "username": "example2",
  "address": "",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6ImV4YW1wbGUyQGdtYWlsLmNvbSIsImV4cCI6MTYxMTU1Mjk3NywiaXNzIjoiQXV0aFNlcnZpY2UifQ.PxsSvfRD8OizUvW1LLWN-T8nowx4EIUIGecEohXOS2E"
}
```

### Get All Users
request _(require token)_
```
curl --request GET \
  --url http://localhost:8080/users \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6ImV4YW1wbGUyQGdtYWlsLmNvbSIsImV4cCI6MTYxMTU1Mjk3NywiaXNzIjoiQXV0aFNlcnZpY2UifQ.PxsSvfRD8OizUvW1LLWN-T8nowx4EIUIGecEohXOS2E'
```
response
```json
[
  {
    "id": 79,
    "email": "example2@gmail.com",
    "username": "example2",
    "address": ""
  },
  {
    "id": 60,
    "email": "example@gmail.com",
    "username": "example",
    "address": ""
  }
]
```

### Get A User Detail
request _(require token)_
```
curl --request GET \
  --url http://localhost:8080/users/79 \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6ImV4YW1wbGUyQGdtYWlsLmNvbSIsImV4cCI6MTYxMTU1Mjk3NywiaXNzIjoiQXV0aFNlcnZpY2UifQ.PxsSvfRD8OizUvW1LLWN-T8nowx4EIUIGecEohXOS2E'
```
response
```json
{
  "id": 79,
  "email": "example2@gmail.com",
  "username": "example2",
  "address": ""
}
```

### Update User
request _(require token)_
```
curl --request PUT \
  --url http://localhost:8080/users/79 \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6ImhlbGxvQGdtYWlsLmNvbSIsImV4cCI6MTYxMTUxMzkxMiwiaXNzIjoiQXV0aFNlcnZpY2UifQ.FTxkbajHj7MwLwkKx3wC18YlqbW6C0bms-rvPHetZMY' \
  --header 'Content-Type: application/json' \
  --data '{
	  "email": "example2@gmail.com",
	  "username": "example2",
	  "address": "Diponegoro Street"
  }'
```
response
```json
{
  "message": "Successfully udpate user!",
  "status": "OK"
}
```

### Delete User
request _(require token)_
```
curl --request DELETE \
  --url http://localhost:8080/users/79 \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6ImhlbGxvQGdtYWlsLmNvbSIsImV4cCI6MTYxMTUxMzkxMiwiaXNzIjoiQXV0aFNlcnZpY2UifQ.FTxkbajHj7MwLwkKx3wC18YlqbW6C0bms-rvPHetZMY' \
  --header 'Content-Type: application/json' \
  --cookie '__profilin=p%253Dt; _sweetescape_session=2de809db56fb64c55344f8feb942d64f; _fotto_rails2_session=Dsmr1PnezgjkCWqgmJn0MXYx2T9AO9833F9k92IcNKpbg1xJjQZPPdPqfA57mas3dTR%252BbcyoaDJLwgvbswWhf%252Fh%252BScowb0sJ9%252BZW5nYjTJkFzIcic0BGXQQo84dePbHep9lEZ8TBQ%252Fc%252Ftxkhssq%252BwOiFan0ZwgscNP0LjxqV1RS%252F%252F66tfR1st%252FaW1icIj1fDLFAgd4h746fVJ6uQsA%253D%253D--NekFlWjSwATgLZ6Q--WYnqUMXrRD7mYrkvIV0HUg%253D%253D'
```
response
```json
{
  "message": "Successfully delete user!",
  "status": "OK"
}

```
