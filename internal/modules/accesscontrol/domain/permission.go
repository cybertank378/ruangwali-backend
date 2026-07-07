package domain

type PermissionCode string

const (
	DashboardRead PermissionCode = "dashboard.read"

	StudentRead   PermissionCode = "student.read"
	StudentCreate PermissionCode = "student.create"
	StudentUpdate PermissionCode = "student.update"
	StudentDelete PermissionCode = "student.delete"
	StudentImport PermissionCode = "student.import"
	StudentExport PermissionCode = "student.export"

	GuardianRead   PermissionCode = "guardian.read"
	GuardianCreate PermissionCode = "guardian.create"
	GuardianUpdate PermissionCode = "guardian.update"
	GuardianDelete PermissionCode = "guardian.delete"

	ClassroomRead               PermissionCode = "classroom.read"
	ClassroomUpdate             PermissionCode = "classroom.update"
	ClassroomOrganizationRead   PermissionCode = "classroom.organization.read"
	ClassroomOrganizationManage PermissionCode = "classroom.organization.manage"
	ClassroomScheduleRead       PermissionCode = "classroom.schedule.read"
	ClassroomScheduleManage     PermissionCode = "classroom.schedule.manage"
	ClassroomDutyRead           PermissionCode = "classroom.duty.read"
	ClassroomDutyManage         PermissionCode = "classroom.duty.manage"
	ClassroomSeatingRead        PermissionCode = "classroom.seating.read"
	ClassroomSeatingManage      PermissionCode = "classroom.seating.manage"
	ClassroomAgreementRead      PermissionCode = "classroom.agreement.read"
	ClassroomAgreementManage    PermissionCode = "classroom.agreement.manage"

	AttendanceRead   PermissionCode = "attendance.read"
	AttendanceRecord PermissionCode = "attendance.record"
	AttendanceUpdate PermissionCode = "attendance.update"
	AttendanceDelete PermissionCode = "attendance.delete"
	AttendanceExport PermissionCode = "attendance.export"

	AssessmentRead   PermissionCode = "assessment.read"
	AssessmentRecord PermissionCode = "assessment.record"
	AssessmentUpdate PermissionCode = "assessment.update"
	AssessmentDelete PermissionCode = "assessment.delete"
	AssessmentExport PermissionCode = "assessment.export"

	AchievementRead   PermissionCode = "achievement.read"
	AchievementCreate PermissionCode = "achievement.create"
	AchievementUpdate PermissionCode = "achievement.update"
	AchievementDelete PermissionCode = "achievement.delete"

	ViolationRead   PermissionCode = "violation.read"
	ViolationCreate PermissionCode = "violation.create"
	ViolationUpdate PermissionCode = "violation.update"
	ViolationDelete PermissionCode = "violation.delete"

	JournalRead   PermissionCode = "journal.read"
	JournalCreate PermissionCode = "journal.create"
	JournalUpdate PermissionCode = "journal.update"
	JournalDelete PermissionCode = "journal.delete"

	SchoolRead   PermissionCode = "school.read"
	SchoolUpdate PermissionCode = "school.update"

	UserRead       PermissionCode = "user.read"
	UserCreate     PermissionCode = "user.create"
	UserUpdate     PermissionCode = "user.update"
	UserDeactivate PermissionCode = "user.deactivate"

	RoleRead   PermissionCode = "role.read"
	RoleManage PermissionCode = "role.manage"
	RoleAssign PermissionCode = "role.assign"

	AuditRead    PermissionCode = "audit.read"
	BackupCreate PermissionCode = "backup.create"
	BackupRestore PermissionCode = "backup.restore"
	SystemReset  PermissionCode = "system.reset"
)
