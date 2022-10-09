package main

import (
	"encoding/json"
	"fmt"
	"log"
	"middleware/models"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var userGlobal models.User

var mySigningKey = "my-secret"

func main() {
	http.HandleFunc("/sign-in", signIn)
	http.HandleFunc("/log-in", logIn)
	http.HandleFunc("/bar", handleMiddleware(bar))
	http.HandleFunc("/foo", handleMiddleware(foo))

	log.Println(mySigningKey)

	log.Println("Server start in port 8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		authorization := request.Header.Get("Authorization")
		if authorization == "" {
			writer.WriteHeader(http.StatusExpectationFailed)
			return
		}
		token, err := jwt.Parse(authorization, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(mySigningKey), nil
		})

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			f(writer, request)
		} else {
			log.Println(err)
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
}

func logIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		panic(err)
	}

	passwordMatch := CheckPasswordHash(user.Password, userGlobal.Password)

	if !passwordMatch {
		w.WriteHeader(http.StatusUnauthorized)
		errorResponse := map[string]string{
			"messageError": "User or password incorrect",
		}
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(1 * time.Minute).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(mySigningKey))

	response := map[string]string{
		"userEmail":    userGlobal.Email,
		"access_token": tokenString,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func signIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		panic(err)
	}

	if user == (models.User{}) || user.Email == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := map[string]string{
			"status":       "Failed",
			"messageError": "Require body not present",
		}
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	hash, err := HashPassword(user.Password)
	user.Password = hash

	userGlobal = user

	response := make(map[string]string)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response["status"] = "Failed"
	} else {
		w.WriteHeader(http.StatusCreated)
		response["status"] = "Successful"
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func foo(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode("foo"); err != nil {
		log.Println("Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func bar(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode("bar"); err != nil {
		log.Println("Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
