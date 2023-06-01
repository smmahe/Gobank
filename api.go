package main

import ("fmt"
		"net/http"
		"github.com/gorilla/mux"
		"encoding/json"
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
	http.ListenAndServe(sr.listenaddr, router)

}

func (sr *Apiserver) handleAccount(w http.ResponseWriter, r *http.Request)error{
	if r.Method == "GET"{
		return sr.handleGetAccount(w, r)
	}else if r.Method == "POST"{
		return sr.handleCreateAccount(w,r)
	}else if r.Method == "DELETE"{
		return sr.handleDeleteAccount(w, r)
	}
	return nil
}

func (sr *Apiserver) handleGetAccount(w http.ResponseWriter, r *http.Request)error{
	return nil
	
}

func (sr *Apiserver) handleCreateAccount(w http.ResponseWriter, r *http.Request)error{
	jstemp := createaccreq{}
	if err:= json.NewDecoder(r.Body).Decode(&jstemp); err != nil{
		return err
	}

	acc := NewAccount(jstemp.Firstname,jstemp.Lastname,jstemp.Dob)

	err := sr.store.CreateAccount(acc)

	if err != nil{
		return err
	}

	return writeJson(w,http.StatusOK,acc)
}

func (sr *Apiserver) handleDeleteAccount(w http.ResponseWriter, r *http.Request)error{
	return nil
}
