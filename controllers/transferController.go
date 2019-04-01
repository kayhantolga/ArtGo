package controllers

import (
	"encoding/json"
	"go-backend/models"
	u "go-backend/utils"
	"net/http"

	"github.com/gorilla/mux"
)

//TransferCertificate This endpoint will create a new transfer request.s
var TransferCertificate = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(string)
	params := mux.Vars(r)
	id := params["Id"]
	if id == "" {
		u.Respond(w, u.Message(false, "Certificate Id should be on the payload "))
		return
	}
	transfer := &models.Transfer{}
	err := json.NewDecoder(r.Body).Decode(transfer)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	certificate := models.GetUserCertificate(user, id)
	if certificate == nil {
		u.Respond(w, u.Message(false, "We couldn't find any certificate belong to you with given ID"))
		return
	}

	data := certificate.TransferTo(transfer.To)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

//AcceptTransferCertificate This endpoint is used to complete a transfer for an object specified by Id by the authentication code.
var AcceptTransferCertificate = func(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["Id"]
	code := params["Code"]

	certificate := models.GetCertificate(code)
	if certificate == nil {
		u.Respond(w, u.Message(false, "We couldn't find any transfer invitation"))
		return
	}
	if id != certificate.ID {
		u.Respond(w, u.Message(false, "Invalid code"))
		return
	}

	data := certificate.TransferAccept(code)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
