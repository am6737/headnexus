package repository

import (
	"context"
	"github.com/am6737/headnexus/domain/rule/entity"
)

type RuleRepository interface {
	Create(ctx context.Context, rule *entity.Rule) (*entity.Rule, error)
	Get(ctx context.Context, id string) (*entity.Rule, error)
	Update(ctx context.Context, rule *entity.Rule) error
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, options *entity.FindOptions) ([]*entity.Rule, error)
}
