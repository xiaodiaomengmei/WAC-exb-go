package model

type Useraccount struct {
	Id          int    `json:"id"`
	Name        string `json:"user_name"`
	Password    string `json:"password"`
	Human       string `json:"human_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Department  string `json:"department"`
	Role        string `json:"role"`
}
type UserDisp struct {
	PageTotal int           `json:"pagetotal"`
	User      []Useraccount `json:"list"`
}