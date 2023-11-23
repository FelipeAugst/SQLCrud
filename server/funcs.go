package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mycrud/bank"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type user struct {
	ID    int64  `json: "id"`
	Name  string `json: "name" `
	Email string `json: "email" `
}

// Criacao de um novo usuario

func CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Failed to load request body"))
		return

	}
	var user user

	if err := json.Unmarshal(body, &user); err != nil {
		w.Write([]byte("Failed to Unmarshal json"))
		return

	}
	db, err := bank.Connect()
	if err != nil {
		w.Write([]byte("Failed to connect Db."))
		return
	}

	defer db.Close()

	statement, err := db.Prepare("insert into usuarios(name,email) values(?,?)")
	if err != nil {
		w.Write([]byte("Failed to prepare statement"))
		return

	}
	defer statement.Close()

	result, err := statement.Exec(user.Name, user.Email)
	if err != nil {
		w.Write([]byte("Failed to execute statement"))
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		w.Write([]byte("Failed to recover last Id"))
		return
	}
	w.Write([]byte(fmt.Sprintf("User %d created", id)))
	w.WriteHeader(http.StatusCreated)
}

// Busca por um usuario

func SearchUser(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	ID, err := strconv.ParseUint(params["id"], 10, 32)

	if err != nil {
		w.Write([]byte("Failed to convert ID"))
		return
	}
	db, err := bank.Connect()
	if err != nil {
		w.Write([]byte("Failed to load database"))
		return
	}
	defer db.Close()
	result, err := db.Query("select * from usuarios where id=?", ID)
	if err != nil {
		w.Write([]byte("Failed to perform SQL Query"))
		return

	}
	var user user
	for result.Next() {
		error := result.Scan(&user.ID, &user.Name, &user.Email)
		if error != nil {
			w.Write([]byte("Failed to scan query results"))
			return
		}
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.Write([]byte("Failed to encode Json"))
		return
	}

}

// Mostra todos os usuarios

func ShowUsers(w http.ResponseWriter, r *http.Request) {
	var users []user

	db, err := bank.Connect()
	if err != nil {
		w.Write([]byte("Fail to Connect Db"))
		return
	}
	defer db.Close()
	results, err := db.Query("select * from usuarios")
	if err != nil {

		w.Write([]byte("Error in SQL query()"))
		return
	}

	for results.Next() {
		var user user
		err := results.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			w.Write([]byte("Failed to scan query results"))
			return
		}

		users = append(users, user)

	}

	json.NewEncoder(w).Encode(users)

}

// Altera um usuario

func AlterUser(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	ID, err := strconv.ParseUint(parameters["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Failed to convert ID into uint"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Failed to read request body"))
		return
	}
	var user user
	err = json.Unmarshal(body, &user)
	if err != nil {
		w.Write([]byte("Failed to unmarshal json"))
		return
	}
	db, err := bank.Connect()
	if err != nil {
		w.Write([]byte("Failed to connect Db"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("update usuarios set name = ?,email= ? where id = ?")
	if err != nil {
		w.Write([]byte("Failed to create statement"))
		return
	}
	if _, err := statement.Exec(user.Name, user.Email, ID); err != nil {
		w.Write([]byte("Failed to execute Sql Query"))
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Deletando Usuarios

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	ID, err := strconv.ParseUint(parameters["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Failed to convert ID"))
		return
	}

	db, err := bank.Connect()
	if err != nil {
		w.Write([]byte("Failed to connect Db"))
		return
	}
	defer db.Close()
	statement, err := db.Prepare("delete from usuarios where id=?")
	if err != nil {
		w.Write([]byte("Failed to create Sql statement"))
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(ID); err != nil {
		w.Write([]byte("Failed to execute Sql query"))
		return
	}
	w.WriteHeader(http.StatusOK)
}
