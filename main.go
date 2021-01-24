package main

import (
  "log"
  "net/http"
  "github.com/gorilla/mux"
  "nickyv/miniapps/handler"
)

func main() {
  router := mux.NewRouter()

  router.HandleFunc("/users", handler.GetAllUsers).Methods("GET")
  router.HandleFunc("/update/{id}", handler.UpdateUser).Methods("PUT")
  router.HandleFunc("/delete/{id}", handler.DeleteUser).Methods("DELETE")
  router.HandleFunc("/signup", handler.Signup).Methods("POST")
  router.HandleFunc("/login", handler.Login).Methods("POST")

  log.Fatal(http.ListenAndServe(":8080", router))
}
