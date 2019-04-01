package main

import (
	"Go-Backend/app"
	"Go-Backend/controllers"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// main function to boot up everything
func main() {
	router := mux.NewRouter()

	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/certificates", controllers.CreateCertificate).Methods("POST")
	router.HandleFunc("/certificates/{id}", controllers.UpdateCertificate).Methods("PATCH")
	router.HandleFunc("/certificates/{id}", controllers.DeleteCertificate).Methods("DELETE")
	router.HandleFunc("/user/{userId}/certificates", controllers.GetCertificatesFor).Methods("GET")
	router.HandleFunc("/certificates/{Id}/transfers", controllers.TransferCertificate).Methods("POST")
	router.HandleFunc("/certificates/{Id}/transfers/{Code}", controllers.AccpetTransferCertificate).Methods("GET")

	//Cors suppport for  all urls
	corsObj := handlers.AllowedOrigins([]string{"*"})
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(corsObj)(router)))
}
