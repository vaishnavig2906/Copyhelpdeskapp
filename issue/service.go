package issue

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type Service interface {
	ListIssues(db *sqlx.DB)
	ShowIssueStatus(db *sqlx.DB)
	HandleNewIssue(db *sqlx.DB)
	AssignCustomerCare(db *sqlx.DB)
	UpdateIssueStatus(db *sqlx.DB)
}

//HandleNewIssue Example post data (Json raw)
// {
//     "id":"1",
//     "title":"Server Issue",
//     "description":"my server is not running",
//     "reported_by":"1",
//     "created_by":"1",
//     "belongs_to":"1"
// }
func HandleNewIssue(db *sqlx.DB) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var IssueDetails IssueRequest

		err := json.NewDecoder(req.Body).Decode(&IssueDetails)
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
	})
}

func AssignCustomerCare(db *sqlx.DB) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		query := `UPDATE public.issue
		SET status='Inprogress', resolved_by='#10', updated_at=CURRENT_TIMESTAMP
		WHERE status='New';`

		_, err := db.ExecContext(ctx, query)
		if err != nil {
			log.Fatal(err)
			return
		}

		fmt.Fprintf(res, "Status set to Inprogress for all the new requests.")
		db.Close()
	})
}

func ShowIssueStatus(db *sqlx.DB) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["Issue_Id"]

		ctx := req.Context()
		query := `Select * FROM "issue" WHERE id=$1;`
		var IssueInfo Issue

		err := db.GetContext(ctx, &IssueInfo, query, id)
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
	})
}

//UpdateIssueStatus Example put data (Json raw)
// {
// 	"id":"1",
// 	"description":"please restart the service"
// }
func UpdateIssueStatus(db *sqlx.DB) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var IssueDetails IssueRequest

		err := json.NewDecoder(req.Body).Decode(&IssueDetails)
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
	})
}

func ListIssues(db *sqlx.DB) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("Complete Issue Data")
		var IssueInfo Issue

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
	})
}
