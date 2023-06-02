package main

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"log"
	)

// An interface is created
// Current db is postgres but with an interface easy to swap out to any db

type Store interface {
	CreateAccount(*Account) error
	DeleteAccount(*Account) error
	UpdateAccount(*transferreq) error
	GetAccountByID(int) (Account,error)
	GetAccounts()([]Account,error)
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

	createstmt := `CREATE TABLE IF NOT EXISTS ACCOUNT (Id serial primary key, Firstname varchar(50), Lastname varchar(50),Email varchar(50) unique,Balance int,Dob Date,created_at timestamp);`
	_,err := pg.db.Exec(createstmt)
	if err != nil{
		return fmt.Errorf("Error in creating account table")
	}

	// On delete cascade
	createTransfer := `CREATE TABLE IF NOT EXISTS Transfer (Id serial primary key,FromAccId int, ToAccId int,
									 CONSTRAINT frm FOREIGN KEY(FromAccId) references account(id) ON DELETE CASCADE,
									 CONSTRAINT trm FOREIGN KEY(ToAccId) references account(id) ON DELETE CASCADE);`
 	_,err = pg.db.Exec(createTransfer)
	if err != nil{
		return fmt.Errorf(" Error in creating transfer table")
	}
	return nil
}



func(pg *PostgresStore)CreateAccount(acc *Account)error{
	
	q:=fmt.Sprintf(`INSERT INTO ACCOUNT VALUES(%d,'%s','%s','%s',%d,'%s','%s');`, acc.Id,acc.Firstname,acc.Lastname,acc.Email,acc.Balance,acc.Dob.Format("2006-01-02"),acc.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println("QUERY ",q)
	_,err := pg.db.Exec(q)
	if err != nil{
		fmt.Println("Error inserting values")
		fmt.Println(err)
	}
	return nil
}

func(pg *PostgresStore)DeleteAccount(acc *Account)error{
	_,err := pg.db.Exec("DELETE FROM ACCOUNT WHERE ID = $1",acc.Id)
	if err != nil{
		return err
	}
	return nil

}
func(pg *PostgresStore)UpdateAccount(treq *transferreq)error{

	var ctx context.Context
	ctx = context.Background()

	tx,err := pg.db.BeginTx(ctx)
	if err != nil{
		fmt.Println(err)
	}

	defer tx.Rollback()

	

	qacc := `SELECT id,firstname,lastname,email,dob,balance,created_at FROM account WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE;`

	//upd := `UPDATE accounts SET balance = $2 WHERE id = $1 RETURNING *;`

	var fromacc,toacc Account

	row := tx.QueryRow(qacc,treq.FromAccId)
	err = scanrow(row, &fromacc)

	if err != nil{
		fmt.Errorf("Account %d does not exist", treq.FromAccId)
		return err
	}

	row = tx.QueryRow(qacc,treq.ToAccId)

	err = scanrow(row, &toacc)

	if err != nil{
		fmt.Errorf("Account %d does not exist", treq.ToAccId)
		return err
	}
	


	if fromacc.Balance - treq.Amount < 0{
		
	    return fmt.Errorf("No sufficient balance in Account %d ",treq.FromAccId, "for transfer")
	}

	

	insert :=fmt.Sprintf(`INSERT INTO Transfer(fromaccid,toaccid,amount) VALUES(%d,%d,%d);`,treq.FromAccId,treq.ToAccId,treq.Amount)

	transaction := `UPDATE account SET balance = $2 WHERE id = $1 RETURNING *;`

	_,err = tx.ExecContext(ctx,insert)

	if err != nil{
		return err
	}


	_,err = tx.ExecContext(ctx, transaction, fromacc.Id,fromacc.Balance - treq.Amount)
	
	
	if err != nil{
		if rollbackErr := tx.Rollback(); rollbackErr!=nil{
			log.Fatalf("update failed: %v, unable to rollback: %v\n", err, rollbackErr)
		}
		return err
	}


	_,err = tx.ExecContext(ctx, transaction, toacc.Id,toacc.Balance + treq.Amount)
	
	if err != nil{
		if rollbackErr := tx.Rollback(); rollbackErr!=nil{
			log.Fatalf("update failed: %v, unable to rollback: %v\n", err, rollbackErr)
		}
		fmt.Println(err)
		return err
	}


	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	return nil

}

func(pg *PostgresStore)GetAccountByID(id int)(Account,error) {

	row := pg.db.QueryRow("Select id,firstname,lastname,email,dob,balance,created_at from account where id = $1",id)
	var acc Account
	err:= row.Scan(&acc.Id,&acc.Firstname,&acc.Lastname,&acc.Email,&acc.Dob,&acc.Balance,&acc.CreatedAt)

	switch err {
		case sql.ErrNoRows:
  			return acc,fmt.Errorf("No rows")

  		case nil:
  			return acc,nil

  		default:
  			fmt.Println(err)
	}

	return acc,nil

}

func(pg *PostgresStore)GetAccounts()(acc []Account,err error) {
	rows,err := pg.db.Query("Select id,firstname,lastname,email,dob,balance,created_at from account")

	if err != nil{
		return nil,err
	}

	defer rows.Close()

	var accounts []Account

	for ;rows.Next();{
		var acc Account
		if err := rows.Scan(&acc.Id,&acc.Firstname,&acc.Lastname,&acc.Email,&acc.Dob,&acc.Balance,&acc.CreatedAt); err!=nil{
			return accounts,err
		}
		accounts = append(accounts,acc)
	}

	if err = rows.Err(); err != nil {
        return accounts, err
    }

	return accounts,err


}


func scanrow(row *sql.Row ,acc *Account)error{
	fmt.Println("scanrow")
	err := row.Scan(&acc.Id,&acc.Firstname,&acc.Lastname,&acc.Email,&acc.Dob,&acc.Balance,&acc.CreatedAt)
	switch err {
		case sql.ErrNoRows:
  			fmt.Println("No rows were returned!")
  			return err
  		case nil:
  			return nil
  		default:
  			fmt.Println("Default")
  			panic(err)
	}

	return nil
	
}

