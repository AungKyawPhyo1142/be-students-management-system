package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AungKyawPhyo1142/be-students-management-system/config"
	"github.com/AungKyawPhyo1142/be-students-management-system/helpers"
	"github.com/AungKyawPhyo1142/be-students-management-system/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// Claims define the structure of the JWT payload
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// register user
func Register(w http.ResponseWriter, r *http.Request) {
	CreateAdmin(w, r)
}

// login user
func Login(w http.ResponseWriter, r *http.Request) {
	var admin models.Admin
	var dbAdmin models.Admin

	body, err := io.ReadAll(r.Body)

	if err != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err.Error()))
		return
	}

	if err := json.Unmarshal(body, &admin); err != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err.Error()))
		return
	}

	// get the db user with same username
	if err := config.DB.Where("username=?", admin.Username).First(&dbAdmin).Error; err != nil {
		helpers.RespondWithErr(w, http.StatusUnauthorized, fmt.Sprintf("User not found: %v", err.Error()))
		return
	}

	log.Printf("Db Admin: %v", dbAdmin)

	// compare the hashedPassword with provided password
	if err := bcrypt.CompareHashAndPassword([]byte(dbAdmin.Password), []byte(admin.Password)); err != nil {
		helpers.RespondWithErr(w, http.StatusUnauthorized, fmt.Sprintf("Invalid password: %v", err.Error()))
		return
	}

	// create jwt token
	expirationTime := time.Now().Add(1 * time.Hour) // set 1hr as expiration time
	claims := &Claims{
		Username: dbAdmin.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Error generating JWT token: %v", err.Error()))
		return
	}

	type jwtResponse struct {
		Token string `json:"token"`
	}

	helpers.RespondWithJSON(w, http.StatusOK, jwtResponse{Token: tokenString})

}
