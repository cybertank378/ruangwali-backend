package requestcontext

import "context"

type Principal struct {
	UserID       string
	TenantID     string
	MembershipID string
	Permissions  map[string]struct{}
}

func (p Principal) HasPermission(code string) bool {
	_, ok := p.Permissions[code]
	return ok
}

type key struct{}

func WithPrincipal(ctx context.Context, principal Principal) context.Context {
	return context.WithValue(ctx, key{}, principal)
}

func PrincipalFrom(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(key{}).(Principal)
	return p, ok
}
