package main

import ("database/sql"
	_ "github.com/lib/pq"
	"fmt"
	)

// An interface is created
// Current db is postgres but with an interface easy to swap out to any db

type Store interface {
	CreateAccount(*Account) error
	DeleteAccount(*Account) error
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account,error)
}


type PostgresStore struct{
	db *sql.DB
}

func NewPostgresStore()(*PostgresStore,error){
	db,err := sql.Open("postgres",connstr)
	if err != nil{
		return nil,err
	}

	if err := db.Ping();err!=nil{
		return nil,err
	}


	return &PostgresStore{
		db : db,
	},nil
}


func (pg *PostgresStore)INIT()error{

	createstmt := `CREATE TABLE IF NOT EXISTS ACCOUNT (Id serial primary key, Firstname varchar(50), Lastname varchar(50),Balance int,Dob Date,created_at timestamp);`
	_,err := pg.db.Exec(createstmt)
	if err != nil{
		fmt.Println(err)
		fmt.Println("Error is table creation")
	}
	return nil
}



func(pg *PostgresStore)CreateAccount(acc *Account)error{
	
	q:=fmt.Sprintf(`INSERT INTO ACCOUNT VALUES(%d,'%s','%s',%d,'%s','%s');`, acc.Id,acc.Firstname,acc.Lastname,acc.Balance,acc.Dob.Format("2006-01-02"),acc.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println("QUERY ",q)
	_,err := pg.db.Exec(q)
	if err != nil{
		fmt.Println("Error inserting values")
		fmt.Println(err)
	}
	return nil
}

func(pg *PostgresStore)DeleteAccount(acc *Account)error{
	return nil

}
func(pg *PostgresStore)UpdateAccount(acc *Account)error{
	return nil

}
func(pg *PostgresStore)GetAccountByID(int)(acc *Account,err error) {
	return nil,nil

}
