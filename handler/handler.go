package handler

import (
  "fmt"
  "strings"
  "net/http"
  "database/sql"
  "encoding/json"
  "io/ioutil"
  "github.com/gorilla/mux"
  "golang.org/x/crypto/bcrypt"
  "nickyv/miniapps/auth"
  "nickyv/miniapps/user"
)

type NullString struct {
  sql.NullString
}

func (ns NullString) MarshalText() ([]byte, error) {
  if !ns.Valid {
    return []byte(nil), nil
  }
  return []byte(ns.String), nil
}

func CheckToken(r *http.Request) bool {
  authHeader := r.Header.Get("Authorization")
  if authHeader == "" {
    return false
  }

  extractedToken := strings.Split(authHeader, "Bearer ")

  var token string
  if len(extractedToken) == 2 {
    token = strings.TrimSpace(extractedToken[1])
  } else {
    return false
  }

  jwtWrapper := auth.JwtWrapper{
    SecretKey: "mysecretkey",
    Issuer: "AuthService",
  }

  _, err := jwtWrapper.ValidateToken(token)

  if err != nil {
    return false
  }

  return true
}

func GenerateToken(email string) (string, error) {
  jwtWrapper := auth.JwtWrapper{
    SecretKey: "mysecretkey",
    Issuer: "AuthService",
    ExpirationHours: 24,
  }

  signedToken, err := jwtWrapper.GenerateToken(email)

  if err != nil {
    return "", err
  }

  return signedToken, nil
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
  validToken := CheckToken(r)

  if !validToken {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }

  users, err := user.GetAllUsers()
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(err.Error()))
    return
  }

  usersBytes, err := json.Marshal(users)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(err.Error()))
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  w.Write(usersBytes)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
  validToken := CheckToken(r)

  if !validToken {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }

  params := mux.Vars(r)
  userExists := user.UserExists(params["id"])

  if !userExists {
    w.WriteHeader(http.StatusUnprocessableEntity)
    return
  }

  bodyBytes, err := ioutil.ReadAll(r.Body)
  defer r.Body.Close()

  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte(err.Error()))
    return
  }

  var data user.PostPayload

  err = json.Unmarshal(bodyBytes, &data)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte(err.Error()))
    return
  }

  _, err = user.UpdateUser(params["id"], data)

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(err.Error()))
    return
  }
  var res = map[string]string {
    "status": "OK",
    "message": "Successfully udpate user!",
  }

  resBytes, err := json.Marshal(res)

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(err.Error()))
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  w.Write(resBytes)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
  validToken := CheckToken(r)

  if !validToken {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }

  params := mux.Vars(r)
  _, err := user.DeleteUser(params["id"])

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(err.Error()))
    return
  }

  var res = map[string]string {
    "status": "OK",
    "message": "Successfully delete user!",
  }

  resBytes, err := json.Marshal(res)

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(err.Error()))
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  w.Write(resBytes)
}

func Signup(w http.ResponseWriter, r *http.Request) {
  contentType := r.Header.Get("Content-Type")
  if contentType != "application/json" {
    w.WriteHeader(http.StatusUnsupportedMediaType)
    w.Write([]byte(fmt.Sprintf("Need content-type 'application/json', but got '%s'", contentType)))
    return
  }

  bodyBytes, err := ioutil.ReadAll(r.Body)
  defer r.Body.Close()

  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte(err.Error()))
    return
  }

  var data user.PostPayload
  err = json.Unmarshal(bodyBytes, &data)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte(err.Error()))
    return
  }

  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 8)
  data["password"] = string(hashedPassword)

  user, err := user.CreateUser(data)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(err.Error()))
    return
  }

  user.Token, err = GenerateToken(user.Email)
  userBytes, err := json.Marshal(user)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(err.Error()))
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  w.Write(userBytes)
}

func Login(w http.ResponseWriter, r *http.Request) {
  contentType := r.Header.Get("Content-Type")
  if contentType != "application/json" {
    w.WriteHeader(http.StatusUnsupportedMediaType)
    w.Write([]byte(fmt.Sprintf("Need content-type 'application/json', but got '%s'", contentType)))
    return
  }

  bodyBytes, err := ioutil.ReadAll(r.Body)
  defer r.Body.Close()
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte(err.Error()))
    return
  }

  var data user.PostPayload
  err = json.Unmarshal(bodyBytes, &data)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte(err.Error()))
    return
  }

  userDB, err := user.GetUserByUsername(data["username"])
  if err != nil {
    w.WriteHeader(http.StatusUnprocessableEntity)
    w.Write([]byte(err.Error()))
    return
  }

  err = userDB.CheckPassword(data["password"])
  if err != nil {
    w.WriteHeader(http.StatusUnprocessableEntity)
    w.Write([]byte(err.Error()))
    return
  }

  user, err := user.GetAuthResponse(data["username"])
  user.Token, err = GenerateToken(user.Email)
  userBytes, err := json.Marshal(user)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(err.Error()))
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  w.Write(userBytes)
}
