package main

import ("fmt"
		"net/http"
		"github.com/gorilla/mux"
		"encoding/json"
		"strconv"
		)

type Apiserver struct{
	listenaddr string
	store Store
}

func NewApiServer(addr string,store Store)*Apiserver{
	fmt.Println("creating new api server")
	return &Apiserver{listenaddr : addr,store: store,}
}



func writeJson(w http.ResponseWriter,status int,v any) error{
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter,*http.Request) error

type ApiError struct{
	Error string
}


func makeHandlerFunc(f apiFunc)http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		if err := f(w,r); err!=nil{
			writeJson(w,http.StatusBadRequest,ApiError{Error : err.Error()})
			// to handle error
		}
	}

}


func (sr *Apiserver)Run(){
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHandlerFunc(sr.handleAccount))
	router.HandleFunc("/account/{id}", makeHandlerFunc(sr.handleGetAccountById))
	router.HandleFunc("/transfer", makeHandlerFunc(sr.handleTransfer))
	router.HandleFunc("/listaccounts",sr.handleGetAccounts).Methods("GET")
	http.ListenAndServe(sr.listenaddr, router)

}

func (sr *Apiserver) handleAccount(w http.ResponseWriter, r *http.Request)error{
	if r.Method == "GET"{
		return sr.handleGetAccountById(w, r)
	}else if r.Method == "POST"{
		return sr.handleCreateAccount(w,r)
	}else if r.Method == "DELETE"{
		return sr.handleDeleteAccount(w, r)
	}
	return nil
}

func (sr *Apiserver) handleGetAccountById(w http.ResponseWriter, r *http.Request)error{
	if r.Method == "GET"{
		idstr := mux.Vars(r)["id"]
		id,err := strconv.Atoi(idstr)
		if err != nil{
			return fmt.Errorf("Invalid id given %s", idstr)
		}

		acc,err:= sr.store.GetAccountByID(id)
		if err != nil{
			return err
		}

		return writeJson(w, http.StatusOK, acc)
	}else if r.Method == "DELETE"{
		return sr.handleDeleteAccount(w, r)
	}else{
		return fmt.Errorf("Method not allowed %s",r.Method)
	}
}


func (sr *Apiserver) handleGetAccounts(w http.ResponseWriter, r *http.Request){
	accounts, err:=sr.store.GetAccounts()
	if err != nil{
		fmt.Println(err)

	}
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "application/json")
		jsoned,err:=json.Marshal(accounts)
		if err != nil{
			fmt.Println(err)
			fmt.Println(string(jsoned))
		}
		w.Write(jsoned)
	}


func (sr *Apiserver) handleCreateAccount(w http.ResponseWriter, r *http.Request)error{
	jstemp := createaccreq{}
	if err:= json.NewDecoder(r.Body).Decode(&jstemp); err != nil{
		return err
	}

	acc := NewAccount(jstemp.Firstname,jstemp.Lastname,jstemp.Email,jstemp.Dob)

	err := sr.store.CreateAccount(acc)

	if err != nil{
		return err
	}

	return writeJson(w,http.StatusOK,acc)
}

func (sr *Apiserver) handleDeleteAccount(w http.ResponseWriter, r *http.Request)error{
	
	idstr := mux.Vars(r)["id"]
	id,err := strconv.Atoi(idstr)
	if err != nil{
		return fmt.Errorf("Invalid id given %s", id)
	}

	acc,err := sr.store.GetAccountByID(id)
	if err != nil{
		return err
	}
	sr.store.DeleteAccount(&acc)
	return writeJson(w, http.StatusOK, map[string]int{"deleted":id})
}

func(sr *Apiserver) handleTransfer(w http.ResponseWriter, r *http.Request) error{
	if r.Method == "POST"{
		treq := transferreq{}
		if err:= json.NewDecoder(r.Body).Decode(&treq); err != nil{
			return err
		}

		err := sr.store.UpdateAccount(&treq)
		if err != nil{
			return err
		}

		var acc []Account

		// acc1,err := sr.store.GetAccountByID(treq.FromAccId)
		// acc2,err := sr.store.GetAccountByID(treq.ToAccId)

		// acc = append(acc,acc1)
		// acc = append(acc,acc2)
		
		return writeJson(w, http.StatusOK,acc)

	}else{
		return fmt.Errorf("Invalid Method %s",r.Method)
	}

}


// func writeTransferJson(w http.ResponseWriter,status int,v any) error{
// 	w.WriteHeader(status)
// 	w.Header().Add("Content-Type", "application/json")
// 	return json.NewEncoder(w).Encode(v)
// }