package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/RajatPrak/students/internal/storage"
	"github.com/RajatPrak/students/internal/types"
	"github.com/RajatPrak/students/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a student")
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		}

		// request validation
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)

		slog.Info("user created successfully with", slog.String("userId", fmt.Sprint(lastId)))
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Getting a student", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		student, err := storage.GetStudentById(intId)

		if err != nil {
			slog.Error("error getting user", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, student)

	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Getting all students")

		students, err := storage.GetStudents()

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}
		response.WriteJson(w, http.StatusOK, students)
	}
}

func UpdateStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get id from URL
		id := r.PathValue("id")
		intId, err := strconv.ParseInt(id, 10, 64)

		slog.Info("Student Updated with", slog.String("id", id))

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Decode request body
		var student types.Student

		err = json.NewDecoder(r.Body).Decode(&student)

		// Empty body
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest,
				response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		// Invalid JSON
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest,
				response.GeneralError(err))
			return
		}

		// Validate request
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)

			response.WriteJson(w,
				http.StatusBadRequest,
				response.ValidationError(validateErrs))
			return
		}

		// Update student
		rowsAffected, err := storage.UpdateStudent(
			student.Name,
			student.Email,
			student.Age,
			intId,
		)

		if err != nil {
			response.WriteJson(w,
				http.StatusInternalServerError,
				response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]any{
			"message": "Student updated successfully",
			"rows":    rowsAffected,
		})
	}
}

func PartiallyUpdateStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.PathValue("id")

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		slog.Info("Partially Updating Student with", slog.String("id", id))

		// Get existing student from database
		student, err := storage.GetStudentById(intId)
		if err != nil {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
			return
		}

		// Decode request body
		var req types.UpdateStudentRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest,
				response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest,
				response.GeneralError(err))
			return
		}

		// Update only the fields provided
		if req.Name != nil {
			student.Name = *req.Name
		}

		if req.Email != nil {
			student.Email = *req.Email
		}

		if req.Age != nil {
			student.Age = *req.Age
		}

		// Save updated student
		rowsAffected, err := storage.UpdateStudent(
			student.Name,
			student.Email,
			student.Age,
			student.Id,
		)

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError,
				response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]any{
			"message": "Student updated successfully",
			"rows":    rowsAffected,
		})
	}
}

func DeleteStudentById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get ID from URL
		id := r.PathValue("id")

		slog.Info("Deleting student with", slog.String("id", id))

		// Convert string to int64
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w,
				http.StatusBadRequest,
				response.GeneralError(err))
			return
		}

		// Delete student
		rowsAffected, err := storage.DeleteStudentById(intId)
		if err != nil {
			slog.Error("Failed to delete student",
				slog.String("id", id),
				slog.String("error", err.Error()))

			response.WriteJson(w,
				http.StatusInternalServerError,
				response.GeneralError(err))
			return
		}

		// Success response
		response.WriteJson(w, http.StatusOK, map[string]any{
			"message": "Student deleted successfully",
			"rows":    rowsAffected,
		})
	}
}

// UseLess function made for just testing
func PromoteStudentById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		slog.Info("api hit successfully", slog.String("id", id))
		storage.PromoteStudentById(5)
		response.WriteJson(w, http.StatusOK, map[string]any{
			"message": "Testing api",
			"goal":    "Become a good developer",
		})
	}
}
