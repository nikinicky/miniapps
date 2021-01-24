package main

import (
  "log"
  "net/http"
  "github.com/gorilla/mux"
  "nickyv/miniapps/handler"
)

func main() {
  router := mux.NewRouter()

  router.HandleFunc("/signup", handler.Signup).Methods("POST")
  router.HandleFunc("/login", handler.Login).Methods("POST")
  router.HandleFunc("/users", handler.GetAllUsers).Methods("GET")
  router.HandleFunc("/users/{id}", handler.GetUser).Methods("GET")
  router.HandleFunc("/users/{id}", handler.UpdateUser).Methods("PUT")
  router.HandleFunc("/users/{id}", handler.DeleteUser).Methods("DELETE")

  log.Fatal(http.ListenAndServe(":8080", router))
}
