package main
import ("math/rand"
		"time")

type Account struct{
	Id int				`json:"id"`
	Firstname string 	`json:"firstname"`
	Lastname string		`json:"lastname"`
	Email string		`json:"email"`
	Balance int     	`json:"balance"`
	Dob time.Time    	`json:"date"`
	CreatedAt time.Time  `json:"createdat"`
}

type createaccreq struct{
	Firstname string
	Lastname string
	Email string
	Dob string
}

func NewAccount(firstname string,lastname string,email string,dob string) *Account{
	
	const lay = "2006-01-02"
	const DateTime = "2006-01-02 15:04:05"

	par,_ := time.Parse(lay, dob)

	
	return &Account{
		Id:rand.Intn(100000),
		Firstname: firstname,
		Lastname: lastname,
		Email: email,
		Balance: 0,
		Dob: par,
		CreatedAt :time.Now().UTC(),
	}
}