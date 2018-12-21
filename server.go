package main

import (
  "github.com/gorilla/mux"
  "net/http"
  "log"
  "encoding/json"
  "sync"
)

// JSON Format
type Response struct {
  Message string `json:"message"`
  Status bool `json:"status"`
}

// dont use lists (bad for concurrency - security)
var Users = struct {
  sync.RWMutex
  m map[string] User
}{m: make(map[string] User)}

// user struct
type User struct {
  User_Name string
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

func Validate(w http.ResponseWriter, r *http.Request){
  // parsing form data
  r.ParserForm()
  user_name := r.FormValue("user_name")
}

func UserExist {}

func main(){

  // handling css
  cssHandle := http.FileServer(http.Dir("./css"))
  jsHandle := http.FileServer(http.Dir("./js"))
  // mux URL Router
  mux := mux.NewRouter()
  // HTTP Route
  mux.HandleFunc("/hello", Hello).Methods("GET")
  mux.HandleFunc("/hello-json", HelloJson).Methods("GET")
  mux.HandleFunc("/", LoadStatic).Methods("GET")
  mux.HandleFunc("/validate", Validate).Methods("POST")
  // we tell http module to use MUX Routes
  http.Handle("/", mux)
  // redirecting CSS routes
  http.Handle("/css/", http.StripPrefix("/css/", cssHandle))
  // redirecting JS routes
  http.Handle("/js/", http.StripPrefix("/js/", jsHandle))
  // message
  log.Println("Server running on port 8000")
  // starting server
  log.Fatal(http.ListenAndServe(":8000", nil))
}
