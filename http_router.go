package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func init_router() {
	r := mux.NewRouter().StrictSlash(true)

	//GET requests
	r.HandleFunc("/hello", Hello).Methods("GET")                             //Welcome to the helpdesk
	r.HandleFunc("/list_users", ListUsers).Methods("GET")                    //List all the users
	r.HandleFunc("/list_issues", ListIssues).Methods("GET")                  //List all the issues
	r.HandleFunc("/user/{User_Id}", GetDetailsByID).Methods("GET")           //Get User Status
	r.HandleFunc("/issue_status/{Issue_Id}", ShowIssueStatus).Methods("GET") //Get issue status

	//POST requests
	r.HandleFunc("/new_user", HandleNewUSer).Methods("POST")    //Register as a new User
	r.HandleFunc("/post_issue", HandleNewIssue).Methods("POST") //Submit a Issue

	//PUT requests
	r.HandleFunc("/assing_customer_care", AssignCustomerCare).Methods("PUT") //Assing Customer Care to a Query
	r.HandleFunc("/update_issue_status", UpdateIssueStatus).Methods("PUT")   //Solve query and change status, description and update time

	log.Fatal(http.ListenAndServe(":1001", r))
}
