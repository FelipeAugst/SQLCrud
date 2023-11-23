package server

import (
	"encoding/json"
	"fmt"
	"io"
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Failed to load request body"))
		w.WriteHeader(http.StatusBadRequest)
		return

	}
	var user user

	if err := json.Unmarshal(body, &user); err != nil {
		w.Write([]byte("Failed to Unmarshal json"))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return

	}
	db, err := bank.Connect()
	if err != nil {
		w.Write([]byte("Failed to connect Db."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	statement, err := db.Prepare("insert into usuarios(name,email) values(?,?)")
	if err != nil {
		w.Write([]byte("Failed to prepare statement"))
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	defer statement.Close()

	result, err := statement.Exec(user.Name, user.Email)
	if err != nil {
		w.Write([]byte("Failed to execute statement"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		w.Write([]byte("Failed to recover last Id"))
		w.WriteHeader(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	db, err := bank.Connect()
	if err != nil {
		w.Write([]byte("Failed to load database"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()
	result, err := db.Query("select * from usuarios where id=?", ID)
	if err != nil {
		w.Write([]byte("Failed to perform SQL Query"))
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	var user user
	for result.Next() {
		error := result.Scan(&user.ID, &user.Name, &user.Email)
		if error != nil {
			w.Write([]byte("Failed to scan query results"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.Write([]byte("Failed to encode Json"))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

}

// Mostra todos os usuarios

func ShowUsers(w http.ResponseWriter, r *http.Request) {
	var users []user

	db, err := bank.Connect()
	if err != nil {
		w.Write([]byte("Fail to Connect Db"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()
	results, err := db.Query("select * from usuarios")
	if err != nil {

		w.Write([]byte("Error in SQL query()"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for results.Next() {
		var user user
		err := results.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			w.Write([]byte("Failed to scan query results"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		users = append(users, user)

	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

}

// Altera um usuario

func AlterUser(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	ID, err := strconv.ParseUint(parameters["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Failed to convert ID into uint"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Failed to read request body"))
		w.WriteHeader(http.StatusBadRequest)
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	statement, err := db.Prepare("update usuarios set name = ?,email= ? where id = ?")
	if err != nil {
		w.Write([]byte("Failed to create statement"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err := statement.Exec(user.Name, user.Email, ID); err != nil {
		w.Write([]byte("Failed to execute Sql Query"))
		w.WriteHeader(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := bank.Connect()
	if err != nil {
		w.Write([]byte("Failed to connect Db"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()
	statement, err := db.Prepare("delete from usuarios where id=?")
	if err != nil {
		w.Write([]byte("Failed to create Sql statement"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(ID); err != nil {
		w.Write([]byte("Failed to execute Sql query"))
		return
	}
	w.WriteHeader(http.StatusOK)
}
