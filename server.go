package main

import (
  "github.com/gorilla/mux"
  "net/http"
  "log"
)

func Hello(w http.ResponseWriter, r *http.Request){
  w.Write([]byte ("Hello, World"))
}

func main(){
  // mux URL Router
  mux := mux.NewRouter()
  // HTTP Route
  mux.HandleFunc("/hello", Hello).Methods("GET")
  // we tell http module to use MUX Routes
  http.Handle("/", mux)
  // message
  log.Println("Server running on port 8000")
  // starting server
  log.Fatal(http.ListenAndServe(":8000", nil))
}
