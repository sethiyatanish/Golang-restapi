package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sethiyatanish/student-api/internal/types"
	"github.com/sethiyatanish/student-api/internal/utils/response"
)

func New() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			// responseBody := map[string]string{"error": err.Error()}
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}



		slog.Info("creating a student")

			if err := validator.New().Struct(student); err != nil {

				validateErrs := err.(validator.ValidationErrors)
				response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
				return 
			}




		responseBody := map[string]string{"success": "OK"}
		response.WriteJson(w, http.StatusAccepted, responseBody)


	}
}