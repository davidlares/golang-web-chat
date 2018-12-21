package main

import (
  "github.com/gorilla/mux"
  "net/http"
  "log"
  "encoding/json"
)

// JSON Format
type Response struct {
  Message string `json:"message"`
  Status bool `json:"status"`
}

func Hello(w http.ResponseWriter, r *http.Request){
  w.Write([]byte ("Hello, World"))
}

func HelloJson(w http.ResponseWriter, r *http.Request){
  // create a Response Struct
  response := CreateResponse()
  // Encoding response
  json.NewEncoder(w).Encode(response)
}

func CreateResponse() Response {
  return Response{"This is json format", true}
}

func LoadStatic(w http.ResponseWriter, r *http.Request){
  // serving HTML File
  http.ServeFile(w,r, "./index.html")
}

func main(){

  // handling css
  cssHandle := http.FileServer(http.Dir("./css"))

  // mux URL Router
  mux := mux.NewRouter()
  // HTTP Route
  mux.HandleFunc("/hello", Hello).Methods("GET")
  mux.HandleFunc("/hello-json", HelloJson).Methods("GET")
  mux.HandleFunc("/", LoadStatic).Methods("GET")
  // we tell http module to use MUX Routes
  http.Handle("/", mux)
  // redirecting CSS routes
  http.Handle("/css/", http.StripPrefix("/css/", cssHandle))
  // message
  log.Println("Server running on port 8000")
  // starting server
  log.Fatal(http.ListenAndServe(":8000", nil))
}
