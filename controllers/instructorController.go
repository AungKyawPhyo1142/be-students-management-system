package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/AungKyawPhyo1142/be-students-management-system/config"
	"github.com/AungKyawPhyo1142/be-students-management-system/helpers"
	"github.com/AungKyawPhyo1142/be-students-management-system/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func CreateInstructor(w http.ResponseWriter, r *http.Request) {
	requestBody, Berror := io.ReadAll(r.Body)
	if Berror != nil {
		http.Error(w, Berror.Error(), http.StatusInternalServerError)
		return
	}

	var instructor models.Instructor
	err := json.Unmarshal(requestBody, &instructor)
	if err != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	//validation
	validator := validator.New()
	Verr := validator.Struct(instructor)

	if Verr != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("Validation failed: %v", Verr.Error()))
		return
	}
	result := config.DB.Create(&instructor)

	if result.Error != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Error creating instructor: %v", result.Error.Error()))
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, models.Instructor.ToInstructorResponse(instructor))
}

// update instructor in db
func EditInstructor(w http.ResponseWriter, r *http.Request) {

	//getting the id from url params
	idParams := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParams)
	if err != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("Error parsing instructor id: %v", err.Error()))
		return
	}

	// get o from db by id
	var instrcutor models.Instructor
	result := config.DB.First(&instrcutor, id)

	if result.Error != nil {
		helpers.RespondWithErr(w, http.StatusNotFound, "Instructor not found!")
		return
	}

	requestBody, rError := io.ReadAll(r.Body)

	if rError != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("error reading request: %v", rError.Error()))
		return
	}

	error := json.Unmarshal(requestBody, &instrcutor)

	if error != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("invalid request: %v", error.Error()))
		return
	}

	// save the instructor's updated data
	if err := config.DB.Save(&instrcutor).Error; err != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update instructor's data: %v", err.Error()))
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, models.Instructor.ToInstructorResponse(instrcutor))
}

// get instructor by id
func GetInstructorByID(w http.ResponseWriter, r *http.Request) {
	idParams := chi.URLParam(r, "id")
	id, error := strconv.Atoi(idParams)
	if error != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("Error parsing student id: %v", error.Error()))
		return
	}

	var instructor models.Instructor
	if err := config.DB.First(&instructor, id).Error; err != nil {
		helpers.RespondWithErr(w, http.StatusNotFound, "Instructor not found!")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, models.Instructor.ToInstructorResponse(instructor))
}

// delete instructor
func DeleteInstructor(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.GetIDFromParams(r)
	if err != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("Error parsing id: %v", err.Error()))
		return
	}

	// find the id in the student table
	if err := config.DB.Delete(&models.Instructor{}, id).Error; err != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting instructor: %v", err.Error()))
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, fmt.Sprintf("Instructor id %v is deleted successfully!", id))
}

