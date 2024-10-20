package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/AungKyawPhyo1142/be-students-management-system/config"
	"github.com/AungKyawPhyo1142/be-students-management-system/helpers"
	"github.com/AungKyawPhyo1142/be-students-management-system/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// create instructor in db
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

// get all instructors
func GetAllInstructors(w http.ResponseWriter, r *http.Request) {
	var instructors []models.Instructor

	if err := config.DB.Find(&instructors).Error; err != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Could not retrieve instructors: %v", err))
		return
	}

	if len(instructors) == 0 {
		helpers.RespondWithJSON(w, http.StatusOK, models.GetAllInstructorsResponse{Data: []models.InstructorResponse{}})
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, models.Instructor.GetAllInstructorsResponse(instructors[0], instructors))

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

// assign instructor to a class
func AssignInstructorToClass(w http.ResponseWriter, r *http.Request) {

	// get the instructorID and ClassID from params
	instructorID, err := helpers.GetIDFromParams(r)
	classCode := chi.URLParam(r, "code")

	if err != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("Error parsing instructor id: %v", err.Error()))
		return
	}

	// requestBody, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Error reading request body: %v", err.Error()))
	// 	return
	// }

	instructor_class := models.InstructorClass{
		InstructorID: uint(instructorID),
		ClassID:      classCode,
	}

	// if err := json.Unmarshal(requestBody, &instructor_class); err != nil {
	// 	helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Error parsing request body: %v", err.Error()))
	// 	return
	// }

	// validation
	validator := validator.New()
	if err := validator.Struct(instructor_class); err != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("Validation Failed: %v", err.Error()))
		return
	}

	result := config.DB.Create(&instructor_class)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key values") {
			helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("instructor is already assigned to the particular class: %v", result.Error.Error()))
			return
		}
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Error assigning instructor to class: %v", result.Error.Error()))
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, models.InstructorClass.ToInstructorClassResponse(instructor_class))
}

// get all instructors from a class
func GetAllInstructorsFromClass(w http.ResponseWriter, r *http.Request) {
	classCode := chi.URLParam(r, "code")

	if classCode == "" {
		helpers.RespondWithErr(w, http.StatusBadRequest, "Empty class-code in URL Param")
		return
	}

	var data []models.InstructorClass

	if err := config.DB.Where("class_id=?", classCode).Find(&data).Error; err != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Error retrieving instructor data: %v", err.Error()))
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, models.InstructorClass.ToAllInstructorClassResponse(data[0], data))
}

// delete instructor from a class
func DeleteInstructorFromClass(w http.ResponseWriter, r *http.Request) {
	instructorID, err := helpers.GetIDFromParams(r)
	classCode := chi.URLParam(r, "code")

	if err != nil {
		helpers.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("Error parsing instructor id: %v", err.Error()))
		return
	}

	if err := config.DB.Where("instructor_id=? AND class_id=?", instructorID, classCode).Delete(&models.InstructorClass{}).Error; err != nil {
		helpers.RespondWithErr(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting data: %v", err.Error()))
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, fmt.Sprintf("Instructor-ID (%d) is removed from Class-Code (%s)", instructorID, classCode))
}
