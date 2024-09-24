package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/AungKyawPhyo1142/be-students-management-system/config"
	"github.com/AungKyawPhyo1142/be-students-management-system/helpers"
	"github.com/AungKyawPhyo1142/be-students-management-system/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// create a class in db
func CreateClass(w http.ResponseWriter, r *http.Request) {
	reqBody, err := io.ReadAll(r.Body)

	if err != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err.Error()))
		return
	}

	var class models.Class
	if err := json.Unmarshal(reqBody, &class); err != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Error parsing request body: %v", err.Error()))
		return
	}

	// validation
	validator := validator.New()
	if err := validator.Struct(class); err != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("Validation Error: %v", err.Error()))
		return
	}

	if err := config.DB.Create(&class).Error; err != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Error parsing request body: %v", err))
		return
	}
	helpers.RespondWithJSON(w, http.StatusCreated, class.ToClassResponse())

}

// update class
func UpdateClass(w http.ResponseWriter, r *http.Request) {
	classCode := chi.URLParam(r, "code")

	// get class from db by class_code
	var class models.Class
	if err := config.DB.Where("class_code=?", classCode).First(&class).Error; err != nil {
		helpers.RespondWithErr(w, http.StatusNotFound, fmt.Sprintf("Class not found: %v", err.Error()))
		return
	}

	reqBody, err := io.ReadAll(r.Body)

	if err != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("Error reading request body: %v", err.Error()))
		return
	}

	if err := json.Unmarshal(reqBody, &class); err != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("Error parsing request body: %v", err.Error()))
		return
	}

	if err := config.DB.Save(&class).Error; err != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Error updating course data: %v", err.Error()))
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, models.Class.ToClassResponse(class))

}

// delete class
func DeleteClass(w http.ResponseWriter, r *http.Request) {
	classCode := chi.URLParam(r, "code")

	if err := config.DB.Where("class_code=?", classCode).Delete(&models.Class{}).Error; err != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting class code %s: %v", classCode, err.Error()))
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, fmt.Sprintf("%s deleted successfully!", classCode))
}

// get class by classCode
func GetClassByID(w http.ResponseWriter, r *http.Request) {
	// getting the class code
	classCode := chi.URLParam(r, "code")

	var class models.Class

	// search class by classCode
	if err := config.DB.Where("class_code=?", classCode).First(&class).Error; err != nil {
		helpers.RespondWithErr(w, http.StatusNotFound, fmt.Sprintf("Class not found: %v", err.Error()))
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, class.ToClassResponse())
}

// get all classes
func GetAllClasses(w http.ResponseWriter, r *http.Request) {
	var classes []models.Class

	if err := config.DB.Find(&classes).Error; err != nil {
		helpers.RespondWithErr(w, http.StatusNotFound, err.Error())
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, models.Class.GetAllClassesResponse(classes[0], classes))
}
