package main

import (
	"Go-Backend/controllers"
	"go-backend/models"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// main function to boot up everything
func main() {
	router := mux.NewRouter()

	router.Use(models.JwtAuthentication) //attach JWT auth middleware

	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/certificates", controllers.CreateCertificate).Methods("POST")
	router.HandleFunc("/api/certificates/{id}", controllers.UpdateCertificate).Methods("PATCH")
	router.HandleFunc("/api/certificates/{id}", controllers.DeleteCertificate).Methods("DELETE")
	router.HandleFunc("/api/user/{userId}/certificates", controllers.GetCertificatesFor).Methods("GET")
	router.HandleFunc("/api/certificates/{Id}/transfers", controllers.TransferCertificate).Methods("POST")
	router.HandleFunc("/api/certificates/{Id}/transfers/{Code}", controllers.AcceptTransferCertificate).Methods("PATCH")

	//Cors suppport for  all urls
	corsObj := handlers.AllowedOrigins([]string{"*"})
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(corsObj)(router)))
}
