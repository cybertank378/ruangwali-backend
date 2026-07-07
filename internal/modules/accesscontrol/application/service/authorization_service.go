package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/ruangwali/internal/modules/accesscontrol/application/ports"
	accessdomain "github.com/ruangwali/internal/modules/accesscontrol/domain"
)

type AuthorizationService struct {
	permissionResolver ports.PermissionResolver
}

func NewAuthorizationService(
	permissionResolver ports.PermissionResolver,
) *AuthorizationService {
	if permissionResolver == nil {
		panic(
			"authorization service: permission resolver nil",
		)
	}

	return &AuthorizationService{
		permissionResolver: permissionResolver,
	}
}

func (s *AuthorizationService) Authorize(
	ctx context.Context,
	userID uuid.UUID,
	permission string,
) error {
	if userID == uuid.Nil {
		return accessdomain.ErrAccessDenied
	}

	permission = normalizePermission(
		permission,
	)
	if permission == "" {
		return accessdomain.ErrAccessDenied
	}

	permissions, err := s.resolvePermissions(
		ctx,
		userID,
	)
	if err != nil {
		return err
	}

	if _, exists := permissions[permission]; !exists {
		return accessdomain.ErrAccessDenied
	}

	return nil
}

func (s *AuthorizationService) AuthorizeAny(
	ctx context.Context,
	userID uuid.UUID,
	requiredPermissions ...string,
) error {
	if userID == uuid.Nil {
		return accessdomain.ErrAccessDenied
	}

	required := normalizePermissions(
		requiredPermissions,
	)
	if len(required) == 0 {
		return accessdomain.ErrAccessDenied
	}

	permissions, err := s.resolvePermissions(
		ctx,
		userID,
	)
	if err != nil {
		return err
	}

	for _, permission := range required {
		if _, exists := permissions[permission]; exists {
			return nil
		}
	}

	return accessdomain.ErrAccessDenied
}

func (s *AuthorizationService) AuthorizeAll(
	ctx context.Context,
	userID uuid.UUID,
	requiredPermissions ...string,
) error {
	if userID == uuid.Nil {
		return accessdomain.ErrAccessDenied
	}

	required := normalizePermissions(
		requiredPermissions,
	)
	if len(required) == 0 {
		return accessdomain.ErrAccessDenied
	}

	permissions, err := s.resolvePermissions(
		ctx,
		userID,
	)
	if err != nil {
		return err
	}

	for _, permission := range required {
		if _, exists := permissions[permission]; !exists {
			return accessdomain.ErrAccessDenied
		}
	}

	return nil
}

func (s *AuthorizationService) resolvePermissions(
	ctx context.Context,
	userID uuid.UUID,
) (map[string]struct{}, error) {
	resolved, err := s.permissionResolver.ResolveByUserID(
		ctx,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal menyelesaikan permission user: %w",
			err,
		)
	}

	permissions := make(
		map[string]struct{},
		len(resolved),
	)

	for _, permission := range resolved {
		normalized := normalizePermission(
			permission,
		)
		if normalized == "" {
			continue
		}

		permissions[normalized] = struct{}{}
	}

	return permissions, nil
}

func normalizePermissions(
	permissions []string,
) []string {
	normalized := make(
		[]string,
		0,
		len(permissions),
	)

	seen := make(
		map[string]struct{},
		len(permissions),
	)

	for _, permission := range permissions {
		value := normalizePermission(
			permission,
		)
		if value == "" {
			continue
		}

		if _, exists := seen[value]; exists {
			continue
		}

		seen[value] = struct{}{}

		normalized = append(
			normalized,
			value,
		)
	}

	return normalized
}

func normalizePermission(
	permission string,
) string {
	return strings.ToLower(
		strings.TrimSpace(permission),
	)
}
