package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"errors"
	"time"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"os"
	"log"
	"flag"
	//"github.com/shadowshot-x/micro-product-go/authservice"
)
var (
	Log      *log.Logger
)

var errorBootingServer string = "Error Booting the Server"

type user struct {
	email string
	username string
	password string
	Fullname string
	createDate string
	role int
}

var userList []user

func TestHandler(rw http.ResponseWriter, r *http.Request) {
	var emailMissing string = "Email Missing"
	if _, ok := r.Header["Email"]; !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(emailMissing))
		return
	}
	var myString string = "Testing Done."
	fmt.Println(myString)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(myString))

}

// adds the user to the database of users
func SignupHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header)
	// extra error handling should be done at server side to prevent malicious attacks
	var emailMissing string = "Email Missing"
	fmt.Println(r.Header["Email"])
	if _, ok := r.Header["Email"]; !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(emailMissing))
		return
	}
	var userNameMissing string = "Username Missing"
	fmt.Println(r.Header["Username"])
	if _, ok := r.Header["Username"]; !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(userNameMissing))
		return
	}
	var passwordMissing string = "Password Missing"
	fmt.Println(r.Header["Password"])
	if _, ok := r.Header["Password"]; !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(passwordMissing))
		return
	}
	var fullNameMissing string = "FullName Missing"
	fmt.Println(r.Header["Fullname"])
	if _, ok := r.Header["Fullname"]; !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(fullNameMissing))
		return
	}

	// validate and then add the user
	var emailUserNameAlreadyExist string = "Email or Username already exists"
	check := AddUserObject(r.Header["Email"][0], r.Header["Username"][0], r.Header["Password"][0],
		r.Header["Fullname"][0], 0)
	// if false means username already exists
	if !check {
		rw.WriteHeader(http.StatusConflict)
		rw.Write([]byte(emailUserNameAlreadyExist))
		return
	}
	var userCreatedSuccess string = "User Created Successfully"
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(userCreatedSuccess))
}

func AddUserObject(email string, username string, password string, Fullname string, role int) bool {
	// declare the new user object
	newUser := user{
		email:     email,
		password: password,
		username:  username,
		Fullname:  Fullname,
		role:      role,
	}
	// check if a user already exists
	for _, ele := range userList {
		if ele.email == email || ele.username == username {
			return false
		}
	}
	userList = append(userList, newUser)
	return true
}

func SigninHandler(rw http.ResponseWriter, r *http.Request) {
	// validate the request first.
	var emailMissing string = "Email Missing"
	if _, ok := r.Header["Email"]; !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(emailMissing))
		return
	}
	var passwordMissing string = "Password Missing"
	if _, ok := r.Header["Password"]; !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(passwordMissing))
		return
	}
	// let’s see if the user exists
	var userNotExist string = "User Does not Exist"
	valid, err := validateUser(r.Header["Email"][0], r.Header["Password"][0])
	if err != nil {
		// this means either the user does not exist
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte(userNotExist))
		return
	}
	var passwordIncorrect string = "Incorrect Password"
	if !valid {
		// this means the password is wrong
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte(passwordIncorrect))
		return
	}

	var internalServerError string = "Internal Server Error"
	tokenString, err := getSignedToken()
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(internalServerError))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(tokenString))
}

// searches the user in the database.
func validateUser(email string, passwordHash string) (bool, error) {
	usr, exists := GetUserObject(email)
	fmt.Println("start")
	fmt.Println(usr)
	fmt.Println("end")
	var userNotExist string = "User Does not Exist"
	if !exists {
		return false, errors.New(userNotExist)
	}
	passwordCheck := usr.ValidatePasswordHash(passwordHash)

	if !passwordCheck {
		return false, nil
	}
	return true, nil
}

func GetUserObject(email string) (user, bool) {
	//needs to be replaces using Database
	for _, user := range userList {
		if user.email == email {
			return user, true
		}
	}
	return user{}, false
}
// checks if the password hash is valid
func (u *user) ValidatePasswordHash(pswdhash string) bool {
	return u.password == pswdhash
}

