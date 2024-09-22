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

// create student in db
func CreateStudent(w http.ResponseWriter, r *http.Request) {
	requestBody, Berror := io.ReadAll(r.Body)
	if Berror != nil {
		http.Error(w, Berror.Error(), http.StatusInternalServerError)
		return
	}
	var student models.Student
	err := json.Unmarshal(requestBody, &student)
	if err != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	//validation
	validator := validator.New()
	Verr := validator.Struct(student)

	if Verr != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("Validation failed: %v", Verr.Error()))
		return
	}

	result := config.DB.Create(&student)

	if result.Error != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Error creating student: %v", result.Error.Error()))
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, models.Student.ToStudentResponse(student))
}

// update student in db
func EditStudent(w http.ResponseWriter, r *http.Request) {

	//getting the id from url params
	idParams := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParams)
	if err != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("Error parsing student id: %v", err.Error()))
		return
	}

	// get student from db by id
	var student models.Student
	result := config.DB.First(&student, id)

	if result.Error != nil {
		helpers.RespondWithErr(w, http.StatusNotFound, "Student not found!")
		return
	}

	requestBody, rError := io.ReadAll(r.Body)

	if rError != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("error reading request: %v", rError.Error()))
		return
	}

	error := json.Unmarshal(requestBody, &student)

	if error != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("invalid request: %v", error.Error()))
		return
	}

	// save the student's updated data
	if err := config.DB.Save(&student).Error; err != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update student's data: %v", err.Error()))
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, models.Student.ToStudentResponse(student))
}

// get student by id
func GetStudentByID(w http.ResponseWriter, r *http.Request) {
	idParams := chi.URLParam(r, "id")
	id, error := strconv.Atoi(idParams)
	if error != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("Error parsing student id: %v", error.Error()))
		return
	}

	var student models.Student
	if err := config.DB.First(&student, id).Error; err != nil {
		helpers.RespondWithErr(w, http.StatusNotFound, "Student not found!")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, models.Student.ToStudentResponse(student))
}

// get all students
func GetAllStudents(w http.ResponseWriter, r *http.Request) {
	var students []models.Student

	if err := config.DB.Find(&students).Error; err != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Error getting students: %v", err.Error()))
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, models.Student.GetAllStudentsResponse(students[0], students))
}

// delete student
func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.GetIDFromParams(r)
	if err != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("Error parsing id: %v", err.Error()))
		return
	}

	// find the id in the student table
	if err := config.DB.Delete(&models.Student{}, id).Error; err != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting student: %v", err.Error()))
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, fmt.Sprintf("Student id %v is deleted successfully!", id))
}
