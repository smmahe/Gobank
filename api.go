package main

import ("fmt"
		"net/http"
		"github.com/gorilla/mux"
		"encoding/json"
		"strconv"
		jwt "github.com/golang-jwt/jwt"
		"os"
		"golang.org/x/crypto/bcrypt"
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
	router.HandleFunc("/login", makeHandlerFunc(sr.handleLogin))
	router.HandleFunc("/account", makeHandlerFunc(sr.handleAccount))
	router.HandleFunc("/account/{id}", withJWTAuth(makeHandlerFunc(sr.handleGetAccountById)))
	router.HandleFunc("/transfer", makeHandlerFunc(sr.handleTransfer))
	router.HandleFunc("/listaccounts",sr.handleGetAccounts).Methods("GET")
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
	if r.Method == "GET"{
		tokenstring := r.Header.Get("jwt-token")
		fmt.Println(tokenstring)
		token,err := validateJWT(tokenstring)
		if err != nil{
			return writeJson(w,http.StatusForbidden,ApiError{Error:"Invalid Token "})
			 
		}
		
		claims, _ := token.Claims.(jwt.MapClaims)


		var id int
		if token.Valid {
			id,_ = strconv.Atoi(fmt.Sprint(claims["accountnumber"]))
		}
        		

		acc,err:= sr.store.GetAccountByID(int(id))
		
		if err != nil{
			return err
		}

			return writeJson(w, http.StatusOK, acc)
		}
		return writeJson(w, http.StatusForbidden, "Error")
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
		w.WriteHeader(http.StatusOK)
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

	acc := NewAccount(jstemp.Firstname,jstemp.Lastname,jstemp.Email,jstemp.Dob,jstemp.Password)

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

		
		return writeJson(w, http.StatusOK,acc)

	}else{
		return fmt.Errorf("Invalid Method %s",r.Method)
	}
}

// JWT Middleware

func withJWTAuth(httphandlerfunc http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		tokenstring := r.Header.Get("jwt-token")
		token,err := validateJWT(tokenstring)
		if err != nil{
			writeJson(w,http.StatusForbidden,ApiError{Error:"Invalid Token "})
			return 
		}
		idstr := mux.Vars(r)["id"]
		id,err := strconv.Atoi(idstr)
		claims, _ := token.Claims.(jwt.MapClaims)

        
        if token.Valid && claims["accountnumber"] == float64(id){
        	httphandlerfunc(w,r)
    	}else{
			writeJson(w,http.StatusForbidden,ApiError{Error:"Account number cannot be accessed"})
			return
		}
		
		
	}
}

func validateJWT(tokenstring string)(*jwt.Token,error){
	secret:=os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
		return []byte(secret), nil
	})

	return token,err

}

func createJWT(acc *Account)(string,error){
	claims := &jwt.MapClaims{
		"expiresAt":15000,
		"accountnumber":acc.Id,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret)) 	

}	

func (sr *Apiserver) handleLogin(w http.ResponseWriter, r *http.Request)error{
	var req loginreq
	if r.Method != "POST"{
		return writeJson(w, http.StatusUnauthorized, "Invalid Request Method")
	}
	if err:= json.NewDecoder(r.Body).Decode(&req); err != nil{
			return err
	}

	acc,err := sr.store.GetAccountForLogin(req.Email)

	if err != nil{
		return fmt.Errorf("Invalid Email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(acc.Password),[]byte(req.Password))

	if err != nil{
		return fmt.Errorf("Invalid Password")
	}

	token,err := createJWT(&acc)
	w.Header().Add("bearer-token",token)

	return writeJson(w, http.StatusOK,"Login Success!!")
}


