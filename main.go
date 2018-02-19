package main

import (
  "database/sql"
  "html/template"
  "log"
  "net/http"

  _ "github.com/go-sql-driver/mysql"
  bcrypt "golang.org/x/crypto/bcrypt"
)


var db *sql.DB
var err error
var tpl *template.Template


// users schema
type user struct {
  ID        int64
  Username  string
  Firstname string
  Lastname  string
  Password  []byte
}


// Handler Function

// Initialize function
func init() {
  tpl = template.Must(template.ParseGlob("templates/*"))
}


// Index func
func index(w http.ResponseWriter, req *http.Request) {
  rows, e := db.Query(
    `SELECT id,
     username,
     first_name,
     last_name,
     password
     FROM users;`
  )
  if e != nil {
    log.Println(e)
    http.Error(w, e.Error(), http.StatusInternalServerError)
    return
  }
  users := make([]user, 0)
  for rows.Next() {
    usr := user{}
    rows.Scan(&usr.ID, &usr.Username, &usr.Firstname, &usr.Lastname, &usr.Password)
    users = append(users, usr)
  }
  log.Println(users)
  tpl.ExecuteTemplate(w, "index.gohtml", users)
}


// UserForm func
func userForm(w http.ResponseWriter, req *http.Request) {
  err = tpl.ExecuteTemplate(w, "userForm.gohtml", nil)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}


// createUsers func
func createUsers(w http.ResponseWriter, req *http.Request) {}


// editUsers func
func editUsers(w http.ResponseWriter, req *http.Request) {}


// deleteUsers func
func deleteUsers(w http.ResponseWriter, req *http.Request) {}


// updateUsers func
func updateUsers(w http.ResponseWriter, req *http.Request) {}

func main() {
  defer db.Close()
  http.HandleFunc("/", index)
  http.HandleFunc("/userForm", userForm)
  http.HandleFunc("/createUsers", createUsers)
  http.HandleFunc("/editUsers", editUsers)
  http.HandleFunc("/deleteUsers", deleteUsers)
  http.HandleFunc("/updateUsers", updateUsers)
  log.Println("Server is up on port 8080")
  log.Fatalln(http.listenAndServe(":8080", nil))
}

