package repository

import (
	"context"
	"github.com/am6737/headnexus/domain/host/entity"
)

type RuleRepository interface {
	Create(ctx context.Context, userID string, rule *entity.Rule) (*entity.Rule, error)
	Get(ctx context.Context, userID, id string) (*entity.Rule, error)
	Update(ctx context.Context, userID string, rule *entity.Rule) error
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, userID string, options *entity.RuleFindOptions) ([]*entity.Rule, error)
}
