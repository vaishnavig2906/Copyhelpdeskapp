package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vaishnavi2906/helpdesk/db"
)

type Service interface {
	ListUsers(res http.ResponseWriter, req *http.Request)
	GetDetailsByID(res http.ResponseWriter, req *http.Request)
	HandleNewUSer(res http.ResponseWriter, req *http.Request)
	Hello(res http.ResponseWriter, req *http.Request)
}

func Hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Hello Welcome to helpdesk")
}

func ListUsers(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Complete User Data")
	var UserInfo User

	db, err := db.Init_DB()
	if err != nil {
		log.Fatal(err)
		return
	}

	query := `SELECT * from "user";`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
		return
	}

	for rows.Next() {
		err := rows.Scan(&UserInfo.Id, &UserInfo.Email, &UserInfo.Usertype)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println("\n", UserInfo.Id, UserInfo.Email, UserInfo.Usertype)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return
	}
	db.Close()
}

//Handle new User syntax
// {
//     "id":"1",
//     "email":"a@example.com",
//     "type":"Customer"
// }
func HandleNewUSer(res http.ResponseWriter, req *http.Request) {
	var UserDetails UserRequest

	err := json.NewDecoder(req.Body).Decode(&UserDetails)
	if err != nil {
		return
	}

	db, err := db.Init_DB()
	if err != nil {
		log.Fatal(err)
		return
	}

	ctx := req.Context()

	query := `INSERT INTO public."user"(id, email, user_type) VALUES ($1,$2,$3);`

	_, err = db.ExecContext(ctx, query, UserDetails.Id, UserDetails.Email, UserDetails.Usertype)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Fprintf(res, "Succesfully Registered")
	db.Close()
}

func GetDetailsByID(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["User_Id"]

	db, err := db.Init_DB()
	if err != nil {
		log.Fatal(err)
		return
	}
	ctx := req.Context()
	query := `Select * FROM "user" WHERE id=$1;`
	var UserInfo User

	err = db.GetContext(ctx, &UserInfo, query, id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Given the details to the user")
	fmt.Fprintf(res, "email: %v\n", UserInfo.Email)
	db.Close()
}
