package login

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Akmyrat17/carm/models"
	"github.com/dgrijalva/jwt-go"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	valid := (credentials.Password == "admin" && credentials.Username == "admin")

	if !valid {
		http.Error(w, "Incorrect username or password", http.StatusUnauthorized)
		return
	}

	tokenString, err := GenerateToken(credentials.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		log.Println("error Generating token: ", err)
		return
	}

	response := map[string]string{"token": tokenString}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GenerateToken(username string) (string, error) {
	expiration := time.Now().Add(24 * time.Hour)
	claims := &jwt.StandardClaims{
		ExpiresAt: expiration.Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   username,
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	signedToken, err := token.SignedString([]byte("some_valeu"))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
