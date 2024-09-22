package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/AungKyawPhyo1142/be-students-management-system/config"
	"github.com/AungKyawPhyo1142/be-students-management-system/helpers"
	"github.com/AungKyawPhyo1142/be-students-management-system/models"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// create user in db
func CreateAdmin(w http.ResponseWriter, r *http.Request) {

	// read/print out the data from request body
	// Read the request body
	body, Berr := io.ReadAll(r.Body)
	if Berr != nil {
		http.Error(w, Berr.Error(), http.StatusInternalServerError)
		return
	}

	var user models.Admin
	err := json.Unmarshal(body, &user)

	if err != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// validation
	validator := validator.New()
	if err := validator.Struct(user); err != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("Validation Failed: %v", err.Error()))
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Error hashing password: %v", err.Error()))
		return
	}

	// update the password with hashed password
	user.Password = string(hashedPassword)

	result := config.DB.Create(&user)

	if result.Error != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("error creating user: %v", result.Error))
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, models.Admin.ToAdminResponse(user))

}

// get all users from db
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.Admin
	result := config.DB.Find(&users)
	if result.Error != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("error fetching users: %v", result.Error))
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, models.Admin.GetAllUsersResponse(users[0], users))
}
