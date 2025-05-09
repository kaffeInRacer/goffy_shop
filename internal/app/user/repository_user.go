package user

import "context"

type UserRepository interface {
	Create(ctx context.Context, d *UserEntities) (*UserEntities, error)
	ShowAll(ctx context.Context) ([]*UserEntities, error)
	Get(ctx context.Context, slug string) (*UserEntities, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, d *UserEntities) (*UserEntities, error)
}
