package controllers

import (
	"encoding/json"
	"go-backend/models"
	u "go-backend/utils"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// CreateCertificate This endpoint will create a Verisart certificate using the provided request body fields.
var CreateCertificate = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(string)
	certificate := &models.Certificate{}

	err := json.NewDecoder(r.Body).Decode(certificate)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	certificate.OwnerID = user
	certificate.CreatedAt = time.Now()
	resp := certificate.Create()
	u.Respond(w, resp)
}

//UpdateCertificate This endpoint will update an existing Verisart certificate using the provided request body fields.
var UpdateCertificate = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(string)
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		u.Respond(w, u.Message(false, "Certificate Id should be on the payload "))
		return
	}

	certificate := &models.Certificate{}

	err := json.NewDecoder(r.Body).Decode(certificate)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}
	if user != certificate.OwnerID {
		u.Respond(w, u.Message(false, "Changing owner is not allowed"))
		return
	}
	certificate.ID = id

	resp := certificate.Update()
	u.Respond(w, resp)
}

//DeleteCertificate This endpoint will delete an existing Verisart certificate that is identified by the id in the URL
var DeleteCertificate = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(string)
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		u.Respond(w, u.Message(false, "Certificate Id should be on the payload "))
		return
	}

	resp := models.Delete(id, user)
	u.Respond(w, resp)
}

//GetCertificatesFor This endpoint is used to retrieve a specific certificate, specified by object_id
var GetCertificatesFor = func(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["userId"]
	if id == "" {
		id = r.Context().Value("user").(string)
	}

	data := models.GetUserCertificates(id)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
