package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	if len(students) == 0 {
		helpers.RespondWithJSON(w, http.StatusOK, models.GetAllStudentsResponse{Data: []models.StudentResponse{}})
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

// assign student to class
func AssignStudentToClass(w http.ResponseWriter, r *http.Request) {
	studentID, err := helpers.GetIDFromParams(r)
	classCode := chi.URLParam(r, "code")
	if err != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("Error parsing student id: %v", err.Error()))
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Error reading request body: %v", err.Error()))
		return
	}

	// we can't unmarshal primitive data-types
	type EnrollmentDateRequest struct {
		EnrollmentDate string `json:"enrollment_date"`
	}

	var enrollmentDate EnrollmentDateRequest

	if err := json.Unmarshal(reqBody, &enrollmentDate); err != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Error parsing request body: %v", err.Error()))
		return
	}

	// define the layout of time
	const layout = "02-01-2006"

	// parse the string to time
	parsedDate, err := time.Parse(layout, enrollmentDate.EnrollmentDate)
	if err != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Error parsing enrollment date: %v", err.Error()))
		return
	}

	studentClass := models.StudentClass{
		StudentID:      uint(studentID),
		ClassID:        classCode,
		EnrollmentDate: parsedDate,
	}

	if err := config.DB.Create(&studentClass).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key values") {
			helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("student is already assigned to the particular class: %v", err.Error()))
			return
		}
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Error saving data: %v", err.Error()))
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, studentClass)

}

// get all students from a class
func GetAllStudentsFromClass(w http.ResponseWriter, r *http.Request) {
	classCode := chi.URLParam(r, "code")

	if classCode == "" {
		helpers.RespondWithErr(w, http.StatusBadRequest, "Empty class-code in URL Param")
		return
	}

	var data []models.StudentClass

	if err := config.DB.Where("class_id=?", classCode).Find(&data).Error; err != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Error retrieving student data: %v", err.Error()))
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, models.StudentClass.ToAllStudentClassResponse(data[0], data))

}

// delete student from class
func DeleteStudentFromClass(w http.ResponseWriter, r *http.Request) {
	studentID, err := helpers.GetIDFromParams(r)
	classCode := chi.URLParam(r, "code")
	if err != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("Error parsing student id: %v", err.Error()))
		return
	}

	if err := config.DB.Where("student_id=? AND class_id=?", studentID, classCode).Delete(&models.StudentClass{}).Error; err != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting data: %v", err.Error()))
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, fmt.Sprintf("Student-ID (%d) is removed from Class-Code (%s)", studentID, classCode))
}
