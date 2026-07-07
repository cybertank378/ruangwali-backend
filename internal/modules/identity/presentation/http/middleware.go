package httpidentity

import (
	"net/http"
	"strings"

	acapp "github.com/ruangwali/internal/modules/accesscontrol/application"
	acdomain "github.com/ruangwali/internal/modules/accesscontrol/domain"
	"github.com/ruangwali/internal/modules/identity/infrastructure/security"
	"github.com/ruangwali/internal/shared/application/requestcontext"
)

func Authenticate(tokens *security.TokenService, authorizer *acapp.Authorizer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := tokens.Parse(r.Context(), strings.TrimPrefix(header, "Bearer "))
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			// Resolve once per request. Production dapat ditambah Redis cache.
			permissions, err := resolvePermissions(r, authorizer, claims)
			if err != nil {
				http.Error(w, "authorization resolution failed", http.StatusInternalServerError)
				return
			}

			principal := requestcontext.Principal{
				UserID:       claims.UserID,
				TenantID:     claims.TenantID,
				MembershipID: claims.MembershipID,
				Permissions:  permissions,
			}

			next.ServeHTTP(w, r.WithContext(
				requestcontext.WithPrincipal(r.Context(), principal),
			))
		})
	}
}

func resolvePermissions(r *http.Request, authorizer *acapp.Authorizer, claims security.Claims) (map[string]struct{}, error) {
	// Temporary adapter: permission catalog dicek satu per satu.
	// Ganti dengan Authorizer.ResolveAll pada iterasi berikutnya.
	catalog := []acdomain.PermissionCode{
		acdomain.DashboardRead,
		acdomain.StudentRead, acdomain.StudentCreate, acdomain.StudentUpdate, acdomain.StudentDelete,
		acdomain.AttendanceRead, acdomain.AttendanceRecord, acdomain.AttendanceUpdate,
		acdomain.AssessmentRead, acdomain.AssessmentRecord, acdomain.AssessmentUpdate,
		acdomain.AchievementRead, acdomain.AchievementCreate,
		acdomain.ViolationRead, acdomain.ViolationCreate,
		acdomain.JournalRead, acdomain.JournalCreate,
		acdomain.SchoolRead, acdomain.SchoolUpdate,
		acdomain.UserRead, acdomain.UserCreate,
		acdomain.RoleRead, acdomain.RoleManage, acdomain.RoleAssign,
		acdomain.AuditRead, acdomain.BackupCreate, acdomain.BackupRestore, acdomain.SystemReset,
	}

	out := make(map[string]struct{})
	for _, permission := range catalog {
		ok, err := authorizer.Can(r.Context(), claims.UserID, claims.TenantID, string(permission))
		if err != nil {
			return nil, err
		}
		if ok {
			out[string(permission)] = struct{}{}
		}
	}
	return out, nil
}

func RequirePermission(code acdomain.PermissionCode) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			principal, ok := requestcontext.PrincipalFrom(r.Context())
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			if !principal.HasPermission(string(code)) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
