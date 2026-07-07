package domain

type RoleCode string

const (
	RoleSuperAdmin      RoleCode = "SUPER_ADMIN"
	RoleSchoolAdmin     RoleCode = "SCHOOL_ADMIN"
	RoleHomeroomTeacher RoleCode = "HOMEROOM_TEACHER"
	RoleSubjectTeacher  RoleCode = "SUBJECT_TEACHER"
)

type Role struct {
	ID          string
	TenantID    *string
	Code        RoleCode
	Name        string
	Description string
	IsSystem    bool
	IsActive    bool
}
