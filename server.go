package main

import (
  "github.com/gorilla/mux"
  "github.com/gorilla/websocket"
  "net/http"
  "log"
  "encoding/json"
  "sync"
)

// JSON Format
type Response struct {
  Message string `json:"message"`
  Status int `json:"status"`
  IsValid bool  `json:"isvalid"`
}

// dont use lists (bad for concurrency - security)
var Users = struct {
  sync.RWMutex
  m map[string] User
}{m: make(map[string] User)}

// user struct
type User struct {
  User_Name string
  WebSocket *websocket.Conn
}

func Hello(w http.ResponseWriter, r *http.Request){
  w.Write([]byte ("Hello, World"))
}

func HelloJson(w http.ResponseWriter, r *http.Request){
  // create a Response Struct
  response := CreateResponse("This is json format", 200, true)
  // Encoding response
  json.NewEncoder(w).Encode(response)
}

func CreateResponse(message string, status int, valid bool) Response {
  return Response{message, status, valid}
}

func LoadStatic(w http.ResponseWriter, r *http.Request){
  // serving HTML File
  http.ServeFile(w,r, "./index.html")
}

func Validate(w http.ResponseWriter, r *http.Request){
  // parsing form data
  r.ParseForm()
  user_name := r.FormValue("user_name")
  response := Response{}
  // checking existance
  if UserExist(user_name) {
    // denied
    response.IsValid = false
    // response := CreateResponse("User not valid", false)
  } else {
    // granted
    response.IsValid = true
    // response := CreateResponse("User valid", true)
  }
  json.NewEncoder(w).Encode(response)
}

func UserExist(user_name string) bool{
  // while validating, it blocks the structure
  Users.RLock()
  defer Users.RUnlock()
  // validating existance
  if _,ok := Users.m[user_name]; ok {
    return true
  }
  return false
}

func CreateUser(user_name string, ws *websocket.Conn) User {
  // going to the map
  return User {user_name, ws}
}

func AddUser(user User){
  Users.Lock()
  defer Users.Unlock()
  Users.m[user.User_Name] = user
}

func RemoveUser(user_name string){
  Users.Lock()
  defer Users.Unlock()
  delete(Users.m, user_name)
}

func SendMessage(type_message int, message []byte){
  Users.RLock()
  defer Users.RUnlock()
  for _, user := range Users.m{
    if err := user.WebSocket.WriteMessage(type_message, message); err != nil {
      return
    }
  }
}

func ToArrayByte(value string) []byte{
  return []byte(value)
}

func ConcatMessage(user_name string, array[]byte) string {
  return user_name + ":" + string(array[:])
}

func WebSocket(w http.ResponseWriter, r *http.Request){
  // param
  vars := mux.Vars(r)
  user_name := vars["user_name"]
  // creating websocket
  ws,err := websocket.Upgrade(w,r,nil,1024,1024)
  if err != nil {
    log.Println(err)
    return
  }
  // getting current user
  current_user := CreateUser(user_name,ws)
  AddUser(current_user)
  log.Println("New User created")

  for {
    type_message, message,err := ws.ReadMessage()
    if err != nil {
      RemoveUser(user_name)
      return
    }
    final_message := ConcatMessage(user_name, message)
    SendMessage(type_message, ToArrayByte(final_message))
  }
}


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
  mux.HandleFunc("/chat/{user_name}", WebSocket)
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
