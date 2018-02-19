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
     FROM users;`)
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
func createUsers(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		usr := user{}
		usr.Username = req.FormValue("username")
		usr.FirstName = req.FormValue("firstName")
		usr.LastName = req.FormValue("lastName")
		bPass, e := bcrypt.GenerateFromPassword([]byte(req.FormValue("password")), bcrypt.MinCost)
		if e != nil {
			http.Error(w, e.Error(), http.StatusInternalServerError)
			return
		}
		usr.Password = bPass
		_, e = db.Exec(
			"INSERT INTO users (username, first_name, last_name, password) VALUES (?, ?, ?, ?,)",
			usr.Username,
			usr.FirstName,
			usr.LastName,
			usr.Password,
		)
		if e != nil {
			http.Error(w, e.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, req, '/', http.StatusSeeOther)
		return
	}
	http.Error(w, "Method Not Supported", http.StatusMethodNotAllowed)
}

// editUsers func
func editUsers(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	rows, err := db.Query(
		`SELECT id,
				username,
				first_name,
				last_name
		FROM users
		WHERE id = ` + id + `;`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	usr := user{}
	for rows.Next() {
		rows.Scan(&usr.ID, &usr.Username, &usr.FirstName, &usr.LastName)
	}
	tpl.ExecuteTemplate(w, "editUser.gohtml", usr)
}

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
