package entity

type Useraccount struct {
	User_id  string `json:"user_id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Status   string `json:"status"`
	Auditable
	Verification string `json:"verification_code"`
}
