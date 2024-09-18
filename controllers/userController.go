package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"net/http"

	"github.com/AungKyawPhyo1142/be-students-management-system/config"
	"github.com/AungKyawPhyo1142/be-students-management-system/helpers"
	"github.com/AungKyawPhyo1142/be-students-management-system/models"
)

// create user in db
func CreateUser(w http.ResponseWriter, r *http.Request) {

	// read/print out the data from request body
	// Read the request body
	body, Berr := io.ReadAll(r.Body)
	if Berr != nil {
		http.Error(w, Berr.Error(), http.StatusInternalServerError)
		return
	}

	// Print the request body
	fmt.Println("Request Body:", string(body))

	var user models.User
	err := json.Unmarshal(body, &user)

	if err != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	log.Printf("User: %+v", user)

	result := config.DB.Create(&user)

	if result.Error != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("error creating user: %v", result.Error))
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, models.User.ToUserResponse(user))

}
