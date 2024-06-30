package handlers

import (
	"net/http"

	"github.com/AungKyawPhyo1142/be-students-management-system/helpers"
)

func HandlerReady(w http.ResponseWriter, r *http.Request) {

	type Ready struct {
		Message string `json:"message"`
	}

	helpers.RespondWithJSON(w, 200, Ready{
		Message: "Ready!!!",
	})

}
