package models

import u "go-backend/utils"

//Transfer a struct to rep Certificate Transfer
type Transfer struct {
	To     string `json:"to,omitempty"`
	Status string `json:"status,omitempty"`
	Code   string `json:"-"`
}

//Validate This struct function validate the required parameters sent through the http request body
//returns message and true if the requirement is met
func (transfer *Transfer) Validate() (map[string]interface{}, bool) {

	if transfer.To == "" {
		return u.Message(false, "Transfer To should be on the payload"), false
	}
	//All the required parameters are present
	return u.Message(true, "success"), true
}
