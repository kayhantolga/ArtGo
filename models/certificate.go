package models

import (
	u "go-backend/utils"
	"time"

	uuid "github.com/satori/go.uuid"
)

//Certificate a struct to rep user account
type Certificate struct {
	ID        string    `json:"id,omitempty"`
	Title     string    `json:"title,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	OwnerID   string    `json:"ownerId,omitempty"`
	Year      int       `json:"year,omitempty"`
	Note      string    `json:"note,omitempty"`
	Transfer  *Transfer `json:"transfer,omitempty"`
}

var certificates []Certificate

//Validate This struct function validate the required parameters sent through the http request body
//returns message and true if the requirement is met
func (certificate *Certificate) Validate() (map[string]interface{}, bool) {

	if certificate.Title == "" {
		return u.Message(false, "Phone number should be on the payload"), false
	}

	//All the required parameters are present
	return u.Message(true, "success"), true
}

// Create will create a Verisart certificate using the provided fields.
func (certificate *Certificate) Create() map[string]interface{} {

	if resp, ok := certificate.Validate(); !ok {
		return resp
	}
	generatedUIId, err := uuid.NewV4()
	if err != nil {
		return u.Message(false, "Something went wrong.")
	}
	certificate.ID = generatedUIId.String()
	certificates = append(certificates, *certificate)

	resp := u.Message(true, "success")
	resp["certificate"] = certificate
	return resp
}

//Update will update an existing Verisart certificate using the provided  fields.
func (certificate *Certificate) Update() map[string]interface{} {

	if resp, ok := certificate.Validate(); !ok {
		return resp
	}
	oldcertIndex := GetIndex(certificate.ID)
	if oldcertIndex == -1 {
		return u.Message(false, "We couldn't find any certificate belong to you with given ID")
	}
	if certificates[oldcertIndex].OwnerID != certificate.OwnerID {
		return u.Message(false, "We couldn't find any certificate belong to you with given ID")
	}

	certificate.CreatedAt = certificates[oldcertIndex].CreatedAt
	certificates[oldcertIndex] = *certificate

	resp := u.Message(true, "success")
	resp["certificate"] = certificate
	return resp
}

//Delete will delete an existing Verisart certificate that is identified by the ID
func Delete(id string, userID string) map[string]interface{} {

	oldcertIndex := GetIndex(id)
	if oldcertIndex == -1 {
		return u.Message(false, "We couldn't find any certificate belong to you with given ID")
	}
	if certificates[oldcertIndex].OwnerID != userID {
		return u.Message(false, "We couldn't find any certificate belong to you with given ID")
	}
	certificates = append(certificates[:oldcertIndex], certificates[oldcertIndex+1:]...)
	resp := u.Message(true, "success")
	return resp
}

//TransferTo will create a new transfer request.
func (certificate *Certificate) TransferTo(userEmail string) map[string]interface{} {

	oldcertIndex := GetIndex(certificate.ID)
	if oldcertIndex == -1 {
		return u.Message(false, "We couldn't find any certificate belong to you with given ID")
	}

	generatedUUID, err := uuid.NewV4()
	if err != nil {
		return u.Message(false, "Something went wrong.")
	}

	certificate.Transfer = &Transfer{
		Code:   generatedUUID.String(),
		Status: "pending",
		To:     userEmail,
	}
	user := GetAccount(certificate.Transfer.To)
	certificates[oldcertIndex] = *certificate
	//TODO get url from env
	u.SendTransformInvitation(user.Name, "http://localhost:8000/certificates/"+certificate.ID+"/transfers/"+certificate.Transfer.Code, certificate.Transfer.To)
	resp := u.Message(true, "success")
	resp["certificate"] = certificate
	return resp
}

//TransferAccept will complete a transfer
func (certificate *Certificate) TransferAccept(code string) map[string]interface{} {

	oldcertIndex := GetIndexFromTransferCode(code)
	if oldcertIndex == -1 {
		return u.Message(false, "We couldn't find any transfer invitation")
	}
	if certificate.Transfer == nil {
		return u.Message(false, "We couldn't find any transfer invitation")
	}

	account := GetAccount(certificates[oldcertIndex].Transfer.To)
	certificates[oldcertIndex].OwnerID = account.ID
	certificates[oldcertIndex].Transfer.Status = "done"

	resp := u.Message(true, "success")
	resp["certificate"] = certificates[oldcertIndex]
	return resp
}

//GetUserCertificates Retrieve a collection of certificates that belong to the user.
func GetUserCertificates(userID string) []Certificate {
	var userCertificates []Certificate
	for _, item := range certificates {
		if item.OwnerID == userID {
			userCertificates = append(userCertificates, item)
		}
	}
	return userCertificates
}

//GetUserCertificate Retrieve a certificate that belong to the user with given ID.
func GetUserCertificate(userID string, certificateID string) *Certificate {
	for _, item := range certificates {
		if item.OwnerID == userID && item.ID == certificateID {
			return &item
		}
	}
	return nil
}

//GetCertificate Retrieve a certificate with given ID
func GetCertificate(certificateID string) *Certificate {
	for _, item := range certificates {
		if item.ID == certificateID {
			return &item
		}
	}
	return nil
}

//GetIndex Get index of certificate with given ID
func GetIndex(id string) int {
	for index, item := range certificates {
		if item.ID == id {
			return index
		}
	}
	return -1
}

//GetIndexFromTransferCode Get index of certificate with given transaction Codes
func GetIndexFromTransferCode(code string) int {
	for index, item := range certificates {
		if item.Transfer != nil && item.Transfer.Code == code && item.Transfer.Status == "pending" {
			return index
		}
	}
	return -1
}

//GetCertificateFromTransferCode Get certificate with given transaction Codes
func GetCertificateFromTransferCode(code string) int {
	for index, item := range certificates {
		if item.Transfer != nil && item.Transfer.Code == code && item.Transfer.Status == "pending" {
			return index
		}
	}
	return -1
}
