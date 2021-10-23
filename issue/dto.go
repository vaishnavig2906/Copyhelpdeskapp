package issue

import "time"

type Issue struct {
	Id          string    `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Reported_by string    `db:"reported_by"`
	Resolved_by string    `db:"resolved_by"`
	Status      string    `db:"status"`
	Resolved_at time.Time `db:"resolved_at"`
	Created_by  string    `db:"created_by"`
	Created_at  time.Time `db:"created_at"`
	Updated_at  time.Time `db:"updated_at"`
	Belongs_to  string    `db:"belongs_to"`
}

type IssueRequest struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ReportedBy  string    `json:"reported_by"`
	ResolvedBy  string    `json:"resolved_by"`
	Status      string    `json:"status"`
	ResolvedAt  time.Time `json:"resolved_at"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	BelongsTo   string    `json:"belongs_to"`
}
