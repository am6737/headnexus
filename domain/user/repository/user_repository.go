package repository

import (
	"context"
	"github.com/am6737/headnexus/domain/user/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	Get(ctx context.Context, id string) (*entity.User, error)
	Find(ctx context.Context, opt *entity.FindOptions) ([]*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string) error
}
