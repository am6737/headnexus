package repository

import (
	"context"
	"github.com/am6737/headnexus/domain/user/entity"
)

type UserRepository interface {
	Create(ctx context.Context, network *entity.User) (*entity.User, error)
	Get(ctx context.Context, id string) (*entity.User, error)
	//Find(ctx context.Context, query *entity.QueryNetwork) ([]*entity.Network, error)
	Update(ctx context.Context, network *entity.User) (*entity.User, error)
	Delete(ctx context.Context, id string) error
}
