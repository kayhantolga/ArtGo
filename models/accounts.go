package models

import (
	"context"
	u "go-backend/utils"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

//Token a struct to rep user token
type Token struct {
	UserID string
	jwt.StandardClaims
}

//Account a struct to rep user account
type Account struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
	ID       string `json:"to,omitempty"`
	Name     string `json:"name,omitempty"`
}

//JwtAuthentication Standart JWT token service
var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/api/user/login"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path              //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			response = u.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			response = u.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in
		tk := &Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			response = u.Message(false, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			response = u.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		//	fmt.Sprintf("User %", tk.UserID) //Useful for monitoring
		ctx := context.WithValue(r.Context(), "user", tk.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}

//Validate incoming user details
func (account *Account) Validate() (map[string]interface{}, bool) {

	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	//Email must be unique
	temp := &Account{}

	//check for errors and duplicate emails
	// err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error

	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

//Login Login with static user DB
func Login(email, password string) map[string]interface{} {

	user := GetAccount(email)
	if user == nil {
		return u.Message(false, "Invalid login credentials. Please try again")
	}

	if password != "123" {
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	account := &Account{
		Email: user.Email,
	}

	//Create JWT token
	tk := &Token{UserID: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString //Store the token in the response
	account.ID = "1"
	account.CreatedAt = time.Now()

	resp := u.Message(true, "Logged In")
	resp["account"] = account
	return resp
}

//GetAccount Get user account via email
func GetAccount(email string) *Account {
	email = strings.ToLower(email)
	if email == "kayhantolga@hotmail.com" {
		return &Account{
			Email: "kayhantolga@hotmail.com",
			ID:    "1",
			Name:  "Tolga",
		}
	}
	if email == "kayhantolga@gmail.com" {
		return &Account{
			Email: "kayhantolga@gmail.com",
			ID:    "2",
			Name:  "Sedat",
		}
	}

	return nil
}
