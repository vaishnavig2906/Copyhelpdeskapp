package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vaishnavi2906/helpdesk/issue"
	"github.com/vaishnavi2906/helpdesk/user"
)

func init_router() {
	r := mux.NewRouter().StrictSlash(true)

	db := InitDB()

	//GET requests
	r.HandleFunc("/hello", user.Hello).Methods("GET")                                  //Welcome to the helpdesk
	r.HandleFunc("/list_users", user.ListUsers(db)).Methods("GET")                     //List all the users
	r.HandleFunc("/list_issues", issue.ListIssues(db)).Methods("GET")                  //List all the issues
	r.HandleFunc("/user/{User_Id}", user.GetDetailsByID(db)).Methods("GET")            //Get User Status
	r.HandleFunc("/issue_status/{Issue_Id}", issue.ShowIssueStatus(db)).Methods("GET") //Get issue status

	//POST requests
	r.HandleFunc("/new_user", user.HandleNewUSer(db)).Methods("POST")     //Register as a new User
	r.HandleFunc("/post_issue", issue.HandleNewIssue(db)).Methods("POST") //Submit a Issue

	//PUT requests
	r.HandleFunc("/assing_customer_care", issue.AssignCustomerCare(db)).Methods("PUT") //Assing Customer Care to a Query
	r.HandleFunc("/update_issue_status", issue.UpdateIssueStatus(db)).Methods("PUT")   //Solve query and change status, description and update time

	log.Fatal(http.ListenAndServe(":1001", r))
}
