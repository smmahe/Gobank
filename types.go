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
	Password string		`json:"-"`
}

type loginreq struct{
	Email string	`json:"email"`
	Password string  `json:"password"`
}

type createaccreq struct{
	Firstname string
	Lastname string
	Email string
	Dob string
	Password string
}

type transferreq struct{
	FromAccId int `json:"fromaccount"`
	ToAccId	int 	`json:"toaccount"`
	Amount int     `json:"amount"`
}


func NewAccount(firstname string,lastname string,email string,dob string,password string) *Account{
	
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
		Password: password,
	}
}
