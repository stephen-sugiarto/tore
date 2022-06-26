package main

import (
	"database/sql"  //Go core package untuk handling komunikasi sql
	"encoding/json" //Go core package untuk handling JSON
	"fmt"           //Untuk print message
	"io/ioutil"
	"log"      //untuk logging error
	"net/http" //Go HTTP package untuk handling hTTP request saat API Go di hit

	"github.com/gorilla/mux" //untuk URL match dan routing. berguna untuk implement request route dan match tiap kali API di hit
	_ "github.com/lib/pq"    //untuk handling DB/SQL package
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "admin"
	DB_NAME     = "tore"
)

func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)

	checkErr(err)

	return db
}

//fect json data
type UserLogin struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type JsonResponse struct {
	Type    string      `json:"type"`
	Data    []UserLogin `json:"data"`
	Message string      `json:"message"`
}

//main function to endpoint
func main() {
	// init the mux router
	router := mux.NewRouter()

	//route handles and endpoint

	//get ALL username
	router.HandleFunc("/userLogin", getAllUser).Methods("GET")

	//create new user
	router.HandleFunc("/newUser", setUser).Methods("POST")

	//delete user
	router.HandleFunc("/deleteUser/{id}", deleteUser).Methods("DELETE")

	//delete ALL user
	router.HandleFunc("/deleteAllUser", deleteAllUser).Methods("DELETE")

	//serve the app
	fmt.Println("Server at 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

//function to handling messages
func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Print("")
}

//function to handling errors
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// get all user
//response dan request handler
func getAllUser(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("getting all users")

	//get all user from table
	rows, err := db.Query("select * from public.user")

	//check errors
	checkErr(err)

	// var response []json response
	var listUserLogin []UserLogin
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	//foreach user
	for rows.Next() {
		var id int
		var username string
		var password string
		var name string

		err = rows.Scan(&id, &username, &password, &name)

		//check error
		checkErr(err)

		listUserLogin = append(listUserLogin, UserLogin{Id: id, Username: username, Password: password, Name: name})
	}

	var response = JsonResponse{Type: "success", Data: listUserLogin}

	json.NewEncoder(w).Encode(response)
}

//create a user
//response and request handlers
func setUser(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var userLogin UserLogin
	json.Unmarshal(reqBody, &userLogin)

	// r.ParseForm()
	// username := r.FormValue("username")
	// password := r.FormValue("password")
	// name := r.FormValue("name")

	username := userLogin.Username
	password := userLogin.Password
	name := userLogin.Name

	var response = JsonResponse{}

	if username == "" || password == "" || name == "" {
		printMessage("Username: " + username)
		printMessage("Password: " + password)
		printMessage("Name: " + name)
		response = JsonResponse{Type: "error", Message: "Missing username, password, or name parameter"}
	} else {
		db := setupDB()

		printMessage("Insert users into DB")

		fmt.Println("Insert new user with username: " + username + " and name: " + password)

		var lastInsertID int

		err := db.QueryRow("INSERT INTO public.user(username, password, name) VALUES ($1, $2, $3) returning id;", username, password, name).Scan(&lastInsertID)

		checkErr(err)

		response = JsonResponse{Type: "success", Message: "User has been inserted successfully"}
	}

	json.NewEncoder(w).Encode(response)
}

//delete single records
func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]

	var response = JsonResponse{}

	if id == "" {
		response = JsonResponse{Type: "error", Message: "No id provided"}
	} else {
		db := setupDB()

		printMessage("Deleting id from user...")

		_, err := db.Exec("DELETE FROM public.user WHERE id = $1", id)

		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The user was deleted"}
	}

	json.NewEncoder(w).Encode(response)
}

//function delete all users
func deleteAllUser(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Deleting all users")

	_, err := db.Exec("DELETE FROM public.user")

	//check error
	checkErr(err)

	var response = JsonResponse{Type: "success", Message: "All users have been deleted"}

	json.NewEncoder(w).Encode(response)
}
