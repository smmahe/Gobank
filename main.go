package main
import ("log"
		"fmt")
func main(){
	store,err := NewPostgresStore()
	if err != nil{
		fmt.Println(err)
		log.Fatal(" No store")
	}
	fmt.Printf("%+v", store)
	store.INIT()
	srv := NewApiServer(":8080",store)
	srv.Run()
}

// To run go run .