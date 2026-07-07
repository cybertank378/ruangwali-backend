package httpstudent

import (
	"encoding/json"
	"net/http"

	studentapp "github.com/ruangwali/internal/modules/student/application"
	"github.com/ruangwali/internal/shared/application/requestcontext"
)

type Handler struct {
	create *studentapp.CreateStudentUseCase
}

func NewHandler(create *studentapp.CreateStudentUseCase) *Handler {
	return &Handler{create: create}
}

type createRequest struct {
	FullName string `json:"fullName"`
	Gender   string `json:"gender"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	principal, ok := requestcontext.PrincipalFrom(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var request createRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	output, err := h.create.Execute(r.Context(), studentapp.CreateStudentInput{
		TenantID: principal.TenantID,
		FullName: request.FullName,
		Gender:   request.Gender,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"data": output})
}
