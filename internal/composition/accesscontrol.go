// =========================================================
// File: internal/composition/accesscontrol.go.go
// =========================================================
package composition

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	accesscontrolservice "github.com/ruangwali/internal/modules/accesscontrol/application/service"
	accesscontrolpostgres "github.com/ruangwali/internal/modules/accesscontrol/infrastructure/postgres"
	accesscontrolhttp "github.com/ruangwali/internal/modules/accesscontrol/presentation/http"
	"github.com/ruangwali/internal/shared/application/requestcontext"
)

type AccessControlModule struct {
	AuthorizationService *accesscontrolservice.AuthorizationService
	Middleware           *accesscontrolhttp.Middleware
}

type requestContextUserIDResolver struct{}

func (requestContextUserIDResolver) ResolveUserID(
	ctx context.Context,
) (uuid.UUID, bool) {
	return requestcontext.UserID(
		ctx,
	)
}

func buildAccessControl(
	db *pgxpool.Pool,
) (*AccessControlModule, error) {
	if db == nil {
		return nil, errors.New(
			"build access control: db nil",
		)
	}

	permissionResolver :=
		accesscontrolpostgres.NewPermissionResolver(
			db,
		)

	authorizationService :=
		accesscontrolservice.NewAuthorizationService(
			permissionResolver,
		)

	userIDResolver :=
		requestContextUserIDResolver{}

	middleware :=
		accesscontrolhttp.NewMiddleware(
			authorizationService,
			userIDResolver,
		)

	return &AccessControlModule{
		AuthorizationService: authorizationService,
		Middleware:           middleware,
	}, nil
}
