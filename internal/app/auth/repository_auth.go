package auth

import "context"

type AuthRepository interface {
	Login(ctx context.Context, a *AuthEntities)
}
