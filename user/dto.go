package user

type User struct {
	Id       string `db:"id"`
	Email    string `db:"email"`
	Usertype string `db:"user_type"`
}

type UserRequest struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Usertype string `json:"user_type"`
}
