package user

import (
  "fmt"
  "strings"
  "errors"
  "database/sql"
  "golang.org/x/crypto/bcrypt"
  _ "github.com/lib/pq"
)

type DBUserData struct {
  Id         int         `json:"id"`
  Email      string      `json:"email"`
  Username   string  `json:"username"`
  Password   string      `json:"password"`
  Address    string  `json:"address"`
}

type User struct {
  Id         int         `json:"id"`
  Email      string  `json:"email"`
  Username   string  `json:"username"`
  Address    string  `json:"address"`
}

type AuthResponse struct {
  User
  Token string `json:"token"`
}

type PostPayload map[string]string

type NullString struct {
  sql.NullString
}

const (
  host   = "localhost"
  port   = 5432
  user   = "kyv"
  dbname = "miniapps"
)

func OpenConnection() *sql.DB {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)

  db, err := sql.Open("postgres", psqlInfo)

  if err != nil {
    panic(err)
  }

  return db
}

func GetAllUsers() ([]User, error) {
  db := OpenConnection()
  rows, err := db.Query("SELECT id, email, username, address FROM users")

  var users []User

  if err != nil {
    return users, err
  }

  for rows.Next() {
    var user User
    rows.Scan(&user.Id, &user.Email, &user.Username, &user.Address)
    users = append(users, user)
  }

  defer rows.Close()
  defer db.Close()

  return users, err
}

func UserExists(id string) bool {
  db := OpenConnection()
  rows, err := db.Query("SELECT id, email, username, address FROM users WHERE id = $1", id)

  if err != nil {
    return false
  }

  defer db.Close()

  return rows.Next()
}

func UpdateUser(id string, data PostPayload) (bool, error) {
  statement := "UPDATE users SET "
  var dataset []string
  for i, d := range data {
    dataset = append(dataset, fmt.Sprintf("%s = '%s'", i, d))
  }
  statement += strings.Join(dataset, ", ")
  statement += fmt.Sprintf("WHERE id = %s", id)

  db := OpenConnection()
  _, err := db.Query(statement)
  defer db.Close()

  if err != nil {
    return false, err
  }

  return true, err
}

func DeleteUser(id string) (bool, error) {
  db := OpenConnection()
  _, err := db.Query("DELETE FROM users WHERE id = $1", id)
  defer db.Close()

  if err != nil {
    return false, err
  }

  return true, err
}

func CreateUser(data PostPayload) (AuthResponse, error) {
  var column []string
  var insertValue []string
  for i, d := range data {
    column = append(column, i)
    insertValue = append(insertValue, fmt.Sprintf("'%s'", d))
  }

  var user AuthResponse

  db := OpenConnection()
  statement := "INSERT INTO users (" + strings.Join(column, ", ") + ") VALUES (" + strings.Join(insertValue, ", ") + ") RETURNING id, email, username"
  err := db.QueryRow(statement).Scan(&user.Id, &user.Email, &user.Username)
  defer db.Close()

  return user, err
}

func GetUserByUsername(username string) (DBUserData, error) {
  var user DBUserData

  db := OpenConnection()
  rows, err := db.Query("SELECT id, email, username, password, address FROM users WHERE username = $1", username)
  if err != nil {
    return user, err
  }

  if rows.Next() {
    rows.Scan(&user.Id, &user.Email, &user.Username, &user.Password, &user.Address)
  }

  defer rows.Close()
  defer db.Close()

  if user.Id == 0 {
    return user, errors.New("User doesn't exists.")
  } else {
    return user, err
  }
}

func (u *DBUserData) CheckPassword(providedPassword string) error {
  err := bcrypt.CompareHashAndPassword([]byte(string(u.Password)), []byte(providedPassword))

  if err != nil {
    return err
  }

  return nil
}

func GetAuthResponse(username string) (AuthResponse, error) {
  var user AuthResponse

  db := OpenConnection()
  rows, err := db.Query("SELECT id, email, username, address FROM users WHERE username = $1", username)
  if err != nil {
    return user, err
  }

  if rows.Next() {
    rows.Scan(&user.Id, &user.Email, &user.Username, &user.Address)
  }

  defer rows.Close()
  defer db.Close()

  if user.Id == 0 {
    return user, errors.New("User doesn't exists.")
  } else {
    return user, err
  }
}