func getSignedToken() (string, error) {
	// we make a JWT Token here with signing method of ES256 and claims.
	// claims are attributes.
	// aud - audience
	// iss - issuer
	// exp - expiration of the Token
	claimsMap := map[string]string{
		"aud": "frontend.knowsearch.ml",
		"iss": "knowsearch.ml",
		"exp": fmt.Sprint(time.Now().Add(time.Minute * 1).Unix()),
	}
	// here we provide the shared secret. It should be very complex.
	// Also, it should be passed as a System Environment variable

	secret := "Secure_Random_String"
	header := "HS256"
	tokenString, err := GenerateToken(header, claimsMap, secret)
	if err != nil {
		return tokenString, err
	}
	return tokenString, nil
}
func GenerateToken(header string, payload map[string]string, secret string) (string, error) {
	// create a new hash of type sha256. We pass the secret key to it
	h := hmac.New(sha256.New, []byte(secret))
	header64 := base64.StdEncoding.EncodeToString([]byte(header))
	// We then Marshal the payload which is a map. This converts it to a string of JSON.
	payloadstr, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error generating Token")
		return string(payloadstr), err
	}
	payload64 := base64.StdEncoding.EncodeToString(payloadstr)

	// Now add the encoded string.
	message := header64 + "." + payload64

	// We have the unsigned message ready.
	unsignedStr := header + string(payloadstr)

	// We write this to the SHA256 to hash it.
	h.Write([]byte(unsignedStr))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	//Finally we have the token
	tokenStr := message + "." + signature
	return tokenStr, nil
}

func tokenValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// check if token is present
		if _, ok := r.Header["Token"]; !ok {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Token Missing"))
			return
		}
		token := r.Header["Token"][0]
		check, err := ValidateToken(token, "Secure_Random_String")

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Token Validation Failed"))
			return
		}
		if !check {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Token Invalid"))
			return
		}
		rw.WriteHeader(http.StatusOK)
		//rw.Write([]byte("Authorized Token"))
		next.ServeHTTP(rw, r)

	})
}

// This helps in validating the token
func ValidateToken(token string, secret string) (bool, error) {
	// JWT has 3 parts separated by '.'
	splitToken := strings.Split(token, ".")
	// if length is not 3, we know that the token is corrupt
	if len(splitToken) != 3 {
		return false, nil
	}

	// decode the header and payload back to strings
	header, err := base64.StdEncoding.DecodeString(splitToken[0])
	if err != nil {
		return false, err
	}
	payload, err := base64.StdEncoding.DecodeString(splitToken[1])
	if err != nil {
		return false, err
	}
	//again create the signature
	unsignedStr := string(header) + string(payload)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(unsignedStr))

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	fmt.Println(signature)

	// if both the signature don’t match, this means token is wrong
	if signature != splitToken[2] {
		return false, nil
	}
	// This means the token matches
	return true, nil
}

func main() {

	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mydir)

	// set location of log file
	var logpath =  "G:/xampp/htdocs/Go/GoPracticeREST/GoAuthMicroService/logger/info.log"

	flag.Parse()
	var file, err1 = os.Create(logpath)

	if err1 != nil {
		panic(err1)
	}
	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
	Log.Println("LogFile : " + logpath)

	mainRouter := mux.NewRouter()
	authRouter := mainRouter.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/signup", SignupHandler)
	authRouter.HandleFunc("/signin", SigninHandler)

	authCheckRouter := mainRouter.PathPrefix("/v1").Subrouter()
	authCheckRouter.HandleFunc("/test", TestHandler)
	authCheckRouter.Use(tokenValidationMiddleware)

	// Add the middleware to different subrouter
	// HTTP server
	// Add time outs
	server := &http.Server{
		Addr:    "127.0.0.1:9090",
		Handler: mainRouter,
	}
	err = server.ListenAndServe()
	fmt.Println("Server start",server.Addr)

	if err != nil {
		fmt.Println(errorBootingServer)
	}
}