package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/vaishnavi2906/helpdesk/issue"
	"github.com/vaishnavi2906/helpdesk/user"
)

type HelpDesk interface {
	ListUsers(res http.ResponseWriter, req *http.Request)
	ListIssues(res http.ResponseWriter, req *http.Request)
	GetDetailsByID(res http.ResponseWriter, req *http.Request)
	ShowIssueStatus(res http.ResponseWriter, req *http.Request)
	HandleNewUSer(res http.ResponseWriter, req *http.Request)
	HandleNewIssue(res http.ResponseWriter, req *http.Request)
	Hello(res http.ResponseWriter, req *http.Request)
	AssignCustomerCare(res http.ResponseWriter, req *http.Request)
	UpdateIssueStatus(res http.ResponseWriter, req *http.Request)
}

const (
	host     = "localhost"
	port     = 5432
	user_    = "postgres"
	password = "postgres"
	dbname   = "dpay_helpdesk"
)

func init_DB() (*sqlx.DB, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user_, password, host, port, dbname)
	db, err := sqlx.Open("postgres", dbURL)
	return db, err
}

//Handle new Issue
// {
//     "id":"1",
//     "title":"Server Issue",
//     "description":"my server is not running",
//     "reported_by":"1",
//     "created_by":"1",
//     "belongs_to":"1"
// }
func HandleNewIssue(res http.ResponseWriter, req *http.Request) {
	var IssueDetails issue.IssueRequest

	err := json.NewDecoder(req.Body).Decode(&IssueDetails)
	if err != nil {
		log.Fatal(err)
		return
	}

	db, err := init_DB()
	if err != nil {
		log.Fatal(err)
		return
	}

	ctx := req.Context()
	query := `INSERT INTO public.issue(
		id, title, description, reported_by, resolved_by, status ,resolved_at, created_by, created_at, updated_at, belongs_to)
		VALUES ($1, $2, $3, $4, 'Not assinged', DEFAULT, CURRENT_TIMESTAMP, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $6);`

	_, err = db.ExecContext(ctx, query, IssueDetails.Id, IssueDetails.Title, IssueDetails.Description, IssueDetails.ReportedBy, IssueDetails.CreatedBy, IssueDetails.BelongsTo)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Fprintf(res, "Succesfully Submitted the issue")
	db.Close()
}

func AssignCustomerCare(res http.ResponseWriter, req *http.Request) {
	db, err := init_DB()
	if err != nil {
		log.Fatal(err)
		return
	}

	ctx := req.Context()
	query := `UPDATE public.issue
	SET status='Inprogress', resolved_by='#10', updated_at=CURRENT_TIMESTAMP
	WHERE status='New';`

	_, err = db.ExecContext(ctx, query)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Fprintf(res, "Status set to Inprogress for all the new requests.")
	db.Close()
}

func ListUsers(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Complete User Data")
	var UserInfo user.User

	db, err := init_DB()
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
	var UserDetails user.UserRequest

	err := json.NewDecoder(req.Body).Decode(&UserDetails)
	if err != nil {
		return
	}

	db, err := init_DB()
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

	db, err := init_DB()
	if err != nil {
		log.Fatal(err)
		return
	}
	ctx := req.Context()
	query := `Select * FROM "user" WHERE id=$1;`
	var UserInfo user.User

	err = db.GetContext(ctx, &UserInfo, query, id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Given the details to the user")
	fmt.Fprintf(res, "email: %v\n", UserInfo.Email)
	db.Close()
}

func ShowIssueStatus(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["Issue_Id"]

	db, err := init_DB()
	if err != nil {
		log.Fatal(err)
		return
	}
	ctx := req.Context()
	query := `Select * FROM "issue" WHERE id=$1;`
	var IssueInfo issue.Issue

	err = db.GetContext(ctx, &IssueInfo, query, id)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Fprintf(res, "Issue_id: %v\n", IssueInfo.Id)
	fmt.Fprintf(res, "title: %v\n", IssueInfo.Title)
	fmt.Fprintf(res, "Description: %v\n", IssueInfo.Description)
	fmt.Fprintf(res, "User_id: %v\n", IssueInfo.Belongs_to)
	fmt.Fprintf(res, "Issue Submitted: %v\n", IssueInfo.Created_at)
	fmt.Fprintf(res, "Status: %v\n", IssueInfo.Status)
	fmt.Fprintf(res, "Updated At: %v\n", IssueInfo.Updated_at)
	db.Close()
}

//UpdateIssueStatus
// {
// 	"id":"1",
// 	"description":"please restart the service"
// }
func UpdateIssueStatus(res http.ResponseWriter, req *http.Request) {
	var IssueDetails issue.IssueRequest

	err := json.NewDecoder(req.Body).Decode(&IssueDetails)
	if err != nil {
		log.Fatal(err)
		return
	}

	db, err := init_DB()
	if err != nil {
		log.Fatal(err)
		return
	}

	ctx := req.Context()
	query := `UPDATE public.issue
	SET description=$2, status='Closed', resolved_at=CURRENT_TIMESTAMP, updated_at=CURRENT_TIMESTAMP
	WHERE id=$1;`

	_, err = db.ExecContext(ctx, query, IssueDetails.Id, IssueDetails.Description)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Fprintf(res, "Resolved Description: %v\n", IssueDetails.Description)
	db.Close()
}

func ListIssues(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Complete Issue Data")
	var IssueInfo issue.Issue

	db, err := init_DB()
	if err != nil {
		log.Fatal(err)
		return
	}

	query := `SELECT * from "issue";`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
		return
	}

	for rows.Next() {
		err := rows.Scan(&IssueInfo.Id, &IssueInfo.Title, &IssueInfo.Description, &IssueInfo.Reported_by, &IssueInfo.Resolved_by, &IssueInfo.Status, &IssueInfo.Resolved_at, &IssueInfo.Created_by, &IssueInfo.Created_at, &IssueInfo.Updated_at, &IssueInfo.Belongs_to)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println("\n", IssueInfo.Id, IssueInfo.Title, IssueInfo.Description, IssueInfo.Reported_by, IssueInfo.Resolved_by, IssueInfo.Status, IssueInfo.Resolved_at, IssueInfo.Created_by, IssueInfo.Created_at, IssueInfo.Updated_at, IssueInfo.Belongs_to)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return
	}
	db.Close()
}

func Hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Hello Welcome to helpdesk")
}

func main() {
	init_router()
}
