package application

import "context"

type PermissionResolver interface {
	ResolveEffectivePermissions(
		ctx context.Context,
		userID string,
		tenantID string,
	) ([]string, error)
}

type Authorizer struct {
	resolver PermissionResolver
}

func NewAuthorizer(resolver PermissionResolver) *Authorizer {
	return &Authorizer{resolver: resolver}
}

func (a *Authorizer) Can(
	ctx context.Context,
	userID string,
	tenantID string,
	permission string,
) (bool, error) {
	permissions, err := a.resolver.ResolveEffectivePermissions(ctx, userID, tenantID)
	if err != nil {
		return false, err
	}
	for _, code := range permissions {
		if code == permission {
			return true, nil
		}
	}
	return false, nil
}
