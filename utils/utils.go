package utils

import (
	"encoding/json"
	"net/http"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// func SendEmail(body string) {
// 	from := "kayhantolga@gmail.com"
// 	pass := "dummyIbo13"
// 	to := "kayhantolga@gmail.com"

// 	msg := "From: " + from + "\n" +
// 		"To: " + to + "\n" +
// 		"Subject: Hello there\n\n" +
// 		body

// 	err := smtp.SendMail("smtp.gmail.com:587",
// 		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
// 		from, []string{to}, []byte(msg))

// 	if err != nil {
// 		log.Printf("smtp error: %s", err)
// 		return
// 	}

// 	log.Print("sent")
// }
