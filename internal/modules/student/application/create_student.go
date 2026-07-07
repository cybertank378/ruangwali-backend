package application

import (
	"context"

	"github.com/google/uuid"
	studentdomain "github.com/ruangwali/internal/modules/student/domain"
)

type CreateStudentInput struct {
	TenantID string
	FullName string
	Gender   string
}

type CreateStudentOutput struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Gender   string `json:"gender"`
}

type CreateStudentUseCase struct {
	repository studentdomain.Repository
}

func NewCreateStudentUseCase(repository studentdomain.Repository) *CreateStudentUseCase {
	return &CreateStudentUseCase{repository: repository}
}

func (uc *CreateStudentUseCase) Execute(ctx context.Context, input CreateStudentInput) (CreateStudentOutput, error) {
	tenantID, err := uuid.Parse(input.TenantID)
	if err != nil {
		return CreateStudentOutput{}, err
	}

	student, err := studentdomain.NewStudent(
		tenantID,
		input.FullName,
		studentdomain.Gender(input.Gender),
	)
	if err != nil {
		return CreateStudentOutput{}, err
	}

	if err := uc.repository.Create(ctx, student); err != nil {
		return CreateStudentOutput{}, err
	}

	return CreateStudentOutput{
		ID: student.ID.String(), FullName: student.FullName, Gender: string(student.Gender),
	}, nil
}
